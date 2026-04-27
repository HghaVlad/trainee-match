import { zodResolver } from '@hookform/resolvers/zod'
import { useForm, type UseFormProps, type DefaultValues } from 'react-hook-form'
import { type ZodSchema, type z } from 'zod'
import { Form } from '@/shared/ui/form'
import { type ReactNode } from 'react'

interface FormWrapperProps<T extends ZodSchema> {
  schema: T
  defaultValues?: DefaultValues<z.infer<T>>
  onSubmit: (data: z.infer<T>) => void | Promise<void>
  children: ReactNode | ((form: ReturnType<typeof useForm<z.infer<T>>>) => ReactNode)
  formProps?: UseFormProps<z.infer<T>>
}

export function FormWrapper<T extends ZodSchema>({
  schema, defaultValues, onSubmit, children, formProps
}: FormWrapperProps<T>) {
  const form = useForm<z.infer<T>>({
    resolver: zodResolver(schema),
    defaultValues,
    ...formProps,
  })
  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} noValidate>
        {typeof children === 'function' ? children(form) : children}
      </form>
    </Form>
  )
}
