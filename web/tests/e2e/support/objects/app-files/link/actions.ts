import { Page, expect } from '@playwright/test'
import util from 'util'
import { sidebar } from '../utils'
import { getActualExpiryDate } from '../../../utils/datePicker'
import { clickResource } from '../resource/actions'
import { objects } from '../../../index'

export interface createLinkArgs {
  page: Page
  role?: string
  resource?: string
  name?: string
  space?: boolean
  password?: string
}

export interface copyLinkArgs {
  page: Page
  resource: string
  resourceType?: string
  name?: string
  via?: string
}

export type changeNameArgs = {
  page: Page
  resource?: string
  newName: string
  space?: boolean
}

export type addExpirationArgs = {
  page: Page
  resource: string
  linkName: string
  expireDate: string
}

export type addPasswordArgs = {
  page: Page
  resource: string
  linkName: string
  newPassword: string
}

export type changeRoleArgs = {
  page: Page
  resource?: string
  linkName: string
  role: string
  space?: boolean
}

export type deleteLinkArgs = {
  page: Page
  resourceName: string
  name: string
}

export type publicLinkAndItsEditButtonVisibilityArgs = {
  page: Page
  linkName: string
  resource?: string
  space?: boolean
}
const publicLinkRoleToggle = `//button[contains(@class, "link-role-dropdown-toggle")]`
const publicLinkSetRoleButton = `//span[contains(@class,"role-dropdown-list-option-label") and text()='%s']`
const linkExpiryDatepicker = '.oc-modal-body .oc-date-picker input'
const linkExpiryDatepickerConfirmButton = '.oc-modal-body-actions-confirm'
const publicLinkEditRoleButton = `//span[contains(@class, "files-links-name") and text()="%s"]//ancestor::li//button[contains(@class, "link-role-dropdown-toggle")]`
const addPublicLinkButton = '#files-file-link-add'
const publicLinkNameList =
  '//div[@id="oc-files-file-link"]//ul//span[contains(@class, "files-links-name")]'
const publicLink = `//ul//h4[text()='%s']/following-sibling::div//p`
const publicLinkCurrentRole =
  '//button[contains(@class,"link-role-dropdown-toggle")]//span[contains(@class,"link-current-role")]'
const linkUpdateDialog = '//div[contains(@class,"oc-notification-message-title")]'
const editPublicLinkButton =
  '//span[contains(@class, "files-links-name") and text()="%s"]//ancestor::li//button[contains(@class, "edit-drop-trigger")]'
const editPublicLinkRenameButton =
  '//div[contains(@id,"edit-public-link-dropdown")]//button/span[text()="Rename"]'
const editPublicLinkSetExpirationButton =
  '//div[contains(@id,"edit-public-link-dropdown")]//button/span[text()="Set expiration date"]'
const editPublicLinkAddPasswordButton =
  '//div[contains(@id,"edit-public-link-dropdown")]//button/span[text()="Edit password"]'
const editPublicLinkNameInput = '.oc-modal-body input.oc-text-input'
const editPublicLinkPasswordInput = '.oc-modal-body .oc-text-input-password-wrapper input'
const editPublicLinkRenameConfirm = '.oc-modal-body-actions-confirm'
const deleteLinkButton =
  `//span[contains(@class, "files-links-name") and text()="%s"]` +
  `//ancestor::li//button/span[text()="Delete link"]`
const confirmDeleteButton = `//button[contains(@class,"oc-modal-body-actions-confirm") and text()="Delete"]`
const notificationContainer = 'div.oc-notification'
const publicLinkPasswordErrorMessage = `//div[contains(@class, "oc-text-input-message oc-text-input-danger")]/span`
const cancelButton = '.oc-modal-body-actions-cancel'
const showOrHidePasswordButton = '.oc-text-input-show-password-toggle'
const copyPasswordButton = '.oc-text-input-copy-password-button'
const generatePasswordButton = '.oc-text-input-generate-password-button'
const expectedRegexForGeneratedPassword = /^[A-Za-z0-9\s\S]{12}$/
const passwordInputDescription = '.oc-text-input-description .oc-text-input-description'
const copyLinkButton =
  '//span[contains(@class, "files-links-name") and text()="%s"]//ancestor::li//button[contains(@class, "oc-files-public-link-copy-url")]'
