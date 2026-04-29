import { redirect, type LoaderFunctionArgs } from 'react-router'
import { useSessionStore } from '@/shared/session/sessionStore'
import { readActiveCompanyId, writeActiveCompanyId } from '@/shared/session/types'
import { refreshCompanies } from '@/shared/session/refreshCompanies'

function loginRedirect(request: Request): Response {
  const url = new URL(request.url)
  return redirect(`/login?next=${encodeURIComponent(url.pathname + url.search)}`)
}

export function requireAuth({ request }: LoaderFunctionArgs) {
  const { status } = useSessionStore.getState()
  if (status !== 'authed') return loginRedirect(request)
  return null
}

export function requireRole(role: 'Candidate' | 'Company') {
  return function roleLoader({ request }: LoaderFunctionArgs) {
    const { status, user } = useSessionStore.getState()
    if (status !== 'authed') return loginRedirect(request)
    if (user?.role !== role) return redirect('/403')
    return null
  }
}

export async function requireCompanyMember({ request, params }: LoaderFunctionArgs) {
  const initial = useSessionStore.getState()
  if (initial.status !== 'authed') return loginRedirect(request)
  if (initial.user?.role !== 'Company') return redirect('/403')
  const companyId = params.companyId
  if (!companyId || companyId === 'me') return redirect('/company')

  let { companies, user } = useSessionStore.getState()
  if (!companies.some((c) => c.id === companyId)) {
    try {
      await refreshCompanies()
    } catch {
      void 0
    }
    ;({ companies, user } = useSessionStore.getState())
  }
  const membership = companies.find((c) => c.id === companyId)
  if (!membership) {
    if (user && readActiveCompanyId(user.id) === companyId) return null
    return redirect('/403')
  }
  if (user) {
    useSessionStore.getState().setActiveCompany(companyId)
    writeActiveCompanyId(user.id, companyId)
  }
  return null
}

export async function requireCompanyAdmin(args: LoaderFunctionArgs) {
  const memberCheck = await requireCompanyMember(args)
  if (memberCheck) return memberCheck
  const { companies } = useSessionStore.getState()
  const membership = companies.find((c) => c.id === args.params.companyId)
  if (membership?.role !== 'admin') return redirect('/403')
  return null
}

export async function resolveActiveCompany({ request }: LoaderFunctionArgs) {
  const initial = useSessionStore.getState()
  if (initial.status !== 'authed') return loginRedirect(request)
  if (initial.user?.role !== 'Company') return redirect('/403')

  let { companies, activeCompanyId } = useSessionStore.getState()
  const { user } = useSessionStore.getState()
  const stored = user ? readActiveCompanyId(user.id) : undefined
  if (stored) return redirect(`/company/${stored}/dashboard`)
  if (companies.length === 0) {
    try {
      await refreshCompanies()
    } catch {
      void 0
    }
    ;({ companies, activeCompanyId } = useSessionStore.getState())
  }
  if (companies.length === 0) return redirect('/company/new')
  const target =
    activeCompanyId && companies.some((c) => c.id === activeCompanyId)
      ? activeCompanyId
      : companies[0]!.id
  return redirect(`/company/${target}/dashboard`)
}

export function redirectIfAuth() {
  const { status, user, companies, activeCompanyId } = useSessionStore.getState()
  if (status !== 'authed' || !user) return null
  if (user.role === 'Candidate') return redirect('/me/profile')
  if (companies.length === 0) return redirect('/company/new')
  const target =
    activeCompanyId && companies.some((c) => c.id === activeCompanyId)
      ? activeCompanyId
      : companies[0]!.id
  return redirect(`/company/${target}/dashboard`)
}
