import { createHash } from 'crypto'
import fsExtra from 'fs-extra'
import fg from 'fast-glob'
import dayjs from 'dayjs'
import type { Plugin, ResolvedConfig } from 'vite'
import path from 'path'

const { readFileSync, writeFileSync } = fsExtra

// 全局缓存构建哈希
let cachedBuildHash: string | null = null

// 同步计算构建哈希的函数
export function getBuildHash(): string {
  if (cachedBuildHash) {
    return cachedBuildHash
  }

  try {
    const projectFiles = fg.sync(['src/**/*', 'public/**/*'], {
      dot: true,
      ignore: [],
      absolute: true  // 使用绝对路径确保稳定性
    }).sort()

    const hash = createHash('sha256', { outputLength: 32 })
    for (const file of projectFiles) {
      const content = readFileSync(file, 'utf-8')
        .replace(/\r\n/g, '\n')  // 统一换行符
        .trim()  // 去除首尾空白
      hash.update(file + content) // 包含文件名确保唯一性
    }

    cachedBuildHash = hash.digest('hex').slice(0, 8)
    return cachedBuildHash
  } catch (error) {
    console.warn('Failed to generate build hash:', error)
    // 降级方案：使用时间戳
    return dayjs().format('YYYYMMDDHHmmss')
  }
}

export default function mataPlugin(distFilePath: string = 'dist/meta.json') {
  // 通过 configResolved 获取最终 outDir 和 root
  let resolvedOutDir = 'dist'
  let projectRoot = process.cwd()
  const apply: Plugin['apply'] = (_, { command }) => {
    return command === 'build' // 仅打包环境调用
  }
  return {
    name: 'generate-meta-file',
    apply,
    configResolved(resolvedConfig: ResolvedConfig) {
      projectRoot = resolvedConfig.root
      resolvedOutDir = resolvedConfig.build?.outDir || 'dist'
    },
    async closeBundle() {
      // 获取已经计算好的构建哈希
      const buildHash = getBuildHash()

      const metaData = {
        buildHash,
        buildTime: dayjs().format('YYYY-MM-DD HH:mm:ss')
      }
      // 兼容：若传入的是相对路径，按 outDir/meta.json 写入；若用户传入了绝对路径，则直接使用
      const outputFile = path.isAbsolute(distFilePath)
        ? distFilePath
        : path.resolve(projectRoot, path.extname(distFilePath) ? distFilePath : path.join(resolvedOutDir, 'meta.json'))
      try {
        // 确保目录存在
        fsExtra.ensureDirSync(path.dirname(outputFile))
  
        writeFileSync(outputFile, JSON.stringify(metaData, null, 2))
        console.log('Meta info generated:', JSON.stringify(metaData))
      } catch (error) {
        console.error('Failed to write meta file:', error)
      }
    },
  }
}