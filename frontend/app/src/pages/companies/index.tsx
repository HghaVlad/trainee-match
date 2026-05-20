import { useMemo, useState } from 'react'
import { Link } from 'react-router'
import { useQueries } from '@tanstack/react-query'
import { useQueryClient } from '@tanstack/react-query'
import {
  getGetCompaniesQueryOptions,
} from '@/api/generated/company/company/company'
import { usePatchAdminCompaniesIdModeration } from '@/api/generated/company/admin-company/admin-company'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { Button } from '@/shared/ui/button'
import { useSession } from '@/shared/session/useSession'
import { useToast } from '@/shared/hooks/use-toast'
import { AppError } from '@/shared/api/http/client'

const PAGE_SIZE = 20

export default function CompaniesPage() {
  const [cursors, setCursors] = useState<Array<string | undefined>>([undefined])
  const { user } = useSession()
  const { toast } = useToast()
  const qc = useQueryClient()
  const archive = usePatchAdminCompaniesIdModeration()

  const queries = useQueries({
    queries: cursors.map((cursor) =>
      getGetCompaniesQueryOptions({ limit: PAGE_SIZE, cursor }),
    ),
  })

  const isLoading = queries[0]?.isLoading ?? false
  const isFetching = queries.some((q) => q.isFetching)
  const error = queries.find((q) => q.error)?.error
  const items = useMemo(() => {
    const seen = new Set<string>()
    const out: Array<{ id?: string; name?: string; openVacanciesCount?: number }> = []
    for (const q of queries) {
      for (const c of q.data?.companies ?? []) {
        if (!c.id || seen.has(c.id)) continue
        seen.add(c.id)
        out.push(c)
      }
    }
    return out
  }, [queries])
  const lastQuery = queries[queries.length - 1]
  const nextCursor = lastQuery?.data?.nextCursor ?? undefined

  async function onArchive(companyId: string) {
    try {
      await archive.mutateAsync({ id: companyId, data: { status: 'hidden' } })
      await qc.invalidateQueries({ queryKey: ['getGetCompaniesQueryOptions'] })
      toast({ title: 'Компания скрыта' })
    } catch (e) {
      const msg = e instanceof AppError ? e.message : 'Не удалось скрыть компанию'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  if (isLoading && items.length === 0) return <LoadingState />
  if (error && items.length === 0) {
    return <ErrorState onRetry={() => queries.forEach((q) => void q.refetch())} />
  }
  if (!isLoading && items.length === 0) {
    return <EmptyState title="Компании не найдены" />
  }

  const isPlatformAdmin = user?.role === 'PlatformAdmin'

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <h1 className="text-2xl font-bold">Компании</h1>
      <ul className="space-y-2">
        {items.map((c) => (
          <li key={c.id} className="rounded-lg border bg-card p-4">
            <div className="flex items-start justify-between gap-3">
              <div className="flex-1">
                <Link
                  to={`/companies/${c.id ?? ''}`}
                  className="text-lg font-medium text-primary underline"
                >
                  {c.name ?? '—'}
                </Link>
                <p className="text-sm text-muted-foreground">
                  Открытых вакансий: {c.openVacanciesCount ?? 0}
                </p>
              </div>
              {isPlatformAdmin && c.id && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => onArchive(c.id as string)}
                  disabled={archive.isPending}
                >
                  {archive.isPending ? '...' : 'Скрыть'}
                </Button>
              )}
            </div>
          </li>
        ))}
      </ul>
      {nextCursor && (
        <Button
          variant="outline"
          onClick={() => setCursors((prev) => [...prev, nextCursor])}
          disabled={isFetching}
        >
          {isFetching ? 'Загрузка…' : 'Загрузить ещё'}
        </Button>
      )}
    </div>
  )
}