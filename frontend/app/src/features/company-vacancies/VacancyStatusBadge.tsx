import { Badge } from '@/shared/ui/badge'
import { DtoVacancyFullResponseStatus } from '@/api/generated/company/schemas'

const LABELS: Record<DtoVacancyFullResponseStatus, string> = {
  draft: 'Черновик',
  published: 'Опубликовано',
  archived: 'В архиве',
}

const VARIANTS: Record<
  DtoVacancyFullResponseStatus,
  'default' | 'secondary' | 'outline' | 'destructive'
> = {
  draft: 'outline',
  published: 'default',
  archived: 'secondary',
}

interface Props {
  status?: DtoVacancyFullResponseStatus | string
}

function isStatus(v: unknown): v is DtoVacancyFullResponseStatus {
  return v === 'draft' || v === 'published' || v === 'archived'
}

export function VacancyStatusBadge({ status }: Props) {
  if (!isStatus(status)) {
    return <Badge variant="outline">—</Badge>
  }
  return <Badge variant={VARIANTS[status]}>{LABELS[status]}</Badge>
}
