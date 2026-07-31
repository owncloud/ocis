import path from 'node:path'
import { defineConfig, devices, ReporterDescription } from '@playwright/test'
import { config } from './config'

const __dirname = path.dirname(new URL(import.meta.url).pathname)
const reportsDir = path.resolve(__dirname, '../../', config.reportDir)

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  // Set default test timeout to 60 seconds
  timeout: config.timeout * 1000,

  // Look for test files in the following directory, relative to this configuration file.
  testDir: 'specs',

  // Run all tests in parallel.
  fullyParallel: true,

  // Fail the build on CI if you accidentally left test.only in the source code.
  forbidOnly: !!process.env.CI,

  // Retry on CI only.
  retries: config.retry,

  // Run tests in parallel - use CI-determined workers or auto-detect locally
  //   workers: process.env.CI ? 2 : undefined,
  workers: undefined,

  // Reporter to use
  reporter: [
    ['list'] as ReporterDescription,
    ['./reporters/a11y.ts', { outputFile: path.join(reportsDir, 'a11y-report.json') }]
  ].filter(Boolean) as ReporterDescription[],
  outputDir: reportsDir,

  use: {
    ignoreHTTPSErrors: true,

    // Collect trace when retrying the failed test.
    trace: config.reportTracing ? 'on' : 'on-first-retry',
    headless: config.headless
  },
  // Configure projects for major browsers.
  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        contextOptions: {
          permissions: ['clipboard-read', 'clipboard-write', 'camera', 'microphone']
        }
      }
    },
    {
      name: 'firefox',
      use: {
        ...devices['Desktop Firefox'],
        launchOptions: {
          firefoxUserPrefs: {
            'dom.events.testing.asyncClipboard': true,
            'dom.events.asyncClipboard.readText': true,
            'dom.events.asyncClipboard.clipboardItem': true,
            'dom.events.asyncClipboard.writeText': true,
            'permissions.default.clipboard-read': 1,
            'permissions.default.clipboard-write': 1
          }
        }
      }
    },
    {
      name: 'webkit',
      use: {
        ...devices['Desktop Safari'],
        contextOptions: {
          permissions: ['clipboard-read']
        }
      }
    }
  ]
})
