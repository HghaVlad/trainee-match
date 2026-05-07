import { useEffect, useMemo, useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router'
import { useForm, useFieldArray, type UseFormReturn } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  useGetResumeId,
  usePatchResumeId,
} from '@/api/generated/candidate/resume/resume'
import type { DtoResumeData } from '@/api/generated/candidate/schemas'
import { useDeleteResume } from '@/shared/api/resume/deleteResume'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/shared/ui/form'
import { Input } from '@/shared/ui/input'
import { Textarea } from '@/shared/ui/textarea'
import { Button } from '@/shared/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/ui/select'
import { useToast } from '@/shared/hooks/use-toast'
import { AppError } from '@/shared/api/http/client'
import { SkillsMultiSelect } from '@/features/skill-catalog'

const today = new Date().toISOString().slice(0, 10)
const currentYear = new Date().getFullYear()
const YEARS: number[] = Array.from(
  { length: currentYear + 10 - 1950 + 1 },
  (_, i) => 1950 + i,
)
const MONTHS: { value: string; label: string }[] = [
  { value: '01', label: 'Январь' },
  { value: '02', label: 'Февраль' },
  { value: '03', label: 'Март' },
  { value: '04', label: 'Апрель' },
  { value: '05', label: 'Май' },
  { value: '06', label: 'Июнь' },
  { value: '07', label: 'Июль' },
  { value: '08', label: 'Август' },
  { value: '09', label: 'Сентябрь' },
  { value: '10', label: 'Октябрь' },
  { value: '11', label: 'Ноябрь' },
  { value: '12', label: 'Декабрь' },
]

const FORMAT_OPTIONS: { value: string; label: string }[] = [
  { value: 'remote', label: 'Удалёнка' },
  { value: 'onsite', label: 'Офис' },
  { value: 'hybrid', label: 'Гибрид' },
]
const ENGLISH_OPTIONS: { value: string; label: string }[] = [
  { value: 'A1', label: 'A1 — начальный' },
  { value: 'A2', label: 'A2 — элементарный' },
  { value: 'B1', label: 'B1 — средний' },
  { value: 'B2', label: 'B2 — выше среднего' },
  { value: 'C1', label: 'C1 — продвинутый' },
  { value: 'C2', label: 'C2 — профессиональный' },
  { value: 'native', label: 'Носитель' },
]

const educationSchema = z.object({
  university: z.string().optional().or(z.literal('')),
  faculty: z.string().optional().or(z.literal('')),
  specialization: z.string().optional().or(z.literal('')),
  level: z.string().optional().or(z.literal('')),
  format: z.string().optional().or(z.literal('')),
  start_year: z
    .union([z.literal(''), z.coerce.number().int().min(1900).max(2100)])
    .optional(),
  end_year: z
    .union([z.literal(''), z.coerce.number().int().min(1900).max(2100)])
    .optional(),
})
const workSchema = z.object({
  company: z.string().optional().or(z.literal('')),
  position: z.string().optional().or(z.literal('')),
  start_month: z.string().optional().or(z.literal('')),
  start_year: z.string().optional().or(z.literal('')),
  end_month: z.string().optional().or(z.literal('')),
  end_year: z.string().optional().or(z.literal('')),
  responsibilities: z.string().optional().or(z.literal('')),
})

const resumeSchema = z.object({
  name: z.string().min(1, 'Введите название').max(200),
  first_name: z.string().optional().or(z.literal('')),
  last_name: z.string().optional().or(z.literal('')),
  middle_name: z.string().optional().or(z.literal('')),
  email: z.string().email().optional().or(z.literal('')),
  phone: z.string().optional().or(z.literal('')),
  city: z.string().optional().or(z.literal('')),
  date_of_birth: z
    .string()
    .optional()
    .or(z.literal(''))
    .refine((v) => !v || v <= today, 'Дата не может быть в будущем'),
  citizenship: z.string().optional().or(z.literal('')),
  desired_format: z.string().optional().or(z.literal('')),
  english_level: z.string().optional().or(z.literal('')),
  additional_info: z.string().max(5000).optional().or(z.literal('')),
  portfolio_link: z.string().url().optional().or(z.literal('')),
  education: z.array(educationSchema),
  work_experiences: z.array(workSchema),
  skills_list: z.array(z.string()),
})
type ResumeFormData = z.infer<typeof resumeSchema>

