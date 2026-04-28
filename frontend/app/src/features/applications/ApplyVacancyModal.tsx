import { useEffect, useMemo, useRef, useState } from 'react'
import { z } from 'zod'
import { useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/shared/ui/dialog'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shared/ui/form'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'
import { Textarea } from '@/shared/ui/textarea'
import { Button } from '@/shared/ui/button'
import { useToast } from '@/shared/hooks/use-toast'
import { useGetResume } from '@/api/generated/candidate/resume/resume'
import {
  useCreateApplication,
  getListMyApplicationsQueryKey,
} from '@/api/generated/application/candidate-applications/candidate-applications'
import { AppError } from '@/shared/api/http/client'
import { useDefaultResumeId } from '@/features/resume-default'

const applySchema = z.object({
  resumeId: z.string().min(1, 'Выберите резюме'),
  coverLetter: z.string().max(2000).optional(),
})

type ApplyFormData = z.infer<typeof applySchema>

interface Props {
  vacancyId: string
  open: boolean
  onOpenChange: (open: boolean) => void
  onApplied?: () => void
}

export function ApplyVacancyModal({
  vacancyId,
  open,
  onOpenChange,
  onApplied,
}: Props) {
  const qc = useQueryClient()
  const { toast } = useToast()
  const resumes = useGetResume({ query: { enabled: open } })
  const create = useCreateApplication()
  const { defaultResumeId } = useDefaultResumeId()
  const [serverError, setServerError] = useState<string | null>(null)

  const publishedResumes = useMemo(
    () => (resumes.data ?? []).filter((r) => (r.status ?? 0) === 1 && r.id),
    [resumes.data],
  )

  const initialResumeId =
    defaultResumeId && publishedResumes.some((r) => r.id === defaultResumeId)
      ? defaultResumeId
      : (publishedResumes[0]?.id ?? '')

  const form = useForm<ApplyFormData>({
    resolver: zodResolver(applySchema),
    defaultValues: { resumeId: initialResumeId, coverLetter: '' },
  })

  const formRef = useRef(form)
  useEffect(() => {
    formRef.current = form
  })

  useEffect(() => {
    if (open) {
      formRef.current.reset({ resumeId: initialResumeId, coverLetter: '' })
    }
  }, [open, initialResumeId])

  function handleOpenChange(next: boolean) {
    if (next) setServerError(null)
    onOpenChange(next)
  }

  async function onSubmit(values: ApplyFormData) {
    setServerError(null)
    try {
      await create.mutateAsync({
        data: {
          vacancyId,
          resumeId: values.resumeId,
          coverLetter: values.coverLetter ? values.coverLetter : undefined,
        },
      })
      await qc.invalidateQueries({
        queryKey: getListMyApplicationsQueryKey(),
      })
      toast({ title: 'Отклик отправлен' })
      onOpenChange(false)
      onApplied?.()
    } catch (e) {
      const msg =
        e instanceof AppError
          ? e.message || 'Не удалось отправить отклик'
          : 'Не удалось отправить отклик'
      setServerError(msg)
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  const noPublished = !resumes.isLoading && publishedResumes.length === 0

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Отклик на вакансию</DialogTitle>
          <DialogDescription>
            Выберите опубликованное резюме и при желании добавьте сопроводительное письмо.
          </DialogDescription>
        </DialogHeader>

        {resumes.isLoading ? (
          <p className="text-sm text-muted-foreground">Загрузка резюме…</p>
        ) : noPublished ? (
          <p className="text-sm text-destructive">
            Нет опубликованных резюме. Создайте и опубликуйте резюме, чтобы откликнуться.
          </p>
        ) : (
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(onSubmit)}
              noValidate
              className="space-y-4"
            >
              {serverError && (
                <p role="alert" className="text-sm text-destructive">
                  {serverError}
                </p>
              )}
              <FormField
                control={form.control}
                name="resumeId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Резюме</FormLabel>
                    <FormControl>
                      <Select
                        value={field.value}
                        onValueChange={field.onChange}
                      >
                        <SelectTrigger>
                          <SelectValue placeholder="Выберите резюме" />
                        </SelectTrigger>
                        <SelectContent>
                          {publishedResumes.map((r) => (
                            <SelectItem key={r.id} value={r.id ?? ''}>
                              {r.name ?? '—'}
                              {defaultResumeId === r.id ? ' (основное)' : ''}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="coverLetter"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Сопроводительное письмо (необязательно)</FormLabel>
                    <FormControl>
                      <Textarea
                        rows={5}
                        maxLength={2000}
                        placeholder="Расскажите, почему вы подходите на эту роль"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                    <p className="text-xs text-muted-foreground">
                      {(field.value ?? '').length} / 2000
                    </p>
                  </FormItem>
                )}
              />
              <DialogFooter>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => onOpenChange(false)}
                  disabled={create.isPending}
                >
                  Отмена
                </Button>
                <Button type="submit" disabled={create.isPending}>
                  {create.isPending ? 'Отправка…' : 'Откликнуться'}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        )}
      </DialogContent>
    </Dialog>
  )
}
