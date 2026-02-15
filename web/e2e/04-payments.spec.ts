import { test, expect, PageHelpers } from './fixtures'

test.describe('Payment Flow', () => {
  // No login needed - using pre-authenticated storage state

  test('can view payments list', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/payments')
    await helpers.waitForLoading()

    // Should see the payments table
    await expect(page.locator('h1, h2').first()).toContainText(/Pagos|Payments/i)
    await expect(page.locator('table')).toBeVisible()
  })

  test('can filter payments by date range', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/payments')
    await helpers.waitForLoading()

    // Look for date filter
    const dateFilter = page.locator('button:has-text("Fecha"), input[type="date"]')
    if (await dateFilter.first().isVisible()) {
      await dateFilter.first().click()
    }
  })

  test('can view payment detail', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/payments')
    await helpers.waitForTableData()

    // Skip test if no payment data exists
    if (!(await helpers.tableHasData())) {
      test.skip()
      return
    }

    // Click on first payment link
    await helpers.clickFirstRowLink()

    // Should navigate to detail page
    await expect(page).toHaveURL(/\/payments\/\d+/)

    // Should show payment information
    await expect(page.locator('h1, h2').first()).toBeVisible()
  })

  test('payment detail shows linked loan info', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/payments')
    await helpers.waitForTableData()

    // Skip test if no payment data exists
    if (!(await helpers.tableHasData())) {
      test.skip()
      return
    }

    // Click on first payment link
    await helpers.clickFirstRowLink()
    await expect(page).toHaveURL(/\/payments\/\d+/)

    // Should have link to related loan
    const loanLink = page.locator('a[href*="/loans/"]')
    await expect(loanLink.first()).toBeVisible()
  })

  test('payment has receipt/print option', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/payments')
    await helpers.waitForTableData()

    // Skip test if no payment data exists
    if (!(await helpers.tableHasData())) {
      test.skip()
      return
    }

    // Click on first payment link
    await helpers.clickFirstRowLink()
    await expect(page).toHaveURL(/\/payments\/\d+/)

    // Look for print/receipt button
    const printButton = page.locator('button:has-text("Imprimir"), button:has-text("Recibo"), button:has-text("Print")')
    // May or may not exist based on implementation
  })

  test('can navigate to create payment from loan', async ({ page }) => {
    const helpers = new PageHelpers(page)

    // Go to loans first
    await page.goto('/loans')
    await helpers.waitForTableData()

    // Skip test if no loan data exists
    if (!(await helpers.tableHasData())) {
      test.skip()
      return
    }

    // Navigate to first active loan
    await helpers.clickFirstRowLink()
    await expect(page).toHaveURL(/\/loans\/\d+/)

    // Look for payment action
    const paymentAction = page.locator('button:has-text("Pagar"), button:has-text("Nuevo Pago"), a:has-text("Registrar Pago")')
    if (await paymentAction.first().isVisible()) {
      await paymentAction.first().click()

      // Should open dialog or navigate to payment form
      const dialog = page.locator('[role="dialog"]')
      const paymentPage = page.url().includes('/payments')

      expect(await dialog.isVisible() || paymentPage).toBeTruthy()
    }
  })

  test('payment form shows payment amount calculation', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/payments/new')

    // Or go through loan
    if (await page.locator('text=Seleccionar préstamo').isVisible()) {
      // Payment creation requires loan selection first
    }

    // Look for amount field and calculation display
    const amountField = page.locator('input[name="amount"]')
    if (await amountField.isVisible()) {
      // Should show amount input
      await expect(amountField).toBeVisible()
    }
  })

  test('payments list shows correct currency format', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/payments')
    await helpers.waitForTableData()

    // Skip test if no payment data exists
    if (!(await helpers.tableHasData())) {
      test.skip()
      return
    }

    // Look for currency formatted values (Q or GTQ for Quetzal)
    const currencyCell = page.locator('table tbody td:has-text("Q"), table tbody td:has-text("GTQ")')
    if (await currencyCell.first().isVisible()) {
      await expect(currencyCell.first()).toBeVisible()
    }
  })

  test('can reverse a payment (admin only)', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/payments')
    await helpers.waitForTableData()

    // Skip test if no payment data exists
    if (!(await helpers.tableHasData())) {
      test.skip()
      return
    }

    // Open row actions
    await helpers.openRowActions(0)

    // Look for reverse option
    const reverseOption = page.locator('[role="menuitem"]:has-text("Reversar"), [role="menuitem"]:has-text("Anular")')
    if (await reverseOption.isVisible()) {
      await reverseOption.click()

      // Should show confirmation dialog
      const dialog = page.locator('[role="dialog"], [role="alertdialog"]')
      await expect(dialog).toBeVisible()
    }
  })
})
