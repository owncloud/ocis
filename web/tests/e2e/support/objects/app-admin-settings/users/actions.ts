import { Page } from '@playwright/test'
import util from 'util'
import { UsersEnvironment } from '../../../environment'
import { getWorld } from '../../../../environment/world'
import { createdGroupStore } from '../../../store/group'
import { objects } from '../../../index'
import { fileAction } from '../../../../environment/constants'

const userIdSelector = `[data-item-id="%s"] .users-table-btn-action-dropdown`
const editActionBtnContextMenu = '.context-menu .oc-users-actions-edit-trigger'
const editActionBtnQuickActions =
  '[data-item-id="%s"] .oc-table-data-cell-actions .users-table-btn-edit'
const editPanel = '.sidebar-panel__body-EditPanel:visible'
const closeEditPanel = '.sidebar-panel__header .header__close'
const deleteActionBtn = '.oc-users-actions-delete-trigger'
const loginDropDown = '.vs__dropdown-menu'
const dropdownOption = '.vs__dropdown-option'
const loginInputDropdownToggle = '.vs__dropdown-toggle:has(input[id="login-input"])' // login input dropdown toggle selector with dropdown icon
const compareDialogConfirmButton = '.compare-save-dialog-confirm-btn'
const addToGroupsBatchAction = '.oc-users-actions-add-to-groups-trigger'
const removeFromGroupsBatchAction = '.oc-users-actions-remove-from-groups-trigger'
const groupsModalInput = '.oc-modal .vs__search'
const actionConfirmButton = '.oc-modal-body-actions-confirm'
const userTrSelector = 'tr'
const userFilter = '.item-filter-%s'
const userFilterOption = '//ul[contains(@class, "item-filter-list")]//button[@data-test-value="%s"]'
const usersTable = '.users-table'
const quotaInput = '#quota-select-form .vs__search'
const quotaValueDropDown = 'ul.vs__dropdown-menu'
const userCheckboxSelector = `[data-item-id="%s"] input[type=checkbox]`
const editQuotaBtn = '.oc-users-actions-edit-quota-trigger'
const quotaInputBatchAction = '.quota-select-batch-action-form .vs__search'
const userInput = '#%s-input'
const roleValueDropDown = `.vs__dropdown-menu :text-is("%s")`
const groupsInput = '#user-group-select-form .vs__search'
const createUserButton = '#create-user-btn'
const userNameInput = '#create-user-input-user-name'
const displayNameInput = '#create-user-input-display-name'
const emailInput = '#create-user-input-email'
const passwordInput = '#create-user-input-password'

export interface UserInterface {
  displayName: string
  givenName: string
  id: string
  mail: string
  onPremisesSamAccountName: string
  surname: string
  userType: string
}

export const createUser = async (args: {
  page: Page
  name: string
  displayname: string
  email: string
  password: string
}): Promise<UserInterface> => {
  const { page, name, displayname, email, password } = args
  await page.locator(createUserButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'create user modal'
  )
  await page.locator(userNameInput).fill(name)
  await page.locator(displayNameInput).fill(displayname)
  await page.locator(emailInput).fill(email)
  await page.locator(passwordInput).fill(password)

  const [response] = await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith('users') && resp.status() === 201 && resp.request().method() === 'POST'
    ),
    page.locator(actionConfirmButton).click()
  ])

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['usersTable'],
    'users table after creating a new user'
  )

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['notificationContainer'],
    'notification after creating a new user'
  )

  return await response.json()
}
export const changeAccountEnabled = async (args: {
  page: Page
  uuid: string
  value: boolean
}): Promise<void> => {
  const { page, value, uuid } = args
  await page.locator(loginInputDropdownToggle).waitFor()
  await page.locator(loginInputDropdownToggle).click()
  await page.locator(loginDropDown).waitFor()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['loginToggleWithDropdown', 'loginDropDown'],
    'login dropdown to change account enabled status'
  )

  await page
    .locator(dropdownOption)
    .getByText(value === false ? 'Forbidden' : 'Allowed')
    .click()

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['loginToggleWithDropdown'],
    'login toggle after changing account enabled status'
  )

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['compareSaveDialog'],
    'compare save dialog after changing account enabled status'
  )

  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith(encodeURIComponent(uuid)) &&
        resp.status() === 200 &&
        resp.request().method() === 'PATCH'
    ),
    page.locator(compareDialogConfirmButton).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['compareSaveDialog'],
    'compare save dialog after changing account enabled status'
  )
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    [`[data-item-id="${uuid}"]`],
    `user table row for UUID ${uuid} after changing account enabled status`
  )
}

