import {
  HrAllowedAction,
  ChangeApplicationStatusRequestStatus,
} from '@/api/generated/application/schemas'

export const ACTION_TO_STATUS: Record<
  HrAllowedAction,
  ChangeApplicationStatusRequestStatus
> = {
  [HrAllowedAction.markSeen]: ChangeApplicationStatusRequestStatus.seen,
  [HrAllowedAction.moveToInterview]:
    ChangeApplicationStatusRequestStatus.interview,
  [HrAllowedAction.reject]: ChangeApplicationStatusRequestStatus.rejected,
  [HrAllowedAction.makeOffer]: ChangeApplicationStatusRequestStatus.offer,
}

export const ACTION_LABEL: Record<HrAllowedAction, string> = {
  [HrAllowedAction.markSeen]: 'Отметить просмотренным',
  [HrAllowedAction.moveToInterview]: 'Пригласить на интервью',
  [HrAllowedAction.reject]: 'Отклонить',
  [HrAllowedAction.makeOffer]: 'Сделать оффер',
}

export const ACTION_DIALOG_TITLE: Record<HrAllowedAction, string> = {
  [HrAllowedAction.markSeen]: 'Отметить как просмотренный?',
  [HrAllowedAction.moveToInterview]: 'Пригласить на интервью?',
  [HrAllowedAction.reject]: 'Отклонить отклик?',
  [HrAllowedAction.makeOffer]: 'Сделать оффер?',
}

export function isDestructiveAction(a: HrAllowedAction): boolean {
  return a === HrAllowedAction.reject
}
