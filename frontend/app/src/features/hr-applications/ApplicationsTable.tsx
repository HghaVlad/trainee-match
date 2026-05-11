import type { ColumnDef } from '@tanstack/react-table'
import { DataTable } from '@/shared/ui/DataTable'
import { CursorPagination } from '@/shared/ui/CursorPagination'
import { Button } from '@/shared/ui/button'
import { ApplicationStatusBadge } from '@/shared/ui/ApplicationStatusBadge'
import type { HrApplicationListItem } from '@/api/generated/application/schemas'

function formatDate(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleDateString()
}

interface Props {
  items: HrApplicationListItem[]
  nextCursor?: string | null
  isFetching?: boolean
  onLoadMore: () => void
  onOpen: (id: string) => void
  hideVacancyColumn?: boolean
}

export function ApplicationsTable({
  items,
  nextCursor,
  isFetching,
  onLoadMore,
  onOpen,
  hideVacancyColumn,
}: Props) {
  const columns: ColumnDef<HrApplicationListItem>[] = [
    {
      header: 'Кандидат',
      accessorKey: 'snapshot',
      cell: ({ row }) => {
        const s = row.original.snapshot
        const contact = s.email || s.telegram || ''
        return (
          <div className="flex flex-col">
            <span className="font-medium">{s.fullName}</span>
            {contact && (
              <span className="text-xs text-muted-foreground">{contact}</span>
            )}
          </div>
        )
      },
    },
    ...(hideVacancyColumn
      ? []
      : ([
          {
            header: 'Вакансия',
            accessorKey: 'vacancyTitle',
            cell: ({ row }) => row.original.vacancyTitle,
          },
        ] satisfies ColumnDef<HrApplicationListItem>[])),
    {
      header: 'Статус',
      accessorKey: 'status',
      cell: ({ row }) => <ApplicationStatusBadge status={row.original.status} />,
    },
    {
      header: 'Создано',
      accessorKey: 'createdAt',
      cell: ({ row }) => formatDate(row.original.createdAt),
    },
    {
      id: 'actions',
      header: '',
      cell: ({ row }) => (
        <Button
          size="sm"
          variant="outline"
          onClick={() => onOpen(row.original.id)}
        >
          Открыть
        </Button>
      ),
    },
  ]

  return (
    <div className="space-y-3">
      <DataTable columns={columns} data={items} />
      <CursorPagination
        nextCursor={nextCursor}
        onNext={onLoadMore}
        isLoading={isFetching}
      />
    </div>
  )
}
