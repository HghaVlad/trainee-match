import { QueryClient, QueryClientProvider, QueryCache, MutationCache } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import { type ReactNode } from 'react'
import { AppError, SessionExpiredError } from '@/shared/api/http/client'
import { toast } from '@/shared/hooks/use-toast'

function describeError(error: unknown): string {
  if (error instanceof AppError) return error.message || 'Что-то пошло не так'
  if (error instanceof Error) return error.message
  return 'Что-то пошло не так'
}

function createQueryClient() {
  return new QueryClient({
    queryCache: new QueryCache({
      onError: (error) => {
        if (error instanceof SessionExpiredError) return
        console.error('[QueryCache]', error)
      },
    }),
    mutationCache: new MutationCache({
      onError: (error, _vars, _ctx, mutation) => {
        if (error instanceof SessionExpiredError) return
        if (mutation.options.onError) return
        toast({
          title: 'Ошибка',
          description: describeError(error),
          variant: 'destructive',
        })
      },
    }),
    defaultOptions: {
      queries: {
        staleTime: 30_000,
        gcTime: 5 * 60_000,
        retry: (failureCount, error) => {
          if (error instanceof AppError && error.status < 500) return false
          return failureCount < 2
        },
        refetchOnWindowFocus: false,
      },
      mutations: {
        retry: 0,
      },
    },
  })
}

const queryClient = createQueryClient()

export function QueryProvider({ children }: { children: ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      {children}
      {import.meta.env.DEV && <ReactQueryDevtools />}
    </QueryClientProvider>
  )
}
