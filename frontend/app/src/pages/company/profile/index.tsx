import { Navigate, useParams } from 'react-router'
import { CompanyProfile } from '@/features/company-profile'

export default function CompanyProfilePage() {
  const { companyId } = useParams<{ companyId: string }>()
  if (!companyId) return <Navigate to="/company" replace />
  return (
    <div className="mx-auto max-w-3xl p-6">
      <CompanyProfile companyId={companyId} />
    </div>
  )
}
