/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL: string
  readonly VITE_USE_MSW?: string
  readonly VITE_APP_ENV?: 'development' | 'staging' | 'production'
  readonly VITE_AUTH_ME_AVAILABLE?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
