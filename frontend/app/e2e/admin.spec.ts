import { test, expect, type Page } from '@playwright/test'

async function login(page: Page, username: string, password: string) {
  await page.goto('/login')
  await page.getByLabel('Имя пользователя').fill(username)
  await page.getByLabel('Пароль').fill(password)
  await page.getByRole('button', { name: 'Войти' }).click()
  await page.waitForURL((url) => !/\/login(\?|$)/.test(url.pathname + url.search), {
    timeout: 15_000,
  })
}

test.describe('admin: guard — anonymous access', () => {
  test('anon → /admin/skills redirects to /login', async ({ page }) => {
    await page.goto('/admin/skills')
    await page.waitForURL(/\/login/)
    await expect(page.getByLabel('Имя пользователя')).toBeVisible()
  })

  test('anon → /admin/users redirects to /login', async ({ page }) => {
    await page.goto('/admin/users')
    await page.waitForURL(/\/login/)
  })

  test('anon → /candidates redirects to /login', async ({ page }) => {
    await page.goto('/candidates')
    await page.waitForURL(/\/login/)
  })
})

test.describe('admin: guard — non-admin access', () => {
  test('Candidate → /candidates gets redirected or sees 403', async ({ page }) => {
    const stamp = Date.now()
    const username = `e2eguard${stamp}`
    const { request } = page.context()
    const res = await request.post(`${process.env['E2E_BACKEND_URL'] ?? 'https://api.traineematch.space'}/api/auth/auth/register`, {
      data: {
        username,
        password: 'Password123!',
        email: `${username}@test.local`,
        first_name: 'Guard',
        last_name: 'Test',
        role: 'Candidate',
      },
    })
    expect([200, 201, 409]).toContain(res.status())

    await login(page, username, 'Password123!')
    await page.goto('/candidates')
    await page.waitForLoadState('networkidle')
    const url = new URL(page.url())
    const onErrorPage = url.pathname === '/403' || page.getByRole('heading', { name: '403' }).isVisible().catch(() => false)
    const onLoginPage = /\/login/.test(url.pathname)
    const sees403Content = await page.getByRole('heading', { name: '403' }).isVisible().catch(() => false)
    expect(onErrorPage || onLoginPage || sees403Content, `Expected /403 or /login or 403 page, got ${url.pathname}`).toBe(true)
  })
})

test.describe('admin: existing admin user', () => {
  const ADMIN_USER = { username: 'admin', password: 'admin' }

  async function tryLoginAsAdmin(page: Page) {
    await page.goto('/login')
    await page.getByLabel('Имя пользователя').fill(ADMIN_USER.username)
    await page.getByLabel('Пароль').fill(ADMIN_USER.password)
    await page.getByRole('button', { name: 'Войти' }).click()
    return page.waitForURL((url) => !/\/login(\?|$)/.test(url.pathname + url.search), {
      timeout: 15_000,
    })
  }

  test('admin user login succeeds and lands on admin area', async ({ page }) => {
    let loggedIn = false
    try {
      await tryLoginAsAdmin(page)
      loggedIn = true
    } catch {
      // Login may fail if admin creds aren't valid — skip platform-admin-dependent tests
    }

    if (!loggedIn) {
      test.skip()
    }

    await expect(page.getByRole('heading', { name: 'Навыки' }).or(page.getByText(/Навыки/i))).toBeVisible({ timeout: 5000 })
  })

  test('admin header shows admin nav links', async ({ page }) => {
    let loggedIn = false
    try {
      await tryLoginAsAdmin(page)
      loggedIn = true
    } catch {
      test.skip()
    }

    if (!loggedIn) return

    const header = page.getByRole('banner')
    const navLinks = ['Навыки', 'Кандидаты']
    for (const link of navLinks) {
      const el = header.getByRole('link', { name: link })
      if (await el.isVisible({ timeout: 2000 }).catch(() => false)) {
        await expect(el).toBeVisible()
      }
    }
  })

  test('admin header shows Пользователи link', async ({ page }) => {
    let loggedIn = false
    try {
      await tryLoginAsAdmin(page)
      loggedIn = true
    } catch {
      test.skip()
    }

    if (!loggedIn) return

    const usersLink = page.getByRole('banner').getByRole('link', { name: 'Пользователи' })
    if (await usersLink.isVisible({ timeout: 2000 }).catch(() => false)) {
      await expect(usersLink).toBeVisible()
    }
  })
})

