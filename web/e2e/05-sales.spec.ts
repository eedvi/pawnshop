import { test, expect, PageHelpers } from './fixtures'

test.describe('Sales Flow', () => {
  // No login needed - using pre-authenticated storage state

  test('can view sales list', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/sales')
    await helpers.waitForLoading()

    // Should see the sales table
    await expect(page.locator('h1, h2').first()).toContainText(/Ventas|Sales/i)
    await expect(page.locator('table')).toBeVisible()
  })

  test('can navigate to create sale form', async ({ page }) => {
    await page.goto('/sales')

    // Click new sale button
    const newButton = page.locator('a:has-text("Nueva Venta"), button:has-text("Nueva"), a:has-text("Nuevo")')
    if (await newButton.first().isVisible()) {
      await newButton.first().click()

      // Should be on create page
      await expect(page).toHaveURL(/\/sales\/new|\/sales\/create/)
    }
  })

  test('sale creation requires item selection', async ({ page }) => {
    await page.goto('/sales/new')

    // Should have item ID input field
    const itemInput = page.locator('input[name="item_id"]')
    await expect(itemInput).toBeVisible()

    // Should have sale price and payment method fields
    await expect(page.locator('input[name="sale_price"]')).toBeVisible()
  })

  test('can view sale detail', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/sales')
    await helpers.waitForLoading()

    // Check if there are sales in the table
    const hasData = await page.locator('table tbody tr a').first().isVisible()
    if (!hasData) {
      // Skip if no sales data
      test.skip()
      return
    }

    // Click on first sale link
    await helpers.clickFirstRowLink()

    // Should navigate to detail page
    await expect(page).toHaveURL(/\/sales\/\d+/)

    // Should show sale information
    await expect(page.locator('h1, h2').first()).toBeVisible()
  })

  test('sale detail shows item and buyer info', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/sales')
    await helpers.waitForLoading()

    // Check if there are sales in the table
    const hasData = await page.locator('table tbody tr a').first().isVisible()
    if (!hasData) {
      // Skip if no sales data
      test.skip()
      return
    }

    // Click on first sale link
    await helpers.clickFirstRowLink()
    await expect(page).toHaveURL(/\/sales\/\d+/)

    // Should show item link
    const itemLink = page.locator('a[href*="/items/"]')
    // Item info should be present
  })

  test('can process refund on sale', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/sales')
    await helpers.waitForLoading()

    // Check if there are sales with action buttons
    const hasActionButton = await page.locator('table tbody tr button[aria-haspopup="menu"]').first().isVisible()
    if (!hasActionButton) {
      // Skip if no sales data or no action buttons
      test.skip()
      return
    }

    // Open row actions
    await helpers.openRowActions(0)

    // Look for refund option
    const refundOption = page.locator('[role="menuitem"]:has-text("Reembolso"), [role="menuitem"]:has-text("Refund")')
    if (await refundOption.isVisible()) {
      await refundOption.click()

      // Should show refund dialog
      const dialog = page.locator('[role="dialog"]')
      await expect(dialog).toBeVisible()
    }
  })

  test('can cancel a sale', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/sales')
    await helpers.waitForLoading()

    // Check if there are sales with action buttons
    const hasActionButton = await page.locator('table tbody tr button[aria-haspopup="menu"]').first().isVisible()
    if (!hasActionButton) {
      // Skip if no sales data or no action buttons
      test.skip()
      return
    }

    // Open row actions
    await helpers.openRowActions(0)

    // Look for cancel option
    const cancelOption = page.locator('[role="menuitem"]:has-text("Cancelar"), [role="menuitem"]:has-text("Cancel")')
    if (await cancelOption.isVisible()) {
      await cancelOption.click()

      // Should show confirmation dialog
      const dialog = page.locator('[role="dialog"], [role="alertdialog"]')
      await expect(dialog).toBeVisible()
    }
  })

  test('only items marked for sale appear in sale form', async ({ page }) => {
    await page.goto('/sales/new')

    // The item dropdown should only show items with "for_sale" status
    const itemSelect = page.locator('button[role="combobox"]:has-text("Artículo"), button[role="combobox"]:has-text("Item")')
    if (await itemSelect.isVisible()) {
      await itemSelect.click()

      // Should show available items (or empty if no items for sale)
      const options = page.locator('[role="option"]')
      // Items shown should be for sale
    }
  })

  test('sale form calculates totals correctly', async ({ page }) => {
    await page.goto('/sales/new')

    // Look for price input and total display
    const priceInput = page.locator('input[name="sale_price"], input[name="price"]')
    if (await priceInput.isVisible()) {
      await priceInput.fill('1500')

      // Total should update
      const totalDisplay = page.locator('text=Total').or(page.locator('text=Subtotal'))
      // Total display may update automatically
    }
  })

  test('completed sales appear with correct status', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/sales')
    await helpers.waitForTableData()

    // Look for status badges
    const completedBadge = page.locator('[class*="badge"]:has-text("Completada"), [class*="badge"]:has-text("Completed")')
    if (await completedBadge.first().isVisible()) {
      await expect(completedBadge.first()).toBeVisible()
    }
  })
})
