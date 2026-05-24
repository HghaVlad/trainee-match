export interface CompanyMembership {
  id: string
  name: string
  logoKey?: string
  openVacanciesCount: number
  createdAt: string
  role?: 'admin' | 'recruiter'
}

export interface CompaniesMeResponse {
  data: CompanyMembership[]
  nextCursor?: string | null
  hasNext?: boolean
}

export type CompanyRole = 'admin' | 'recruiter'

const ACTIVE_COMPANY_KEY_PREFIX = 'tm.activeCompanyId.'

export function activeCompanyKey(userId: string | number): string {
  return `${ACTIVE_COMPANY_KEY_PREFIX}${userId}`
}

export function readActiveCompanyId(userId: string | number): string | undefined {
  if (typeof window === 'undefined') return undefined
  try {
    return window.localStorage.getItem(activeCompanyKey(userId)) ?? undefined
  } catch {
    return undefined
  }
}

export function writeActiveCompanyId(
  userId: string | number,
  companyId: string | undefined,
): void {
  if (typeof window === 'undefined') return
  try {
    const key = activeCompanyKey(userId)
    if (companyId) {
      window.localStorage.setItem(key, companyId)
    } else {
      window.localStorage.removeItem(key)
    }
  } catch {
    return
  }
}

const COMPANY_ROLES_KEY_PREFIX = 'tm.companyRoles.'

function companyRolesKey(userId: string | number): string {
  return `${COMPANY_ROLES_KEY_PREFIX}${userId}`
}

interface StoredRole {
  companyId: string
  role: 'admin' | 'recruiter'
}

export function saveCompanyRoles(
  userId: string | number,
  roles: StoredRole[],
): void {
  if (typeof window === 'undefined') return
  try {
    window.localStorage.setItem(companyRolesKey(userId), JSON.stringify(roles))
  } catch {
    // quota exceeded, ignore
  }
}

export function readCompanyRoles(
  userId: string | number,
): StoredRole[] {
  if (typeof window === 'undefined') return []
  try {
    const raw = window.localStorage.getItem(companyRolesKey(userId))
    if (!raw) return []
    const parsed = JSON.parse(raw)
    if (!Array.isArray(parsed)) return []
    return parsed.filter(
      (r): r is StoredRole =>
        typeof r === 'object' &&
        r !== null &&
        typeof r.companyId === 'string' &&
        (r.role === 'admin' || r.role === 'recruiter'),
    )
  } catch {
    return []
  }
}

export function clearCompanyRoles(userId: string | number): void {
  if (typeof window === 'undefined') return
  try {
    window.localStorage.removeItem(companyRolesKey(userId))
  } catch {
    return
  }
}
