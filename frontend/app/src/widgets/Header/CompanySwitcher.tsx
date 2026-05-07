import { useNavigate, useParams, useLocation } from 'react-router'
import { useSession } from '@/shared/session/useSession'
import { useSessionStore } from '@/shared/session/sessionStore'
import { writeActiveCompanyId } from '@/shared/session/types'

const NEW_COMPANY_OPTION = '__new_company__'

export function CompanySwitcher() {
  const navigate = useNavigate()
  const params = useParams()
  const location = useLocation()
  const { user, companies, activeCompanyId } = useSession()

  if (!user || user.role !== 'Company') return null

  function handleChange(e: React.ChangeEvent<HTMLSelectElement>) {
    const nextId = e.target.value
    if (!nextId || !user) return
    if (nextId === NEW_COMPANY_OPTION) {
      navigate('/company/new')
      return
    }
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

  const selectValue = activeCompanyId && companies.some((c) => c.id === activeCompanyId)
    ? activeCompanyId
    : ''

  return (
    <label className="flex items-center gap-2 text-sm">
      <span className="text-muted-foreground">Активная компания:</span>
      <select
        aria-label="Активная компания"
        value={selectValue}
        onChange={handleChange}
        className="h-9 rounded-md border border-input bg-background px-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
      >
        {companies.length === 0 && (
          <option value="" disabled>
            Нет компаний
          </option>
        )}
        {companies.map((c) => (
          <option key={c.id} value={c.id}>
            {c.name}
            {c.role ? ` · ${c.role}` : ''}
          </option>
        ))}
        <option disabled>──────────</option>
        <option value={NEW_COMPANY_OPTION}>+ Новая компания</option>
      </select>
    </label>
  )
}
