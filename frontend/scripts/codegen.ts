import { exec } from 'node:child_process'
import { mkdirSync, readFileSync, writeFileSync } from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { promisify } from 'node:util'

const execAsync = promisify(exec)
const __dirname = path.dirname(fileURLToPath(import.meta.url))
const rootDir = path.resolve(__dirname, '..')
const swaggerDir = path.join(rootDir, 'swagger')
const appDir = path.join(rootDir, 'app')
const cacheDir = path.join(appDir, '.codegen-cache', 'openapi')

/** @see sync-swagger.ts — keeps definition names short. */
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
      renameMap.set(key, `dto.${match[1]}`)
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

const swagger2Specs = [
  { name: 'auth', input: path.join(swaggerDir, 'swagger-auth.yaml') },
  { name: 'candidate', input: path.join(swaggerDir, 'swagger-candidate.yaml') },
  { name: 'company', input: path.join(swaggerDir, 'swagger-company.yaml') },
] as const

const openapi3Specs = [
  {
    name: 'application',
    input: path.join(swaggerDir, 'openapi-application.yaml'),
  },
] as const

async function convertSwagger(input: string, output: string): Promise<void> {
  const bin = path.join(appDir, 'node_modules', '.bin', 'swagger2openapi')
  const { stderr } = await execAsync(`"${bin}" --patch --outfile "${output}" "${input}"`)
  if (stderr && !/warn/i.test(stderr)) {
    console.warn(stderr)
  }
}

async function runOrval(): Promise<void> {
  const bin = path.join(appDir, 'node_modules', '.bin', 'orval')
  const { stdout, stderr } = await execAsync(`"${bin}"`, { cwd: appDir })
  if (stdout) console.log(stdout)
  if (stderr) console.warn(stderr)
}

async function main(): Promise<void> {
  mkdirSync(cacheDir, { recursive: true })
  console.log('Converting Swagger 2.0 -> OpenAPI 3...')
  for (const spec of swagger2Specs) {
    const raw = readFileSync(spec.input, 'utf8')
    const cleaned = cleanSwaggerDefinitions(raw)
    if (cleaned !== raw) {
      writeFileSync(spec.input, cleaned, 'utf8')
    }
    const out = path.join(cacheDir, `${spec.name}.json`)
    console.log(`  ${spec.name}: ${spec.input} -> ${out}`)
    await convertSwagger(spec.input, out)
  }
  console.log('Copying OpenAPI 3 specs...')
  for (const spec of openapi3Specs) {
    const out = path.join(cacheDir, `${spec.name}.yaml`)
    console.log(`  ${spec.name}: ${spec.input} -> ${out}`)
    const raw = readFileSync(spec.input, 'utf8')
    const stripped = raw.replace(/^(\s+)\/api\/v1\//gm, '$1/')
    writeFileSync(out, stripped)
  }
  console.log('Running orval...')
  await runOrval()
  console.log('Codegen complete.')
}

main().catch((err: unknown) => {
  console.error('Codegen failed:', err)
  process.exit(1)
})
