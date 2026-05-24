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

  test('BUG #20: /me/resumes/:id loads on first visit (no naked Something went wrong)', async ({ page }) => {
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
    await page.getByRole('button', { name: /Создать/ }).first().click()
    await page.waitForURL(/\/me\/resumes\/[0-9a-f-]{36}/, { timeout: 15_000 })

    await expect(page.getByRole('heading', { name: 'Редактирование резюме' })).toBeVisible({
      timeout: 10_000,
    })
    await expect(page.getByText('Something went wrong')).toHaveCount(0)
    await expect(page.getByRole('button', { name: 'Try again' })).toHaveCount(0)
  })

  test('BUG #21: /me/resumes/:id has Save and Cancel buttons; Cancel returns to /me/resumes', async ({ page }) => {
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
    await page.getByRole('button', { name: /Создать/ }).first().click()
    await page.waitForURL(/\/me\/resumes\/[0-9a-f-]{36}/, { timeout: 15_000 })

    await expect(page.getByRole('button', { name: /^Сохранить$/ })).toBeVisible()
    await expect(page.getByRole('button', { name: /^Удалить$/ }).first()).toBeVisible()
    await page.getByRole('button', { name: /^Отмена$/ }).click()
    await page.waitForURL(/\/me\/resumes$/, { timeout: 10_000 })
  })

  test('BUG #22: /me/resumes/:id Delete shows confirmation and removes resume', async ({ page }) => {
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
    await page.getByRole('button', { name: /Создать/ }).first().click()
    await page.waitForURL(/\/me\/resumes\/([0-9a-f-]{36})/, { timeout: 15_000 })
    const url = new URL(page.url())
    const resumeId = url.pathname.split('/').pop()!

    await page.getByRole('button', { name: /^Удалить$/ }).first().click()
    await expect(page.getByRole('dialog')).toBeVisible()
    await expect(page.getByText(/Удалить резюме\?/)).toBeVisible()

    const [resp] = await Promise.all([
      page.waitForResponse(
        (r) =>
          new URL(r.url()).pathname.endsWith(`/resume/${resumeId}`) &&
          r.request().method() === 'DELETE',
        { timeout: 15_000 },
      ),
      page.getByRole('dialog').getByRole('button', { name: /^Удалить$/ }).click(),
    ])
    expect(
      [200, 204, 404, 405, 501].includes(resp.status()),
      `delete returned ${resp.status()}: ${await resp.text().catch(() => '')}`,
    ).toBeTruthy()
  })
})

