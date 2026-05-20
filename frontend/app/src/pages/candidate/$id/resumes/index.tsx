import { useParams, Link } from 'react-router'

interface MockResume {
  id: string
  name: string
  status: string
}

const MOCK_RESUMES: MockResume[] = [
  { id: '1', name: 'Junior Developer', status: 'published' },
  { id: '2', name: 'Intern Position', status: 'draft' },
]

export default function CandidateResumesPage() {
  const { id } = useParams<{ id: string }>()

  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <div className="flex items-center gap-2">
        <Link
          to={`/candidate/${id}`}
          className="text-sm text-muted-foreground underline"
        >
          ← Профиль кандидата
        </Link>
      </div>

      <h1 className="text-2xl font-bold">Резюме кандидата #{id}</h1>

      <ul className="space-y-2">
        {MOCK_RESUMES.map((resume) => (
          <li key={resume.id} className="rounded-lg border bg-card p-4">
            <Link
              to={`/candidate/${id}/resumes/${resume.id}`}
              className="text-lg font-medium text-primary underline"
            >
              {resume.name}
            </Link>
            <p className="text-sm text-muted-foreground">Статус: {resume.status}</p>
          </li>
        ))}
      </ul>
    </div>
  )
}