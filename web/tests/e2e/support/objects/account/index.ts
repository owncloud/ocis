import { Page } from '@playwright/test'
import * as po from './actions'

export class Account {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  getQuotaValue(): Promise<string> {
    return po.getQuotaValue({ page: this.#page })
  }

  getUserInfo(key: string): Promise<string> {
    return po.getUserInfo({ page: this.#page, key })
  }

  async openAccountPage(): Promise<void> {
    await po.openAccountPage({ page: this.#page })
  }

  async requestGdprExport(): Promise<void> {
    await po.requestGdprExport({ page: this.#page })
  }

  downloadGdprExport(): Promise<string> {
    return po.downloadGdprExport({ page: this.#page })
  }

  async changeLanguage(language: string, isAnonymousUser = false): Promise<void> {
    await po.changeLanguage({ page: this.#page, language, isAnonymousUser })
  }

  getTitle(): Promise<string> {
    return po.getTitle({ page: this.#page })
  }

  async disableNotificationEvent(event: string): Promise<void> {
    await po.disableNotificationEvent({ page: this.#page, event })
  }
}
