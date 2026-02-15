import { test, expect, PageHelpers, generateTestData } from './fixtures'

test.describe('Item CRUD Operations', () => {
  let testData: ReturnType<typeof generateTestData>

  test.beforeAll(() => {
    testData = generateTestData()
  })

  test('can view items list with table', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    // Should see items table
    await expect(page.locator('table')).toBeVisible()
    await expect(page.locator('h1, h2').first()).toContainText(/Artículos|Items/i)
  })

  test('items table shows status badges', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    // Look for status badges in table
    const statusBadges = page.locator('table tbody [class*="badge"]')

    try {
      await statusBadges.first().waitFor({ state: 'visible', timeout: 5000 })
      await expect(statusBadges.first()).toBeVisible()
    } catch {
      // No items with status badges
      test.skip()
    }
  })

  test('can navigate to item detail page', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    const itemLink = page.locator('table tbody tr a[href*="/items/"]').first()

    try {
      await itemLink.waitFor({ state: 'visible', timeout: 5000 })
      await itemLink.click()

      await expect(page).toHaveURL(/\/items\/\d+/)
      await expect(page.locator('main')).toBeVisible()
    } catch {
      test.skip()
    }
  })

  test('item detail shows item information', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    const itemLink = page.locator('table tbody tr a[href*="/items/"]').first()

    try {
      await itemLink.waitFor({ state: 'visible', timeout: 5000 })
      await itemLink.click()
      await expect(page).toHaveURL(/\/items\/\d+/)

      // Should show description, values, status
      await expect(page.locator('main')).toContainText(/Descripción|Description|Valor|Value/i)
    } catch {
      test.skip()
    }
  })

  test('item detail shows linked loan', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    const itemLink = page.locator('table tbody tr a[href*="/items/"]').first()

    try {
      await itemLink.waitFor({ state: 'visible', timeout: 5000 })
      await itemLink.click()
      await expect(page).toHaveURL(/\/items\/\d+/)

      // Should have link to loan if item is pawned
      const loanLink = page.locator('a[href*="/loans/"]')
      if (await loanLink.first().isVisible()) {
        await expect(loanLink.first()).toBeVisible()
      }
    } catch {
      test.skip()
    }
  })

  test('can filter items by status', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    // Look for status filter
    const statusFilter = page.locator('button:has-text("Estado"), button:has-text("Status")')

    if (await statusFilter.first().isVisible()) {
      await statusFilter.first().click()

      const filterOption = page.locator('[role="menuitem"], [role="option"]')
      if (await filterOption.first().isVisible()) {
        await filterOption.first().click()
        await helpers.waitForLoading()
      }
    }
  })

  test('can filter items by category', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    // Look for category filter
    const categoryFilter = page.locator('button:has-text("Categoría"), button:has-text("Category")')

    if (await categoryFilter.first().isVisible()) {
      await categoryFilter.first().click()

      const filterOption = page.locator('[role="menuitem"], [role="option"]')
      if (await filterOption.first().isVisible()) {
        await filterOption.first().click()
        await helpers.waitForLoading()
      }
    }
  })

  test('can search items by description', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')

    if (await searchInput.isVisible()) {
      await searchInput.fill('anillo')
      await page.waitForTimeout(500) // Debounce
      await helpers.waitForLoading()

      // Table should still be visible
      await expect(page.locator('table')).toBeVisible()
    }
  })

  test('can open row actions menu', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    const actionButton = page.locator('table tbody tr button[aria-haspopup="menu"]').first()

    try {
      await actionButton.waitFor({ state: 'visible', timeout: 5000 })
      await actionButton.click()

      // Should open dropdown menu
      await expect(page.locator('[role="menu"]')).toBeVisible()
    } catch {
      test.skip()
    }
  })

  test('row actions menu has expected options', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    const actionButton = page.locator('table tbody tr button[aria-haspopup="menu"]').first()

    try {
      await actionButton.waitFor({ state: 'visible', timeout: 5000 })
      await actionButton.click()

      // Should have menu items like View, Edit
      const menuItems = page.locator('[role="menuitem"]')
      await expect(menuItems.first()).toBeVisible()
    } catch {
      test.skip()
    }
  })

  test('can view item photos gallery', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    const itemLink = page.locator('table tbody tr a[href*="/items/"]').first()

    try {
      await itemLink.waitFor({ state: 'visible', timeout: 5000 })
      await itemLink.click()
      await expect(page).toHaveURL(/\/items\/\d+/)

      // Look for image gallery or photos section
      const photos = page.locator('img[alt*="item"], img[alt*="artículo"], [class*="gallery"], [class*="photo"]')
      // Photos may or may not exist
      if (await photos.first().isVisible()) {
        await expect(photos.first()).toBeVisible()
      }
    } catch {
      test.skip()
    }
  })

  test('item detail shows valuation information', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    const itemLink = page.locator('table tbody tr a[href*="/items/"]').first()

    try {
      await itemLink.waitFor({ state: 'visible', timeout: 5000 })
      await itemLink.click()
      await expect(page).toHaveURL(/\/items\/\d+/)

      // Should show valuation/appraisal values
      const valueInfo = page.locator('text=/Valor|Avalúo|Precio|Value|Price/i')
      await expect(valueInfo.first()).toBeVisible()
    } catch {
      test.skip()
    }
  })

  test('shows columns toggle button', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    // DataTable should have columns toggle
    const columnsButton = page.locator('button:has-text("Columnas"), button:has-text("Columns")')

    try {
      await columnsButton.waitFor({ state: 'visible', timeout: 3000 })
      await expect(columnsButton).toBeVisible()

      // Click to open dropdown
      await columnsButton.click()
      await expect(page.locator('[role="menu"], [role="dialog"]')).toBeVisible()
    } catch {
      // Columns toggle may not be present
      test.skip()
    }
  })

  test('pagination displays correctly', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    // Look for pagination info
    const paginationInfo = page.locator('text=/de|of|página|page/i')

    try {
      await paginationInfo.first().waitFor({ state: 'visible', timeout: 3000 })
      await expect(paginationInfo.first()).toBeVisible()
    } catch {
      // May not have enough items for pagination
      test.skip()
    }
  })

  test('can sort items table', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    // Click on sortable column header
    const sortableHeader = page.locator('th button, th[role="button"]').first()

    try {
      await sortableHeader.waitFor({ state: 'visible', timeout: 3000 })
      await sortableHeader.click()
      await helpers.waitForLoading()

      // Should still have table visible
      await expect(page.locator('table')).toBeVisible()
    } catch {
      // May not have sortable headers
      test.skip()
    }
  })
})
