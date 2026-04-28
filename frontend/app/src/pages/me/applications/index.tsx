import { Link } from 'react-router'
import { useListMyApplications } from '@/api/generated/application/candidate-applications/candidate-applications'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { STATUS_LABEL } from '@/features/applications'

function formatDate(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleDateString()
}

export default function MyApplicationsPage() {
  const { data, isLoading, error, refetch } = useListMyApplications()

  if (isLoading) return <LoadingState />
  if (error) return <ErrorState onRetry={() => refetch()} />

  const items = data?.data ?? []

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <h1 className="text-2xl font-bold">Мои отклики</h1>
      {items.length === 0 ? (
        <EmptyState
          title="Откликов пока нет"
          description="Найдите вакансию и отправьте свой первый отклик."
        />
      ) : (
        <ul className="space-y-2">
          {items.map((it) => (
            <li
              key={it.id}
              className="rounded-lg border bg-card p-4"
            >
              <div className="flex items-start justify-between gap-3">
                <div className="min-w-0 flex-1">
                  <Link
                    to={`/me/applications/${it.id}`}
                    className="text-lg font-medium text-primary underline"
                  >
                    {it.vacancyTitle}
                  </Link>
                  <p className="text-sm text-muted-foreground">
                    {it.companyName}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    Отправлено: {formatDate(it.createdAt)}
                  </p>
                </div>
                <span className="shrink-0 rounded-full border px-2 py-0.5 text-xs">
                  {STATUS_LABEL[it.status] ?? it.status}
                </span>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
