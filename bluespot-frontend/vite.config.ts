import { fileURLToPath, URL } from 'node:url'
import path from 'path'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from '@tailwindcss/vite'
import metaPlugin, { getBuildHash } from './build/meta'
import dayjs from 'dayjs'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'

// https://vite.dev/config/
export default defineConfig(({ mode, command }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const isServe = command === 'serve'
  console.log('env.VITE_APP_MODE_ENV: ', env.VITE_APP_MODE_ENV)
  console.log('isServe: ', isServe)
  const now = dayjs().format('YYYY-MM-DD HH:mm:ss')
  return {
    define: {
      // https://github.com/vitejs/vite/issues/2605#issuecomment-803276660
      LOCAL_BUILD_HASH: env.VITE_APP_MODE_ENV !== 'dev' ? JSON.stringify(getBuildHash()) : '""',
      LOCAL_BUILD_TIME: env.VITE_APP_MODE_ENV !== 'dev' ? JSON.stringify(now) : '""',
    },
    base: '/',
    server: {
      port: 5175,
      // host: '0.0.0.0',
      proxy: {
        '/bluespot/api': {
          target: 'http://localhost:6306',
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/bluespot/, ''),
        },
      },
    },
    plugins: [
      vue(),
      vueJsx(),
      vueDevTools(),
      tailwindcss(),
      metaPlugin(),
      createSvgIconsPlugin({
        // 指定需要缓存的图标文件夹
        iconDirs: [path.resolve(process.cwd(), 'src/assets/icons')],
        // 指定symbolId格式
        symbolId: 'icon-[dir]-[name]',
      }),
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    build: {
      rollupOptions: {
        output: {
          // Group JS files in js directory
          entryFileNames: `js/[name]-[hash].js`,
          chunkFileNames: `js/[name]-[hash].js`,
          // Group CSS files in css directory
          assetFileNames: (assetInfo) => {
            const extType = assetInfo.name!.split('.').pop()?.toLowerCase()
            if (extType === 'css') {
              return `css/[name]-[hash].[ext]`
            }
            // Keep other assets in assets directory
            return `assets/[name]-[hash].[ext]`
          },
        },
      },
    },
  }
})
