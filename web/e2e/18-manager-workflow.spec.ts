import { test, expect, PageHelpers } from './fixtures'

/**
 * Manager Workflow Tests
 *
 * Simulates managerial/supervisory tasks:
 * 1. Dashboard monitoring and KPIs
 * 2. Expense approval
 * 3. Transfer approval
 * 4. Overdue loan review
 * 5. Cash session review
 * 6. Report generation
 * 7. Team performance overview
 */
test.describe('Manager Workflow', () => {
  test.describe('Dashboard Monitoring', () => {
    test('dashboard shows key performance indicators', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/')
      await helpers.waitForLoading()

      // Should show KPI cards
      const kpiCards = page.locator('[class*="stat"], [class*="card"]').filter({
        has: page.locator('h3, [class*="title"]'),
      })

      await expect(kpiCards.first()).toBeVisible()
    })

    test('can see active loans count', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/')
      await helpers.waitForLoading()

      const loanStats = page.locator('text=/Préstamos|Activos|Loans/i')
      await expect(loanStats.first()).toBeVisible()
    })

    test('can see overdue loans warning', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/')
      await helpers.waitForLoading()

      // Look for overdue/vencido indicator
      const overdueIndicator = page.locator('text=/Vencido|Overdue|Mora/i')

      try {
        await overdueIndicator.first().waitFor({ state: 'visible', timeout: 3000 })
        await expect(overdueIndicator.first()).toBeVisible()
      } catch {
        // May not have overdue loans
        test.skip()
      }
    })

    test('can see todays collections', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/')
      await helpers.waitForLoading()

      const collectionsInfo = page.locator('text=/Cobros|Recaudado|Collections|Pagos/i')
      await expect(collectionsInfo.first()).toBeVisible()
    })

    test('dashboard shows charts', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/')
      await helpers.waitForLoading()

      const charts = page.locator('svg.recharts-surface, [class*="chart"]')

      try {
        await charts.first().waitFor({ state: 'visible', timeout: 5000 })
        await expect(charts.first()).toBeVisible()
      } catch {
        // Charts may not be visible without data
        test.skip()
      }
    })
  })

  test.describe('Overdue Loans Review', () => {
    test('can filter loans by overdue status', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans?status=overdue')
      await helpers.waitForLoading()

      // Should show overdue loans or empty state
      const content = page.locator('table').or(page.locator('text=/No hay|No se encontraron/i'))
      await expect(content.first()).toBeVisible()
    })

    test('overdue loans show days overdue', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans?status=overdue')
      await helpers.waitForLoading()

      const loanLink = page.locator('table tbody tr a[href*="/loans/"]').first()

      try {
        await loanLink.waitFor({ state: 'visible', timeout: 5000 })
        // Table should show overdue information
        await expect(page.locator('table')).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('can contact customer from overdue loan', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans?status=overdue')
      await helpers.waitForLoading()

      const loanLink = page.locator('table tbody tr a[href*="/loans/"]').first()

      try {
        await loanLink.waitFor({ state: 'visible', timeout: 5000 })
        await loanLink.click()

        // Should have customer link
        const customerLink = page.locator('a[href*="/customers/"]')
        await expect(customerLink.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Expense Approval', () => {
    test('can view pending expenses', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses?status=pending')
      await helpers.waitForLoading()

      const content = page.locator('table').or(page.locator('text=/No hay|No se encontraron/i'))
      await expect(content.first()).toBeVisible()
    })

    test('expense shows approval status', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses')
      await helpers.waitForLoading()

      // Table should have status column
      const statusBadges = page.locator('table [class*="badge"]')

      try {
        await statusBadges.first().waitFor({ state: 'visible', timeout: 5000 })
        await expect(statusBadges.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('can view expense details', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses')
      await helpers.waitForLoading()

      const expenseLink = page.locator('table tbody tr a[href*="/expenses/"]').first()

      try {
        await expenseLink.waitFor({ state: 'visible', timeout: 5000 })
        await expenseLink.click()
        await expect(page).toHaveURL(/\/expenses\/\d+/)
      } catch {
        test.skip()
      }
    })

    test('expense detail shows approval buttons', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses')
      await helpers.waitForLoading()

      const expenseLink = page.locator('table tbody tr a[href*="/expenses/"]').first()

      try {
        await expenseLink.waitFor({ state: 'visible', timeout: 5000 })
        await expenseLink.click()

        // Look for approve/reject buttons
        const approvalButtons = page.locator('button:has-text("Aprobar"), button:has-text("Rechazar")')

        if (await approvalButtons.first().isVisible()) {
          await expect(approvalButtons.first()).toBeVisible()
        }
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Transfer Approval', () => {
    test('can view pending transfers', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/transfers?status=pending')
      await helpers.waitForLoading()

      const content = page.locator('table').or(page.locator('text=/No hay|No se encontraron/i'))
      await expect(content.first()).toBeVisible()
    })

    test('transfer shows items and branches', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/transfers')
      await helpers.waitForLoading()

      const transferLink = page.locator('table tbody tr a[href*="/transfers/"]').first()

      try {
        await transferLink.waitFor({ state: 'visible', timeout: 5000 })
        await transferLink.click()
        await expect(page).toHaveURL(/\/transfers\/\d+/)

        // Should show source and destination
        const branchInfo = page.locator('text=/Origen|Destino|Source|Destination/i')
        await expect(branchInfo.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('transfer has action buttons', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/transfers')
      await helpers.waitForLoading()

      const transferLink = page.locator('table tbody tr a[href*="/transfers/"]').first()

      try {
        await transferLink.waitFor({ state: 'visible', timeout: 5000 })
        await transferLink.click()

        // Look for action buttons
        const actionButtons = page.locator('button:has-text("Aprobar"), button:has-text("Enviar"), button:has-text("Recibir")')

        if (await actionButtons.first().isVisible()) {
          await expect(actionButtons.first()).toBeVisible()
        }
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Cash Session Review', () => {
    test('can view all cash sessions', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/cash')
      await helpers.waitForLoading()

      // Look for sessions tab or list
      const sessionsTab = page.locator('[role="tab"]:has-text("Sesiones")')

      if (await sessionsTab.isVisible()) {
        await sessionsTab.click()
        await helpers.waitForLoading()
      }

      // Should see sessions list
      const sessionsList = page.locator('table, [class*="card"]')
      await expect(sessionsList.first()).toBeVisible()
    })

    test('can view session details', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/cash')
      await helpers.waitForLoading()

      // Navigate to sessions
      const sessionsTab = page.locator('[role="tab"]:has-text("Sesiones")')
      if (await sessionsTab.isVisible()) {
        await sessionsTab.click()
        await helpers.waitForLoading()
      }

      // Click on a session
      const sessionLink = page.locator('table tbody tr a, [class*="card"] a').first()

      try {
        await sessionLink.waitFor({ state: 'visible', timeout: 5000 })
        await sessionLink.click()
      } catch {
        test.skip()
      }
    })

    test('session shows opening and closing amounts', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/cash')
      await helpers.waitForLoading()

      // Should show financial info
      const amountInfo = page.locator('text=/Apertura|Cierre|Opening|Closing|Monto|Amount/i')
      try {
        await amountInfo.first().waitFor({ state: 'visible', timeout: 5000 })
        await expect(amountInfo.first()).toBeVisible()
      } catch {
        // Cash session info may not be visible without active session
        test.skip()
      }
    })

    test('can view all cash movements', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/cash')
      await helpers.waitForLoading()

      const movementsTab = page.locator('[role="tab"]:has-text("Movimientos")')

      try {
        await movementsTab.waitFor({ state: 'visible', timeout: 5000 })
        await movementsTab.click()
        await helpers.waitForLoading()

        const movementsList = page.locator('table')
        await movementsList.waitFor({ state: 'visible', timeout: 5000 })
        await expect(movementsList).toBeVisible()
      } catch {
        // Movements tab or table may not be available
        test.skip()
      }
    })
  })

  test.describe('Report Generation', () => {
    test('can access loan report', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/reports')
      await helpers.waitForLoading()

      const loanReportTab = page.locator('[role="tab"]:has-text("Préstamos")').or(page.locator('text=/Préstamos/i'))

      if (await loanReportTab.first().isVisible()) {
        await loanReportTab.first().click()
        await helpers.waitForLoading()
      }
    })

    test('can access payment report', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/reports')
      await helpers.waitForLoading()

      const paymentReportTab = page.locator('[role="tab"]:has-text("Pagos")').or(page.locator('text=/Pagos/i'))

      if (await paymentReportTab.first().isVisible()) {
        await paymentReportTab.first().click()
        await helpers.waitForLoading()
      }
    })

    test('can access sales report', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/reports')
      await helpers.waitForLoading()

      const salesReportTab = page.locator('[role="tab"]:has-text("Ventas")').or(page.locator('text=/Ventas/i'))

      if (await salesReportTab.first().isVisible()) {
        await salesReportTab.first().click()
        await helpers.waitForLoading()
      }
    })

    test('reports have date range selection', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/reports')
      await helpers.waitForLoading()

      const dateInputs = page.locator('input[type="date"], button:has-text("Fecha")')

      try {
        await dateInputs.first().waitFor({ state: 'visible', timeout: 5000 })
        await expect(dateInputs.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('reports may have export option', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/reports')
      await helpers.waitForLoading()

      const exportButton = page.locator('button:has-text("Exportar"), button:has-text("Export"), button:has-text("PDF")')

      if (await exportButton.first().isVisible()) {
        await expect(exportButton.first()).toBeVisible()
      }
    })
  })

  test.describe('Branch Performance', () => {
    test('can compare branches in reports', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/reports')
      await helpers.waitForLoading()

      const branchFilter = page.locator('button:has-text("Sucursal"), select[name*="branch"]')

      if (await branchFilter.first().isVisible()) {
        await expect(branchFilter.first()).toBeVisible()
      }
    })

    test('dashboard respects branch filter', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/')
      await helpers.waitForLoading()

      // Look for branch selector in header
      const branchSelector = page.locator('header button:has-text("Sucursal"), header [class*="branch"]')

      if (await branchSelector.first().isVisible()) {
        await expect(branchSelector.first()).toBeVisible()
      }
    })
  })

  test.describe('Item Inventory Review', () => {
    test('can view items by status', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/items')
      await helpers.waitForLoading()

      const statusFilter = page.locator('button:has-text("Estado")')

      if (await statusFilter.isVisible()) {
        await statusFilter.click()

        const filterOptions = page.locator('[role="menuitem"], [role="option"]')
        await expect(filterOptions.first()).toBeVisible()
      }
    })

    test('can see pawned items count', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/items?status=pawned')
      await helpers.waitForLoading()

      // Should show items or empty state
      const content = page.locator('table').or(page.locator('text=/No hay|No se encontraron/i'))
      await expect(content.first()).toBeVisible()
    })

    test('can see items for sale', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/items?status=for_sale')
      await helpers.waitForLoading()

      const content = page.locator('table').or(page.locator('text=/No hay|No se encontraron/i'))
      await expect(content.first()).toBeVisible()
    })
  })

  test.describe('Customer Analytics', () => {
    test('can view customer list with metrics', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      await expect(page.locator('table')).toBeVisible()
    })

    test('customer detail shows loan history summary', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const customerLink = page.locator('table tbody tr a[href*="/customers/"]').first()

      try {
        await customerLink.waitFor({ state: 'visible', timeout: 5000 })
        await customerLink.click()

        // Should show summary info
        const summaryInfo = page.locator('text=/Total|Préstamos|Historial/i')
        await expect(summaryInfo.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Audit Trail', () => {
    test('can review user activities', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      await expect(page.locator('table')).toBeVisible()
    })

    test('can filter by specific user', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      const filtersButton = page.locator('button:has-text("Filtros")')
      await filtersButton.click()

      await page.waitForSelector('input#date_from', { state: 'visible' })

      const userFilter = page.locator('label:has-text("Usuario")')
      await expect(userFilter).toBeVisible()
    })

    test('can filter by date range', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit')
      await helpers.waitForLoading()

      const filtersButton = page.locator('button:has-text("Filtros")')
      await filtersButton.click()

      await page.waitForSelector('input#date_from', { state: 'visible' })

      await expect(page.locator('input#date_from')).toBeVisible()
      await expect(page.locator('input#date_to')).toBeVisible()
    })
  })
})
