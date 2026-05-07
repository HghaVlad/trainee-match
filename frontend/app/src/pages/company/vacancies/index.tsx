import { useMemo, useState } from 'react'
import { Link, Navigate, useNavigate, useParams } from 'react-router'
import type { ColumnDef } from '@tanstack/react-table'
import { Button } from '@/shared/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/ui/card'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { DataTable } from '@/shared/ui/DataTable'
import { CursorPagination } from '@/shared/ui/CursorPagination'
import { useGetCompaniesCompanyIdVacancies } from '@/api/generated/company/vacancy/vacancy'
import {
  DtoVacancyFullResponseStatus,
  type DtoVacancyByCompListItemResponse,
  type DtoVacancyFullResponseStatus as VacancyStatus,
} from '@/api/generated/company/schemas'
import { useSession } from '@/shared/session/useSession'
import {
  VacancyActions,
  VacancyStatusBadge,
} from '@/features/company-vacancies'

const PAGE_SIZE = 20

type StatusFilter = 'all' | VacancyStatus

function getStatus(item: DtoVacancyByCompListItemResponse): VacancyStatus {
  return (item.status as VacancyStatus | undefined) ??
    DtoVacancyFullResponseStatus.draft
}

function formatDate(value: string | undefined): string {
  if (!value) return '—'
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return '—'
  return d.toLocaleDateString('ru-RU')
}

export default function CompanyVacanciesPage() {
  const { companyId } = useParams<{ companyId: string }>()
  if (!companyId) return <Navigate to="/company" replace />
  return <VacanciesList companyId={companyId} />
}

function VacanciesList({ companyId }: { companyId: string }) {
  const { companies } = useSession()
  const isAdmin =
    companies.find((c) => c.id === companyId)?.role === 'admin'

  const [cursor, setCursor] = useState<string | undefined>(undefined)
  const [statusFilter, setStatusFilter] = useState<StatusFilter>('all')

  const params = { cursor, limit: PAGE_SIZE }
  const query = useGetCompaniesCompanyIdVacancies(companyId, params)

  const allItems = useMemo(
    () => query.data?.vacancies ?? [],
    [query.data?.vacancies],
  )
  const items = useMemo(
    () =>
      statusFilter === 'all'
        ? allItems
        : allItems.filter((v) => getStatus(v) === statusFilter),
    [allItems, statusFilter],
  )

  const columns: ColumnDef<DtoVacancyByCompListItemResponse>[] = useMemo(
    () => [
      {
        header: 'Название',
        accessorKey: 'title',
        cell: ({ row }) => {
          const v = row.original
          return (
            <Link
              to={`/company/${companyId}/vacancies/${v.id ?? ''}`}
              className="font-medium text-primary underline"
            >
              {v.title || '—'}
            </Link>
          )
        },
      },
      {
        header: 'Статус',
        id: 'status',
        cell: ({ row }) => <VacancyStatusBadge status={getStatus(row.original)} />,
      },
      {
        header: 'Создана',
        accessorKey: 'createdAt',
        cell: ({ row }) => formatDate(row.original.createdAt),
      },
      {
        id: 'actions',
        header: '',
        cell: ({ row }) => {
          const v = row.original
          if (!v.id) return null
          return (
            <div className="flex items-center justify-end gap-2">
              <VacancyActions
                companyId={companyId}
                vacancyId={v.id}
                status={getStatus(v)}
                isAdmin={isAdmin}
                variant="row"
              />
            </div>
          )
        },
      },
    ],
    [companyId, isAdmin],
  )

  if (query.isLoading) return <LoadingState />
  if (query.isError) return <ErrorState onRetry={() => query.refetch()} />

  return (
    <div className="mx-auto max-w-5xl space-y-6 p-6">
      <Card>
        <CardHeader className="flex-row items-start justify-between gap-4">
          <div>
            <CardTitle>Вакансии</CardTitle>
            <CardDescription>
              Управление вакансиями вашей компании.
            </CardDescription>
          </div>
          <div className="flex items-center gap-2">
            <Select
              value={statusFilter}
              onValueChange={(v) => setStatusFilter(v as StatusFilter)}
            >
              <SelectTrigger className="w-44">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Все статусы</SelectItem>
                <SelectItem value={DtoVacancyFullResponseStatus.draft}>
                  Черновики
                </SelectItem>
                <SelectItem value={DtoVacancyFullResponseStatus.published}>
                  Опубликованные
                </SelectItem>
                <SelectItem value={DtoVacancyFullResponseStatus.archived}>
                  В архиве
                </SelectItem>
              </SelectContent>
            </Select>
            <CreateButton companyId={companyId} />
          </div>
        </CardHeader>
        <CardContent>
          {items.length === 0 ? (
            <EmptyState
              title="Вакансий нет"
              description={
                allItems.length === 0
                  ? 'Создайте первую вакансию, чтобы начать поиск кандидатов.'
                  : 'По выбранному фильтру ничего не найдено.'
              }
            />
          ) : (
            <DataTable columns={columns} data={items} />
          )}
          <div className="mt-4">
            <CursorPagination
              nextCursor={query.data?.nextCursor}
              onNext={() => setCursor(query.data?.nextCursor ?? undefined)}
              isLoading={query.isFetching}
            />
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

function CreateButton({ companyId }: { companyId: string }) {
  const navigate = useNavigate()
  return (
    <Button onClick={() => navigate(`/company/${companyId}/vacancies/new`)}>
      Создать
    </Button>
  )
}
