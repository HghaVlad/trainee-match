import { useState } from 'react'
import { Link } from 'react-router'
import { useQueryClient } from '@tanstack/react-query'
import { useGetSkillList } from '@/api/generated/candidate/skill/skill'
import { useDeleteAdminSkillsId } from '@/api/generated/candidate/admin/admin'
import { Button } from '@/shared/ui/button'
import { LoadingState } from '@/shared/ui/LoadingState'
import { ErrorState } from '@/shared/ui/ErrorState'
import { EmptyState } from '@/shared/ui/EmptyState'
import { useToast } from '@/shared/hooks/use-toast'
import { AppError } from '@/shared/api/http/client'

const PAGE_SIZE = 20

export default function AdminSkillsPage() {
  const [page, setPage] = useState(1)
  const [removedIds, setRemovedIds] = useState<Set<string>>(new Set())
  const qc = useQueryClient()
  const { toast } = useToast()
  const deleteSkill = useDeleteAdminSkillsId()
  const { data, isLoading, error, refetch, isFetching } = useGetSkillList({ page, size: PAGE_SIZE })

  const visibleSkills = data?.filter((s) => s.id && !removedIds.has(s.id)) ?? []

  async function onDelete(skillId: string, skillName?: string) {
    setRemovedIds((prev) => new Set(prev).add(skillId))
    try {
      await deleteSkill.mutateAsync({ id: skillId })
      await qc.invalidateQueries({ queryKey: ['useGetSkillList'] })
      toast({ title: `Удалён: ${skillName ?? skillId}` })
    } catch (e) {
      setRemovedIds((prev) => {
        const next = new Set(prev)
        next.delete(skillId)
        return next
      })
      const msg = e instanceof AppError ? e.message : 'Не удалось удалить навык'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  if (isLoading) return <LoadingState />
  if (error) return <ErrorState onRetry={() => refetch()} />

  const hasMore = data != null && data.length >= PAGE_SIZE

  if (visibleSkills.length === 0 && !data?.length) {
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
        {visibleSkills.map((skill) => (
          <li
            key={skill.id}
            className="rounded-lg border bg-card p-4 flex items-center justify-between"
          >
            <div>
              <span className="font-medium">{skill.name}</span>
              <span className="text-sm text-muted-foreground ml-2">{skill.id}</span>
            </div>
            {skill.id && (
              <Button
                variant="destructive"
                size="sm"
                onClick={() => onDelete(skill.id as string, skill.name)}
                disabled={deleteSkill.isPending}
              >
                {deleteSkill.isPending ? '...' : 'Удалить'}
              </Button>
            )}
          </li>
        ))}
      </ul>

      <div className="flex items-center justify-center gap-4">
        <Button
          variant="outline"
          size="sm"
          onClick={() => setPage((p) => Math.max(1, p - 1))}
          disabled={page <= 1 || isFetching}
        >
          Назад
        </Button>
        <span className="text-sm text-muted-foreground">
          Страница {page}
        </span>
        <Button
          variant="outline"
          size="sm"
          onClick={() => setPage((p) => p + 1)}
          disabled={!hasMore || isFetching}
        >
          {isFetching ? 'Загрузка…' : 'Вперёд'}
        </Button>
      </div>
    </div>
  )
}