type WorkRow = z.infer<typeof workSchema>

function parsePeriod(period?: string): {
  start_month: string
  start_year: string
  end_month: string
  end_year: string
} {
  const empty = { start_month: '', start_year: '', end_month: '', end_year: '' }
  if (!period) return empty
  const m = /^(\d{2})\.(\d{4})\s*[-–—]\s*(\d{2})\.(\d{4})$/.exec(period.trim())
  if (m) {
    return {
      start_month: m[1] ?? '',
      start_year: m[2] ?? '',
      end_month: m[3] ?? '',
      end_year: m[4] ?? '',
    }
  }
  const y = /^(\d{4})\s*[-–—]\s*(\d{4})$/.exec(period.trim())
  if (y) {
    return {
      start_month: '',
      start_year: y[1] ?? '',
      end_month: '',
      end_year: y[2] ?? '',
    }
  }
  return empty
}

function buildPeriod(row: WorkRow): string | undefined {
  const sm = row.start_month
  const sy = row.start_year
  const em = row.end_month
  const ey = row.end_year
  const left = sm && sy ? `${sm}.${sy}` : sy ? sy : ''
  const right = em && ey ? `${em}.${ey}` : ey ? ey : ''
  if (!left && !right) return undefined
  return `${left}${left || right ? ' — ' : ''}${right}`.trim()
}

function backendDateToInput(value?: string): string {
  if (!value) return ''
  const dotMatch = /^(\d{2})\.(\d{2})\.(\d{4})$/.exec(value)
  if (dotMatch) return `${dotMatch[3]}-${dotMatch[2]}-${dotMatch[1]}`
  const isoMatch = /^(\d{4})-(\d{2})-(\d{2})/.exec(value)
  if (isoMatch) return value.slice(0, 10)
  return ''
}

function inputDateToBackend(value: string): string | undefined {
  if (!value) return undefined
  const m = /^(\d{4})-(\d{2})-(\d{2})$/.exec(value)
  if (!m) return undefined
  return `${m[3]}.${m[2]}.${m[1]}`
}

function toForm(name?: string, d?: DtoResumeData): ResumeFormData {
  const works = (d?.work_experiences ?? []).map((w) => {
    const parts = parsePeriod(w.period)
    return {
      company: w.company ?? '',
      position: w.position ?? '',
      responsibilities: w.responsibilities ?? '',
      ...parts,
    }
  })
  return {
    name: name ?? '',
    first_name: d?.first_name ?? '',
    last_name: d?.last_name ?? '',
    middle_name: d?.middle_name ?? '',
    email: d?.email ?? '',
    phone: d?.phone ?? '',
    city: d?.city ?? '',
    date_of_birth: backendDateToInput(d?.date_of_birth),
    citizenship: d?.citizenship ?? '',
    desired_format: d?.desired_format ?? '',
    english_level: d?.english_level ?? '',
    additional_info: d?.additional_info ?? '',
    portfolio_link: d?.portfolio_link ?? '',
    education: (d?.education ?? []).map((e) => ({
      university: e.university ?? '',
      faculty: e.faculty ?? '',
      specialization: e.specialization ?? '',
      level: e.level ?? '',
      format: e.format ?? '',
      start_year: typeof e.start_year === 'number' ? e.start_year : '',
      end_year: typeof e.end_year === 'number' ? e.end_year : '',
    })),
    work_experiences: works,
    skills_list: d?.skills_list ?? [],
  }
}

