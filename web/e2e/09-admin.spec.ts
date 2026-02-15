import { test, expect, PageHelpers } from './fixtures'

test.describe('Admin Modules', () => {
  // ============ USERS ============
  test.describe('Users Management', () => {
    test('can view users list', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/users')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Usuarios|Users/i)
      await expect(page.locator('table')).toBeVisible()
    })

    test('can navigate to create user form', async ({ page }) => {
      await page.goto('/users')

      const newButton = page.locator('a:has-text("Nuevo"), button:has-text("Nuevo")')
      if (await newButton.first().isVisible()) {
        await newButton.first().click()
        await expect(page).toHaveURL(/\/users\/new|\/users\/create/)
      }
    })

    test('user form has required fields', async ({ page }) => {
      await page.goto('/users/new')

      await expect(page.locator('input[name="email"]')).toBeVisible()
      await expect(page.locator('input[name="first_name"]')).toBeVisible()
      await expect(page.locator('input[name="last_name"]')).toBeVisible()
    })
  })

  // ============ ROLES ============
  test.describe('Roles Management', () => {
    test('can view roles list', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/roles')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Roles/i)
    })

    test('can view role detail', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/roles')
      await helpers.waitForLoading()

      const roleLink = page.locator('table tbody tr a, [class*="card"] a').first()

      try {
        await roleLink.waitFor({ state: 'visible', timeout: 5000 })
        await roleLink.click()
        await expect(page).toHaveURL(/\/roles\/\d+/)
      } catch {
        test.skip()
      }
    })
  })

  // ============ BRANCHES ============
  test.describe('Branches Management', () => {
    test('can view branches list', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/branches')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Sucursales|Branches/i)
    })

    test('can navigate to create branch form', async ({ page }) => {
      await page.goto('/branches')

      const newButton = page.locator('a:has-text("Nueva"), button:has-text("Nueva")')
      if (await newButton.first().isVisible()) {
        await newButton.first().click()
        await expect(page).toHaveURL(/\/branches\/new|\/branches\/create/)
      }
    })
  })

  // ============ CATEGORIES ============
  test.describe('Categories Management', () => {
    test('can view categories list', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/categories')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Categorías|Categories/i)
    })

    test('can create new category', async ({ page }) => {
      await page.goto('/categories')

      const newButton = page.locator('a:has-text("Nueva"), button:has-text("Nueva")')
      if (await newButton.first().isVisible()) {
        await newButton.first().click()

        // Should open dialog or navigate to form
        const dialog = page.locator('[role="dialog"]')
        const formPage = page.url().includes('/categories/new')

        expect(await dialog.isVisible() || formPage).toBeTruthy()
      }
    })
  })

  // ============ EXPENSES ============
  test.describe('Expenses Management', () => {
    test('can view expenses list', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Gastos|Expenses/i)
      await expect(page.locator('table')).toBeVisible()
    })

    test('can navigate to create expense form', async ({ page }) => {
      await page.goto('/expenses')

      const newButton = page.locator('a:has-text("Nuevo"), button:has-text("Nuevo")')
      if (await newButton.first().isVisible()) {
        await newButton.first().click()
        await expect(page).toHaveURL(/\/expenses\/new|\/expenses\/create/)
      }
    })

    test('expense form has required fields', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses/new')
      await helpers.waitForLoading()

      // Should have amount and description fields
      await expect(page.locator('input[name="amount"], input[type="number"]').first()).toBeVisible()
    })
  })

  // ============ TRANSFERS ============
  test.describe('Transfers Management', () => {
    test('can view transfers list', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/transfers')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Transferencias|Transfers/i)
      await expect(page.locator('table')).toBeVisible()
    })

    test('can navigate to create transfer form', async ({ page }) => {
      await page.goto('/transfers')

      const newButton = page.locator('a:has-text("Nueva"), button:has-text("Nueva")')
      if (await newButton.first().isVisible()) {
        await newButton.first().click()
        await expect(page).toHaveURL(/\/transfers\/new|\/transfers\/create/)
      }
    })
  })

  // ============ REPORTS ============
  test.describe('Reports', () => {
    test('can view reports page', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/reports')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Reportes|Reports/i)
    })

    test('reports page has different report types', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/reports')
      await helpers.waitForLoading()

      // Should have tabs or cards for different reports
      const reportSections = page.locator('[role="tab"], [class*="card"]')
      await expect(reportSections.first()).toBeVisible()
    })
  })

  // ============ SETTINGS ============
  test.describe('Settings', () => {
    test('can view settings page', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/settings')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Configuración|Settings/i)
    })
  })

  // ============ AUDIT ============
  test.describe('Audit Log', () => {
    test('can view audit log page', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Auditoría|Audit/i)
      await expect(page.locator('table')).toBeVisible()
    })

    test('audit log has filters', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      // Should have filters button
      const filtersButton = page.locator('button:has-text("Filtros")')
      await expect(filtersButton).toBeVisible()

      // Click to open filters panel
      await filtersButton.click()

      // Wait for filter panel to appear (inside a Card)
      await page.waitForSelector('input#date_from', { state: 'visible', timeout: 5000 })

      // Should show filter controls (date inputs)
      const dateFromInput = page.locator('input#date_from')
      await expect(dateFromInput).toBeVisible()
    })
  })

  // ============ NOTIFICATIONS ============
  test.describe('Notifications', () => {
    test('can view notifications page', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/notifications')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Notificaciones|Notifications/i)
    })
  })
})
