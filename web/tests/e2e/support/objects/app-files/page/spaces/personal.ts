import { expect, Page } from '@playwright/test'

const personalSpaceNavSelector = '//a[@data-nav-name="files-spaces-generic"]'
// rendered by GenericSpace.vue when the current folder holds nothing
const emptySpaceSelector = '#files-space-empty'

export class Personal {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  async navigate(): Promise<void> {
    await this.#page.locator(personalSpaceNavSelector).click()
  }

  async expectToBeEmpty(): Promise<void> {
    await expect(this.#page.locator(emptySpaceSelector)).toBeVisible()
  }
}
