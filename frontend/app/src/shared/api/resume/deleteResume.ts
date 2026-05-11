import { useMutation, useQueryClient } from '@tanstack/react-query'
import { httpClient } from '@/shared/api/http/client'
import { getGetResumeQueryKey } from '@/api/generated/candidate/resume/resume'

export async function deleteResume(id: string): Promise<void> {
  await httpClient.delete(`/resume/${id}`)
}

export function useDeleteResume() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => deleteResume(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: getGetResumeQueryKey() })
    },
  })
}
