import { Link } from 'react-router'
import { useGetAdminCandidates } from '@/api/generated/candidate/admin/admin'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'

export default function CandidatesPage() {
  const { data, isLoading, error, refetch } = useGetAdminCandidates()

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

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <h1 className="text-2xl font-bold">Кандидаты</h1>

      <ul className="space-y-2">
        {data.map((candidate) => (
          <li key={candidate.id} className="rounded-lg border bg-card p-4">
            <Link
              to={`/candidate/${candidate.id}`}
              className="text-lg font-medium text-primary underline"
            >
              {candidate.full_name ?? candidate.id}
            </Link>
            <p className="text-sm text-muted-foreground">
              {candidate.city ?? '—'} • {candidate.telegram ? `@${candidate.telegram}` : '—'}
            </p>
          </li>
        ))}
      </ul>
    </div>
  )
}