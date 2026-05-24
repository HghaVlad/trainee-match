import { test, expect, type Page } from '@playwright/test'
import { makeUser, registerAndLogin, registerViaApi } from './helpers'

async function setupCandidateAndProfile(page: Page) {
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
  return u
}

async function createCompany(page: Page): Promise<string> {
  await page.goto('/company/new')
  await page.getByLabel('Название').fill(`E2E Co ${Date.now()}`)
  await page.getByRole('button', { name: 'Создать' }).click()
  await page.waitForURL(/\/company\/[0-9a-f-]+\/dashboard/, { timeout: 15_000 })
  const m = new URL(page.url()).pathname.match(/\/company\/([^/]+)\/dashboard/)
  if (!m) throw new Error('failed to extract companyId from url')
  return m[1]
}

async function createResumeViaUi(page: Page) {
  await page.goto('/me/resumes')
  await expect(page.getByRole('heading', { name: 'Мои резюме' })).toBeVisible()
  await page.getByRole('button', { name: /^Создать$/ }).first().click()
  await page.waitForURL(/\/me\/resumes\/[0-9a-f-]{36}/, { timeout: 15_000 })
}

test.describe('wave 5 — candidate resumes', () => {
  test('#3: resume cards show status badge and "Основное" tag instead of dash', async ({
    page,
  }) => {
    await setupCandidateAndProfile(page)
    await createResumeViaUi(page)

    await page.goto('/me/resumes')
    const row = page.getByTestId('resume-row').first()
    await expect(row).toBeVisible()
    await expect(row.getByText('Черновик')).toBeVisible()
    await expect(row.getByText('—', { exact: true })).toHaveCount(0)
  })

  test('#1: resume publish flow — button publishes draft, star unlocks for published', async ({
    page,
  }) => {
    await setupCandidateAndProfile(page)
    await createResumeViaUi(page)

    await page.goto('/me/resumes')
    const row = page.getByTestId('resume-row').first()

    const star = row.getByRole('button', { name: /Опубликуйте резюме/ })
    await expect(star).toBeDisabled()

    await row.getByRole('button', { name: /^Опубликовать$/ }).click()
    const dialog = page.getByRole('dialog')
    await expect(dialog.getByText(/Опубликовать резюме\?/)).toBeVisible()
    await Promise.all([
      page.waitForResponse(
        (r) =>
          /\/resume\/[0-9a-f-]{36}/.test(new URL(r.url()).pathname) &&
          r.request().method() === 'PATCH',
        { timeout: 15_000 },
      ),
      dialog.getByRole('button', { name: /^Опубликовать$/ }).click(),
    ])

    await expect(row.getByText('Опубликовано')).toBeVisible({ timeout: 10_000 })
    const unlockedStar = row.getByRole('button', { name: /Сделать основным/ })
    await expect(unlockedStar).toBeEnabled()

    await unlockedStar.click()
    await expect(row.getByText('Основное')).toBeVisible()
  })

  test('#1: making published resume a draft removes "Основное" tag and star', async ({
    page,
  }) => {
    await setupCandidateAndProfile(page)
    await createResumeViaUi(page)

    await page.goto('/me/resumes')
    const row = page.getByTestId('resume-row').first()

    await row.getByRole('button', { name: /^Опубликовать$/ }).click()
    await Promise.all([
      page.waitForResponse(
        (r) =>
          /\/resume\/[0-9a-f-]{36}/.test(new URL(r.url()).pathname) &&
          r.request().method() === 'PATCH',
        { timeout: 15_000 },
      ),
      page
        .getByRole('dialog')
        .getByRole('button', { name: /^Опубликовать$/ })
        .click(),
    ])
    await row.getByRole('button', { name: /Сделать основным/ }).click()
    await expect(row.getByText('Основное')).toBeVisible()

    await row.getByRole('button', { name: /В черновик/ }).click()
    await Promise.all([
      page.waitForResponse(
        (r) =>
          /\/resume\/[0-9a-f-]{36}/.test(new URL(r.url()).pathname) &&
          r.request().method() === 'PATCH',
        { timeout: 15_000 },
      ),
      page
        .getByRole('dialog')
        .getByRole('button', { name: /В черновик/ })
        .click(),
    ])

    await expect(row.getByText('Черновик')).toBeVisible({ timeout: 10_000 })
    await expect(row.getByText('Основное')).toHaveCount(0)
  })
})

