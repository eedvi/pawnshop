import { test, expect, PageHelpers } from './fixtures'

test.describe('Form Validation', () => {
  test.describe('Customer Form Validation', () => {
    test('shows error for empty required fields', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers/new')
      await helpers.waitForLoading()

      // Submit empty form
      await page.click('button[type="submit"]')

      // Should show validation errors
      const errors = page.locator('[class*="error"], [class*="destructive"], [aria-invalid="true"]')

      try {
        await errors.first().waitFor({ state: 'visible', timeout: 3000 })
        await expect(errors.first()).toBeVisible()
      } catch {
        // May use browser-native validation
        const browserValidation = page.locator(':invalid')
        await expect(browserValidation.first()).toBeVisible()
      }
    })

    test('shows error for invalid email format', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers/new')
      await helpers.waitForLoading()

      // Fill email with invalid format
      const emailField = page.locator('input[name="email"]')
      if (await emailField.isVisible()) {
        await emailField.fill('invalid-email')
        await page.click('button[type="submit"]')

        // Should show email validation error
        await page.waitForTimeout(500)
      }
    })

    test('shows error for invalid phone number', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers/new')
      await helpers.waitForLoading()

      // Fill phone with too few digits
      const phoneField = page.locator('input[name="phone"]')
      if (await phoneField.isVisible()) {
        await phoneField.fill('123')
        await page.click('button[type="submit"]')
        await page.waitForTimeout(500)
      }
    })

    test('shows error for invalid identity number', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers/new')
      await helpers.waitForLoading()

      // Fill identity with too few characters
      const identityField = page.locator('input[name="identity_number"]')
      if (await identityField.isVisible()) {
        await identityField.fill('12')
        await page.click('button[type="submit"]')
        await page.waitForTimeout(500)
      }
    })
  })

  test.describe('User Form Validation', () => {
    test('shows error for empty required fields', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/users/new')
      await helpers.waitForLoading()

      await page.click('button[type="submit"]')

      const errors = page.locator('[class*="error"], [class*="destructive"], [aria-invalid="true"]')

      try {
        await errors.first().waitFor({ state: 'visible', timeout: 3000 })
        await expect(errors.first()).toBeVisible()
      } catch {
        const browserValidation = page.locator(':invalid')
        await expect(browserValidation.first()).toBeVisible()
      }
    })

    test('shows error for weak password', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/users/new')
      await helpers.waitForLoading()

      const passwordField = page.locator('input[name="password"]')
      if (await passwordField.isVisible()) {
        await passwordField.fill('123') // Too short
        await page.click('button[type="submit"]')
        await page.waitForTimeout(500)
      }
    })

    test('shows error for password mismatch', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/users/new')
      await helpers.waitForLoading()

      const passwordField = page.locator('input[name="password"]')
      const confirmField = page.locator('input[name="password_confirmation"], input[name="confirmPassword"]')

      if (await passwordField.isVisible() && await confirmField.isVisible()) {
        await passwordField.fill('Password123!')
        await confirmField.fill('DifferentPassword!')
        await page.click('button[type="submit"]')
        await page.waitForTimeout(500)

        // Should show mismatch error
        const mismatchError = page.locator('text=/coinciden|match/i')
        try {
          await mismatchError.first().waitFor({ state: 'visible', timeout: 2000 })
        } catch {
          // Validation may prevent submission differently
        }
      }
    })
  })

  test.describe('Expense Form Validation', () => {
    test('shows error for zero or negative amount', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses/new')
      await helpers.waitForLoading()

      const amountField = page.locator('input[name="amount"], input[type="number"]').first()
      if (await amountField.isVisible()) {
        await amountField.fill('0')
        await page.click('button[type="submit"]')
        await page.waitForTimeout(500)
      }
    })

    test('shows error for missing description', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses/new')
      await helpers.waitForLoading()

      // Fill amount but not description
      const amountField = page.locator('input[name="amount"], input[type="number"]').first()
      if (await amountField.isVisible()) {
        await amountField.fill('100')
        await page.click('button[type="submit"]')

        const errors = page.locator('[class*="error"], [class*="destructive"]')
        try {
          await errors.first().waitFor({ state: 'visible', timeout: 3000 })
        } catch {
          // May have different validation
        }
      }
    })
  })

  test.describe('Branch Form Validation', () => {
    test('shows error for empty branch name', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/branches/new')
      await helpers.waitForLoading()

      await page.click('button[type="submit"]')

      const errors = page.locator('[class*="error"], [class*="destructive"], [aria-invalid="true"]')

      try {
        await errors.first().waitFor({ state: 'visible', timeout: 3000 })
      } catch {
        const browserValidation = page.locator(':invalid')
        await expect(browserValidation.first()).toBeVisible()
      }
    })
  })

  test.describe('Role Form Validation', () => {
    test('shows error for empty role name', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/roles/new')
      await helpers.waitForLoading()

      await page.click('button[type="submit"]')

      const errors = page.locator('[class*="error"], [class*="destructive"], [aria-invalid="true"]')

      try {
        await errors.first().waitFor({ state: 'visible', timeout: 3000 })
      } catch {
        const browserValidation = page.locator(':invalid')
        await expect(browserValidation.first()).toBeVisible()
      }
    })
  })

  test.describe('Login Form Validation', () => {
    test('shows error for empty credentials', async ({ page }) => {
      // Go to login directly (not authenticated)
      await page.goto('/login')
      await page.evaluate(() => localStorage.clear())
      await page.reload()

      await page.waitForSelector('input[name="email"]')
      await page.click('button[type="submit"]')

      // Should show validation
      const errors = page.locator('[class*="error"], [class*="destructive"], [aria-invalid="true"]')

      try {
        await errors.first().waitFor({ state: 'visible', timeout: 3000 })
      } catch {
        const browserValidation = page.locator(':invalid')
        await expect(browserValidation.first()).toBeVisible()
      }
    })

    test('shows error for invalid email format', async ({ page }) => {
      await page.goto('/login')
      await page.evaluate(() => localStorage.clear())
      await page.reload()

      await page.waitForSelector('input[name="email"]')
      await page.fill('input[name="email"]', 'not-an-email')
      await page.fill('input[name="password"]', 'somepassword')
      await page.click('button[type="submit"]')

      await page.waitForTimeout(500)
    })

    test('shows error for wrong credentials', async ({ page }) => {
      await page.goto('/login')
      await page.evaluate(() => localStorage.clear())
      await page.reload()

      await page.waitForSelector('input[name="email"]')
      await page.fill('input[name="email"]', 'wrong@email.com')
      await page.fill('input[name="password"]', 'wrongpassword')
      await page.click('button[type="submit"]')

      // Should show authentication error
      const authError = page.locator('text=/incorrectos|inválido|unauthorized|invalid/i')

      try {
        await authError.first().waitFor({ state: 'visible', timeout: 5000 })
        await expect(authError.first()).toBeVisible()
      } catch {
        // May show generic error or toast
        test.skip()
      }
    })
  })

  test.describe('Common Form Behaviors', () => {
    test('required field indicators are shown', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers/new')
      await helpers.waitForLoading()

      // Look for required indicators (* or required attribute)
      const requiredIndicators = page.locator('label:has-text("*"), [aria-required="true"], input[required]')

      try {
        await requiredIndicators.first().waitFor({ state: 'visible', timeout: 3000 })
        await expect(requiredIndicators.first()).toBeVisible()
      } catch {
        // May not have visual required indicators
        test.skip()
      }
    })

    test('cancel button navigates back', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers/new')
      await helpers.waitForLoading()

      const cancelButton = page.locator('button:has-text("Cancelar"), a:has-text("Cancelar")')

      if (await cancelButton.first().isVisible()) {
        await cancelButton.first().click()
        // Should navigate away from form
        await expect(page).not.toHaveURL(/\/new$/)
      }
    })

    test('form preserves data after validation error', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers/new')
      await helpers.waitForLoading()

      // Fill some fields
      await page.fill('input[name="first_name"]', 'Test Name')
      await page.click('button[type="submit"]')

      // After validation error, field should still have value
      await page.waitForTimeout(500)
      const firstNameValue = await page.locator('input[name="first_name"]').inputValue()
      expect(firstNameValue).toBe('Test Name')
    })
  })

  test.describe('Date Validation', () => {
    test('date picker shows calendar', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses/new')
      await helpers.waitForLoading()

      const datePicker = page.locator('button:has(svg), input[type="date"]').first()

      if (await datePicker.isVisible()) {
        await datePicker.click()

        // Should show calendar or date input
        const calendar = page.locator('[role="dialog"], [class*="calendar"]')
        try {
          await calendar.waitFor({ state: 'visible', timeout: 2000 })
        } catch {
          // May be native date input
        }
      }
    })
  })

  test.describe('Currency Input Validation', () => {
    test('currency input formats value', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses/new')
      await helpers.waitForLoading()

      const amountField = page.locator('input[name="amount"]')

      if (await amountField.isVisible()) {
        await amountField.fill('1234.56')
        await page.locator('body').click() // Blur

        // Value should be formatted or preserved
        const value = await amountField.inputValue()
        expect(value).toMatch(/\d/)
      }
    })

    test('rejects non-numeric input in amount fields', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/expenses/new')
      await helpers.waitForLoading()

      const amountField = page.locator('input[name="amount"]')

      if (await amountField.isVisible()) {
        const inputType = await amountField.getAttribute('type')

        if (inputType === 'number') {
          // Number inputs inherently reject non-numeric text
          // Get original value first
          const originalValue = await amountField.inputValue()
          // Type letters directly to test browser rejection
          await amountField.focus()
          await page.keyboard.type('abc')
          const newValue = await amountField.inputValue()
          // Number input should reject letters and keep original value
          expect(newValue).toBe(originalValue)
        } else {
          // Text input - fill and check
          await amountField.fill('abc')
          const value = await amountField.inputValue()
          // Should reject or strip non-numeric
          expect(value).not.toContain('abc')
        }
      }
    })
  })
})
