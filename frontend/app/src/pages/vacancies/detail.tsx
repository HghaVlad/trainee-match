import { useParams, Link } from 'react-router'
import { useGetVacanciesVacancyId } from '@/api/generated/company/vacancy/vacancy'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { Button } from '@/shared/ui/button'

export default function VacancyDetailPage() {
  const { id = '' } = useParams<{ id: string }>()
  const { data, isLoading, error, refetch } = useGetVacanciesVacancyId(id, {
    query: { enabled: Boolean(id) },
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
      <Button disabled title="Скоро будет доступно">
        Откликнуться
      </Button>
    </div>
  )
}
