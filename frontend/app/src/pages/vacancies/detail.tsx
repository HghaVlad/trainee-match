import { useParams, Link } from 'react-router'
import { useGetVacanciesVacancyId } from '@/api/generated/company/vacancy/vacancy'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { ApplyVacancyButton } from '@/features/applications'

export default function VacancyDetailPage() {
  const { vacancyId = '' } = useParams<{ vacancyId: string }>()
  const { data, isLoading, error, refetch } = useGetVacanciesVacancyId(vacancyId, {
    query: { enabled: Boolean(vacancyId) },
  })

  if (isLoading) return <LoadingState />
  if (error || !data) return <ErrorState onRetry={() => refetch()} />

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <Link to="/vacancies" className="text-sm text-muted-foreground underline">
        ← Все вакансии
      </Link>
      <h1 className="text-2xl font-bold">{data.title ?? '—'}</h1>
      <p className="text-sm text-muted-foreground">
        {data.companyName ?? '—'} • {data.city ?? '—'} • {data.workFormat ?? ''}
      </p>
      {(data.salaryFrom ?? data.salaryTo) && (
        <p>
          Зарплата: {data.salaryFrom ?? ''}
          {data.salaryFrom && data.salaryTo ? '–' : ''}
          {data.salaryTo ?? ''} ₽
        </p>
      )}
      {data.description && (
        <div className="prose max-w-none whitespace-pre-wrap">
          {data.description}
        </div>
      )}
      {vacancyId && <ApplyVacancyButton vacancyId={vacancyId} />}
    </div>
  )
}
