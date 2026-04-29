import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tsconfigPaths from 'vite-tsconfig-paths'

const useMsw = process.env['VITE_USE_MSW'] === 'true'

const authTarget = process.env['VITE_AUTH_URL'] ?? 'http://localhost:8000'
const candidateTarget = process.env['VITE_CANDIDATE_URL'] ?? 'http://localhost:8081'
const companyTarget = process.env['VITE_COMPANY_URL'] ?? 'http://localhost:8088'

export default defineConfig({
  plugins: [react(), tsconfigPaths()],
  server: {
    proxy: useMsw
      ? {}
      : {
          '/api/v1/auth': { target: authTarget, changeOrigin: true },
          '/api/v1/candidate': { target: candidateTarget, changeOrigin: true },
          '/api/v1/resume': { target: candidateTarget, changeOrigin: true },
          '/api/v1/skill': { target: candidateTarget, changeOrigin: true },
          '/api/v1/companies': { target: companyTarget, changeOrigin: true },
          '/api/v1/vacancies': { target: companyTarget, changeOrigin: true },
        },
  },
})
