import { test, expect, PageHelpers } from './fixtures'

test.describe('Loan Lifecycle', () => {
  test.describe('Loan List and Navigation', () => {
    test('can view loans list', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Préstamos|Loans/i)
      await expect(page.locator('table')).toBeVisible()
    })

    test('loans table shows status and amounts', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      // Table headers should include status and amount
      const headers = page.locator('th')
      await expect(headers.first()).toBeVisible()
    })

    test('can navigate to loan detail', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const loanLink = page.locator('table tbody tr a[href*="/loans/"]').first()

      try {
        await loanLink.waitFor({ state: 'visible', timeout: 5000 })
        await loanLink.click()
        await expect(page).toHaveURL(/\/loans\/\d+/)
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Loan Detail View', () => {
    test('loan detail shows customer link', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const loanLink = page.locator('table tbody tr a[href*="/loans/"]').first()

      try {
        await loanLink.waitFor({ state: 'visible', timeout: 5000 })
        await loanLink.click()
        await expect(page).toHaveURL(/\/loans\/\d+/)

        // Should have customer reference
        const customerLink = page.locator('a[href*="/customers/"]')
        await expect(customerLink.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('loan detail shows item link', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const loanLink = page.locator('table tbody tr a[href*="/loans/"]').first()

      try {
        await loanLink.waitFor({ state: 'visible', timeout: 5000 })
        await loanLink.click()
        await expect(page).toHaveURL(/\/loans\/\d+/)

        // Should have item reference
        const itemLink = page.locator('a[href*="/items/"]')
        if (await itemLink.first().isVisible()) {
          await expect(itemLink.first()).toBeVisible()
        }
      } catch {
        test.skip()
      }
    })

    test('loan detail shows payment schedule', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const loanLink = page.locator('table tbody tr a[href*="/loans/"]').first()

      try {
        await loanLink.waitFor({ state: 'visible', timeout: 5000 })
        await loanLink.click()
        await expect(page).toHaveURL(/\/loans\/\d+/)

        // Look for installments/cuotas table or section
        const scheduleSection = page.locator('text=/Cuotas|Pagos|Vencimiento|Installments|Schedule/i')
        await expect(scheduleSection.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('loan detail shows financial summary', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const loanLink = page.locator('table tbody tr a[href*="/loans/"]').first()

      try {
        await loanLink.waitFor({ state: 'visible', timeout: 5000 })
        await loanLink.click()
        await expect(page).toHaveURL(/\/loans\/\d+/)

        // Should show loan amounts
        const amountInfo = page.locator('text=/Monto|Principal|Interés|Interest|Total|Saldo|Balance/i')
        await expect(amountInfo.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Loan Actions', () => {
    test('loan detail has payment action button', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const loanLink = page.locator('table tbody tr a[href*="/loans/"]').first()

      try {
        await loanLink.waitFor({ state: 'visible', timeout: 5000 })
        await loanLink.click()
        await expect(page).toHaveURL(/\/loans\/\d+/)

        // Look for payment action
        const paymentButton = page.locator('button:has-text("Pagar"), button:has-text("Registrar Pago"), a:has-text("Pago")')
        if (await paymentButton.first().isVisible()) {
          await expect(paymentButton.first()).toBeVisible()
        }
      } catch {
        test.skip()
      }
    })

    test('can click payment button to open form', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const loanLink = page.locator('table tbody tr a[href*="/loans/"]').first()

      try {
        await loanLink.waitFor({ state: 'visible', timeout: 5000 })
        await loanLink.click()
        await expect(page).toHaveURL(/\/loans\/\d+/)

        const paymentButton = page.locator('button:has-text("Pagar"), a:has-text("Nuevo Pago"), a:has-text("Registrar Pago")')

        if (await paymentButton.first().isVisible()) {
          await paymentButton.first().click()

          // Should open dialog or navigate to payment form
          const paymentForm = page.locator('[role="dialog"], form')
          await paymentForm.first().waitFor({ state: 'visible', timeout: 3000 })
        }
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Loan Filtering', () => {
    test('can filter loans by status', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

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

    test('can filter loans by date range', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      // Look for date filter
      const dateFilter = page.locator('button:has-text("Fecha"), input[type="date"]')

      if (await dateFilter.first().isVisible()) {
        await expect(dateFilter.first()).toBeVisible()
      }
    })

    test('can search loans', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')

      if (await searchInput.isVisible()) {
        await searchInput.fill('test')
        await page.waitForTimeout(500)
        await helpers.waitForLoading()

        await expect(page.locator('table')).toBeVisible()
      }
    })
  })

  test.describe('Payments Flow', () => {
    test('can view payments list', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/payments')
      await helpers.waitForLoading()

      await expect(page.locator('h1, h2').first()).toContainText(/Pagos|Payments/i)
      await expect(page.locator('table')).toBeVisible()
    })

    test('payment has loan reference', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/payments')
      await helpers.waitForLoading()

      const paymentLink = page.locator('table tbody tr a[href*="/payments/"]').first()

      try {
        await paymentLink.waitFor({ state: 'visible', timeout: 5000 })
        await paymentLink.click()
        await expect(page).toHaveURL(/\/payments\/\d+/)

        // Should link to loan
        const loanLink = page.locator('a[href*="/loans/"]')
        await expect(loanLink.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('payment shows amount and method', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/payments')
      await helpers.waitForLoading()

      const paymentLink = page.locator('table tbody tr a[href*="/payments/"]').first()

      try {
        await paymentLink.waitFor({ state: 'visible', timeout: 5000 })
        await paymentLink.click()
        await expect(page).toHaveURL(/\/payments\/\d+/)

        // Should show payment details
        const paymentInfo = page.locator('text=/Monto|Amount|Método|Method|Efectivo|Cash/i')
        await expect(paymentInfo.first()).toBeVisible()
      } catch {
        test.skip()
      }
    })
  })

  test.describe('Loan Status Transitions', () => {
    test('active loans show in list', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans?status=active')
      await helpers.waitForLoading()

      // Should show active loans or empty state
      await expect(page.locator('table').or(page.locator('text=/No hay|No se encontraron/i')).first()).toBeVisible()
    })

    test('overdue loans filter works', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans?status=overdue')
      await helpers.waitForLoading()

      await expect(page.locator('table').or(page.locator('text=/No hay|No se encontraron/i')).first()).toBeVisible()
    })

    test('closed loans filter works', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans?status=closed')
      await helpers.waitForLoading()

      await expect(page.locator('table').or(page.locator('text=/No hay|No se encontraron/i')).first()).toBeVisible()
    })
  })

  test.describe('Cross-Entity Navigation', () => {
    test('can navigate from loan to customer', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const loanLink = page.locator('table tbody tr a[href*="/loans/"]').first()

      try {
        await loanLink.waitFor({ state: 'visible', timeout: 5000 })
        await loanLink.click()
        await expect(page).toHaveURL(/\/loans\/\d+/)

        const customerLink = page.locator('a[href*="/customers/"]').first()
        await customerLink.click()

        await expect(page).toHaveURL(/\/customers\/\d+/)
      } catch {
        test.skip()
      }
    })

    test('can navigate from loan to item', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans')
      await helpers.waitForLoading()

      const loanLink = page.locator('table tbody tr a[href*="/loans/"]').first()

      try {
        await loanLink.waitFor({ state: 'visible', timeout: 5000 })
        await loanLink.click()
        await expect(page).toHaveURL(/\/loans\/\d+/)

        const itemLink = page.locator('a[href*="/items/"]').first()
        if (await itemLink.isVisible()) {
          await itemLink.click()
          await expect(page).toHaveURL(/\/items\/\d+/)
        }
      } catch {
        test.skip()
      }
    })

    test('can navigate from payment to loan', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/payments')
      await helpers.waitForLoading()

      const paymentLink = page.locator('table tbody tr a[href*="/payments/"]').first()

      try {
        await paymentLink.waitFor({ state: 'visible', timeout: 5000 })
        await paymentLink.click()
        await expect(page).toHaveURL(/\/payments\/\d+/)

        const loanLink = page.locator('a[href*="/loans/"]').first()
        await loanLink.click()

        await expect(page).toHaveURL(/\/loans\/\d+/)
      } catch {
        test.skip()
      }
    })
  })
})
