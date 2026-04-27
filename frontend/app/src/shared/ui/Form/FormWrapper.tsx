import { zodResolver } from '@hookform/resolvers/zod'
import { useForm, type UseFormProps, type DefaultValues, type FieldValues } from 'react-hook-form'
import { type ZodType, type z } from 'zod'
import { Form } from '@/shared/ui/form'
import { type ReactNode } from 'react'

// reason: ZodType uses 3 generics whose variance doesn't line up with FieldValues — see https://github.com/react-hook-form/resolvers/issues/768
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
    // reason: zodResolver returns Resolver<z.infer<T>, any, any> but RHF's UseFormProps wants Resolver<TFieldValues, any, TFieldValues> — generic chain breaks
    resolver: zodResolver(schema) as any,
    // reason: DefaultValues<z.infer<T>> is structurally fine but RHF rejects optional fields without `as any`
    defaultValues: defaultValues as any,
    ...formProps,
  })
  return (
    <Form {...form}>
      {/* reason: handleSubmit signature mismatch with z.infer<T> when T contains transforms/refinements */}
      <form onSubmit={form.handleSubmit(onSubmit as any)} noValidate>
        {typeof children === 'function'
          // reason: children render-prop receives the typed form; cast bridges the same generic mismatch
          ? children(form as any)
          : children}
      </form>
    </Form>
  )
}
