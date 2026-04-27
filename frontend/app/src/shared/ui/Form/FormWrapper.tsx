import { zodResolver } from '@hookform/resolvers/zod'
import { useForm, type UseFormProps, type DefaultValues, type FieldValues } from 'react-hook-form'
import { type ZodType, type z } from 'zod'
import { Form } from '@/shared/ui/form'
import { type ReactNode } from 'react'

export interface FormWrapperProps<T extends ZodType<any, any, any>> {
  schema: T
  defaultValues?: DefaultValues<z.infer<T>>
  onSubmit: (data: z.infer<T>) => void | Promise<void>
  children: ReactNode | ((form: ReturnType<typeof useForm<z.infer<T>>>) => ReactNode)
  formProps?: Omit<UseFormProps<z.infer<T>>, 'resolver' | 'defaultValues'>
}

export function FormWrapper<T extends ZodType<any, any, any>>({
  schema, defaultValues, onSubmit, children, formProps
}: FormWrapperProps<T>) {
  const form = useForm<z.infer<T>>({
    resolver: zodResolver(schema) as any,
    defaultValues: defaultValues as any,
    ...formProps,
  })
  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit as any)} noValidate>
        {typeof children === 'function' ? children(form as any) : children}
      </form>
    </Form>
  )
}
