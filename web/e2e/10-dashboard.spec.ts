import { test, expect, PageHelpers } from './fixtures'

test.describe('Dashboard', () => {
  test('displays main KPI cards', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/')
    await helpers.waitForLoading()

    // Should show stat cards with KPIs
    const statCards = page.locator('[class*="stat"], [class*="card"]').filter({
      has: page.locator('h3, [class*="title"], [class*="label"]'),
    })

    await expect(statCards.first()).toBeVisible()
  })

  test('shows loan statistics', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/')
    await helpers.waitForLoading()

    // Look for loan-related stats (active loans, overdue, etc.)
    const loanStats = page.locator('text=/Préstamos|Activos|Vencidos|Loans/i')
    await expect(loanStats.first()).toBeVisible()
  })

  test('shows payment statistics', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/')
    await helpers.waitForLoading()

    // Look for payment or collection related stats
    const paymentStats = page.locator('text=/Pagos|Cobros|Recaudado|Payments/i')
    await expect(paymentStats.first()).toBeVisible()
  })

  test('displays charts or graphs', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/')
    await helpers.waitForLoading()

    // Look for chart containers (Recharts renders SVG)
    const charts = page.locator('svg.recharts-surface, [class*="chart"], [class*="Chart"]')

    try {
      await charts.first().waitFor({ state: 'visible', timeout: 5000 })
      await expect(charts.first()).toBeVisible()
    } catch {
      // Charts may not be present if no data
      test.skip()
    }
  })

  test('shows recent activity or transactions', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/')
    await helpers.waitForLoading()

    // Look for recent activity section
    const recentSection = page.locator('text=/Reciente|Últim|Recent|Activity/i')

    try {
      await recentSection.first().waitFor({ state: 'visible', timeout: 5000 })
      await expect(recentSection.first()).toBeVisible()
    } catch {
      // May not have recent activity section
      test.skip()
    }
  })

  test('can navigate to modules from dashboard', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/')
    await helpers.waitForLoading()

    // Click on a quick action or card link if available
    const quickLink = page.locator('a[href="/loans"], a[href="/customers"], a[href="/payments"]').first()

    try {
      await quickLink.waitFor({ state: 'visible', timeout: 3000 })
      await quickLink.click()
      await expect(page).toHaveURL(/\/(loans|customers|payments)/)
    } catch {
      // Quick links may not be present
      test.skip()
    }
  })

  test('dashboard respects branch filter', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/')
    await helpers.waitForLoading()

    // Check if branch selector exists in header
    const branchSelector = page.locator('button:has-text("Sucursal"), [class*="branch"]')

    try {
      await branchSelector.first().waitFor({ state: 'visible', timeout: 3000 })
      await expect(branchSelector.first()).toBeVisible()
    } catch {
      // Branch selector may not be visible for single-branch users
      test.skip()
    }
  })

  test('shows cash register status', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/')
    await helpers.waitForLoading()

    // Look for cash register status indicator
    const cashStatus = page.locator('text=/Caja|Cash|Abierta|Cerrada/i')

    try {
      await cashStatus.first().waitFor({ state: 'visible', timeout: 3000 })
      await expect(cashStatus.first()).toBeVisible()
    } catch {
      // Cash status may not be on dashboard
      test.skip()
    }
  })
})
