import { dirname } from 'path'
import { isLocationTrashActive } from '../../../router'

import {
  Resource,
  isProjectSpaceResource,
  extractExtensionFromFile,
  SpaceResource,
  isTrashResource
} from '@ownclouders/web-client'
import {
  ResolveStrategy,
  ResolveConflict,
  resolveFileNameDuplicate,
  ConflictDialog
} from '../../../helpers/resource'
import { urlJoin } from '@ownclouders/web-client'
import { useClientService } from '../../clientService'
import { useRouter } from '../../router'
import { computed, unref } from 'vue'
import { useGettext } from 'vue3-gettext'
import type { FileAction, FileActionOptions } from '../types'
import {
  useMessages,
  useSpacesStore,
  useUserStore,
  useResourcesStore,
  useSharesStore
} from '../../piniaStores'
import { useRestoreWorker } from '../../webWorkers/restoreWorker'

export const useFileActionsRestore = () => {
  const { showMessage, showErrorMessage } = useMessages()
  const userStore = useUserStore()
  const router = useRouter()
  const { $gettext, $ngettext } = useGettext()
  const clientService = useClientService()
  const spacesStore = useSpacesStore()
  const sharesStore = useSharesStore()
  const resourcesStore = useResourcesStore()
  const { startWorker } = useRestoreWorker()

  // FIXME: use ConflictDialog class for this
  const collectConflicts = async (space: SpaceResource, sortedResources: Resource[]) => {
    const existingResourcesCache: Record<string, Resource[]> = {}
    const conflicts: Resource[] = []
    const resolvedResources: Resource[] = []
    const missingFolderPaths: string[] = []
    for (const resource of sortedResources) {
      const parentPath = dirname(resource.path)

      let existingResources: Resource[] = []
      if (parentPath in existingResourcesCache) {
        existingResources = existingResourcesCache[parentPath]
      } else {
        try {
          existingResources = (
            await clientService.webdav.listFiles(space, {
              path: parentPath
            })
          ).children
        } catch {
          missingFolderPaths.push(parentPath)
        }
        existingResourcesCache[parentPath] = existingResources
      }
      // Check for naming conflict in parent folder and between resources batch
      const hasConflict =
        existingResources.some((r) => r.name === resource.name) ||
        resolvedResources.filter((r) => r.id !== resource.id).some((r) => r.path === resource.path)
      if (hasConflict) {
        conflicts.push(resource)
      } else {
        resolvedResources.push(resource)
      }
    }
    return {
      existingResourcesByPath: existingResourcesCache,
      conflicts,
      resolvedResources,
      missingFolderPaths: missingFolderPaths.filter((path) => !existingResourcesCache[path]?.length)
    }
  }

  // FIXME: use ConflictDialog class for this
  const collectResolveStrategies = async (conflicts: Resource[]) => {
    let count = 0
    const resolvedConflicts = []
    const allConflictsCount = conflicts.length
    let doForAllConflicts = false
    let allConflictsStrategy: ResolveStrategy
    for (const conflict of conflicts) {
      const isFolder = conflict.type === 'folder'
      if (doForAllConflicts) {
        resolvedConflicts.push({
          resource: conflict,
          strategy: allConflictsStrategy
        })
        continue
      }
      const remainingConflictCount = allConflictsCount - count
      const conflictDialog = new ConflictDialog($gettext, $ngettext)
      const resolvedConflict: ResolveConflict = await conflictDialog.resolveFileExists(
        { name: conflict.name, isFolder } as Resource,
        remainingConflictCount,
        false
      )
      count++
      if (resolvedConflict.doForAllConflicts) {
        doForAllConflicts = true
        allConflictsStrategy = resolvedConflict.strategy
      }
      resolvedConflicts.push({
        resource: conflict,
        strategy: resolvedConflict.strategy
      })
    }
    return resolvedConflicts
  }

  const restoreResources = (
    space: SpaceResource,
    resources: Resource[],
    missingFolderPaths: string[]
  ) => {
    const originalRoute = unref(router.currentRoute)

    startWorker({ space, resources, missingFolderPaths }, async ({ successful, failed }) => {
      if (successful.length) {
        let title: string
        if (successful.length === 1) {
          title = $gettext('%{resource} was restored successfully', {
            resource: successful[0].name
          })
        } else {
          title = $gettext('%{resourceCount} files restored successfully', {
            resourceCount: successful.length.toString()
          })
        }
        showMessage({ title })

        // user hasn't navigated to another location meanwhile
        if (
          originalRoute.name === unref(router.currentRoute).name &&
          originalRoute.query?.fileId === unref(router.currentRoute).query?.fileId
        ) {
          resourcesStore.removeResources(successful)
          resourcesStore.resetSelection()
        }

        // Reload quota
        const graphClient = clientService.graphAuthenticated
        const updatedSpace = await graphClient.drives.getDrive(space.id, sharesStore.graphRoles)
        spacesStore.updateSpaceField({
          id: updatedSpace.id,
          field: 'spaceQuota',
          value: updatedSpace.spaceQuota
        })
      }

      if (failed.length) {
        let translated: string
        const translateParams: Record<string, string> = {}
        if (failed.length === 1) {
          translateParams.resource = failed[0].resource.name
          translated = $gettext('Failed to restore "%{resource}"', translateParams, true)
        } else {
          translateParams.resourceCount = failed.length.toString()
          translated = $gettext('Failed to restore %{resourceCount} files', translateParams, true)
        }
        showErrorMessage({ title: translated, errors: failed.map(({ error }) => error) })
      }
    })
  }

  const handler = async ({ space, resources }: FileActionOptions) => {
    // resources need to be sorted by path ASC to recover the parents first in case of deep nested folder structure
    const sortedResources = resources.sort((a, b) => a.path.length - b.path.length)

    // collect and request existing files in associated parent folders of each resource
    const { existingResourcesByPath, conflicts, resolvedResources, missingFolderPaths } =
      await collectConflicts(space, sortedResources)

    // iterate through conflicts and collect resolve strategies
    const resolvedConflicts = await collectResolveStrategies(conflicts)

    // iterate through conflicts and behave according to strategy
    const filesToOverwrite = resolvedConflicts
      .filter((e) => e.strategy === ResolveStrategy.REPLACE)
      .map((e) => e.resource)
    resolvedResources.push(...filesToOverwrite)
    const filesToKeepBoth = resolvedConflicts
      .filter((e) => e.strategy === ResolveStrategy.KEEP_BOTH)
      .map((e) => e.resource)

    for (let resource of filesToKeepBoth) {
      resource = { ...resource }
      const parentPath = dirname(resource.path)
      const existingResources = existingResourcesByPath[parentPath] || []
      const extension = extractExtensionFromFile(resource)
      const resolvedName = resolveFileNameDuplicate(resource.name, extension, [
        ...existingResources,
        ...resolvedConflicts.map((e) => e.resource),
        ...resolvedResources
      ])
      resource.name = resolvedName
      resource.path = urlJoin(parentPath, resolvedName)
      resolvedResources.push(resource)
    }

    return restoreResources(space, resolvedResources, missingFolderPaths)
  }

  const actions = computed((): FileAction[] => [
    {
      name: 'restore',
      icon: 'arrow-go-back',
      label: () => $gettext('Restore'),
      handler,
      isVisible: ({ space, resources }) => {
        if (!isLocationTrashActive(router, 'files-trash-generic')) {
          return false
        }
        if (!resources.every((r) => isTrashResource(r) && r.canBeRestored())) {
          return false
        }

        if (
          isProjectSpaceResource(space) &&
          !space.canRestoreFromTrashbin({ user: userStore.user })
        ) {
          return false
        }

        return resources.length > 0
      },
      class: 'oc-files-actions-restore-trigger'
    }
  ])

  return {
    actions,
    // HACK: exported for unit tests:
    restoreResources,
    collectConflicts
  }
}
