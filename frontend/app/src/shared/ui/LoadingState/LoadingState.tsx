import { Skeleton } from '@/shared/ui/skeleton'

export function LoadingState() {
  return (
    <div className="flex h-[450px] w-full flex-col items-center justify-center space-y-3 rounded-md border p-8">
      <Skeleton className="h-10 w-[200px]" />
      <Skeleton className="h-4 w-[300px]" />
      <Skeleton className="h-4 w-[250px]" />
      <div className="mt-8 flex w-full flex-col space-y-3">
        <Skeleton className="h-12 w-full" />
        <Skeleton className="h-12 w-full" />
        <Skeleton className="h-12 w-full" />
      </div>
    </div>
  )
}
