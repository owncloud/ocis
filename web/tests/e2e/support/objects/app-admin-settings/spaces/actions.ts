import { Page } from '@playwright/test'
import util from 'util'
import { locatorUtils } from '../../../utils'
import { objects } from '../../../index'
import { fileAction } from '../../../../environment/constants'

const spaceTrSelector = '.spaces-table tbody > tr'
const actionConfirmButton = '.oc-modal-body-actions-confirm'
const contextMenuSelector = `[data-item-id="%s"] .spaces-table-btn-action-dropdown`
const spaceCheckboxSelector = `[data-item-id="%s"] input[type=checkbox]`
const contextMenuActionButton = `.oc-files-actions-%s-trigger`
const inputFieldSelector =
  '//div[@class="oc-modal-body-input"]//input[contains(@class,"oc-text-input")]'
const modalConfirmBtn = `.oc-modal-body-actions-confirm`
const quotaValueDropDown = `.vs__dropdown-option :text-is("%s")`
const selectedQuotaValueField = '.vs__dropdown-toggle'
const spacesQuotaSearchField = '.oc-modal .vs__search'
const appSidebarDiv = '#app-sidebar'
const toggleSidebarButton = '#files-toggle-sidebar'
const sideBarActive = '.sidebar-panel.is-active-root-panel'
const sideBarCloseButton = '.sidebar-panel .header__close:visible'
const sideBarBackButton = '.sidebar-panel .header__back:visible'
const sideBarActionButtons = `#sidebar-panel-%s-select`
const siderBarActionPanel = `#sidebar-panel-%s`
const spaceMembersDiv = '[data-testid="space-members"]'
const spaceMemberList =
  '[data-testid="space-members-role-%s"] ul [data-testid="space-members-list"]'

export const getDisplayedSpaces = async (page: Page): Promise<string[]> => {
  const spaces = []
  const result = page.locator(spaceTrSelector)

  const count = await result.count()
  for (let i = 0; i < count; i++) {
    spaces.push(await result.nth(i).getAttribute('data-item-id'))
  }

  return spaces
}

const performSpaceAction = async (args: {
  page: Page
  action: string
  via: typeof fileAction.contextMenu | typeof fileAction.batchAction
  id?: string
}): Promise<void> => {
  const { page, action, via, id } = args

  let actionButtonSelector = '.batch-actions '

  if (id && via === fileAction.contextMenu) {
    await page.locator(util.format(contextMenuSelector, id)).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['contextMenuContainer'],
      'context menu container'
    )
    actionButtonSelector = '.context-menu '
  }

  switch (action) {
    case 'rename':
      actionButtonSelector += util.format(contextMenuActionButton, action)
      break
    case 'edit-description':
      actionButtonSelector += util.format(contextMenuActionButton, action)
      break
    case 'edit-quota':
      actionButtonSelector += util.format(contextMenuActionButton, action)
      break
    case 'delete':
      actionButtonSelector += util.format(contextMenuActionButton, action)
      break
    case 'disable':
      actionButtonSelector += util.format(contextMenuActionButton, action)
      break
    case 'restore':
      actionButtonSelector += util.format(contextMenuActionButton, action)
      break
    default:
      throw new Error(`${action} not implemented`)
  }
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['filesView'], 'space page')
  await page.locator(actionButtonSelector).click()
}

export const changeSpaceQuota = async (args: {
  page: Page
  spaceIds: string[]
  value: string
  via: typeof fileAction.contextMenu | typeof fileAction.batchAction
}): Promise<void> => {
  const { page, value, spaceIds, via } = args
  await performSpaceAction({ page, action: 'edit-quota', via, id: spaceIds[0] })

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    `change quota for space ${spaceIds[0]} modal`
  )

  const searchLocator = page.locator(spacesQuotaSearchField)
  await searchLocator.pressSequentially(value)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    `change quota for space ${spaceIds[0]} modal`
  )
  await page.locator(selectedQuotaValueField).waitFor()
  await page.locator(util.format(quotaValueDropDown, `${value} GB`)).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    `change quota for space ${spaceIds[0]} modal after selecting quota`
  )
  await confirmAction({
    page,
    method: 'PATCH',
    statusCode: 200,
    spaceIds,
    actionConfirm: true
  })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['body'],
    'body after changing space quota'
  )
}

export const disableSpace = async (args: {
  page: Page
  spaceIds: string[]
  via: typeof fileAction.contextMenu | typeof fileAction.batchAction
}): Promise<void> => {
  const { page, spaceIds, via } = args
  await performSpaceAction({ page, action: 'disable', via, id: spaceIds[0] })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    `disable space ${spaceIds[0]} modal`
  )
  await confirmAction({
    page,
    method: 'DELETE',
    statusCode: 204,
    spaceIds,
    actionConfirm: false
  })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['body'],
    'body after changing space quota'
  )
}

export const enableSpace = async (args: {
  page: Page
  spaceIds: string[]
  via: typeof fileAction.contextMenu | typeof fileAction.batchAction
}): Promise<void> => {
  const { page, spaceIds, via } = args
  await performSpaceAction({ page, action: 'restore', via, id: spaceIds[0] })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    `enable space ${spaceIds[0]} modal`
  )
  await confirmAction({
    page,
    method: 'PATCH',
    statusCode: 200,
    spaceIds,
    actionConfirm: false
  })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['body'],
    'body after enabling space'
  )
}

