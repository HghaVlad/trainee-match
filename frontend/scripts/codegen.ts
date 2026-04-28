import { exec } from 'node:child_process'
import { copyFileSync, mkdirSync } from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { promisify } from 'node:util'

const execAsync = promisify(exec)
const __dirname = path.dirname(fileURLToPath(import.meta.url))
const rootDir = path.resolve(__dirname, '..')
const swaggerDir = path.join(rootDir, 'swagger')
const appDir = path.join(rootDir, 'app')
const cacheDir = path.join(appDir, '.codegen-cache', 'openapi')

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
    const out = path.join(cacheDir, `${spec.name}.json`)
    console.log(`  ${spec.name}: ${spec.input} -> ${out}`)
    await convertSwagger(spec.input, out)
  }
  console.log('Copying OpenAPI 3 specs...')
  for (const spec of openapi3Specs) {
    const out = path.join(cacheDir, `${spec.name}.yaml`)
    console.log(`  ${spec.name}: ${spec.input} -> ${out}`)
    copyFileSync(spec.input, out)
  }
  console.log('Running orval...')
  await runOrval()
  console.log('Codegen complete.')
}

main().catch((err: unknown) => {
  console.error('Codegen failed:', err)
  process.exit(1)
})
