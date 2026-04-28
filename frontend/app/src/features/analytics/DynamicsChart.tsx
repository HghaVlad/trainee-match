import type { ApplicationDynamicsPoint } from '@/api/generated/application/schemas'

interface Props {
  points: ApplicationDynamicsPoint[]
}

function formatBucket(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleDateString()
}

const SERIES: Array<{
  key: keyof Omit<ApplicationDynamicsPoint, 'bucketStart'>
  label: string
  color: string
}> = [
  { key: 'createdCount', label: 'Создано', color: 'bg-sky-500' },
  { key: 'seenCount', label: 'Просмотрено', color: 'bg-indigo-400' },
  { key: 'interviewCount', label: 'Интервью', color: 'bg-indigo-600' },
  { key: 'offerCount', label: 'Офферы', color: 'bg-emerald-500' },
  { key: 'rejectedCount', label: 'Отказы', color: 'bg-rose-500' },
  { key: 'withdrawnCount', label: 'Отозваны', color: 'bg-zinc-400' },
]

export function DynamicsChart({ points }: Props) {
  if (points.length === 0) {
    return (
      <p className="text-sm text-muted-foreground">
        Нет данных за выбранный период.
      </p>
    )
  }

  const max = Math.max(
    1,
    ...points.flatMap((p) => SERIES.map((s) => p[s.key])),
  )

  return (
    <div className="space-y-3">
      <div className="flex flex-wrap gap-3 text-xs">
        {SERIES.map((s) => (
          <div key={s.key} className="flex items-center gap-1.5">
            <span className={`inline-block h-3 w-3 rounded-sm ${s.color}`} />
            <span>{s.label}</span>
          </div>
        ))}
      </div>
      <div className="overflow-x-auto">
        <div
          className="flex items-end gap-2"
          style={{ minHeight: '160px' }}
        >
          {points.map((p) => (
            <div
              key={p.bucketStart}
              className="flex shrink-0 flex-col items-center gap-1"
              style={{ width: '64px' }}
            >
              <div className="flex h-40 items-end gap-0.5">
                {SERIES.map((s) => {
                  const v = p[s.key]
                  const heightPct = (v / max) * 100
                  return (
                    <div
                      key={s.key}
                      className="flex w-2 flex-col-reverse"
                      title={`${s.label}: ${v}`}
                    >
                      <div
                        className={`w-2 rounded-t ${s.color}`}
                        style={{ height: `${heightPct}%` }}
                        role="img"
                        aria-label={`${s.label}: ${v}`}
                      />
                    </div>
                  )
                })}
              </div>
              <div className="text-[10px] text-muted-foreground">
                {formatBucket(p.bucketStart)}
              </div>
            </div>
          ))}
        </div>
      </div>
      <div className="overflow-x-auto rounded-md border">
        <table className="w-full text-sm">
          <thead className="bg-muted/50 text-xs">
            <tr>
              <th className="px-2 py-1 text-left">Период</th>
              {SERIES.map((s) => (
                <th key={s.key} className="px-2 py-1 text-right">
                  {s.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {points.map((p) => (
              <tr key={p.bucketStart} className="border-t">
                <td className="px-2 py-1">{formatBucket(p.bucketStart)}</td>
                {SERIES.map((s) => (
                  <td
                    key={s.key}
                    className="px-2 py-1 text-right tabular-nums"
                  >
                    {p[s.key]}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
