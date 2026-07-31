import { Page } from '@playwright/test'
import * as po from './actions'

export class URLNavigation {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  async navigateToDetailsPanelOfResource(
    args: Omit<po.navigateToDetailsPanelOfResourceArgs, 'page'>
  ): Promise<void> {
    await po.navigateToDetailsPanelOfResource({ ...args, page: this.#page })
  }

  async openResourceViaUrl(args: Omit<po.openResourceViaUrlArgs, 'page'>): Promise<void> {
    await po.openResourceViaUrl({ ...args, page: this.#page })
  }

  async openSpaceViaUrl(args: Omit<po.openResourceViaUrlArgs, 'page'>): Promise<void> {
    await po.openSpaceViaUrl({ ...args, page: this.#page })
  }

  async navigateToNonExistingPage(): Promise<void> {
    await po.navigateToNonExistingPage({ page: this.#page })
  }
  async waitForNotFoundPageToBeVisible(): Promise<void> {
    await po.waitForNotFoundPageToBeVisible({ page: this.#page })
  }
}
