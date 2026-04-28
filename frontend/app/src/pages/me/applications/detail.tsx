import { useState } from 'react'
import { Link, useParams } from 'react-router'
import {
  useGetMyApplication,
  useGetMyApplicationHistory,
} from '@/api/generated/application/candidate-applications/candidate-applications'
import { CandidateAllowedAction } from '@/api/generated/application/schemas'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { Button } from '@/shared/ui/button'
import { STATUS_LABEL, WithdrawApplicationDialog } from '@/features/applications'

function formatDateTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleString()
}

const ROLE_LABEL: Record<string, string> = {
  candidate: 'Вы',
  hr: 'HR',
  system: 'Система',
}

export default function MyApplicationDetailPage() {
  const { applicationId = '' } = useParams<{ applicationId: string }>()
  const detailQ = useGetMyApplication(applicationId)
  const historyQ = useGetMyApplicationHistory(applicationId)
  const [withdrawOpen, setWithdrawOpen] = useState(false)

  if (detailQ.isLoading) return <LoadingState />
  if (detailQ.error || !detailQ.data)
    return <ErrorState onRetry={() => detailQ.refetch()} />

  const app = detailQ.data.data
  const history = historyQ.data?.data ?? app.statusHistory ?? []
  const canWithdraw = (app.allowedActions ?? []).includes(
    CandidateAllowedAction.withdraw,
  )

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <Link
        to="/me/applications"
        className="text-sm text-muted-foreground underline"
      >
        ← Все отклики
      </Link>
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold">{app.vacancyTitle}</h1>
          <p className="text-sm text-muted-foreground">{app.companyName}</p>
          <p className="text-xs text-muted-foreground">
            Создан: {formatDateTime(app.createdAt)}
          </p>
        </div>
        <span className="rounded-full border px-2 py-0.5 text-xs">
          {STATUS_LABEL[app.status] ?? app.status}
        </span>
      </div>

      {app.coverLetter && (
        <section className="rounded-lg border bg-card p-4">
          <h2 className="mb-2 text-lg font-semibold">Сопроводительное письмо</h2>
          <p className="whitespace-pre-wrap text-sm">{app.coverLetter}</p>
        </section>
      )}

      <section className="rounded-lg border bg-card p-4">
        <h2 className="mb-2 text-lg font-semibold">История статусов</h2>
        {historyQ.isLoading ? (
          <p className="text-sm text-muted-foreground">Загрузка…</p>
        ) : history.length === 0 ? (
          <p className="text-sm text-muted-foreground">История пуста.</p>
        ) : (
          <ul className="space-y-2">
            {history.map((h, i) => (
              <li
                key={`${h.createdAt}-${i}`}
                className="flex items-center justify-between gap-3 text-sm"
              >
                <span>
                  {STATUS_LABEL[h.status] ?? h.status}{' '}
                  <span className="text-muted-foreground">
                    · {ROLE_LABEL[h.changedByRole] ?? h.changedByRole}
                  </span>
                </span>
                <span className="text-xs text-muted-foreground">
                  {formatDateTime(h.createdAt)}
                </span>
              </li>
            ))}
          </ul>
        )}
      </section>

      {canWithdraw && (
        <div>
          <Button
            type="button"
            variant="destructive"
            onClick={() => setWithdrawOpen(true)}
          >
            Отозвать отклик
          </Button>
          <WithdrawApplicationDialog
            applicationId={app.id}
            open={withdrawOpen}
            onOpenChange={setWithdrawOpen}
          />
        </div>
      )}
    </div>
  )
}
