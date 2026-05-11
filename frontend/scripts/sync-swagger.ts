/**
 * sync-swagger.ts — fetches backend Swagger / OpenAPI specs into frontend/swagger/.
 *
 * For each known service we try sources in order until one succeeds:
 *   1. Live HTTP endpoint (e.g. http://localhost:8081/swagger/doc.json)
 *   2. Filesystem path under backend/.../docs/swagger.yaml
 *
 * Outputs are written next to the existing frontend/swagger/* files used by
 * codegen.ts so that running `pnpm codegen` afterwards regenerates clients.
 */

import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs'
import { request as httpRequest } from 'node:http'
import { request as httpsRequest } from 'node:https'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const rootDir = path.resolve(__dirname, '..')
const swaggerDir = path.join(rootDir, 'swagger')
const repoRoot = path.resolve(rootDir, '..')

interface SpecSource {
  /** logical name used in output filenames */
  name: 'auth' | 'candidate' | 'company' | 'application'
  /** ordered list of HTTP URLs to try */
  http: string[]
  /** ordered list of filesystem paths to try (relative to repo root) */
  files: string[]
  /** target filename(s) under frontend/swagger/ */
  outputs: { yaml?: string; json?: string }
}

const sources: SpecSource[] = [
  {
    name: 'auth',
    http: ['http://localhost:8000/swagger/doc.json'],
    files: ['backend/auth/docs/swagger.yaml', 'backend/auth/docs/swagger.json'],
    outputs: { yaml: 'swagger-auth.yaml', json: 'swagger-auth.json' },
  },
  {
    name: 'candidate',
    http: ['http://localhost:8081/swagger/doc.json'],
    files: [
      'backend/candidate/docs/swagger.yaml',
      'backend/candidate/docs/swagger.json',
    ],
    outputs: { yaml: 'swagger-candidate.yaml', json: 'swagger-candidate.json' },
  },
  {
    name: 'company',
    http: ['http://localhost:8088/swagger/doc.json'],
    files: [
      'backend/company/api/docs/swagger.yaml',
      'backend/company/docs/swagger.yaml',
    ],
    outputs: { yaml: 'swagger-company.yaml' },
  },
  {
    name: 'application',
    http: ['http://localhost:8086/swagger/doc.json'],
    files: [
      'backend/application/api/contract.yaml',
      'backend/application/docs/openapi.yaml',
    ],
    outputs: { yaml: 'openapi-application.yaml' },
  },
]

function fetchHttp(url: string, timeoutMs = 4000): Promise<string> {
  return new Promise((resolve, reject) => {
    const lib = url.startsWith('https:') ? httpsRequest : httpRequest
    const req = lib(url, (res) => {
      if (!res.statusCode || res.statusCode < 200 || res.statusCode >= 300) {
        reject(new Error(`HTTP ${res.statusCode} for ${url}`))
        res.resume()
        return
      }
      let data = ''
      res.setEncoding('utf8')
      res.on('data', (chunk: string) => {
        data += chunk
      })
      res.on('end', () => resolve(data))
    })
    req.on('error', reject)
    req.setTimeout(timeoutMs, () => {
      req.destroy(new Error(`timeout ${url}`))
    })
    req.end()
  })
}

function jsonToYaml(jsonText: string): string {
  return jsonText
}

/**
 * cleanSwaggerDefinitions — shortens Go‑package‑path definition names produced by swaggo.
 *
 * swaggo emits definition names like
 *   github_com_HghaVlad_trainee-match_backend_candidate_internal_delivery_http_dto.CandidateResponse
 * which orval translates into ~100‑char TypeScript identifiers.
 *
 * This function renames every definition matching a long Go‑import prefix
 * to the shorter form `dto.{ShortName}` and updates all $ref pointers.
 */
function cleanSwaggerDefinitions(raw: string): string {
  let obj: Record<string, unknown>
  try {
    obj = JSON.parse(raw)
  } catch {
    return raw
  }

  const defs = obj?.definitions as Record<string, unknown> | undefined
  if (!defs) return raw

  const goPkgPattern = /^github_com_HghaVlad_trainee-match_backend_.+\.(.+)$/
  const renameMap = new Map<string, string>()

  for (const key of Object.keys(defs)) {
    const match = key.match(goPkgPattern)
    if (match) {
      const short = `dto.${match[1]}`
      renameMap.set(key, short)
    }
  }

  if (renameMap.size === 0) return raw

  for (const [oldKey, newKey] of renameMap) {
    defs[newKey] = defs[oldKey]
    delete defs[oldKey]
    console.log(`  ↻ definition: ${oldKey} → ${newKey}`)
  }

  function walkRefs(node: unknown): void {
    if (Array.isArray(node)) {
      for (const item of node) walkRefs(item)
    } else if (node && typeof node === 'object') {
      const record = node as Record<string, unknown>
      if (typeof record['$ref'] === 'string') {
        for (const [oldKey, newKey] of renameMap) {
          const oldRef = `#/definitions/${oldKey}`
          if (record['$ref'] === oldRef) {
            record['$ref'] = `#/definitions/${newKey}`
            break
          }
        }
      }
      for (const val of Object.values(record)) walkRefs(val)
    }
  }
  walkRefs(obj)

  return JSON.stringify(obj, null, 2)
}

interface Resolved {
  source: string
  body: string
  isJson: boolean
}

async function resolve(spec: SpecSource): Promise<Resolved | null> {
  for (const url of spec.http) {
    try {
      const body = await fetchHttp(url)
      return { source: url, body, isJson: true }
    } catch (err) {
      console.warn(`  ! ${url}: ${(err as Error).message}`)
    }
  }
  for (const rel of spec.files) {
    const abs = path.join(repoRoot, rel)
    if (existsSync(abs)) {
      const body = readFileSync(abs, 'utf8')
      const isJson = abs.endsWith('.json')
      return { source: abs, body, isJson }
    }
  }
  return null
}

async function main(): Promise<void> {
  mkdirSync(swaggerDir, { recursive: true })
  let failed = 0
  for (const spec of sources) {
    console.log(`\n[${spec.name}] resolving...`)
    const resolved = await resolve(spec)
    if (!resolved) {
      console.warn(`  ✗ no source available for ${spec.name}`)
      failed++
      continue
    }
    console.log(`  ✓ source: ${resolved.source}`)

    const cleanedBody =
      resolved.isJson ? cleanSwaggerDefinitions(resolved.body) : resolved.body

    if (spec.outputs.json) {
      const out = path.join(swaggerDir, spec.outputs.json)
      const text = resolved.isJson
        ? cleanedBody
        : JSON.stringify({ note: 'binary YAML — keep yaml output' }, null, 2)
      if (resolved.isJson) writeFileSync(out, text)
      else console.warn(`    skip ${spec.outputs.json} (source is YAML)`)
    }

    if (spec.outputs.yaml) {
      const out = path.join(swaggerDir, spec.outputs.yaml)
      const text = resolved.isJson ? jsonToYaml(cleanedBody) : resolved.body
      writeFileSync(out, text)
      console.log(`  → ${out}`)
    }
  }
  if (failed > 0) {
    console.error(`\nDone with ${failed} unresolved spec(s).`)
    process.exit(failed === sources.length ? 1 : 0)
  }
  console.log('\nAll specs synced. Run `pnpm codegen` to regenerate clients.')
}

main().catch((err: unknown) => {
  console.error('sync-swagger failed:', err)
  process.exit(1)
})