const folderModalIframe = '#iframe-folder-view'
const closeFolderModalButton = '.oc-modal-body-actions-cancel'

const getRecentLinkUrl = async (page: Page, name: string): Promise<string> => {
  await page.locator(util.format(copyLinkButton, name)).click()
  const handle = await page.evaluateHandle(() => navigator.clipboard.readText())
  return handle.jsonValue()
}

const getRecentLinkName = async (page: Page): Promise<string> => {
  return await page.locator(publicLinkNameList).first().textContent()
}

export const createLink = async (args: createLinkArgs): Promise<string> => {
  const { space, page, resource, password, role } = args
  if (!space) {
    const resourcePaths = resource.split('/')
    const resourceName = resourcePaths.pop()
    if (resourcePaths.length) {
      await clickResource({ page: page, path: resourcePaths.join('/') })
    }
    await sidebar.open({ page: page, resource: resourceName })
    await sidebar.openPanel({ page: page, name: 'sharing' })
  }
  await page.locator(addPublicLinkButton).click()

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'create public link modal'
  )

  if (role) {
    await page.locator(publicLinkRoleToggle).click()
    await page.locator(util.format(publicLinkSetRoleButton, role)).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['tippyBox'],
      'link role dropdown tippy box'
    )
  }

  await page.locator(editPublicLinkPasswordInput).fill(password)

  await Promise.all([
    page.waitForResponse(
      (res) =>
        res.url().includes('createLink') &&
        res.request().method() === 'POST' &&
        res.status() === 200
    ),
    page.locator(editPublicLinkRenameConfirm).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appSidebar'],
    'sidebar panel'
  )
  await clearCurrentPopup(page)
  return await getRecentLinkUrl(page, 'Unnamed link')
}

export const changeRole = async (args: changeRoleArgs): Promise<string> => {
  const { page, resource, linkName, role, space } = args

  // clear all popups
  await clearAllPopups(page)

  let resourceName = null
  let shareType = 'space-share'
  if (!space) {
    const resourcePaths = resource.split('/')
    resourceName = resourcePaths.pop()
    shareType = 'sharing'
    if (resourcePaths.length) {
      await clickResource({ page: page, path: resourcePaths.join('/') })
    }
  }
  await sidebar.open({ page: page, resource: resourceName })
  await sidebar.openPanel({ page: page, name: shareType })
  await page.locator(util.format(publicLinkEditRoleButton, linkName)).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'Edit public link role dropdown'
  )

  await Promise.all([
    page.waitForResponse(
      (res) =>
        res.url().includes('permissions') &&
        res.request().method() === 'PATCH' &&
        res.status() === 200
    ),
    page.locator(util.format(publicLinkSetRoleButton, role)).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    [linkUpdateDialog],
    'Edit public link role success message'
  )
  const message = await page.locator(linkUpdateDialog).textContent()
  expect(message.trim()).toBe('Link was updated successfully')
  return await page.locator(publicLinkCurrentRole).textContent()
}

export const changeName = async (args: changeNameArgs): Promise<string> => {
  const { page, resource, space, newName } = args

  // clear all popups
  await clearAllPopups(page)

  if (!space) {
    const resourcePaths = resource.split('/')
    const resourceName = resourcePaths.pop()
    if (resourcePaths.length) {
      await clickResource({ page: page, path: resourcePaths.join('/') })
    }
    await sidebar.open({ page: page, resource: resourceName })
    await sidebar.openPanel({ page: page, name: 'sharing' })
  }
  await page.locator(util.format(editPublicLinkButton, 'Unnamed link')).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'edit public link tippy box'
  )
  await page.locator(editPublicLinkRenameButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'update link modal'
  )
  await page.locator(editPublicLinkNameInput).fill(newName)
  await page.locator(editPublicLinkRenameConfirm).click()
  const message = await page.locator(linkUpdateDialog).textContent()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['linkUpdateDialog'],
    'edit public link success message'
  )
  expect(message.trim()).toBe('Link was updated successfully')
  return await getRecentLinkName(page)
}

