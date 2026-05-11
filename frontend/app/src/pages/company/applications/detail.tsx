import { useParams } from 'react-router'
import { ApplicationDetail } from '@/features/hr-applications'

export default function CompanyApplicationDetailPage() {
  const { companyId = '', applicationId = '' } = useParams<{
    companyId: string
    applicationId: string
  }>()
  return (
    <div className="mx-auto max-w-4xl p-6">
      <ApplicationDetail companyId={companyId} applicationId={applicationId} />
    </div>
  )
}
