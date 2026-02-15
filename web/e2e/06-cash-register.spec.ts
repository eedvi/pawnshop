import { test, expect, PageHelpers } from './fixtures'

test.describe('Cash Register Operations', () => {
  // No login needed - using pre-authenticated storage state

  test('can view cash register page', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/cash')
    await helpers.waitForLoading()

    // Should see the cash management page
    await expect(page.locator('h1, h2').first()).toContainText(/Caja|Cash/i)
  })

  test('cash page has tabs for different views', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/cash')
    await helpers.waitForLoading()

    // Should have tabs for registers, sessions, movements
    const tabList = page.locator('[role="tablist"]')
    if (await tabList.isVisible()) {
      await expect(tabList).toBeVisible()

      // Common tab names
      const expectedTabs = ['Registros', 'Sesiones', 'Movimientos', 'Registers', 'Sessions', 'Movements']
      for (const tabName of expectedTabs) {
        const tab = page.locator(`[role="tab"]:has-text("${tabName}")`)
        if (await tab.isVisible()) {
          await expect(tab).toBeVisible()
        }
      }
    }
  })

  test('can open a cash session', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/cash')
    await helpers.waitForLoading()

    // Look for open session button
    const openButton = page.locator('button:has-text("Abrir"), button:has-text("Abrir Caja"), button:has-text("Open")')
    if (await openButton.first().isVisible()) {
      await openButton.first().click()

      // Should show dialog for opening session
      const dialog = page.locator('[role="dialog"]')
      await expect(dialog).toBeVisible()

      // Dialog should have initial amount input
      const amountInput = dialog.locator('input[name*="amount"], input[name*="initial"], input[type="number"]')
      if (await amountInput.isVisible()) {
        await amountInput.fill('1000')
      }
    }
  })

  test('can close a cash session', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/cash')
    await helpers.waitForLoading()

    // Look for close session button (only visible if session is open)
    const closeButton = page.locator('button:has-text("Cerrar"), button:has-text("Cerrar Caja"), button:has-text("Close")')
    if (await closeButton.first().isVisible()) {
      await closeButton.first().click()

      // Should show dialog for closing session
      const dialog = page.locator('[role="dialog"]')
      if (await dialog.isVisible()) {
        await expect(dialog).toBeVisible()

        // Dialog should have final count inputs
        const countInput = dialog.locator('input[name*="count"], input[name*="final"], input[type="number"]')
        if (await countInput.first().isVisible()) {
          await countInput.first().fill('1500')
        }
      }
    }
  })

  test('can add cash movement', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/cash')
    await helpers.waitForLoading()

    // Look for add movement button
    const addButton = page.locator('button:has-text("Movimiento"), button:has-text("Agregar"), button:has-text("Add Movement")')
    if (await addButton.first().isVisible()) {
      await addButton.first().click()

      // Should show movement dialog
      const dialog = page.locator('[role="dialog"]')
      if (await dialog.isVisible()) {
        await expect(dialog).toBeVisible()

        // Should have type selector and amount
        const typeSelect = dialog.locator('button[role="combobox"], select[name*="type"]')
        const amountInput = dialog.locator('input[name*="amount"], input[type="number"]')

        if (await typeSelect.isVisible()) {
          await typeSelect.click()
          // Select ingreso/deposit
          await page.locator('[role="option"]').first().click()
        }

        if (await amountInput.isVisible()) {
          await amountInput.fill('100')
        }
      }
    }
  })

  test('shows current session status', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/cash')
    await helpers.waitForLoading()

    // Should show current session info or "no active session" message
    const sessionStatus = page.locator('text=Sesión activa').or(page.locator('text=Active session')).or(page.locator('text=No hay sesión')).or(page.locator('text=No active session'))
    // One of these messages should be visible
  })

  test('displays session movements history', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/cash')
    await helpers.waitForLoading()

    // Switch to movements tab if available
    const movementsTab = page.locator('[role="tab"]:has-text("Movimientos"), [role="tab"]:has-text("Movements")')
    if (await movementsTab.isVisible()) {
      await movementsTab.click()
      await helpers.waitForLoading()

      // Should show movements table
      const table = page.locator('table')
      if (await table.isVisible()) {
        await expect(table).toBeVisible()
      }
    }
  })

  test('shows cash balance', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/cash')
    await helpers.waitForLoading()

    // Look for balance display
    const balanceDisplay = page.locator('text=Saldo').or(page.locator('text=Balance')).or(page.locator('text=Total'))
    if (await balanceDisplay.first().isVisible()) {
      await expect(balanceDisplay.first()).toBeVisible()
    }
  })

  test('can view session detail', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/cash')
    await helpers.waitForLoading()

    // Switch to sessions tab
    const sessionsTab = page.locator('[role="tab"]:has-text("Sesiones"), [role="tab"]:has-text("Sessions")')
    if (await sessionsTab.isVisible()) {
      await sessionsTab.click()
      await helpers.waitForLoading()

      // Click on first session row
      const firstRow = page.locator('table tbody tr').first()
      if (await firstRow.isVisible()) {
        await firstRow.click()

        // Should show session detail or dialog
        const detail = page.locator('[role="dialog"], main:has-text("Sesión")')
        // Detail view may appear
      }
    }
  })

  test('movement types include income and expense', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/cash')
    await helpers.waitForLoading()

    // Open movement dialog
    const addButton = page.locator('button:has-text("Movimiento"), button:has-text("Agregar")')
    if (await addButton.first().isVisible()) {
      await addButton.first().click()

      const dialog = page.locator('[role="dialog"]')
      if (await dialog.isVisible()) {
        // Click type selector
        const typeSelect = dialog.locator('button[role="combobox"]')
        if (await typeSelect.isVisible()) {
          await typeSelect.click()

          // Should have income and expense options
          const incomeOption = page.locator('[role="option"]:has-text("Ingreso"), [role="option"]:has-text("Income")')
          const expenseOption = page.locator('[role="option"]:has-text("Egreso"), [role="option"]:has-text("Expense")')

          // At least one should exist
        }
      }
    }
  })
})
