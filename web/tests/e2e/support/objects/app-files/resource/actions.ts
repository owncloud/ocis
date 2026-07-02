import { Download, Locator, Page, Response, expect } from '@playwright/test'
import util from 'util'
import path from 'path'
import { resourceExists, waitForResources } from './utils'
import {
  fillCollaboraDocumentContent,
  fillOnlyOfficeDocumentContent,
  canEditCollaboraDocument,
  canEditOnlyOfficeDocument,
  getOfficeDocumentContent,
  focusCollaboraEditor,
  focusOnlyOfficeEditor,
  waitForOnlyOfficeEditor,
  waitForCollaboraEditor
} from './webOffice'
import { editor, sidebar } from '../utils'
import { environment, objects, utils } from '../../../index'
import { config } from '../../../../config'
import { File, Space } from '../../../types'
import { substitute } from '../../../utils/substitute'
import { fileAction, application, searchScope } from '../../../../environment/constants'

const topbarFilenameSelector = '#app-top-bar-resource .oc-resource-name'
const downloadFileButtonSingleShareView = '.oc-files-actions-download-file-trigger'
const downloadFolderButtonSingleShareView = '.oc-files-actions-download-archive-trigger'
const filesView = '#files-view'
const downloadFileButtonSideBar =
  '#oc-files-actions-sidebar .oc-files-actions-download-file-trigger'
const downloadFolderButtonSideBar =
  '#oc-files-actions-sidebar .oc-files-actions-download-archive-trigger'
const downloadButtonBatchAction = '.oc-files-actions-download-archive-trigger'
const appBarContextMenu = '#oc-openfile-contextmenu-trigger'
const appBarDownloadFileButton = '#oc-openfile-contextmenu .oc-files-actions-download-file-trigger'
const deleteButtonBatchAction = '.batch-actions .oc-files-actions-delete-trigger'
const createSpaceFromResourceAction = '.oc-files-actions-create-space-from-resource-trigger'
const checkBox = `//*[@data-test-resource-name="%s"]//ancestor::tr//input`
const checkBoxForTrashbin = `//*[@data-test-resource-path="%s"]//ancestor::tr//input`
const filesSelector = '//*[@data-test-resource-name="%s"]'
export const fileRow =
  '//ancestor::*[(contains(@class, "oc-tile-card") or contains(@class, "oc-tbody-tr"))]'
export const resourceNameSelector =
  ':is(#files-files-table, .oc-tiles-item, #files-shared-with-me-accepted-section, .files-table) [data-test-resource-name="%s"]'
// following breadcrumb selectors is passed to buildXpathLiteral function as the content to be inserted might contain quotes
const breadcrumbResourceNameSelector =
  '//span[contains(@class, "oc-breadcrumb-item-text") and text()=%s]'
const breadcrumbLastResourceNameSelector = '.oc-breadcrumb-item-text-last'
const breadcrumbResourceSelector = '//*[@id="files-breadcrumb"]//span[text()=%s]//ancestor::li'
const addNewResourceButton = `#new-file-menu-btn`
const createNewFolderButton = '#new-folder-btn'
const passwordProtectedFolderButton = '.new-file-btn-psec'
const createNewTxtFileButton = '.new-file-btn-txt'
const createNewMdFileButton = '.new-file-btn-md'
const createNewOfficeDocumentFileButton = '//div[@id="new-file-menu-drop"]//span[text()="%s"]'
const createNewShortcutButton = '#new-shortcut-btn'
const shortcutResorceInput = '#create-shortcut-modal-url-input'
const saveTextFileInEditorButton = '#app-save-action:visible'
const textEditor = '#text-editor #text-editor-container'
const textEditorPlainTextInput = '#text-editor #text-editor-container .cm-content'
const textEditorMarkdownInput = '#text-editor #text-editor-container .cm-content'
const resourceNameInput = '.oc-modal input'
const passwordProtectedFolderNameInput = '#input-folder-name'
const passwordProtectedFolderPasswordInput = '#input-folder-password'
const resourceUploadButton = '#upload-menu-btn'
const fileUploadInput = '#files-file-upload-input'
const folderUploadInput = '#files-folder-upload-input'
const uploadInfoCloseButton = '//button[contains(@id, "close-upload")]'
const uploadErrorCloseButton = '.oc-notification-message-danger button[aria-label="Close"]'
const filesBatchAction = '.files-app-bar-actions .oc-files-actions-%s-trigger'
const pasteButton = '.paste-files-btn'
const breadcrumbRoot = '//nav[@id="files-breadcrumb"]//li[1]'
const fileRenameInput = '.oc-text-input'
const deleteButtonSidebar = '#oc-files-actions-sidebar .oc-files-actions-delete-trigger'
const actionConfirmationButton =
  '//button[contains(@class,"oc-modal-body-actions-confirm") and text()="%s"]'
const actionSkipButton = '.oc-modal-body-actions-cancel'
const actionSecondaryConfirmationButton = '.oc-modal-body-actions-secondary'
const versionRevertButton = '//*[@data-testid="file-versions-revert-button"]'
const sideBarActionButton =
  '//div[contains(@class, "files-side-bar")]//*[contains(@data-testid, "action-handler")]//span[text()="%s"]'
const notificationMessageDialog = '.oc-notification-message-title'
const notificationMessage = '.oc-notification-message'
const permanentDeleteButton = '.oc-files-actions-delete-permanent-trigger'
const restoreResourceButton = '.oc-files-actions-restore-trigger'
const globalSearchInput = '.oc-search-input'
const globalSearchBarFilter = '.oc-search-bar-filter'
const globalSearchDirFilterDropdown =
  '//div[@id="files-global-search"]//button[contains(@id, "oc-filter")]'
const globalSearchBarFilterAllFiles = '//*[@data-test-id="all-files"]'
const globalSearchBarFilterCurrentFolder = '//*[@data-test-id="current-folder"]'
const searchList =
  '//div[@id="files-global-search-options"]//li[contains(@class,"preview")]//span[@class="oc-resource-name"]'
const globalSearchOptions = '#files-global-search-options'
const loadingSpinner = '#files-global-search-options .loading'
const filesViewOptionButton = '#files-view-options-btn'
const hiddenFilesToggleButton = '//*[@data-testid="files-switch-hidden-files"]/button'
const flatListToggleButton = '//*[@data-testid="files-switch-flat-list"]/button'
const fileNameSelector = '.oc-resource-basename'
const previewImage = '//main[@id="preview"]//div[contains(@class,"stage_media")]//img'
const previewAudio = '//main[@id="preview"]//div[contains(@class,"stage_media")]//audio//source'
const previewVideo = '//main[@id="preview"]//div[contains(@class,"stage_media")]//video//source'
const tagTableCell =
  '//*[@data-test-resource-name="%s"]/ancestor::tr//td[contains(@class, "oc-table-data-cell-tags")]'
const tagInFilesTable = '//*[contains(@class, "oc-tag")]//span[text()="%s"]//ancestor::a'
const tagInDetailsPanel = '//*[@data-testid="tags"]/td//span[text()="%s"]'
const tagInInputForm =
  '//span[contains(@class, "tags-select-tag")]//span[text()="%s"]//ancestor::span//button[contains(@class, "vs__deselect")]'
const tagFormInput = '//*[@data-testid="tags"]//input'
const tagValidationMessageSelector = '.vs__dropdown-menu .oc-text-input-danger'
const resourcesAsTiles = '#files-view .oc-tiles'
const fileVersionSidebar = '#oc-file-versions-sidebar'
const versionsPanelSelect = '//*[@data-testid="sidebar-panel-versions-select"]'
const noLinkMessage = '#web .oc-link-resolve-error-message'
const listItemPageSelector = '//*[contains(@class,"oc-pagination-list-item-page") and text()="%s"]'
const currentPageSelector =
  '//*[contains(@class,"oc-pagination-list-item-current") and text()="%s"]'
const itemsPerPageDropDownOptionSelector =
  '//li[contains(@class,"vs__dropdown-option") and text()="%s"]'
const footerTextSelector = '//*[@data-testid="files-list-footer-info"]'
const filesTableRowSelector = 'tbody tr'
const itemsPerPageDropDownSelector = '.vs__actions'
const filesPaginationNavSelector = '.files-pagination'
const uploadInfoLabel = '//*[text()="%s"]'
const pauseResumeUploadButton = '//button[contains(@id, "pause-upload")]'
const pauseUploadButton = '//button[contains(@id, "pause-upload") and @aria-label="Pause upload"]'
const resumeUploadButton = '//button[contains(@id, "pause-upload") and @aria-label="Resume upload"]'
const cancelUploadButton = '//button[contains(@id, "cancel-upload")]'
const filesContextMenuAction = 'div[id^="context-menu-drop"] button.oc-files-actions-%s-trigger'
const highlightedFileRowSelector = '#files-space-table tr.oc-table-highlighted'
const emptyTrashbinButtonSelector = '.oc-files-actions-empty-trash-bin-trigger'
const resourceLockIcon =
  '//*[@data-test-resource-name="%s"]/ancestor::tr//td//span[@data-test-indicator-type="resource-locked"]'
const sharesNavigationButtonSelector = '.oc-sidebar-nav [data-nav-name="files-shares"]'
const keepBothButton = '.oc-modal-body-actions-confirm'
const mediaNavigationButton = `//button[contains(@class, "preview-controls-%s")]`
const sideBarActions =
  '//ul[@id="oc-files-actions-sidebar"]//span[contains(@class,"oc-files-context-action-label")]/span'

