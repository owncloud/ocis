import { expect, Page } from '@playwright/test'
import { objects } from '../../../../index'

// rendered by Projects.vue when the user has access to no project space
const emptySpacesSelector = '#files-spaces-empty'

export class Projects {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  async expectToBeEmpty(): Promise<void> {
    await expect(this.#page.locator(emptySpacesSelector)).toBeVisible()
  }

  async navigate(): Promise<void> {
    await this.#page.locator('//a[@data-nav-name="files-spaces-projects"]').click()
    await this.#page.locator('#app-loading-spinner').waitFor({ state: 'detached' })
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      ['filesView'],
      'spaces page'
    )
  }
}
