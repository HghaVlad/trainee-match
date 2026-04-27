import { useParams, Link } from 'react-router'
import { useGetCompaniesId } from '@/api/generated/company/company/company'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'

export default function CompanyDetailPage() {
  const { id = '' } = useParams<{ id: string }>()
  const { data, isLoading, error, refetch } = useGetCompaniesId(id, {
    query: { enabled: Boolean(id) },
  })

  if (isLoading) return <LoadingState />
  if (error || !data) return <ErrorState onRetry={() => refetch()} />

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <Link to="/companies" className="text-sm text-muted-foreground underline">
        ← Все компании
      </Link>
      <h1 className="text-2xl font-bold">{data.name ?? '—'}</h1>
      {data.website && (
        <a
          href={data.website}
          target="_blank"
          rel="noreferrer"
          className="text-primary underline"
        >
          {data.website}
        </a>
      )}
      {data.description && <p>{data.description}</p>}
      <p className="text-sm text-muted-foreground">
        Открытых вакансий: {data.openVacanciesCount ?? 0}
      </p>
    </div>
  )
}
