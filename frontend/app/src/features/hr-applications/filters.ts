import { ApplicationStatus, HrSortQueryParameter } from '@/api/generated/application/schemas'

export interface FiltersValue {
  statuses: ApplicationStatus[]
  vacancyId?: string
  createdFrom?: string
  createdTo?: string
  sort: HrSortQueryParameter
}

export const DEFAULT_FILTERS: FiltersValue = {
  statuses: [],
  vacancyId: undefined,
  createdFrom: undefined,
  createdTo: undefined,
  sort: HrSortQueryParameter.createdAtDesc,
}

export const ALL_STATUSES: ApplicationStatus[] = [
  ApplicationStatus.submitted,
  ApplicationStatus.seen,
  ApplicationStatus.interview,
  ApplicationStatus.rejected,
  ApplicationStatus.offer,
  ApplicationStatus.withdrawn,
]

export const SORT_LABEL: Record<HrSortQueryParameter, string> = {
  createdAtDesc: 'Сначала новые',
  updatedAtDesc: 'Недавно обновлённые',
}
