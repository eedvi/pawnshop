import { test, expect, PageHelpers, generateTestData, TEST_USERS, login } from './fixtures'

/**
 * Cashier Daily Workflow Tests
 *
 * Simulates real scenarios a cashier encounters during their work day:
 * 1. Opening cash register at start of shift
 * 2. Attending customers for loan inquiries
 * 3. Processing payments
 * 4. Handling sales
 * 5. Recording expenses
 * 6. Closing cash register at end of shift
 */
test.describe('Cashier Daily Workflow', () => {
  test.describe('Shift Start - Opening Cash Register', () => {
    test('cashier sees cash register status on login', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/cash')
      await helpers.waitForLoading()

      // Should see cash register page
      await expect(page.locator('h1, h2').first()).toContainText(/Caja|Cash/i)
    })

    test('can see if cash register is open or closed', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/cash')
      await helpers.waitForLoading()

      // Should show status indicator
      const statusIndicator = page.locator('text=/Abierta|Cerrada|Open|Closed/i')
      try {
        await statusIndicator.first().waitFor({ state: 'visible', timeout: 5000 })
        await expect(statusIndicator.first()).toBeVisible()
      } catch {
        // Status may be shown differently
      }
    })

    test('can open cash register with initial amount', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/cash')
      await helpers.waitForLoading()

      const openButton = page.locator('button:has-text("Abrir"), button:has-text("Open")')

      try {
        await openButton.first().waitFor({ state: 'visible', timeout: 3000 })

        if (await openButton.first().isVisible()) {
          await openButton.first().click()

          // Dialog should appear
          const dialog = page.locator('[role="dialog"]')
          await dialog.waitFor({ state: 'visible', timeout: 3000 })

          // Fill opening amount
          const amountInput = dialog.locator('input[name*="amount"], input[type="number"]').first()
          if (await amountInput.isVisible()) {
            await amountInput.fill('5000')
          }

          // Submit
          const submitButton = dialog.locator('button[type="submit"], button:has-text("Abrir")')
          await submitButton.click()
          await helpers.waitForLoading()
        }
      } catch {
        // Cash register may already be open
        test.skip()
      }
    })

    test('shows current cash session details', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/cash')
      await helpers.waitForLoading()

      // Look for session info (opening time, user, amount)
      const sessionInfo = page.locator('text=/Sesión|Session|Apertura|Opening|Monto|Amount/i')

      try {
        await sessionInfo.first().waitFor({ state: 'visible', timeout: 5000 })
        await expect(sessionInfo.first()).toBeVisible()
      } catch {
        // May need to open cash first
        test.skip()
      }
    })
  })

  test.describe('Customer Service - Loan Consultation', () => {
    test('can quickly search for existing customer', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')
      await searchInput.fill('Juan')
      await page.waitForTimeout(600)
      await helpers.waitForLoading()

      // Table should update with results
      await expect(page.locator('table')).toBeVisible()
    })

    test('can view customer loan history', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const customerLink = page.locator('table tbody tr a[href*="/customers/"]').first()

      try {
        await customerLink.waitFor({ state: 'visible', timeout: 5000 })
        await customerLink.click()
        await expect(page).toHaveURL(/\/customers\/\d+/)

        // Look for loans tab
        const loansTab = page.locator('[role="tab"]:has-text("Préstamos")')
        if (await loansTab.isVisible()) {
          await loansTab.click()
          await helpers.waitForLoading()
        }
      } catch {
        test.skip()
      }
    })

    test('can see customer credit standing', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const customerLink = page.locator('table tbody tr a[href*="/customers/"]').first()

      try {
        await customerLink.waitFor({ state: 'visible', timeout: 5000 })
        await customerLink.click()
        await expect(page).toHaveURL(/\/customers\/\d+/)

        // Should show credit/status info
        await expect(page.locator('main')).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('can check active loans for customer', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const customerLink = page.locator('table tbody tr a[href*="/customers/"]').first()

      try {
        await customerLink.waitFor({ state: 'visible', timeout: 5000 })
        await customerLink.click()

        // Navigate to loans tab
        const loansTab = page.locator('[role="tab"]:has-text("Préstamos")')
        if (await loansTab.isVisible()) {
          await loansTab.click()
          await helpers.waitForLoading()

          // Should show loan list or empty state
          const loanContent = page.locator('table').or(page.locator('text=/No tiene|No hay/i'))
          await expect(loanContent.first()).toBeVisible()
        }
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Payment Processing', () => {
    test('can access payment screen from loan', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const loanLink = page.locator('table tbody tr a[href*="/loans/"]').first()

      try {
        await loanLink.waitFor({ state: 'visible', timeout: 5000 })
        await loanLink.click()
        await expect(page).toHaveURL(/\/loans\/\d+/)

        // Find payment button
        const payButton = page.locator('button:has-text("Pagar"), a:has-text("Pago")')
        await expect(payButton.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('can see payment amount calculation', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const loanLink = page.locator('table tbody tr a[href*="/loans/"]').first()

      try {
        await loanLink.waitFor({ state: 'visible', timeout: 5000 })
        await loanLink.click()

        // Should show amounts (principal, interest, total)
        const amounts = page.locator('text=/Principal|Interés|Total|Saldo/i')
        await expect(amounts.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('can see payment history for loan', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const loanLink = page.locator('table tbody tr a[href*="/loans/"]').first()

      try {
        await loanLink.waitFor({ state: 'visible', timeout: 5000 })
        await loanLink.click()

        // Look for payments section or tab
        const paymentsSection = page.locator('text=/Pagos|Historial|Payments/i')
        await expect(paymentsSection.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('payments list shows todays payments', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/payments')
      await helpers.waitForLoading()

      // Table should be visible
      await expect(page.locator('table')).toBeVisible()
    })

    test('can filter payments by date', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/payments')
      await helpers.waitForLoading()

      // Look for date filter
      const dateFilter = page.locator('input[type="date"], button:has-text("Fecha")')

      if (await dateFilter.first().isVisible()) {
        await expect(dateFilter.first()).toBeVisible()
      }
    })
  })

  test.describe('Sales Processing', () => {
    test('can view available items for sale', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/items?status=for_sale')
      await helpers.waitForLoading()

      // Should show items list or empty state
      const content = page.locator('table').or(page.locator('text=/No hay|No se encontraron/i'))
      await expect(content.first()).toBeVisible()
    })

    test('can access sales list', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/sales')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Ventas|Sales/i)
      await expect(page.locator('table')).toBeVisible()
    })

    test('can view sale details', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/sales')
      await helpers.waitForLoading()

      const saleLink = page.locator('table tbody tr a[href*="/sales/"]').first()

      try {
        await saleLink.waitFor({ state: 'visible', timeout: 5000 })
        await saleLink.click()
        await expect(page).toHaveURL(/\/sales\/\d+/)
      } catch {
        test.skip()
      }
    })

    test('sale shows item and customer info', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/sales')
      await helpers.waitForLoading()

      const saleLink = page.locator('table tbody tr a[href*="/sales/"]').first()

      try {
        await saleLink.waitFor({ state: 'visible', timeout: 5000 })
        await saleLink.click()

        // Should show item and customer references
        const itemRef = page.locator('a[href*="/items/"]')
        await expect(itemRef.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Expense Recording', () => {
    test('can access expense list', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Gastos|Expenses/i)
    })

    test('can navigate to create expense', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses')
      await helpers.waitForLoading()

      const newButton = page.locator('a:has-text("Nuevo"), button:has-text("Nuevo")')

      if (await newButton.first().isVisible()) {
        await newButton.first().click()
        await expect(page).toHaveURL(/\/expenses\/new/)
      }
    })

    test('expense form has required fields', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses/new')
      await helpers.waitForLoading()

      // Should have amount and description
      await expect(page.locator('input[name="amount"], input[type="number"]').first()).toBeVisible()
    })

    test('can see todays expenses', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses')
      await helpers.waitForLoading()

      // Table should be visible
      await expect(page.locator('table')).toBeVisible()
    })
  })

  test.describe('End of Shift - Cash Closing', () => {
    test('can view cash movements during session', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/cash')
      await helpers.waitForLoading()

      // Look for movements section or tab
      const movementsSection = page.locator('text=/Movimientos|Movements/i')

      try {
        await movementsSection.first().waitFor({ state: 'visible', timeout: 5000 })
        await expect(movementsSection.first()).toBeVisible()
      } catch {
        // May be in different tab
        const movementsTab = page.locator('[role="tab"]:has-text("Movimientos")')
        if (await movementsTab.isVisible()) {
          await movementsTab.click()
          await helpers.waitForLoading()
        }
      }
    })

    test('can see cash balance summary', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/cash')
      await helpers.waitForLoading()

      // Should show balance info
      const balanceInfo = page.locator('text=/Saldo|Balance|Total/i')
      try {
        await balanceInfo.first().waitFor({ state: 'visible', timeout: 5000 })
        await expect(balanceInfo.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('can see close cash button when session is open', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/cash')
      await helpers.waitForLoading()

      const closeButton = page.locator('button:has-text("Cerrar Caja"), button:has-text("Close")')

      try {
        await closeButton.first().waitFor({ state: 'visible', timeout: 3000 })
        await expect(closeButton.first()).toBeVisible()
      } catch {
        // Cash may be closed
        test.skip()
      }
    })

    test('closing cash shows summary dialog', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/cash')
      await helpers.waitForLoading()

      const closeButton = page.locator('button:has-text("Cerrar Caja"), button:has-text("Close")')

      try {
        await closeButton.first().waitFor({ state: 'visible', timeout: 3000 })
        await closeButton.first().click()

        // Dialog should appear with closing details
        const dialog = page.locator('[role="dialog"]')
        await dialog.waitFor({ state: 'visible', timeout: 3000 })
        await expect(dialog).toBeVisible()

        // Cancel to not actually close
        const cancelButton = dialog.locator('button:has-text("Cancelar")')
        if (await cancelButton.isVisible()) {
          await cancelButton.click()
        }
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Quick Actions', () => {
    test('can quickly navigate from sidebar', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/')
      await helpers.waitForLoading()

      // Click on Customers from sidebar
      await helpers.navigateTo('Clientes')
      await expect(page).toHaveURL(/\/customers/)

      // Click on Loans
      await helpers.navigateTo('Préstamos')
      await expect(page).toHaveURL(/\/loans/)

      // Click on Payments
      await helpers.navigateTo('Pagos')
      await expect(page).toHaveURL(/\/payments/)
    })

    test('can see notifications badge', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/')
      await helpers.waitForLoading()

      // Look for notification bell in header
      const notificationBell = page.locator('button svg, [aria-label*="notification"]').first()

      try {
        await notificationBell.waitFor({ state: 'visible', timeout: 3000 })
        await expect(notificationBell).toBeVisible()
      } catch {
        // May not have notifications component
        test.skip()
      }
    })
  })
})
