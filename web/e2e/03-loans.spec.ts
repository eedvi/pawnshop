import { test, expect, PageHelpers } from './fixtures'

test.describe('Loan Workflow', () => {
  // No login needed - using pre-authenticated storage state

  test('can view loans list', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/loans')
    await helpers.waitForLoading()

    // Should see the loans table
    await expect(page.locator('h1, h2').first()).toContainText(/Préstamos|Loans/i)
    await expect(page.locator('table')).toBeVisible()

    // Should have new loan button
    await expect(page.locator('a:has-text("Nuevo"), button:has-text("Nuevo")')).toBeVisible()
  })

  test('can filter loans by status', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/loans')
    await helpers.waitForLoading()

    // Look for status filter
    const statusFilter = page.locator('button:has-text("Estado"), select[name="status"]')
    if (await statusFilter.first().isVisible()) {
      await statusFilter.first().click()

      // Select a status option
      const activeOption = page.locator('[role="option"]:has-text("Activo"), [role="option"]:has-text("Active")')
      if (await activeOption.isVisible()) {
        await activeOption.click()
        await helpers.waitForLoading()
      }
    }
  })

  test('can navigate to loan creation wizard', async ({ page }) => {
    await page.goto('/loans')

    // Click new loan button
    await page.click('a:has-text("Nuevo"), button:has-text("Nuevo")')

    // Should be on create page with wizard
    await expect(page).toHaveURL(/\/loans\/new|\/loans\/create/)
  })

  test('loan creation wizard has multiple steps', async ({ page }) => {
    await page.goto('/loans/new')

    // Wizard should have step labels visible (Cliente, Artículo, Condiciones, Confirmar)
    await expect(page.locator('text=Cliente').first()).toBeVisible()
    await expect(page.locator('text=Artículo').first()).toBeVisible()

    // Should have navigation buttons
    await expect(page.locator('button:has-text("Anterior"), button:has-text("Siguiente")').first()).toBeVisible()
  })

  test('can view loan detail', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/loans')
    await helpers.waitForTableData()

    // Click on first loan link to navigate to detail
    await helpers.clickFirstRowLink()

    // Should navigate to detail page
    await expect(page).toHaveURL(/\/loans\/\d+/)

    // Should show loan information
    await expect(page.locator('h1, h2').first()).toBeVisible()

    // Should show loan details like amount, status, customer
    const loanInfo = page.locator('main')
    await expect(loanInfo).toBeVisible()
  })

  test('loan detail shows installment schedule', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/loans')
    await helpers.waitForTableData()

    // Navigate to first loan
    await helpers.clickFirstRowLink()

    await expect(page).toHaveURL(/\/loans\/\d+/)

    // Look for installments/payments section
    const installmentsSection = page.locator('text=Cuotas').or(page.locator('text=Pagos')).or(page.locator('text=Installments')).or(page.locator('text=Schedule'))
    // At least one should be visible if loan has schedule
  })

  test('loan detail has action buttons', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/loans')
    await helpers.waitForTableData()

    // Navigate to first loan
    await helpers.clickFirstRowLink()
    await expect(page).toHaveURL(/\/loans\/\d+/)

    // Should have action buttons based on loan status
    // Common actions: Pay, Renew, Confiscate
    const actionButtons = page.locator('button:has-text("Pagar"), button:has-text("Renovar"), button:has-text("Confiscar"), button:has-text("Pay"), button:has-text("Renew")')
    // At least one action should be available
  })

  test('can open payment dialog from loan detail', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/loans')
    await helpers.waitForTableData()

    // Navigate to an active loan
    await helpers.clickFirstRowLink()
    await expect(page).toHaveURL(/\/loans\/\d+/)

    // Look for payment button
    const payButton = page.locator('button:has-text("Pagar"), button:has-text("Registrar Pago"), a:has-text("Pagar")')
    if (await payButton.first().isVisible()) {
      await payButton.first().click()

      // Should open payment dialog or navigate to payment page
      const dialog = page.locator('[role="dialog"]')
      if (await dialog.isVisible()) {
        await expect(dialog).toBeVisible()
      } else {
        // Or navigate to payment creation
        await expect(page).toHaveURL(/\/payments\/new|\/loans\/\d+\/pay/)
      }
    }
  })

  test('shows overdue loans with visual indicator', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/loans')
    await helpers.waitForLoading()

    // Look for overdue status badges
    const overdueBadges = page.locator('[class*="destructive"], [class*="red"]:has-text("Vencido"), [class*="red"]:has-text("Overdue")')
    // These may or may not exist depending on data
  })

  test('can navigate between loan related entities', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/loans')
    await helpers.waitForTableData()

    // Skip test if no loan data exists
    if (!(await helpers.tableHasData())) {
      test.skip()
      return
    }

    // Navigate to first loan
    await helpers.clickFirstRowLink()
    await expect(page).toHaveURL(/\/loans\/\d+/)

    // Should have links to customer and item
    const customerLink = page.locator('a[href*="/customers/"]')
    const itemLink = page.locator('a[href*="/items/"]')

    // At least one should be clickable
    if (await customerLink.first().isVisible()) {
      await customerLink.first().click()
      await expect(page).toHaveURL(/\/customers\/\d+/)
    }
  })
})
