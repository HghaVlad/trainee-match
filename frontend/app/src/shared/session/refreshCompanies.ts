import { fetchCompaniesMe } from '@/shared/api/companies/companiesMe'
import { useSessionStore } from './sessionStore'
import { writeActiveCompanyId } from './types'

interface RefreshOptions {
  setActiveId?: string | undefined
}

export async function refreshCompanies(options?: RefreshOptions): Promise<void> {
  const { user, setCompanies, setActiveCompany, activeCompanyId } =
    useSessionStore.getState()
  const { data } = await fetchCompaniesMe({ limit: 100 })
  setCompanies(data)

  if (!user) return

  if (options && 'setActiveId' in options) {
    const next = options.setActiveId
    if (next && data.some((c) => c.id === next)) {
      setActiveCompany(next)
      writeActiveCompanyId(user.id, next)
    } else {
      setActiveCompany(undefined)
      writeActiveCompanyId(user.id, undefined)
    }
    return
  }

  if (activeCompanyId && !data.some((c) => c.id === activeCompanyId)) {
    setActiveCompany(undefined)
    writeActiveCompanyId(user.id, undefined)
  }
}
