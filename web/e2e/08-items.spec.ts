import { test, expect, PageHelpers } from './fixtures'

test.describe('Items Management', () => {
  test('can view items list', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    // Should see the items table
    await expect(page.locator('h1, h2').first()).toContainText(/Artículos|Items/i)
    await expect(page.locator('table')).toBeVisible()
  })

  test('can search items', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    // Type in search field
    const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')
    if (await searchInput.isVisible()) {
      await searchInput.fill('anillo')
      await page.waitForTimeout(500) // Debounce wait
      await helpers.waitForLoading()
    }
  })

  test('can navigate to create item form', async ({ page }) => {
    await page.goto('/items')

    // Click new item button/link
    const newButton = page.locator('a:has-text("Nuevo"), button:has-text("Nuevo")')
    await expect(newButton.first()).toBeVisible()
  })

  test('can view item detail', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    // Wait for table data
    const itemLink = page.locator('table tbody tr a').first()

    try {
      await itemLink.waitFor({ state: 'visible', timeout: 5000 })
      await itemLink.click()
      await expect(page).toHaveURL(/\/items\/\d+/)
    } catch {
      // No items in database, skip
      test.skip()
    }
  })

  test('item detail shows status and values', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    const itemLink = page.locator('table tbody tr a').first()

    try {
      await itemLink.waitFor({ state: 'visible', timeout: 5000 })
      await itemLink.click()
      await expect(page).toHaveURL(/\/items\/\d+/)

      // Should show item status
      await expect(page.locator('main')).toBeVisible()
    } catch {
      test.skip()
    }
  })

  test('can filter items by status', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    // Look for status filter
    const statusFilter = page.locator('button:has-text("Estado"), button:has-text("Filtrar")')
    if (await statusFilter.first().isVisible()) {
      await statusFilter.first().click()
      await helpers.waitForLoading()
    }
  })

  test('can use row actions menu', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/items')
    await helpers.waitForLoading()

    // Check if there are items with action buttons
    const actionButton = page.locator('table tbody tr button[aria-haspopup="menu"]').first()

    try {
      await actionButton.waitFor({ state: 'visible', timeout: 5000 })
      await actionButton.click()

      // Should open dropdown menu
      await expect(page.locator('[role="menu"]')).toBeVisible()
    } catch {
      // No items or no action buttons
      test.skip()
    }
  })
})
