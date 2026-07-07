import { Page, expect } from '@playwright/test'
import util from 'util'
import { config } from '../../../config'
import { objects } from '../../index'

const accountMenuButton = '.oc-topbar-avatar'
const quotaValue = '.quota-information-text'
const accountManageButton = '#oc-topbar-account-manage'
const infoValue = '.account-page-info-%s td:nth-child(2)'
const requestExportButton = '[data-testid="request-export-btn"]'
const downloadExportButton = '[data-testid="download-export-btn"]'
const languageInput = '[data-testid="language"] .vs__search'
const languageValueDropDown = `.vs__dropdown-menu :text-is("%s")`
const languageValue = '[data-testid="language"] .vs__selected'
const accountPageTitle = '#account-page-title'
const notificationEventCheckBox =
  '//div[@class="account-table"]//td[text()="%s"]//following-sibling::td//span//input[@aria-label="In-App"]'

export const getQuotaValue = async (args: { page: Page }): Promise<string> => {
  const { page } = args
  await page.reload()
  await page.locator(accountMenuButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['accountInfoContainer'],
    'account info modal'
  )
  const quotaText = await page.locator(quotaValue).textContent()
  await page.locator(quotaValue).click()

  // parse "0 B of 10 GB used"
  return quotaText.match(/\d+/g)?.[1]
}

export const getUserInfo = async (args: { page: Page; key: string }): Promise<string> => {
  const { page, key } = args
  await page.locator(accountMenuButton).click()
  await page.locator(accountManageButton).click()
  return await page.locator(util.format(infoValue, key)).textContent()
}

export const openAccountPage = async (args: { page: Page }): Promise<void> => {
  const { page } = args
  await page.locator(accountMenuButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['accountInfoContainer'],
    'account info modal'
  )
  await page.locator(accountManageButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['account'], 'account page')
}

export const requestGdprExport = async (args: { page: Page }): Promise<void> => {
  const { page } = args
  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith('exportPersonalData') &&
        resp.status() === 202 &&
        resp.request().method() === 'POST'
    ),
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith('.personal_data_export.json') &&
        resp.status() === 207 &&
        resp.request().method() === 'PROPFIND' &&
        resp.text().then((text) => text.includes('HTTP/1.1 200 OK')),
      // generating GDPR report can take a while
      // so we need to increase the timeout to 60 seconds
      { timeout: config.timeout * 1000 }
    ),
    page.locator(requestExportButton).click()
  ])
}

export const downloadGdprExport = async (args: { page: Page }): Promise<string> => {
  const { page } = args

  const [download] = await Promise.all([
    page.waitForEvent('download'),
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith('.personal_data_export.json') &&
        resp.status() === 200 &&
        resp.request().method() === 'HEAD'
    ),
    page.locator(downloadExportButton).click()
  ])
  await page.locator(requestExportButton).waitFor()
  return download.suggestedFilename()
}

export const changeLanguage = async (args: {
  page: Page
  language: string
  isAnonymousUser: boolean
}): Promise<void> => {
  const { page, language, isAnonymousUser } = args
  await page.locator(languageInput).waitFor()
  await page.locator(languageInput).click()
  await page.locator(languageInput).pressSequentially(language)
  const promises = []

  if (!isAnonymousUser) {
    promises.push(
      page.waitForResponse(
        (res) =>
          res.url().includes('graph/v1.0/me') &&
          res.request().method() === 'PATCH' &&
          res.status() === 200
      )
    )
  }

  promises.push(page.locator(util.format(languageValueDropDown, language)).press('Enter'))
  await Promise.all(promises)

  await expect(page.locator(languageValue)).toHaveText(language)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['account'], 'account page')
}

export const getTitle = async (args: { page: Page }): Promise<string> => {
  const { page } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['account'], 'account page')
  return page.locator(accountPageTitle).textContent()
}

export const disableNotificationEvent = async (args: {
  page: Page
  event: string
}): Promise<void> => {
  const { page, event } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['account'], 'account page')
  await page.locator(util.format(notificationEventCheckBox, event)).click()
}
