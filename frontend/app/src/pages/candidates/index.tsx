import { useState } from 'react'
import { Link } from 'react-router'
import { useGetAdminCandidates } from '@/api/generated/candidate/admin/admin'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { Button } from '@/shared/ui/button'

const PAGE_SIZE = 20

export default function CandidatesPage() {
  const [page, setPage] = useState(1)
  const { data, isLoading, error, refetch, isFetching } = useGetAdminCandidates(
    { page, size: PAGE_SIZE },
  )

  if (isLoading) return <LoadingState />
  if (error) return <ErrorState onRetry={() => refetch()} />
  if (!data || data.length === 0) {
    return (
      <div className="mx-auto max-w-3xl p-6 space-y-4">
        <h1 className="text-2xl font-bold">Кандидаты</h1>
        <EmptyState title="Кандидаты не найдены" />
      </div>
    )
  }

  const hasMore = data.length >= PAGE_SIZE

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <h1 className="text-2xl font-bold">Кандидаты</h1>

      <ul className="space-y-2">
        {data.map((candidate) => (
          <li key={candidate.id} className="rounded-lg border bg-card p-4">
            <div className="flex items-start justify-between gap-3">
              <div className="flex-1">
                <Link
                  to={`/candidates/${candidate.id}`}
                  className="text-lg font-medium text-primary underline"
                >
                  {candidate.full_name ?? '—'}
                </Link>
                <p className="text-sm text-muted-foreground">
                  {candidate.city ?? '—'} • {candidate.telegram ? `${candidate.telegram}` : '—'}
                </p>
                <p className="text-xs text-muted-foreground mt-1">
                  ID: {candidate.id ?? '—'}
                </p>
              </div>
            </div>
          </li>
        ))}
      </ul>

      <div className="flex items-center justify-center gap-4">
        <Button
          variant="outline"
          size="sm"
          onClick={() => setPage((p) => Math.max(1, p - 1))}
          disabled={page <= 1 || isFetching}
        >
          Назад
        </Button>
        <span className="text-sm text-muted-foreground">
          Страница {page}
        </span>
        <Button
          variant="outline"
          size="sm"
          onClick={() => setPage((p) => p + 1)}
          disabled={!hasMore || isFetching}
        >
          {isFetching ? 'Загрузка…' : 'Вперёд'}
        </Button>
      </div>
    </div>
  )
}