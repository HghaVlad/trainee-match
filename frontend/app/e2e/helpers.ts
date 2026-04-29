import { expect, type Page, type APIRequestContext, request } from '@playwright/test'

const AUTH_URL = process.env['E2E_AUTH_URL'] ?? 'http://localhost:8000'

export type Role = 'Candidate' | 'Company'

export interface TestUser {
  username: string
  password: string
  email: string
  firstName: string
  lastName: string
  role: Role
}

function rand(): string {
  return `${Date.now()}${Math.floor(Math.random() * 1e6)}`
}

export function makeUser(role: Role): TestUser {
  const id = rand()
  const prefix = role === 'Company' ? 'cmp' : 'cnd'
  const username = `e2e${prefix}${id}`
  return {
    username,
    password: 'Password123!',
    email: `${username}@test.local`,
    firstName: role === 'Company' ? 'Comp' : 'Cand',
    lastName: 'User',
    role,
  }
}

export async function registerViaApi(user: TestUser, api?: APIRequestContext): Promise<void> {
  const ctx = api ?? (await request.newContext())
  const res = await ctx.post(`${AUTH_URL}/api/v1/auth/register`, {
    data: {
      username: user.username,
      password: user.password,
      email: user.email,
      first_name: user.firstName,
      last_name: user.lastName,
      role: user.role,
    },
  })
  if (!res.ok()) {
    throw new Error(`register failed ${res.status()}: ${await res.text()}`)
  }
  if (!api) await ctx.dispose()
}

export async function registerAndLogin(page: Page, user: TestUser): Promise<void> {
  await registerViaApi(user)
  await loginViaUi(page, user)
}

export async function loginViaUi(page: Page, user: TestUser): Promise<void> {
  await page.goto('/login')
  await page.getByLabel('Имя пользователя').fill(user.username)
  await page.getByLabel('Пароль').fill(user.password)
  await page.getByRole('button', { name: 'Войти' }).click()
  await page.waitForURL((url) => !/\/login(\?|$)/.test(url.pathname + url.search), {
    timeout: 15_000,
  })
}

export async function expectAuthedHeader(page: Page, username: string): Promise<void> {
  await expect(page.getByRole('banner').getByText(username)).toBeVisible()
  await expect(page.getByRole('banner').getByRole('button', { name: 'Logout' })).toBeVisible()
}

export async function expectAnonHeader(page: Page): Promise<void> {
  await expect(page.getByRole('banner').getByRole('link', { name: 'Login' })).toBeVisible()
  await expect(page.getByRole('banner').getByRole('link', { name: 'Register' })).toBeVisible()
  await expect(page.getByRole('banner').getByRole('button', { name: 'Logout' })).toHaveCount(0)
}