export const changeQuota = async (args: {
  page: Page
  uuid: string
  value: string
}): Promise<void> => {
  const { page, value, uuid } = args
  await page.locator(quotaInput).pressSequentially(value)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['batchQuotaInputDropdownToggle', 'quotaValueDropDown'],
    'quota input field and value dropdown during quota value entry'
  )
  await page.locator(quotaValueDropDown).first().click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['batchQuotaInputDropdownToggle'],
    'quota input field after changing quota value'
  )
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['compareSaveDialog'],
    'compare save dialog after changing quota value'
  )
  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith(encodeURIComponent(uuid)) &&
        resp.status() === 200 &&
        resp.request().method() === 'PATCH'
    ),
    page.locator(compareDialogConfirmButton).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['compareSaveDialog'],
    'compare save dialog after changing quota value'
  )
}

export const changeQuotaUsingBatchAction = async (args: {
  page: Page
  value: string
  userIds: string[]
}): Promise<void> => {
  const { page, value, userIds } = args
  await page.locator(editQuotaBtn).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'quota edit modal'
  )
  await page.locator(quotaInputBatchAction).pressSequentially(value, { delay: 100 })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['batchQuotaInputDropdownToggle', 'quotaValueDropDown'],
    'quota input field and value dropdown during quota value entry in batch action'
  )
  await page.locator(quotaValueDropDown).first().click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['batchQuotaInputDropdownToggle'],
    'quota input field after changing quota value in batch action'
  )
  await page.locator(quotaInputBatchAction).press('Enter')

  const checkResponses = []
  for (const id of userIds) {
    checkResponses.push(
      page.waitForResponse(
        (resp) =>
          resp.url().endsWith(encodeURIComponent(id)) &&
          resp.status() === 200 &&
          resp.request().method() === 'PATCH'
      )
    )
  }

  await Promise.all([...checkResponses, page.locator(actionConfirmButton).click()])

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['notificationContainer'],
    'notification after changing quota value in batch action'
  )
}

export const getDisplayedUsers = async (args: { page: Page }): Promise<string[]> => {
  const { page } = args
  const users = []
  await page.locator(usersTable).waitFor()
  const result = page.locator(userTrSelector)

  const count = await result.count()
  for (let i = 0; i < count; i++) {
    users.push(await result.nth(i).getAttribute('data-item-id'))
  }

  return users
}

export const selectUser = async (args: { page: Page; uuid: string }): Promise<void> => {
  const { page, uuid } = args
  const checkbox = page.locator(util.format(userCheckboxSelector, uuid))
  const checkBoxAlreadySelected = await checkbox.isChecked()
  if (checkBoxAlreadySelected) {
    return
  }
  await checkbox.click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    [util.format(userCheckboxSelector, uuid)],
    `Select user checkbox for user with UUID ${uuid}`
  )
}

export const addSelectedUsersToGroups = async (args: {
  page: Page
  userIds: string[]
  groups: string[]
}): Promise<void> => {
  const { page, userIds, groups } = args
  const usersEnvironment = new UsersEnvironment()
  const groupIds = []

  await page.locator(addToGroupsBatchAction).click()

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'add users to groups modal'
  )

  for (const group of groups) {
    const groupObj = usersEnvironment.getCreatedGroup({ key: group })
    groupIds.push(groupObj.uuid)
    await page.locator(groupsModalInput).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['groupsDropdownMenu'],
      'groups dropdown in add users to groups modal'
    )
    await page.locator(groupsModalInput).pressSequentially(groupObj.displayName)
    await page.keyboard.press('Enter')
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['groupsModalInput'],
      'groups input in add users to groups modal after selecting a group' + group
    )
  }

  const checkResponses = []
  for (const userId of userIds) {
    for (const groupId of groupIds) {
      checkResponses.push(
        page.waitForResponse((resp) => {
          if (
            resp.url().endsWith(`groups/${encodeURIComponent(groupId)}/members/$ref`) &&
            resp.status() === 204 &&
            resp.request().method() === 'POST'
          ) {
            return JSON.parse(resp.request().postData())['@odata.id'].endsWith(
              `/users/${encodeURIComponent(userId)}`
            )
          }
          return false
        })
      )
    }
  }

  await Promise.all([...checkResponses, page.locator(actionConfirmButton).click()])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['notificationContainer'],
    'notification after adding users to groups'
  )
}

