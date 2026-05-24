import { fetchCompaniesMe } from '@/shared/api/companies/companiesMe'
import { httpClient } from '@/shared/api/http/client'
import { useSessionStore } from './sessionStore'
import {
  readActiveCompanyId,
  writeActiveCompanyId,
  saveCompanyRoles,
  readCompanyRoles,
  type CompanyMembership,
} from './types'

interface MembersMeResponse {
  role?: 'admin' | 'recruiter'
  userId?: string
  username?: string
  email?: string
  companyId?: string
}

interface RefreshOptions {
  setActiveId?: string | undefined
}

export async function refreshCompanies(options?: RefreshOptions): Promise<void> {
  const stateBefore = useSessionStore.getState()
  const { user, setCompanies, setActiveCompany, activeCompanyId } = stateBefore

  const storedRoles = user ? readCompanyRoles(user.id) : []

  const { data } = await fetchCompaniesMe({ limit: 100 })

  const fromServer: CompanyMembership[] = await Promise.all(
    data.map(async (srv) => {
      const local = stateBefore.companies.find((l) => l.id === srv.id)
      const role = srv.role ?? local?.role
      if (role) {
        return { ...srv, role }
      }
      // role missing — try /companies/{id}/members/me
      try {
        const { data: me } = await httpClient.get<MembersMeResponse>(
          `/companies/${srv.id}/members/me`,
        )
        const resolvedRole = me.role ?? storedRoles.find((r) => r.companyId === srv.id)?.role
        return { ...srv, role: resolvedRole }
      } catch {
        // fallback to localStorage if API fails (e.g. 404 after session expiry)
        const storedRole = storedRoles.find((r) => r.companyId === srv.id)?.role
        return { ...srv, role: storedRole }
      }
    }),
  )
  const localOptimistic = stateBefore.companies.filter(
    (local) => !data.some((srv) => srv.id === local.id),
  )
  const merged = [...fromServer, ...localOptimistic]

  if (merged.length === 0 && user) {
    const storedId = readActiveCompanyId(user.id)
    if (storedId) {
      try {
        const [profileRes, memberRes] = await Promise.all([
          httpClient.get<{
            id: string
            name: string
            openVacanciesCount?: number
            createdAt?: string
          }>(`/companies/${storedId}`),
          httpClient.get<{ role?: 'admin' | 'recruiter' }>(
            `/companies/${storedId}/members/me`,
          ),
        ])
        const p = profileRes.data
        const resolvedRole = memberRes.data.role ?? storedRoles.find((r) => r.companyId === storedId)?.role
        merged.push({
          id: p.id,
          name: p.name,
          openVacanciesCount: p.openVacanciesCount ?? 0,
          createdAt: p.createdAt ?? new Date(0).toISOString(),
          role: resolvedRole,
        })
      } catch {
        void 0
      }
    }
  }

  setCompanies(merged)

  // persist roles to localStorage so they survive session expiry
  if (user) {
    const roles = merged
      .filter((c): c is CompanyMembership & { role: 'admin' | 'recruiter' } => c.role !== undefined)
      .map((c) => ({ companyId: c.id, role: c.role }))
    if (roles.length > 0) {
      saveCompanyRoles(user.id, roles)
    }
  }

  if (!user) return

  if (options && 'setActiveId' in options) {
    const next = options.setActiveId
    if (next && merged.some((c) => c.id === next && c.name)) {
      setActiveCompany(next)
      writeActiveCompanyId(user.id, next)
    } else if (!merged.some((c) => c.id === next) && next) {
      setActiveCompany(undefined)
      writeActiveCompanyId(user.id, undefined)
    }
    return
  }

  if (activeCompanyId && !merged.some((c) => c.id === activeCompanyId)) {
    setActiveCompany(undefined)
    writeActiveCompanyId(user!.id, undefined)
  }
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
