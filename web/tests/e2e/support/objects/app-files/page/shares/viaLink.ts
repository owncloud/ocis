import { Page } from '@playwright/test'
import { objects } from '../../../../index'

const sharesNavSelector = '//a[@data-nav-name="files-shares"]'

export class ViaLink {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  async navigate(): Promise<void> {
    await this.#page.locator(sharesNavSelector).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      ['files', 'sidebarNavigationMenu'],
      'shares via link page after navigation'
    )
    await this.#page.getByText('Shared via link').click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      this.#page,
      ['filesView'],
      'shares via link page after clicking shared via link'
    )
  }
}