export const removeSelectedUsersFromGroups = async (args: {
  page: Page
  userIds: string[]
  groups: string[]
}): Promise<void> => {
  const { page, userIds, groups } = args
  const usersEnvironment = new UsersEnvironment()
  const groupIds = []

  await page.locator(removeFromGroupsBatchAction).click()

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'remove users from groups modal'
  )

  for (const group of groups) {
    const groupObj = usersEnvironment.getCreatedGroup({ key: group })
    groupIds.push(groupObj.uuid)
    await page.locator(groupsModalInput).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['groupsDropdownMenu'],
      'groups dropdown in remove users from groups modal'
    )
    await page.locator(groupsModalInput).fill(groupObj.displayName)
    await page.keyboard.press('Enter')
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['groupsModalInput'],
      'groups input in remove users from groups modal after selecting a group' + group
    )
  }

  const checkResponses = []
  for (const userId of userIds) {
    for (const groupId of groupIds) {
      checkResponses.push(
        page.waitForResponse(
          (resp) =>
            resp
              .url()
              .endsWith(
                `groups/${encodeURIComponent(groupId)}/members/${encodeURIComponent(userId)}/$ref`
              ) &&
            resp.status() === 204 &&
            resp.request().method() === 'DELETE'
        )
      )
    }
  }

  await Promise.all([...checkResponses, page.locator(actionConfirmButton).click()])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['notificationContainer'],
    'notification after removing users from groups'
  )
}

export const filterUsers = async (args: {
  page: Page
  filter: string
  values: string[]
}): Promise<void> => {
  const { page, filter, values } = args
  const world = getWorld()

  await page.locator(util.format(userFilter, filter)).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBoxVisible'],
    `user filter dropdown for filter ${filter}`
  )

  for (const value of values) {
    // For group filters, convert test value to actual displayName
    let filterValue = value
    if (filter === 'groups' && world) {
      // Get all created groups and find matching one by displayName prefix
      const allGroups = Array.from(createdGroupStore.values())
      const matchingGroup = allGroups.find((g) => g.displayName?.startsWith(value))
      if (matchingGroup) {
        filterValue = matchingGroup.displayName
      }
    }

    await Promise.all([
      page.waitForResponse(
        (resp) =>
          resp.url().includes('/users') &&
          resp.status() === 200 &&
          resp.request().method() === 'GET'
      ),
      page.locator(util.format(userFilterOption, filterValue)).click()
    ])
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['usersList'],
      `users list after applying user filter ${filter} with value ${filterValue}`
    )
  }
}

export const changeUser = async (args: {
  page: Page
  uuid: string
  attribute: string
  value: string
}): Promise<void> => {
  const { page, attribute, value, uuid } = args
  await page.locator(util.format(userInput, attribute)).fill(value)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    [util.format(userInput, attribute)],
    `user input for attribute ${attribute} while changing user with UUID ${uuid}`
  )
  await page.locator(util.format(userInput, attribute)).press('Enter')

  if (attribute === 'role') {
    await page.locator(util.format(roleValueDropDown, value)).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['roleDropdownMenu'],
      `role value dropdown while changing user role to ${value} for user with UUID ${uuid}`
    )
    await Promise.all([
      page.waitForResponse(
        (resp) =>
          resp.url().endsWith(`${encodeURIComponent(uuid)}/appRoleAssignments`) &&
          resp.status() === 201 &&
          resp.request().method() === 'POST'
      ),
      page.locator(compareDialogConfirmButton).click()
    ])
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['compareSaveDialog'],
      `compare save dialog after changing user role to ${value} for user with UUID ${uuid}`
    )
  } else {
    await Promise.all([
      page.waitForResponse(
        (resp) =>
          resp.url().endsWith(encodeURIComponent(uuid)) &&
          resp.status() === 200 &&
          resp.request().method() === 'PATCH'
      ),
      page.locator(compareDialogConfirmButton).click()
    ])
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['compareSaveDialog'],
      `compare save dialog after saving changes after changing user with UUID ${uuid}`
    )
  }
}

