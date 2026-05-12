import { Badge } from '@/shared/ui/badge'
import { ApplicationStatus } from '@/api/generated/application/schemas'

const LABELS: Record<ApplicationStatus, string> = {
  submitted: 'Отправлено',
  seen: 'Просмотрено',
  interview: 'Интервью',
  rejected: 'Отклонено',
  offer: 'Оффер',
  withdrawn: 'Отозвано',
}

const VARIANTS: Record<
  ApplicationStatus,
  'default' | 'secondary' | 'outline' | 'destructive'
> = {
  submitted: 'outline',
  seen: 'secondary',
  interview: 'default',
  rejected: 'destructive',
  offer: 'default',
  withdrawn: 'outline',
}

function isStatus(v: unknown): v is ApplicationStatus {
  return (
    v === 'submitted' ||
    v === 'seen' ||
    v === 'interview' ||
    v === 'rejected' ||
    v === 'offer' ||
    v === 'withdrawn'
  )
}

interface Props {
  status?: ApplicationStatus | string
}

export function ApplicationStatusBadge({ status }: Props) {
  if (!isStatus(status)) {
    return <Badge variant="outline">—</Badge>
  }
  return <Badge variant={VARIANTS[status]}>{LABELS[status]}</Badge>
}

export const APPLICATION_STATUS_LABEL = LABELS
