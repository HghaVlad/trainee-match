import { useParams, Link } from 'react-router'
import { useGetCompaniesId } from '@/api/generated/company/company/company'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'

function dash(value: unknown): string {
  if (value === undefined || value === null || value === '') return '—'
  return String(value)
}

export default function CompanyDetailPage() {
  const { companyId = '' } = useParams<{ companyId: string }>()
  const id = companyId
  const { data, isLoading, error, refetch } = useGetCompaniesId(id, {
    query: { enabled: Boolean(id) },
  })

  if (isLoading) return <LoadingState />
  if (error || !data) return <ErrorState onRetry={() => refetch()} />

  const openCount = data.openVacanciesCount ?? 0
  const vacanciesHref = data.id
    ? `/vacancies?company_id=${encodeURIComponent(data.id)}`
    : '/vacancies'

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <Link to="/companies" className="text-sm text-muted-foreground underline">
        ← Все компании
      </Link>
      <h1 className="text-2xl font-bold">{data.name ?? '—'}</h1>

      <dl className="grid grid-cols-1 gap-3 rounded-lg border bg-card p-4 sm:grid-cols-[200px,1fr]">
        <dt className="text-sm font-medium text-muted-foreground">Название</dt>
        <dd>{dash(data.name)}</dd>

        <dt className="text-sm font-medium text-muted-foreground">Сайт</dt>
        <dd>
          {data.website ? (
            <a
              href={data.website}
              target="_blank"
              rel="noreferrer"
              className="text-primary underline"
            >
              {data.website}
            </a>
          ) : (
            '—'
          )}
        </dd>

        <dt className="text-sm font-medium text-muted-foreground">Описание</dt>
        <dd className="whitespace-pre-line">{dash(data.description)}</dd>

        <dt className="text-sm font-medium text-muted-foreground">
          Открытых вакансий
        </dt>
        <dd>
          {openCount > 0 ? (
            <Link to={vacanciesHref} className="text-primary underline">
              {openCount}
            </Link>
          ) : (
            '0'
          )}
        </dd>
      </dl>
    </div>
  )
}
