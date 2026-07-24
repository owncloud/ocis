import { expect } from '@playwright/test'
import { objects } from '../../support'
import { createResourceTypes, shortcutType } from '../../support/objects/app-files/resource/actions'
import { editor } from '../../support/objects/app-files/utils'
import path from 'path'
import { Public } from '../../support/objects/app-files/page'
import { Resource } from '../../support/objects/app-files/resource'
import { config } from '../../config'
import * as runtimeFs from '../../support/utils/runtimeFs'
import { substitute } from '../../support/utils'
import { getWorld } from '../../environment/world'
import {
  fileAction,
  application,
  searchScope,
  resourcePage,
  shareIndicator
} from '../../environment/constants'

export async function userUploadsResources({
  stepUser,
  resources
}: {
  stepUser: string
  resources: { name: string; to?: string; type?: string; option?: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    await resourceObject.upload({
      to: resource.to,
      resources: [world.filesEnvironment.getFile({ name: resource.name })],
      option: resource.option,
      type: resource.type
    })
  }
}

export async function userShouldBeAbleToEditResource({
  stepUser,
  resource
}: {
  stepUser: string
  resource: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const userCanEdit = await resourceObject.canManageResource({ resource })
  expect(userCanEdit).toBe(true)
}

export async function userShouldNotBeAbleToEditResource({
  stepUser,
  resource
}: {
  stepUser: string
  resource: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const userCanEdit = await resourceObject.canManageResource({ resource })
  expect(userCanEdit).toBe(false)
}

export async function userCreatesResources({
  stepUser,
  resources
}: {
  stepUser: string
  resources: { name: string; type: createResourceTypes; content?: string; password?: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    await resourceObject.create({
      name: resource.name,
      type: resource.type,
      content: resource.content,
      password: resource.password
    })
  }
}

// export async function userSearchesGloballyWithFilter({
//   stepUser,
//   keyword,
//   filter,
//   command
// }: {
//   stepUser: string
//   keyword: string
//   filter: typeof searchScope.allFiles | typeof searchScope.currentFolder
//   command?: string
// }): Promise<void> {
//   const world = getWorld()
//   keyword = keyword ?? ''
//   const pressEnter = !!command && command.endsWith('presses enter')
//   const { page } = world.actorsEnvironment.getActor({ key: stepUser })
//   // let search indexing to complete
//   await page.waitForTimeout(1000)
//   const resourceObject = new objects.applicationFiles.Resource({ page })
//   await resourceObject.searchResource({
//     keyword,
//     filter: filter,
//     pressEnter
//   })
// }

export async function userSearchesGloballyWithFilter({
  stepUser,
  keyword = '',
  filter,
  command
}: {
  stepUser: string
  keyword: string
  filter: typeof searchScope.allFiles | typeof searchScope.currentFolder
  command?: string
}): Promise<void> {
  const { page } = getWorld().actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })

  await expect
    .poll(
      async () => {
        await resourceObject.searchResource({
          keyword,
          filter,
          pressEnter: command?.endsWith('presses enter') ?? false
        })

        return (await resourceObject.reSearchAndGetDisplayedResources()).includes(keyword)
      },
      {
        message: `Waiting for search to index "${keyword}"`,
        timeout: 30_000,
        intervals: [500, 1000, 2000]
      }
    )
    .toBe(true)
}

export async function userSwitchesToTilesViewMode({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.switchToTilesViewMode()
}

export async function userShouldSeeResourcesAsTiles({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.expectThatResourcesAreTiles()
}

export async function userOpensResource({
  stepUser,
  resource
}: {
  stepUser: string
  resource: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.openFolder(resource)
}

export async function userOpensResourceInViewer({
  stepUser,
  resource,
  viewer
}: {
  stepUser: string
  resource: string
  viewer:
    | typeof application.textEditor
    | typeof application.pdfViewer
    | typeof application.mediaViewer
    | typeof application.collabora
    | typeof application.onlyOffice
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.openFileInViewer({
    name: resource,
    actionType: viewer
  })
}

export async function userShouldSeeResources({
  listType,
  stepUser,
  resources
}: {
  listType:
    | typeof resourcePage.searchList
    | typeof resourcePage.filesList
    | typeof resourcePage.shares
    | typeof resourcePage.trashbin
  stepUser: string
  resources: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })

  // search list waits longer for tika full-text indexing; other lists only need UI render time
  const isSearchList = listType === resourcePage.searchList
  const timeout = isSearchList ? 30000 : 10000

  for (const resource of resources) {
    await expect
      .poll(
        async () => {
          // the global search dropdown is a one-shot query, so for the search list we
          // re-issue the search each poll to pick up resources that finish indexing after
          // the initial query instead of repeatedly reading a stale result list
          const actualList = isSearchList
            ? await resourceObject.reSearchAndGetDisplayedResources()
            : await resourceObject.getDisplayedResources({ keyword: listType })
          return actualList.includes(resource)
        },
        {
          message: `Waiting for resource "${resource}" to appear in ${listType}`,
          timeout,
          intervals: [500, 1000, 2000]
        }
      )
      .toBe(true)
  }
}

export async function userShouldNotSeeTheResources({
  listType,
  stepUser,
  resources
}: {
  listType:
    | typeof resourcePage.searchList
    | typeof resourcePage.filesList
    | typeof resourcePage.shares
    | typeof resourcePage.trashbin
  stepUser: string
  resources: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const actualList = await resourceObject.getDisplayedResources({
    keyword: listType
  })
  for (const resource of resources) {
    expect(actualList).not.toContain(resource)
  }
}

export async function userNavigatesToPageNumber({
  stepUser,
  pageNumber
}: {
  pageNumber: string
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.changePage({ pageNumber })
}

export async function userShouldSeeFooterText({
  stepUser,
  expectedText
}: {
  stepUser: string
  expectedText: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const actualText = await resourceObject.getFileListFooterText()
  expect(actualText).toBe(expectedText)
}

export async function userShouldSeeNumberOfResources({
  stepUser,
  expectedNumberOfResources
}: {
  stepUser: string
  expectedNumberOfResources: number
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const actualNumberOfResources = await resourceObject.countNumberOfResourcesInThePage()
  expect(actualNumberOfResources).toBe(expectedNumberOfResources)
}

export async function userEnablesShowHiddenFilesOption({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.showHiddenFiles()
}

export async function userShouldBeOnPage({
  stepUser,
  pageNumber
}: {
  stepUser: string
  pageNumber: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const currentPage = await resourceObject.getCurrentPageNumber({ pageNumber })
  expect(currentPage).toBe(pageNumber)
}

export async function userChangesItemsPerPage({
  stepUser,
  itemsPerPage
}: {
  stepUser: string
  itemsPerPage: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.changeItemsPerPage({ itemsPerPage })
}

export async function userShouldNotSeePagination({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.expectPageNumberNotToBeVisible()
}

export async function userEnablesFlatList({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.toggleFlatList()
}

export async function userShouldSeeFilesSortedAlphabetically({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const allFiles = await resourceObject.getAllFiles()
  const sortedFiles = [...allFiles].sort((a, b) =>
    a.localeCompare(b, 'en-us', { numeric: true, ignorePunctuation: true })
  )
  expect(allFiles).toEqual(sortedFiles)
}

export async function userCreatesSpaceFromFolderUsingContexMenu({
  stepUser,
  spaceName,
  folderName
}: {
  stepUser: string
  spaceName: string
  folderName: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const space = await resourceObject.createSpaceFromFolder({
    folderName: folderName,
    spaceName: spaceName
  })
  world.spacesEnvironment.createSpace({
    key: space.name,
    space: { name: space.name, id: space.id }
  })
}

export async function userCreatesSpaceFromResourcesUsingContexMenu({
  stepUser,
  spaceName,
  resources
}: {
  stepUser: string
  spaceName: string
  resources: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const space = await resourceObject.createSpaceFromSelection({ resources, spaceName })
  world.spacesEnvironment.createSpace({
    key: space.name,
    space: { name: space.name, id: space.id }
  })
}

export async function userAddsContentInTextEditor({
  stepUser,
  text
}: {
  stepUser: string
  text: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.Public({ page })
  await pageObject.fillDocumentContent({
    page,
    text,
    editor: application.textEditor
  })
}

export async function userSavesTextEditor({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  await editor.save(page)
}

export async function userClosesFileViewer({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  await editor.close(page)
}

export async function userDeletesResources({
  stepUser,
  actionType = fileAction.sideBarPanel,
  resources
}: {
  stepUser: string
  actionType: typeof fileAction.batchAction | typeof fileAction.sideBarPanel
  resources: { name: string; from?: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    await pageObject.delete({
      folder: resource.from === '' ? null : resource.from,
      resourcesWithInfo: [{ name: resource.name }],
      via: actionType
    })
  }
}

export async function userShouldBeAbleToDeleteResourceFromTrashbin({
  stepUser,
  resources
}: {
  stepUser: string
  resources: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    const message = await resourceObject.deleteTrashBin({ resource: resource })
    const paths = resource.split('/')
    expect(message).toBe(`"${paths[paths.length - 1]}" was deleted successfully`)
  }
}

export async function userShouldNotBeAbleToDeleteResourceFromTrashbin({
  stepUser,
  resources
}: {
  stepUser: string
  resources: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    await resourceObject.expectThatDeleteTrashBinButtonIsNotVisible({ resource: resource })
  }
}

export async function userShouldBeAbleToRestoreResourceFromTrashbin({
  stepUser,
  resources
}: {
  stepUser: string
  resources: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    const message = await resourceObject.restoreTrashBin({
      resource: resource
    })
    const paths = resource.split('/')
    expect(message).toBe(`${paths[paths.length - 1]} was restored successfully`)
  }
}

export async function userShouldNotBeAbleToRestoreResourceFromTrashbin({
  stepUser,
  resources
}: {
  stepUser: string
  resources: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    await resourceObject.expectThatRestoreTrashBinButtonIsNotVisible({
      resource: resource
    })
  }
}

export async function userCreatesShortcutForResources({
  stepUser,
  resources
}: {
  stepUser: string
  resources: { resource: string; name: string; type: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })

  for (const resource of resources) {
    await resourceObject.createShotcut({
      resource: resource.resource,
      name: resource.name,
      type: resource.type as shortcutType
    })
  }
}

type resourceToDownload = {
  resource: string
  type: string
  from?: string
}

export async function userDownloadsResource({
  stepUser,
  resourceToDownload,
  actionType
}: {
  stepUser: string
  resourceToDownload: resourceToDownload[]
  actionType:
    | typeof fileAction.sideBarPanel
    | typeof fileAction.batchAction
    | typeof fileAction.singleShareView
    | typeof fileAction.previewTopBar
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await processDownload(resourceObject, actionType, resourceToDownload)
}

export const processDownload = async (
  pageObject: Public | Resource,
  actionType:
    | typeof fileAction.sideBarPanel
    | typeof fileAction.batchAction
    | typeof fileAction.singleShareView
    | typeof fileAction.previewTopBar,
  resourceToDownload: resourceToDownload[]
) => {
  let downloads, files, parentFolder
  const downloadedResources: string[] = []

  const downloadInfo = resourceToDownload.reduce<Record<string, { name: string; type: string }[]>>(
    (acc, stepRow) => {
      const { resource, from, type } = stepRow
      const resourceInfo = {
        name: resource,
        type: type
      }
      if (!acc[from]) {
        acc[from] = []
      }
      acc[from].push(resourceInfo)
      return acc
    },
    {}
  )

  for (const folder of Object.keys(downloadInfo)) {
    files = downloadInfo[folder]
    parentFolder = folder !== 'undefined' ? folder : null

    downloads = await pageObject.download({
      folder: parentFolder,
      resources: files,
      via: actionType
    })

    downloads.forEach((download) => {
      const { name } = path.parse(download.suggestedFilename())
      downloadedResources.push(name)
    })

    if (actionType === fileAction.sideBarPanel || actionType === fileAction.previewTopBar) {
      expect(downloads.length).toBe(files.length)
      for (const resource of files) {
        const fileOrFolderName = path.parse(resource.name).name
        if (resource.type === 'file') {
          expect(downloadedResources).toContain(fileOrFolderName)
        } else {
          expect(downloadedResources).toContain('download')
        }
      }
    }
  }

  if (actionType === fileAction.batchAction) {
    expect(downloads.length).toBe(1)
    downloads.forEach((download) => {
      const { name } = path.parse(download.suggestedFilename())
      expect(name).toBe('download')
    })
  }
}

export async function userOpensShortcut({
  stepUser,
  name
}: {
  stepUser: string
  name: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.openShotcut({ name })
}
export async function userCanOpenShortcutWithExternalUrl({
  stepUser,
  name,
  url
}: {
  stepUser: string
  name: string
  url: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.openShotcut({ name, url })
}

export async function userShouldSeeContentInEditor({
  stepUser,
  expectedContent,
  editor
}: {
  stepUser: string
  expectedContent: string
  editor: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.Public({ page })
  const actualFileContent = await pageObject.getDocumentContent({
    page,
    editor
  })
  expect(actualFileContent.trim()).toBe(expectedContent)
}

export async function resourceShouldBeLockedForUser({
  stepUser,
  resource
}: {
  stepUser: string
  resource: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const lockLocator = await resourceObject.getLockLocator({ resource })
  expect(lockLocator).toBeVisible()
}

export async function resourceShouldNotBeLockedForUser({
  stepUser,
  resource
}: {
  stepUser: string
  resource: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const lockLocator = await resourceObject.getLockLocator({ resource })

  // can take more than 5 seconds for lock to be released in case of OnlyOffice
  expect(lockLocator).not.toBeVisible({ timeout: config.timeout * 1000 })
}

export async function userNavigatesToFolderViaBreadcrumb({
  stepUser,
  resource
}: {
  stepUser: string
  resource: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.openFolderViaBreadcrumb(resource)
}

export async function userEditsResources({
  stepUser,
  resources
}: {
  stepUser: string
  resources: { name: string; content: string; type?: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    await resourceObject.editResource({
      name: resource.name,
      type: resource.type,
      content: resource.content
    })
  }
}

export async function userShouldSeeThumbnailAndPreviewForFile({
  stepUser,
  resource
}: {
  stepUser: string
  resource: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await expect(resourceObject.getFileThumbnailLocator(resource)).toBeVisible()
  await resourceObject.shouldSeeFilePreview({ resource })
}

export async function userShouldNotSeePreviewForFile({
  stepUser,
  resource
}: {
  stepUser: string
  resource: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.shouldNotSeeFilePreview({ resource })
}

export async function userShouldNotSeeThumbnailAndPreviewForFile({
  stepUser,
  resource
}: {
  stepUser: string
  resource: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await expect(resourceObject.getFileThumbnailLocator(resource)).not.toBeVisible()
  await resourceObject.shouldNotSeeFilePreview({ resource })
}

export async function userOpensMediaUsingSidebarPanel({
  stepUser,
  resource
}: {
  stepUser: string
  resource: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.previewMediaFromSidebarPanel(resource)
}

export async function userNavigatesToMediaResource({
  stepUser,
  navigationType
}: {
  stepUser: string
  navigationType: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.navigateMediaFile(navigationType)
}

export async function userRenamesResource({
  stepUser,
  resource,
  newResourceName
}: {
  stepUser: string
  resource: string
  newResourceName: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.rename({ resource, newName: newResourceName })
}

export async function userShouldNotSeeVersionPanelForFiles({
  stepUser,
  file,
  to
}: {
  stepUser: string
  file: string
  to: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const fileInfo = world.filesEnvironment.getFile({ name: file })
  await resourceObject.checkThatFileVersionPanelIsNotAvailable({
    folder: to,
    files: [fileInfo]
  })
}

export async function userUploadsMultipleFilesInPersonalSpace({
  stepUser,
  numberOfFiles
}: {
  stepUser: string
  numberOfFiles: number
}): Promise<void> {
  const world = getWorld()
  const files = []
  for (let i = 0; i < numberOfFiles; i++) {
    const file = `file${i}.txt`
    runtimeFs.createFile(file, 'test content')

    files.push(
      world.filesEnvironment.getFile({
        name: path.join(runtimeFs.getTempUploadPath().replace(config.assets, ''), file)
      })
    )
  }

  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.uploadLargeNumberOfResources({ resources: files })
}

export async function userTriesToUploadResource({
  stepUser,
  resource,
  error,
  to
}: {
  stepUser: string
  resource: string
  error: string
  to: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.tryToUpload({
    to: to,
    resources: [world.filesEnvironment.getFile({ name: resource })],
    error: error
  })
}

export async function userUploadsResourcesViaDragNDrop({
  stepUser,
  resourceNames
}: {
  stepUser: string
  resourceNames: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const resources = resourceNames.map((name) => world.filesEnvironment.getFile({ name }))
  await resourceObject.dropUpload({ resources })
}

export async function userRestoresResourceVersion({
  stepUser,
  file,
  to,
  version,
  openDetailsPanel
}: {
  stepUser: string
  file: string
  to: string
  version: number
  openDetailsPanel: boolean
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })

  const fileInfo: Record<string, any> = {}
  if (!fileInfo[to]) {
    fileInfo[to] = []
  }

  fileInfo[to].push(world.filesEnvironment.getFile({ name: file }))

  if (version !== 1) {
    throw new Error('restoring is only supported for the most recent version')
  }
  fileInfo[to]['openDetailsPanel'] = openDetailsPanel

  for (const folder of Object.keys(fileInfo)) {
    await resourceObject.restoreVersion({
      folder,
      files: fileInfo[folder],
      openDetailsPanel: fileInfo[folder]['openDetailsPanel']
    })
  }
}

export async function userStartsUploadingFileFromTheTempUploadDir({
  stepUser,
  file,
  to,
  option
}: {
  stepUser: string
  file: string
  to?: string
  option?: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.startUpload({
    to: to,
    resources: [
      world.filesEnvironment.getFile({
        name: path.join(runtimeFs.getTempUploadPath().replace(config.assets, ''), file)
      })
    ],
    option: option
  })
}

export async function userPausesUpload({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.pauseUpload()
}

export async function userCancelsUpload({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.cancelUpload()
}

export async function userResumesUpload({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.resumeUpload()
}

export async function userShouldNotSeeAnyActivityOfResources({
  stepUser,
  resources
}: {
  stepUser: string
  resources: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })

  for (const resource of resources) {
    await resourceObject.checkEmptyActivity({ resource })
  }
}

export async function userShouldSeeActivityOfResources({
  stepUser,
  resources
}: {
  stepUser: string
  resources: { resource: string; activity: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })

  for (const info of resources) {
    await resourceObject.checkActivity({
      resource: info.resource,
      activity: substitute(info.activity)
    })
  }
}

export async function userCopiesResources({
  stepUser,
  actionType,
  resources
}: {
  stepUser: string
  actionType:
    | typeof fileAction.keyboard
    | typeof fileAction.sideBarPanel
    | typeof fileAction.dropDownMenu
    | typeof fileAction.batchAction
  resources: { resource: string; to: string; option?: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })

  for (const { resource, to, option } of resources) {
    await resourceObject.copy({
      resource,
      newLocation: to,
      method: actionType,
      option: option
    })
  }
}

export async function userMovesResources({
  stepUser,
  actionType,
  resources
}: {
  stepUser: string
  actionType:
    | typeof fileAction.keyboard
    | typeof fileAction.sideBarPanel
    | typeof fileAction.dropDownMenu
    | typeof fileAction.batchAction
    | typeof fileAction.dragDrop
    | typeof fileAction.dragDropBreadcrumb
  resources: { resource: string; to: string; option?: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })

  for (const { resource, to, option } of resources) {
    await resourceObject.move({
      resource,
      newLocation: to,
      method: actionType,
      option: option
    })
  }
}

export async function userDownloadsThePublicLinkResources({
  stepUser,
  actionType,
  resources
}: {
  stepUser: string
  actionType:
    | typeof fileAction.sideBarPanel
    | typeof fileAction.batchAction
    | typeof fileAction.singleShareView
  resources: resourceToDownload[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const pageObject = new objects.applicationFiles.page.Public({ page })
  await processDownload(pageObject, actionType, resources)
}

export async function userShouldSeeActionsForResource({
  stepUser,
  resource,
  actions
}: {
  stepUser: string
  resource: string
  actions: string[]
}) {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const action of actions) {
    const actions = await resourceObject.getAllAvailableActions({ resource })
    expect(actions.some((act) => act.startsWith(action))).toBe(true)
  }
}

export async function userShouldNotSeeActionsForResource({
  stepUser,
  resource,
  actions
}: {
  stepUser: string
  resource: string
  actions: string[]
}) {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const action of actions) {
    const actions = await resourceObject.getAllAvailableActions({ resource })
    expect(actions.some((act) => act.startsWith(action))).toBe(false)
  }
}

export async function userDeletesResourcesFromTrashbinUsingBatchAction({
  stepUser,
  resources
}: {
  stepUser: string
  resources: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.deleteTrashbinMultipleResources({ resources })
}

export async function userEmptiesTrashbin({ stepUser }: { stepUser: string }): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.emptyTrashbin({ page })
}

export async function userRestoresResourcesFromTrashbin({
  stepUser,
  resources,
  actionType
}: {
  stepUser: string
  resources: string[]
  actionType?: typeof fileAction.batchAction
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  if (actionType === fileAction.batchAction) {
    const message = await resourceObject.batchRestoreTrashBin({ resources })
    expect(message).toBe(`${resources.length} files restored successfully`)
  } else {
    for (const resource of resources) {
      const message = await resourceObject.restoreTrashBin({ resource: resource })
      const paths = resource.split('/')
      expect(message).toBe(`${paths[paths.length - 1]} was restored successfully`)
    }
  }
}

export async function userShouldNotBeAbleToEditContentOfResources({
  stepUser,
  resources
}: {
  stepUser: string
  resources: { name: string; type: string; content: string }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    const canEdit = await resourceObject.canEditDocumentContent({ type: resource.type })
    expect(canEdit).toBe(false)
  }
}

export async function userCreatesFileFromTemplateFile({
  stepUser,
  file,
  webOffice,
  actionType
}: {
  stepUser: string
  file: string
  webOffice: string
  actionType: typeof fileAction.sideBarPanel | typeof fileAction.contextMenu
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.createFileFromTemplate(file, webOffice, actionType)
}

export async function userOpensTemplateFileUsingContextMenu({
  stepUser,
  file,
  webOffice
}: {
  stepUser: string
  file: string
  webOffice: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.openTemplateFile(file, webOffice)
}

export async function userDuplicatesResources({
  stepUser,
  method,
  resources
}: {
  stepUser: string
  method:
    typeof fileAction.sideBarPanel | typeof fileAction.dropDownMenu | typeof fileAction.batchAction
  resources: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    await resourceObject.duplicate(resource, method)
  }
}

async function performCopyOrMoveMultipleResources({
  stepUser,
  actionType,
  newLocation,
  method,
  resources
}: {
  stepUser: string
  actionType: 'copy' | 'move'
  newLocation: string
  method:
    | typeof fileAction.dropDownMenu
    | typeof fileAction.batchAction
    | typeof fileAction.dragDrop
    | typeof fileAction.dragDropBreadcrumb
    | typeof fileAction.keyboard
  resources: string[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })

  // drag-n-drop always does MOVE
  if (method.includes(fileAction.dragDrop)) {
    expect(actionType).toBe('move')
  }

  const normalizedResources = [].concat(...resources)
  await resourceObject[actionType === 'copy' ? 'copyMultipleResources' : 'moveMultipleResources']({
    newLocation,
    method,
    resources: normalizedResources
  })
}

export const userCopiesResourcesAtOnce = (args) =>
  performCopyOrMoveMultipleResources({ ...args, actionType: 'copy' })

export const userMovesResourcesAtOnce = (args) =>
  performCopyOrMoveMultipleResources({ ...args, actionType: 'move' })

export async function userShouldSeeShareIndicatorOnResource({
  stepUser,
  buttonLabel,
  resource
}: {
  stepUser: string
  buttonLabel:
    | typeof shareIndicator.linkDirect
    | typeof shareIndicator.linkIndirect
    | typeof shareIndicator.userDirect
    | typeof shareIndicator.userIndirect
  resource: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const showShareIndicator = resourceObject.showShareIndicatorSelector({
    buttonLabel,
    resource
  })

  await expect(showShareIndicator).toBeVisible()
}

export async function userShouldNotSeeShareIndicatorOnResource({
  stepUser,
  buttonLabel,
  resource
}: {
  stepUser: string
  buttonLabel:
    | typeof shareIndicator.linkDirect
    | typeof shareIndicator.linkIndirect
    | typeof shareIndicator.userDirect
    | typeof shareIndicator.userIndirect
  resource: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const shareIndicator = resourceObject.showShareIndicatorSelector({
    buttonLabel,
    resource
  })
  await expect(shareIndicator).not.toBeVisible()
}

export async function userAddsFollowingTagsForResourcesUsingSidebarPanel({
  stepUser,
  resources
}: {
  stepUser: string
  resources: { name: string; tags: string[] }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    await resourceObject.addTags({
      resource: resource.name,
      tags: resource.tags.map((tag) => tag.trim().toLowerCase())
    })
  }
}

export async function resourceShouldContainTagsInFileList({
  stepUser,
  resources
}: {
  stepUser: string
  resources: { name: string; tags: string[] }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    const isVisible = await resourceObject.areTagsVisibleForResourceInFilesTable({
      resource: resource.name,
      tags: resource.tags.map((tag) => tag.trim().toLowerCase())
    })
    expect(isVisible).toBe(true)
  }
}

export async function resourceShouldContainTagsInDetailPanel({
  stepUser,
  resources
}: {
  stepUser: string
  resources: { name: string; tags: string[] }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    const isVisible = await resourceObject.areTagsVisibleForResourceInDetailsPanel({
      resource: resource.name,
      tags: resource.tags.map((tag) => tag.trim().toLowerCase())
    })
  }
}

export async function userRemovesTagsFromResourcesUsingSideBar({
  stepUser,
  resources
}: {
  stepUser: string
  resources: { name: string; tags: string[] }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    await resourceObject.removeTags({
      resource: resource.name,
      tags: resource.tags.map((tag) => tag.trim().toLowerCase())
    })
  }
}

export async function userTriesToAddTagForResourceUsingSidebarPanel({
  stepUser,
  resources
}: {
  stepUser: string
  resources: { name: string; tags: string[] }[]
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  for (const resource of resources) {
    await resourceObject.tryToAddTags({
      resource: resource.name,
      tags: resource.tags.map((tag) => tag.trim().toLowerCase())
    })
  }
}

export async function userShouldSeeFollowingTagValidationMessages({
  stepUser,
  message
}: {
  stepUser: string
  message: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const actualMessage = await resourceObject.getTagValidationMessage()
  expect(actualMessage).toBe(message)
}

export async function userClicksTheTagOnResource({
  stepUser,
  resourceName,
  tagName
}: {
  stepUser: string
  resourceName: string
  tagName: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  await resourceObject.clickTag({ resource: resourceName, tag: tagName.trim().toLowerCase() })
}

export async function userDownloadsPreviousVersionOfResource({
  stepUser,
  resource,
  to
}: {
  stepUser: string
  resource: string
  to: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const resourceObject = new objects.applicationFiles.Resource({ page })
  const fileInfo = world.filesEnvironment.getFile({ name: resource })
  await resourceObject.downloadVersion({ folder: to, files: [fileInfo] })
}
