import { Locator, Page } from '@playwright/test'

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
    this.driveBreadcrumb = page.getByRole('link', { name: 'Drive', exact: true })
    this.vaultBreadcrumb = page.getByRole('link', { name: 'Vault', exact: true })
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
    await this.vaultBreadcrumb.waitFor({ state: 'visible' })
  }

  public async userEntersDriveMode(): Promise<void> {
    await this.vaultOption.click()
    await this.driveOption.click()
    await this.driveBreadcrumb.waitFor({ state: 'visible' })
  }
}
