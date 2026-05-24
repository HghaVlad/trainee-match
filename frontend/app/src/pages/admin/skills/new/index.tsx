import { useState } from 'react'
import { useNavigate } from 'react-router'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'
import { useToast } from '@/shared/hooks/use-toast'
import { usePostAdminSkills } from '@/api/generated/candidate/admin/admin'
import { useQueryClient } from '@tanstack/react-query'
import { AppError } from '@/shared/api/http/client'

export default function AdminSkillNewPage() {
  const [skillName, setSkillName] = useState('')
  const navigate = useNavigate()
  const { toast } = useToast()
  const qc = useQueryClient()
  const createSkill = usePostAdminSkills()

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!skillName.trim()) return

    try {
      await createSkill.mutateAsync({ data: { name: skillName } })
      await qc.invalidateQueries({
        predicate: (query) => {
          const first = query.queryKey[0]
          return typeof first === 'string' && first === '/skill/list'
        },
        refetchType: 'all',
      })
      toast({ title: 'Навык создан' })
      navigate('/admin/skills')
    } catch (err) {
      const msg = err instanceof AppError ? err.message : 'Не удалось создать навык'
      toast({ title: 'Ошибка', description: msg, variant: 'destructive' })
    }
  }

  return (
    <div className="mx-auto max-w-md p-6 space-y-4">
      <div className="flex items-center gap-2">
        <Button variant="ghost" size="sm" onClick={() => navigate('/admin/skills')}>
          ← Назад
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Создание навыка</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="space-y-4">
            <div>
              <Input
                placeholder="Название навыка"
                value={skillName}
                onChange={(e) => setSkillName(e.target.value)}
                disabled={createSkill.isPending}
              />
            </div>

            <div className="flex gap-2">
              <Button type="submit" disabled={createSkill.isPending || !skillName.trim()}>
                {createSkill.isPending ? 'Создание...' : 'Создать'}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}