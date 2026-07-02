import { Page } from '@playwright/test'
import util from 'util'
import { objects } from '../../../index'
import { a11y } from '../../index'

const searchResultMessageSelector = '//p[@class="oc-text-muted"]'
const selectTagDropdownSelector =
  '//div[contains(@class,"files-search-filter-tags")]//button[contains(@class,"oc-filter-chip-button")]'
const tagFilterChipSelector = '//button[contains(@data-test-value,"%s")]'
const mediaTypeFilterSelector = '.item-filter-mediaType'
const mediaTypeFilterItem = '[data-test-id="media-type-%s"]'
const mediaTypeOutside = '.files-search-result-filter'
const clearFilterSelector = '.item-filter-%s .oc-filter-chip-clear'
const lastModifiedFilterSelector = '.item-filter-lastModified'
const lastModifiedFilterItem = '[data-test-value="%s"]'
const enableSearchTitleOnlySelector =
  '//div[contains(@class,"files-search-filter-title-only")]//button[contains(@class,"item-inline-filter-option") and contains(@id,"true")]'
const disableSearchTitleOnlySelector =
  '//div[contains(@class,"files-search-filter-title-only")]//button[contains(@class,"item-inline-filter-option") and contains(@id,"false")]'

export const getSearchResultMessage = ({ page }: { page: Page }): Promise<string> => {
  return page.locator(searchResultMessageSelector).innerText()
}

export const selectTagFilter = async ({
  tag,
  page
}: {
  tag: string
  page: Page
}): Promise<void> => {
  await page.locator(selectTagDropdownSelector).click()
  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().includes('/dav/spaces') &&
        resp.status() === 207 &&
        resp.request().method() === 'REPORT'
    ),
    page.locator(util.format(tagFilterChipSelector, tag)).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'Text editor Save button is disabled after saving'
  )
}

export const selectMediaTypeFilter = async ({
  mediaType,
  page
}: {
  mediaType: string
  page: Page
}): Promise<void> => {
  await page.locator(mediaTypeFilterSelector).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'Media type filter dropdown'
  )
  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().includes('/dav/spaces') &&
        resp.status() === 207 &&
        resp.request().method() === 'REPORT'
    ),
    page.locator(util.format(mediaTypeFilterItem, mediaType.toLowerCase())).click()
  ])
  await page.locator(mediaTypeOutside).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesView'],
    'Files view after selecting media type filter'
  )
}

export const selectLastModifiedFilter = async ({
  lastModified,
  page
}: {
  lastModified: string
  page: Page
}): Promise<void> => {
  await page.locator(lastModifiedFilterSelector).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'Last modified filter dropdown'
  )
  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().includes('/dav/spaces') &&
        resp.status() === 207 &&
        resp.request().method() === 'REPORT'
    ),
    page.locator(util.format(lastModifiedFilterItem, lastModified)).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesView'],
    'Files view after selecting last modified filter'
  )
}

export const clearFilter = async ({
  page,
  filter
}: {
  page: Page
  filter: string
}): Promise<void> => {
  await page.locator(util.format(clearFilterSelector, filter)).click()
  await a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesView'],
    `Files view after clearing ${filter} filter`
  )
}

export const toggleSearchTitleOnly = async ({
  enableOrDisable,
  page
}: {
  enableOrDisable: string
  page: Page
}): Promise<void> => {
  const selector =
    enableOrDisable === 'enable' ? enableSearchTitleOnlySelector : disableSearchTitleOnlySelector
  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().includes('/dav/spaces') &&
        resp.status() === 207 &&
        resp.request().method() === 'REPORT'
    ),
    page.locator(selector).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['files'],
    'search title only toggle button before toggling'
  )
}