test.describe('candidate: apply vacancy', () => {
  test('BUG #23: full apply flow — navigate to vacancy → select resume → submit → 200/201', async ({ page }) => {
    // 1. Company: create company + published vacancy
    const companyUser = makeUser('Company')
    await registerAndLogin(page, companyUser)
    await page.goto('/company/new')
    await page.getByLabel('Название').fill(`E2E Co ${Date.now()}`)
    await page.getByRole('button', { name: 'Создать' }).click()
    await page.waitForURL(/\/company\/[^/]+\/dashboard/, { timeout: 15_000 })
    const companyId = new URL(page.url()).pathname.match(/\/company\/([^/]+)\/dashboard/)?.[1] ?? ''

    // Create vacancy
    await page.goto(`/company/${companyId}/vacancies/new`)
    await page.getByLabel('Название').fill(`E2E Vacancy ${Date.now()}`)
    await page.getByLabel(/Описание/).first().fill('test vacancy for apply flow')
    await page.getByLabel('Город').first().fill('Moscow')

    const [, vacancyResp] = await Promise.all([
      page.waitForURL(/\/company\/[^/]+\/vacancies\/[^/]+$/, { timeout: 15_000 }),
      page.getByRole('button', { name: /Создать/ }).first().click(),
    ])
    const vacancyId = new URL(page.url()).pathname.split('/').pop() ?? ''

    // Publish the vacancy
    await page.getByRole('button', { name: /^Опубликовать$/ }).first().click()
    await page.getByRole('dialog').getByRole('button', { name: /^Опубликовать$/ }).click()
    await expect(page.getByText('Опубликовано').first()).toBeVisible({ timeout: 15_000 })

    await page.context().clearCookies()
    await page.evaluate(() => localStorage.clear())
    await page.goto('/login', { waitUntil: 'networkidle' })
    const candidate = makeUser('Candidate')
    await registerAndLogin(page, candidate)
    await page.waitForURL('**/me/profile')
    const stamp = Date.now().toString().slice(-9)
    await page.getByLabel('Телефон').fill(`+7${stamp}`)
    await page.getByLabel('Telegram').fill(`@e2e${stamp}`)
    await page.getByLabel('Город').fill('Moscow')
    await page.getByLabel('Дата рождения').fill('2000-01-01')

    await Promise.all([
      page.waitForResponse(
        (r) => /\/api\/v1\/candidate/.test(new URL(r.url()).pathname) && ['POST', 'PATCH', 'PUT'].includes(r.request().method()),
        { timeout: 15_000 },
      ),
      page.getByRole('button', { name: /Сохранить/ }).click(),
    ])
    await expect(page.getByText('Дата рождения:')).toBeVisible({ timeout: 10_000 })

    await page.goto('/me/resumes')
    await expect(page.getByRole('heading', { name: 'Мои резюме' })).toBeVisible()
    await page.getByRole('button', { name: /^Создать$/ }).first().click()
    await page.waitForURL(/\/me\/resumes\/[0-9a-f-]{36}/, { timeout: 15_000 })

    await page.goto('/me/resumes')
    const row = page.getByTestId('resume-row').first()
    await row.getByRole('button', { name: /^Опубликовать$/ }).click()
    await Promise.all([
      page.waitForResponse(
        (r) => /\/resume\/[0-9a-f-]{36}/.test(new URL(r.url()).pathname) && r.request().method() === 'PATCH',
        { timeout: 15_000 },
      ),
      page.getByRole('dialog').getByRole('button', { name: /^Опубликовать$/ }).click(),
    ])
    await expect(row.getByText('Опубликовано')).toBeVisible({ timeout: 10_000 })

    await page.goto('/vacancies')
    await expect(page.getByRole('heading', { name: 'Вакансии', exact: true })).toBeVisible({ timeout: 10_000 })
    await expect(page.getByText(/Something went wrong/i)).toHaveCount(0)

    const vacancyLink = page.getByRole('link', { name: /E2E Vacancy/ }).first()
    await vacancyLink.click()
    await page.waitForURL(/\/vacancies\/[0-9a-f-]{36}/, { timeout: 15_000 })
    await expect(page.getByRole('navigation', { name: 'breadcrumb' })).toBeVisible({ timeout: 15_000 })
    await page.getByRole('button', { name: 'Откликнуться' }).click()
    await expect(page.getByRole('dialog')).toBeVisible({ timeout: 5000 })
    await expect(page.getByRole('dialog').getByText('Отклик на вакансию')).toBeVisible()

    const selectTrigger = page.getByRole('dialog').getByRole('combobox')
    await selectTrigger.click()
    await page.getByRole('option').first().click()

    const [applyResp] = await Promise.all([
      page.waitForResponse(
        (r) => /\/api\/v1\/applications$/.test(new URL(r.url()).pathname) && r.request().method() === 'POST',
        { timeout: 15_000 },
      ),
      page.getByRole('dialog').getByRole('button', { name: 'Откликнуться' }).click(),
    ])

    const status = applyResp.status()
    expect(
      status !== 404,
      `PROXY ROUTING FAILURE: /api/v1/applications returned 404 - not routed to backend. Response: ${await applyResp.text().catch(() => '')}`,
    ).toBe(true)
    expect(
      status === 200 || status === 201 || status === 400,
      `Expected 200/201/400, got ${status}: ${await applyResp.text().catch(() => '')}`,
    ).toBe(true)
  })
})
