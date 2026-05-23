import { useParams, Link } from 'react-router'
import { useGetAdminCandidatesIdResumes } from '@/api/generated/candidate/admin/admin'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'

export default function CandidateResumesPage() {
  const { id = '' } = useParams<{ id: string }>()
  const { data, isLoading, error, refetch } = useGetAdminCandidatesIdResumes(id, undefined, { query: { enabled: Boolean(id) } })

  if (isLoading) return <LoadingState />
  if (error) return <ErrorState onRetry={() => refetch()} />

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <div className="flex items-center gap-2">
        <Link
          to={`/candidate/${id}`}
          className="text-sm text-muted-foreground underline"
        >
          ← Профиль кандидата
        </Link>
      </div>

      <h1 className="text-2xl font-bold">Резюме кандидата</h1>

      {!data || data.length === 0 ? (
        <EmptyState title="Резюме не найдены" />
      ) : (
        <ul className="space-y-2">
          {data.map((resume) => (
            <li key={resume.id} className="rounded-lg border bg-card p-4">
              <Link
                to={`/candidate/${id}/resumes/${resume.id}`}
                className="text-lg font-medium text-primary underline"
              >
                {resume.name}
              </Link>
              <p className="text-sm text-muted-foreground">
                Статус: {resume.status ?? '—'} • Модерация: {resume.moderation_status ?? '—'}
              </p>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}