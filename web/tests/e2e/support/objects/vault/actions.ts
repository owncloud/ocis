import { Page } from '@playwright/test'
import { VaultPage } from './page/VaultPage'
import { config } from '../../../config'

export interface userAuthenticatesWithOTPArgs {
  page: Page
  otp: string
}

export const userEntersVaultMode = async ({ page }: { page: Page }): Promise<void> => {
  const vaultPage = new VaultPage({ page })

  await vaultPage.driveOption.click()
  await vaultPage.vaultOption.click()
  await vaultPage.qrImage.waitFor({ state: 'visible' })
}

export const waitForVaultMode = async ({ page }: { page: Page }): Promise<void> => {
  const vaultPage = new VaultPage({ page })
  const vaultUrl = `${config.baseUrl}/vault`

  await page.waitForURL((url) => url.href.startsWith(vaultUrl))
  await vaultPage.vaultBreadcrumb.waitFor()
}

export const userEntersDriveMode = async ({ page }: { page: Page }): Promise<void> => {
  const vaultPage = new VaultPage({ page })

  await vaultPage.vaultOption.click()
  await vaultPage.driveOption.click()
  await vaultPage.driveBreadcrumb.waitFor({ state: 'visible' })
}
