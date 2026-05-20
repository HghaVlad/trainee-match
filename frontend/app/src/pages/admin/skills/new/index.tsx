import { useState } from 'react'
import { useNavigate } from 'react-router'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card'

export default function AdminSkillNewPage() {
  const [skillName, setSkillName] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)
  const navigate = useNavigate()

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!skillName.trim()) return

    setIsLoading(true)
    setError(null)
    setSuccess(false)

    try {
      await new Promise((resolve) => setTimeout(resolve, 500))
      setSuccess(true)
      setSkillName('')
    } catch {
      setError('Не удалось создать навык')
    } finally {
      setIsLoading(false)
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
                disabled={isLoading}
              />
            </div>

            {error && <p className="text-sm text-destructive">{error}</p>}
            {success && (
              <p className="text-sm text-green-600">Навык успешно создан!</p>
            )}

            <div className="flex gap-2">
              <Button type="submit" disabled={isLoading || !skillName.trim()}>
                {isLoading ? 'Создание...' : 'Создать'}
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}