const fileThumbnail = `//span[@data-test-resource-name="%s"]/ancestor::tr[contains(@class, "oc-tbody-tr")]//img[contains(@class,"oc-resource-thumbnail")]`
const fileIconWrapper = '#oc-file-details-sidebar .details-icon-wrapper'
const fileIconPreview = '#oc-file-details-sidebar .details-preview'
const activitySidebarPanel = 'sidebar-panel-activities'
const activitySidebarPanelBodyContent = '#sidebar-panel-activities .sidebar-panel__body-content'
const contextMenuAction = '//*[@id="oc-files-context-actions-context"]//span[text()="%s"]'

export const clickResource = async ({
  page,
  path
}: {
  page: Page
  path: string
}): Promise<void> => {
  const paths = path.split('/')
  for (const name of paths) {
    // if resource name consists of single or double quotes, add an escape character
    const folder = name.replace(/'/g, "\\'").replace(/"/g, '\\"')
    const resource = page.locator(util.format(resourceNameSelector, folder))
    const itemId = await resource.locator(fileRow).getAttribute('data-item-id')
    await Promise.all([
      page.waitForResponse((resp) => {
        // check with resource name or file-id
        if (
          ([207, 200].includes(resp.status()) &&
            ['PROPFIND', 'GET'].includes(resp.request().method()) &&
            resp.url().endsWith(encodeURIComponent(name))) ||
          resp.url().endsWith(itemId) ||
          resp.url().endsWith(encodeURIComponent(itemId))
        ) {
          return true
        }
        // pass for PROPFIND request wiht 207 status code
        // required for ocm shares where resource name and file-id are not determinable
        if (resp.request().method() === 'PROPFIND' && resp.status() === 207) {
          return true
        }
      }),
      resource.click()
    ])
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['breadcrumb'],
      'Personal Page Breadcrumb after navigating into a folder'
    )
  }
}

export const clickResourceFromBreadcrumb = async ({
  page,
  resource
}: {
  page: Page
  resource: string
}): Promise<void> => {
  const folder = utils.locatorUtils.buildXpathLiteral(resource)
  const itemId = await page
    .locator(util.format(breadcrumbResourceSelector, folder))
    .getAttribute('data-item-id')
  await Promise.all([
    page.waitForResponse(
      (resp) =>
        (resp.status() === 207 &&
          resp.request().method() === 'PROPFIND' &&
          resp.url().endsWith(encodeURIComponent(resource))) ||
        resp.url().endsWith(itemId) ||
        resp.url().endsWith(encodeURIComponent(itemId))
    ),
    page.locator(util.format(breadcrumbResourceNameSelector, folder)).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesSpaceTable'],
    'Personal Page child folder rows after breadcrumb navigation'
  )
  await expect(page.locator(breadcrumbLastResourceNameSelector)).toHaveText(resource)
}

/**/

export type createResourceTypes =
  | 'folder'
  | 'txtFile'
  | 'mdFile'
  | 'OpenDocument'
  | 'Microsoft Word'
  | 'Password Protected Folder'

export interface createResourceArgs {
  page: Page
  name: string
  type: createResourceTypes
  content?: string
  password?: string
}

export const createSpaceFromFolder = async ({
  page,
  folderName,
  spaceName
}: {
  page: Page
  folderName: string
  spaceName: string
}): Promise<Space> => {
  await page.locator(util.format(resourceNameSelector, folderName)).click({ button: 'right' })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'contex menu tippy box'
  )
  await page.locator(createSpaceFromResourceAction).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'create space from folder modal'
  )
  await page.locator(resourceNameInput).fill(spaceName)
  const [response] = await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.status() === 201 &&
        resp.request().method() === 'POST' &&
        resp.url().endsWith('/drives?template=default')
    ),
    page.locator(util.format(actionConfirmationButton, 'Create')).click()
  ])

  await page.locator(notificationMessage).waitFor()
  return (await response.json()) as Space
}

export const openTemplateFile = async ({
  page,
  resource,
  webOffice
}: {
  page: Page
  resource: string
  webOffice: string
}): Promise<void> => {
  await page.locator(util.format(resourceNameSelector, resource)).click({ button: 'right' })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'context menu tippy box'
  )
  await page.locator(util.format(contextMenuAction, `Open in ${webOffice}`)).click()
}

export const createFileFromTemplate = async ({
  page,
  resource,
  webOffice,
  actionType
}: {
  page: Page
  resource: string
  webOffice: string
  actionType: string
}): Promise<void> => {
  const menuItem = `Create from template via ${webOffice}`
  if (actionType.startsWith('sidebar')) {
    await sidebar.open({ page, resource })
    await sidebar.openPanel({ page, name: 'actions' })
    await page.locator(util.format(sideBarActionButton, menuItem)).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['appSidebar'],
      'sidebar panel'
    )
    return
  } else if (actionType.startsWith('context')) {
    await page.locator(util.format(resourceNameSelector, resource)).click({ button: 'right' })
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['tippyBox'],
      'context menu tippy box'
    )
    await page.locator(util.format(contextMenuAction, menuItem)).click()
    return
  }
  throw new Error(`Invalid action '${actionType}' was provided`)
}

export const createSpaceFromSelection = async ({
  page,
  resources,
  spaceName
}: {
  page: Page
  resources: string[]
  spaceName: string
}): Promise<Space> => {
  await selectOrDeselectResources({
    page,
    resources: resources.map((r) => ({name: r}) as resourceArgs), // prettier-ignore
    select: true
  })
  await page.locator(util.format(resourceNameSelector, resources[0])).click({ button: 'right' })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'contex menu tippy box'
  )
  await page.locator(createSpaceFromResourceAction).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'create space from folder modal'
  )
  await page.locator(resourceNameInput).fill(spaceName)
  const [response] = await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.status() === 201 &&
        resp.request().method() === 'POST' &&
        resp.url().endsWith('/drives?template=default')
    ),
    page.locator(util.format(actionConfirmationButton, 'Create')).click()
  ])
  await page.locator(notificationMessage).waitFor()
  return (await response.json()) as Space
}

export const createNewFolder = async ({
  page,
  resource
}: {
  page: Page
  resource: string
}): Promise<void> => {
  await page.locator(createNewFolderButton).click()

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'create new folder modal'
  )

  await page.locator(resourceNameInput).fill(resource)
  await Promise.all([
    page.waitForResponse((resp) => resp.status() === 207 && resp.request().method() === 'PROPFIND'),
    page.locator(util.format(actionConfirmationButton, 'Create')).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesSpaceTable'],
    'Personal Page new folder row'
  )
}

export const createPasswordProtectedFolder = async ({
  page,
  resource,
  password
}: {
  page: Page
  resource: string
  password: string
}): Promise<void> => {
  password = substitute(password)
  await page.locator(passwordProtectedFolderButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'create new folder modal'
  )
  await page.locator(passwordProtectedFolderNameInput).fill(resource)
  await page.locator(passwordProtectedFolderPasswordInput).fill(password)
  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith(`/${encodeURIComponent(resource)}`) &&
        resp.status() === 201 &&
        resp.request().method() === 'MKCOL'
    ),
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith('/createLink') &&
        resp.status() === 200 &&
        resp.request().method() === 'POST'
    ),
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith(`${encodeURIComponent(resource)}.psec`) &&
        resp.status() === 207 &&
        resp.request().method() === 'PROPFIND'
    ),
    page.locator(util.format(actionConfirmationButton, 'Create')).click()
  ])
}

export const createNewFileOrFolder = async (args: createResourceArgs): Promise<void> => {
  const { page, name, type, content, password } = args
  await page.locator(addNewResourceButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'create new tippy box'
  )
  switch (type) {
    case 'folder': {
      await createNewFolder({ page, resource: name })
      break
    }
    case 'Password Protected Folder': {
      await createPasswordProtectedFolder({ page, resource: name, password: password })
      break
    }
    case 'txtFile': {
      await page.locator(createNewTxtFileButton).click()
      await page.locator(resourceNameInput).clear()
      await page.locator(resourceNameInput).fill(name)
      await Promise.all([
        page.waitForResponse((resp) => resp.status() === 201 && resp.request().method() === 'PUT'),
        page.locator(util.format(actionConfirmationButton, 'Create')).click()
      ])
      await editTextDocument({ page, content, name })
      break
    }
    case 'mdFile': {
      await page.locator(createNewMdFileButton).click()
      await page.locator(resourceNameInput).clear()
      await page.locator(resourceNameInput).fill(name)
      await Promise.all([
        page.waitForResponse((resp) => resp.status() === 201 && resp.request().method() === 'PUT'),
        page.locator(util.format(actionConfirmationButton, 'Create')).click()
      ])
      await editTextDocument({ page, content, name })
      break
    }
    case 'OpenDocument': {
      // By Default when OpenDocument is created, it is opened with collabora if both app-provider services are running together
      await createDocumentFile(args, 'Collabora')
      break
    }
    case 'Microsoft Word': {
      // By Default when Microsoft Word document is created, it is opened with OnlyOffice if both app-provider services are running together
      await createDocumentFile(args, 'OnlyOffice')
      break
    }
  }
}