export default function ResumeEditPage() {
  const { resumeId = '' } = useParams<{ resumeId: string }>()
  const id = resumeId
  const navigate = useNavigate()
  const { toast } = useToast()
  const detail = useGetResumeId(id, { query: { enabled: Boolean(id) } })
  const patch = usePatchResumeId()
  const del = useDeleteResume()
  const [savedAt, setSavedAt] = useState<Date | null>(null)
  const [confirmingDelete, setConfirmingDelete] = useState(false)

  const form = useForm<ResumeFormData>({
    resolver: zodResolver(resumeSchema) as never,
    defaultValues: toForm(),
  })

  useEffect(() => {
    if (detail.data) {
      form.reset(toForm(detail.data.name, detail.data.data))
    }
  }, [detail.data, form])

  function buildPayload(values: ResumeFormData) {
    const works = values.work_experiences.map((w) => ({
      company: w.company || undefined,
      position: w.position || undefined,
      period: buildPeriod(w),
      responsibilities: w.responsibilities || undefined,
    }))
    const education = values.education.map((e) => ({
      university: e.university || undefined,
      faculty: e.faculty || undefined,
      specialization: e.specialization || undefined,
      level: e.level || undefined,
      format: e.format || undefined,
      start_year:
        typeof e.start_year === 'number' ? e.start_year : undefined,
      end_year: typeof e.end_year === 'number' ? e.end_year : undefined,
    }))
    return {
      id,
      data: {
        name: values.name,
        data: {
          first_name: values.first_name || undefined,
          last_name: values.last_name || undefined,
          middle_name: values.middle_name || undefined,
          email: values.email || undefined,
          phone: values.phone || undefined,
          city: values.city || undefined,
          date_of_birth: inputDateToBackend(values.date_of_birth ?? ''),
          citizenship: values.citizenship || undefined,
          desired_format: values.desired_format || undefined,
          english_level: values.english_level || undefined,
          additional_info: values.additional_info || undefined,
          portfolio_link: values.portfolio_link || undefined,
          education,
          work_experiences: works,
          skills_list: values.skills_list,
        },
      },
    }
  }

  async function onManualSave(values: ResumeFormData) {
    try {
      await patch.mutateAsync(buildPayload(values))
      setSavedAt(new Date())
      toast({ title: 'Резюме сохранено' })
    } catch (e) {
      const msg = e instanceof AppError ? e.message : 'Не удалось сохранить резюме'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  async function onConfirmDelete() {
    try {
      await del.mutateAsync(id)
      toast({ title: 'Резюме удалено' })
      navigate('/me/resumes')
    } catch (e) {
      const msg = e instanceof AppError ? e.message : 'Не удалось удалить резюме'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    } finally {
      setConfirmingDelete(false)
    }
  }

  if (detail.isLoading) return <LoadingState />
  if (detail.error || !detail.data) return <ErrorState onRetry={() => detail.refetch()} />

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <Link to="/me/resumes" className="text-sm text-muted-foreground underline">
        ← Все резюме
      </Link>
      <h1 className="text-2xl font-bold">Редактирование резюме</h1>
      {savedAt && (
        <p className="text-xs text-muted-foreground">
          Сохранено в {savedAt.toLocaleTimeString()}
        </p>
      )}
      <Form {...form}>
        <form noValidate onSubmit={form.handleSubmit(onManualSave)} className="space-y-4">
          <BasicFields form={form} />
          <EducationList form={form} />
          <WorkList form={form} />
          <SkillsField form={form} />
          <div className="sticky bottom-0 -mx-2 flex items-center justify-end gap-2 border-t bg-background/95 p-3 backdrop-blur">
            <Button
              type="button"
              variant="destructive"
              onClick={() => setConfirmingDelete(true)}
              disabled={del.isPending}
            >
              Удалить
            </Button>
            <div className="flex-1" />
            <Button
              type="button"
              variant="outline"
              onClick={() => navigate('/me/resumes')}
            >
              Отмена
            </Button>
            <Button type="submit" disabled={patch.isPending}>
              {patch.isPending ? 'Сохранение…' : 'Сохранить'}
            </Button>
          </div>
        </form>
      </Form>

      {confirmingDelete && (
        <div
          role="dialog"
          aria-modal="true"
          aria-labelledby="resume-delete-title"
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
        >
          <div className="w-full max-w-md rounded-lg border bg-card p-6 shadow-lg">
            <h2 id="resume-delete-title" className="text-lg font-semibold">
              Удалить резюме?
            </h2>
            <p className="mt-2 text-sm text-muted-foreground">
              Это действие нельзя отменить. Все данные резюме будут удалены безвозвратно.
            </p>
            <div className="mt-4 flex justify-end gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => setConfirmingDelete(false)}
                disabled={del.isPending}
              >
                Отмена
              </Button>
              <Button
                type="button"
                variant="destructive"
                onClick={onConfirmDelete}
                disabled={del.isPending}
              >
                {del.isPending ? 'Удаление…' : 'Удалить'}
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

function BasicFields({ form }: { form: UseFormReturn<ResumeFormData> }) {
  return (
    <div className="space-y-3 rounded-lg border bg-card p-4">
      <FormField
        control={form.control}
        name="name"
        render={({ field }) => (
          <FormItem>
            <FormLabel>Название резюме</FormLabel>
            <FormControl><Input {...field} /></FormControl>
            <FormMessage />
          </FormItem>
        )}
      />
      <div className="grid grid-cols-2 gap-3">
        <FormField control={form.control} name="last_name" render={({ field }) => (
          <FormItem><FormLabel>Фамилия</FormLabel><FormControl><Input {...field} /></FormControl><FormMessage /></FormItem>
        )} />
        <FormField control={form.control} name="first_name" render={({ field }) => (
          <FormItem><FormLabel>Имя</FormLabel><FormControl><Input {...field} /></FormControl><FormMessage /></FormItem>
        )} />
        <FormField control={form.control} name="middle_name" render={({ field }) => (
          <FormItem><FormLabel>Отчество</FormLabel><FormControl><Input {...field} /></FormControl><FormMessage /></FormItem>
        )} />
        <FormField control={form.control} name="date_of_birth" render={({ field }) => (
          <FormItem>
            <FormLabel>Дата рождения</FormLabel>
            <FormControl><Input type="date" max={today} {...field} /></FormControl>
            <FormMessage />
          </FormItem>
        )} />
        <FormField control={form.control} name="citizenship" render={({ field }) => (
          <FormItem><FormLabel>Гражданство</FormLabel><FormControl><Input {...field} placeholder="Россия" /></FormControl><FormMessage /></FormItem>
        )} />
        <FormField control={form.control} name="city" render={({ field }) => (
          <FormItem><FormLabel>Город</FormLabel><FormControl><Input {...field} /></FormControl><FormMessage /></FormItem>
        )} />
        <FormField control={form.control} name="email" render={({ field }) => (
          <FormItem><FormLabel>Email</FormLabel><FormControl><Input type="email" {...field} /></FormControl><FormMessage /></FormItem>
        )} />
        <FormField control={form.control} name="phone" render={({ field }) => (
          <FormItem><FormLabel>Телефон</FormLabel><FormControl><Input {...field} /></FormControl><FormMessage /></FormItem>
        )} />
        <FormField control={form.control} name="desired_format" render={({ field }) => (
          <FormItem>
            <FormLabel>Желаемый формат</FormLabel>
            <Select value={field.value ?? ''} onValueChange={field.onChange}>
              <FormControl>
                <SelectTrigger><SelectValue placeholder="Не указан" /></SelectTrigger>
              </FormControl>
              <SelectContent>
                {FORMAT_OPTIONS.map((o) => (
                  <SelectItem key={o.value} value={o.value}>
                    {o.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>
        )} />
        <FormField control={form.control} name="english_level" render={({ field }) => (
          <FormItem>
            <FormLabel>Английский</FormLabel>
            <Select value={field.value ?? ''} onValueChange={field.onChange}>
              <FormControl>
                <SelectTrigger><SelectValue placeholder="Не указан" /></SelectTrigger>
              </FormControl>
              <SelectContent>
                {ENGLISH_OPTIONS.map((o) => (
                  <SelectItem key={o.value} value={o.value}>
                    {o.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>
        )} />
      </div>
      <FormField control={form.control} name="portfolio_link" render={({ field }) => (
        <FormItem><FormLabel>Портфолио</FormLabel><FormControl><Input type="url" {...field} placeholder="https://" /></FormControl><FormMessage /></FormItem>
      )} />
      <FormField control={form.control} name="additional_info" render={({ field }) => (
        <FormItem><FormLabel>Дополнительно</FormLabel><FormControl><Textarea rows={3} {...field} /></FormControl><FormMessage /></FormItem>
      )} />
    </div>
  )
}

function YearSelect({
  value,
  onChange,
  placeholder = 'Год',
}: {
  value: number | '' | undefined
  onChange: (next: number | '') => void
  placeholder?: string
}) {
  const display = typeof value === 'number' ? String(value) : ''
  return (
    <Select
      value={display}
      onValueChange={(v) => onChange(v ? Number(v) : '')}
    >
      <SelectTrigger>
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        {YEARS.map((y) => (
          <SelectItem key={y} value={String(y)}>
            {y}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}

function EducationList({ form }: { form: UseFormReturn<ResumeFormData> }) {
  const arr = useFieldArray({ control: form.control, name: 'education' })
  return (
    <div className="space-y-3 rounded-lg border bg-card p-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">Образование</h2>
        <Button type="button" variant="outline" onClick={() => arr.append({})}>
          Добавить
        </Button>
      </div>
      {arr.fields.map((f, i) => (
        <div key={f.id} className="grid grid-cols-2 gap-2 border-t pt-2">
          <FormField control={form.control} name={`education.${i}.university`} render={({ field }) => (
            <FormItem><FormLabel>Вуз</FormLabel><FormControl><Input {...field} /></FormControl></FormItem>
          )} />
          <FormField control={form.control} name={`education.${i}.faculty`} render={({ field }) => (
            <FormItem><FormLabel>Факультет</FormLabel><FormControl><Input {...field} /></FormControl></FormItem>
          )} />
          <FormField control={form.control} name={`education.${i}.specialization`} render={({ field }) => (
            <FormItem><FormLabel>Специальность</FormLabel><FormControl><Input {...field} /></FormControl></FormItem>
          )} />
          <FormField control={form.control} name={`education.${i}.level`} render={({ field }) => (
            <FormItem><FormLabel>Уровень</FormLabel><FormControl><Input {...field} /></FormControl></FormItem>
          )} />
          <FormField control={form.control} name={`education.${i}.start_year`} render={({ field }) => (
            <FormItem>
              <FormLabel>Начало обучения</FormLabel>
              <YearSelect
                value={
                  typeof field.value === 'number'
                    ? field.value
                    : field.value === ''
                      ? ''
                      : undefined
                }
                onChange={field.onChange}
              />
              <FormMessage />
            </FormItem>
          )} />
          <FormField control={form.control} name={`education.${i}.end_year`} render={({ field }) => (
            <FormItem>
              <FormLabel>Конец обучения</FormLabel>
              <YearSelect
                value={
                  typeof field.value === 'number'
                    ? field.value
                    : field.value === ''
                      ? ''
                      : undefined
                }
                onChange={field.onChange}
              />
              <FormMessage />
            </FormItem>
          )} />
          <Button type="button" variant="ghost" onClick={() => arr.remove(i)} className="col-span-2 justify-self-start">
            Удалить
          </Button>
        </div>
      ))}
    </div>
  )
}

function MonthYearSelect({
  monthValue,
  yearValue,
  onMonthChange,
  onYearChange,
  monthPlaceholder,
  yearPlaceholder,
}: {
  monthValue: string
  yearValue: string
  onMonthChange: (v: string) => void
  onYearChange: (v: string) => void
  monthPlaceholder: string
  yearPlaceholder: string
}) {
  return (
    <div className="grid grid-cols-2 gap-2">
      <Select value={monthValue} onValueChange={onMonthChange}>
        <SelectTrigger>
          <SelectValue placeholder={monthPlaceholder} />
        </SelectTrigger>
        <SelectContent>
          {MONTHS.map((m) => (
            <SelectItem key={m.value} value={m.value}>
              {m.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <Select value={yearValue} onValueChange={onYearChange}>
        <SelectTrigger>
          <SelectValue placeholder={yearPlaceholder} />
        </SelectTrigger>
        <SelectContent>
          {YEARS.map((y) => (
            <SelectItem key={y} value={String(y)}>
              {y}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  )
}

function WorkList({ form }: { form: UseFormReturn<ResumeFormData> }) {
  const arr = useFieldArray({ control: form.control, name: 'work_experiences' })
  return (
    <div className="space-y-3 rounded-lg border bg-card p-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">Опыт работы</h2>
        <Button type="button" variant="outline" onClick={() => arr.append({})}>
          Добавить
        </Button>
      </div>
      {arr.fields.map((f, i) => (
        <div key={f.id} className="space-y-2 border-t pt-2">
          <div className="grid grid-cols-2 gap-2">
            <FormField control={form.control} name={`work_experiences.${i}.company`} render={({ field }) => (
              <FormItem><FormLabel>Компания</FormLabel><FormControl><Input {...field} /></FormControl></FormItem>
            )} />
            <FormField control={form.control} name={`work_experiences.${i}.position`} render={({ field }) => (
              <FormItem><FormLabel>Должность</FormLabel><FormControl><Input {...field} /></FormControl></FormItem>
            )} />
          </div>
          <div className="grid grid-cols-2 gap-2">
            <div className="space-y-1">
              <span className="text-sm font-medium">Начало</span>
              <MonthYearSelect
                monthValue={form.watch(`work_experiences.${i}.start_month`) ?? ''}
                yearValue={form.watch(`work_experiences.${i}.start_year`) ?? ''}
                onMonthChange={(v) =>
                  form.setValue(`work_experiences.${i}.start_month`, v, {
                    shouldDirty: true,
                  })
                }
                onYearChange={(v) =>
                  form.setValue(`work_experiences.${i}.start_year`, v, {
                    shouldDirty: true,
                  })
                }
                monthPlaceholder="Месяц"
                yearPlaceholder="Год"
              />
            </div>
            <div className="space-y-1">
              <span className="text-sm font-medium">Окончание</span>
              <MonthYearSelect
                monthValue={form.watch(`work_experiences.${i}.end_month`) ?? ''}
                yearValue={form.watch(`work_experiences.${i}.end_year`) ?? ''}
                onMonthChange={(v) =>
                  form.setValue(`work_experiences.${i}.end_month`, v, {
                    shouldDirty: true,
                  })
                }
                onYearChange={(v) =>
                  form.setValue(`work_experiences.${i}.end_year`, v, {
                    shouldDirty: true,
                  })
                }
                monthPlaceholder="Месяц"
                yearPlaceholder="Год"
              />
              <p className="text-xs text-muted-foreground">
                Оставьте пустым, если работаете сейчас
              </p>
            </div>
          </div>
          <FormField control={form.control} name={`work_experiences.${i}.responsibilities`} render={({ field }) => (
            <FormItem><FormLabel>Обязанности</FormLabel><FormControl><Textarea rows={3} {...field} /></FormControl></FormItem>
          )} />
          <Button type="button" variant="ghost" onClick={() => arr.remove(i)}>
            Удалить
          </Button>
        </div>
      ))}
    </div>
  )
}

function SkillsField({ form }: { form: UseFormReturn<ResumeFormData> }) {
  const watched = form.watch('skills_list')
  const value = useMemo(() => watched ?? [], [watched])
  return (
    <div className="space-y-3 rounded-lg border bg-card p-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">Навыки</h2>
      </div>
      <p className="text-sm text-muted-foreground">
        Выберите навыки из общего каталога. Скрытые навыки можно искать по названию.
      </p>
      <SkillsMultiSelect
        value={value}
        onChange={(next) =>
          form.setValue('skills_list', next, { shouldDirty: true })
        }
      />
    </div>
  )
}
