import { type ReactNode } from 'react'
import { Button } from '@/shared/ui/button'

interface Props {
  title: string
  description?: string
  icon?: ReactNode
  action?: {
    label: string
    onClick: () => void
  }
}

export function EmptyState({ title, description, icon, action }: Props) {
  return (
    <div className="flex h-[450px] shrink-0 items-center justify-center rounded-md border border-dashed">
      <div className="mx-auto flex max-w-[420px] flex-col items-center justify-center text-center">
        {icon && <div className="flex h-20 w-20 items-center justify-center rounded-full bg-muted">{icon}</div>}
        <h3 className="mt-4 text-lg font-semibold">{title}</h3>
        {description && <p className="mb-4 mt-2 text-sm text-muted-foreground">{description}</p>}
        {action && (
          <Button onClick={action.onClick} className="mt-4">
            {action.label}
          </Button>
        )}
      </div>
    </div>
  )
}
