import {
  ApplicationStatus,
  type StatusFunnel,
} from '@/api/generated/application/schemas'
import { APPLICATION_STATUS_LABEL } from '@/shared/ui/ApplicationStatusBadge'

interface Props {
  funnel: StatusFunnel
}

const ORDER: ApplicationStatus[] = [
  ApplicationStatus.submitted,
  ApplicationStatus.seen,
  ApplicationStatus.interview,
  ApplicationStatus.offer,
  ApplicationStatus.rejected,
  ApplicationStatus.withdrawn,
]

const BAR_COLOR: Record<ApplicationStatus, string> = {
  submitted: 'bg-slate-400',
  seen: 'bg-sky-500',
  interview: 'bg-indigo-500',
  offer: 'bg-emerald-500',
  rejected: 'bg-rose-500',
  withdrawn: 'bg-zinc-400',
}

function valueOf(funnel: StatusFunnel, key: ApplicationStatus): number {
  switch (key) {
    case ApplicationStatus.submitted:
      return funnel.submitted
    case ApplicationStatus.seen:
      return funnel.seen
    case ApplicationStatus.interview:
      return funnel.interview
    case ApplicationStatus.offer:
      return funnel.offer
    case ApplicationStatus.rejected:
      return funnel.rejected
    case ApplicationStatus.withdrawn:
      return funnel.withdrawn
  }
}

export function StatusFunnelChart({ funnel }: Props) {
  const total = ORDER.reduce((sum, k) => sum + valueOf(funnel, k), 0)
  const max = Math.max(1, ...ORDER.map((k) => valueOf(funnel, k)))

  return (
    <div className="space-y-2">
      {ORDER.map((status) => {
        const count = valueOf(funnel, status)
        const widthPct = (count / max) * 100
        const pctOfTotal = total === 0 ? 0 : (count / total) * 100
        return (
          <div key={status} className="flex items-center gap-3">
            <div className="w-28 shrink-0 text-sm">
              {APPLICATION_STATUS_LABEL[status]}
            </div>
            <div className="relative flex-1">
              <div className="h-6 rounded bg-muted">
                <div
                  className={`h-6 rounded ${BAR_COLOR[status]}`}
                  style={{ width: `${widthPct}%` }}
                  role="progressbar"
                  aria-valuenow={count}
                  aria-valuemin={0}
                  aria-valuemax={max}
                  aria-label={APPLICATION_STATUS_LABEL[status]}
                />
              </div>
            </div>
            <div className="w-32 shrink-0 text-right text-sm tabular-nums">
              {count}{' '}
              <span className="text-xs text-muted-foreground">
                ({pctOfTotal.toFixed(1)}%)
              </span>
            </div>
          </div>
        )
      })}
      <div className="pt-2 text-xs text-muted-foreground">
        Всего: <span className="tabular-nums">{total}</span>
      </div>
    </div>
  )
}
