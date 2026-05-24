import { useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router'
import { useListCompanyApplications } from '@/api/generated/application/hr-applications/hr-applications'
import type { ListCompanyApplicationsParams } from '@/api/generated/application/schemas'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { AppError } from '@/shared/api/http/client'
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
    statuses: filters.statuses.length > 0
      ? (filters.statuses.join(',') as unknown as ListCompanyApplicationsParams['statuses'])
      : undefined,
    vacancyId: filters.vacancyId,
    createdFrom: filters.createdFrom ? `${filters.createdFrom}T00:00:00Z` : undefined,
    createdTo: filters.createdTo ? `${filters.createdTo}T23:59:59Z` : undefined,
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
  const query = useListCompanyApplications(companyId, params, {
    query: { retry: false },
  })

  function applyFilters(next: FiltersValue) {
    setCursor(undefined)
    setFilters(next)
  }

  const queryError: unknown = query.error
  const notFound = queryError instanceof AppError && queryError.status === 404
  const items = notFound ? [] : query.data?.data ?? []
  const showError = query.isError && !notFound

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
      ) : showError ? (
        <ErrorState onRetry={() => query.refetch()} />
      ) : items.length === 0 ? (
        <EmptyState
          title="Откликов нет"
          description="По заданным фильтрам ничего не найдено."
        />
      ) : (
        <ApplicationsTable
          items={items}
          nextCursor={query.data?.nextCursor}
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
