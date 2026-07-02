import { basename } from 'path'
import { Page, expect } from '@playwright/test'
import { objects } from '../../../index'
import { getOtpFromImage } from '../../../utils/mfa'
import { Jimp } from 'jimp'

export const uploadLogo = async (path: string, page: Page): Promise<void> => {
  await page.click('#logo-context-btn')

  // wait for the visible context menu and run accessibility scan on that menu
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBoxVisible'],
    'logo menu'
  )

  const logoInput = page.locator('#logo-upload-input')
  await logoInput.setInputFiles(path)

  await page.locator('.oc-notification-message').waitFor()
  await page.reload()
  const selectors = new objects.a11y.Accessibility({ page }).getSelectors()
  // run accessibility scan on the logo area after upload
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['logoWrapper'],
    'logo area after upload'
  )

  const logoImg = page.locator(`${selectors.logoWrapper} img`)
  const logoSrc = await logoImg.getAttribute('src')
  expect(logoSrc).toContain(basename(path))
}

export const resetLogo = async (page: Page): Promise<void> => {
  const a11yObject = new objects.a11y.Accessibility({ page })
  const selectors = a11yObject.getSelectors()

  const imgBefore = page.locator(`${selectors.logoWrapper} img`)
  const srcBefore = await imgBefore.getAttribute('src')
  await page.click('#logo-context-btn')

  // wait for the visible context menu and run accessibility scan on that menu
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBoxVisible'],
    'logo menu'
  )

  await page.click('.oc-general-actions-reset-logo-trigger')

  await page.locator('.oc-notification-message').waitFor()
  await page.reload()

  // run accessibility scan on the logo area after reset
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['logoWrapper'],
    'logo area after reset'
  )

  const imgAfter = page.locator(`${selectors.logoWrapper} img`)
  const srcAfter = await imgAfter.getAttribute('src')
  expect(srcAfter).not.toEqual(srcBefore)
}

export const userAuthenticatesWithOTP = async (page: Page, deviceName: string): Promise<void> => {
  const element = page.locator('#kc-totp-secret-qr-code')
  await element.screenshot({ path: 'qr.png' })
  const image = await Jimp.read('./qr.png')
  const { data, width, height } = image.bitmap
  const otp = await getOtpFromImage(data, width, height)
  await page.locator('#totp').fill(String(otp))
  await page.locator('#userLabel').fill(deviceName)
  await page.locator('#saveTOTPBtn').click()
}
