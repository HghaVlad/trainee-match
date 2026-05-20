import { Link } from 'react-router'

interface MockCandidate {
  id: string
  name: string
  city: string
  status: string
}

const MOCK_CANDIDATES: MockCandidate[] = [
  { id: '1', name: 'Иван Иванов', city: 'Москва', status: 'active' },
  { id: '2', name: 'Петр Петров', city: 'Санкт-Петербург', status: 'active' },
  { id: '3', name: 'Анна Сидорова', city: 'Москва', status: 'inactive' },
]

export default function CandidatesPage() {
  return (
    <div className="mx-auto max-w-3xl p-6 space-y-4">
      <h1 className="text-2xl font-bold">Кандидаты</h1>

      <ul className="space-y-2">
        {MOCK_CANDIDATES.map((candidate) => (
          <li key={candidate.id} className="rounded-lg border bg-card p-4">
            <Link
              to={`/candidate/${candidate.id}`}
              className="text-lg font-medium text-primary underline"
            >
              {candidate.name}
            </Link>
            <p className="text-sm text-muted-foreground">
              {candidate.city} • {candidate.status}
            </p>
          </li>
        ))}
      </ul>
    </div>
  )
}