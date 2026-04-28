import {
  ActorType,
  type HrApplicationStatusHistoryItem,
} from '@/api/generated/application/schemas'
import { ApplicationStatusBadge } from '@/shared/ui/ApplicationStatusBadge'
import { EmptyState } from '@/shared/ui/EmptyState'

const ROLE_LABEL: Record<ActorType, string> = {
  candidate: 'Кандидат',
  hr: 'HR',
  system: 'Система',
}

function formatDateTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleString()
}

interface Props {
  items: HrApplicationStatusHistoryItem[]
  isLoading?: boolean
}

export function HistoryTimeline({ items, isLoading }: Props) {
  if (isLoading) {
    return <p className="text-sm text-muted-foreground">Загрузка…</p>
  }
  if (items.length === 0) {
    return (
      <EmptyState title="История пуста" description="Ещё нет изменений статуса." />
    )
  }
  return (
    <ol className="space-y-3">
      {items.map((h, i) => (
        <li
          key={`${h.createdAt}-${i}`}
          className="rounded-md border bg-card p-3"
        >
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div className="flex items-center gap-2">
              <ApplicationStatusBadge status={h.status} />
              <span className="text-sm text-muted-foreground">
                · {ROLE_LABEL[h.changedByRole] ?? h.changedByRole}
              </span>
            </div>
            <span className="text-xs text-muted-foreground">
              {formatDateTime(h.createdAt)}
            </span>
          </div>
          {h.comment && (
            <p className="mt-2 whitespace-pre-wrap text-sm">{h.comment}</p>
          )}
        </li>
      ))}
    </ol>
  )
}