export const addUserToGroups = async (args: {
  page: Page
  userId: string
  groups: string[]
}): Promise<void> => {
  const { page, userId, groups } = args
  const usersEnvironment = new UsersEnvironment()
  const groupIds = []
  for (const group of groups) {
    const groupObj = usersEnvironment.getCreatedGroup({ key: group })
    groupIds.push(groupObj.uuid)
    await page.locator(groupsInput).pressSequentially(groupObj.displayName)
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['addUserToGroupForm', 'groupsDropdownMenu'],
      'add user to group form after selecting a group ' + group
    )
    await page.keyboard.press('Enter')
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['addUserToGroupForm'],
      'add user to group form after selecting a group ' + group
    )
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['compareSaveDialog'],
      'compare save dialog after adding user to groups'
    )
  }

  const checkResponses = []
  for (const groupId of groupIds) {
    checkResponses.push(
      page.waitForResponse((resp) => {
        if (
          resp.url().endsWith(`groups/${encodeURIComponent(groupId)}/members/$ref`) &&
          resp.status() === 204 &&
          resp.request().method() === 'POST'
        ) {
          return JSON.parse(resp.request().postData())['@odata.id'].endsWith(
            `/users/${encodeURIComponent(userId)}`
          )
        }
        return false
      })
    )
  }

  await Promise.all([...checkResponses, page.locator(compareDialogConfirmButton).click()])

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['compareSaveDialog'],
    'compare save dialog after saving changes after adding user to groups'
  )
}

export const removeUserFromGroups = async (args: {
  page: Page
  userId: string
  groups: string[]
}): Promise<void> => {
  const { page, userId, groups } = args
  const usersEnvironment = new UsersEnvironment()
  const groupIds = []
  for (const group of groups) {
    const groupObj = usersEnvironment.getCreatedGroup({ key: group })
    groupIds.push(groupObj.uuid)
    await page.getByTitle(groupObj.displayName).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['addUserToGroupForm'],
      'add user to group form after selecting a group ' + group
    )
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['compareSaveDialog'],
      'compare save dialog after removing user from groups'
    )
  }

  const checkResponses = []
  for (const groupId of groupIds) {
    checkResponses.push(
      page.waitForResponse(
        (resp) =>
          resp
            .url()
            .endsWith(
              `groups/${encodeURIComponent(groupId)}/members/${encodeURIComponent(userId)}/$ref`
            ) &&
          resp.status() === 204 &&
          resp.request().method() === 'DELETE'
      )
    )
  }

  await Promise.all([...checkResponses, page.locator(compareDialogConfirmButton).click()])

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['compareSaveDialog'],
    'compare save dialog after saving changes after removing user from groups'
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
  }
  switch (action) {
    case fileAction.contextMenu:
      await page.locator(util.format(userIdSelector, uuid)).click()
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['tippyBoxVisible'],
        'user context menu'
      )
      await page.locator(editActionBtnContextMenu).click()
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['editPanel'],
        'user edit modal'
      )
      break
    case fileAction.quickAction:
      await selectUser({ page, uuid })
      await page.locator(util.format(editActionBtnQuickActions, uuid)).click()
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['editPanel'],
        'user edit modal'
      )
      break
    default:
      throw new Error(`${action} not implemented`)
  }
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appSidebar'],
    'account page'
  )
}

export const deleteUserUsingContextMenu = async (args: {
  page: Page
  uuid: string
}): Promise<void> => {
  const { page, uuid } = args
  await page.locator(util.format(userIdSelector, uuid)).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    [`util.format(userIdSelector, uuid)`],
    'selected user row for deletion'
  )
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBoxVisible'],
    'user context menu'
  )
  await page.locator(`.context-menu`).locator(deleteActionBtn).click()

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'delete user confirmation modal'
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
    ['usersTable'],
    'users table after deleting user'
  )

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['notificationContainer'],
    'notification after deleting user'
  )
}

export const deleteUserUsingBatchAction = async (args: {
  page: Page
  userIds: string[]
}): Promise<void> => {
  const { page, userIds } = args
  await page.locator(deleteActionBtn).click()

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'delete users confirmation modal'
  )

  const checkResponses = []
  for (const id of userIds) {
    checkResponses.push(
      page.waitForResponse(
        (resp) =>
          resp.url().endsWith(encodeURIComponent(id)) &&
          resp.status() === 204 &&
          resp.request().method() === 'DELETE'
      )
    )
  }

  await Promise.all([...checkResponses, page.locator(actionConfirmButton).click()])

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['usersTable'],
    'users table after deleting users'
  )

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['notificationContainer'],
    'notification after deleting users'
  )
}

export const waitForEditPanelToBeVisible = async (args: { page: Page }): Promise<void> => {
  const { page } = args
  await page.locator(editPanel).waitFor()
}
