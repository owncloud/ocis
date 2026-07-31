import { Page, expect } from '@playwright/test'
import util from 'util'

import { sidebar, editor } from '../utils'
import Collaborator, { ICollaborator } from '../share/collaborator'
import { createLink } from '../link/actions'
import { File } from '../../../types'
import { objects } from '../../../index'

const newSpaceMenuButton = '#new-space-menu-btn'
const spaceNameInputField = '.oc-modal input'
const actionConfirmButton = '.oc-modal-body-actions-confirm'
const spaceIdSelector = `[data-item-id="%s"] .oc-resource-basename`
const spacesRenameOptionSelector = '.oc-files-actions-rename-trigger:visible'
const editSpacesSubtitleOptionSelector = '.oc-files-actions-edit-description-trigger:visible'
const editQuotaOptionSelector = '.oc-files-actions-edit-quota-trigger:visible'
const editImageOptionSelector = '.oc-files-actions-upload-space-image-trigger:visible'
const downloadSpaceSelector = '.oc-files-actions-download-archive-trigger:visible'
const spacesQuotaSearchField = '.oc-modal .vs__search'
const selectedQuotaValueField = '.vs--open'
const quotaValueDropDown = `.vs__dropdown-option :text-is("%s")`
const editSpacesDescription = '.oc-files-actions-edit-readme-content-trigger:visible'
const spacesDescriptionInputArea = '.cm-content'
const spacesDescriptionSaveTextFileInEditorButton = '#app-save-action:visible'
const activitySidebarPanel = 'sidebar-panel-activities'
const activitySidebarPanelBodyContent = '#sidebar-panel-activities .sidebar-panel__body-content'

export const openActionsPanel = async (page: Page): Promise<void> => {
  await sidebar.open({ page })
  await sidebar.openPanel({ page, name: 'space-actions' })
}

export const openSharingPanel = async (page: Page): Promise<void> => {
  await sidebar.open({ page })
  await sidebar.openPanel({ page, name: 'space-share' })
}

export const openActivitiesPanel = async (page: Page): Promise<void> => {
  await sidebar.open({ page })
  await sidebar.openPanel({ page, name: 'activities' })
}

/**/

export interface createSpaceArgs {
  name: string
  page: Page
}

export const createSpace = async (args: createSpaceArgs): Promise<string> => {
  const { page, name } = args

  await page.locator(newSpaceMenuButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['ocModal'], 'spaces page')
  await page.locator(spaceNameInputField).fill(name)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['ocModal'], 'spaces page')

  const [responses] = await Promise.all([
    page.waitForResponse(
      (postResp) =>
        postResp.status() === 201 &&
        postResp.request().method() === 'POST' &&
        postResp.url().endsWith('drives?template=default')
    ),
    page.locator(actionConfirmButton).click()
  ])

  // createSpace runs from both the Files app spaces view (#files-view) and the
  // admin-settings spaces page (#admin-settings-wrapper), so scan the always-present body
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['body'], 'spaces page')

  const { id } = await responses.json()
  return id
}

/**/

export interface openSpaceArgs {
  id: string
  page: Page
}

export const openSpace = async (args: openSpaceArgs): Promise<void> => {
  const { page, id } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['filesView'], 'spaces page')
  const locator = page.locator(util.format(spaceIdSelector, id))
  await locator.click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['filesView'], 'spaces page')
}
/**/

export const changeSpaceName = async (args: {
  page: Page
  id: string
  value: string
}): Promise<void> => {
  const { page, value, id } = args
  await openActionsPanel(page)

  await page.locator(spacesRenameOptionSelector).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'rename space modal'
  )
  await page.locator(spaceNameInputField).fill(value)
  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith(encodeURIComponent(id)) &&
        resp.status() === 200 &&
        resp.request().method() === 'PATCH'
    ),
    page.locator(actionConfirmButton).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['notificationContainer'],
    'notification popup after renaming space'
  )
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['spaceInfoContainer'],
    'space information container after renaming space'
  )

  await sidebar.close({ page: page })
}

/**/

export const changeSpaceSubtitle = async (args: {
  page: Page
  id: string
  value: string
}): Promise<void> => {
  const { page, value, id } = args
  await openActionsPanel(page)

  await page.locator(editSpacesSubtitleOptionSelector).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'edit space subtitle modal'
  )
  await page.locator(spaceNameInputField).fill(value)
  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith(encodeURIComponent(id)) &&
        resp.status() === 200 &&
        resp.request().method() === 'PATCH'
    ),
    page.locator(actionConfirmButton).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['notificationContainer'],
    'notification popup after changing space subtitle'
  )
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['spaceInfoContainer'],
    'space information container after changing space subtitle'
  )

  await sidebar.close({ page: page })
}

/**/

