import { useEffect, useMemo, useRef, useState } from 'react'
import { Link, Navigate, useNavigate, useParams } from 'react-router'
import type { ColumnDef } from '@tanstack/react-table'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
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
import {
  useGetCompaniesCompanyIdVacanciesSearch,
} from '@/api/generated/company/vacancy/vacancy'
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
import { SearchIcon, XIcon } from 'lucide-react'

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

function useDebounce<T>(value: T, delay: number): T {
  const [debounced, setDebounced] = useState(value)
  useEffect(() => {
    const id = setTimeout(() => setDebounced(value), delay)
    return () => clearTimeout(id)
  }, [value, delay])
  return debounced
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
  const [searchInput, setSearchInput] = useState('')
  const searchQuery = useDebounce(searchInput, 300)
  const inputRef = useRef<HTMLInputElement>(null)

  // Reset cursor when filters change
  useEffect(() => {
    setCursor(undefined)
  }, [searchQuery, statusFilter])

  const params = {
    cursor,
    limit: PAGE_SIZE,
    query: searchQuery || undefined,
    status: statusFilter === 'all' ? undefined : statusFilter,
  }
  const query = useGetCompaniesCompanyIdVacanciesSearch(companyId, params)

  const items = useMemo(
    () => query.data?.vacancies ?? [],
    [query.data?.vacancies],
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
            <div className="relative">
              <SearchIcon className="pointer-events-none absolute left-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                ref={inputRef}
                type="search"
                placeholder="Поиск вакансий…"
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
                className="w-56 pl-8 pr-8"
              />
              {searchInput && (
                <button
                  type="button"
                  onClick={() => {
                    setSearchInput('')
                    inputRef.current?.focus()
                  }}
                  className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                >
                  <XIcon className="h-4 w-4" />
                </button>
              )}
            </div>
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
                searchQuery || statusFilter !== 'all'
                  ? 'По заданным параметрам ничего не найдено.'
                  : 'Создайте первую вакансию, чтобы начать поиск кандидатов.'
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
