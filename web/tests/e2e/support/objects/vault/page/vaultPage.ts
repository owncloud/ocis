import { Locator, Page } from '@playwright/test'
import { config } from '../../../../config'

export class VaultPage {
  private readonly page: Page

  private readonly driveOption: Locator
  private readonly vaultOption: Locator
  public readonly authenticatorHeading: Locator
  private readonly qrImage: Locator
  private readonly oneTimeCodeTextbox: Locator
  private readonly otpSubmitButton: Locator
  public readonly driveBreadcrumb: Locator
  public readonly vaultBreadcrumb: Locator

  public constructor(page: Page) {
    this.page = page

    this.driveOption = page.getByRole('button', { name: 'DRIVE' })
    this.vaultOption = page.getByRole('button', { name: 'VAULT' })
    this.authenticatorHeading = page.getByRole('heading', { name: 'Mobile Authenticator Setup' })
    this.qrImage = page.getByRole('img', { name: 'Figure: Barcode' })
    this.oneTimeCodeTextbox = page.locator('#totp')
    this.otpSubmitButton = page.getByRole('button', { name: 'Submit' })
    this.driveBreadcrumb = page.getByRole('link', { name: 'Drive' })
    this.vaultBreadcrumb = page.getByRole('link', { name: 'Vault' })
  }

  public async userEntersVaultMode(): Promise<void> {
    await this.driveOption.click()
    await this.vaultOption.click()
    await this.qrImage.waitFor({ state: 'visible' })
  }

  public async captureQrCodeScreenshot(): Promise<Buffer> {
    return await this.qrImage.screenshot()
  }

  public async userAuthenticatesWithOTP(otp: string): Promise<void> {
    await this.oneTimeCodeTextbox.fill(otp)
    await this.otpSubmitButton.click()
    await this.waitForVaultMode()
  }

  public async waitForVaultMode(): Promise<void> {
    const vaultUrl = `${config.baseUrlOcis}/vault`

    await this.page.waitForURL(
      (url) => url.href.startsWith(vaultUrl),
      {
        timeout: 60000,
        waitUntil: 'domcontentloaded'
      }
    )
    await this.vaultBreadcrumb.waitFor({
      state: 'visible',
      timeout: 60000
    })
  }

  public async userEntersDriveMode(): Promise<void> {
    await this.vaultOption.click()
    await this.driveOption.click()
    await this.driveBreadcrumb.waitFor({ state: 'visible' })
  }
}
