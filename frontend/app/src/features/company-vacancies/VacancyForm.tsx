import { z } from 'zod'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Textarea } from '@/shared/ui/textarea'
import {
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
import { FormWrapper } from '@/shared/ui/Form'
import {
  DtoVacancyCreateRequestWorkFormat,
  DtoVacancyUpdateRequestEmploymentType,
} from '@/api/generated/company/schemas'

export interface VacancyFormPayload {
  title: string
  description?: string
  city?: string
  workFormat: typeof DtoVacancyCreateRequestWorkFormat[keyof typeof DtoVacancyCreateRequestWorkFormat]
  employmentType?: typeof DtoVacancyUpdateRequestEmploymentType[keyof typeof DtoVacancyUpdateRequestEmploymentType]
  salaryFrom?: number
  salaryTo?: number
  isPaid?: boolean
}

const optionalNonNegative = z
  .union([
    z.literal(''),
    z.coerce.number().min(0, 'Не меньше 0'),
  ])
  .optional()

const baseSchema = z
  .object({
    title: z
      .string()
      .min(3, 'Минимум 3 символа')
      .max(200, 'Максимум 200 символов'),
    description: z
      .string()
      .max(10000, 'Максимум 10000 символов')
      .optional()
      .or(z.literal('')),
    city: z
      .string()
      .max(200, 'Максимум 200 символов')
      .optional()
      .or(z.literal('')),
    workFormat: z.enum([
      DtoVacancyCreateRequestWorkFormat.onsite,
      DtoVacancyCreateRequestWorkFormat.remote,
      DtoVacancyCreateRequestWorkFormat.hybrid,
    ]),
    employmentType: z
      .enum([
        DtoVacancyUpdateRequestEmploymentType.internship,
        DtoVacancyUpdateRequestEmploymentType.full_time,
        DtoVacancyUpdateRequestEmploymentType.part_time,
      ])
      .optional()
      .or(z.literal('')),
    salaryFrom: optionalNonNegative,
    salaryTo: optionalNonNegative,
    isPaid: z.boolean().optional(),
  })
  .superRefine((val, ctx) => {
    const from = val.salaryFrom
    const to = val.salaryTo
    if (typeof from === 'number' && typeof to === 'number' && to < from) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ['salaryTo'],
        message: '«До» должно быть не меньше «От»',
      })
    }
  })

type FormValues = z.infer<typeof baseSchema>

interface Props {
  mode: 'create' | 'edit'
  defaultValues?: Partial<VacancyFormPayload>
  onSubmit: (payload: VacancyFormPayload) => Promise<void> | void
  isSubmitting?: boolean
  submitLabel?: string
  onCancel?: () => void
  serverError?: string | null
}

function toFormValues(v: Partial<VacancyFormPayload> | undefined): FormValues {
  return {
    title: v?.title ?? '',
    description: v?.description ?? '',
    city: v?.city ?? '',
    workFormat:
      (v?.workFormat as FormValues['workFormat'] | undefined) ??
      DtoVacancyCreateRequestWorkFormat.remote,
    employmentType:
      (v?.employmentType as FormValues['employmentType'] | undefined) ?? '',
    salaryFrom: v?.salaryFrom ?? '',
    salaryTo: v?.salaryTo ?? '',
    isPaid: v?.isPaid ?? false,
  }
}