const createDocumentFile = async (
  args: createResourceArgs,
  editorToOpen: string
): Promise<void> => {
  const { page, name, type, content } = args
  // for creating office suites documents we need the external app provider services to be ready
  // though the service is ready it takes some time for the list of office suites documents to be visible in the dropdown in the webUI
  // which requires a retry to check if the service is ready and the office suites documents is visible in the dropdown
  const isAppProviderServiceReadyInWebUI = await isAppProviderServiceForOfficeSuitesReadyInWebUI(
    page,
    type
  )
  if (isAppProviderServiceReadyInWebUI === false) {
    throw new Error(
      `The document of type ${type} did not appear in the webUI for ${editorToOpen}. Possible reason could be the app provider service for ${editorToOpen} was not ready yet.`
    )
  }
  const a11yObject = new objects.a11y.Accessibility({ page })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'create new folder modal'
  )
  await page.locator(util.format(createNewOfficeDocumentFileButton, type)).click()
  await page.locator(resourceNameInput).clear()
  await page.locator(resourceNameInput).fill(name)
  await Promise.all([
    page.waitForLoadState(),
    page.waitForURL(/.*\/external-.*/),
    page.waitForResponse(
      (resp) =>
        resp.status() === 200 &&
        resp.request().method() === 'POST' &&
        resp.request().url().includes('/app/open?')
    ),
    page.locator(util.format(actionConfirmationButton, 'Create')).click()
  ])
  switch (editorToOpen) {
    case application.collabora:
      await fillCollaboraDocumentContent(page, content)
      break
    case application.onlyOffice:
      await fillOnlyOfficeDocumentContent(page, content)
      break
    default:
      throw new Error(
        "Editor should be either 'Collabora' or 'OnlyOffice' but found " + editorToOpen
      )
  }
  await Promise.all([
    page.waitForResponse((res) => res.status() === 207 && res.request().method() === 'PROPFIND'),
    editor.close(page)
  ])
}

export const fillDocumentContent = async ({
  page,
  text,
  editor
}: {
  page: Page
  text: string
  editor: string
}): Promise<void> => {
  switch (editor) {
    case application.textEditor:
      await page.locator(textEditorPlainTextInput).fill(text)
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['saveTextEditorOrViewerButton'],
        'Text editor Save button is enabled after content change'
      )
      break
    case application.collabora:
      await fillCollaboraDocumentContent(page, text)
      break
    case application.onlyOffice:
      await fillOnlyOfficeDocumentContent(page, text)
      break
    default:
      throw new Error("Editor should be 'TextEditor' but found " + editor)
  }
}

export const getDocumentContent = async ({
  page,
  editor
}: {
  page: Page
  editor: string
}): Promise<string> => {
  await page.waitForLoadState()
  await page.waitForURL(/.*\/external-.*/)
  switch (editor) {
    case application.collabora:
      await focusCollaboraEditor(page)
      break
    case application.onlyOffice:
      await focusOnlyOfficeEditor(page)
      break
    default:
      throw new Error("Editor should be either 'Collabora' or 'OnlyOffice' but found " + editor)
  }
  return getOfficeDocumentContent(page)
}

const isAppProviderServiceForOfficeSuitesReadyInWebUI = async (page: Page, type: string) => {
  let retry = 1
  let isCreateNewOfficeDocumentFileButtonVisible
  while (retry <= 5) {
    isCreateNewOfficeDocumentFileButtonVisible = await page
      .locator(util.format(createNewOfficeDocumentFileButton, type))
      .isVisible()
    if (isCreateNewOfficeDocumentFileButtonVisible === true) {
      break
    }
    await new Promise((resolve) => setTimeout(resolve, 3000))
    await page.reload()
    await page.locator(addNewResourceButton).click()
    retry++
  }
  return isCreateNewOfficeDocumentFileButtonVisible
}

export const createResources = async (args: createResourceArgs): Promise<void> => {
  const { page, name, type, content, password } = args
  const paths = name.split('/')
  const resource = paths.pop()

  for (const path of paths) {
    const resourcesExists = await resourceExists({
      page: page,
      name: path
    })

    if (!resourcesExists) {
      await page.locator(addNewResourceButton).click()
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['tippyBoxVisible'],
        'Tippy box after button click'
      )
      await createNewFolder({ page, resource: path })
    }
    await clickResource({ page, path })
  }
  await createNewFileOrFolder({ page, name: resource, type, content, password })
}

export const editTextDocument = async ({
  page,
  name,
  content
}: {
  page: Page
  name: string
  content: string
}): Promise<void> => {
  const isMarkdownMode = await page.locator(textEditor).getAttribute('data-markdown-mode')
  const inputLocator =
    isMarkdownMode === 'true' ? textEditorMarkdownInput : textEditorPlainTextInput

  await page.locator(inputLocator).fill(content)
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['textEditor', 'topBar'],
    'text editor'
  )
  await Promise.all([
    page.waitForResponse((resp) => resp.status() === 204 && resp.request().method() === 'PUT'),
    page.waitForResponse((resp) => resp.status() === 207 && resp.request().method() === 'PROPFIND'),
    page.locator(saveTextFileInEditorButton).click()
  ])
  await editor.close(page)
  await page.locator(util.format(resourceNameSelector, name)).waitFor()
}

/**/

export interface uploadResourceArgs {
  page: Page
  resources: File[]
  to?: string
  option?: string
  error?: string
  expectToFail?: boolean
  type?: string
}

const performUpload = async (args: uploadResourceArgs): Promise<void> => {
  const { page, resources, to, option, error, expectToFail, type } = args
  if (to) {
    await clickResource({ page, path: to })
  }

  await page.locator(resourceUploadButton).click()

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['uploadMenuDropdown'],
    'upload menu dropdown'
  )

  let uploadAction: Promise<void> = page
    .locator(type === 'folder' ? folderUploadInput : fileUploadInput)
    .setInputFiles(resources.map((file) => file.path))

  if (option) {
    await uploadAction
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['ocModal'],
      'file already exist modal'
    )
    switch (option) {
      case 'skip': {
        await page.locator(actionSkipButton).click()
        return
      }
      case 'merge':
      case 'replace': {
        uploadAction = page.locator(actionSecondaryConfirmationButton).click()
        break
      }
      case 'keep both': {
        uploadAction = page.locator(util.format(actionConfirmationButton, 'Keep both')).click()
        break
      }
    }
  }

  if (expectToFail) {
    expect(await page.locator(notificationMessageDialog).textContent()).toBe(error)
    return
  }

  await Promise.all([
    page.waitForResponse(
      (resp) =>
        [201, 204].includes(resp.status()) &&
        ['POST', 'PUT', 'PATCH'].includes(resp.request().method())
    ),
    uploadAction
  ])

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['uploadInfoSnackBar'],
    'Upload information modal while uploading'
  )
}

export const uploadLargeNumberOfResources = async (args: uploadResourceArgs): Promise<void> => {
  const { page, resources } = args
  await performUpload(args)
  await page.locator(uploadInfoCloseButton).waitFor()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['.snackbars'],
    'upload menu dropdown'
  )
  await page.locator(util.format(uploadInfoLabel, `${resources.length} items uploaded`)).waitFor()
}

export const uploadResource = async (args: uploadResourceArgs): Promise<void> => {
  const { page, resources, option } = args

  await performUpload(args)

  if (option !== 'skip') {
    await page.locator(uploadInfoCloseButton).click()
  }

  await waitForResources({
    page,
    names: resources.map((file) => path.basename(file.name))
  })
}

export const tryToUploadResource = async (args: uploadResourceArgs): Promise<void> => {
  const { page } = args
  await performUpload({ ...args, expectToFail: true })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['body'],
    'upload menu dropdown'
  )
  await page.locator(uploadErrorCloseButton).click()
}

export const dropUploadFiles = async (args: uploadResourceArgs): Promise<void> => {
  const { page, resources } = args

  // waiting to files view
  await page.locator(addNewResourceButton).waitFor()
  await utils.dragDropFiles(page, resources, filesView)

  await page.locator(util.format(uploadInfoLabel, 'Upload complete')).waitFor()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['.snackbars'],
    'upload menu dropdown'
  )
  await page.locator(uploadInfoCloseButton).click()
  await Promise.all(
    resources.map((file) =>
      page.locator(util.format(resourceNameSelector, path.basename(file.name))).waitFor()
    )
  )
}

// uploads the file without other checks
export const startResourceUpload = (args: uploadResourceArgs): Promise<void> => {
  return performUpload(args)
}

const pauseResumeUpload = async (page: Page): Promise<void> => {
  await page.locator(pauseResumeUploadButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['uploadInfoSnackBar'],
    'Upload information modal while upload is paused'
  )
}

export const navigateMediaFile = async ({
  page,
  navigationType
}: {
  page: Page
  navigationType: string
}): Promise<void> => {
  const oldFileInMediaViewer = await page
    .locator(topbarFilenameSelector)
    .getAttribute('data-test-resource-name')

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['previewControlBar'],
    'preview controls action bar'
  )
  await page.locator(util.format(mediaNavigationButton, navigationType)).click()
  const fileViewerLocator = editor.fileViewerLocator({
    page,
    fileViewerType: application.mediaViewer
  })
  await expect(fileViewerLocator).toBeVisible()

  const currentFileInMediaViewer = await page
    .locator(topbarFilenameSelector)
    .getAttribute('data-test-resource-name')
  expect(currentFileInMediaViewer).not.toEqual(oldFileInMediaViewer)
}

export const pauseResourceUpload = async (page: Page): Promise<void> => {
  await pauseResumeUpload(page)
  await page.locator(resumeUploadButton).waitFor()
}

export const resumeResourceUpload = async (page: Page): Promise<void> => {
  await pauseResumeUpload(page)
  await page.locator(pauseUploadButton).waitFor()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['uploadInfoSnackBar'],
    'Upload information modal while upload is resumed'
  )
  await page.locator(util.format(uploadInfoLabel, 'Upload complete')).waitFor()
  await page.locator(uploadInfoCloseButton).click()
}

export const cancelResourceUpload = async (page: Page): Promise<void> => {
  await page.locator(cancelUploadButton).click()
  await page.locator(util.format(uploadInfoLabel, 'Upload cancelled')).waitFor()
  await page.locator(util.format(uploadInfoLabel, '0 items uploaded')).waitFor()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['uploadInfoSnackBar'],
    'Upload information modal while upload is canceled'
  )
}

/**/

interface resourceArgs {
  name: string
  type?: string
}

