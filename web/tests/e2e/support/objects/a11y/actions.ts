import { Page, test } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'
import { AxeResults } from 'axe-core'
import { config } from '../../../config'

export const selectors = {
  files: '#files',
  resourceTableEditName: '.resource-table-edit-name',
  resourceIconWrapper: '.oc-resource-icon-wrapper',
  resourceTableCondensedIcon: '.resource-table-condensed',
  filesSpaceTableCondensed: '#files-space-table.condensed', // '.condensed.files-table',
  resourceTiles: '.resource-tiles',
  tilesView: '#tiles-view',
  cardMediaTop: '.oc-card-media-top',
  resourceTableIcon: '.resource-table',
  filesSpaceTable: '#files-space-table',
  filesViewOptionsBtn: '#files-view-options-btn',
  displayOptionsMenu: '#files-app-bar-controls-right .tippy-content',
  webContentMain: '#web-content-main',
  filesContextMenu: '#context-menu-drop-whitespace',
  newFileMenuBtn: '#new-file-menu-btn',
  newResourceContextMenu: '.files-app-bar-actions .tippy-content',
  newFolderBtn: '#new-folder-btn',
  ocModal: '.oc-modal',
  ocModalCancel: '.oc-modal-body-actions-cancel',
  uploadMenuBtn: '#upload-menu-btn',
  uploadContextMenu: '#upload-menu-drop',
  appbarBatchActions: '#oc-appbar-batch-actions',
  filesSpaceTableCheckbox: '#files-space-table .oc-checkbox',
  uploadMenuDropdown: '#upload-menu-drop',
  appSidebar: '#app-sidebar',
  sidebarPanelActions: '#sidebar-panel-actions',
  accountInfoContainer: '#account-info-container',
  account: '#account',
  removeUserModal: '.oc-modal.oc-modal-danger',
  appSwitcherDropdown: '#app-switcher-dropdown',
  tippyBox: '.tippy-box',
  textEditor: '#text-editor-component',
  topBar: '#oc-topbar',
  displayOptions: '#files-app-bar-controls-right',
  ocCard: '.oc-card',
  body: 'body',
  appStore: '#app-store',
  appDetails: '.app-details',
  filesView: '#files-view',
  sidebarNavigationMenu: '.oc-sidebar-nav',
  sidebarPaneSharing: '#sidebar-panel-sharing',
  filesAppBarActions: '.files-app-bar-actions',
  // visible tippy/popover (active)
  tippyBoxVisible: '.tippy-box[data-state="visible"]',
  logoWrapper: '.logo-wrapper',
  fileAppBar: '#files-app-bar',
  pageNotFound: '.page-not-found',
  adminSettingsWrapper: '#admin-settings-wrapper',
  adminSettingAppBar: '#admin-settings-app-bar',
  createGroupInput: '#create-group-input-display-name',
  actionConfirmButton: '.oc-modal-body-actions-confirm',
  contextMenuContainer: '#oc-files-context-menu',
  groupList: '.group-list',
  editPanel: '.sidebar-panel__body-EditPanel:visible',
  breadcrumb: '#files-breadcrumb',
  previewControlBar: '.preview-controls-action-bar',
  uploadInfoSnackBar: '#upload-info-snackbar',
  folderViewModal: '.folder-view-modal',
  ocNotificationSuccessMessage: '.oc-notification-message-success',
  ocModalGenerateNewToken: '.oc-modal.oc-modal-passive',
  scienceMesh: '.sciencemesh',
  filesTable: '.files-table',
  notificationContainer: 'div.oc-notification',
  spaceInfoContainer: '.space-header-infos',
  saveTextEditorOrViewerButton: '#app-save-action',
  spaceDescriptionPreview: '#space-description-preview',
  quotaValueDropDown: 'ul.vs__dropdown-menu',
  compareSaveDialog: '.compare-save-dialog',
  batchQuotaInputDropdownToggle: '.quota-select-batch-action-form .vs__dropdown-toggle', //dropdown with icon
  roleDropdownMenu: 'ul.vs__dropdown-menu',
  groupsDropdownMenu: 'ul.vs__dropdown-menu',
  usersTable: '.users-table',
  usersList: '#user-list', //users list with filter options included
  addUserToGroupForm: '#user-group-select-form',
  loginErrorMessageLocator: '#oc-login-error-message'
}

const a11yRuleTags = ['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa', 'best-practice']
// decide which tags should be included in the default configuration of axebuilder

