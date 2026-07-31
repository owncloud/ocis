import { join } from 'path'
import { Resource } from '@ownclouders/web-client'
import { ResolveConflict, ResolveStrategy } from './types'
import { useModals } from '../../../composables'
import SpaceMoveInfoModal from '../../../components/Modals/SpaceMoveInfoModal.vue'
import ResourceConflictModal from '../../../components/Modals/ResourceConflictModal.vue'
import type { Language } from 'vue3-gettext'

export interface FileConflict {
  resource: Resource
  strategy?: ResolveStrategy
}

export class ConflictDialog {
  constructor(
    protected $gettext: Language['$gettext'],
    protected $ngettext: Language['$ngettext']
  ) {}

  async resolveAllConflicts(
    resourcesToMove: Resource[],
    targetFolder: Resource,
    targetFolderResources: Resource[]
  ): Promise<FileConflict[]> {
    // Collect all conflicting resources
    const allConflicts: FileConflict[] = []
    for (const resource of resourcesToMove) {
      const targetFilePath = join(targetFolder.path, resource.name)
      const exists = targetFolderResources.some((r) => r.path === targetFilePath)
      if (exists) {
        allConflicts.push({ resource, strategy: null })
      }
    }
    let count = 0
    let doForAllConflicts = false
    let doForAllConflictsStrategy = null
    const resolvedConflicts: FileConflict[] = []
    for (const conflict of allConflicts) {
      // Resolve conflicts accordingly
      if (doForAllConflicts) {
        conflict.strategy = doForAllConflictsStrategy
        resolvedConflicts.push(conflict)
        continue
      }

      // Resolve next conflict
      const conflictsLeft = allConflicts.length - count
      const result: ResolveConflict = await this.resolveFileExists(conflict.resource, conflictsLeft)
      conflict.strategy = result.strategy
      resolvedConflicts.push(conflict)
      count += 1

      // User checked 'do for all x conflicts'
      if (!result.doForAllConflicts) {
        continue
      }
      doForAllConflicts = true
      doForAllConflictsStrategy = result.strategy
    }
    return resolvedConflicts
  }

  resolveFileExists(
    resource: Resource,
    conflictCount: number,
    suggestMerge = false,
    separateSkipHandling = false // separate skip-handling between files and folders
  ): Promise<ResolveConflict> {
    const { dispatchModal } = useModals()

    return new Promise<ResolveConflict>((resolve) => {
      dispatchModal({
        variation: 'danger',
        title: resource.isFolder
          ? this.$gettext('Folder already exists')
          : this.$gettext('File already exists'),
        hideActions: true,
        customComponent: ResourceConflictModal,
        customComponentAttrs: () => ({
          resource,
          conflictCount,
          suggestMerge,
          separateSkipHandling,
          callbackFn: (conflict: ResolveConflict) => {
            resolve(conflict)
          }
        })
      })
    })
  }

  resolveDoCopyInsteadOfMoveForSpaces(): Promise<boolean> {
    const { dispatchModal } = useModals()

    return new Promise<boolean>((resolve) => {
      dispatchModal({
        variation: 'danger',
        title: this.$gettext('Copy here?'),
        customComponent: SpaceMoveInfoModal,
        confirmText: this.$gettext('Copy here'),
        onCancel: () => {
          resolve(false)
        },
        onConfirm: () => Promise.resolve(resolve(true))
      })
    })
  }
}
