import { Button } from '@/shared/ui/button'

interface Props {
  title?: string
  message?: string
  onRetry?: () => void
}

export function ErrorState({ title = 'Something went wrong', message = 'An error occurred while loading the data.', onRetry }: Props) {
  return (
    <div className="flex h-[450px] shrink-0 items-center justify-center rounded-md border border-destructive bg-destructive/10">
      <div className="mx-auto flex max-w-[420px] flex-col items-center justify-center text-center">
        <h3 className="mt-4 text-lg font-semibold text-destructive">{title}</h3>
        <p className="mb-4 mt-2 text-sm text-destructive/80">{message}</p>
        {onRetry && (
          <Button variant="outline" onClick={onRetry} className="mt-4 border-destructive text-destructive hover:bg-destructive hover:text-destructive-foreground">
            Try again
          </Button>
        )}
      </div>
    </div>
  )
}