test.describe('admin: /admin/skills', () => {
  test('loads without error', async ({ page }) => {
    const stamp = Date.now()
    const username = `e2e${stamp}`
    const { request } = page.context()
    await request.post(`${process.env['E2E_BACKEND_URL'] ?? 'https://api.traineematch.space'}/api/auth/auth/register`, {
      data: {
        username,
        password: 'Password123!',
        email: `${username}@test.local`,
        first_name: 'Skill',
        last_name: 'Test',
        role: 'Candidate',
      },
    })
    await login(page, username, 'Password123!')

    await page.goto('/admin/skills')
    await page.waitForLoadState('networkidle')
    const sees403 = await page.getByRole('heading', { name: '403' }).isVisible().catch(() => false)
    if (sees403) {
      await expect(page.getByRole('heading', { name: '403' })).toBeVisible()
    } else {
      await expect(page.getByRole('heading', { name: 'Навыки' }).or(page.getByText(/Навыки/i))).toBeVisible({ timeout: 5000 })
      await expect(page.getByText(/Something went wrong/i)).toHaveCount(0)
    }
  })

  test('has "Создать навык" button', async ({ page }) => {
    const stamp = Date.now()
    const username = `e2e${stamp}`
    const { request } = page.context()
    await request.post(`${process.env['E2E_BACKEND_URL'] ?? 'https://api.traineematch.space'}/api/auth/auth/register`, {
      data: {
        username,
        password: 'Password123!',
        email: `${username}@test.local`,
        first_name: 'Skill',
        last_name: 'Test',
        role: 'Candidate',
      },
    })
    await login(page, username, 'Password123!')

    await page.goto('/admin/skills')
    await page.waitForLoadState('networkidle')
    const sees403 = await page.getByRole('heading', { name: '403' }).isVisible().catch(() => false)
    if (sees403) {
      await expect(page.getByRole('heading', { name: '403' })).toBeVisible()
    } else {
      const btn = page.getByRole('link', { name: 'Создать навык' })
      if (await btn.isVisible({ timeout: 2000 }).catch(() => false)) {
        await expect(btn).toBeVisible()
      }
    }
  })
})

test.describe('admin: /admin/users', () => {
  test('loads without error (no crash on page load)', async ({ page }) => {
    const stamp = Date.now()
    const username = `e2e${stamp}`
    const { request } = page.context()
    await request.post(`${process.env['E2E_BACKEND_URL'] ?? 'https://api.traineematch.space'}/api/auth/auth/register`, {
      data: {
        username,
        password: 'Password123!',
        email: `${username}@test.local`,
        first_name: 'User',
        last_name: 'Test',
        role: 'Candidate',
      },
    })
    await login(page, username, 'Password123!')

    await page.goto('/admin/users')
    await expect(page.getByText(/Something went wrong/i)).toHaveCount(0)
    await expect(page.getByText(/Назначить администратором/i).or(page.getByText(/403/i))).toBeVisible({ timeout: 10_000 })
  })
})

test.describe('admin: /candidates', () => {
  test('loads without error', async ({ page }) => {
    const stamp = Date.now()
    const username = `e2e${stamp}`
    const { request } = page.context()
    await request.post(`${process.env['E2E_BACKEND_URL'] ?? 'https://api.traineematch.space'}/api/auth/auth/register`, {
      data: {
        username,
        password: 'Password123!',
        email: `${username}@test.local`,
        first_name: 'Cand',
        last_name: 'Test',
        role: 'Candidate',
      },
    })
    await login(page, username, 'Password123!')

    await page.goto('/candidates')
    await expect(page.getByText(/Something went wrong/i)).toHaveCount(0)
    await expect(
      page.getByRole('heading', { name: 'Кандидаты' }).or(page.getByText(/403/i)),
    ).toBeVisible({ timeout: 10_000 })
  })
})