export const deleteSpace = async (args: {
  page: Page
  spaceIds: string[]
  via: typeof fileAction.contextMenu | typeof fileAction.batchAction
}): Promise<void> => {
  const { page, spaceIds, via } = args
  await performSpaceAction({ page, action: 'delete', via, id: spaceIds[0] })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    `delete space ${spaceIds[0]} modal`
  )
  await confirmAction({
    page,
    method: 'DELETE',
    statusCode: 204,
    spaceIds,
    actionConfirm: false
  })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['body'],
    'body after deleting space'
  )
}

export const selectSpace = async (args: { page: Page; id: string }): Promise<void> => {
  const { page, id } = args
  const checkbox = page.locator(util.format(spaceCheckboxSelector, id))
  const checkBoxAlreadySelected = await checkbox.isChecked()
  if (checkBoxAlreadySelected) {
    return
  }
  await checkbox.click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    [util.format(spaceCheckboxSelector, id)],
    `Select space checkbox for space with ID ${id}`
  )
}

export const renameSpaceUsingContextMenu = async (args: {
  page: Page
  id: string
  value: string
}): Promise<void> => {
  const { page, id, value } = args
  await performSpaceAction({ page, action: 'rename', via: fileAction.contextMenu, id })
  await page.locator(inputFieldSelector).fill(value)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'space rename modal'
  )
  await confirmAction({
    page,
    method: 'PATCH',
    statusCode: 200,
    spaceIds: [id],
    actionConfirm: true
  })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['body'],
    'body after renaming space'
  )
}

export const changeSpaceSubtitleUsingContextMenu = async (args: {
  page: Page
  id: string
  value: string
}): Promise<void> => {
  const { page, id, value } = args
  await performSpaceAction({
    page,
    action: 'edit-description',
    via: fileAction.contextMenu,
    id
  })
  await page.locator(inputFieldSelector).fill(value)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    `changing subtitle for space ${id} modal`
  )
  await confirmAction({
    page,
    method: 'PATCH',
    statusCode: 200,
    spaceIds: [id],
    actionConfirm: true
  })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['body'],
    'body after changing space subtitle'
  )
}

const confirmAction = async (args: {
  page: Page
  method: string
  statusCode: number
  spaceIds: string[]
  actionConfirm: boolean
}): Promise<void> => {
  const { page, method, statusCode, spaceIds, actionConfirm } = args
  let confirmButton = modalConfirmBtn
  if (actionConfirm) {
    confirmButton = actionConfirmButton
  }

  const checkResponses = []
  for (const id of spaceIds) {
    checkResponses.push(
      page.waitForResponse(
        (resp) =>
          resp.url().endsWith(encodeURIComponent(id)) &&
          resp.status() === statusCode &&
          resp.request().method() === method
      )
    )
  }
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['tippyBox'], 'account page')
  await page.locator(confirmButton).waitFor()
  await Promise.all([...checkResponses, page.locator(confirmButton).click()])
}

export const openSpaceAdminSidebarPanel = async (args: {
  page: Page
  id: string
}): Promise<void> => {
  const { page, id } = args
  if (await page.locator(appSidebarDiv).count()) {
    await page.locator(sideBarCloseButton).click()
  }
  await selectSpace({ page, id })
  await page.click(toggleSidebarButton)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appSidebarDiv'],
    'app sidebar after opening space admin sidebar panel'
  )
}

export const openSpaceAdminActionSidebarPanel = async (args: {
  page: Page
  action: string
}): Promise<void> => {
  const { page, action } = args
  const currentPanel = page.locator(sideBarActive)
  const backButton = currentPanel.locator(sideBarBackButton)
  if (await backButton.count()) {
    await backButton.click()
    await locatorUtils.waitForEvent(currentPanel, 'transitionend')
  }
  const panelSelector = page.locator(util.format(sideBarActionButtons, action))
  const nextPanel = page.locator(util.format(siderBarActionPanel, action))
  await panelSelector.click()
  await locatorUtils.waitForEvent(nextPanel, 'transitionend')
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appSidebarDiv'],
    `app sidebar after opening space admin ${action} sidebar panel`
  )
}

export const listSpaceMembers = async (args: {
  page: Page
  filter: string
}): Promise<Array<string>> => {
  const { page, filter } = args
  await page.locator(spaceMembersDiv).waitFor()
  let users: string[] = []
  const names = []
  switch (filter) {
    case 'Can manage':
      users = await page.locator(util.format(spaceMemberList, filter)).allTextContents()
      break
    case 'Can view':
      users = await page.locator(util.format(spaceMemberList, filter)).allTextContents()
      break
    case 'Can edit with versions and trash bin':
      users = await page.locator(util.format(spaceMemberList, filter)).allTextContents()
      break
  }

  for (const user of users) {
    // the value comes in "['initials firstName secondName lastName',..]" format so only get the first name
    const [, name] = user.split(' ')
    names.push(name)
  }
  return names
}
