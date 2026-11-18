import { basename, join } from 'path'
import type { Resource, SpaceResource } from '@ownclouders/web-client'
import { ResolveStrategy, TransferType, type TransferData, type TransferFailure } from './types'
import { ConflictDialog } from './conflictDialog'
import { resolveFileNameDuplicate, isResourceBeeingMovedToSameLocation } from './conflictUtils'
import type { ClientService } from '../../../services'
import { getXRequestIdsLog, useMessages } from '../../../composables'
import { Ref, unref } from 'vue'
import type { Language } from 'vue3-gettext'
import { HttpError } from '@ownclouders/web-client'

/** At most this many lines of what wasn't copied are shown in the details, the rest is counted. */
const NOT_COPIED_LOG_LIMIT = 100

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
    errors: TransferFailure[],
    movedResources: Array<Resource>,
    transferType: TransferType
  ) {
    if (errors.length === 0) {
      const count = movedResources.length
      if (count === 0) {
        return
      }
      // a duplicate is a copy too
      const title =
        transferType === TransferType.COPY || transferType === TransferType.DUPLICATE
          ? this.$ngettext(
              '%{count} item was copied successfully',
              '%{count} items were copied successfully',
              count,
              { count: count.toString() }
            )
          : this.$ngettext(
              '%{count} item was moved successfully',
              '%{count} items were moved successfully',
              count,
              { count: count.toString() }
            )
      const messageStore = useMessages()
      messageStore.showMessage({ title, status: 'success' })
      return
    }
    const incomplete = this.isIncomplete(errors, transferType)
    const title = this.getFailureTitle(errors, transferType, incomplete)
    const descriptions: string[] = []
    if (incomplete) {
      descriptions.push(this.getCopiedItemsDescription(errors))
    }
    if (errors.some(({ error }) => error instanceof HttpError && error.statusCode === 507)) {
      descriptions.push(this.getInsufficientQuotaDescription(transferType))
    }
    const description = descriptions.filter(Boolean).join(' ')
    const messageStore = useMessages()
    messageStore.showErrorMessage({
      title,
      ...(description && { desc: description }),
      errors: errors.map(({ error }) => error),
      ...(incomplete && { errorLogContent: this.getNotCopiedLog(errors) })
    })
  }

  /** The copied folders exist, but some of their content is missing. Moves are one request. */
  private isIncomplete(errors: TransferFailure[], transferType: TransferType): boolean {
    return transferType !== TransferType.MOVE && errors.every(({ incomplete }) => incomplete)
  }

  /**
   * E.g. "87 of 100 items copied.", empty when the folders didn't contain any known item. There is no
   * verb, so the plural form only has to agree with the total ("of 100 items").
   */
  private getCopiedItemsDescription(errors: TransferFailure[]): string {
    const copied = errors.reduce((sum, { incomplete }) => sum + incomplete.copiedItems, 0)
    const total = errors.reduce((sum, { incomplete }) => sum + incomplete.totalItems, 0)
    if (total === 0) {
      return ''
    }
    return this.$ngettext(
      '%{copied} of %{total} item copied.',
      '%{copied} of %{total} items copied.',
      total,
      { copied: copied.toString(), total: total.toString() }
    )
  }

  /** What the "Details" of the message show: what wasn't copied, then the request ids. */
  private getNotCopiedLog(errors: TransferFailure[]): string {
    const lines = errors
      .flatMap(({ incomplete }) => incomplete.notCopied)
      .sort((a, b) => a.localeCompare(b))

    // a big failure, e.g. running out of storage, can leave thousands of items behind
    const shownLines = lines.slice(0, NOT_COPIED_LOG_LIMIT)
    const hiddenCount = lines.length - shownLines.length
    const requestIds = getXRequestIdsLog(errors.map(({ error }) => error))

    return [
      this.$gettext('Not copied:'),
      ...shownLines,
      ...(hiddenCount
        ? [
            this.$ngettext('…and %{count} more item', '…and %{count} more items', hiddenCount, {
              count: hiddenCount.toString()
            })
          ]
        : []),
      ...(requestIds ? ['', requestIds] : [])
    ].join('\r\n')
  }

  private getFailureTitle(
    errors: TransferFailure[],
    transferType: TransferType,
    incomplete: boolean
  ): string {
    const errorCount = errors.length
    const params = { name: errors[0]?.resourceName, count: errorCount.toString() }
    switch (transferType) {
      case TransferType.COPY:
        if (incomplete) {
          return this.$ngettext(
            'Copy of "%{name}" incomplete',
            'Copy of %{count} resources incomplete',
            errorCount,
            params
          )
        }
        return this.$ngettext(
          'Failed to copy "%{name}"',
          'Failed to copy %{count} resources',
          errorCount,
          params
        )
      case TransferType.DUPLICATE:
        if (incomplete) {
          return this.$ngettext(
            'Duplication of "%{name}" incomplete',
            'Duplication of %{count} resources incomplete',
            errorCount,
            params
          )
        }
        return this.$ngettext(
          'Failed to duplicate "%{name}"',
          'Failed to duplicate %{count} resources',
          errorCount,
          params
        )
      default:
        return this.$ngettext(
          'Failed to move "%{name}"',
          'Failed to move %{count} resources',
          errorCount,
          params
        )
    }
  }

  private getInsufficientQuotaDescription(transferType: TransferType): string {
    const spaceName =
      this.targetSpace.driveType === 'personal' ? this.$gettext('Personal') : this.targetSpace.name

    switch (transferType) {
      case TransferType.COPY:
        return this.$gettext(
          'The file cannot be copied because there is not enough storage left in "%{spaceName}".',
          { spaceName }
        )
      case TransferType.DUPLICATE:
        return this.$gettext(
          'The file cannot be duplicated because there is not enough storage left in "%{spaceName}".',
          { spaceName }
        )
      default:
        return this.$gettext(
          'The file cannot be moved because there is not enough storage left in "%{spaceName}".',
          { spaceName }
        )
    }
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