export const analyzeAccessibilityConformityViolations = async (args: {
  page: Page
  include: string
}): Promise<AxeResults['violations']> => {
  if (config.skipA11y) {
    return []
  }

  const { page, include } = args

  const a11yResult = await new AxeBuilder({ page })
    .withTags(a11yRuleTags)
    .include(include)
    .analyze()

  if (config.testType === 'playwright') {
    test.info().attach('accessibility-scan', {
      body: JSON.stringify(a11yResult, null, 2),
      contentType: 'application/json'
    })
  }

  return a11yResult.violations
}

export const analyzeAccessibilityConformityViolationsWithExclusions = async (args: {
  page: Page
  include: string
  exclude: string | string[]
}): Promise<AxeResults['violations']> => {
  if (config.skipA11y) {
    return []
  }

  const { page, include, exclude } = args

  const axeBuilder = new AxeBuilder({ page }).withTags(a11yRuleTags).include(include)

  if (typeof exclude == 'string') {
    // excluding single selector
    axeBuilder.exclude(exclude)
  } else {
    // excluding multiple selectors
    for (const e in exclude) {
      axeBuilder.exclude(exclude[e])
    }
  }

  const a11yResult = await axeBuilder.analyze()

  return a11yResult.violations
}

export const switchToCondensedTableView = async (args: { page: Page }): Promise<void> => {
  const { page } = args

  await page.locator(selectors.files).waitFor()
  await page.locator(selectors.resourceTableCondensedIcon).click()
  await page.locator(selectors.filesSpaceTableCondensed).waitFor()
}

export const switchToDefaultTableView = async (args: { page: Page }): Promise<void> => {
  const { page } = args

  await page.locator(selectors.files).waitFor()
  await page.locator(selectors.resourceTableIcon).click()
  await page.locator(selectors.filesSpaceTable).waitFor()
}

export const showDisplayOptions = async (args: { page: Page }): Promise<void> => {
  const { page } = args

  await page.locator(selectors.files).waitFor()
  await page.locator(selectors.filesViewOptionsBtn).click()
  await page.locator(selectors.displayOptionsMenu).last().waitFor() // first element contains the invisible state, last the visible state
}

export const closeDisplayOptions = async (args: { page: Page }): Promise<void> => {
  const { page } = args

  await page.locator(selectors.filesViewOptionsBtn).click()
}

export const openFilesContextMenu = async (args: { page: Page }): Promise<void> => {
  const { page } = args

  await page.locator(selectors.files).waitFor()
  // right click to get context menu with "new folder" and "details" context menu
  await page.locator(selectors.webContentMain).click({ button: 'right' })
  await page.locator(selectors.filesContextMenu).waitFor()
}

export const exitContextMenu = async (args: { page: Page }): Promise<void> => {
  const { page } = args

  await page.locator(selectors.files).waitFor()
  await page.locator(selectors.webContentMain).click()
}

export const selectNew = async (args: { page: Page }): Promise<void> => {
  const { page } = args

  await page.locator(selectors.files).waitFor()
  await page.locator(selectors.newFileMenuBtn).click()
  await page.locator(selectors.newResourceContextMenu).waitFor()
}

export const selectFolderOptionWithinNew = async (args: { page: Page }): Promise<void> => {
  const { page } = args

  await page.locator(selectors.newResourceContextMenu).waitFor()
  await page.locator(selectors.newFolderBtn).click()
  await page.locator(selectors.ocModal).waitFor()
}

export const cancelCreatingNewFolder = async (args: { page: Page }): Promise<void> => {
  const { page } = args

  await page.locator(selectors.ocModalCancel).click()
}

export const selectUpload = async (args: { page: Page }): Promise<void> => {
  const { page } = args

  await page.locator(selectors.files).waitFor()
  await page.locator(selectors.uploadMenuBtn).click()
  await page.locator(selectors.uploadContextMenu).waitFor()
}

export const checkFileCheckbox = async (args: { page: Page }): Promise<void> => {
  const { page } = args

  await page.locator(selectors.files).waitFor()
  // check checkbox of the first file in the list
  await page.locator(selectors.filesSpaceTableCheckbox).first().check()
  await page.locator(selectors.appbarBatchActions).waitFor()
}

export const uncheckFileCheckbox = async (args: { page: Page }): Promise<void> => {
  const { page } = args

  await page.locator(selectors.files).waitFor()
  await page.locator(selectors.filesSpaceTableCheckbox).first().uncheck()
}
