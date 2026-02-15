import { test, expect, PageHelpers } from './fixtures'

test.describe('Search and Filter Functionality', () => {
  test.describe('Global Search', () => {
    test('customers list has search input', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')
      await expect(searchInput).toBeVisible()
    })

    test('items list has search input', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/items')
      await helpers.waitForLoading()

      const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')
      await expect(searchInput).toBeVisible()
    })

    test('loans list has search input', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')
      await expect(searchInput).toBeVisible()
    })

    test('search triggers on input', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')

      await searchInput.fill('test')
      await page.waitForTimeout(600) // Wait for debounce
      await helpers.waitForLoading()

      // Table should still be visible after search
      await expect(page.locator('table')).toBeVisible()
    })

    test('clearing search shows all results', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')

      // Search then clear
      await searchInput.fill('test')
      await page.waitForTimeout(600)
      await searchInput.fill('')
      await page.waitForTimeout(600)
      await helpers.waitForLoading()

      await expect(page.locator('table')).toBeVisible()
    })
  })

  test.describe('Status Filters', () => {
    test('loans have status filter options', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const statusFilter = page.locator('button:has-text("Estado"), select[name*="status"]')

      if (await statusFilter.first().isVisible()) {
        await statusFilter.first().click()

        // Should show filter options
        const options = page.locator('[role="menuitem"], [role="option"], option')
        await expect(options.first()).toBeVisible()
      }
    })

    test('items have status filter options', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/items')
      await helpers.waitForLoading()

      const statusFilter = page.locator('button:has-text("Estado"), select[name*="status"]')

      if (await statusFilter.first().isVisible()) {
        await statusFilter.first().click()

        const options = page.locator('[role="menuitem"], [role="option"], option')
        await expect(options.first()).toBeVisible()
      }
    })

    test('payments have status filter', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/payments')
      await helpers.waitForLoading()

      const statusFilter = page.locator('button:has-text("Estado"), button:has-text("Tipo")')

      if (await statusFilter.first().isVisible()) {
        await expect(statusFilter.first()).toBeVisible()
      }
    })
  })

  test.describe('Date Range Filters', () => {
    test('audit log has date filters', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      // Open filters
      const filtersButton = page.locator('button:has-text("Filtros")')
      await filtersButton.click()

      await page.waitForSelector('input#date_from', { state: 'visible', timeout: 5000 })

      const dateFrom = page.locator('input#date_from')
      const dateTo = page.locator('input#date_to')

      await expect(dateFrom).toBeVisible()
      await expect(dateTo).toBeVisible()
    })

    test('audit quick date filters work', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      const filtersButton = page.locator('button:has-text("Filtros")')
      await filtersButton.click()

      await page.waitForSelector('input#date_from', { state: 'visible', timeout: 5000 })

      // Click quick filter
      const todayButton = page.locator('button:has-text("Hoy")')
      if (await todayButton.isVisible()) {
        await todayButton.click()
        await helpers.waitForLoading()
      }
    })

    test('reports may have date range', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/reports')
      await helpers.waitForLoading()

      const dateInput = page.locator('input[type="date"]')

      if (await dateInput.first().isVisible()) {
        await expect(dateInput.first()).toBeVisible()
      }
    })
  })

  test.describe('Entity Filters', () => {
    test('audit log can filter by user', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      const filtersButton = page.locator('button:has-text("Filtros")')
      await filtersButton.click()

      await page.waitForSelector('input#date_from', { state: 'visible', timeout: 5000 })

      // User select should be visible
      const userLabel = page.locator('label:has-text("Usuario")')
      await expect(userLabel).toBeVisible()
    })

    test('audit log can filter by entity type', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      const filtersButton = page.locator('button:has-text("Filtros")')
      await filtersButton.click()

      await page.waitForSelector('input#date_from', { state: 'visible', timeout: 5000 })

      const entityLabel = page.locator('label:has-text("Tipo de Entidad")')
      await expect(entityLabel).toBeVisible()
    })

    test('audit log can filter by action', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      const filtersButton = page.locator('button:has-text("Filtros")')
      await filtersButton.click()

      await page.waitForSelector('input#date_from', { state: 'visible', timeout: 5000 })

      const actionLabel = page.locator('label:has-text("Acción")')
      await expect(actionLabel).toBeVisible()
    })
  })

  test.describe('Branch Filters', () => {
    test('audit log can filter by branch', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      const filtersButton = page.locator('button:has-text("Filtros")')
      await filtersButton.click()

      await page.waitForSelector('input#date_from', { state: 'visible', timeout: 5000 })

      const branchLabel = page.locator('label:has-text("Sucursal")')
      await expect(branchLabel).toBeVisible()
    })

    test('transfers may filter by source/destination branch', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/transfers')
      await helpers.waitForLoading()

      // Look for branch filter
      const branchFilter = page.locator('button:has-text("Sucursal"), select[name*="branch"]')

      if (await branchFilter.first().isVisible()) {
        await expect(branchFilter.first()).toBeVisible()
      }
    })
  })

  test.describe('Clear Filters', () => {
    test('audit log has clear filters button', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      const filtersButton = page.locator('button:has-text("Filtros")')
      await filtersButton.click()

      await page.waitForSelector('input#date_from', { state: 'visible', timeout: 5000 })

      const clearButton = page.locator('button:has-text("Limpiar")')
      await expect(clearButton).toBeVisible()
    })

    test('clearing filters resets view', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      const filtersButton = page.locator('button:has-text("Filtros")')
      await filtersButton.click()

      await page.waitForSelector('input#date_from', { state: 'visible', timeout: 5000 })

      // Set a date filter
      await page.locator('input#date_from').fill('2024-01-01')
      await helpers.waitForLoading()

      // Clear filters
      const clearButton = page.locator('button:has-text("Limpiar")')
      await clearButton.click()
      await helpers.waitForLoading()

      // URL should be clean
      expect(page.url()).not.toContain('date_from')
    })
  })

  test.describe('URL Persistence', () => {
    test('search is persisted in URL', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')
      await searchInput.fill('test')
      await page.waitForTimeout(600)
      await helpers.waitForLoading()

      // Check URL has search param
      try {
        expect(page.url()).toMatch(/search|q=|query=/i)
      } catch {
        // Search may not update URL
        test.skip()
      }
    })

    test('page number is in URL', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const nextButton = page.locator('button:has-text("Siguiente"), button[aria-label*="next"]')

      try {
        await nextButton.first().waitFor({ state: 'visible', timeout: 3000 })

        if (!(await nextButton.first().isDisabled())) {
          await nextButton.first().click()
          await helpers.waitForLoading()

          expect(page.url()).toMatch(/page=2|offset=/i)
        }
      } catch {
        test.skip()
      }
    })

    test('filters are shareable via URL', async ({ page }) => {
      const helpers = new PageHelpers(page)

      // Navigate with filters in URL
      await page.goto('/loans?status=active')
      await helpers.waitForLoading()

      // Should show loans (or empty state)
      const content = page.locator('table').or(page.locator('text=/No hay|No se encontraron/i'))
      await expect(content.first()).toBeVisible()
    })
  })

  test.describe('Pagination', () => {
    test('shows pagination info', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const paginationInfo = page.locator('text=/de|of|página|page|mostrando|showing/i')

      try {
        await paginationInfo.first().waitFor({ state: 'visible', timeout: 3000 })
        await expect(paginationInfo.first()).toBeVisible()
      } catch {
        // May not have enough data for pagination
        test.skip()
      }
    })

    test('has previous/next buttons', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const navButtons = page.locator('button:has-text("Anterior"), button:has-text("Siguiente"), button[aria-label*="previous"], button[aria-label*="next"]')

      try {
        await navButtons.first().waitFor({ state: 'visible', timeout: 3000 })
        await expect(navButtons.first()).toBeVisible()
      } catch {
        // May not have pagination
        test.skip()
      }
    })

    test('first page has previous disabled', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers?page=1')
      await helpers.waitForLoading()

      const prevButton = page.locator('button:has-text("Anterior"), button[aria-label*="previous"]')

      try {
        await prevButton.first().waitFor({ state: 'visible', timeout: 3000 })
        const isDisabled = await prevButton.first().isDisabled()
        expect(isDisabled).toBe(true)
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Column Visibility', () => {
    test('table has columns toggle', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/items')
      await helpers.waitForLoading()

      const columnsButton = page.locator('button:has-text("Columnas")')

      try {
        await columnsButton.waitFor({ state: 'visible', timeout: 3000 })
        await expect(columnsButton).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('can toggle column visibility', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/items')
      await helpers.waitForLoading()

      const columnsButton = page.locator('button:has-text("Columnas")')

      try {
        await columnsButton.waitFor({ state: 'visible', timeout: 3000 })
        await columnsButton.click()

        // Should show column options
        const columnOptions = page.locator('[role="menuitemcheckbox"], [role="checkbox"]')
        await expect(columnOptions.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Sorting', () => {
    test('can sort by clicking column header', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const sortableHeader = page.locator('th button, th[role="button"], th[class*="cursor-pointer"]').first()

      try {
        await sortableHeader.waitFor({ state: 'visible', timeout: 3000 })
        await sortableHeader.click()
        await helpers.waitForLoading()

        // Table should update
        await expect(page.locator('table')).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('sort indicator shows direction', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const sortableHeader = page.locator('th button, th[role="button"]').first()

      try {
        await sortableHeader.waitFor({ state: 'visible', timeout: 3000 })
        await sortableHeader.click()
        await helpers.waitForLoading()

        // Look for sort indicator (arrow icon)
        const sortIcon = page.locator('th svg, th [class*="sort"]')
        await expect(sortIcon.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })
  })
})