export const fillPassword = async (args: addPasswordArgs): Promise<void> => {
  const { page, resource, linkName, newPassword } = args

  // clear all popups
  await clearAllPopups(page)

  const resourcePaths = resource.split('/')
  const resourceName = resourcePaths.pop()
  if (resourcePaths.length) {
    await clickResource({ page: page, path: resourcePaths.join('/') })
  }
  await sidebar.open({ page: page, resource: resourceName })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appSidebar'],
    'sidebar after opening the resource'
  )
  await sidebar.openPanel({ page: page, name: 'sharing' })
  await page.locator(util.format(editPublicLinkButton, linkName)).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'edit public link tippy box'
  )
  await page.locator(editPublicLinkAddPasswordButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'edit public link password modal'
  )
  await page.locator(editPublicLinkPasswordInput).fill(newPassword)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'edit public link password modal after filling password'
  )
  await page.locator(editPublicLinkRenameConfirm).click()
}

export const addPassword = async (args: addPasswordArgs): Promise<void> => {
  const { page } = args

  await fillPassword(args)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['notificationContainer'],
    'notification popup after updating password'
  )
  const message = await page.locator(linkUpdateDialog).textContent()
  expect(message.trim()).toBe('Link was updated successfully')
}

export const showOrHidePassword = async (args: {
  page: Page
  showOrHide: string
}): Promise<void> => {
  const { page, showOrHide } = args
  await page.locator(showOrHidePasswordButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    `${showOrHide} password of public link modal`
  )
  showOrHide === 'reveals'
    ? await expect(page.locator(editPublicLinkPasswordInput)).toHaveAttribute('type', 'text')
    : await expect(page.locator(editPublicLinkPasswordInput)).toHaveAttribute('type', 'password')
}

export const copyEnteredPassword = async (page: Page): Promise<string> => {
  const enteredPassword = await page.locator(editPublicLinkPasswordInput).inputValue()
  await page.locator(copyPasswordButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'copy password of public link modal'
  )
  // Return entered password directly (clipboard may have wrong data in parallel tests)
  return enteredPassword
}

export const generatePassword = async (page: Page): Promise<void> => {
  await page.locator(generatePasswordButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'generate password for public link modal'
  )
  const generatedPassword = await page.locator(editPublicLinkPasswordInput).inputValue()
  expect(generatedPassword).toMatch(expectedRegexForGeneratedPassword)
}

export const setPassword = async (page: Page): Promise<void> => {
  await Promise.all([
    page.waitForResponse(
      (res) =>
        res.url().includes('permissions') &&
        res.request().method() === 'POST' &&
        res.status() === 200
    ),
    page.locator(editPublicLinkRenameConfirm).click()
  ])
}

export const addExpiration = async (args: addExpirationArgs): Promise<void> => {
  const { page, resource, linkName, expireDate } = args
  const resourcePaths = resource.split('/')
  const resourceName = resourcePaths.pop()
  if (resourcePaths.length) {
    await clickResource({ page: page, path: resourcePaths.join('/') })
  }
  await sidebar.open({ page: page, resource: resourceName })
  await sidebar.openPanel({ page: page, name: 'sharing' })
  await page.locator(util.format(editPublicLinkButton, linkName)).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'Edit public link dropdown'
  )
  await page.locator(editPublicLinkSetExpirationButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'Set expiration date modal'
  )

  const newExpiryDate = getActualExpiryDate(
    expireDate.toLowerCase().match(/[dayrmonthwek]+/)[0] as any,
    expireDate
  )

  await page.locator(linkExpiryDatepicker).fill(newExpiryDate.toISOString().split('T')[0])
  await page.locator(linkExpiryDatepickerConfirmButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    [linkUpdateDialog],
    'Set expiration date success message'
  )
  const message = await page.locator(linkUpdateDialog).textContent()
  expect(message.trim()).toBe('Link was updated successfully')
}