export type ActionViaType =
  | typeof fileAction.sideBarPanel
  | typeof fileAction.batchAction
  | typeof fileAction.singleShareView
  | typeof fileAction.previewTopBar

export interface downloadResourcesArgs {
  page: Page
  resources: resourceArgs[]
  folder?: string
  via: ActionViaType
}

export const downloadResources = async (args: downloadResourcesArgs): Promise<Download[]> => {
  const { page, resources, folder, via } = args
  const downloads = []

  // When tracing is enabled, the HEAD request (download action) triggers a credentials prompt
  // which blocks next actions in the test.
  // See https://github.com/owncloud/web/issues/11541
  // As a workaround, we fulfill the HEAD requests with an empty response to fix the issue.
  await page.route('*/**/*.*', async (route, req) => {
    if (req.method() === 'HEAD') {
      await route.fulfill({ body: '' })
      return
    }
    await route.continue()
  })

  switch (via) {
    case fileAction.sideBarPanel: {
      if (folder) {
        await clickResource({ page, path: folder })
      }
      for (const resource of resources) {
        await sidebar.open({ page, resource: resource.name })
        await sidebar.openPanel({ page, name: 'actions' })
        const downloadResourceSelector =
          resource.type === 'file' ? downloadFileButtonSideBar : downloadFolderButtonSideBar
        await objects.a11y.Accessibility.assertNoSevereA11yViolations(
          page,
          ['sidebarPanel'],
          'sidebar panel'
        )
        const [download] = await Promise.all([
          page.waitForEvent('download'),
          page.locator(downloadResourceSelector).click()
        ])

        await sidebar.close({ page })

        downloads.push(download)
      }
      break
    }

    case fileAction.batchAction: {
      await selectOrDeselectResources({ page, resources, folder, select: true })
      if (resources.length === 1) {
        throw new Error('Single resource cannot be downloaded with batch action')
      }
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['fileAppBar'],
        'files top bar'
      )
      const [download] = await Promise.all([
        page.waitForEvent('download'),
        page.locator(downloadButtonBatchAction).click()
      ])
      downloads.push(download)
      break
    }

    case fileAction.singleShareView: {
      if (folder) {
        await clickResource({ page, path: folder })
      }
      for (const resource of resources) {
        const downloadResourceSelector =
          resource.type === 'file'
            ? downloadFileButtonSingleShareView
            : downloadFolderButtonSingleShareView
        await objects.a11y.Accessibility.assertNoSevereA11yViolations(
          page,
          ['fileAppBar'],
          'files top bar'
        )
        const [download] = await Promise.all([
          page.waitForEvent('download'),
          page.locator(downloadResourceSelector).click()
        ])

        downloads.push(download)
      }
      break
    }

    case fileAction.previewTopBar:
      await Promise.all([
        page.locator(appBarDownloadFileButton).waitFor(),
        page.locator(appBarContextMenu).click()
      ])

      const [download] = await Promise.all([
        page.waitForEvent('download'),
        page.locator(appBarDownloadFileButton).click()
      ])
      downloads.push(download)
      break
  }

  return downloads
}

export type selectResourcesArgs = {
  page: Page
  resources: resourceArgs[]
  folder?: string
  select: boolean
}

export const selectOrDeselectResources = async (args: selectResourcesArgs): Promise<void> => {
  const { page, folder, resources, select } = args
  if (folder) {
    await clickResource({ page, path: folder })
  }
  for (const resource of resources) {
    const resourceCheckbox = page.locator(util.format(checkBox, resource.name))
    if (!(await resourceCheckbox.isChecked()) && select) {
      await resourceCheckbox.check()
    } else if (await resourceCheckbox.isChecked()) {
      await resourceCheckbox.uncheck()
    }
  }
}

/**/

export interface moveOrCopyResourceArgs {
  page: Page
  resource: string
  newLocation: string
  action: 'copy' | 'move'
  method: string
  option?: string
}

export interface moveOrCopyMultipleResourceArgs extends Omit<moveOrCopyResourceArgs, 'resource'> {
  resources: string[]
}

export const pasteResource = async (args: moveOrCopyResourceArgs): Promise<void> => {
  const { page, resource, newLocation, action, method, option } = args

  await page.locator(breadcrumbRoot).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['breadcrumb', 'filesSpaceTable'],
    'Personal Page and breadcrumb before pasting resource'
  )
  const newLocationPath = newLocation.split('/')

  for (const path of newLocationPath) {
    if (path !== 'Personal') {
      await clickResource({ page, path: path })
    }
  }

  if (method === fileAction.dropDownMenu) {
    await page.locator(filesView).click({ button: 'right' })
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['filesContextMenu'],
      'Personal page'
    )
    await page.locator(util.format(filesContextMenuAction, 'copy')).click()
  } else {
    await page.locator(pasteButton).click()
  }
  if (option) {
    await Promise.all([
      page.waitForResponse(
        (resp) =>
          resp.url().endsWith(resource) &&
          resp.ok &&
          resp.request().method() === action.toUpperCase()
      ),
      option === 'replace'
        ? page.locator(actionSecondaryConfirmationButton).click()
        : page.locator(keepBothButton).click()
    ])
  } else {
    await waitForResources({
      page,
      names: [resource]
    })
  }
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesSpaceTable', 'breadcrumb', 'filesAppBarActions'],
    `Personal Page and breadcrumb after pasting resource with ${action} action`
  )
}

const selectBatchAction = async (page: Page, action: string): Promise<void> => {
  await page.locator(util.format(filesBatchAction, action)).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesAppBarActions'],
    'files app bar'
  )
  await page.mouse.move(0, 0)
}

export const moveOrCopyMultipleResources = async (
  args: moveOrCopyMultipleResourceArgs
): Promise<void> => {
  const { page, newLocation, action, method, resources } = args

  for (const resource of resources) {
    await page.locator(util.format(checkBox, resource)).click()
  }

  const waitForMoveResponses = []
  if (method === fileAction.dragDropBreadcrumb || method === fileAction.dragDrop) {
    for (const resource of resources) {
      waitForMoveResponses.push(
        page.waitForResponse(
          (resp) =>
            resp.url().endsWith(resource) &&
            resp.status() === 201 &&
            resp.request().method() === 'MOVE'
        )
      )
    }
  }

  switch (method) {
    case fileAction.dropDownMenu: {
      // after selecting multiple resources, resources can be copied or moved by clicking on any of the selected resources
      await page.locator(highlightedFileRowSelector).first().click({ button: 'right' })
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['contextMenuContainer'],
        'Personal page'
      )
      await page.locator(util.format(filesContextMenuAction, action)).click()
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['filesAppBarActions'],
        'files app bar'
      )

      await page.locator(breadcrumbRoot).click()
      const newLocationPath = newLocation.split('/')
      for (const path of newLocationPath) {
        if (path !== 'Personal') {
          await clickResource({ page, path: path })
        }
      }
      await page.locator(filesView).click({ button: 'right' })
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['filesContextMenu'],
        'Personal page'
      )
      await page.locator(util.format(filesContextMenuAction, 'copy')).click()
      break
    }
    case fileAction.batchAction: {
      await selectBatchAction(page, action)

      await page.locator(breadcrumbRoot).click()
      const newLocationPath = newLocation.split('/')
      for (const path of newLocationPath) {
        if (path !== 'Personal') {
          await clickResource({ page, path: path })
        }
      }
      await page.locator(pasteButton).click()
      break
    }
    case fileAction.keyboard: {
      const keyValue = action === 'copy' ? 'c' : 'x'
      await page.keyboard.press(`ControlOrMeta+${keyValue}`)
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['filesAppBarActions'],
        'files app bar'
      )
      await page.locator(breadcrumbRoot).click()
      const newLocationPath = newLocation.split('/')
      for (const path of newLocationPath) {
        if (path !== 'Personal') {
          await clickResource({ page, path: path })
        }
      }
      await page.keyboard.press('ControlOrMeta+v')
      break
    }
    case fileAction.dragDrop: {
      const source = page.locator(highlightedFileRowSelector).first()
      const target = page.locator(util.format(resourceNameSelector, newLocation))

      await Promise.all([...waitForMoveResponses, source.dragTo(target)])

      await target.click()
      break
    }
    case fileAction.dragDropBreadcrumb: {
      const source = page.locator(highlightedFileRowSelector).first()
      const target = page.locator(
        util.format(
          breadcrumbResourceNameSelector,
          utils.locatorUtils.buildXpathLiteral(newLocation)
        )
      )

      await Promise.all([...waitForMoveResponses, source.dragTo(target)])

      await target.click()
      break
    }
  }
  await waitForResources({
    page,
    names: resources
  })
}

