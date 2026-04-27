import { test, expect } from '@playwright/test'

test.describe('login flow', () => {
  test('candidate logs in and lands on /me/profile', async ({ page }) => {
    await page.goto('/login')
    await expect(page.getByRole('heading', { name: 'Вход' })).toBeVisible()

    await page.getByLabel('Имя пользователя').fill('candidate')
    await page.getByLabel('Пароль').fill('password')
    await page.getByRole('button', { name: 'Войти' }).click()

    await page.waitForURL('**/me/profile')
    expect(page.url()).toContain('/me/profile')
  })

  test('company logs in and lands on /company/me', async ({ page }) => {
    await page.goto('/login')
    await page.getByLabel('Имя пользователя').fill('company-acme')
    await page.getByLabel('Пароль').fill('password')
    await page.getByRole('button', { name: 'Войти' }).click()

    await page.waitForURL('**/company/me')
    expect(page.url()).toContain('/company/me')
  })
})
