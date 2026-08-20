import { Locator, Page } from '@playwright/test'

export class VaultPage {
  private readonly page: Page

  public readonly driveOption: Locator
  public readonly vaultOption: Locator
  public readonly authenticatorHeading: Locator
  public readonly qrImage: Locator
  public readonly oneTimeCodeTextbox: Locator
  public readonly otpSubmitButton: Locator
  public readonly driveBreadcrumb: Locator
  public readonly vaultBreadcrumb: Locator

  public constructor({ page }: { page: Page }) {
    this.page = page

    this.driveOption = page.getByRole('button', { name: 'DRIVE' })
    this.vaultOption = page.getByRole('button', { name: 'VAULT' })
    this.authenticatorHeading = page.getByRole('heading', {
      name: 'Mobile Authenticator Setup'
    })
    this.qrImage = page.locator('#kc-totp-secret-qr-code')
    this.oneTimeCodeTextbox = page.locator('#totp')
    this.otpSubmitButton = page.getByRole('button', { name: 'Submit' })
    this.driveBreadcrumb = page.getByRole('link', { name: 'Drive' })
    this.vaultBreadcrumb = page.getByRole('link', { name: 'Vault' })
  }
}
