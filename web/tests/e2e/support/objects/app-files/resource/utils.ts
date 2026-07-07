import { errors, Page, Locator } from '@playwright/test'
import util from 'util'

const resourceNameSelector = '#files-space-table [data-test-resource-name="%s"]'
const showLinkShareButton =
  '//span[@data-test-resource-name="%s"]/ancestor::tr[contains(@class, "oc-tbody-tr")]//button[contains(@data-test-indicator-type, "%s")]'

/**
 * one of the few places where timeout should be used, as we also use this to detect the absence of an element
 * it is not possible to differentiate between `element not there yet` and `element not loaded yet`.
 *
 * @param page
 * @param name
 * @param timeout
 */
export const resourceExists = async ({
  page,
  name,
  timeout = 500
}: {
  page: Page
  name: string
  timeout?: number
}): Promise<boolean> => {
  let exist = true
  await page
    .locator(util.format(resourceNameSelector, name))
    .waitFor({ timeout })
    .catch((e) => {
      if (!(e instanceof errors.TimeoutError)) {
        throw e
      }

      exist = false
    })

  return exist
}

export const waitForResources = async ({
  page,
  names
}: {
  page: Page
  names: string[]
}): Promise<void> => {
  await Promise.all(
    names.map((name) => page.locator(util.format(resourceNameSelector, name)).waitFor())
  )
}

export const showShareIndicator = (args: {
  page: Page
  buttonLabel: string
  resource: string
}): Locator => {
  const { page, buttonLabel, resource } = args
  return page.locator(util.format(showLinkShareButton, resource, buttonLabel))
}
