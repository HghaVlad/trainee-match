import { test, expect } from '@playwright/test'
import { makeUser, registerAndLogin, expectAnonHeader } from './helpers'

test.describe('public pages', () => {
  test('BUG #6: /vacancies renders without naked error', async ({ page }) => {
    await page.goto('/vacancies')
    await expect(page.getByRole('heading', { name: 'Вакансии', exact: true })).toBeVisible()
    await expect(page.getByText(/Something went wrong/i)).toHaveCount(0)
    await expect(page.getByText(/Произошла ошибка/i)).toHaveCount(0)
  })

  test('BUG #6: /companies renders without naked error', async ({ page }) => {
    await page.goto('/companies')
    await expect(page.getByRole('heading', { name: 'Компании' })).toBeVisible()
    await expect(page.getByText(/Something went wrong/i)).toHaveCount(0)
    await expect(page.getByText(/Произошла ошибка/i)).toHaveCount(0)
  })
})

test.describe('company flow', () => {
  test('BUG #8: anon → /company/new → login → returns to /company/new', async ({ page }) => {
    const u = makeUser('Company')

    await page.goto('/company/new')
    await page.waitForURL(/\/login/)

    await page.evaluate(() => localStorage.clear())

    const { request: req } = page.context()
    await req.post('http://localhost:8000/api/v1/auth/register', {
      data: {
        username: u.username,
        password: u.password,
        email: u.email,
        first_name: u.firstName,
        last_name: u.lastName,
        role: u.role,
      },
    })

    await page.getByLabel('Имя пользователя').fill(u.username)
    await page.getByLabel('Пароль').fill(u.password)
    await page.getByRole('button', { name: 'Войти' }).click()

    await page.waitForURL('**/company/new', { timeout: 15_000 })
    await expect(page.getByText('Создание компании')).toBeVisible()
  })

  test('BUG #9: create company succeeds (no 404)', async ({ page }) => {
    const u = makeUser('Company')
    await registerAndLogin(page, u)

    await page.goto('/company/new')
    await expect(page.getByText('Создание компании')).toBeVisible()

    await page.getByLabel('Название').fill(`E2E Co ${Date.now()}`)
    await page.getByLabel('Описание (необязательно)').fill('e2e test company')

    const [resp] = await Promise.all([
      page.waitForResponse(
        (r) =>
          /\/api\/v1\/companies(\b|\?|$)/.test(new URL(r.url()).pathname) &&
          r.request().method() === 'POST',
        { timeout: 15_000 },
      ),
      page.getByRole('button', { name: 'Создать' }).click(),
    ])

    expect(resp.status(), `create returned ${resp.status()}: ${await resp.text().catch(() => '')}`).toBeLessThan(400)
  })

  test('BUG #10: cancel button on /company/new navigates away', async ({ page }) => {
    const u = makeUser('Company')
    await registerAndLogin(page, u)

    await page.goto('/company/new')
    await expect(page.getByText('Создание компании')).toBeVisible()

    await page.getByRole('button', { name: 'Отмена' }).click()
    await expect.poll(() => new URL(page.url()).pathname).not.toBe('/company/new')
  })

  test('BUG #18: after creating company stay logged-in (no /403 redirect)', async ({ page }) => {
    const u = makeUser('Company')
    await registerAndLogin(page, u)

    await page.goto('/company/new')
    await expect(page.getByText('Создание компании')).toBeVisible()

    await page.getByLabel('Название').fill(`E2E Co ${Date.now()}`)
    await page.getByLabel('Описание (необязательно)').fill('e2e test company')
    await page.getByRole('button', { name: 'Создать' }).click()

    await page.waitForURL(/\/company\/[^/]+\/dashboard/, { timeout: 15_000 })
    expect(new URL(page.url()).pathname).not.toBe('/403')
  })

  test('BUG #19: /company/me resolves to active company dashboard (no /403)', async ({ page }) => {
    const u = makeUser('Company')
    await registerAndLogin(page, u)

    await page.goto('/company/new')
    await page.getByLabel('Название').fill(`E2E Me ${Date.now()}`)
    await page.getByRole('button', { name: 'Создать' }).click()
    await page.waitForURL(/\/company\/[^/]+\/dashboard/, { timeout: 15_000 })

    await page.goto('/company/me')
    await page.waitForURL((url) => /\/company\/[^/]+\/dashboard/.test(url.pathname), {
      timeout: 15_000,
    })
    expect(new URL(page.url()).pathname).not.toBe('/403')
  })
})

test.describe('logout regression', () => {
  test('logout from company area also clears session', async ({ page }) => {
    const u = makeUser('Company')
    await registerAndLogin(page, u)

    await page.getByRole('banner').getByRole('button', { name: 'Logout' }).click()
    await expectAnonHeader(page)
    await expect.poll(() => new URL(page.url()).pathname).not.toMatch(/^\/company/)
  })
})
