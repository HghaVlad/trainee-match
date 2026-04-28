import { useQuery } from '@tanstack/react-query'
import type { UseQueryOptions } from '@tanstack/react-query'
import { httpClient } from '@/shared/api/http/client'
import type { CompanyRole } from '@/shared/session/types'

export interface CompanyMember {
  userId: string
  username: string
  role: CompanyRole
}

export interface CompanyMembersPage {
  data: CompanyMember[]
  nextCursor: string | null
  hasNext: boolean
}

export interface FetchCompanyMembersParams {
  cursor?: string
  limit?: number
}

interface RawCompanyMember {
  userId?: string
  user_id?: string
  userID?: string
  id?: string
  username?: string
  userName?: string
  user_name?: string
  login?: string
  role?: CompanyRole
}

interface RawCompanyMembersResponse {
  data?: RawCompanyMember[]
  items?: RawCompanyMember[]
  members?: RawCompanyMember[]
  content?: RawCompanyMember[]
  nextCursor?: string | null
  next_cursor?: string | null
  hasNext?: boolean
  has_next?: boolean
}

function normalizeMember(raw: RawCompanyMember): CompanyMember {
  const userId = raw.userId ?? raw.user_id ?? raw.userID ?? raw.id ?? ''
  const username = raw.username ?? raw.userName ?? raw.user_name ?? raw.login ?? ''
  const role: CompanyRole = raw.role === 'admin' ? 'admin' : 'recruiter'
  return { userId, username, role }
}

export async function fetchCompanyMembers(
  companyId: string,
  params?: FetchCompanyMembersParams,
): Promise<CompanyMembersPage> {
  const { data } = await httpClient.get<RawCompanyMembersResponse>(
    `/companies/${companyId}/members`,
    {
      params: {
        cursor: params?.cursor,
        limit: params?.limit,
      },
    },
  )
  const list = data.data ?? data.items ?? data.members ?? data.content ?? []
  return {
    data: list.map(normalizeMember),
    nextCursor: data.nextCursor ?? data.next_cursor ?? null,
    hasNext: data.hasNext ?? data.has_next ?? false,
  }
}

export function getCompanyMembersQueryKey(
  companyId: string,
  params?: FetchCompanyMembersParams,
): readonly unknown[] {
  return ['company-members', companyId, params ?? {}] as const
}

export function useCompanyMembersQuery(
  companyId: string,
  params?: FetchCompanyMembersParams,
  options?: Omit<
    UseQueryOptions<CompanyMembersPage, Error>,
    'queryKey' | 'queryFn'
  >,
) {
  return useQuery<CompanyMembersPage, Error>({
    queryKey: getCompanyMembersQueryKey(companyId, params),
    queryFn: () => fetchCompanyMembers(companyId, params),
    enabled: Boolean(companyId),
    ...options,
  })
}
