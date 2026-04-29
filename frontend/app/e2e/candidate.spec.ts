import { test, expect } from '@playwright/test'
import { makeUser, registerAndLogin } from './helpers'

test.describe('candidate flow', () => {
  test('BUG #2: save profile from /me/profile succeeds (no 404 on /candidate)', async ({ page }) => {
    const u = makeUser('Candidate')
    await registerAndLogin(page, u)
    await page.waitForURL('**/me/profile')

    const stamp = Date.now().toString().slice(-9)
    await page.getByLabel('Телефон').fill(`+7${stamp}`)
    await page.getByLabel('Telegram').fill(`@e2e${stamp}`)
    await page.getByLabel('Город').fill('Moscow')
    await page.getByLabel('Дата рождения').fill('2000-01-01')

    const [resp] = await Promise.all([
      page.waitForResponse(
        (r) => /candidate/.test(new URL(r.url()).pathname) && ['POST', 'PATCH', 'PUT'].includes(r.request().method()),
        { timeout: 15_000 },
      ),
      page.getByRole('button', { name: /Сохранить/ }).click(),
    ])

    expect(resp.status(), `save returned ${resp.status()}: ${await resp.text().catch(() => '')}`).toBeLessThan(400)
  })

  test('BUG #4: /me/resumes shows usable UI (no naked QueryProvider error, has create CTA)', async ({ page }) => {
    const u = makeUser('Candidate')
    await registerAndLogin(page, u)

    await page.goto('/me/resumes')
    await expect(page.getByRole('heading', { name: 'Мои резюме' })).toBeVisible()

    await expect(page.getByText(/Something went wrong/i)).toHaveCount(0)
    await expect(page.getByText(/QueryProvider/i)).toHaveCount(0)

    const createCta = page.getByRole('link', { name: /(Создать|Новое резюме|Add)/ }).or(
      page.getByRole('button', { name: /(Создать|Новое резюме|Add)/ }),
    )
    await expect(createCta.first()).toBeVisible()
  })

  test('BUG #5: /me/applications shows usable UI (no naked error)', async ({ page }) => {
    const u = makeUser('Candidate')
    await registerAndLogin(page, u)

    await page.goto('/me/applications')
    await expect(page.getByRole('heading', { name: 'Мои отклики' })).toBeVisible()
    await expect(page.getByText(/Something went wrong/i)).toHaveCount(0)
  })

  test('BUG #15: /me/profile edit mode preloads birthday and has Cancel button', async ({ page }) => {
    const u = makeUser('Candidate')
    await registerAndLogin(page, u)
    await page.waitForURL('**/me/profile')

    const stamp = Date.now().toString().slice(-9)
    await page.getByLabel('Телефон').fill(`+7${stamp}`)
    await page.getByLabel('Telegram').fill(`@e2e${stamp}`)
    await page.getByLabel('Город').fill('Moscow')
    await page.getByLabel('Дата рождения').fill('2000-01-01')
    await page.getByRole('button', { name: /Сохранить/ }).click()

    await expect(page.getByText('Дата рождения:')).toBeVisible({ timeout: 10_000 })
    await page.getByRole('button', { name: 'Редактировать' }).click()

    const birthday = page.getByLabel('Дата рождения')
    await expect(birthday).toHaveValue('2000-01-01')

    const cancel = page.getByRole('button', { name: 'Отмена' })
    await expect(cancel).toBeVisible()
    await cancel.click()
    await expect(page.getByRole('button', { name: 'Редактировать' })).toBeVisible()
  })

  test('BUG #16: /me/profile save in edit mode succeeds (no invalid request body)', async ({ page }) => {
    const u = makeUser('Candidate')
    await registerAndLogin(page, u)
    await page.waitForURL('**/me/profile')

    const stamp = Date.now().toString().slice(-9)
    await page.getByLabel('Телефон').fill(`+7${stamp}`)
    await page.getByLabel('Telegram').fill(`@e2e${stamp}`)
    await page.getByLabel('Город').fill('Moscow')
    await page.getByLabel('Дата рождения').fill('2000-01-01')
    await page.getByRole('button', { name: /Сохранить/ }).click()

    await page.getByRole('button', { name: 'Редактировать' }).click()
    await page.getByLabel('Город').fill('Saint-Petersburg')

    const [resp] = await Promise.all([
      page.waitForResponse(
        (r) =>
          /candidate/.test(new URL(r.url()).pathname) &&
          ['PATCH', 'PUT', 'POST'].includes(r.request().method()),
        { timeout: 15_000 },
      ),
      page.getByRole('button', { name: /Сохранить/ }).click(),
    ])
    expect(
      resp.status(),
      `edit save returned ${resp.status()}: ${await resp.text().catch(() => '')}`,
    ).toBeLessThan(400)
  })

  test('BUG #17: /me/resumes create succeeds after profile completed', async ({ page }) => {
    const u = makeUser('Candidate')
    await registerAndLogin(page, u)
    await page.waitForURL('**/me/profile')

    const stamp = Date.now().toString().slice(-9)
    await page.getByLabel('Телефон').fill(`+7${stamp}`)
    await page.getByLabel('Telegram').fill(`@e2e${stamp}`)
    await page.getByLabel('Город').fill('Moscow')
    await page.getByLabel('Дата рождения').fill('2000-01-01')
    await page.getByRole('button', { name: /Сохранить/ }).click()
    await expect(page.getByText('Дата рождения:')).toBeVisible({ timeout: 10_000 })

    await page.goto('/me/resumes')
    await expect(page.getByRole('heading', { name: 'Мои резюме' })).toBeVisible()

    const [resp] = await Promise.all([
      page.waitForResponse(
        (r) =>
          /\/resume(\b|\/|\?)/.test(new URL(r.url()).pathname) &&
          r.request().method() === 'POST',
        { timeout: 15_000 },
      ),
      page.getByRole('button', { name: /Создать/ }).first().click(),
    ])
    expect(
      resp.status(),
      `resume create returned ${resp.status()}: ${await resp.text().catch(() => '')}`,
    ).toBeLessThan(400)
  })
})
