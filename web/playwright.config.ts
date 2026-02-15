import { defineConfig, devices } from '@playwright/test'

/**
 * PawnShop Admin Panel - E2E Test Configuration
 *
 * Tests simulate real user workflows through the admin panel.
 * Requires backend running on localhost:8090 and database seeded.
 */
export default defineConfig({
  testDir: './e2e',
  fullyParallel: false, // Run tests sequentially for realistic flow
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1, // Single worker for consistent state
  reporter: [
    ['html', { open: 'never' }],
    ['list'],
  ],
  timeout: 60000, // 60 second timeout per test
  expect: {
    timeout: 10000, // 10 second timeout for assertions
  },

  // Global setup - authenticate once before all tests
  globalSetup: './e2e/global-setup.ts',

  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'on-first-retry',
    actionTimeout: 15000,
  },

  projects: [
    // Auth tests run without stored state
    {
      name: 'auth-tests',
      testMatch: /01-auth\.spec\.ts/,
      use: { ...devices['Desktop Chrome'] },
    },
    // All other tests use pre-authenticated state
    {
      name: 'authenticated-tests',
      testIgnore: /01-auth\.spec\.ts/,
      use: {
        ...devices['Desktop Chrome'],
        storageState: './e2e/.auth/admin.json',
      },
    },
  ],

  // Start local dev server before running tests
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
    timeout: 120000, // 2 minutes to start dev server
  },
})
