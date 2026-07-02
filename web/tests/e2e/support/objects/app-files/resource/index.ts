import { Download, Locator, Page, Response } from '@playwright/test'
import * as po from './actions'
import { Space } from '../../../types'
import { showShareIndicator } from './utils'
import { resourcePage } from '../../../../environment/constants'

export class Resource {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  async create(args: Omit<po.createResourceArgs, 'page'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.createResources({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
    await this.#page.locator('#files-view').waitFor()
  }

  async upload(args: Omit<po.uploadResourceArgs, 'page'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.uploadResource({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
  }

  async tryToUpload(args: Omit<po.uploadResourceArgs, 'page'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.tryToUploadResource({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
  }

  async uploadLargeNumberOfResources(args: Omit<po.uploadResourceArgs, 'page'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.uploadLargeNumberOfResources({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
  }

  async dropUpload(args: Omit<po.uploadResourceArgs, 'page'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.dropUploadFiles({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
  }

  // uploads the file but does check if the upload was successful
  // and does not navigate back to the startUrl
  startUpload(args: Omit<po.uploadResourceArgs, 'page'>): Promise<void> {
    return po.startResourceUpload({ ...args, page: this.#page })
  }

  pauseUpload(): Promise<void> {
    return po.pauseResourceUpload(this.#page)
  }

  async resumeUpload(): Promise<void> {
    const startUrl = this.#page.url()
    await po.resumeResourceUpload(this.#page)
    await this.#page.goto(startUrl)
  }

  cancelUpload(): Promise<void> {
    return po.cancelResourceUpload(this.#page)
  }

  async download(args: Omit<po.downloadResourcesArgs, 'page'>): Promise<Download[]> {
    const startUrl = this.#page.url()
    const downloads = await po.downloadResources({ ...args, page: this.#page })
    await this.#page.goto(startUrl)

    return downloads
  }

  async rename(args: Omit<po.renameResourceArgs, 'page'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.renameResource({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
  }

  async copy(args: Omit<po.moveOrCopyResourceArgs, 'page' | 'action'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.moveOrCopyResource({ ...args, page: this.#page, action: 'copy' })
    await this.#page.goto(startUrl)
  }

  async move(args: Omit<po.moveOrCopyResourceArgs, 'page' | 'action'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.moveOrCopyResource({ ...args, page: this.#page, action: 'move' })
    await this.#page.goto(startUrl)
  }

  async copyMultipleResources(args: Omit<po.moveOrCopyMultipleResourceArgs, 'page' | 'action'>) {
    const startUrl = this.#page.url()
    await po.moveOrCopyMultipleResources({ ...args, page: this.#page, action: 'copy' })
    await this.#page.goto(startUrl)
  }

  async moveMultipleResources(args: Omit<po.moveOrCopyMultipleResourceArgs, 'page' | 'action'>) {
    const startUrl = this.#page.url()
    await po.moveOrCopyMultipleResources({ ...args, page: this.#page, action: 'move' })
    await this.#page.goto(startUrl)
  }

  async delete(args: Omit<po.deleteResourceArgs, 'page'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.deleteResource({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
  }

  async open(): Promise<void> {}

  async restoreVersion(args: Omit<po.resourceVersionArgs, 'page'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.restoreResourceVersion({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
  }

  async downloadVersion(args: Omit<po.downloadResourceVersionArgs, 'page'>): Promise<Response[]> {
    const startUrl = this.#page.url()
    const downloads = await po.downloadResourceVersion({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
    return downloads
  }

  async deleteTrashBin(args: Omit<po.deleteResourceTrashbinArgs, 'page'>): Promise<string> {
    const startUrl = this.#page.url()
    const message = await po.deleteResourceTrashbin({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
    return message
  }

  async deleteTrashbinMultipleResources(
    args: Omit<po.deleteTrashbinMultipleResourcesArgs, 'page'>
  ): Promise<void> {
    const startUrl = this.#page.url()
    await po.deleteTrashbinMultipleResources({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
  }

  async emptyTrashbin({ page }: { page: Page }): Promise<void> {
    const startUrl = this.#page.url()
    await po.emptyTrashbin({ page })
    await this.#page.goto(startUrl)
  }

  async expectThatDeleteTrashBinButtonIsNotVisible(
    args: Omit<po.deleteResourceTrashbinArgs, 'page'>
  ): Promise<void> {
    return await po.expectThatDeleteButtonIsNotVisible({ ...args, page: this.#page })
  }

  async restoreTrashBin(args: Omit<po.restoreResourceTrashbinArgs, 'page'>): Promise<string> {
    const startUrl = this.#page.url()
    const message = await po.restoreTrashBinResource({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
    return message
  }

  async batchRestoreTrashBin(
    args: Omit<po.batchRestoreTrashbinResourcesArgs, 'page'>
  ): Promise<string> {
    const startUrl = this.#page.url()
    const message = await po.batchRestoreTrashBinResources({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
    return message
  }

  async expectThatRestoreTrashBinButtonIsNotVisible(
    args: Omit<po.restoreResourceTrashbinArgs, 'page'>
  ): Promise<void> {
    return await po.expectThatRestoreResourceButtonVisibility({ ...args, page: this.#page })
  }

  async areTagsVisibleForResourceInFilesTable(
    args: Omit<po.resourceTagsArgs, 'page'>
  ): Promise<boolean> {
    return await po.getTagsForResourceVisibilityInFilesTable({ ...args, page: this.#page })
  }

  async areTagsVisibleForResourceInDetailsPanel(
    args: Omit<po.resourceTagsArgs, 'page'>
  ): Promise<boolean> {
    return await po.getTagsForResourceVisibilityInDetailsPanel({
      ...args,
      page: this.#page
    })
  }

  async searchResource(args: Omit<po.searchResourceGlobalSearchArgs, 'page'>): Promise<void> {
    await po.searchResourceGlobalSearch({ ...args, page: this.#page })
  }

  // re-issues the global search before reading, so resources still being indexed when the
  // initial query ran are picked up instead of polling a stale one-shot result list
  reSearchAndGetDisplayedResources(): Promise<string[]> {
    return po.reSearchAndGetDisplayedResourcesFromSearch(this.#page)
  }

  getDisplayedResources(args: Omit<po.getDisplayedResourcesArgs, 'page'>): Promise<string[]> {
    switch (args.keyword) {
      case resourcePage.filesList:
        return po.getDisplayedResourcesFromFilesList(this.#page)
      case resourcePage.searchList:
        return po.getDisplayedResourcesFromSearch(this.#page)
      case resourcePage.shares:
        return po.getDisplayedResourcesFromShares(this.#page)
      case resourcePage.trashbin:
        return po.getDisplayedResourcesFromTrashbin(this.#page)
      default:
        throw new Error('Unknown keyword')
    }
  }

  async openFolder(resource: string): Promise<void> {
    await po.clickResource({ page: this.#page, path: resource })
  }

  async openFolderViaBreadcrumb(resource: string): Promise<void> {
    await po.clickResourceFromBreadcrumb({ page: this.#page, resource })
  }

  async switchToTilesViewMode(): Promise<void> {
    await po.clickViewModeToggle({ page: this.#page, target: 'resource-tiles' })
  }

  async expectThatResourcesAreTiles(): Promise<void> {
    await po.expectThatResourcesAreTiles({ page: this.#page })
  }

  async showHiddenFiles(): Promise<void> {
    await po.showHiddenResources(this.#page)
  }

  async toggleFlatList(): Promise<void> {
    await po.toggleFlatList(this.#page)
  }

  async getAllFiles(): Promise<string[]> {
    return po.getAllFiles(this.#page)
  }

  async editResource(args: Omit<po.editResourcesArgs, 'page'>): Promise<void> {
    await po.editResource({ ...args, page: this.#page })
  }

  async openFileInViewer(args: Omit<po.openFileInViewerArgs, 'page'>): Promise<void> {
    await po.openFileInViewer({ ...args, page: this.#page })
  }

  async addTags(args: Omit<po.resourceTagsArgs, 'page'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.addTagsToResource({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
  }

  async tryToAddTags(args: Omit<po.resourceTagsArgs, 'page'>): Promise<void> {
    await po.tryToAddTagsToResource({ ...args, page: this.#page })
  }

  async removeTags(args: Omit<po.resourceTagsArgs, 'page'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.removeTagsFromResource({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
  }

  async clickTag(args: Omit<po.clickTagArgs, 'page'>): Promise<void> {
    await po.clickResourceTag({ ...args, page: this.#page })
  }

  createSpaceFromFolder(args: Omit<po.createSpaceFromFolderArgs, 'page'>): Promise<Space> {
    return po.createSpaceFromFolder({ ...args, page: this.#page })
  }

  createSpaceFromSelection(args: Omit<po.createSpaceFromSelectionArgs, 'page'>): Promise<Space> {
    return po.createSpaceFromSelection({ ...args, page: this.#page })
  }

  async checkThatFileVersionIsNotAvailable(
    args: Omit<po.resourceVersionArgs, 'page'>
  ): Promise<void> {
    const startUrl = this.#page.url()
    await po.checkThatFileVersionIsNotAvailable({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
  }

  async checkThatFileVersionPanelIsNotAvailable(
    args: Omit<po.resourceVersionArgs, 'page'>
  ): Promise<void> {
    const startUrl = this.#page.url()
    await po.checkThatFileVersionPanelIsNotAvailable({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
  }

  async changePage(args: Omit<po.changePageArgs, 'page'>): Promise<void> {
    await po.changePage({ ...args, page: this.#page })
  }

  async getCurrentPageNumber(args: Omit<po.changePageArgs, 'page'>): Promise<string> {
    return await po.getCurrentPageNumber({ ...args, page: this.#page })
  }

  async changeItemsPerPage(args: Omit<po.changeItemsPerPageArgs, 'page'>): Promise<void> {
    await po.changeItemsPerPage({ ...args, page: this.#page })
  }

  getFileListFooterText(): Promise<string> {
    return po.getFileListFooterText({ page: this.#page })
  }

  countNumberOfResourcesInThePage(): Promise<number> {
    return po.countNumberOfResourcesInThePage({ page: this.#page })
  }

  async expectPageNumberNotToBeVisible(): Promise<void> {
    await po.expectPageNumberNotToBeVisible({ page: this.#page })
  }

  async expectFileToBeSelected(args: Omit<po.expectFileToBeSelectedArgs, 'page'>): Promise<void> {
    await po.expectFileToBeSelected({ ...args, page: this.#page })
  }

  async createShotcut(args: Omit<po.shortcutArgs, 'page'>): Promise<void> {
    const startUrl = this.#page.url()
    await po.createShotcut({ ...args, page: this.#page })
    await this.#page.goto(startUrl)
  }

  async openShotcut({ name, url }: { name: string; url?: string }): Promise<void> {
    await po.openShotcut({ page: this.#page, name: name, url: url })
  }

  async getLockLocator(args: Omit<po.expectFileToBeLockedArgs, 'page'>): Promise<Locator> {
    return await po.getLockLocator({ ...args, page: this.#page })
  }

  navigateMediaFile(navigationType: string): Promise<void> {
    return po.navigateMediaFile({ page: this.#page, navigationType })
  }

  async previewMediaFromSidebarPanel(resource: string): Promise<void> {
    await po.previewMediaFromSidebarPanel({ page: this.#page, resource })
  }

  showShareIndicatorSelector({
    buttonLabel,
    resource
  }: {
    buttonLabel: string
    resource: string
  }): Locator {
    return showShareIndicator({ page: this.#page, buttonLabel, resource })
  }

  async canManageResource(args: Omit<po.canManageResourceArgs, 'page'>): Promise<boolean> {
    return await po.canManageResource({ ...args, page: this.#page })
  }

  async canEditDocumentContent({ type }: { type: string }): Promise<boolean> {
    return await po.canEditDocumentContent({ page: this.#page, type })
  }

  async getAllAvailableActions({ resource }: { resource: string }): Promise<string[]> {
    return await po.getAllAvailableActions({ page: this.#page, resource })
  }

  getFileThumbnailLocator(resource: string): Locator {
    return po.getFileThumbnailLocator({ page: this.#page, resource })
  }

  async shouldSeeFilePreview({ resource }: { resource: string }): Promise<void> {
    await po.shouldSeeFilePreview({ page: this.#page, resource })
  }

  async shouldNotSeeFilePreview({ resource }: { resource: string }): Promise<void> {
    await po.shouldNotSeeFilePreview({ page: this.#page, resource })
  }

  async checkActivity({
    resource,
    activity
  }: {
    resource: string
    activity: string
  }): Promise<void> {
    await po.checkActivity({ page: this.#page, resource, activity })
  }

  async checkEmptyActivity({ resource }: { resource: string }): Promise<void> {
    await po.checkEmptyActivity({ page: this.#page, resource })
  }

  async openTemplateFile(resource: string, actionName: string): Promise<void> {
    await po.openTemplateFile({ page: this.#page, resource, webOffice: actionName })
  }

  async createFileFromTemplate(
    resource: string,
    webOffice: string,
    actionType: string
  ): Promise<void> {
    await po.createFileFromTemplate({ page: this.#page, resource, webOffice, actionType })
  }

  async duplicate(resource: string, method: string): Promise<void> {
    const startUrl = this.#page.url()
    await po.duplicateResource({ page: this.#page, resource, method })
    await this.#page.goto(startUrl)
  }

  async getTagValidationMessage(): Promise<string> {
    return po.getTagValidationMessage({ page: this.#page })
  }

  async duplicateMultipleResources(resources: string[], method: string): Promise<void> {
    const startUrl = this.#page.url()
    await po.duplicateMultipleResources({ page: this.#page, resources, method })
    await this.#page.goto(startUrl)
  }
}
