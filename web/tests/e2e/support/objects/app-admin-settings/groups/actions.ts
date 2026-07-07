import { Page } from '@playwright/test'
import util from 'util'
import { selectUser } from '../users/actions'
import { objects } from '../../../index'
import { fileAction } from '../../../../environment/constants'

const newGroupBtn = '#create-group-btn'
const createGroupInput = '#create-group-input-display-name'
const actionConfirmButton = '.oc-modal-body-actions-confirm'
const editActionBtnContextMenu = '.context-menu .oc-groups-actions-edit-trigger'
const editActionBtnQuickActions =
  '[data-item-id="%s"] .oc-table-data-cell-actions .groups-table-btn-edit'
const groupTrSelector = 'tr'
const groupNameSelector =
  '//div[@id="group-list"]//td[contains(@class,"oc-table-data-cell-displayName")]'
const groupIdSelector = `[data-item-id="%s"] .groups-table-btn-action-dropdown`
const groupCheckboxSelector = `[data-item-id="%s"]:not(.oc-table-highlighted) input[type=checkbox]`
const deleteBtnContextMenu = '.context-menu .oc-groups-actions-delete-trigger'
const deleteBtnBatchAction = '#oc-appbar-batch-actions'
const editPanel = '.sidebar-panel__body-EditPanel:visible'
const closeEditPanel = '.sidebar-panel__header .header__close'
const userInput = '#%s-input'
const compareDialogConfirm = '.compare-save-dialog-confirm-btn'

export const createGroup = async (args: { page: Page; key: string }) => {
  const { page, key } = args
  await page.locator(newGroupBtn).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'Create group modal'
  )
  await page.locator(createGroupInput).fill(key)

  const [response] = await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith('groups') && resp.status() === 201 && resp.request().method() === 'POST'
    ),
    page.locator(actionConfirmButton).click()
  ])

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['body'],
    'body after creating group'
  )

  return await response.json()
}

export const getDisplayedGroupsIds = async (args: { page: Page }): Promise<string[]> => {
  const { page } = args
  const groups = []
  const result = page.locator(groupTrSelector)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['adminSettingsWrapper'],
    'group page'
  )
  const count = await result.count()
  for (let i = 0; i < count; i++) {
    groups.push(await result.nth(i).getAttribute('data-item-id'))
  }
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['groupList'],
    'displayed groups list'
  )
  return groups
}

export const getGroupsDisplayName = async (args: { page: Page }): Promise<string> => {
  const { page } = args
  const groups = []
  const result = page.locator(groupNameSelector)

  const count = await result.count()
  for (let i = 0; i < count; i++) {
    groups.push(await result.nth(i).textContent())
  }
  return groups.join(', ')
}

export const selectGroup = async (args: { page: Page; uuid: string }): Promise<void> => {
  const { page, uuid } = args
  const checkbox = page.locator(util.format(groupCheckboxSelector, uuid))
  const checkBoxAlreadySelected = await checkbox.isChecked()

  if (checkBoxAlreadySelected) {
    return
  }
  await checkbox.click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['groupList'],
    `group list after selecting group ${uuid}`
  )
}

export const deleteGroupUsingContextMenu = async (args: {
  page: Page
  uuid: string
}): Promise<void> => {
  const { page, uuid } = args
  await page.locator(util.format(groupIdSelector, uuid)).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippbBox'],
    'group contex menu'
  )
  await page.locator(deleteBtnContextMenu).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'delete group modal'
  )

  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith(encodeURIComponent(uuid)) &&
        resp.status() === 204 &&
        resp.request().method() === 'DELETE'
    ),
    page.locator(actionConfirmButton).click()
  ])

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['body'],
    'body after deleting group'
  )
}

export const deleteGrouprUsingBatchAction = async (args: {
  page: Page
  groupIds: string[]
}): Promise<void> => {
  const { page, groupIds } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['adminSettingAppBar'],
    'admin setting app bar'
  )
  await page.locator(deleteBtnBatchAction).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'delete group modal'
  )

  const checkResponses = []
  for (const id of groupIds) {
    checkResponses.push(
      page.waitForResponse(
        (resp) =>
          resp.url().endsWith(encodeURIComponent(id)) &&
          resp.status() === 204 &&
          resp.request().method() === 'DELETE'
      )
    )
  }
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'group delete modal'
  )
  await Promise.all([...checkResponses, page.locator(actionConfirmButton).click()])

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['body'],
    'body after deleting group'
  )
}

export const changeGroup = async (args: {
  page: Page
  uuid: string
  attribute: string
  value: string
}): Promise<void> => {
  const { page, attribute, value, uuid } = args
  await page.locator(util.format(userInput, attribute)).pressSequentially(value)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['editPanel'],
    `Edit group sidebar panel while changing ${attribute}`
  )

  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith(encodeURIComponent(uuid)) &&
        resp.status() === 204 &&
        resp.request().method() === 'PATCH'
    ),
    page.locator(compareDialogConfirm).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appSidebar'],
    'group contex menu'
  )
}

export const openEditPanel = async (args: {
  page: Page
  uuid: string
  action: string
}): Promise<void> => {
  const { page, uuid, action } = args
  if (await page.locator(editPanel).count()) {
    await page.locator(closeEditPanel).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['editPanel'],
      'Edit group modal'
    )
  }
  switch (action) {
    case fileAction.contextMenu:
      await page.locator(util.format(groupIdSelector, uuid)).click()
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['tippbBox'],
        'group contex menu'
      )
      await page.locator(editActionBtnContextMenu).click()
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['editPanel'],
        'Edit group sidebar panel'
      )
      break
    case fileAction.quickAction:
      await selectUser({ page, uuid })
      await page.locator(util.format(editActionBtnQuickActions, uuid)).click()
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['editPanel'],
        'Edit group sidebar panel'
      )
      break
    default:
      throw new Error(`${action} not implemented`)
  }
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appSidebar'],
    'group contex menu'
  )
}
