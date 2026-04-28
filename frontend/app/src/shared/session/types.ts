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
