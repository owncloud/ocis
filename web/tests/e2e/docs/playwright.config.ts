import { defineConfig } from '@playwright/test'
import { config } from '../config'

// Standalone Playwright config for the documentation screenshot capture tool.
// It reuses the repository's Playwright install and the shared e2e config
// (tests/e2e/config.js): point BASE_URL_OCIS at a running ownCloud Web / oCIS
// instance and set ADMIN_USERNAME / ADMIN_PASSWORD (all default to the same
// values as the rest of the e2e suite).
export default defineConfig({
  testDir: '.',
  testMatch: 'capture.spec.ts',
  outputDir: './test-results',
  fullyParallel: false,
  workers: 1,
  retries: process.env.CI ? 1 : 0,
  reporter: 'list',
  timeout: 120_000,
  use: {
    baseURL: config.baseUrl,
    ignoreHTTPSErrors: true,
    headless: true,
    viewport: { width: 1440, height: 900 }
  }
})
