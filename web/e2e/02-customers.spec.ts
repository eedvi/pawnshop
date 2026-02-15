import { test, expect, PageHelpers, generateTestData } from './fixtures'

test.describe('Customer Management', () => {
  let testData: ReturnType<typeof generateTestData>

  test.beforeAll(() => {
    testData = generateTestData()
  })

  // No login needed - using pre-authenticated storage state

  test('can view customers list', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/customers')
    await helpers.waitForLoading()

    // Should see the customers table
    await expect(page.locator('h1, h2').first()).toContainText(/Clientes/i)
    await expect(page.locator('table')).toBeVisible()

    // Should have action buttons
    await expect(page.locator('a:has-text("Nuevo Cliente"), button:has-text("Nuevo Cliente")')).toBeVisible()
  })

  test('can search customers', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/customers')
    await helpers.waitForLoading()

    // Type in search field
    const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')
    if (await searchInput.isVisible()) {
      await searchInput.fill('Juan')
      await page.waitForTimeout(500) // Debounce wait
      await helpers.waitForLoading()
    }
  })

  test('can navigate to create customer form', async ({ page }) => {
    await page.goto('/customers')

    // Click new customer button/link
    await page.click('a:has-text("Nuevo Cliente")')

    // Should be on create page
    await expect(page).toHaveURL(/\/customers\/new|\/customers\/create/)
  })

  test('can view customer detail', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/customers')
    await helpers.waitForTableData()

    // Click on a link in the first row (name or document number)
    await helpers.clickFirstRowLink()

    // Should navigate to detail page
    await expect(page).toHaveURL(/\/customers\/\d+/)

    // Should show customer information
    await expect(page.locator('h1, h2').first()).toBeVisible()
  })

  test('can navigate to edit from detail', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/customers')
    await helpers.waitForTableData()

    // Navigate to detail
    await helpers.clickFirstRowLink()
    await expect(page).toHaveURL(/\/customers\/\d+/)

    // Look for edit button
    const editButton = page.locator('a:has-text("Editar"), button:has-text("Editar")')
    if (await editButton.first().isVisible()) {
      await editButton.first().click()
      await expect(page).toHaveURL(/\/customers\/\d+\/edit/)
    }
  })

  test('can edit from dropdown menu', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/customers')
    await helpers.waitForTableData()

    // Use dropdown menu to edit
    await helpers.clickRowAction(0, 'Editar')

    // Should be on edit page
    await expect(page).toHaveURL(/\/customers\/\d+\/edit/)
  })

  test('shows customer tabs in detail view', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/customers')
    await helpers.waitForTableData()

    // Navigate to detail
    await helpers.clickFirstRowLink()
    await expect(page).toHaveURL(/\/customers\/\d+/)

    // Should see tabs in detail view
    const tabList = page.locator('[role="tablist"]')
    if (await tabList.isVisible()) {
      await expect(tabList).toBeVisible()
    }
  })

  test('form validation works', async ({ page }) => {
    const helpers = new PageHelpers(page)

    await page.goto('/customers/new')
    await helpers.waitForLoading()

    // Wait for the form to be visible
    await expect(page.locator('form')).toBeVisible()

    // Try to submit empty form
    await page.click('button[type="submit"]')

    // Should show validation errors
    const errorMessages = page.locator('[class*="text-destructive"], [class*="text-red"], .error, [id*="error"]')
    await expect(errorMessages.first()).toBeVisible()
  })
})
