import { Page } from '@playwright/test'
import { VaultPage } from './page/VaultPage'
import { config } from '../../../config'
export interface userEntersVaultModeArgs {
  page: Page
}

export interface captureQrCodeScreenshotArgs {
  page: Page
}

export interface userAuthenticatesWithOTPArgs {
  page: Page
  otp: string
}

export interface waitForVaultModeArgs {
  page: Page
}

export interface userEntersDriveModeArgs {
  page: Page
}

export const userEntersVaultMode = async ({ page }: userEntersVaultModeArgs): Promise<void> => {
  const vaultPage = new VaultPage({ page })

  await vaultPage.driveOption.click()
  await vaultPage.vaultOption.click()
  await vaultPage.qrImage.waitFor({ state: 'visible' })
}

export const captureQrCodeScreenshot = async ({
  page
}: captureQrCodeScreenshotArgs): Promise<Buffer> => {
  const vaultPage = new VaultPage({ page })

  return await vaultPage.qrImage.screenshot()
}

export const userAuthenticatesWithOTP = async ({
  page,
  otp
}: userAuthenticatesWithOTPArgs): Promise<void> => {
  const vaultPage = new VaultPage({ page })

  await vaultPage.oneTimeCodeTextbox.fill(otp)
  await vaultPage.otpSubmitButton.click()

  await waitForVaultMode({ page })
}

export const waitForVaultMode = async ({ page }: waitForVaultModeArgs): Promise<void> => {
  const vaultPage = new VaultPage({ page })
  const vaultUrl = `${config.baseUrl}/vault`

  await page.waitForURL((url) => url.href.startsWith(vaultUrl))
  await vaultPage.vaultBreadcrumb.waitFor()
}

export const userEntersDriveMode = async ({ page }: userEntersDriveModeArgs): Promise<void> => {
  const vaultPage = new VaultPage({ page })

  await vaultPage.vaultOption.click()
  await vaultPage.driveOption.click()
  await vaultPage.driveBreadcrumb.waitFor({ state: 'visible' })
}
