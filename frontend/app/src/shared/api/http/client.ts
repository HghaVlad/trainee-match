import axios, {
  type AxiosError,
  type AxiosRequestConfig,
  type AxiosInstance,
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

// ─── Single-flight refresh state (shared across all clients) ─────────────────

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
  const message =
    (typeof data?.['error'] === 'string' ? data['error'] : undefined) ??
    (typeof data?.['message'] === 'string' ? data['message'] : undefined) ??
    err.message
  return new AppError('HTTP_ERROR', message, status)
}

// ─── Common base URL fallback ─────────────────────────────────────────────────

const commonBaseURL = env.VITE_API_URL || '/api/v1'

function serviceBaseURL(overrideUrl: string | undefined): string {
  return overrideUrl && overrideUrl.length > 0 ? overrideUrl : commonBaseURL
}

// ─── Params serializer ────────────────────────────────────────────────────────

function serializeParams(params: Record<string, unknown>): string {
  const usp = new URLSearchParams()
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === null) continue
    if (Array.isArray(value)) {
      for (const item of value) {
        if (item === undefined || item === null || item === '') continue
        usp.append(key, String(item))
      }
    } else if (value !== '') {
      usp.append(key, String(value))
    }
  }
  return usp.toString()
}

// ─── Auth client (internal, for refresh calls) ────────────────────────────────

const authRefreshBaseURL = serviceBaseURL(env.VITE_AUTH_URL)
export const authRefreshClient = axios.create({
  baseURL: authRefreshBaseURL,
  withCredentials: true,
  paramsSerializer: serializeParams,
})

authRefreshClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    config.headers.set('Accept', 'application/json')
    config.headers.set('X-Request-Id', crypto.randomUUID())
    return config
  },
)

// ─── Client factory ───────────────────────────────────────────────────────────

function createApiClient(baseURL: string): AxiosInstance {
  const client = axios.create({
    baseURL,
    withCredentials: true,
    paramsSerializer: serializeParams,
  })

  client.interceptors.request.use((config: InternalAxiosRequestConfig) => {
    config.headers.set('Accept', 'application/json')
    config.headers.set('X-Request-Id', crypto.randomUUID())
    return config
  })

  client.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
      const originalRequest = error.config as InternalAxiosRequestConfig & {
        _retry?: boolean
      }

      if (error.response?.status !== 401 || originalRequest._retry) {
        return Promise.reject(normalizeAxiosError(error))
      }

      if (originalRequest.url?.includes('/auth/refresh')) {
        window.dispatchEvent(new CustomEvent('session:expired'))
        return Promise.reject(new SessionExpiredError())
      }

      originalRequest._retry = true

      if (refreshPromise !== null) {
        return new Promise<void>((resolve, reject) => {
          queuedRequests.push({ resolve, reject })
        }).then(() => client(originalRequest))
      }

      refreshPromise = authRefreshClient
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

      return refreshPromise.then(() => client(originalRequest))
    },
  )

  return client
}

// ─── Service-specific clients ─────────────────────────────────────────────────

export const httpClient = createApiClient(commonBaseURL)
export const authClient = createApiClient(serviceBaseURL(env.VITE_AUTH_URL))
export const candidateClient = createApiClient(
  serviceBaseURL(env.VITE_CANDIDATE_URL),
)
export const companyClient = createApiClient(serviceBaseURL(env.VITE_COMPANY_URL))
export const applicationClient = createApiClient(
  serviceBaseURL(env.VITE_APPLICATION_URL),
)

// ─── orval mutator exports ────────────────────────────────────────────────────

export const mutatorFn = <T>(config: AxiosRequestConfig): Promise<T> =>
  httpClient.request<T>(config).then((r) => r.data)

export const authMutatorFn = <T>(config: AxiosRequestConfig): Promise<T> =>
  authClient.request<T>(config).then((r) => r.data)

export const candidateMutatorFn = <T>(config: AxiosRequestConfig): Promise<T> =>
  candidateClient.request<T>(config).then((r) => r.data)

export const companyMutatorFn = <T>(config: AxiosRequestConfig): Promise<T> =>
  companyClient.request<T>(config).then((r) => r.data)

export const applicationMutatorFn = <T>(
  config: AxiosRequestConfig,
): Promise<T> => applicationClient.request<T>(config).then((r) => r.data)

export default mutatorFn
