import { test as base, expect, Page } from '@playwright/test'

/**
 * Test Credentials - These should match seeded data
 */
export const TEST_USERS = {
  admin: {
    email: 'admin@pawnshop.com',
    password: 'admin123',
  },
  manager: {
    email: 'gerente@pawnshop.com',
    password: 'admin123',
  },
  cashier: {
    email: 'cajero@pawnshop.com',
    password: 'admin123',
  },
}

/**
 * Helper to login and store auth state
 */
export async function login(page: Page, email: string, password: string) {
  await page.goto('/login')
  await page.waitForSelector('input[name="email"]')
  await page.fill('input[name="email"]', email)
  await page.fill('input[name="password"]', password)
  await page.click('button[type="submit"]')
  // Wait for redirect to dashboard or any authenticated page
  await page.waitForURL(/\/(dashboard|customers|items|loans)?$/, { timeout: 15000 })
}

/**
 * Extended test fixture with authenticated context
 */
type AuthFixtures = {
  authenticatedPage: Page
}

export const test = base.extend<AuthFixtures>({
  authenticatedPage: async ({ page }, use) => {
    await login(page, TEST_USERS.admin.email, TEST_USERS.admin.password)
    await use(page)
  },
})

export { expect }

/**
 * Common page object helpers
 */
export class PageHelpers {
  constructor(private page: Page) {}

  /**
   * Wait for table to load data (actual data rows, not empty state)
   */
  async waitForTableData() {
    // Wait for table to exist
    await this.page.waitForSelector('table tbody tr', { timeout: 10000 })
    // Check if it's actually data (has links) not empty state
    const hasLinks = await this.page.locator('table tbody tr').first().locator('a[href*="/"]').count()
    if (hasLinks === 0) {
      // Table is showing empty state - wait a bit more for data to load
      await this.page.waitForSelector('table tbody tr a[href*="/"]', { timeout: 5000 }).catch(() => {})
    }
  }

  /**
   * Check if table has actual data rows (not just empty state)
   */
  async tableHasData(): Promise<boolean> {
    const links = await this.page.locator('table tbody tr a[href*="/"]').count()
    return links > 0
  }

  /**
   * Wait for loading to complete
   */
  async waitForLoading() {
    // Wait for any loading spinners to disappear
    await this.page.waitForSelector('[class*="animate-spin"]', { state: 'detached', timeout: 5000 }).catch(() => {})
  }

  /**
   * Click on the first link in a table row to navigate to detail
   * Tables use links in cells (not row click) for navigation
   */
  async clickFirstRowLink() {
    // First ensure there's actually data in the table
    const hasData = await this.tableHasData()
    if (!hasData) {
      throw new Error('Cannot click first row link: table has no data rows (empty state)')
    }
    // Find the first link in the table body that goes to a detail page
    const link = this.page.locator('table tbody tr').first().locator('a[href*="/"]').first()
    await link.click()
  }

  /**
   * Click on a specific row's link by row index
   */
  async clickRowLink(rowIndex: number) {
    const link = this.page.locator('table tbody tr').nth(rowIndex).locator('a[href*="/"]').first()
    await link.click()
  }

  /**
   * Open row action menu by row index
   */
  async openRowActions(rowIndex: number) {
    const actionButton = this.page.locator('table tbody tr').nth(rowIndex).locator('button[aria-haspopup="menu"]')
    await actionButton.click()
  }

  /**
   * Click row action menu item
   */
  async clickRowAction(rowIndex: number, actionText: string) {
    await this.openRowActions(rowIndex)
    await this.page.click(`[role="menuitem"]:has-text("${actionText}")`)
  }

  /**
   * Fill form field by label
   */
  async fillField(label: string, value: string) {
    const field = this.page.locator(`label:has-text("${label}")`).locator('..').locator('input, textarea, select')
    await field.fill(value)
  }

  /**
   * Select option from dropdown by label
   */
  async selectOption(label: string, optionText: string) {
    const trigger = this.page.locator(`label:has-text("${label}")`).locator('..').locator('button[role="combobox"]')
    await trigger.click()
    await this.page.locator(`[role="option"]:has-text("${optionText}")`).click()
  }

  /**
   * Click primary submit button
   */
  async submitForm() {
    await this.page.click('button[type="submit"]')
  }

  /**
   * Verify toast notification appears
   */
  async expectToast(text: string) {
    await expect(this.page.locator('[data-sonner-toast]').filter({ hasText: text })).toBeVisible({ timeout: 5000 })
  }

  /**
   * Navigate via sidebar
   */
  async navigateTo(menuText: string) {
    await this.page.click(`nav a:has-text("${menuText}")`)
    await this.waitForLoading()
  }

  /**
   * Confirm dialog action
   */
  async confirmDialog() {
    await this.page.click('button:has-text("Confirmar"), button:has-text("Sí"), button:has-text("Continuar")')
  }

  /**
   * Cancel dialog
   */
  async cancelDialog() {
    await this.page.click('button:has-text("Cancelar"), button:has-text("No")')
  }
}

/**
 * Generate unique test data
 */
export function generateTestData() {
  const timestamp = Date.now()
  return {
    customer: {
      firstName: `Test${timestamp}`,
      lastName: 'Cliente',
      identityNumber: `${1000000000 + timestamp % 900000000}`,
      phone: `5${String(timestamp).slice(-7)}`,
      email: `test${timestamp}@example.com`,
    },
    item: {
      description: `Artículo de prueba ${timestamp}`,
      serialNumber: `SN-${timestamp}`,
      appraisedValue: '1000',
    },
  }
}
