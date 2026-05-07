import { fetchCompaniesMe } from '@/shared/api/companies/companiesMe'
import { useSessionStore } from './sessionStore'
import { writeActiveCompanyId, type CompanyMembership } from './types'

interface RefreshOptions {
  setActiveId?: string | undefined
}

export async function refreshCompanies(options?: RefreshOptions): Promise<void> {
  const stateBefore = useSessionStore.getState()
  const { user, setCompanies, setActiveCompany, activeCompanyId } = stateBefore
  const { data } = await fetchCompaniesMe({ limit: 100 })

  const fromServer: CompanyMembership[] = data.map((srv) => {
    const local = stateBefore.companies.find((l) => l.id === srv.id)
    return {
      ...srv,
      role: srv.role ?? local?.role,
    }
  })
  const localOptimistic = stateBefore.companies.filter(
    (local) => !data.some((srv) => srv.id === local.id),
  )
  const merged = [...fromServer, ...localOptimistic]
  setCompanies(merged)

  if (!user) return

  if (options && 'setActiveId' in options) {
    const next = options.setActiveId
    if (next && merged.some((c) => c.id === next)) {
      setActiveCompany(next)
      writeActiveCompanyId(user.id, next)
    }
    // Do NOT clear active company when requested id is missing - the optimistic
    // entry is preserved above and the server will catch up on next refresh.
    return
  }

  // Do not clear active id when server omits it: /companies/me has a 20s
  // cache that can return stale empty list after POST /companies.
  void user
  void activeCompanyId
}

export function addLocalCompany(membership: CompanyMembership, makeActive: boolean): void {
  const { user, companies, setCompanies, setActiveCompany } = useSessionStore.getState()
  if (!companies.some((c) => c.id === membership.id)) {
    setCompanies([...companies, membership])
  }
  if (makeActive) {
    setActiveCompany(membership.id)
    if (user) writeActiveCompanyId(user.id, membership.id)
  }
}
