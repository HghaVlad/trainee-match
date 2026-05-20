import { useParams, Link } from 'react-router'

export default function CandidateDetailPage() {
  const { id } = useParams<{ id: string }>()

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <div className="flex items-center gap-2">
        <Link to="/candidates" className="text-sm text-muted-foreground underline">
          ← Все кандидаты
        </Link>
      </div>

      <h1 className="text-2xl font-bold">Кандидат #{id}</h1>

      <div className="rounded-lg border bg-card p-4 space-y-2">
        <p><strong>ID:</strong> {id}</p>
        <p><strong>Имя:</strong> Mock Name</p>
        <p><strong>Город:</strong> Mock City</p>
      </div>

      <div className="rounded-lg border bg-card p-4">
        <Link
          to={`/candidate/${id}/resumes`}
          className="text-primary underline"
        >
          Посмотреть резюме →
        </Link>
      </div>
    </div>
  )
}