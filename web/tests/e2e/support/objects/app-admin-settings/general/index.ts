import { Page } from '@playwright/test'
import * as po from './actions'

export class General {
  #page: Page
  constructor({ page }: { page: Page }) {
    this.#page = page
  }
  async uploadLogo({ path }: { path: string }): Promise<void> {
    await po.uploadLogo(path, this.#page)
  }
  async resetLogo(): Promise<void> {
    await po.resetLogo(this.#page)
  }
  async userAuthenticatesWithOTP({ deviceName }: { deviceName: string }): Promise<void> {
    await po.userAuthenticatesWithOTP(this.#page, deviceName)
  }
}
