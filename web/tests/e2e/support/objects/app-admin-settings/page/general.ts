import { Page } from '@playwright/test'
import { objects } from '../../../index'

export class General {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  async navigate(): Promise<void> {
    await this.#page.locator('//a[@data-nav-name="admin-settings-general"]').click()
    await this.#page.locator('#app-loading-spinner').waitFor({ state: 'detached' })
    // run accessibility scan for the general management page body
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      ['body'],
      'General Management page'
    )
  }
}
