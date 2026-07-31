import { Page } from '@playwright/test'
import { objects } from '../../../index'

const groupsNavSelector = '//a[@data-nav-name="admin-settings-groups"]'
const appLoadingSpinnerSelector = '#app-loading-spinner'
export class Groups {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  async navigate(): Promise<void> {
    await this.#page.locator(groupsNavSelector).click()
    await this.#page.locator(appLoadingSpinnerSelector).waitFor({ state: 'detached' })
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      ['body'],
      'Groups Management page'
    )
  }
}
