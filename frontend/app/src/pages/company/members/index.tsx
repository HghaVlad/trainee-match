import { Navigate, useParams } from 'react-router'
import { CompanyMembers } from '@/features/company-members'

export default function CompanyMembersPage() {
  const { companyId } = useParams<{ companyId: string }>()
  if (!companyId) return <Navigate to="/company" replace />
  return (
    <div className="mx-auto max-w-4xl p-6">
      <CompanyMembers companyId={companyId} />
    </div>
  )
}
