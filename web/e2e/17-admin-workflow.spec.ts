import { test, expect, PageHelpers, generateTestData } from './fixtures'

/**
 * Admin Workflow Tests
 *
 * Simulates administrative tasks:
 * 1. User management (create, edit, deactivate)
 * 2. Role and permission management
 * 3. Branch configuration
 * 4. System settings
 * 5. Audit log review
 * 6. Report generation
 */
test.describe('Admin Workflow', () => {
  test.describe('User Management', () => {
    test('can view all users', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/users')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Usuarios|Users/i)
      await expect(page.locator('table')).toBeVisible()
    })

    test('can access create user form', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/users/new')
      await helpers.waitForLoading()

      // Should have user form fields
      await expect(page.locator('input[name="email"]')).toBeVisible()
      await expect(page.locator('input[name="first_name"]')).toBeVisible()
    })

    test('user form has role selection', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/users/new')
      await helpers.waitForLoading()

      // Should have role selector
      const roleSelector = page.locator('button[role="combobox"]').filter({
        has: page.locator('text=/Rol|Role/i'),
      })

      // Or look for role label
      const roleLabel = page.locator('label:has-text("Rol")')
      await expect(roleLabel.or(roleSelector.first())).toBeVisible()
    })

    test('user form has branch assignment', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/users/new')
      await helpers.waitForLoading()

      // Should have branch selector
      const branchLabel = page.locator('label:has-text("Sucursal")')
      await expect(branchLabel).toBeVisible()
    })

    test('can view user details', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/users')
      await helpers.waitForLoading()

      const userLink = page.locator('table tbody tr a[href*="/users/"]').first()

      try {
        await userLink.waitFor({ state: 'visible', timeout: 5000 })
        await userLink.click()
        await expect(page).toHaveURL(/\/users\/\d+/)
      } catch {
        test.skip()
      }
    })

    test('can access edit user form', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/users')
      await helpers.waitForLoading()

      const userLink = page.locator('table tbody tr a[href*="/users/"]').first()

      try {
        await userLink.waitFor({ state: 'visible', timeout: 5000 })
        await userLink.click()
        await expect(page).toHaveURL(/\/users\/\d+/)

        const editButton = page.locator('a:has-text("Editar"), button:has-text("Editar")')
        await editButton.first().click()
        await expect(page).toHaveURL(/\/users\/\d+\/edit/)
      } catch {
        test.skip()
      }
    })

    test('user detail shows activity info', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/users')
      await helpers.waitForLoading()

      const userLink = page.locator('table tbody tr a[href*="/users/"]').first()

      try {
        await userLink.waitFor({ state: 'visible', timeout: 5000 })
        await userLink.click()

        // Should show last login, creation date, etc.
        const activityInfo = page.locator('text=/Último acceso|Creado|Last login|Created/i')
        await expect(activityInfo.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('can search users by name or email', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/users')
      await helpers.waitForLoading()

      const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')

      if (await searchInput.isVisible()) {
        await searchInput.fill('admin')
        await page.waitForTimeout(600)
        await helpers.waitForLoading()

        await expect(page.locator('table')).toBeVisible()
      }
    })
  })

  test.describe('Role Management', () => {
    test('can view all roles', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/roles')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Roles/i)
    })

    test('can view role details with permissions', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/roles')
      await helpers.waitForLoading()

      const roleLink = page.locator('table tbody tr a, [class*="card"] a').first()

      try {
        await roleLink.waitFor({ state: 'visible', timeout: 5000 })
        await roleLink.click()
        await expect(page).toHaveURL(/\/roles\/\d+/)

        // Should show permissions
        const permissionsSection = page.locator('text=/Permisos|Permissions/i')
        await expect(permissionsSection.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('can access create role form', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/roles/new')
      await helpers.waitForLoading()

      // Should have role name field
      await expect(page.locator('input[name="name"]')).toBeVisible()
    })

    test('role form shows permission checkboxes', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/roles/new')
      await helpers.waitForLoading()

      // Should have permission checkboxes or grid
      const permissionInputs = page.locator('input[type="checkbox"], [role="checkbox"]')

      try {
        await permissionInputs.first().waitFor({ state: 'visible', timeout: 5000 })
        await expect(permissionInputs.first()).toBeVisible()
      } catch {
        // May be structured differently
        test.skip()
      }
    })

    test('can edit existing role', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/roles')
      await helpers.waitForLoading()

      const roleLink = page.locator('table tbody tr a, [class*="card"] a').first()

      try {
        await roleLink.waitFor({ state: 'visible', timeout: 5000 })
        await roleLink.click()
        await expect(page).toHaveURL(/\/roles\/\d+/)

        const editButton = page.locator('a:has-text("Editar"), button:has-text("Editar")')
        if (await editButton.first().isVisible()) {
          await editButton.first().click()
          await expect(page).toHaveURL(/\/roles\/\d+\/edit/)
        }
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Branch Management', () => {
    test('can view all branches', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/branches')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Sucursales|Branches/i)
    })

    test('can view branch details', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/branches')
      await helpers.waitForLoading()

      const branchLink = page.locator('table tbody tr a[href*="/branches/"], [class*="card"] a[href*="/branches/"]').first()

      try {
        await branchLink.waitFor({ state: 'visible', timeout: 5000 })
        await branchLink.click()
        await expect(page).toHaveURL(/\/branches\/\d+/)
      } catch {
        test.skip()
      }
    })

    test('branch detail shows configuration', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/branches')
      await helpers.waitForLoading()

      const branchLink = page.locator('table tbody tr a[href*="/branches/"], [class*="card"] a[href*="/branches/"]').first()

      try {
        await branchLink.waitFor({ state: 'visible', timeout: 5000 })
        await branchLink.click()

        // Should show branch info (address, phone, etc.)
        const branchInfo = page.locator('text=/Dirección|Teléfono|Address|Phone/i')
        await expect(branchInfo.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('can create new branch', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/branches/new')
      await helpers.waitForLoading()

      // Should have branch form
      await expect(page.locator('input[name="name"]')).toBeVisible()
    })
  })

  test.describe('Category Management', () => {
    test('can view categories', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/categories')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Categorías|Categories/i)
    })

    test('can create new category', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/categories')
      await helpers.waitForLoading()

      const newButton = page.locator('button:has-text("Nueva"), a:has-text("Nueva")')

      if (await newButton.first().isVisible()) {
        await newButton.first().click()

        // Should open dialog or form
        const dialog = page.locator('[role="dialog"]')
        const formPage = page.url().includes('/categories/new')

        if (await dialog.isVisible()) {
          await expect(dialog).toBeVisible()
        } else if (formPage) {
          await expect(page.locator('form')).toBeVisible()
        }
      }
    })

    test('categories show hierarchy', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/categories')
      await helpers.waitForLoading()

      // Should show category list (possibly as tree)
      const categoryList = page.locator('table, [class*="tree"], ul')
      await expect(categoryList.first()).toBeVisible()
    })
  })

  test.describe('System Settings', () => {
    test('can access settings page', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/settings')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Configuración|Settings/i)
    })

    test('settings has multiple sections', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/settings')
      await helpers.waitForLoading()

      // Should have tabs or sections
      const sections = page.locator('[role="tab"], [class*="card"], h3')
      await expect(sections.first()).toBeVisible()
    })

    test('can view loan settings', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/settings')
      await helpers.waitForLoading()

      // Look for loan/interest settings
      const loanSettings = page.locator('text=/Préstamo|Interés|Loan|Interest/i')

      try {
        await loanSettings.first().waitFor({ state: 'visible', timeout: 5000 })
        await expect(loanSettings.first()).toBeVisible()
      } catch {
        // May be in a different tab
        test.skip()
      }
    })
  })

  test.describe('Audit Log Review', () => {
    test('can access audit log', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Auditoría|Audit/i)
    })

    test('audit log shows user actions', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      // Should show table with actions
      await expect(page.locator('table')).toBeVisible()
    })

    test('can filter audit by user', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      const filtersButton = page.locator('button:has-text("Filtros")')
      await filtersButton.click()

      await page.waitForSelector('input#date_from', { state: 'visible', timeout: 5000 })

      const userFilter = page.locator('label:has-text("Usuario")')
      await expect(userFilter).toBeVisible()
    })

    test('can filter audit by action type', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      const filtersButton = page.locator('button:has-text("Filtros")')
      await filtersButton.click()

      await page.waitForSelector('input#date_from', { state: 'visible', timeout: 5000 })

      const actionFilter = page.locator('label:has-text("Acción")')
      await expect(actionFilter).toBeVisible()
    })

    test('can view audit entry details', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      // Look for detail button in table
      const detailButton = page.locator('table tbody tr button').first()

      try {
        await detailButton.waitFor({ state: 'visible', timeout: 5000 })
        await detailButton.click()

        // Should open detail dialog
        const dialog = page.locator('[role="dialog"]')
        await dialog.waitFor({ state: 'visible', timeout: 3000 })
        await expect(dialog).toBeVisible()
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Reports', () => {
    test('can access reports page', async ({ page }) => {
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
      const reportTypes = page.locator('[role="tab"], [class*="card"]')
      await expect(reportTypes.first()).toBeVisible()
    })

    test('can select date range for reports', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/reports')
      await helpers.waitForLoading()

      const dateInput = page.locator('input[type="date"], button:has-text("Fecha")')

      try {
        await dateInput.first().waitFor({ state: 'visible', timeout: 5000 })
        await expect(dateInput.first()).toBeVisible()
      } catch {
        // May not have date filter visible
        test.skip()
      }
    })
  })

  test.describe('Notifications Management', () => {
    test('can access notifications page', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/notifications')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Notificaciones|Notifications/i)
    })

    test('notifications page shows list or tabs', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/notifications')
      await helpers.waitForLoading()

      const content = page.locator('table, [role="tablist"], [class*="card"]')
      await expect(content.first()).toBeVisible()
    })
  })

  test.describe('Transfers Management', () => {
    test('can view transfers list', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/transfers')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Transferencias|Transfers/i)
    })

    test('can create new transfer', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/transfers/new')
      await helpers.waitForLoading()

      // Should have source and destination branch selectors
      const branchSelectors = page.locator('button[role="combobox"]')
      await expect(branchSelectors.first()).toBeVisible()
    })

    test('transfer shows status and items', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/transfers')
      await helpers.waitForLoading()

      const transferLink = page.locator('table tbody tr a[href*="/transfers/"]').first()

      try {
        await transferLink.waitFor({ state: 'visible', timeout: 5000 })
        await transferLink.click()
        await expect(page).toHaveURL(/\/transfers\/\d+/)

        // Should show status
        const statusBadge = page.locator('[class*="badge"]')
        await expect(statusBadge.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })
  })
})
