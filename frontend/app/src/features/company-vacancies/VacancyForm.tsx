import type { UseFormReturn } from 'react-hook-form'
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
  hoursPerWeekFrom?: number
  hoursPerWeekTo?: number
  durationFromDays?: number
  durationToDays?: number
  flexibleSchedule?: boolean
  internshipToOffer?: boolean
}

const optionalRange = (opts: { min: number; max?: number; label?: string }) =>
  z
    .union([
      z.literal(''),
      z.coerce
        .number()
        .min(opts.min, `Не меньше ${opts.min}`)
        .refine(
          (v) => opts.max === undefined || v <= opts.max,
          `Не больше ${opts.max}`,
        ),
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
    salaryFrom: optionalRange({ min: 0 }),
    salaryTo: optionalRange({ min: 0 }),
    isPaid: z.boolean().optional(),
    hoursPerWeekFrom: optionalRange({ min: 1, max: 168 }),
    hoursPerWeekTo: optionalRange({ min: 1, max: 168 }),
    durationFromDays: optionalRange({ min: 1, max: 730 }),
    durationToDays: optionalRange({ min: 1, max: 730 }),
    flexibleSchedule: z.boolean().optional(),
    internshipToOffer: z.boolean().optional(),
  })
  .superRefine((val, ctx) => {
    function checkRange(
      from: number | '' | undefined,
      to: number | '' | undefined,
      path: string[],
      label: string,
    ) {
      if (typeof from === 'number' && typeof to === 'number' && to < from) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path,
          message: `${label}: «До» должно быть не меньше «От»`,
        })
      }
    }
    checkRange(val.salaryFrom, val.salaryTo, ['salaryTo'], 'Зарплата')
    checkRange(
      val.hoursPerWeekFrom,
      val.hoursPerWeekTo,
      ['hoursPerWeekTo'],
      'Часы в неделю',
    )
    checkRange(
      val.durationFromDays,
      val.durationToDays,
      ['durationToDays'],
      'Длительность',
    )
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
    hoursPerWeekFrom: v?.hoursPerWeekFrom ?? '',
    hoursPerWeekTo: v?.hoursPerWeekTo ?? '',
    durationFromDays: v?.durationFromDays ?? '',
    durationToDays: v?.durationToDays ?? '',
    flexibleSchedule: v?.flexibleSchedule ?? false,
    internshipToOffer: v?.internshipToOffer ?? false,
  }
}

function num(value: number | '' | undefined): number | undefined {
  return typeof value === 'number' ? value : undefined
}

export function VacancyForm({
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
      salaryFrom: num(values.salaryFrom),
      salaryTo: num(values.salaryTo),
      isPaid: values.isPaid,
      hoursPerWeekFrom: num(values.hoursPerWeekFrom),
      hoursPerWeekTo: num(values.hoursPerWeekTo),
      durationFromDays: num(values.durationFromDays),
      durationToDays: num(values.durationToDays),
      flexibleSchedule: values.flexibleSchedule,
      internshipToOffer: values.internshipToOffer,
    }
    if (values.employmentType) {
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
          <div className="grid grid-cols-2 gap-3">
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
                      <SelectItem value={DtoVacancyCreateRequestWorkFormat.onsite}>
                        Офис
                      </SelectItem>
                      <SelectItem value={DtoVacancyCreateRequestWorkFormat.remote}>
                        Удалённо
                      </SelectItem>
                      <SelectItem value={DtoVacancyCreateRequestWorkFormat.hybrid}>
                        Гибрид
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />
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
                        value={DtoVacancyUpdateRequestEmploymentType.internship}
                      >
                        Стажировка
                      </SelectItem>
                      <SelectItem
                        value={DtoVacancyUpdateRequestEmploymentType.full_time}
                      >
                        Полная занятость
                      </SelectItem>
                      <SelectItem
                        value={DtoVacancyUpdateRequestEmploymentType.part_time}
                      >
                        Частичная занятость
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>
          <fieldset className="space-y-3 rounded-lg border p-3">
            <legend className="px-1 text-sm font-medium">Зарплата</legend>
            <div className="flex gap-3">
              <NumberFormField
                form={form}
                name="salaryFrom"
                label="От, ₽"
                min={0}
                step={1000}
              />
              <NumberFormField
                form={form}
                name="salaryTo"
                label="До, ₽"
                min={0}
                step={1000}
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
          </fieldset>
          <fieldset className="space-y-3 rounded-lg border p-3">
            <legend className="px-1 text-sm font-medium">Часы в неделю</legend>
            <div className="flex gap-3">
              <NumberFormField
                form={form}
                name="hoursPerWeekFrom"
                label="От"
                min={1}
                max={168}
                step={1}
              />
              <NumberFormField
                form={form}
                name="hoursPerWeekTo"
                label="До"
                min={1}
                max={168}
                step={1}
              />
            </div>
            <FormField
              control={form.control}
              name="flexibleSchedule"
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
                  <FormLabel className="!m-0">Гибкий график</FormLabel>
                  <FormMessage />
                </FormItem>
              )}
            />
          </fieldset>
          <fieldset className="space-y-3 rounded-lg border p-3">
            <legend className="px-1 text-sm font-medium">
              Длительность (дней)
            </legend>
            <div className="flex gap-3">
              <NumberFormField
                form={form}
                name="durationFromDays"
                label="От"
                min={1}
                max={730}
                step={1}
              />
              <NumberFormField
                form={form}
                name="durationToDays"
                label="До"
                min={1}
                max={730}
                step={1}
              />
            </div>
            <FormField
              control={form.control}
              name="internshipToOffer"
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
                  <FormLabel className="!m-0">
                    С возможностью трудоустройства
                  </FormLabel>
                  <FormMessage />
                </FormItem>
              )}
            />
          </fieldset>
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

type NumberFieldName =
  | 'salaryFrom'
  | 'salaryTo'
  | 'hoursPerWeekFrom'
  | 'hoursPerWeekTo'
  | 'durationFromDays'
  | 'durationToDays'

interface NumberFieldProps {
  form: UseFormReturn<FormValues>
  name: NumberFieldName
  label: string
  min?: number
  max?: number
  step?: number
}

function NumberFormField({
  form,
  name,
  label,
  min,
  max,
  step,
}: NumberFieldProps) {
  return (
    <FormField
      control={form.control}
      name={name}
      render={({ field }) => (
        <FormItem className="flex-1">
          <FormLabel>{label}</FormLabel>
          <FormControl>
            <Input
              type="number"
              min={min}
              max={max}
              step={step}
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
  )
}
