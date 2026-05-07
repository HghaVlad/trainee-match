import { test, expect } from '@playwright/test'
import { makeUser, registerAndLogin } from './helpers'

async function setupCandidate(page: import('@playwright/test').Page) {
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
}

test.describe('wave 4 — public pages', () => {
  test('header uses RU labels', async ({ page }) => {
    await page.goto('/')
    const banner = page.getByRole('banner')
    await expect(banner.getByRole('link', { name: 'Вакансии' })).toBeVisible()
    await expect(banner.getByRole('link', { name: 'Компании' })).toBeVisible()
    await expect(banner.getByRole('link', { name: 'Войти' })).toBeVisible()
  })

  test('/vacancies — numeric filters reject negatives via min attribute', async ({ page }) => {
    await page.goto('/vacancies')
    const salaryMin = page.getByLabel('Зарплата от, ₽')
    await expect(salaryMin).toHaveAttribute('min', '0')
    await expect(salaryMin).toHaveAttribute('step', '1000')

    const hoursMin = page.getByLabel('Часов/неделю от')
    await expect(hoursMin).toHaveAttribute('min', '1')
    await expect(hoursMin).toHaveAttribute('max', '168')

    const durationMin = page.getByLabel('Длительность от (дней)')
    await expect(durationMin).toHaveAttribute('min', '1')
    await expect(durationMin).toHaveAttribute('max', '730')
  })

  test('/vacancies — clamps salary on blur', async ({ page }) => {
    await page.goto('/vacancies')
    const salaryMin = page.getByLabel('Зарплата от, ₽')
    await salaryMin.fill('-500')
    await salaryMin.blur()
    await expect(salaryMin).toHaveValue('0')
  })

  test('/vacancies — flags inverted ranges as errors', async ({ page }) => {
    await page.goto('/vacancies')
    await page.getByLabel('Зарплата от, ₽').fill('100000')
    await page.getByLabel('Зарплата до, ₽').fill('50000')
    await expect(page.getByText(/Зарплата «от» больше «до»/)).toBeVisible()
  })

  test('/vacancies — accepts company_id from URL search param', async ({ page }) => {
    await page.goto('/vacancies?company_id=test-uuid-123')
    await expect(page.getByLabel('ID компании')).toHaveValue('test-uuid-123')
  })
})

test.describe('wave 4 — resume editor', () => {
  test('education uses year-pickers labelled "Начало обучения"/"Конец обучения"', async ({
    page,
  }) => {
    await setupCandidate(page)
    await page.goto('/me/resumes')
    await page.getByRole('button', { name: /Создать/ }).first().click()
    await page.waitForURL(/\/me\/resumes\/[0-9a-f-]{36}/, { timeout: 15_000 })

    await page.getByRole('button', { name: 'Добавить' }).first().click()

    await expect(page.getByText('Начало обучения')).toBeVisible()
    await expect(page.getByText('Конец обучения')).toBeVisible()
    await expect(page.getByText(/^С$/)).toHaveCount(0)
    await expect(page.getByText(/^По$/)).toHaveCount(0)
  })

  test('typing in resume fields does NOT trigger PATCH (no auto-save)', async ({
    page,
  }) => {
    await setupCandidate(page)
    await page.goto('/me/resumes')
    await page.getByRole('button', { name: /Создать/ }).first().click()
    await page.waitForURL(/\/me\/resumes\/[0-9a-f-]{36}/, { timeout: 15_000 })

    let patchCount = 0
    page.on('request', (req) => {
      const url = new URL(req.url())
      if (req.method() === 'PATCH' && /\/resume\//.test(url.pathname)) {
        patchCount += 1
      }
    })

    await page.getByLabel('Название резюме').fill('Auto Save Test')
    await page.waitForTimeout(4_000)
    expect(patchCount, 'auto-save PATCH must NOT fire').toBe(0)
  })

  test('manual Save button triggers PATCH', async ({ page }) => {
    await setupCandidate(page)
    await page.goto('/me/resumes')
    await page.getByRole('button', { name: /Создать/ }).first().click()
    await page.waitForURL(/\/me\/resumes\/[0-9a-f-]{36}/, { timeout: 15_000 })

    await page.getByLabel('Название резюме').fill('Manual Save Test')
    const [resp] = await Promise.all([
      page.waitForResponse(
        (r) =>
          /\/resume\//.test(new URL(r.url()).pathname) &&
          r.request().method() === 'PATCH',
        { timeout: 15_000 },
      ),
      page.getByRole('button', { name: /^Сохранить$/ }).click(),
    ])
    expect(resp.status()).toBeLessThan(400)
  })
})
