import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tsconfigPaths from 'vite-tsconfig-paths'

const useMsw = process.env['VITE_USE_MSW'] === 'true'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tsconfigPaths()],
  server: {
    proxy: useMsw
      ? {}
      : {
          '/auth': {
            target: process.env['VITE_AUTH_URL'] ?? 'http://localhost:8000',
            changeOrigin: true,
            rewrite: (path) => `/api/v1${path}`,
          },
          '/api/v1/candidate': {
            target: process.env['VITE_CANDIDATE_URL'] ?? 'http://localhost:8081',
            changeOrigin: true,
          },
          '/api/v1/companies': {
            target: process.env['VITE_COMPANY_URL'] ?? 'http://localhost:8088',
            changeOrigin: true,
          },
          '/api/v1/vacancies': {
            target: process.env['VITE_COMPANY_URL'] ?? 'http://localhost:8088',
            changeOrigin: true,
          },
        },
  },
})