export function VacancyForm({
  mode,
  defaultValues,
  onSubmit,
  isSubmitting,
  submitLabel,
  onCancel,
  serverError,
}: Props) {
  async function handleSubmit(values: FormValues) {
    const payload: VacancyFormPayload = {
      title: values.title,
      description: values.description ? values.description : undefined,
      city: values.city ? values.city : undefined,
      workFormat: values.workFormat,
      salaryFrom:
        typeof values.salaryFrom === 'number' ? values.salaryFrom : undefined,
      salaryTo:
        typeof values.salaryTo === 'number' ? values.salaryTo : undefined,
      isPaid: values.isPaid,
    }
    if (mode === 'edit' && values.employmentType) {
      payload.employmentType =
        values.employmentType as VacancyFormPayload['employmentType']
    }
    await onSubmit(payload)
  }

  return (
    <FormWrapper
      schema={baseSchema}
      defaultValues={toFormValues(defaultValues)}
      onSubmit={handleSubmit}
    >
      {(form) => (
        <div className="space-y-4">
          {serverError && (
            <p role="alert" className="text-sm text-destructive">
              {serverError}
            </p>
          )}
          <FormField
            control={form.control}
            name="title"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Название</FormLabel>
                <FormControl>
                  <Input maxLength={200} {...field} />
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
                  <Input maxLength={200} {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="workFormat"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Формат работы</FormLabel>
                <Select
                  value={field.value as string}
                  onValueChange={field.onChange}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem
                      value={DtoVacancyCreateRequestWorkFormat.onsite}
                    >
                      Офис
                    </SelectItem>
                    <SelectItem
                      value={DtoVacancyCreateRequestWorkFormat.remote}
                    >
                      Удалённо
                    </SelectItem>
                    <SelectItem
                      value={DtoVacancyCreateRequestWorkFormat.hybrid}
                    >
                      Гибрид
                    </SelectItem>
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />
          {mode === 'edit' && (
            <FormField
              control={form.control}
              name="employmentType"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Тип занятости</FormLabel>
                  <Select
                    value={typeof field.value === 'string' ? field.value : ''}
                    onValueChange={field.onChange}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Не указан" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem
                        value={
                          DtoVacancyUpdateRequestEmploymentType.internship
                        }
                      >
                        Стажировка
                      </SelectItem>
                      <SelectItem
                        value={
                          DtoVacancyUpdateRequestEmploymentType.full_time
                        }
                      >
                        Полная занятость
                      </SelectItem>
                      <SelectItem
                        value={
                          DtoVacancyUpdateRequestEmploymentType.part_time
                        }
                      >
                        Частичная занятость
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />
          )}
          <div className="flex gap-3">
            <FormField
              control={form.control}
              name="salaryFrom"
              render={({ field }) => (
                <FormItem className="flex-1">
                  <FormLabel>Зарплата от</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      min={0}
                      inputMode="numeric"
                      value={
                        field.value === undefined || field.value === null
                          ? ''
                          : String(field.value)
                      }
                      onChange={(e) => field.onChange(e.target.value)}
                      onBlur={field.onBlur}
                      name={field.name}
                      ref={field.ref}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="salaryTo"
              render={({ field }) => (
                <FormItem className="flex-1">
                  <FormLabel>До</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      min={0}
                      inputMode="numeric"
                      value={
                        field.value === undefined || field.value === null
                          ? ''
                          : String(field.value)
                      }
                      onChange={(e) => field.onChange(e.target.value)}
                      onBlur={field.onBlur}
                      name={field.name}
                      ref={field.ref}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>
          <FormField
            control={form.control}
            name="isPaid"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center gap-2 space-y-0">
                <FormControl>
                  <input
                    type="checkbox"
                    checked={Boolean(field.value)}
                    onChange={(e) => field.onChange(e.target.checked)}
                    onBlur={field.onBlur}
                    ref={field.ref}
                    name={field.name}
                    className="h-4 w-4"
                  />
                </FormControl>
                <FormLabel className="!m-0">Оплачиваемая</FormLabel>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="description"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Описание</FormLabel>
                <FormControl>
                  <Textarea rows={6} maxLength={10000} {...field} />
                </FormControl>
                <FormMessage />
                <p className="text-xs text-muted-foreground">
                  {(field.value ?? '').length} / 10000
                </p>
              </FormItem>
            )}
          />
          <div className="flex justify-end gap-2">
            {onCancel && (
              <Button
                type="button"
                variant="outline"
                onClick={onCancel}
                disabled={isSubmitting}
              >
                Отмена
              </Button>
            )}
            <Button
              type="submit"
              disabled={isSubmitting || form.formState.isSubmitting}
            >
              {isSubmitting ? 'Сохранение…' : (submitLabel ?? 'Сохранить')}
            </Button>
          </div>
        </div>
      )}
    </FormWrapper>
  )
}