export const moveOrCopyResource = async (args: moveOrCopyResourceArgs): Promise<void> => {
  const { page, resource, newLocation, action, method, option } = args
  const { dir: resourceDir, base: resourceBase } = path.parse(resource)

  if (resourceDir) {
    await clickResource({ page, path: resourceDir })
  }

  switch (method) {
    case fileAction.dropDownMenu: {
      await page.locator(util.format(resourceNameSelector, resourceBase)).click({ button: 'right' })
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['contextMenuContainer'],
        'Personal page'
      )
      await page.locator(util.format(filesContextMenuAction, action)).click()
      await pasteResource({ page, resource: resourceBase, newLocation, action, method, option })
      break
    }
    case fileAction.batchAction: {
      await page.locator(util.format(checkBox, resourceBase)).click()
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['filesAppBarActions'],
        'files app bar'
      )
      await selectBatchAction(page, action)
      await pasteResource({ page, resource: resourceBase, newLocation, action, method, option })
      break
    }
    case fileAction.sideBarPanel: {
      await sidebar.open({ page: page, resource: resourceBase })
      await sidebar.openPanel({ page: page, name: 'actions' })

      const actionButtonType = action === 'copy' ? 'Copy' : 'Cut'
      await page.locator(util.format(sideBarActionButton, actionButtonType)).click()
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['ocNotificationSuccessMessage'],
        `Success notification after clicking ${actionButtonType} in sidebar`
      )
      await pasteResource({ page, resource: resourceBase, newLocation, action, method, option })
      break
    }
    case fileAction.keyboard: {
      const resourceCheckbox = page.locator(util.format(checkBox, resourceBase))
      await resourceCheckbox.check()
      const keyValue = action === 'copy' ? 'c' : 'x'
      await page.keyboard.press(`ControlOrMeta+${keyValue}`)
      await page.locator(breadcrumbRoot).click()
      const newLocationPath = newLocation.split('/')
      for (const path of newLocationPath) {
        if (path !== 'Personal') {
          await clickResource({ page, path: path })
        }
      }
      await Promise.all([
        page.waitForResponse(
          (resp) =>
            resp.url().endsWith(resource) &&
            resp.status() === 201 &&
            resp.request().method() === action.toUpperCase()
        ),
        page.keyboard.press('ControlOrMeta+v')
      ])
      break
    }
    case fileAction.dragDrop: {
      const source = page.locator(util.format(resourceNameSelector, resourceBase))
      const target = page.locator(util.format(resourceNameSelector, newLocation))

      await Promise.all([
        page.waitForResponse(
          (resp) =>
            resp.url().endsWith(resource) &&
            resp.status() === 201 &&
            resp.request().method() === 'MOVE'
        ),
        utils.dragTo(page, source, target)
      ])

      await Promise.all([
        page.locator(util.format(resourceNameSelector, resourceBase)),
        page.locator(util.format(resourceNameSelector, newLocation)).click()
      ])

      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['breadcrumb', 'filesSpaceTable'],
        'Personal Page and breadcrumb after dragging and dropping resource'
      )

      break
    }
    case fileAction.dragDropBreadcrumb: {
      const source = page.locator(util.format(resourceNameSelector, resourceBase))
      const target = page.locator(
        util.format(
          breadcrumbResourceNameSelector,
          utils.locatorUtils.buildXpathLiteral(newLocation)
        )
      )

      await Promise.all([
        page.waitForResponse(
          (resp) =>
            resp.url().endsWith(resource) &&
            resp.status() === 201 &&
            resp.request().method() === 'MOVE'
        ),
        source.dragTo(target)
      ])

      await Promise.all([
        page.locator(util.format(resourceNameSelector, resourceBase)),
        page
          .locator(
            util.format(
              breadcrumbResourceNameSelector,
              utils.locatorUtils.buildXpathLiteral(newLocation)
            )
          )
          .click()
      ])

      break
    }
  }
  await waitForResources({
    page,
    names: [resourceBase]
  })
}

/**/

export interface renameResourceArgs {
  page: Page
  resource: string
  newName: string
}

export interface resourceTagsArgs {
  page: Page
  resource: string
  tags: string[]
}

export const renameResource = async (args: renameResourceArgs): Promise<void> => {
  const { page, resource, newName } = args
  const { dir: resourceDir, base: resourceBase } = path.parse(resource)

  if (resourceDir) {
    await clickResource({ page, path: resourceDir })
  }

  await page.locator(util.format(resourceNameSelector, resourceBase)).click({ button: 'right' })
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'Rename resource tippy box'
  )
  await page.locator(util.format(filesContextMenuAction, 'rename')).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'Rename resource modal'
  )
  await page.locator(fileRenameInput).clear()
  await page.locator(fileRenameInput).fill(newName)
  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().endsWith(resourceBase) &&
        resp.status() === 201 &&
        resp.request().method() === 'MOVE'
    ),
    page.locator(util.format(actionConfirmationButton, 'Rename')).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesView'],
    'files view after renaming a resource'
  )
}

/**/

export interface resourceVersionArgs {
  page: Page
  files: File[]
  folder?: string
  openDetailsPanel?: boolean
}

export const restoreResourceVersion = async (args: resourceVersionArgs) => {
  const { page, files, folder, openDetailsPanel } = args
  if (openDetailsPanel) {
    const fileName = files.map((file) => path.basename(file.name))
    await clickResource({ page, path: folder })
    await sidebar.open({ page, resource: fileName[0] })
    await sidebar.openPanel({ page, name: 'versions' })
  }
  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().includes('/v/') && resp.status() === 204 && resp.request().method() === 'COPY'
    ),
    page.locator(versionRevertButton).first().click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appSidebar', 'filesView'],
    'Files view with opened sidebar after restoring to a previous version of a resource'
  )
}

/**/
export interface deleteResourceArgs {
  page: Page
  resourcesWithInfo: resourceArgs[]
  via: ActionViaType
  folder?: string
  isPublicLink?: boolean
}

export const deleteResource = async (args: deleteResourceArgs): Promise<void> => {
  const { page, resourcesWithInfo, folder, via, isPublicLink } = args
  switch (via) {
    case fileAction.sideBarPanel: {
      if (folder) {
        await clickResource({ page, path: folder })
      }
      for (const resource of resourcesWithInfo) {
        await sidebar.open({ page, resource: resource.name })
        await sidebar.openPanel({ page, name: 'actions' })
        await objects.a11y.Accessibility.assertNoSevereA11yViolations(
          page,
          ['sidebarPanelActions'],
          'Sidebar actions panel action before deleting resource'
        )
        await Promise.all([
          page.waitForResponse(
            (resp) =>
              resp.url().includes(encodeURIComponent(resource.name)) &&
              resp.status() === 204 &&
              resp.request().method() === 'DELETE'
          ),
          page.locator(deleteButtonSidebar).first().click()
        ])
        await sidebar.close({ page })
      }
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['sidebarPanelActions', 'filesView'],
        'Sidebar actions panel before deleting resource and files view after deleting resource'
      )
      break
    }

    case fileAction.batchAction: {
      await selectOrDeselectResources({ page, resources: resourcesWithInfo, folder, select: true })

      const waitResponses = []
      for (const info of resourcesWithInfo) {
        waitResponses.push(
          page.waitForResponse(
            (resp) =>
              resp.status() === 204 &&
              resp.request().method() === 'DELETE' &&
              resp
                .request()
                .url()
                .endsWith(`/${encodeURIComponent(info.name)}`)
          )
        )
      }
      if (!isPublicLink) {
        // wait for GET response after all the resource are deleted with batch action
        waitResponses.push(
          page.waitForResponse(
            (resp) =>
              resp.url().includes('graph/v1beta1/drives') &&
              resp.status() === 200 &&
              resp.request().method() === 'GET'
          )
        )
        await objects.a11y.Accessibility.assertNoSevereA11yViolations(
          page,
          ['body'],
          'public link of a resource'
        )
      }

      await Promise.all([...waitResponses, page.locator(deleteButtonBatchAction).click()])
      break
    }
  }
}

export interface downloadResourceVersionArgs {
  page: Page
  files: File[]
  folder?: string
}

export const downloadResourceVersion = async (args: downloadResourceVersionArgs) => {
  const { page, files, folder } = args
  const fileName = files.map((file) => path.basename(file.name))
  const downloads: Response[] = []
  await clickResource({ page, path: folder })
  await sidebar.open({ page, resource: fileName[0] })
  await sidebar.openPanel({ page, name: 'versions' })
  const [download] = await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().includes('/v/') && resp.status() === 200 && resp.request().method() === 'HEAD'
    ),
    page.waitForEvent('download'),
    page.locator('//*[@data-testid="file-versions-download-button"]').first().click()
  ])
  await sidebar.close({ page: page })
  downloads.push(download)
  return downloads
}

export interface deleteResourceTrashbinArgs {
  page: Page
  resource: string
}

export interface deleteTrashbinMultipleResourcesArgs extends Omit<
  deleteResourceTrashbinArgs,
  'resource'
> {
  resources: string[]
}

export const deleteResourceTrashbin = async (args: deleteResourceTrashbinArgs): Promise<string> => {
  const { page, resource } = args
  const resourceCheckbox = page.locator(
    util.format(checkBoxForTrashbin, `/${resource.replace(/^\/+/, '')}`)
  )
  await new Promise((resolve) => setTimeout(resolve, 5000))
  if (!(await resourceCheckbox.isChecked())) {
    await resourceCheckbox.check()
  }

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesView'],
    'trashbin page'
  )
  await page.locator(permanentDeleteButton).first().click()
  await Promise.all([
    page.waitForResponse((resp) => resp.status() === 204 && resp.request().method() === 'DELETE'),
    page.locator(util.format(actionConfirmationButton, 'Delete')).click()
  ])
  const message = await page.locator(notificationMessageDialog).textContent()
  return message.trim().toLowerCase()
}

export const deleteTrashbinMultipleResources = async (
  args: deleteTrashbinMultipleResourcesArgs
): Promise<void> => {
  const { page, resources } = args
  for (const resource of resources) {
    await page.locator(util.format(checkBox, resource)).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['filesView'],
      'files view when selecting multiple resources in trashbin for deletion'
    )
  }

  await page.locator(permanentDeleteButton).first().click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'confirmation modal when deleting multiple resources in trashbin'
  )
  await Promise.all([
    page.waitForResponse((resp) => resp.status() === 204 && resp.request().method() === 'DELETE'),
    page.locator(util.format(actionConfirmationButton, 'Delete')).click()
  ])

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesView'],
    'files view after deleting multiple resources in trashbin'
  )

  for (const resource of resources) {
    await expect(page.locator(util.format(filesSelector, resource))).not.toBeVisible()
  }
}

