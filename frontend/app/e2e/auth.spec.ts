import { test, expect } from '@playwright/test'
import {
  makeUser,
  registerAndLogin,
  expectAuthedHeader,
  expectAnonHeader,
} from './helpers'

test.describe('auth: candidate', () => {
  test('register → login → land on profile', async ({ page }) => {
    const u = makeUser('Candidate')
    await registerAndLogin(page, u)

    await page.waitForURL('**/me/profile')
    await expectAuthedHeader(page, u.username)
  })

  test('BUG #1: logout fully clears session and redirects', async ({ page }) => {
    const u = makeUser('Candidate')
    await registerAndLogin(page, u)
    await page.waitForURL('**/me/profile')

    await page.getByRole('banner').getByRole('button', { name: 'Logout' }).click()

    await expectAnonHeader(page)
    await expect.poll(() => new URL(page.url()).pathname).not.toMatch(/^\/me\//)

    await page.goto('/me/profile')
    await page.waitForURL(/\/login/)
  })

  test('BUG #3: refresh on /me/profile keeps session (no kick to /login)', async ({ page }) => {
    const u = makeUser('Candidate')
    await registerAndLogin(page, u)
    await page.waitForURL('**/me/profile')

    await page.reload()
    await expect(page).toHaveURL(/\/me\/profile/)
    await expectAuthedHeader(page, u.username)
  })

  test('BUG #11/#12: deep-link without session → login → return to original', async ({ page }) => {
    const u = makeUser('Candidate')
    await registerAndLogin(page, u)
    await page.getByRole('banner').getByRole('button', { name: 'Logout' }).click()
    await expectAnonHeader(page)

    await page.goto('/me/resumes/new')
    await page.waitForURL(/\/login/)

    await page.getByLabel('Имя пользователя').fill(u.username)
    await page.getByLabel('Пароль').fill(u.password)
    await page.getByRole('button', { name: 'Войти' }).click()

    await page.waitForURL('**/me/resumes/new', { timeout: 15_000 })
    await expectAuthedHeader(page, u.username)
  })
})

test.describe('auth: company', () => {
  test('register → login → land in company area (no 403, header works)', async ({ page }) => {
    const u = makeUser('Company')
    await registerAndLogin(page, u)

    await expect.poll(() => new URL(page.url()).pathname).toMatch(/^\/company(\/|$)/)
    await expect(page).not.toHaveURL(/\/403/)
    await expectAuthedHeader(page, u.username)
  })
})
