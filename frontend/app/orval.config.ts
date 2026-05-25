import { defineConfig } from 'orval'

const sharedOutput = {
  mode: 'tags-split' as const,
  client: 'react-query' as const,
  httpClient: 'axios' as const,
}

export default defineConfig({
  auth: {
    input: { target: './.codegen-cache/openapi/auth.json' },
    output: {
      ...sharedOutput,
      target: './src/api/generated/auth',
      schemas: './src/api/generated/auth/schemas',
      override: {
        mutator: {
          path: './src/shared/api/http/client.ts',
          name: 'authMutatorFn',
        },
      },
    },
  },
  candidate: {
    input: { target: './.codegen-cache/openapi/candidate.json' },
    output: {
      ...sharedOutput,
      target: './src/api/generated/candidate',
      schemas: './src/api/generated/candidate/schemas',
      override: {
        mutator: {
          path: './src/shared/api/http/client.ts',
          name: 'candidateMutatorFn',
        },
      },
    },
  },
  company: {
    input: { target: './.codegen-cache/openapi/company.json' },
    output: {
      ...sharedOutput,
      target: './src/api/generated/company',
      schemas: './src/api/generated/company/schemas',
      override: {
        mutator: {
          path: './src/shared/api/http/client.ts',
          name: 'companyMutatorFn',
        },
      },
    },
  },
  application: {
    input: { target: './.codegen-cache/openapi/application.yaml' },
    output: {
      ...sharedOutput,
      target: './src/api/generated/application',
      schemas: './src/api/generated/application/schemas',
      override: {
        mutator: {
          path: './src/shared/api/http/client.ts',
          name: 'applicationMutatorFn',
        },
      },
    },
  },
})
