import { useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
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
import { useToast } from '@/shared/hooks/use-toast'
import {
  getGetMyApplicationQueryKey,
  getGetMyApplicationHistoryQueryKey,
  getListMyApplicationsQueryKey,
  useWithdrawApplication,
} from '@/api/generated/application/candidate-applications/candidate-applications'
import { AppError } from '@/shared/api/http/client'

interface Props {
  applicationId: string
  open: boolean
  onOpenChange: (open: boolean) => void
  onWithdrawn?: () => void
}

export function WithdrawApplicationDialog({
  applicationId,
  open,
  onOpenChange,
  onWithdrawn,
}: Props) {
  const qc = useQueryClient()
  const { toast } = useToast()
  const [comment, setComment] = useState('')
  const withdraw = useWithdrawApplication()

  async function onConfirm() {
    try {
      await withdraw.mutateAsync({
        applicationId,
        data: comment ? { comment } : undefined,
      })
      await Promise.all([
        qc.invalidateQueries({ queryKey: getListMyApplicationsQueryKey() }),
        qc.invalidateQueries({
          queryKey: getGetMyApplicationQueryKey(applicationId),
        }),
        qc.invalidateQueries({
          queryKey: getGetMyApplicationHistoryQueryKey(applicationId),
        }),
      ])
      toast({ title: 'Отклик отозван' })
      setComment('')
      onOpenChange(false)
      onWithdrawn?.()
    } catch (e) {
      const msg =
        e instanceof AppError
          ? e.message || 'Не удалось отозвать отклик'
          : 'Не удалось отозвать отклик'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Отозвать отклик?</DialogTitle>
          <DialogDescription>
            Действие нельзя отменить. Можно добавить необязательный комментарий.
          </DialogDescription>
        </DialogHeader>
        <Textarea
          rows={3}
          maxLength={2000}
          placeholder="Комментарий (необязательно)"
          value={comment}
          onChange={(e) => setComment(e.target.value)}
        />
        <DialogFooter>
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={withdraw.isPending}
          >
            Отмена
          </Button>
          <Button
            type="button"
            variant="destructive"
            onClick={onConfirm}
            disabled={withdraw.isPending}
          >
            {withdraw.isPending ? 'Отзыв…' : 'Отозвать'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
