import axios, {
  type AxiosError,
  type AxiosRequestConfig,
  type InternalAxiosRequestConfig,
} from 'axios'
import { env } from '@/shared/config/env'

// ─── Error types ─────────────────────────────────────────────────────────────

export class AppError extends Error {
  readonly code: string
  readonly status: number
  readonly fields?: Record<string, string>

  constructor(
    code: string,
    message: string,
    status: number,
    fields?: Record<string, string>,
  ) {
    super(message)
    this.name = 'AppError'
    this.code = code
    this.status = status
    this.fields = fields
  }
}

export class SessionExpiredError extends AppError {
  constructor() {
    super('SESSION_EXPIRED', 'Session expired', 401)
    this.name = 'SessionExpiredError'
  }
}

// ─── Single-flight refresh state ──────────────────────────────────────────────

let refreshPromise: Promise<void> | null = null
const queuedRequests: Array<{
  resolve: () => void
  reject: (e: unknown) => void
}> = []

function processQueue(error: unknown): void {
  const queued = queuedRequests.splice(0)
  for (const q of queued) {
    if (error) {
      q.reject(error)
    } else {
      q.resolve()
    }
  }
}

// ─── Error normalizer ─────────────────────────────────────────────────────────

function normalizeAxiosError(err: AxiosError): AppError {
  const status = err.response?.status ?? 0
  const data = err.response?.data as Record<string, unknown> | undefined
  // Auth/Company shape: { error: string }  — Candidate shape: { message: string }
  const message =
    (typeof data?.['error'] === 'string' ? data['error'] : undefined) ??
    (typeof data?.['message'] === 'string' ? data['message'] : undefined) ??
    err.message
  return new AppError('HTTP_ERROR', message, status)
}

// ─── Axios instance ───────────────────────────────────────────────────────────

export const httpClient = axios.create({
  baseURL: env.VITE_API_URL ?? '', // empty = relative URL on same origin
  withCredentials: true,
})

// Request interceptor: correlation id + Accept header
httpClient.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  config.headers.set('Accept', 'application/json')
  config.headers.set('X-Request-Id', crypto.randomUUID())
  return config
})

// Response interceptor: 401 single-flight refresh + error normalization
httpClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean
    }

    // Non-401 errors or already-retried requests → normalize and reject
    if (error.response?.status !== 401 || originalRequest._retry) {
      return Promise.reject(normalizeAxiosError(error))
    }

    // If the failing request IS the refresh endpoint → session is dead, do not retry
    if (originalRequest.url?.includes('/auth/refresh')) {
      window.dispatchEvent(new CustomEvent('session:expired'))
      return Promise.reject(new SessionExpiredError())
    }

    originalRequest._retry = true

    // If a refresh is already in-flight, queue this request behind it
    if (refreshPromise !== null) {
      return new Promise<void>((resolve, reject) => {
        queuedRequests.push({ resolve, reject })
      }).then(() => httpClient(originalRequest))
    }

    // Start a single refresh call
    refreshPromise = httpClient
      .post('/auth/refresh')
      .then(() => {
        processQueue(null)
      })
      .catch((err: unknown) => {
        processQueue(err)
        window.dispatchEvent(new CustomEvent('session:expired'))
        return Promise.reject(new SessionExpiredError())
      })
      .finally(() => {
        refreshPromise = null
      })

    return refreshPromise.then(() => httpClient(originalRequest))
  },
)

// ─── orval mutator export ─────────────────────────────────────────────────────

export const mutatorFn = <T>(config: AxiosRequestConfig): Promise<T> =>
  httpClient.request<T>(config).then((r) => r.data)

export default mutatorFn
