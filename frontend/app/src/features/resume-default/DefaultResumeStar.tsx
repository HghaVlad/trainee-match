import { Star } from 'lucide-react'
import { cn } from '@/shared/lib/utils'
import { useDefaultResumeId } from './useDefaultResumeId'

interface Props {
  resumeId: string
  isPublished: boolean
  className?: string
}

export function DefaultResumeStar({ resumeId, isPublished, className }: Props) {
  const { defaultResumeId, setDefaultResumeId } = useDefaultResumeId()
  const isDefault = defaultResumeId === resumeId
  const disabled = !isPublished
  const title = disabled
    ? 'Опубликуйте резюме, чтобы сделать его основным'
    : isDefault
      ? 'Основное резюме'
      : 'Сделать основным'

  function handleClick() {
    if (disabled) return
    setDefaultResumeId(isDefault ? undefined : resumeId)
  }

  return (
    <button
      type="button"
      onClick={handleClick}
      disabled={disabled}
      title={title}
      aria-label={title}
      aria-pressed={isDefault}
      className={cn(
        'inline-flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:text-yellow-500 disabled:cursor-not-allowed disabled:opacity-40',
        isDefault && 'text-yellow-500',
        className,
      )}
    >
      <Star
        className="h-5 w-5"
        fill={isDefault ? 'currentColor' : 'none'}
        strokeWidth={1.75}
      />
    </button>
  )
}
