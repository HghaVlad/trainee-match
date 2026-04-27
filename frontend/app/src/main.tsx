import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import '@/app/styles/globals.css'
import { QueryProvider } from '@/app/providers/QueryProvider'
import { ErrorBoundary } from '@/app/providers/ErrorBoundary'
import { AppRouter } from '@/app/router/router'
import { bootstrap } from '@/shared/session/bootstrap'
import { Toaster } from '@/shared/ui/toaster'

void bootstrap().catch(() => {})

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ErrorBoundary>
      <QueryProvider>
        <AppRouter />
        <Toaster />
      </QueryProvider>
    </ErrorBoundary>
  </StrictMode>,
)
