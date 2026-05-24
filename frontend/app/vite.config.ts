import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tsconfigPaths from 'vite-tsconfig-paths'

const useMsw = process.env['VITE_USE_MSW'] === 'true'

const backendTarget = process.env['VITE_BACKEND_URL'] ?? 'https://api.traineematch.space'

export default defineConfig({
  plugins: [react(), tsconfigPaths()],
  server: {
    proxy: useMsw
      ? {}
      : {
          '/api/v1/admin/companies': {
            target: backendTarget,
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api\/v1\/admin\/companies/, '/api/company/admin/companies'),
            secure: true,
          },
          '/api/v1/admin/vacancies': {
            target: backendTarget,
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api\/v1\/admin\/vacancies/, '/api/company/admin/vacancies'),
            secure: true,
          },
          '/api/v1/admin/candidates': {
            target: backendTarget,
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api\/v1\/admin\/candidates/, '/api/candidate/admin/candidates'),
            secure: true,
          },
          '/api/v1/admin/skills': {
            target: backendTarget,
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api\/v1\/admin\/skills/, '/api/candidate/admin/skills'),
            secure: true,
          },
          '/api/v1/admin/resumes': {
            target: backendTarget,
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api\/v1\/admin\/resumes/, '/api/candidate/admin/resumes'),
            secure: true,
          },
          '/api/v1/admin/new': {
            target: backendTarget,
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api\/v1\/admin\/new/, '/api/auth/admin/new'),
            secure: true,
          },
          '/api/v1/auth': {
            target: backendTarget,
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api\/v1\/auth/, '/api/auth/auth'),
            secure: true,
          },
          '/api/v1/candidate': {
            target: backendTarget,
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api\/v1\/candidate/, '/api/candidate/candidate'),
            secure: true,
          },
          '/api/v1/resume': {
            target: backendTarget,
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api\/v1\/resume/, '/api/candidate/resume'),
            secure: true,
          },
          '/api/v1/skill': {
            target: backendTarget,
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api\/v1\/skill/, '/api/candidate/skill'),
            secure: true,
          },
          '/api/v1/companies': {
            target: backendTarget,
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api\/v1\/companies/, '/api/company/companies'),
            secure: true,
          },
          '/api/v1/vacancies': {
            target: backendTarget,
            changeOrigin: true,
            rewrite: (path) => path.replace(/^\/api\/v1\/vacancies/, '/api/company/vacancies'),
            secure: true,
          },
          '/api/v1/applications': {
            target: backendTarget,
            changeOrigin: true,
            rewrite: (path) => path.replace('/api/v1/applications', '/api/application/applications'),
            secure: true,
          },
          '/api/v1/hr': {
            target: backendTarget,
            changeOrigin: true,
            rewrite: (path) => path.replace('/api/v1/hr', '/api/application/hr'),
            secure: true,
          },
        },
  },
})
