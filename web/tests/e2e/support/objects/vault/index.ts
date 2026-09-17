import { Page } from '@playwright/test'
import * as po from './actions'
import { VaultPage } from './page/VaultPage'

export { VaultPage }

export class VaultActions {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  async userEntersVaultMode(): Promise<void> {
    await po.userEntersVaultMode({ page: this.#page })
  }

  async waitForVaultMode(): Promise<void> {
    await po.waitForVaultMode({ page: this.#page })
  }

  async userEntersDriveMode(): Promise<void> {
    await po.userEntersDriveMode({ page: this.#page })
  }
}
