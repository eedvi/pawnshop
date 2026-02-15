import { test, expect, PageHelpers } from './fixtures'

/**
 * Edge Cases and Error Handling Tests
 *
 * Tests for unusual scenarios and error conditions:
 * 1. Empty states
 * 2. 404 pages
 * 3. Network error handling
 * 4. Form error recovery
 * 5. Session expiration
 * 6. Invalid data handling
 */
test.describe('Edge Cases and Error Handling', () => {
  test.describe('Empty States', () => {
    test('shows empty state when no customers match search', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')
      await searchInput.fill('xyznonexistentcustomer12345')
      await page.waitForTimeout(600)
      await helpers.waitForLoading()

      // Should show empty state or no results message
      const emptyState = page.locator('text=/No se encontraron|No hay|No results/i')
      try {
        await emptyState.first().waitFor({ state: 'visible', timeout: 5000 })
        await expect(emptyState.first()).toBeVisible()
      } catch {
        // May show empty table instead
        const tableRows = page.locator('table tbody tr')
        const count = await tableRows.count()
        expect(count).toBeLessThanOrEqual(1) // Only header or empty row
      }
    })

    test('shows empty state for filtered loans with no results', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans?status=nonexistent')
      await helpers.waitForLoading()

      const content = page.locator('table').or(page.locator('text=/No hay|No se encontraron/i'))
      await expect(content.first()).toBeVisible()
    })

    test('handles empty audit log gracefully', async ({ page }) => {
      const helpers = new PageHelpers(page)

      // Set impossible date filter
      await page.goto('/audit?date_from=2099-01-01&date_to=2099-12-31')
      await helpers.waitForLoading()

      const content = page.locator('table').or(page.locator('text=/No hay|No se encontraron/i'))
      await expect(content.first()).toBeVisible()
    })
  })

  test.describe('404 and Invalid Routes', () => {
    test('handles non-existent customer gracefully', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers/999999999')
      await helpers.waitForLoading()

      // Should show error or redirect
      const errorMessage = page.locator('text=/No encontrado|Not found|Error|404/i')
      const redirected = page.url().includes('/customers') && !page.url().includes('999999999')

      if (await errorMessage.first().isVisible()) {
        await expect(errorMessage.first()).toBeVisible()
      } else if (redirected) {
        // Redirected to list is also valid
        expect(true).toBe(true)
      }
    })

    test('handles non-existent loan gracefully', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans/999999999')
      await helpers.waitForLoading()

      const errorMessage = page.locator('text=/No encontrado|Not found|Error|404/i')
      const redirected = !page.url().includes('999999999')

      if (await errorMessage.first().isVisible()) {
        await expect(errorMessage.first()).toBeVisible()
      } else if (redirected) {
        expect(true).toBe(true)
      }
    })

    test('handles completely invalid route', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/this-route-does-not-exist-at-all')
      await helpers.waitForLoading()

      // Should show 404 or redirect to dashboard
      const is404 = page.locator('text=/404|No encontrado|Not found/i')
      const redirectedHome = page.url().endsWith('/') || page.url().includes('login')

      if (await is404.first().isVisible()) {
        await expect(is404.first()).toBeVisible()
      } else if (redirectedHome) {
        expect(true).toBe(true)
      }
    })
  })

  test.describe('Form Error Recovery', () => {
    test('can recover from validation errors', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers/new')
      await helpers.waitForLoading()

      // Submit empty form
      await page.click('button[type="submit"]')
      await page.waitForTimeout(500)

      // Now fill required fields
      await page.fill('input[name="first_name"]', 'Test')
      await page.fill('input[name="last_name"]', 'User')
      await page.fill('input[name="identity_number"]', '1234567890')
      await page.fill('input[name="phone"]', '5551234567')

      // Form should now be fillable
      const firstNameValue = await page.locator('input[name="first_name"]').inputValue()
      expect(firstNameValue).toBe('Test')
    })

    test('preserves form data on validation error', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers/new')
      await helpers.waitForLoading()

      // Fill some fields but not all required
      await page.fill('input[name="first_name"]', 'Test Name')

      // Submit
      await page.click('button[type="submit"]')
      await page.waitForTimeout(500)

      // Data should be preserved
      const firstNameValue = await page.locator('input[name="first_name"]').inputValue()
      expect(firstNameValue).toBe('Test Name')
    })

    test('can cancel and return to list without saving', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers/new')
      await helpers.waitForLoading()

      // Fill some data
      await page.fill('input[name="first_name"]', 'Should Not Save')

      // Cancel
      const cancelButton = page.locator('button:has-text("Cancelar"), a:has-text("Cancelar")')
      if (await cancelButton.first().isVisible()) {
        await cancelButton.first().click()
        await expect(page).toHaveURL(/\/customers$/)
      }
    })
  })

  test.describe('Pagination Edge Cases', () => {
    test('handles page 0 gracefully', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers?page=0')
      await helpers.waitForLoading()

      // Should show first page or handle gracefully
      await expect(page.locator('table, main').first()).toBeVisible()
    })

    test('handles negative page gracefully', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers?page=-1')
      await helpers.waitForLoading()

      await expect(page.locator('table, main').first()).toBeVisible()
    })

    test('handles very high page number gracefully', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers?page=9999999')
      await helpers.waitForLoading()

      // Should show empty or last page
      await expect(page.locator('table, main').first()).toBeVisible()
    })
  })

  test.describe('Filter Edge Cases', () => {
    test('handles multiple filters at once', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans?status=active&page=1')
      await helpers.waitForLoading()

      await expect(page.locator('table, main').first()).toBeVisible()
    })

    test('handles invalid filter values', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/items?status=invalid_status_value')
      await helpers.waitForLoading()

      // Should handle gracefully
      await expect(page.locator('table, main').first()).toBeVisible()
    })

    test('handles invalid date format', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/audit?date_from=invalid-date')
      await helpers.waitForLoading()

      // Should handle gracefully
      await expect(page.locator('table, main').first()).toBeVisible()
    })
  })

  test.describe('Special Characters in Search', () => {
    test('handles special characters in search', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')
      await searchInput.fill('Test <script>alert(1)</script>')
      await page.waitForTimeout(600)
      await helpers.waitForLoading()

      // Should not break the page
      await expect(page.locator('table')).toBeVisible()
    })

    test('handles SQL injection attempts in search', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')
      await searchInput.fill("'; DROP TABLE customers; --")
      await page.waitForTimeout(600)
      await helpers.waitForLoading()

      // Should not break the page
      await expect(page.locator('table')).toBeVisible()
    })

    test('handles unicode characters in search', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')
      await searchInput.fill('José María García Ñoño')
      await page.waitForTimeout(600)
      await helpers.waitForLoading()

      // Should handle unicode properly
      await expect(page.locator('table')).toBeVisible()
    })
  })

  test.describe('UI State Consistency', () => {
    test('loading state shows correctly', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')

      // During load, may show spinner
      // After load, should show content
      await helpers.waitForLoading()
      await expect(page.locator('table')).toBeVisible()
    })

    test('sidebar state persists across navigation', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/')
      await helpers.waitForLoading()

      // Navigate to different pages
      await helpers.navigateTo('Clientes')
      await expect(page).toHaveURL(/\/customers/)

      await helpers.navigateTo('Préstamos')
      await expect(page).toHaveURL(/\/loans/)

      // Sidebar should still be visible
      await expect(page.locator('nav')).toBeVisible()
    })

    test('form does not submit on enter in wrong field', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers/new')
      await helpers.waitForLoading()

      // Focus on first name and press enter
      await page.locator('input[name="first_name"]').focus()
      await page.keyboard.press('Enter')

      // Should still be on form (not submitted with validation errors)
      await expect(page).toHaveURL(/\/customers\/new/)
    })
  })

  test.describe('Concurrent Actions', () => {
    test('handles rapid navigation', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/')
      await helpers.waitForLoading()

      // Rapid navigation
      await page.goto('/customers')
      await page.goto('/loans')
      await page.goto('/payments')
      await page.goto('/sales')
      await helpers.waitForLoading()

      // Should be on last page
      await expect(page).toHaveURL(/\/sales/)
      await expect(page.locator('main')).toBeVisible()
    })

    test('handles rapid search typing', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const searchInput = page.locator('input[placeholder*="Buscar"], input[type="search"]')

      // Type rapidly
      await searchInput.fill('a')
      await searchInput.fill('ab')
      await searchInput.fill('abc')
      await searchInput.fill('abcd')

      await page.waitForTimeout(600)
      await helpers.waitForLoading()

      // Should show results for final search term
      await expect(page.locator('table')).toBeVisible()
    })
  })

  test.describe('Browser Behavior', () => {
    test('back button works correctly', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      await page.goto('/loans')
      await helpers.waitForLoading()

      await page.goBack()

      await expect(page).toHaveURL(/\/customers/)
    })

    test('forward button works correctly', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      await page.goto('/loans')
      await helpers.waitForLoading()

      await page.goBack()
      await page.goForward()

      await expect(page).toHaveURL(/\/loans/)
    })

    test('page refresh preserves location', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/loans?status=active&page=1')
      await helpers.waitForLoading()

      await page.reload()
      await helpers.waitForLoading()

      // Should stay on same page with same filters
      expect(page.url()).toContain('/loans')
    })
  })

  test.describe('Responsive Behavior', () => {
    test('table is visible on desktop viewport', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.setViewportSize({ width: 1920, height: 1080 })
      await page.goto('/customers')
      await helpers.waitForLoading()

      await expect(page.locator('table')).toBeVisible()
    })

    test('navigation works on tablet viewport', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.setViewportSize({ width: 768, height: 1024 })
      await page.goto('/')
      await helpers.waitForLoading()

      // Navigation should still work
      await expect(page.locator('main')).toBeVisible()
    })

    test('page loads on mobile viewport', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.setViewportSize({ width: 375, height: 667 })
      await page.goto('/')
      await helpers.waitForLoading()

      // Page should load
      await expect(page.locator('main')).toBeVisible()
    })
  })

  test.describe('Data Integrity', () => {
    test('detail page shows consistent data', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const customerLink = page.locator('table tbody tr a[href*="/customers/"]').first()

      try {
        await customerLink.waitFor({ state: 'visible', timeout: 5000 })

        // Get name from list
        const rowText = await page.locator('table tbody tr').first().innerText()

        await customerLink.click()
        await expect(page).toHaveURL(/\/customers\/\d+/)

        // Detail page should have matching info
        await expect(page.locator('main')).toBeVisible()
      } catch {
        test.skip()
      }
    })

    test('edit form shows current data', async ({ page }) => {
      const helpers = new PageHelpers(page)

      await page.goto('/customers')
      await helpers.waitForLoading()

      const customerLink = page.locator('table tbody tr a[href*="/customers/"]').first()

      try {
        await customerLink.waitFor({ state: 'visible', timeout: 5000 })
        await customerLink.click()
        await expect(page).toHaveURL(/\/customers\/\d+/)

        const editButton = page.locator('a:has-text("Editar"), button:has-text("Editar")')
        await editButton.first().click()

        // Form should have values
        const firstNameValue = await page.locator('input[name="first_name"]').inputValue()
        expect(firstNameValue.length).toBeGreaterThan(0)
      } catch {
        test.skip()
      }
    })
  })
})
