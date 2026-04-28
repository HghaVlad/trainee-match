import { useCallback, useSyncExternalStore } from 'react'
import { useSessionUser } from '@/shared/session/useSession'
import {
  defaultResumeKey,
  readDefaultResumeId,
  writeDefaultResumeId,
} from './storage'

const listeners = new Set<() => void>()

function notify() {
  for (const fn of listeners) fn()
}

function subscribe(cb: () => void): () => void {
  listeners.add(cb)
  const onStorage = (e: StorageEvent) => {
    if (!e.key || e.key.startsWith('tm.defaultResumeId.')) cb()
  }
  if (typeof window !== 'undefined') {
    window.addEventListener('storage', onStorage)
  }
  return () => {
    listeners.delete(cb)
    if (typeof window !== 'undefined') {
      window.removeEventListener('storage', onStorage)
    }
  }
}

export function useDefaultResumeId(): {
  defaultResumeId: string | undefined
  setDefaultResumeId: (id: string | undefined) => void
  userId: string | number | undefined
} {
  const user = useSessionUser()
  const userId = user?.id

  const getSnapshot = useCallback(() => {
    if (userId === undefined) return undefined
    return readDefaultResumeId(userId)
  }, [userId])

  const defaultResumeId = useSyncExternalStore(
    subscribe,
    getSnapshot,
    () => undefined,
  )

  const setDefaultResumeId = useCallback(
    (id: string | undefined) => {
      if (userId === undefined) return
      writeDefaultResumeId(userId, id)
      notify()
    },
    [userId],
  )

  return { defaultResumeId, setDefaultResumeId, userId }
}

export { defaultResumeKey, readDefaultResumeId, writeDefaultResumeId }