export const emptyTrashbin = async ({ page }: { page: Page }): Promise<void> => {
  await page.locator(emptyTrashbinButtonSelector).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'confirmation modal when emptying the trashbin'
  )
  await Promise.all([
    page.waitForResponse((resp) => resp.status() === 204 && resp.request().method() === 'DELETE'),
    page.locator(util.format(actionConfirmationButton, 'Delete')).click()
  ])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesView'],
    'files view after emptying the trashbin'
  )
  const message = await page.locator(notificationMessageDialog).textContent()
  expect(message).toBe('All deleted files were removed')
}

export const expectThatDeleteButtonIsNotVisible = async (
  args: deleteResourceTrashbinArgs
): Promise<void> => {
  const { page, resource } = args
  const resourceCheckbox = page.locator(
    util.format(checkBoxForTrashbin, `/${resource.replace(/^\/+/, '')}`)
  )
  if (!(await resourceCheckbox.isChecked())) {
    await resourceCheckbox.check()
  }
  const deleteButton = page.locator(permanentDeleteButton)
  await expect(deleteButton).not.toBeVisible()
}

export interface restoreResourceTrashbinArgs {
  resource: string
  page: Page
}

export interface batchRestoreTrashbinResourcesArgs {
  resources: string[]
  page: Page
}

export interface clickTagArgs {
  resource: string
  tag: string
  page: Page
}

export interface createSpaceFromFolderArgs {
  folderName: string
  spaceName: string
  page: Page
}

export interface createSpaceFromSelectionArgs {
  resources: string[]
  spaceName: string
  page: Page
}

export const selectTrashbinResource = async (page: Page, resource: string): Promise<void> => {
  const resourceCheckbox = page.locator(
    util.format(checkBoxForTrashbin, `/${resource.replace(/^\/+/, '')}`)
  )
  if (!(await resourceCheckbox.isChecked())) {
    await resourceCheckbox.check()
  }
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesView'],
    'trashbin page'
  )
}

export const restoreTrashBinResource = async (
  args: restoreResourceTrashbinArgs
): Promise<string> => {
  const { page, resource } = args
  await selectTrashbinResource(page, resource)

  const resourceNameLocator = page.locator(util.format(resourceNameSelector, resource))
  const itemId = await resourceNameLocator.locator(fileRow).getAttribute('data-item-id')
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesView'],
    'trashbin page'
  )

  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.status() === 201 &&
        resp.request().method() === 'MOVE' &&
        resp.request().url().endsWith(`/${itemId}`)
    ),
    page.locator(restoreResourceButton).click()
  ])

  const message = await page.locator(notificationMessageDialog).textContent()
  return message.trim().toLowerCase()
}

export const batchRestoreTrashBinResources = async (
  args: batchRestoreTrashbinResourcesArgs
): Promise<string> => {
  const { page, resources } = args

  const waitResponses = []
  for (const resource of resources) {
    await selectTrashbinResource(page, resource)
    const resourceNameLocator = page.locator(util.format(resourceNameSelector, resource))
    const itemId = await resourceNameLocator.locator(fileRow).getAttribute('data-item-id')
    waitResponses.push(
      page.waitForResponse(
        (resp) =>
          resp.status() === 201 &&
          resp.request().method() === 'MOVE' &&
          resp.request().url().endsWith(`/${itemId}`)
      )
    )
  }

  await Promise.all([...waitResponses, page.locator(restoreResourceButton).click()])
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesView'],
    'files view after restoring multiple resources from trashbin'
  )

  const message = await page.locator(notificationMessageDialog).textContent()
  return message.trim().toLowerCase()
}

export const expectThatRestoreResourceButtonVisibility = async (
  args: restoreResourceTrashbinArgs
): Promise<void> => {
  const { page, resource } = args
  const resourceCheckbox = page.locator(
    util.format(checkBoxForTrashbin, `/${resource.replace(/^\/+/, '')}`)
  )
  if (!(await resourceCheckbox.isChecked())) {
    await resourceCheckbox.check()
  }
  const restoreButton = page.locator(restoreResourceButton)
  await expect(restoreButton).not.toBeVisible()
}

export const getTagsForResourceVisibilityInFilesTable = async (
  args: resourceTagsArgs
): Promise<boolean> => {
  const { page, resource, tags } = args
  const { dir: resourceDir } = path.parse(resource)

  const folderPaths = resource.split('/')
  const resourceName = folderPaths.pop()

  if (resourceDir) {
    await clickResource({ page, path: resourceDir })
  }

  const tagCellSelector = util.format(tagTableCell, resourceName)
  await page.locator(tagCellSelector).waitFor()
  const resourceTagCell = page.locator(tagCellSelector)

  for (const tag of tags) {
    const tagSpan = resourceTagCell.locator(util.format(tagInFilesTable, tag))
    const isVisible = await tagSpan.isVisible()
    if (!isVisible) {
      return false
    }
  }

  return true
}

export const clickResourceTag = async (args: clickTagArgs): Promise<void> => {
  const { page, resource, tag } = args
  const { dir: resourceDir } = path.parse(resource)

  const folderPaths = resource.split('/')
  const resourceName = folderPaths.pop()

  if (resourceDir) {
    await clickResource({ page, path: resourceDir })
  }

  const tagCellSelector = util.format(tagTableCell, resourceName)
  await page.locator(tagCellSelector).waitFor()
  const resourceTagCell = page.locator(tagCellSelector)
  const tagSpan = resourceTagCell.locator(util.format(tagInFilesTable, tag))
  tagSpan.click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(page, ['filesView'], 'files view')
}

export const getTagsForResourceVisibilityInDetailsPanel = async (
  args: resourceTagsArgs
): Promise<boolean> => {
  const { page, resource, tags } = args
  const { dir: resourceDir } = path.parse(resource)

  const folderPaths = resource.split('/')
  const resourceName = folderPaths.pop()

  if (resourceDir) {
    await clickResource({ page, path: resourceDir })
  }

  await sidebar.open({ page: page, resource: resourceName })

  for (const tag of tags) {
    const tagSelector = util.format(tagInDetailsPanel, tag)
    await page.locator(tagSelector).waitFor()
    const tagSpan = page.locator(tagSelector)
    const isVisible = await tagSpan.isVisible()
    if (!isVisible) {
      return false
    }
  }

  return true
}
export type searchFilter = 'all files' | 'current folder'

export interface searchResourceGlobalSearchArgs {
  keyword: string
  filter?: searchFilter
  pressEnter?: boolean
  page: Page
}

export const searchResourceGlobalSearch = async (
  args: searchResourceGlobalSearchArgs
): Promise<void> => {
  const { page, keyword, filter, pressEnter } = args

  // .reload() waits nicely for search indexing to be finished
  await page.reload()

  // select the filter if provided
  if (filter) {
    await page.locator(globalSearchDirFilterDropdown).click()
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['tippyBox'],
      'global search filter dropdown'
    )
    await page
      .locator(
        filter === searchScope.allFiles
          ? globalSearchBarFilterAllFiles
          : globalSearchBarFilterCurrentFolder
      )
      .click()
  }

  await page.locator(globalSearchBarFilter).click()

  if (!keyword) {
    await page.locator(globalSearchInput).click()
    await page.keyboard.press('Enter')
    return
  }

  await Promise.all([
    page.waitForResponse((resp) => resp.status() === 207 && resp.request().method() === 'REPORT'),
    page.locator(globalSearchInput).fill(keyword),
    expect(page.locator(globalSearchOptions)).toBeVisible(),
    expect(page.locator(loadingSpinner)).not.toBeVisible()
  ])

  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    [globalSearchOptions],
    'global search file list'
  )

  if (pressEnter) {
    await page.keyboard.press('Enter')
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      [filesView],
      'global search file list'
    )
  }
}

export type displayedResourceType = 'search list' | 'files list' | 'Shares' | 'trashbin'

export interface getDisplayedResourcesArgs {
  keyword: displayedResourceType
  page: Page
}

export const getDisplayedResourcesFromSearch = async (page: Page): Promise<string[]> => {
  const result = await page.locator(searchList).allInnerTexts()
  // the result has values like `test\n.txt` so remove new line
  return result.map((result) => result.replace('\n', ''))
}

// The global search dropdown only queries the backend when the search term changes
// (see the debounced `watch(term)` in SearchBar.vue). A single query issued right after
// an upload can miss resources that tika is still indexing, and re-reading that one-shot
// dropdown will never surface them. Re-issue the query by clearing and re-filling the
// existing keyword: clearing forces a term change so the watcher re-runs the search even
// for the same keyword, and avoiding a page reload preserves the active location filter.
export const reSearchAndGetDisplayedResourcesFromSearch = async (page: Page): Promise<string[]> => {
  const searchInput = page.locator(globalSearchInput)
  const keyword = await searchInput.inputValue()
  if (keyword) {
    await searchInput.clear()
    await Promise.all([
      page.waitForResponse((resp) => resp.status() === 207 && resp.request().method() === 'REPORT'),
      searchInput.fill(keyword)
    ])
    await expect(page.locator(loadingSpinner)).not.toBeVisible()
  }
  return getDisplayedResourcesFromSearch(page)
}

export const getDisplayedResourcesFromFilesList = async (page: Page): Promise<string[]> => {
  // wait for tika indexing
  await new Promise((resolve) => setTimeout(resolve, 1000))
  const files = []
  const result = page.locator('[data-test-resource-path]')

  try {
    await result.first().waitFor({ timeout: config.minTimeout * 1000 })
  } catch {
    console.log('Files list is empty')
  }

  const count = await result.count()
  for (let i = 0; i < count; i++) {
    files.push(await result.nth(i).getAttribute('data-test-resource-name'))
  }

  return files
}

