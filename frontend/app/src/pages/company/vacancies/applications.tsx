import { useMemo, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router'
import { useListVacancyApplications } from '@/api/generated/application/hr-applications/hr-applications'
import { useGetCompaniesCompanyIdVacanciesVacancyId } from '@/api/generated/company/vacancy/vacancy'
import type { ListVacancyApplicationsParams } from '@/api/generated/application/schemas'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import {
  ApplicationsListFilters,
  ApplicationsTable,
  DEFAULT_FILTERS,
  type FiltersValue,
} from '@/features/hr-applications'

const PAGE_SIZE = 20

function buildParams(
  filters: FiltersValue,
  cursor: string | undefined,
): ListVacancyApplicationsParams {
  return {
    statuses: filters.statuses.length > 0 ? filters.statuses : undefined,
    createdFrom: filters.createdFrom,
    createdTo: filters.createdTo,
    cursor,
    limit: PAGE_SIZE,
    sort: filters.sort,
  }
}

export default function CompanyVacancyApplicationsPage() {
  const { companyId = '', vacancyId = '' } = useParams<{
    companyId: string
    vacancyId: string
  }>()
  const navigate = useNavigate()
  const [filters, setFilters] = useState<FiltersValue>(DEFAULT_FILTERS)
  const [cursor, setCursor] = useState<string | undefined>(undefined)

  const vacancyQ = useGetCompaniesCompanyIdVacanciesVacancyId(
    companyId,
    vacancyId,
  )
  const params = useMemo(() => buildParams(filters, cursor), [filters, cursor])
  const query = useListVacancyApplications(vacancyId, params)

  function applyFilters(next: FiltersValue) {
    setCursor(undefined)
    setFilters({ ...next, vacancyId })
  }

  const vacancyTitle = vacancyQ.data?.title || 'Вакансия'

  return (
    <div className="mx-auto max-w-6xl p-6 space-y-4">
      <div>
        <Link
          to={`/company/${companyId}/vacancies/${vacancyId}`}
          className="text-sm text-muted-foreground underline"
        >
          ← Вакансия
        </Link>
        <h1 className="mt-1 text-2xl font-bold">Отклики · {vacancyTitle}</h1>
      </div>
      <ApplicationsListFilters
        companyId={companyId}
        value={filters}
        onChange={applyFilters}
        hideVacancyFilter
      />
      {query.isLoading ? (
        <LoadingState />
      ) : query.isError || !query.data ? (
        <ErrorState onRetry={() => query.refetch()} />
      ) : query.data.data.length === 0 ? (
        <EmptyState
          title="Откликов нет"
          description="На эту вакансию ещё никто не откликнулся."
        />
      ) : (
        <ApplicationsTable
          items={query.data.data}
          nextCursor={query.data.nextCursor}
          isFetching={query.isFetching}
          onLoadMore={() =>
            setCursor(query.data?.nextCursor ?? undefined)
          }
          onOpen={(id) => navigate(`/company/${companyId}/applications/${id}`)}
          hideVacancyColumn
        />
      )}
    </div>
  )
}
