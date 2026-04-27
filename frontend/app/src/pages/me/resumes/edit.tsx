import { useEffect, useRef, useState } from 'react'
import { useParams, Link } from 'react-router'
import { useForm, useFieldArray, type UseFormReturn } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  useGetResumeId,
  usePatchResumeId,
} from '@/api/generated/candidate/resume/resume'
import type { DtoResumeData } from '@/api/generated/candidate/schemas'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/shared/ui/form'
import { Input } from '@/shared/ui/input'
import { Textarea } from '@/shared/ui/textarea'
import { Button } from '@/shared/ui/button'
import { SkillCombobox } from '@/features/skill-catalog'

const educationSchema = z.object({
  university: z.string().optional().or(z.literal('')),
  faculty: z.string().optional().or(z.literal('')),
  specialization: z.string().optional().or(z.literal('')),
  level: z.string().optional().or(z.literal('')),
  format: z.string().optional().or(z.literal('')),
  start_year: z.coerce.number().int().min(1900).max(2100).optional(),
  end_year: z.coerce.number().int().min(1900).max(2100).optional(),
})
const workSchema = z.object({
  company: z.string().optional().or(z.literal('')),
  position: z.string().optional().or(z.literal('')),
  period: z.string().optional().or(z.literal('')),
  responsibilities: z.string().optional().or(z.literal('')),
})

const resumeSchema = z.object({
  name: z.string().min(1, 'Введите название').max(200),
  first_name: z.string().optional().or(z.literal('')),
  last_name: z.string().optional().or(z.literal('')),
  email: z.string().email().optional().or(z.literal('')),
  phone: z.string().optional().or(z.literal('')),
  city: z.string().optional().or(z.literal('')),
  desired_format: z.string().optional().or(z.literal('')),
  english_level: z.string().optional().or(z.literal('')),
  additional_info: z.string().max(5000).optional().or(z.literal('')),
  portfolio_link: z.string().url().optional().or(z.literal('')),
  education: z.array(educationSchema),
  work_experiences: z.array(workSchema),
  skills_list: z.array(z.string()),
})
type ResumeFormData = z.infer<typeof resumeSchema>

function toForm(name?: string, d?: DtoResumeData): ResumeFormData {
  return {
    name: name ?? '',
    first_name: d?.first_name ?? '',
    last_name: d?.last_name ?? '',
    email: d?.email ?? '',
    phone: d?.phone ?? '',
    city: d?.city ?? '',
    desired_format: d?.desired_format ?? '',
    english_level: d?.english_level ?? '',
    additional_info: d?.additional_info ?? '',
    portfolio_link: d?.portfolio_link ?? '',
    education: d?.education ?? [],
    work_experiences: d?.work_experiences ?? [],
    skills_list: d?.skills_list ?? [],
  }
}