test.describe('wave 5 — public vacancies', () => {
  test('#4: invalid company_id UUID is flagged and not sent as filter', async ({
    page,
  }) => {
    let sentCompanyIdFilter = false
    page.on('request', (req) => {
      const url = new URL(req.url())
      if (
        /\/api\/v1\/vacancies(\b|\?|$)/.test(url.pathname) &&
        url.searchParams.has('company_id')
      ) {
        sentCompanyIdFilter = true
      }
    })

    await page.goto('/vacancies')
    const companyInput = page.getByLabel('ID компании')
    await companyInput.click()
    await companyInput.type('not-a-uuid')
    await companyInput.press('Enter')
    await expect(page.getByText(/ID компании должен быть в формате UUID/)).toBeVisible({
      timeout: 5_000,
    })

    await page.waitForTimeout(800)
    expect(sentCompanyIdFilter).toBe(false)
  })

  test('#5: typing into salary field debounces — single request after pause', async ({
    page,
  }) => {
    const requests: string[] = []
    page.on('request', (req) => {
      const url = new URL(req.url())
      if (/\/api\/v1\/vacancies(\b|\?|$)/.test(url.pathname)) {
        requests.push(url.search)
      }
    })

    await page.goto('/vacancies')
    await page.waitForTimeout(800)
    const baseline = requests.length

    const salaryMin = page.getByLabel('Зарплата от, ₽')
    await salaryMin.click()
    await salaryMin.type('1', { delay: 30 })
    await salaryMin.type('0', { delay: 30 })
    await salaryMin.type('0', { delay: 30 })
    await salaryMin.type('0', { delay: 30 })
    await salaryMin.type('0', { delay: 30 })

    await page.waitForTimeout(800)
    const newRequests = requests.length - baseline
    expect(newRequests).toBeLessThanOrEqual(2)
  })

  test('#wave5-cities: multiple city chips are sent as repeated query params', async ({
    page,
  }) => {
    const sentCities: string[][] = []
    page.on('request', (req) => {
      const url = new URL(req.url())
      if (
        /\/api\/v1\/vacancies(\b|\?|$)/.test(url.pathname) &&
        url.searchParams.has('city')
      ) {
        sentCities.push(url.searchParams.getAll('city'))
      }
    })

    await page.goto('/vacancies')
    const cityInput = page.getByLabel('Город')
    await cityInput.click()
    await cityInput.type('Moscow')
    await cityInput.press('Enter')
    await cityInput.type('Saint-Petersburg')
    await cityInput.press('Enter')

    await page.waitForTimeout(800)
    const last = sentCities[sentCities.length - 1] ?? []
    expect(last).toEqual(expect.arrayContaining(['Moscow', 'Saint-Petersburg']))
    expect(last.length).toBe(2)
  })

  test('#wave5-companies: multiple company_id chips are sent as repeated query params', async ({
    page,
  }) => {
    const sent: string[][] = []
    page.on('request', (req) => {
      const url = new URL(req.url())
      if (
        /\/api\/v1\/vacancies(\b|\?|$)/.test(url.pathname) &&
        url.searchParams.has('company_id')
      ) {
        sent.push(url.searchParams.getAll('company_id'))
      }
    })

    const id1 = '11111111-1111-1111-1111-111111111111'
    const id2 = '22222222-2222-2222-2222-222222222222'

    await page.goto('/vacancies')
    const companyInput = page.getByLabel('ID компании')
    await companyInput.click()
    await companyInput.type(id1)
    await companyInput.press('Enter')
    await companyInput.type(id2)
    await companyInput.press('Enter')

    await page.waitForTimeout(800)
    const last = sent[sent.length - 1] ?? []
    expect(last).toEqual(expect.arrayContaining([id1, id2]))
    expect(last.length).toBe(2)
  })
})

test.describe('wave 5 — HR vacancies', () => {
  test('#8: company vacancies table shows "Создана" column once and no "Открыть" button', async ({
    page,
  }) => {
    const u = makeUser('Company')
    await registerAndLogin(page, u)
    const companyId = await createCompany(page)

    await page.goto(`/company/${companyId}/vacancies/new`)
    await page.getByLabel('Название').fill(`E2E V ${Date.now()}`)
    await page.getByLabel(/Описание/).first().fill('e2e desc')
    await page.getByLabel('Город').first().fill('Moscow')
    await Promise.all([
      page.waitForResponse(
        (r) =>
          /\/api\/v1\/companies\/[^/]+\/vacancies(\b|\?|$)/.test(
            new URL(r.url()).pathname,
          ) && r.request().method() === 'POST',
        { timeout: 15_000 },
      ),
      page.getByRole('button', { name: /Создать|Сохранить/ }).first().click(),
    ])

    await page.goto(`/company/${companyId}/vacancies`)
    await expect(page.getByText('Вакансии', { exact: true }).first()).toBeVisible({
      timeout: 10_000,
    })
    await expect(page.locator('table tbody tr').first()).toBeVisible({
      timeout: 10_000,
    })

    await expect(page.getByRole('button', { name: /^Открыть$/ })).toHaveCount(0)
    await expect(page.getByRole('link', { name: /^Открыть$/ })).toHaveCount(0)

    const headers = page.locator('table thead th')
    const count = await headers.count()
    let createdCount = 0
    for (let i = 0; i < count; i += 1) {
      const text = (await headers.nth(i).textContent())?.trim() ?? ''
      if (text === 'Создана') createdCount += 1
    }
    expect(createdCount).toBe(1)
  })
})

