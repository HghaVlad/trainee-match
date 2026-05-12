import { useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/shared/ui/dialog'
import { Button } from '@/shared/ui/button'
import { Textarea } from '@/shared/ui/textarea'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shared/ui/form'
import { useToast } from '@/shared/hooks/use-toast'
import {
  useChangeApplicationStatus,
  getGetHrApplicationQueryKey,
  getGetHrApplicationHistoryQueryKey,
  getListCompanyApplicationsQueryKey,
  getListVacancyApplicationsQueryKey,
} from '@/api/generated/application/hr-applications/hr-applications'
import {
  type HrAllowedAction,
  type ChangeApplicationStatusRequestStatus,
} from '@/api/generated/application/schemas'
import { AppError } from '@/shared/api/http/client'
import {
  ACTION_DIALOG_TITLE,
  ACTION_LABEL,
  ACTION_TO_STATUS,
  isDestructiveAction,
} from './actionMap'

const schema = z.object({
  comment: z.string().max(2000, 'Не более 2000 символов').optional(),
})

type FormData = z.infer<typeof schema>

interface Props {
  open: boolean
  onOpenChange: (open: boolean) => void
  applicationId: string
  companyId: string
  vacancyId: string
  action: HrAllowedAction | null
}

export function ChangeStatusDialog({
  open,
  onOpenChange,
  applicationId,
  companyId,
  vacancyId,
  action,
}: Props) {
  const qc = useQueryClient()
  const { toast } = useToast()
  const mutation = useChangeApplicationStatus()
  const form = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { comment: '' },
  })

  function handleOpenChange(next: boolean) {
    if (!next) form.reset({ comment: '' })
    onOpenChange(next)
  }

  async function onSubmit(values: FormData) {
    if (!action) return
    const status: ChangeApplicationStatusRequestStatus = ACTION_TO_STATUS[action]
    const comment = values.comment?.trim()
    try {
      await mutation.mutateAsync({
        applicationId,
        data: comment ? { status, comment } : { status },
      })
      await Promise.all([
        qc.invalidateQueries({
          queryKey: getGetHrApplicationQueryKey(applicationId),
        }),
        qc.invalidateQueries({
          queryKey: getGetHrApplicationHistoryQueryKey(applicationId),
        }),
        qc.invalidateQueries({
          queryKey: getListCompanyApplicationsQueryKey(companyId),
        }),
        qc.invalidateQueries({
          queryKey: getListVacancyApplicationsQueryKey(vacancyId),
        }),
      ])
      toast({ title: 'Статус обновлён' })
      handleOpenChange(false)
    } catch (e) {
      const msg =
        e instanceof AppError
          ? e.message || 'Не удалось обновить статус'
          : 'Не удалось обновить статус'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  const title = action ? ACTION_DIALOG_TITLE[action] : 'Изменить статус'
  const confirmLabel = action ? ACTION_LABEL[action] : 'Подтвердить'
  const destructive = action ? isDestructiveAction(action) : false

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>
            Можно добавить необязательный комментарий (до 2000 символов).
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            noValidate
            className="space-y-4"
          >
            <FormField
              control={form.control}
              name="comment"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Комментарий</FormLabel>
                  <FormControl>
                    <Textarea
                      rows={3}
                      maxLength={2000}
                      placeholder="Комментарий (необязательно)"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => handleOpenChange(false)}
                disabled={mutation.isPending}
              >
                Отмена
              </Button>
              <Button
                type="submit"
                variant={destructive ? 'destructive' : 'default'}
                disabled={mutation.isPending || !action}
              >
                {mutation.isPending ? 'Сохранение…' : confirmLabel}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
