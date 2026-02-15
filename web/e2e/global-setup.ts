import { chromium, FullConfig } from '@playwright/test'
import { TEST_USERS } from './fixtures'

/**
 * Global setup - runs once before all tests
 * Authenticates and saves storage state for reuse
 */
async function globalSetup(config: FullConfig) {
  const { baseURL } = config.projects[0].use

  const browser = await chromium.launch()
  const page = await browser.newPage()

  try {
    // Navigate to login
    await page.goto(`${baseURL}/login`)
    await page.waitForSelector('input[name="email"]', { timeout: 30000 })

    // Login as admin
    await page.fill('input[name="email"]', TEST_USERS.admin.email)
    await page.fill('input[name="password"]', TEST_USERS.admin.password)
    await page.click('button[type="submit"]')

    // Wait for successful login (redirect away from login page)
    await page.waitForURL(/\/(dashboard|customers|items|loans)?$/, { timeout: 30000 })

    // Save storage state (localStorage with tokens, cookies)
    await page.context().storageState({ path: './e2e/.auth/admin.json' })

    console.log('Global setup: Admin authentication saved')
  } catch (error) {
    console.error('Global setup failed:', error)
    throw error
  } finally {
    await browser.close()
  }
}

export default globalSetup
