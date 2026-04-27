import { Button } from '@/shared/ui/button'
import { Skeleton } from '@/shared/ui/skeleton'

interface Props {
  nextCursor?: string | null
  onNext: () => void
  isLoading?: boolean
}

export function CursorPagination({ nextCursor, onNext, isLoading }: Props) {
  if (isLoading) return <Skeleton className="h-9 w-24" />
  if (!nextCursor) return null
  return <Button variant="outline" onClick={onNext}>Load more</Button>
}