export const deleteLink = async (args: deleteLinkArgs): Promise<void> => {
  const { page, resourceName, name } = args

  // clear all popups
  await clearAllPopups(page)

  await sidebar.open({ page: page, resource: resourceName })
  await sidebar.openPanel({ page: page, name: 'sharing' })
  await page.locator(util.format(editPublicLinkButton, name)).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'Edit public link dropdown'
  )
  await page.locator(util.format(deleteLinkButton, name)).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'Delete public link confirmation modal'
  )
  await page.locator(confirmDeleteButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    [linkUpdateDialog],
    'Delete public link success message'
  )
  const message = await page.locator(linkUpdateDialog).textContent()
  expect(message.trim()).toBe('Link was deleted successfully')
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['sidebarPaneSharing'],
    'Public link list after deletion'
  )
}

export const getPublicLinkVisibility = async (
  args: publicLinkAndItsEditButtonVisibilityArgs
): Promise<string> => {
  const { page, linkName, resource, space } = args
  let shareType = 'space-share'
  let resourceName = null
  if (!space) {
    shareType = 'sharing'
    const resourcePaths = resource.split('/')
    resourceName = resourcePaths.pop()
    if (resourcePaths.length) {
      await clickResource({ page: page, path: resourcePaths.join('/') })
    }
  }
  await sidebar.open({ page: page, resource: resourceName })
  await sidebar.openPanel({ page: page, name: shareType })
  return await page.locator(util.format(publicLink, linkName)).textContent()
}

export const getLinkEditButtonVisibility = async (
  args: publicLinkAndItsEditButtonVisibilityArgs
): Promise<boolean> => {
  const { page, linkName } = args
  return await page.locator(util.format(editPublicLinkButton, linkName)).isVisible()
}

export const clearAllPopups = async (page: Page): Promise<void> => {
  const count = await page.locator(notificationContainer).evaluate((container) => {
    Object.values(container.children).forEach((child) => {
      container.removeChild(child)
    })
    return container.children.length
  })
  if (count) {
    throw new Error(`Failed to clear ${count} notifications`)
  }
  await expect(page.locator(linkUpdateDialog)).not.toBeVisible()
}

export const clearCurrentPopup = async (page: Page): Promise<void> => {
  await expect(page.locator(linkUpdateDialog)).toBeVisible()
  await clearAllPopups(page)
}

export const getPublicLinkPasswordErrorMessage = async (
  page: Page,
  passwordProtectedFolder?: boolean
): Promise<string> => {
  return passwordProtectedFolder
    ? await page.frameLocator(folderModalIframe).locator(publicLinkPasswordErrorMessage).innerText()
    : await page.locator(publicLinkPasswordErrorMessage).innerText()
}

export const clickOnCancelButton = async (page: Page): Promise<void> => {
  await page.locator(cancelButton).click()
}

export const copyLinkToClipboard = async (args: copyLinkArgs): Promise<string> => {
  const { page, resource, resourceType } = args
  await sidebar.open({ page: page, resource: resource, resourceType })
  await sidebar.openPanel({ page: page, name: 'sharing', resourceType })

  // clear the clipboard
  await page.evaluate(`navigator.clipboard.writeText('')`)

  if (resourceType === 'passwordProtectedFolder') {
    await page.frameLocator(folderModalIframe).getByLabel('Copy link to clipboard').click()
  } else {
    await page.getByLabel('Copy link to clipboard').click()
  }
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['.snackbar'],
    'copy link to clipboard notification'
  )
  return await page.evaluate('navigator.clipboard.readText()')
}

export const closeFolderModal = async (page: Page): Promise<void> => {
  await page.locator(closeFolderModalButton).click()
}
