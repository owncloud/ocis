import { Page } from '@playwright/test'
import { a11y } from '../../../index'

const sharesNavSelector = '//a[@data-nav-name="files-shares"]'

export class WithOthers {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  async navigate(): Promise<void> {
    await this.#page.locator(sharesNavSelector).click()
    await a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      ['files', 'sidebarNavigationMenu'],
      'files view and sidebar after navigating to shares section'
    )
    await this.#page.getByText('Shared with others').click()
    await a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      ['files'],
      'files view after navigating to shared with others'
    )
  }
}
