import { Page, expect } from '@playwright/test'
import { config } from '../../config'

/**
 * Log in to ownCloud Web through the built-in IdP and wait until the Files
 * view is ready (the account menu is present in the top bar).
 */
export async function login(page: Page): Promise<void> {
  await page.goto('/')
  const username = page.locator('input[placeholder="Username"]')
  await username.waitFor({ state: 'visible', timeout: 30_000 })
  await username.fill(config.adminUsername)
  await page.locator('input[placeholder="Password"]').fill(config.adminPassword)
  await page.locator('input[placeholder="Password"]').press('Enter')
  await page
    .getByRole('button', { name: 'My Account' })
    .waitFor({ state: 'visible', timeout: 30_000 })
  await page.waitForTimeout(600)
}

/** Dismiss any open popover/menu so the next step starts from a clean top bar. */
export async function dismissOverlays(page: Page): Promise<void> {
  await page.keyboard.press('Escape').catch(() => {})
  await page.waitForTimeout(200)
}

/**
 * Open Personal, select the seeded `report.md`, and reveal the right sidebar
 * (it is collapsed in a fresh session). Shared by the tours that start from a
 * selected file.
 */
export async function selectReportAndOpenSidebar(page: Page): Promise<void> {
  await page.goto('/files/spaces/personal')
  await page.waitForURL(/files\/spaces\/personal/, { timeout: 30_000 })
  const checkbox = page.getByRole('row', { name: /report\.md/ }).getByLabel('Select file')
  await checkbox.waitFor({ state: 'visible', timeout: 30_000 })
  await checkbox.click()
  const openSidebar = page.getByRole('button', { name: 'Open sidebar to view details' })
  if (await openSidebar.isVisible().catch(() => false)) {
    await openSidebar.click()
  }
}

/** Open the right sidebar's Shares panel and wait for it to render. */
export async function openSharingPanel(page: Page): Promise<void> {
  await page.locator('[data-testid="sidebar-panel-sharing-select"]').click()
  await expect(page.getByRole('heading', { name: 'Share with people' })).toBeVisible({
    timeout: 15_000
  })
}

/** Click a left-sidebar navigation entry and wait for its route to load. */
export async function openSection(page: Page, name: string, url: RegExp): Promise<void> {
  await page.getByRole('link', { name }).click()
  await page.waitForURL(url, { timeout: 30_000 })
  await page.waitForTimeout(900)
}
