import { useState } from 'react'
import { z } from 'zod'
import { useQueryClient } from '@tanstack/react-query'
import { FormWrapper } from '@/shared/ui/Form'
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
} from '@/shared/ui/form'
import { Input } from '@/shared/ui/input'
import { Button } from '@/shared/ui/button'
import {
  usePostCandidate,
  usePatchCandidate,
  getGetCandidateMeQueryKey,
} from '@/api/generated/candidate/candidate/candidate'
import type { DtoCandidateResponse } from '@/api/generated/candidate/schemas'
import { AppError } from '@/shared/api/http/client'

const today = new Date().toISOString().slice(0, 10)

const profileSchema = z.object({
  phone: z.string().min(1, 'Введите телефон').max(32),
  telegram: z
    .string()
    .min(1, 'Введите Telegram')
    .regex(/^@[A-Za-z0-9_]{3,32}$/, 'Формат: @username (3–32 символа)'),
  city: z.string().min(1, 'Введите город').max(100),
  birthday: z
    .string()
    .min(1, 'Укажите дату рождения')
    .refine((v) => v <= today, 'Дата не может быть в будущем'),
})

export type CandidateProfileFormData = z.infer<typeof profileSchema>

interface Props {
  mode: 'create' | 'edit'
  initial?: DtoCandidateResponse
  onSuccess?: () => void
  onCancel?: () => void
}

function backendDateToInput(value?: string): string {
  if (!value) return ''
  const dotMatch = /^(\d{2})\.(\d{2})\.(\d{4})$/.exec(value)
  if (dotMatch) return `${dotMatch[3]}-${dotMatch[2]}-${dotMatch[1]}`
  const isoMatch = /^(\d{4})-(\d{2})-(\d{2})/.exec(value)
  if (isoMatch) return value.slice(0, 10)
  return ''
}

export function CandidateProfileForm({ mode, initial, onSuccess, onCancel }: Props) {
  const qc = useQueryClient()
  const [error, setError] = useState<string | null>(null)
  const createMut = usePostCandidate()
  const updateMut = usePatchCandidate()

  async function handleSubmit(data: CandidateProfileFormData) {
    setError(null)
    try {
      const [y, m, d] = data.birthday.split('-')
      const birthdayApi = `${d}.${m}.${y}`
      const payload = {
        phone: data.phone,
        telegram: data.telegram,
        city: data.city,
        birthday: birthdayApi,
      }
      if (mode === 'create') {
        await createMut.mutateAsync({ data: payload })
      } else {
        await updateMut.mutateAsync({ data: payload })
      }
      await qc.invalidateQueries({ queryKey: getGetCandidateMeQueryKey() })
      onSuccess?.()
    } catch (e) {
      setError(e instanceof AppError ? e.message : 'Не удалось сохранить профиль')
    }
  }

  return (
    <FormWrapper
      schema={profileSchema}
      defaultValues={{
        phone: initial?.phone ?? '',
        telegram: initial?.telegram ?? '',
        city: initial?.city ?? '',
        birthday: backendDateToInput(initial?.birthday),
      }}
      onSubmit={handleSubmit}
    >
      {(form) => (
        <div className="space-y-4">
          {error && (
            <p role="alert" className="text-sm text-destructive">
              {error}
            </p>
          )}
          <FormField
            control={form.control}
            name="phone"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Телефон</FormLabel>
                <FormControl>
                  <Input autoComplete="tel" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="telegram"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Telegram</FormLabel>
                <FormControl>
                  <Input placeholder="@username" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="city"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Город</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="birthday"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Дата рождения</FormLabel>
                <FormControl>
                  <Input type="date" max={today} {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <div className="flex gap-2">
            <Button type="submit" disabled={form.formState.isSubmitting}>
              {form.formState.isSubmitting ? 'Сохранение...' : 'Сохранить'}
            </Button>
            {onCancel && (
              <Button
                type="button"
                variant="outline"
                onClick={onCancel}
                disabled={form.formState.isSubmitting}
              >
                Отмена
              </Button>
            )}
          </div>
        </div>
      )}
    </FormWrapper>
  )
}
