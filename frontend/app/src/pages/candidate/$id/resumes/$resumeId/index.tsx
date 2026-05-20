import { useState } from 'react'
import { useParams, Link } from 'react-router'
import { Button } from '@/shared/ui/button'

export default function CandidateResumeDetailPage() {
  const { id, resumeId } = useParams<{ id: string; resumeId: string }>()
  const [archiving, setArchiving] = useState(false)

  async function onArchive() {
    setArchiving(true)
    await new Promise((resolve) => setTimeout(resolve, 500))
    setArchiving(false)
  }

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <div className="flex items-center gap-2">
        <Link
          to={`/candidate/${id}/resumes`}
          className="text-sm text-muted-foreground underline"
        >
          ← Все резюме
        </Link>
      </div>

      <h1 className="text-2xl font-bold">Резюме #{resumeId}</h1>

      <div className="rounded-lg border bg-card p-4 space-y-2">
        <p><strong>Название:</strong> Mock Resume</p>
        <p><strong>Статус:</strong> published</p>
        <p><strong>Кандидат:</strong> {id}</p>
      </div>

      <div className="flex gap-2">
        <Button variant="destructive" onClick={onArchive} disabled={archiving}>
          {archiving ? 'Архивация...' : 'Архивировать'}
        </Button>
      </div>
    </div>
  )
}