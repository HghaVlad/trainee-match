import { httpClient } from '@/shared/api/http/client'
import type { CompaniesMeResponse, CompanyMembership } from '@/shared/session/types'

export interface FetchCompaniesMeParams {
  cursor?: string
  limit?: number
}

interface RawCompaniesMeItem {
  id: string
  name: string
  logoKey?: string
  logo_key?: string
  openVacanciesCount?: number
  open_vacancies_count?: number
  createdAt?: string
  created_at?: string
  role?: 'admin' | 'recruiter'
}

interface RawCompaniesMeResponse {
  data?: RawCompaniesMeItem[]
  items?: RawCompaniesMeItem[]
  nextCursor?: string | null
  next_cursor?: string | null
  hasNext?: boolean
  has_next?: boolean
}

function normalizeItem(raw: RawCompaniesMeItem): CompanyMembership {
  return {
    id: raw.id,
    name: raw.name,
    logoKey: raw.logoKey ?? raw.logo_key,
    openVacanciesCount: raw.openVacanciesCount ?? raw.open_vacancies_count ?? 0,
    createdAt: raw.createdAt ?? raw.created_at ?? new Date(0).toISOString(),
    role: raw.role,
  }
}

export async function fetchCompaniesMe(
  params?: FetchCompaniesMeParams,
): Promise<CompaniesMeResponse> {
  const { data } = await httpClient.get<RawCompaniesMeResponse>('/companies/me', {
    params: {
      cursor: params?.cursor,
      limit: params?.limit,
    },
  })
  const items = (data.data ?? data.items ?? []).map(normalizeItem)
  return {
    data: items,
    nextCursor: data.nextCursor ?? data.next_cursor ?? null,
    hasNext: data.hasNext ?? data.has_next ?? false,
  }
}