export default function ResumeEditPage() {
  const { id = '' } = useParams<{ id: string }>()
  const detail = useGetResumeId(id, { query: { enabled: Boolean(id) } })
  const patch = usePatchResumeId()
  const [savedAt, setSavedAt] = useState<Date | null>(null)

  const form = useForm<ResumeFormData>({
    resolver: zodResolver(resumeSchema) as never,
    defaultValues: toForm(),
  })

  useEffect(() => {
    if (detail.data) {
      form.reset(toForm(detail.data.name, detail.data.data))
    }
  }, [detail.data, form])

  const debounce = useRef<number | null>(null)
  const lastSaved = useRef<string>('')

  useEffect(() => {
    const sub = form.watch((values) => {
      if (!form.formState.isValid) return
      if (debounce.current) window.clearTimeout(debounce.current)
      debounce.current = window.setTimeout(() => {
        const json = JSON.stringify(values)
        if (json === lastSaved.current) return
        lastSaved.current = json
        patch.mutate(
          {
            id,
            data: {
              name: values.name,
              data: {
                first_name: values.first_name || undefined,
                last_name: values.last_name || undefined,
                email: values.email || undefined,
                phone: values.phone || undefined,
                city: values.city || undefined,
                desired_format: values.desired_format || undefined,
                english_level: values.english_level || undefined,
                additional_info: values.additional_info || undefined,
                portfolio_link: values.portfolio_link || undefined,
                education: values.education,
                work_experiences: values.work_experiences,
                skills_list: values.skills_list,
              },
            },
          },
          { onSuccess: () => setSavedAt(new Date()) },
        )
      }, 3000)
    })
    return () => {
      sub.unsubscribe()
      if (debounce.current) window.clearTimeout(debounce.current)
    }
  }, [form, id, patch])

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
        <form noValidate onSubmit={(e) => e.preventDefault()} className="space-y-4">
          <BasicFields form={form} />
          <EducationList form={form} />
          <WorkList form={form} />
          <SkillsList form={form} />
        </form>
      </Form>
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
        <FormField control={form.control} name="first_name" render={({ field }) => (
          <FormItem><FormLabel>Имя</FormLabel><FormControl><Input {...field} /></FormControl><FormMessage /></FormItem>
        )} />
        <FormField control={form.control} name="last_name" render={({ field }) => (
          <FormItem><FormLabel>Фамилия</FormLabel><FormControl><Input {...field} /></FormControl><FormMessage /></FormItem>
        )} />
        <FormField control={form.control} name="email" render={({ field }) => (
          <FormItem><FormLabel>Email</FormLabel><FormControl><Input type="email" {...field} /></FormControl><FormMessage /></FormItem>
        )} />
        <FormField control={form.control} name="phone" render={({ field }) => (
          <FormItem><FormLabel>Телефон</FormLabel><FormControl><Input {...field} /></FormControl><FormMessage /></FormItem>
        )} />
        <FormField control={form.control} name="city" render={({ field }) => (
          <FormItem><FormLabel>Город</FormLabel><FormControl><Input {...field} /></FormControl><FormMessage /></FormItem>
        )} />
        <FormField control={form.control} name="english_level" render={({ field }) => (
          <FormItem><FormLabel>Английский</FormLabel><FormControl><Input {...field} /></FormControl><FormMessage /></FormItem>
        )} />
      </div>
      <FormField control={form.control} name="portfolio_link" render={({ field }) => (
        <FormItem><FormLabel>Портфолио</FormLabel><FormControl><Input type="url" {...field} /></FormControl><FormMessage /></FormItem>
      )} />
      <FormField control={form.control} name="additional_info" render={({ field }) => (
        <FormItem><FormLabel>Дополнительно</FormLabel><FormControl><Textarea rows={3} {...field} /></FormControl><FormMessage /></FormItem>
      )} />
    </div>
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
            <FormItem><FormLabel>С</FormLabel><FormControl><Input type="number" {...field} value={field.value ?? ''} /></FormControl></FormItem>
          )} />
          <FormField control={form.control} name={`education.${i}.end_year`} render={({ field }) => (
            <FormItem><FormLabel>По</FormLabel><FormControl><Input type="number" {...field} value={field.value ?? ''} /></FormControl></FormItem>
          )} />
          <Button type="button" variant="ghost" onClick={() => arr.remove(i)} className="col-span-2 justify-self-start">
            Удалить
          </Button>
        </div>
      ))}
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
          <FormField control={form.control} name={`work_experiences.${i}.period`} render={({ field }) => (
            <FormItem><FormLabel>Период</FormLabel><FormControl><Input {...field} placeholder="2023-2024" /></FormControl></FormItem>
          )} />
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

function SkillsList({ form }: { form: UseFormReturn<ResumeFormData> }) {
  const skills = form.watch('skills_list') ?? []
  function setAt(i: number, v: string) {
    const next = [...skills]
    next[i] = v
    form.setValue('skills_list', next, { shouldDirty: true })
  }
  function append() {
    form.setValue('skills_list', [...skills, ''], { shouldDirty: true })
  }
  function remove(i: number) {
    form.setValue('skills_list', skills.filter((_, idx) => idx !== i), { shouldDirty: true })
  }
  return (
    <div className="space-y-3 rounded-lg border bg-card p-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">Навыки</h2>
        <Button type="button" variant="outline" onClick={append}>
          Добавить навык
        </Button>
      </div>
      {skills.map((value, i) => (
        <div key={i} className="flex gap-2 border-t pt-2">
          <SkillCombobox value={value} onChange={(v) => setAt(i, v)} />
          <Button type="button" variant="ghost" onClick={() => remove(i)}>
            Удалить
          </Button>
        </div>
      ))}
    </div>
  )
}
