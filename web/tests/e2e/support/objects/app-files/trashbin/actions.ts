import { Page } from '@playwright/test'
import util from 'util'
import { objects } from '../../../index'

const spaceIdSelector = '//tr[@data-item-id="%s"]//a'
export interface openTrashBinArgs {
  id: string
  page: Page
}
export const openTrashbin = async (args: openTrashBinArgs): Promise<void> => {
  const { id, page } = args
  await page.locator(util.format(spaceIdSelector, id)).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesView'],
    'trashbin page of space'
  )
}
