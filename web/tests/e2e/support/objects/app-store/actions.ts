import { Download, Page } from '@playwright/test'
import util from 'util'
import { objects } from '../../index'

const selectors = {
  appLoadingSpinner: '#app-loading-spinner',
  appSwitcherButton: '#_appSwitcherButton',
  appStoreMenuButton: 'data-test-id=app.app-store.menuItem',
  downloadButton: '//a[contains(., "%s")]/ancestor::li//button[.//span[text()="Download"]]',
  downloadVersionButton: '//tr[@data-item-id="%s"]//button[.//span[text()="Download"]]',
  appStoreHeadline: '.app-list-headline',
  appTileTitle: '.app-tile-title',
  selectAppTitle: '//a[contains(.,"%s")]',
  appDetailsBack: '.app-details-back',
  appDetailsTitle: '//h2[contains(@class, "app-details-title")][text()="%s"]',
  appsFilter: '#apps-filter',
  tag: '//button[contains(@class,"oc-tag")][span[text()="%s"]]',
  appTag: '//a[contains(.,"%s")]/following::button[contains(@class,"oc-tag")][span[text()="%s"]][1]'
}

export const openAppStore = async (args: { page: Page }): Promise<void> => {
  const { page } = args
  await page.locator(selectors.appSwitcherButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appSwitcherDropdown'],
    'app switcher dropdown'
  )
  await page.locator(selectors.appStoreMenuButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appStore'],
    'app store page'
  )
  await page.locator(selectors.appLoadingSpinner).waitFor({ state: 'detached' })
}
export const navigateToAppStoreOverview = async (args: { page: Page }): Promise<void> => {
  const { page } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appDetails'],
    'app details page'
  )
  await page.locator(selectors.appDetailsBack).click()
  await page.locator(selectors.appDetailsBack).waitFor({ state: 'detached' })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appStore'],
    'app store page'
  )
}

export const waitForAppStoreToBeVisible = async (args: { page: Page }): Promise<void> => {
  const { page } = args
  await page.locator(selectors.appStoreHeadline).waitFor()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appStore'],
    'app store page'
  )
}

export const getAppsList = async (args: { page: Page }): Promise<string[]> => {
  const { page } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appStore'],
    'app store page'
  )
  return page.locator(selectors.appTileTitle).allTextContents()
}

export const setSearchTerm = async (args: { page: Page; searchTerm: string }): Promise<void> => {
  const { page, searchTerm } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appStore'],
    'app store page'
  )
  return page.locator(selectors.appsFilter).fill(searchTerm)
}

export const selectAppTag = async (args: {
  page: Page
  app: string
  tag: string
}): Promise<void> => {
  const { page, app, tag } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appStore'],
    'app store page'
  )
  return page.locator(util.format(selectors.appTag, app, tag)).click()
}
export const selectTag = async (args: { page: Page; tag: string }): Promise<void> => {
  const { page, tag } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appStore'],
    'app store page'
  )
  return page.locator(util.format(selectors.tag, tag)).click()
}

export const selectApp = async (args: { page: Page; app: string }): Promise<void> => {
  const { page, app } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appStore'],
    'app store page'
  )
  await page.locator(util.format(selectors.selectAppTitle, app)).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appDetails'],
    'app details page'
  )
}

export const waitForAppDetailsToBeVisible = async (args: { page: Page; app }): Promise<void> => {
  const { page, app } = args
  await page.locator(util.format(selectors.appDetailsTitle, app)).waitFor()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appDetails'],
    'app details page'
  )
}

export const downloadAppVersion = async (args: {
  page: Page
  version: string
}): Promise<string> => {
  const { page, version } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appDetails'],
    'app details page'
  )
  const [download] = await Promise.all([
    page.waitForEvent('download'),
    page.locator(util.format(selectors.downloadVersionButton, version)).click()
  ])
  return download.suggestedFilename()
}

export const downloadApp = async (args: { page: Page; app: string }): Promise<Download> => {
  const { page, app } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appDetails'],
    'app details page'
  )
  const [download] = await Promise.all([
    page.waitForEvent('download'),
    page.locator(util.format(selectors.downloadButton, app)).click()
  ])

  return download
}
