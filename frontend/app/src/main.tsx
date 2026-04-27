import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import '@/app/styles/globals.css'
import { QueryProvider } from '@/app/providers/QueryProvider'
import { AppRouter } from '@/app/router/router'
import { bootstrap } from '@/shared/session/bootstrap'

void bootstrap().catch(() => {})

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryProvider>
      <AppRouter />
    </QueryProvider>
  </StrictMode>,
)
