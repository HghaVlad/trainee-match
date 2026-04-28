import { useNavigate, useParams, useLocation } from 'react-router'
import { useSession } from '@/shared/session/useSession'
import { useSessionStore } from '@/shared/session/sessionStore'
import { writeActiveCompanyId } from '@/shared/session/types'

export function CompanySwitcher() {
  const navigate = useNavigate()
  const params = useParams()
  const location = useLocation()
  const { user, companies, activeCompanyId } = useSession()

  if (!user || user.role !== 'Company' || companies.length === 0) return null

  function handleChange(e: React.ChangeEvent<HTMLSelectElement>) {
    const nextId = e.target.value
    if (!nextId || !user) return
    useSessionStore.getState().setActiveCompany(nextId)
    writeActiveCompanyId(user.id, nextId)
    const currentId = params.companyId
    if (currentId && location.pathname.startsWith(`/company/${currentId}`)) {
      const rest = location.pathname.slice(`/company/${currentId}`.length)
      navigate(`/company/${nextId}${rest || '/dashboard'}`)
    } else {
      navigate(`/company/${nextId}/dashboard`)
    }
  }

  return (
    <select
      aria-label="Active company"
      value={activeCompanyId ?? ''}
      onChange={handleChange}
      style={{ padding: '0.25rem 0.5rem' }}
    >
      {companies.map((c) => (
        <option key={c.id} value={c.id}>
          {c.name}
          {c.role ? ` · ${c.role}` : ''}
        </option>
      ))}
    </select>
  )
}
