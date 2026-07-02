import { Page } from '@playwright/test'

const personalSpaceNavSelector = '//a[@data-nav-name="files-spaces-generic"]'

export class Personal {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  async navigate(): Promise<void> {
    await this.#page.locator(personalSpaceNavSelector).click()
  }
}
