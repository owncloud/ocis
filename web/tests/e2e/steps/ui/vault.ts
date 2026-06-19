import { expect } from '@playwright/test'
import { getWorld } from '../../environment/world'
import { generateOtpFromScreenshot } from '../../support/utils/mfa'
import { config } from '../../config'
import { VaultPage } from '../../support/objects/vault/page/vaultPage'

/**
 * Switch user from Drive → Vault mode
 */
export async function userSwitchesToVaultMode({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const vaultPage = new VaultPage(page)
  await vaultPage.userEntersVaultMode()
}

/**
 * Assert user is redirected to the MFA authenticator page
 */
export async function userIsRedirectedToAuthenticatorPage({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const vaultPage = new VaultPage(page)
  console.log("keycloakpage urlllllllllllllllllllll", config.keycloakUrl);
  await expect(page).toHaveURL((url) => url.href.startsWith(config.keycloakUrl))
  await expect(vaultPage.authenticatorHeading).toBeVisible()
}

/**
 * Generate an OTP from the Vault QR code and authenticate the user
 */
export async function userAuthenticatesToVault({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const vaultPage = new VaultPage(page)
  const qrBuffer = await vaultPage.captureQrCodeScreenshot()
  const otp = await generateOtpFromScreenshot(qrBuffer)
  await vaultPage.userAuthenticatesWithOTP(otp)
}

/**
 * Assert user is in Vault mode
 */
export async function userIsInVaultMode({
    stepUser
  }: {
    stepUser: string
  }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const vaultPage = new VaultPage(page)
  const vaultPageUrl = `${config.baseUrlOcis}/vault`
  console.log("vaultpage urlllllllllllllllllllll", vaultPageUrl);
  await expect(page).toHaveURL((url) => url.href.startsWith(vaultPageUrl))
  await expect(vaultPage.vaultBreadcrumb).toBeVisible()
}

/**
 * Switch user from Vault → Drive mode
 */
export async function userSwitchesToDriveMode({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const vaultPage = new VaultPage(page)
  await vaultPage.userEntersDriveMode()
}

/**
 * Assert user is in Drive mode
 */
export async function userIsInDriveMode({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const vaultPage = new VaultPage(page)
  const drivePageUrl = `${config.baseUrlOcis}/files`
  console.log("drivepage urlllllllllllllllllllll", drivePageUrl);
  await expect(page).toHaveURL((url) => url.href.startsWith(drivePageUrl))
  await expect(vaultPage.driveBreadcrumb).toBeVisible()
}
