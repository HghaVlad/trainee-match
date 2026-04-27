import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import tsconfigPaths from 'vite-tsconfig-paths'

export default defineConfig({
  plugins: [react(), tsconfigPaths()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    env: {
      VITE_API_URL: 'http://localhost:8080',
      VITE_USE_MSW: 'true',
      VITE_AUTH_ME_AVAILABLE: 'false',
    },
    coverage: {
      provider: 'v8',
      thresholds: {
        statements: 60,
        branches: 50,
      },
      exclude: ['src/api/generated/**', 'src/test/**'],
    },
  },
})
