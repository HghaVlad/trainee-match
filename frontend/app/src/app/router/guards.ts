import { redirect, type LoaderFunctionArgs } from 'react-router'
import { useSessionStore } from '@/shared/session/sessionStore'
import { writeActiveCompanyId } from '@/shared/session/types'

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

export function requireCompanyMember({ request, params }: LoaderFunctionArgs) {
  const { status, user, companies } = useSessionStore.getState()
  if (status !== 'authed') return loginRedirect(request)
  if (user?.role !== 'Company') return redirect('/403')
  const companyId = params.companyId
  if (!companyId) return redirect('/company')
  const membership = companies.find((c) => c.id === companyId)
  if (!membership) return redirect('/403')
  if (user) {
    useSessionStore.getState().setActiveCompany(companyId)
    writeActiveCompanyId(user.id, companyId)
  }
  return null
}

export function requireCompanyAdmin(args: LoaderFunctionArgs) {
  const memberCheck = requireCompanyMember(args)
  if (memberCheck) return memberCheck
  const { companies } = useSessionStore.getState()
  const membership = companies.find((c) => c.id === args.params.companyId)
  if (membership?.role !== 'admin') return redirect('/403')
  return null
}

export function resolveActiveCompany({ request }: LoaderFunctionArgs) {
  const { status, user, companies, activeCompanyId } = useSessionStore.getState()
  if (status !== 'authed') return loginRedirect(request)
  if (user?.role !== 'Company') return redirect('/403')
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
