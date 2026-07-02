import { basename, join } from 'path'
import type { Resource, SpaceResource } from '@ownclouders/web-client'
import { ResolveStrategy, TransferType, type TransferData } from './types'
import { ConflictDialog } from './conflictDialog'
import { resolveFileNameDuplicate, isResourceBeeingMovedToSameLocation } from './conflictUtils'
import type { ClientService } from '../../../services'
import { useMessages } from '../../../composables'
import { Ref, unref } from 'vue'
import type { Language } from 'vue3-gettext'
import { HttpError } from '@ownclouders/web-client'

export class ResourceTransfer extends ConflictDialog {
  constructor(
    private sourceSpace: SpaceResource,
    private resourcesToMove: Resource[],
    private targetSpace: SpaceResource,
    private targetFolder: Resource,
    private currentFolder: Ref<Resource>,
    private clientService: ClientService,
    $gettext: Language['$gettext'],
    $ngettext: Language['$ngettext']
  ) {
    super($gettext, $ngettext)
  }

  hasRecursion(): boolean {
    if (this.sourceSpace.id !== this.targetSpace.id) {
      return false
    }
    return this.resourcesToMove.some(
      (resource: Resource) => this.targetFolder.path === resource.path
    )
  }

  showRecursionErrorMessage() {
    const count = this.resourcesToMove.length
    const title = this.$ngettext(
      "You can't paste the selected file at this location because you can't paste an item into itself.",
      "You can't paste the selected files at this location because you can't paste an item into itself.",
      count
    )
    const messageStore = useMessages()
    messageStore.showErrorMessage({ title })
  }

  showResultMessage(
    errors: { resourceName: string; error: Error }[],
    movedResources: Array<Resource>,
    transferType: TransferType
  ) {
    if (errors.length === 0) {
      const count = movedResources.length
      if (count === 0) {
        return
      }
      const title =
        transferType === TransferType.COPY
          ? this.$ngettext(
              '%{count} item was copied successfully',
              '%{count} items were copied successfully',
              count,
              { count: count.toString() },
              true
            )
          : this.$ngettext(
              '%{count} item was moved successfully',
              '%{count} items were moved successfully',
              count,
              { count: count.toString() },
              true
            )
      const messageStore = useMessages()
      messageStore.showMessage({ title, status: 'success' })
      return
    }
    let title =
      transferType === TransferType.COPY
        ? this.$gettext(
            'Failed to copy %{count} resources',
            { count: errors.length.toString() },
            true
          )
        : this.$gettext(
            'Failed to move %{count} resources',
            { count: errors.length.toString() },
            true
          )
    if (errors.length === 1) {
      title =
        transferType === TransferType.COPY
          ? this.$gettext('Failed to copy "%{name}"', { name: errors[0]?.resourceName }, true)
          : this.$gettext('Failed to move "%{name}"', { name: errors[0]?.resourceName }, true)
    }
    let description = ''
    if (errors.some(({ error }) => error instanceof HttpError && error.statusCode === 507)) {
      description = this.$gettext('Insufficient quota')
    }
    const messageStore = useMessages()
    messageStore.showErrorMessage({
      title,
      ...(description && { desc: description }),
      errors: errors.map(({ error }) => error)
    })
  }

  /**
   * Gathers transfer data after resolving all potential conflicts.
   * This data can then be used to feed the web worker for pasting resources.
   */
  async getTransferData(transferType: TransferType) {
    if (this.hasRecursion()) {
      this.showRecursionErrorMessage()
      return []
    }
    if (this.sourceSpace.id !== this.targetSpace.id && transferType === TransferType.MOVE) {
      const doCopyInsteadOfMove = await this.resolveDoCopyInsteadOfMoveForSpaces()
      if (!doCopyInsteadOfMove) {
        return []
      }
      transferType = TransferType.COPY
    }

    const targetFolderResources = (
      await this.clientService.webdav.listFiles(this.targetSpace, this.targetFolder)
    ).children

    const resolvedConflicts =
      transferType === TransferType.DUPLICATE
        ? this.resourcesToMove.map((resource) => ({
            resource,
            strategy: ResolveStrategy.KEEP_BOTH
          }))
        : await this.resolveAllConflicts(
            this.resourcesToMove,
            this.targetFolder,
            targetFolderResources
          )

    const result: TransferData[] = []

    for (const resourceToMove of this.resourcesToMove) {
      // shallow copy of resources to prevent modifying existing rows
      const resource = { ...resourceToMove }
      const { id, name, extension } = resource

      const hasConflict = resolvedConflicts.some((e) => e.resource.id === id)
      let targetName = name
      let overwriteTarget = false

      if (hasConflict) {
        const resolveStrategy = resolvedConflicts.find((e) => e.resource.id === id)?.strategy
        if (resolveStrategy === ResolveStrategy.SKIP) {
          continue
        }
        if (resolveStrategy === ResolveStrategy.REPLACE) {
          if (this.isOverwritingParentFolder(resource, this.targetFolder, targetFolderResources)) {
            const error = new Error()
            this.showResultMessage([{ error, resourceName: name }], [], transferType)
            continue
          }
          overwriteTarget = true
        }
        if (resolveStrategy === ResolveStrategy.KEEP_BOTH) {
          targetName = resolveFileNameDuplicate(name, extension, targetFolderResources)
          resource.name = targetName
        }
      }

      if (
        isResourceBeeingMovedToSameLocation(
          this.sourceSpace,
          resource,
          this.targetSpace,
          this.targetFolder
        ) &&
        overwriteTarget
      ) {
        continue
      }

      result.push({
        resource,
        sourceSpace: this.sourceSpace,
        targetSpace: this.targetSpace,
        targetFolder: this.targetFolder,
        path: join(this.targetFolder.path, targetName),
        overwrite: overwriteTarget,
        transferType
      })
    }

    return result
  }

  // This is for an edge case if a user moves a subfolder with the same name as the parent folder into the parent of the parent folder (which is not possible because of the backend)
  public isOverwritingParentFolder(
    resource: Resource,
    targetFolder: Resource,
    targetFolderResources: Resource[]
  ) {
    if (resource.type !== 'folder') {
      return false
    }

    if (targetFolder.path === unref(this.currentFolder)?.path) {
      return false
    }

    const folderName = basename(resource.path)
    const newPath = join(targetFolder.path, folderName)
    return targetFolderResources.some((resource) => resource.path === newPath)
  }
}
