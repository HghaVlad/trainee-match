import { useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router'
import { useListCompanyApplications } from '@/api/generated/application/hr-applications/hr-applications'
import type { ListCompanyApplicationsParams } from '@/api/generated/application/schemas'
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
): ListCompanyApplicationsParams {
  return {
    statuses: filters.statuses.length > 0 ? filters.statuses : undefined,
    vacancyId: filters.vacancyId,
    createdFrom: filters.createdFrom,
    createdTo: filters.createdTo,
    cursor,
    limit: PAGE_SIZE,
    sort: filters.sort,
  }
}

export default function CompanyApplicationsPage() {
  const { companyId = '' } = useParams<{ companyId: string }>()
  const navigate = useNavigate()
  const [filters, setFilters] = useState<FiltersValue>(DEFAULT_FILTERS)
  const [cursor, setCursor] = useState<string | undefined>(undefined)

  const params = useMemo(() => buildParams(filters, cursor), [filters, cursor])
  const query = useListCompanyApplications(companyId, params)

  function applyFilters(next: FiltersValue) {
    setCursor(undefined)
    setFilters(next)
  }

  return (
    <div className="mx-auto max-w-6xl p-6 space-y-4">
      <h1 className="text-2xl font-bold">Отклики</h1>
      <ApplicationsListFilters
        companyId={companyId}
        value={filters}
        onChange={applyFilters}
      />
      {query.isLoading ? (
        <LoadingState />
      ) : query.isError || !query.data ? (
        <ErrorState onRetry={() => query.refetch()} />
      ) : query.data.data.length === 0 ? (
        <EmptyState
          title="Откликов нет"
          description="По заданным фильтрам ничего не найдено."
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
        />
      )}
    </div>
  )
}
