import { Locator, Page } from '@playwright/test'
import util from 'util'

const spaceIdSelector = `[data-item-id="%s"]`

export interface searchForSpacesIdsArgs {
  spaceID: string
  page: Page
}

export const spaceLocator = (args: searchForSpacesIdsArgs): Locator => {
  const { page, spaceID } = args
  return page.locator(util.format(spaceIdSelector, spaceID))
}
