import { Link } from 'react-router'
import { useGetSkillList } from '@/api/generated/candidate/skill/skill'
import { Button } from '@/shared/ui/button'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'

export default function AdminSkillsPage() {
  const { data, isLoading, error, refetch } = useGetSkillList()

  if (isLoading) return <LoadingState />
  if (error) return <ErrorState onRetry={() => refetch()} />
  if (!data || data.length === 0) {
    return (
      <div className="mx-auto max-w-3xl p-6 space-y-4">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold">Навыки</h1>
          <Button asChild>
            <Link to="/admin/skills/new">Создать навык</Link>
          </Button>
        </div>
        <EmptyState title="Навыки не найдены" />
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Навыки</h1>
        <Button asChild>
          <Link to="/admin/skills/new">Создать навык</Link>
        </Button>
      </div>

      <ul className="space-y-2">
        {data.map((skill) => (
          <li
            key={skill.id}
            className="rounded-lg border bg-card p-4 flex items-center justify-between"
          >
            <span className="font-medium">{skill.name}</span>
            <span className="text-sm text-muted-foreground">{skill.id}</span>
          </li>
        ))}
      </ul>
    </div>
  )
}