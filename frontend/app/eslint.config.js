import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import tseslint from 'typescript-eslint'
import { defineConfig, globalIgnores } from 'eslint/config'
import boundaries from 'eslint-plugin-boundaries'
import unusedImports from 'eslint-plugin-unused-imports'

export default defineConfig([
  globalIgnores(['dist/**', 'src/api/generated/**', 'node_modules/**', 'tailwind.config.ts']),
  {
    files: ['**/*.{ts,tsx}'],
    ignores: ['src/shared/ui/**', 'src/shared/hooks/use-toast.ts'],
    extends: [
      js.configs.recommended,
      tseslint.configs.recommended,
      reactHooks.configs.flat.recommended,
      reactRefresh.configs.vite,
    ],
    plugins: {
      boundaries,
      'unused-imports': unusedImports,
    },
    languageOptions: {
      globals: globals.browser,
    },
    settings: {
      'boundaries/elements': [
        { type: 'app', pattern: 'src/app/**' },
        { type: 'pages', pattern: 'src/pages/**' },
        { type: 'widgets', pattern: 'src/widgets/**' },
        { type: 'features', pattern: 'src/features/**' },
        { type: 'entities', pattern: 'src/entities/**' },
        { type: 'shared', pattern: 'src/shared/**' },
        { type: 'api', pattern: 'src/api/**' },
      ],
    },
    rules: {
      'unused-imports/no-unused-imports': 'error',
      'no-console': ['warn', { allow: ['warn', 'error'] }],
      '@typescript-eslint/no-explicit-any': 'error',
      'boundaries/element-types': [
        'warn',
        {
          default: 'disallow',
          rules: [
            { from: 'app', allow: ['app', 'pages', 'widgets', 'features', 'entities', 'shared', 'api'] },
            { from: 'pages', allow: ['widgets', 'features', 'entities', 'shared', 'api'] },
            { from: 'widgets', allow: ['features', 'entities', 'shared', 'api'] },
            { from: 'features', allow: ['entities', 'shared', 'api'] },
            { from: 'entities', allow: ['shared', 'api'] },
            { from: 'shared', allow: ['shared'] },
            { from: 'api', allow: ['shared', 'api'] },
          ],
        },
      ],
    },
  },
])