export const getDisplayedResourcesFromShares = async (page: Page): Promise<string[]> => {
  const files = []
  await page.locator(sharesNavigationButtonSelector).click()
  const result = page.locator('[data-test-resource-path]')

  const count = await result.count()
  for (let i = 0; i < count; i++) {
    files.push(await result.nth(i).getAttribute('data-test-resource-name'))
  }

  return files
}

export const getDisplayedResourcesFromTrashbin = async (page: Page): Promise<string[]> => {
  const files = []
  const result = page.locator('[data-test-resource-path]')
  try {
    await result.first().waitFor({ timeout: config.minTimeout * 1000 })
  } catch {
    console.log('Trashbin is empty')
  }

  const count = await result.count()
  for (let i = 0; i < count; i++) {
    files.push(await result.nth(i).getAttribute('data-test-resource-name'))
  }

  return files
}

export interface switchViewModeArgs {
  page: Page
  target: 'resource-table' | 'resource-tiles'
}

export const clickViewModeToggle = async (args: switchViewModeArgs): Promise<void> => {
  const { page, target } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['displayOptions'],
    'file controls modal'
  )
  await page.locator(`.viewmode-switch-buttons .${target}`).click()
}

export const expectThatResourcesAreTiles = async (args: { page: Page }): Promise<void> => {
  const { page } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tilesView'],
    'file tile view'
  )
  const tiles = page.locator(resourcesAsTiles)
  await expect(tiles).toBeVisible()
}

export const showHiddenResources = async (page: Page): Promise<void> => {
  await page.locator(filesViewOptionButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'files control tippy box'
  )
  await page.locator(hiddenFilesToggleButton).click()
  // close the files view option
  await page.locator(filesViewOptionButton).click()
}

export const toggleFlatList = async (page: Page): Promise<void> => {
  await page.locator(filesViewOptionButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'files control tippy box'
  )
  await page.locator(flatListToggleButton).click()
  // close the files view option
  await page.locator(filesViewOptionButton).click()
}

export const getAllFiles = async (page: Page): Promise<string[]> => {
  await page.locator(footerTextSelector).waitFor()
  const fileLocators = await page.locator(fileNameSelector).all()
  return await Promise.all(fileLocators.map((locator) => locator.textContent()))
}

export interface editResourcesArgs {
  page: Page
  name: string
  type: string
  content: string
}

export const editResource = async (args: editResourcesArgs): Promise<void> => {
  const { page, name, type, content } = args
  const { dir: resourceDir } = path.parse(name)

  const folderPaths = name.split('/')
  const resourceName = folderPaths.pop()

  if (resourceDir) {
    await clickResource({ page, path: resourceDir })
  }

  switch (type) {
    case 'OpenDocument':
      await fillDocumentContent({ page, text: content, editor: 'Collabora' })
      break
    case 'Microsoft Word':
      await fillDocumentContent({ page, text: content, editor: 'OnlyOffice' })
      break
    default:
      await page.locator(util.format(resourceNameSelector, resourceName)).click()
      await editTextDocument({ page, content: content, name: resourceName })
  }
}

const openSideBar = async ({ page, resource }): Promise<void> => {
  const { dir: resourceDir } = path.parse(resource)

  const folderPaths = resource.split('/')
  const resourceName = folderPaths.pop()

  if (resourceDir) {
    await clickResource({ page, path: resourceDir })
  }

  await sidebar.open({ page: page, resource: resourceName })
}

export const addTagsToResource = async (args: resourceTagsArgs): Promise<void> => {
  const { page, resource, tags } = args
  await openSideBar({ page, resource })
  const inputForm = page.locator(tagFormInput)

  for (const tag of tags) {
    await inputForm.pressSequentially(tag)

    await Promise.all([
      page.waitForResponse(
        (resp) =>
          resp.url().endsWith('tags') && resp.status() === 200 && resp.request().method() === 'PUT'
      ),

      page.locator('.vs__dropdown-option').first().press('Enter')
    ])
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['appSidebar'],
      'side bar after adding tag'
    )
  }

  await sidebar.close({ page })
}

export const tryToAddTagsToResource = async (args: resourceTagsArgs): Promise<void> => {
  const { page, resource, tags } = args
  await openSideBar({ page, resource })
  const inputForm = page.locator(tagFormInput)

  for (const tag of tags) {
    await inputForm.pressSequentially(tag)
    await objects.a11y.Accessibility.assertNoSevereA11yViolations(
      page,
      ['appSidebar'],
      'side bar after trying to add too long tag name'
    )
  }
}

export const removeTagsFromResource = async (args: resourceTagsArgs): Promise<void> => {
  const { page, resource, tags } = args
  const { dir: resourceDir } = path.parse(resource)

  const folderPaths = resource.split('/')
  const resourceName = folderPaths.pop()

  if (resourceDir) {
    await clickResource({ page, path: resourceDir })
  }

  await sidebar.open({ page: page, resource: resourceName })

  for (const tag of tags) {
    await page.locator(util.format(tagInInputForm, tag)).click()
  }
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['appSidebar'],
    'sidebar after removing tag'
  )
  await sidebar.close({ page })
}

export interface openFileInViewerArgs {
  page: Page
  name: string
  actionType:
    | 'mediaviewer'
    | 'audioviewer'
    | 'pdfviewer'
    | 'texteditor'
    | 'Collabora'
    | 'OnlyOffice'
}

export const openFileInViewer = async (args: openFileInViewerArgs): Promise<void> => {
  const { page, name, actionType } = args

  switch (actionType) {
    case 'OnlyOffice':
      await Promise.all([
        page.waitForResponse(
          (resp) =>
            resp.url().includes(`app_name=${actionType}`) &&
            resp.status() === 200 &&
            resp.request().method() === 'POST'
        ),
        page.locator(util.format(resourceNameSelector, name)).click()
      ])
      await waitForOnlyOfficeEditor(page)
      break
    case 'Collabora':
      await Promise.all([
        page.waitForResponse(
          (resp) =>
            resp.url().includes(`app_name=${actionType}`) &&
            resp.status() === 200 &&
            resp.request().method() === 'POST'
        ),
        page.locator(util.format(resourceNameSelector, name)).click()
      ])
      await waitForCollaboraEditor(page)
      break
    case 'mediaviewer': {
      await Promise.all([
        page.waitForResponse(
          (resp) =>
            resp.url().includes('preview') &&
            resp.status() === 200 &&
            resp.request().method() === 'GET'
        ),
        page.locator(util.format(resourceNameSelector, name)).click()
      ])
      const extension = name.split('.').pop()
      switch (extension) {
        case 'mp3':
        case 'ogg':
          expect(await page.locator(previewAudio).getAttribute('src')).toContain(name)
          break
        case 'webm':
        case 'mp4':
          expect(await page.locator(previewVideo).getAttribute('src')).toContain(name)
          break
        default:
          // in case of error <img> doesn't contain src="blob:https://url"
          expect(await page.locator(previewImage).getAttribute('src')).toContain('blob')
      }
      break
    }
    case 'pdfviewer':
    case 'texteditor': {
      await Promise.all([
        page.waitForResponse(
          (resp) => resp.status() === 207 && resp.request().method() === 'PROPFIND'
        ),
        page.locator(util.format(resourceNameSelector, name)).click()
      ])
      break
    }
  }
}

export const previewMediaFromSidebarPanel = async ({
  page,
  resource
}: {
  page: Page
  resource: string
}): Promise<void> => {
  await sidebar.open({ page, resource })
  await sidebar.openPanel({ page, name: 'actions' })
  await page.locator(util.format(sideBarActionButton, 'Preview')).first().click()
}

export const checkThatFileVersionIsNotAvailable = async (
  args: resourceVersionArgs
): Promise<void> => {
  const { page, files, folder } = args
  const fileName = files.map((file) => path.basename(file.name))
  await clickResource({ page, path: folder })

  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.url().includes('dav/meta') &&
        resp.status() === 403 &&
        resp.request().method() === 'PROPFIND'
    ),
    sidebar.open({ page, resource: fileName[0] })
  ])

  await sidebar.openPanel({ page, name: 'versions' })
  await expect(page.locator(fileVersionSidebar)).toHaveText('No versions available for this file')
}

export const checkThatFileVersionPanelIsNotAvailable = async (
  args: resourceVersionArgs
): Promise<void> => {
  const { page, files, folder } = args
  const fileName = files.map((file) => path.basename(file.name))
  await clickResource({ page, path: folder })
  await sidebar.open({ page, resource: fileName[0] })

  await expect(page.locator(versionsPanelSelect)).not.toBeVisible()
}

export const expectThatPublicLinkIsDeleted = async (args: {
  page: Page
  url: string
}): Promise<void> => {
  const { page, url } = args
  await Promise.all([
    page.waitForResponse((resp) => resp.status() === 404 && resp.request().method() === 'PROPFIND'),
    page.goto(url)
  ])
  await expect(page.locator(noLinkMessage)).toHaveText(
    'The resource could not be located, it may not exist anymore.'
  )
}

export interface changePageArgs {
  page: Page
  pageNumber: string
}

export const changePage = async (args: changePageArgs): Promise<void> => {
  const { page, pageNumber } = args
  await page.locator(util.format(listItemPageSelector, pageNumber)).click()
}

export const getCurrentPageNumber = (args: changePageArgs): Promise<string> => {
  const { page, pageNumber } = args
  return page.locator(util.format(currentPageSelector, pageNumber)).textContent()
}

export interface changeItemsPerPageArgs {
  page: Page
  itemsPerPage: string
}

export const changeItemsPerPage = async (args: changeItemsPerPageArgs): Promise<void> => {
  const { page, itemsPerPage } = args
  await page.locator(filesViewOptionButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'files control tippy box'
  )
  await page.locator(itemsPerPageDropDownSelector).click()
  await page.locator(util.format(itemsPerPageDropDownOptionSelector, itemsPerPage)).click()
  // close the files view option
  await page.locator(filesViewOptionButton).click()
}

