import { useParams, Link } from 'react-router'
import { useGetAdminCandidatesId } from '@/api/generated/candidate/admin/admin'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'

export default function CandidateDetailPage() {
  const { id = '' } = useParams<{ id: string }>()
  const { data, isLoading, error, refetch } = useGetAdminCandidatesId(id, {
    query: { enabled: Boolean(id) },
  })

  if (isLoading) return <LoadingState />
  if (error || !data) return <ErrorState onRetry={() => refetch()} />

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <div className="flex items-center gap-2">
        <Link to="/candidates" className="text-sm text-muted-foreground underline">
          ← Все кандидаты
        </Link>
      </div>

      <h1 className="text-2xl font-bold">Кандидат {data.full_name ?? id}</h1>

      <div className="rounded-lg border bg-card p-4 space-y-2">
        <p><strong>ID:</strong> {data.id}</p>
        <p><strong>Имя:</strong> {data.full_name ?? '—'}</p>
        <p><strong>Город:</strong> {data.city ?? '—'}</p>
        <p><strong>Telegram:</strong> {data.telegram ? `@${data.telegram}` : '—'}</p>
        <p><strong>Телефон:</strong> {data.phone ?? '—'}</p>
        <p><strong>День рождения:</strong> {data.birthday ?? '—'}</p>
      </div>

      <div className="rounded-lg border bg-card p-4">
        <Link
          to={`/candidate/${id}/resumes`}
          className="text-primary underline"
        >
          Посмотреть резюме →
        </Link>
      </div>
    </div>
  )
}