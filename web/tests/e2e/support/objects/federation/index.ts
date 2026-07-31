import { Page } from '@playwright/test'
import * as po from './actions'

export class Federation {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }
  async generateInvitation(user: string): Promise<void> {
    await po.generateInvitation({ page: this.#page, user })
  }

  async acceptInvitation(sharer: string): Promise<void> {
    await po.acceptInvitation({ page: this.#page, sharer })
  }
  async connectionExists(info): Promise<boolean> {
    return await po.connectionExists({ page: this.#page, info })
  }
}
