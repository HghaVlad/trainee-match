import { useState } from 'react'
import { useNavigate } from 'react-router'
import { useQueryClient } from '@tanstack/react-query'
import { Button } from '@/shared/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/shared/ui/dialog'
import { useToast } from '@/shared/hooks/use-toast'
import {
  useDeleteCompaniesCompanyIdVacanciesVacancyId,
  usePostCompaniesCompanyIdVacanciesVacancyIdArchive,
  usePostCompaniesCompanyIdVacanciesVacancyIdPublish,
  getGetCompaniesCompanyIdVacanciesQueryKey,
  getGetCompaniesCompanyIdVacanciesVacancyIdQueryKey,
} from '@/api/generated/company/vacancy/vacancy'
import {
  DtoVacancyFullResponseStatus,
  type DtoVacancyFullResponseStatus as VacancyStatus,
} from '@/api/generated/company/schemas'
import { AppError } from '@/shared/api/http/client'

interface Props {
  companyId: string
  vacancyId: string
  status?: VacancyStatus | string
  isAdmin: boolean
  variant?: 'row' | 'detail'
}

function errorMessage(e: unknown, fallback: string): string {
  return e instanceof AppError ? e.message || fallback : fallback
}

export function VacancyActions({
  companyId,
  vacancyId,
  status,
  isAdmin,
  variant = 'detail',
}: Props) {
  const qc = useQueryClient()
  const { toast } = useToast()
  const navigate = useNavigate()

  const publishMut = usePostCompaniesCompanyIdVacanciesVacancyIdPublish()
  const archiveMut = usePostCompaniesCompanyIdVacanciesVacancyIdArchive()
  const deleteMut = useDeleteCompaniesCompanyIdVacanciesVacancyId()

  const [confirmDelete, setConfirmDelete] = useState(false)
  const [confirmPublish, setConfirmPublish] = useState(false)
  const [confirmArchive, setConfirmArchive] = useState(false)

  const isDraft = status === DtoVacancyFullResponseStatus.draft
  const isPublished = status === DtoVacancyFullResponseStatus.published

  async function invalidate() {
    await Promise.all([
      qc.invalidateQueries({
        queryKey: getGetCompaniesCompanyIdVacanciesQueryKey(companyId),
      }),
      qc.invalidateQueries({
        queryKey: getGetCompaniesCompanyIdVacanciesVacancyIdQueryKey(
          companyId,
          vacancyId,
        ),
      }),
    ])
  }

  async function onPublish() {
    try {
      await publishMut.mutateAsync({ companyId, vacancyId })
      await invalidate()
      toast({ title: 'Вакансия опубликована' })
      setConfirmPublish(false)
    } catch (e) {
      toast({
        title: 'Ошибка',
        description: errorMessage(e, 'Не удалось опубликовать вакансию'),
        variant: 'destructive',
      })
    }
  }

  async function onArchive() {
    try {
      await archiveMut.mutateAsync({ companyId, vacancyId })
      await invalidate()
      toast({ title: 'Вакансия в архиве' })
      setConfirmArchive(false)
    } catch (e) {
      toast({
        title: 'Ошибка',
        description: errorMessage(e, 'Не удалось архивировать вакансию'),
        variant: 'destructive',
      })
    }
  }

  async function onDelete() {
    try {
      await deleteMut.mutateAsync({ companyId, vacancyId })
      await qc.invalidateQueries({
        queryKey: getGetCompaniesCompanyIdVacanciesQueryKey(companyId),
      })
      toast({ title: 'Вакансия удалена' })
      setConfirmDelete(false)
      if (variant === 'detail') {
        navigate(`/company/${companyId}/vacancies`)
      }
    } catch (e) {
      toast({
        title: 'Ошибка',
        description: errorMessage(e, 'Не удалось удалить вакансию'),
        variant: 'destructive',
      })
    }
  }

  const size = variant === 'row' ? 'sm' : 'default'

  return (
    <div
      className={
        variant === 'row'
          ? 'flex flex-wrap items-center justify-end gap-2'
          : 'flex flex-wrap gap-2 border-t pt-4'
      }
    >
      {isDraft && (
        <Button
          size={size}
          variant="outline"
          onClick={() => setConfirmPublish(true)}
          disabled={publishMut.isPending}
        >
          Опубликовать
        </Button>
      )}
      {isPublished && (
        <Button
          size={size}
          variant="outline"
          onClick={() => setConfirmArchive(true)}
          disabled={archiveMut.isPending}
        >
          В архив
        </Button>
      )}
      {isAdmin && (
        <Button
          size={size}
          variant={variant === 'row' ? 'ghost' : 'destructive'}
          className={variant === 'row' ? 'text-destructive' : undefined}
          onClick={() => setConfirmDelete(true)}
          disabled={deleteMut.isPending}
        >
          Удалить
        </Button>
      )}

      <ConfirmDialog
        open={confirmPublish}
        onOpenChange={setConfirmPublish}
        title="Опубликовать вакансию?"
        description="Вакансия станет видна кандидатам."
        confirmLabel={publishMut.isPending ? 'Публикация…' : 'Опубликовать'}
        confirmDisabled={publishMut.isPending}
        onConfirm={onPublish}
      />
      <ConfirmDialog
        open={confirmArchive}
        onOpenChange={setConfirmArchive}
        title="Архивировать вакансию?"
        description="Вакансия станет неактивной для кандидатов."
        confirmLabel={archiveMut.isPending ? 'Архивирование…' : 'В архив'}
        confirmDisabled={archiveMut.isPending}
        onConfirm={onArchive}
      />
      <ConfirmDialog
        open={confirmDelete}
        onOpenChange={setConfirmDelete}
        title="Удалить вакансию?"
        description="Это действие нельзя отменить."
        confirmLabel={deleteMut.isPending ? 'Удаление…' : 'Удалить'}
        confirmDisabled={deleteMut.isPending}
        destructive
        onConfirm={onDelete}
      />
    </div>
  )
}

function ConfirmDialog({
  open,
  onOpenChange,
  title,
  description,
  confirmLabel,
  confirmDisabled,
  destructive,
  onConfirm,
}: {
  open: boolean
  onOpenChange: (next: boolean) => void
  title: string
  description: string
  confirmLabel: string
  confirmDisabled: boolean
  destructive?: boolean
  onConfirm: () => void
}) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={confirmDisabled}
          >
            Отмена
          </Button>
          <Button
            type="button"
            variant={destructive ? 'destructive' : 'default'}
            onClick={onConfirm}
            disabled={confirmDisabled}
          >
            {confirmLabel}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
