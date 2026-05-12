import { useState } from 'react'
import { Button } from '@/shared/ui/button'
import { useSession } from '@/shared/session/useSession'
import { useListMyApplications } from '@/api/generated/application/candidate-applications/candidate-applications'
import { findActiveApplication } from './status'
import { ApplyVacancyModal } from './ApplyVacancyModal'

interface Props {
  vacancyId: string
}

export function ApplyVacancyButton({ vacancyId }: Props) {
  const { isAuthed, role } = useSession()
  const [open, setOpen] = useState(false)
  const enabled = isAuthed && role === 'Candidate'

  const myApps = useListMyApplications(undefined, {
    query: { enabled },
  })

  if (!enabled) return null

  const active = findActiveApplication(myApps.data?.data, vacancyId)
  const hasActive = Boolean(active)

  const label = hasActive ? 'Вы уже откликнулись' : 'Откликнуться'
  const title = hasActive
    ? 'У вас уже есть активный отклик на эту вакансию'
    : undefined

  return (
    <>
      <Button
        type="button"
        disabled={hasActive || myApps.isLoading}
        title={title}
        onClick={() => setOpen(true)}
      >
        {label}
      </Button>
      <ApplyVacancyModal
        vacancyId={vacancyId}
        open={open}
        onOpenChange={setOpen}
      />
    </>
  )
}
