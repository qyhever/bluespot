/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_APP_MODE_ENV: 'dev' | 'test' | 'prod'
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare const LOCAL_BUILD_HASH: string
declare const LOCAL_BUILD_TIME: string