export const changeSpaceDescription = async (args: {
  page: Page
  value: string
}): Promise<void> => {
  const { page, value } = args
  await openActionsPanel(page)
  const waitForUpdate = () =>
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith('readme.md') &&
        resp.status() === 200 &&
        resp.request().method() === 'GET'
    )
  await Promise.all([waitForUpdate(), page.locator(editSpacesDescription).click()])

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['textEditor'],
    'space description editor'
  )

  await page.locator(spacesDescriptionInputArea).fill(value)
  await Promise.all([
    page.waitForResponse((resp) => resp.status() === 204 && resp.request().method() === 'PUT'),
    page.waitForResponse((resp) => resp.status() === 207 && resp.request().method() === 'PROPFIND'),
    page.locator(spacesDescriptionSaveTextFileInEditorButton).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['saveTextEditorOrViewerButton'],
    'save button in space description editor'
  )
  await editor.close(page)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['spaceDescriptionPreview'],
    'space description preview after editing space description'
  )
}

/**/

export const changeQuota = async (args: {
  id: string
  page: Page
  value: string
}): Promise<void> => {
  const { id, page, value } = args
  await openActionsPanel(page)

  await page.locator(editQuotaOptionSelector).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'edit space quota modal'
  )
  const searchLocator = page.locator(spacesQuotaSearchField)
  await searchLocator.pressSequentially(value)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['quotaValueDropDown'],
    'quota value dropdown after entering quota'
  )
  await page.locator(selectedQuotaValueField).waitFor()
  await page.locator(util.format(quotaValueDropDown, `${value} GB`)).click()

  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith(encodeURIComponent(id)) &&
        resp.status() === 200 &&
        resp.request().method() === 'PATCH'
    ),
    page.locator(actionConfirmButton).click()
  ])

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['notificationContainer'],
    'notification popup after changing space quota'
  )

  await sidebar.close({ page: page })
}

export interface SpaceMembersArgs {
  page: Page
  users: ICollaborator[]
}

export const addSpaceMembers = async (args: SpaceMembersArgs): Promise<void> => {
  const { page, users } = args
  await openSharingPanel(page)

  await Collaborator.inviteCollaborators({ page, collaborators: users })
  await sidebar.close({ page: page })
}

export const changeSpaceImage = async (args: {
  id: string
  page: Page
  resource: File
}): Promise<void> => {
  const { id, page, resource } = args
  await openActionsPanel(page)

  const [fileChooser] = await Promise.all([
    page.waitForEvent('filechooser'),
    page.locator(editImageOptionSelector).click()
  ])

  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith(encodeURIComponent(id)) &&
        resp.status() === 200 &&
        resp.request().method() === 'PATCH'
    ),
    page.waitForResponse(
      (resp) =>
        resp.url().includes(resource.name) &&
        resp.status() === 200 &&
        resp.request().method() === 'GET'
    ),
    fileChooser.setFiles(resource.path)
  ])

  await sidebar.close({ page: page })
}

export interface removeAccessMembersArgs extends Omit<SpaceMembersArgs, 'users'> {
  users: Omit<ICollaborator, 'role'>[]
  removeOwnSpaceAccess?: boolean
}

export const removeAccessSpaceMembers = async (args: removeAccessMembersArgs): Promise<void> => {
  const { page, users, removeOwnSpaceAccess } = args
  await openSharingPanel(page)

  for (const collaborator of users) {
    await Collaborator.removeCollaborator({ page, collaborator, removeOwnSpaceAccess })
  }
}

export const changeSpaceRole = async (args: SpaceMembersArgs): Promise<void> => {
  const { page, users } = args
  await openSharingPanel(page)

  for (const collaborator of users) {
    await Promise.all([
      page.waitForResponse(
        (resp) =>
          resp.url().includes('permissions') &&
          resp.status() === 200 &&
          resp.request().method() === 'PATCH'
      ),
      Collaborator.changeCollaboratorRole({ page, collaborator })
    ])
  }
}

export const createPublicLinkForSpace = async (args: {
  page: Page
  password: string
}): Promise<string> => {
  const { page, password } = args
  await openSharingPanel(page)
  return createLink({ page: page, space: true, password: password })
}

export const addExpirationDateToMember = async (args: {
  page: Page
  member: Omit<ICollaborator, 'role'>
  expirationDate: string
}): Promise<void> => {
  const { page, member: collaborator, expirationDate } = args
  await openSharingPanel(page)
  await Collaborator.setExpirationDateForCollaborator({ page, collaborator, expirationDate })
}

export const removeExpirationDateFromMember = async (args: {
  page: Page
  member: Omit<ICollaborator, 'role'>
}): Promise<void> => {
  const { page, member: collaborator } = args
  await openSharingPanel(page)
  await Collaborator.removeExpirationDateFromCollaborator({ page, collaborator })
}

export const downloadSpace = async (page: Page): Promise<string> => {
  await openActionsPanel(page)
  const [download] = await Promise.all([
    page.waitForEvent('download'),
    page.locator(downloadSpaceSelector).click()
  ])
  await sidebar.close({ page })

  return download.suggestedFilename()
}

export const checkSpaceActivity = async ({
  page,
  activity
}: {
  page: Page
  activity: string
}): Promise<void> => {
  await openActivitiesPanel(page)
  await expect(page.getByTestId(activitySidebarPanel)).toBeVisible()
  await expect(page.locator(activitySidebarPanelBodyContent)).toContainText(activity)
}
