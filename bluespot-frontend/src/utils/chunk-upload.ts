import { uploadVerify, uploadChunk, mergeChunks } from '@/api/chunk'
import type { MergeChunksResponse } from '@/api/chunk'
import { toFixed } from './index'

const flag = true

function getInterval(start: number) {
  const end = new Date().getTime();
  const interval = (end - start) / 1000;
  return toFixed(interval);
}

interface ChunkUploaderOptions {
  chunkSize?: number
  concurrency?: number
  retryCount?: number
  onHashProgress?: (progress: number) => void
  onProgress?: (progress: number) => void
  onSuccess?: (data: MergeChunksResponse) => void
  onError?: (error: unknown) => void
}

interface ChunkItem {
  blob: Blob
  index: number
}

interface UploadChunkResponse {
  url: string
  msg: string
}

class ChunkUploader {
  private file: File | null = null
  private chunkSize: number
  private concurrency: number
  private retryCount: number

  private onHashProgress: (progress: number) => void
  private onProgress: (progress: number) => void
  private onSuccess: (data: UploadChunkResponse) => void
  private onError: (error: unknown) => void

  private fileMd5 = ''
  private uploadId = ''
  private chunkList: ChunkItem[] = []
  private uploadedChunks: number[] = []
  private isStop = false

  constructor(options: ChunkUploaderOptions = {}) {
    this.chunkSize = options.chunkSize ?? 5 * 1024 * 1024 // 5MB 最佳分片
    this.concurrency = options.concurrency ?? 4 // 并发数
    this.retryCount = options.retryCount ?? 2 // 每片重试次数

    this.onHashProgress = options.onHashProgress ?? (() => {})
    this.onProgress = options.onProgress ?? (() => {})
    this.onSuccess = options.onSuccess ?? (() => {})
    this.onError = options.onError ?? (() => {})
  }

  /**
   * 开始上传
   */
  async start(file: File): Promise<void> {
    try {
      this.isStop = false
      this.file = file

      // 1. 计算全文件 MD5
      const hashStart = new Date().getTime()
      this.fileMd5 = await this.computeFileMd5(file)
      const hashTime = +getInterval(hashStart)
      console.log('fileMd5: ', this.fileMd5, 'hashTime: ', hashTime)
      if (flag) {
        return
      }
      if (!this.fileMd5) {
        throw new Error('文件 MD5 计算失败')
      }

      // 2. 预检接口：秒传 + 获取已上传分片
      const verifyRes = await uploadVerify({
        fileMd5: this.fileMd5,
        fileName: this.file.name,
        fileSize: this.file.size
      })

      if (verifyRes.isExists) {
        this.onSuccess({
          url: verifyRes.url,
          msg: '秒传成功'
        })
        return
      }

      this.uploadId = verifyRes.uploadId
      this.uploadedChunks = verifyRes.uploadedChunks ?? []

      // 3. 切分文件
      this.chunkList = this.splitFile(file)

      // 4. 过滤已上传分片
      const waitList = this.chunkList.filter(
        item => !this.uploadedChunks.includes(item.index)
      )

      if (waitList.length === 0) {
        await this.merge()
        return
      }

      // 5. 并发上传
      await this.runPool(waitList)

      if (this.isStop) return

      // 6. 合并分片
      await this.merge()
    } catch (error) {
      this.onError(error)
    }
  }

  /**
   * 停止上传
   */
  stop(): void {
    this.isStop = true
  }

  /**
   * 增量计算文件 MD5
   */
  private computeFileMd5(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const worker = new Worker(
      new URL('./md5.worker.ts', import.meta.url),
      {
        type: 'module'
      }
    )

    worker.postMessage({
      file,
      chunkSize: this.chunkSize
    })

    worker.onmessage = event => {
      const data = event.data as
        | { type: 'progress'; progress: number }
        | { type: 'success'; md5: string }
        | { type: 'error'; message: string }

      if (data.type === 'progress') {
        console.log('正在计算文件指纹：', data.progress)
        this.onHashProgress(data.progress)
      }

      if (data.type === 'success') {
        worker.terminate()
        resolve(data.md5)
      }

      if (data.type === 'error') {
        worker.terminate()
        reject(new Error(data.message))
      }
    }

    worker.onerror = error => {
      worker.terminate()
      reject(error)
    }
  })
}

  /**
   * 文件切分
   */
  private splitFile(file: File): ChunkItem[] {
    const list: ChunkItem[] = []
    const chunks = Math.ceil(file.size / this.chunkSize)

    for (let i = 0; i < chunks; i++) {
      const start = i * this.chunkSize
      const end = Math.min(start + this.chunkSize, file.size)

      list.push({
        blob: file.slice(start, end),
        index: i
      })
    }

    return list
  }

  /**
   * 上传单个分片，自带重试
   */
  private async uploadChunk(item: ChunkItem): Promise<boolean> {
    if (!this.file) {
      throw new Error('文件不存在')
    }

    let times = 0

    while (times <= this.retryCount) {
      try {
        await uploadChunk({
          chunk: item.blob,
          uploadId: this.uploadId,
          fileMd5: this.fileMd5,
          chunkIndex: item.index,
          fileName: this.file.name
        })
        return true
      } catch (error) {
        times++

        if (times > this.retryCount) {
          throw error
        }
      }
    }

    return false
  }

  /**
   * Promise 并发池
   */
  private async runPool(tasks: ChunkItem[]): Promise<void> {
    let finished = 0
    let currentIndex = 0

    const totalAll = this.chunkList.length

    const worker = async (): Promise<void> => {
      while (currentIndex < tasks.length && !this.isStop) {
        const task = tasks[currentIndex]
        currentIndex++

        await this.uploadChunk(task!)

        finished++

        const progress = Number(
          (
            ((this.uploadedChunks.length + finished) / totalAll) *
            100
          ).toFixed(2)
        )

        this.onProgress(progress)
      }
    }

    const workers = Array.from(
      {
        length: Math.min(this.concurrency, tasks.length)
      },
      () => worker()
    )

    await Promise.all(workers)
  }

  /**
   * 合并分片
   */
  private async merge(): Promise<void> {
    const res = await mergeChunks({
      uploadId: this.uploadId,
      fileMd5: this.fileMd5,
      chunkLength: this.chunkList.length
    })

    this.onSuccess(res)
  }
}

export default ChunkUploader