test.describe('wave 5 — company switcher', () => {
  test('#12: dropdown contains "+ Новая компания" option that navigates to /company/new', async ({
    page,
  }) => {
    const u = makeUser('Company')
    await registerAndLogin(page, u)
    await createCompany(page)

    const switcher = page.getByLabel('Активная компания')
    await expect(switcher).toBeVisible()
    const optionTexts = await switcher.locator('option').allTextContents()
    expect(optionTexts.some((t) => t.includes('+ Новая компания'))).toBe(true)

    await switcher.selectOption({ label: '+ Новая компания' })
    await page.waitForURL('**/company/new', { timeout: 10_000 })
    await expect(page.getByText('Создание компании')).toBeVisible()
  })
})

test.describe('wave 5 — HR vacancy invalidation', () => {
  test('#9: publishing draft from list triggers immediate UI update', async ({
    page,
  }) => {
    const u = makeUser('Company')
    await registerAndLogin(page, u)
    const companyId = await createCompany(page)

    await page.goto(`/company/${companyId}/vacancies/new`)
    await expect(
      page.getByText(/Новая вакансия|Создание вакансии/).first(),
    ).toBeVisible({ timeout: 10_000 })
    await page.getByLabel('Название').fill(`E2E V ${Date.now()}`)
    const desc = page.getByLabel(/Описание/).first()
    await desc.fill('e2e description for vacancy publish flow')
    await page.getByLabel('Город').first().fill('Moscow')

    await Promise.all([
      page.waitForResponse(
        (r) =>
          /\/api\/v1\/companies\/[^/]+\/vacancies(\b|\?|$)/.test(
            new URL(r.url()).pathname,
          ) && r.request().method() === 'POST',
        { timeout: 15_000 },
      ),
      page.getByRole('button', { name: /Создать|Сохранить/ }).first().click(),
    ])

    await page.goto(`/company/${companyId}/vacancies`)
    await expect(page.getByText('Черновик').first()).toBeVisible({ timeout: 10_000 })

    await page.getByRole('button', { name: /^Опубликовать$/ }).first().click()
    const publishResp = page.waitForResponse(
      (r) =>
        /\/api\/v1\/companies\/[^/]+\/vacancies\/[^/]+\/publish/.test(
          new URL(r.url()).pathname,
        ),
      { timeout: 15_000 },
    )
    await page
      .getByRole('dialog')
      .getByRole('button', { name: /^Опубликовать$/ })
      .click()
    await publishResp

    await expect(page.getByText('Опубликовано').first()).toBeVisible({
      timeout: 15_000,
    })
  })
})

test.describe('wave 5 — company members', () => {
  test('#14: admin adds a recruiter by username', async ({ page }) => {
    const admin = makeUser('Company')
    await registerAndLogin(page, admin)
    const companyId = await createCompany(page)

    const target = makeUser('Company')
    await registerViaApi(target)

    const requests: string[] = []
    page.on('request', (req) => {
      const url = new URL(req.url())
      if (url.pathname.includes('/companies/')) {
        requests.push(`${req.method()} ${url.pathname}`)
      }
    })
    page.on('response', (res) => {
      const url = new URL(res.url())
      if (url.pathname.includes('/companies/') && res.status() >= 400) {
        requests.push(`FAIL ${res.status()} ${url.pathname}`)
      }
    })

    await page.goto(`/company/${companyId}/members`)
    await expect(page.getByText('Команда')).toBeVisible({
      timeout: 15_000,
    })
    await page.waitForTimeout(500)
    console.log('=== COMPANY REQUESTS ===', requests.join(' | '))

    await page.getByRole('button', { name: 'Добавить' }).click()
    await expect(page.getByRole('dialog')).toBeVisible()

    const usernameInput = page.getByRole('dialog').getByLabel('Username')
    await usernameInput.fill(target.username)

    await Promise.all([
      page.waitForResponse(
        (r) =>
          new URL(r.url()).pathname ===
          `/api/v1/companies/${companyId}/members` &&
          r.request().method() === 'POST',
        { timeout: 15_000 },
      ),
      page.getByRole('dialog').getByRole('button', { name: 'Добавить' }).click(),
    ])

    await expect(page.getByRole('dialog')).toHaveCount(0, { timeout: 15_000 })
    await expect(page.getByText(target.username)).toBeVisible({ timeout: 10_000 })
  })
})
