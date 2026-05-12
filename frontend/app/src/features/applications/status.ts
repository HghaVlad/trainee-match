import { ApplicationStatus } from '@/api/generated/application/schemas'
import type { CandidateApplicationListItem } from '@/api/generated/application/schemas'

export const ACTIVE_APPLICATION_STATUSES: readonly ApplicationStatus[] = [
  ApplicationStatus.submitted,
  ApplicationStatus.seen,
  ApplicationStatus.interview,
] as const

export const TERMINAL_APPLICATION_STATUSES: readonly ApplicationStatus[] = [
  ApplicationStatus.rejected,
  ApplicationStatus.offer,
  ApplicationStatus.withdrawn,
] as const

export function isActiveStatus(s: ApplicationStatus): boolean {
  return ACTIVE_APPLICATION_STATUSES.includes(s)
}

export function findActiveApplication(
  items: CandidateApplicationListItem[] | undefined,
  vacancyId: string,
): CandidateApplicationListItem | undefined {
  if (!items) return undefined
  return items.find(
    (it) => it.vacancyId === vacancyId && isActiveStatus(it.status),
  )
}

export const STATUS_LABEL: Record<ApplicationStatus, string> = {
  submitted: 'Отправлено',
  seen: 'Просмотрено',
  interview: 'Интервью',
  rejected: 'Отклонено',
  offer: 'Оффер',
  withdrawn: 'Отозвано',
}
