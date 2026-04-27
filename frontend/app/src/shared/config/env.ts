import { z } from 'zod'

const envSchema = z.object({
  VITE_API_URL: z.string().url('VITE_API_URL must be a valid URL'),
  VITE_USE_MSW: z.string().optional().transform((v) => v === 'true'),
  VITE_APP_ENV: z
    .enum(['development', 'staging', 'production'])
    .optional()
    .default('development'),
  VITE_AUTH_ME_AVAILABLE: z.string().optional().transform((v) => v === 'true'),
})

export type Env = z.infer<typeof envSchema>

const result = envSchema.safeParse(import.meta.env)

if (!result.success) {
  const missing = result.error.issues
    .map((i) => `  ${i.path.join('.')}: ${i.message}`)
    .join('\n')
  throw new Error(`Invalid environment configuration:\n${missing}`)
}

export const env: Readonly<Env> = Object.freeze(result.data)