export const getFileListFooterText = ({ page }: { page: Page }): Promise<string> => {
  return page.locator(footerTextSelector).textContent()
}

export interface expectNumberOfResourcesInThePageToBeArgs {
  page: Page
  numberOfResources: number
}

export const countNumberOfResourcesInThePage = async ({
  page
}: {
  page: Page
}): Promise<number> => {
  // playwright's default count function is not used here because count only counts
  // elements that are visible in the page but in this case we want to get
  // all the elements present
  await page.locator(footerTextSelector).waitFor()
  return page.evaluate(
    ([filesTableRowSelector]) => {
      return Promise.resolve(document.querySelectorAll(filesTableRowSelector).length)
    },
    [filesTableRowSelector]
  )
}

export const expectPageNumberNotToBeVisible = async ({ page }: { page: Page }): Promise<void> => {
  await expect(page.locator(filesPaginationNavSelector)).not.toBeVisible()
}

export interface expectFileToBeSelectedArgs {
  page: Page
  fileName: string
}

export const expectFileToBeSelected = async ({
  page,
  fileName
}: {
  page: Page
  fileName: string
}): Promise<void> => {
  await expect(page.locator(util.format(checkBox, fileName))).toBeChecked()
}

export const createShotcut = async (args: shortcutArgs): Promise<void> => {
  const { page, resource, name, type } = args
  await page.locator(addNewResourceButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['tippyBox'],
    'create new resource tippy box'
  )
  await page.locator(createNewShortcutButton).click()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['ocModal'],
    'create a shortcut modal'
  )

  switch (type) {
    case 'folder':
    case 'space':
    case 'file': {
      await page.locator(shortcutResorceInput).fill(resource)
      const searchResult = page.locator('#create-shortcut-modal-contextmenu .oc-resource-name')
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['tippyBox'],
        'shortcut preview button'
      )
      await expect(searchResult).toHaveText(resource)
      await searchResult.click()
      break
    }
    case 'public link':
      const link = new environment.LinksEnvironment()
      await page.locator(shortcutResorceInput).fill(link.getLink({ name: resource }).url)
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['tippyBox'],
        'shortcut preview button'
      )
      break
    case 'website': {
      await page.locator(shortcutResorceInput).fill(resource)
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['tippyBox'],
        'shortcut preview button'
      )
      await page.locator('#create-shortcut-modal-contextmenu').click()
      break
    }
  }

  if (name) {
    await page.getByRole('textbox', { name: 'Shortcut name' }).fill(name)
  }
  await Promise.all([
    page.waitForResponse(
      (resp) =>
        resp.status() === 201 && resp.request().method() === 'PUT' && resp.url().endsWith('url')
    ),
    page.locator(util.format(actionConfirmationButton, 'Create')).click()
  ])
}

export interface shortcutArgs {
  page: Page
  resource: string
  name: string
  type: shortcutType
}

export type shortcutType = 'folder' | 'file' | 'public link' | 'space' | 'website'

export const openShotcut = async ({
  page,
  name,
  url
}: {
  page: Page
  name: string
  url?: string
}): Promise<void> => {
  const resource = page.locator(util.format(resourceNameSelector, name))
  if (url) {
    const popupPromise = page.waitForEvent('popup')
    await resource.click()
    const popup = await popupPromise
    await popup.waitForURL(url)
  } else {
    const itemId = await resource.locator(fileRow).getAttribute('data-item-id')
    await Promise.all([
      page.waitForResponse(
        (resp) => resp.url().endsWith(encodeURIComponent(name)) || resp.url().endsWith(itemId)
      ),
      resource.click()
    ])
  }
}

export interface expectFileToBeLockedArgs {
  page: Page
  resource: string
}

export const getLockLocator = async (args: expectFileToBeLockedArgs): Promise<Locator> => {
  const { page, resource } = args
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['filesSpaceTable'],
    'files table with lock icons'
  )
  return page.locator(util.format(resourceLockIcon, resource))
}

export interface canManageResourceArgs {
  resource: string
  page: Page
}

export const canManageResource = async (args: canManageResourceArgs): Promise<boolean> => {
  const { resource, page } = args
  const notExpectedActions = ['move', 'rename', 'delete']
  await sidebar.open({ page: page, resource })
  await sidebar.openPanel({ page: page, name: 'actions' })
  const presentActions = await page.locator(sideBarActions).allTextContents()
  const presentActionsToLower = presentActions.map((actions) => actions.toLowerCase())
  for (const actions of notExpectedActions) {
    if (presentActionsToLower.includes(actions)) {
      return true
    }
  }
  return false
}

export const canEditDocumentContent = async ({
  page,
  type
}: {
  page: Page
  type: string
}): Promise<boolean> => {
  switch (type) {
    case 'OpenDocument':
      return canEditCollaboraDocument(page)
    case 'Microsoft Word':
      return canEditOnlyOfficeDocument(page)
  }
}

export const getAllAvailableActions = async ({
  page,
  resource
}: {
  page: Page
  resource: string
}): Promise<string[]> => {
  await sidebar.open({ page: page, resource })
  await sidebar.openPanel({ page: page, name: 'actions' })
  return await page.getByTestId('action-label').allTextContents()
}

export const getFileThumbnailLocator = (args: { page: Page; resource: string }): Locator => {
  const { page, resource } = args
  return page.locator(util.format(fileThumbnail, resource))
}

export const shouldSeeFilePreview = async ({
  page,
  resource
}: {
  page: Page
  resource: string
}): Promise<void> => {
  await sidebar.open({ page: page, resource })
  await expect(page.locator(fileIconPreview)).toHaveCSS('background-image', /blob/)
  await sidebar.close({ page: page })
}

export const shouldNotSeeFilePreview = async ({
  page,
  resource
}: {
  page: Page
  resource: string
}): Promise<void> => {
  await sidebar.open({ page: page, resource })
  await expect(page.locator(fileIconWrapper)).toBeVisible()
  await sidebar.close({ page: page })
}

export const checkActivity = async ({
  page,
  resource,
  activity
}: {
  page: Page
  resource: string
  activity: string
}): Promise<void> => {
  const paths = resource.split('/')
  const finalResource = paths.pop()
  for (const path of paths) {
    await clickResource({ page, path })
  }
  await sidebar.open({ page: page, resource: finalResource })
  await sidebar.openPanel({ page: page, name: 'activities' })
  await expect(page.getByTestId(activitySidebarPanel)).toBeVisible()
  await expect(page.locator(activitySidebarPanelBodyContent)).toContainText(activity)
}

export const checkEmptyActivity = async ({
  page,
  resource
}: {
  page: Page
  resource: string
}): Promise<void> => {
  const paths = resource.split('/')
  const finalResource = paths.pop()
  for (const path of paths) {
    await clickResource({ page, path })
  }
  await sidebar.open({ page: page, resource: finalResource })
  await sidebar.openPanel({ page: page, name: 'activities' })
  await expect(page.getByTestId(activitySidebarPanel)).toBeVisible()
  await expect(page.locator(activitySidebarPanelBodyContent)).toContainText('No activities')
}

export const duplicateResource = async ({
  page,
  resource,
  method
}: {
  page: Page
  resource: string
  method: string
}): Promise<void> => {
  switch (method) {
    case fileAction.dropDownMenu: {
      await page.locator(util.format(resourceNameSelector, resource)).click({ button: 'right' })
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['tippyBoxVisible', 'filesAppBarActions'],
        'file context menu'
      )
      await Promise.all([
        waitForResourceDuplicationResponse(page, resource),
        page.locator(util.format(filesContextMenuAction, 'duplicate')).click()
      ])
      break
    }
    case fileAction.batchAction: {
      await page.locator(util.format(checkBox, resource)).click()
      await objects.a11y.Accessibility.assertNoSevereA11yViolations(
        page,
        ['filesAppBarActions'],
        'files app bar'
      )
      await Promise.all([
        waitForResourceDuplicationResponse(page, resource),
        selectBatchAction(page, 'duplicate')
      ])
      break
    }
    case fileAction.sideBarPanel: {
      await sidebar.open({ page: page, resource: resource })
      await sidebar.openPanel({ page: page, name: 'actions' })

      await Promise.all([
        waitForResourceDuplicationResponse(page, resource),
        page.locator(util.format(sideBarActionButton, 'Duplicate')).click()
      ])
      break
    }
  }
}

const waitForResourceDuplicationResponse = (page: Page, resource: string): Promise<Response> => {
  return page.waitForResponse(
    (resp) =>
      resp.url().endsWith(resource) && resp.status() === 201 && resp.request().method() === 'COPY'
  )
}

export const getTagValidationMessage = async ({ page }: { page: Page }): Promise<string> => {
  return page.locator(tagValidationMessageSelector).innerText()
}

export const duplicateMultipleResources = async ({
  page,
  resources,
  method
}: {
  page: Page
  resources: string[]
  method: string
}): Promise<void> => {
  // tick all the checkboxes first
  for (const resource of resources) {
    await page.locator(util.format(checkBox, resource)).click()
  }

  switch (method) {
    case fileAction.dropDownMenu: {
      await page.locator(highlightedFileRowSelector).first().click({ button: 'right' })
      await Promise.all([
        ...resources.map((resource) => waitForResourceDuplicationResponse(page, resource)),
        page.locator(util.format(filesContextMenuAction, 'duplicate')).click()
      ])
      break
    }
    case fileAction.batchAction: {
      await Promise.all([
        ...resources.map((resource) => waitForResourceDuplicationResponse(page, resource)),
        selectBatchAction(page, 'duplicate')
      ])
      break
    }
  }
}
