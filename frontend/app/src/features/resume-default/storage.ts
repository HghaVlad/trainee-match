const DEFAULT_RESUME_KEY_PREFIX = 'tm.defaultResumeId.'

export function defaultResumeKey(userId: string | number): string {
  return `${DEFAULT_RESUME_KEY_PREFIX}${userId}`
}

export function readDefaultResumeId(
  userId: string | number,
): string | undefined {
  if (typeof window === 'undefined') return undefined
  try {
    return window.localStorage.getItem(defaultResumeKey(userId)) ?? undefined
  } catch {
    return undefined
  }
}

export function writeDefaultResumeId(
  userId: string | number,
  resumeId: string | undefined,
): void {
  if (typeof window === 'undefined') return
  try {
    const key = defaultResumeKey(userId)
    if (resumeId) {
      window.localStorage.setItem(key, resumeId)
    } else {
      window.localStorage.removeItem(key)
    }
  } catch {
    return
  }
}
