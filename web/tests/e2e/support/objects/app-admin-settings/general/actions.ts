import { Page, expect, test } from '@playwright/test'
import { readFileSync } from 'node:fs'
import { objects } from '../../../index'
import { getOtpFromImage } from '../../../utils/mfa'
import { Jimp } from 'jimp'

const waitForLogoWrapper = async (page: Page) => {
  const selectors = new objects.a11y.Accessibility({ page }).getSelectors()
  await page.locator(selectors.logoWrapper).waitFor()
  return selectors
}

export const uploadLogo = async (logoPath: string, page: Page): Promise<void> => {
  const selectors = await waitForLogoWrapper(page)
  const logoImg = page.locator(`${selectors.logoWrapper} img`)
  const srcBefore = await logoImg.getAttribute('src')

  await page.click('#logo-context-btn')

  // wait for the visible context menu and run accessibility scan on that menu
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBoxVisible'],
    'logo menu'
  )

  const logoInput = page.locator('#logo-upload-input')
  const uploadResponsePromise = page.waitForResponse(
    (resp) => resp.url().includes('/branding/logo') && resp.request().method() === 'POST'
  )
  const fileName = logoPath.split('/').pop() ?? 'logo.png'
  const extension = fileName.split('.').pop()?.toLowerCase()
  const mimeTypeByExtension: Record<string, string> = {
    png: 'image/png',
    jpg: 'image/jpeg',
    jpeg: 'image/jpeg',
    gif: 'image/gif'
  }

  await Promise.all([
    uploadResponsePromise,
    logoInput.setInputFiles({
      name: fileName,
      mimeType: mimeTypeByExtension[extension ?? ''] ?? 'application/octet-stream',
      buffer: readFileSync(logoPath)
    })
  ])
  const uploadResponse = await uploadResponsePromise

  const notification = page.locator('.oc-notification-message')
  await notification.waitFor()
  const notificationText = (await notification.textContent())?.trim() ?? ''

  if (!uploadResponse.ok() || !notificationText.includes('Logo was uploaded successfully')) {
    let responseText = ''
    try {
      responseText = await uploadResponse.text()
    } catch {
      responseText = '<unable to read response body>'
    }

    throw new Error(
      `Logo upload failed. status=${uploadResponse.status()} message="${notificationText}" body="${responseText}"`
    )
  }

  // The app triggers a delayed router.go(0) after successful upload.
  await page.waitForTimeout(1500)
  await page.reload({ waitUntil: 'load' })
  await waitForLogoWrapper(page)

  // run accessibility scan on the logo area after upload
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['logoWrapper'],
    'logo area after upload'
  )

  await expect(async () => {
    const logoSrc = await logoImg.getAttribute('src')
    expect(logoSrc).not.toEqual(srcBefore)
    expect(logoSrc).not.toContain('themes/owncloud/assets/oc_white.svg')
  }).toPass({ timeout: 30000 })
}

export const resetLogo = async (page: Page): Promise<void> => {
  const selectors = await waitForLogoWrapper(page)

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
  await waitForLogoWrapper(page)

  // run accessibility scan on the logo area after reset
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['logoWrapper'],
    'logo area after reset'
  )

  const imgAfter = page.locator(`${selectors.logoWrapper} img`)
  await expect(async () => {
    const srcAfter = await imgAfter.getAttribute('src')
    expect(srcAfter).not.toEqual(srcBefore)
  }).toPass({ timeout: 15000 })
}

export const userAuthenticatesWithOTP = async (page: Page, deviceName: string): Promise<void> => {
  const element = page.locator('#kc-totp-secret-qr-code')
  const qrCodePath = test.info().outputPath('qr.png')
  await element.screenshot({ path: qrCodePath })
  const image = await Jimp.read(qrCodePath)
  const { data, width, height } = image.bitmap
  const otp = await getOtpFromImage(data, width, height)
  await page.locator('#totp').fill(String(otp))
  await page.locator('#userLabel').fill(deviceName)
  await page.locator('#saveTOTPBtn').click()
}
