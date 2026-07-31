import {
  isLocationCommonActive,
  isLocationPublicActive,
  isLocationSpacesActive
} from '../../../router'
import { computed, unref } from 'vue'

import { useGettext } from 'vue3-gettext'
import { FileAction, FileActionOptions } from '../types'
import { isProjectSpaceResource, Resource, SpaceResource } from '@ownclouders/web-client'
import { useRouter } from '../../router'
import { useResourcesStore } from '../../piniaStores'
import { storeToRefs } from 'pinia'
import { ResourceTransfer, TransferType } from '../../../helpers/resource'
import { useClientService } from '../../clientService'
import { usePasteWorker } from '../../webWorkers/pasteWorker'

export const useFileActionsDuplicate = () => {
  const router = useRouter()
  const { $pgettext, $gettext, $ngettext } = useGettext()
  const clientService = useClientService()
  const { startWorker } = usePasteWorker()

  const resourcesStore = useResourcesStore()
  const { currentFolder } = storeToRefs(resourcesStore)

  const duplicateResources = async ({
    targetSpace,
    sourceSpace,
    resources
  }: {
    targetSpace: SpaceResource
    sourceSpace: SpaceResource
    resources: Resource[]
  }) => {
    const resourceTransfer = new ResourceTransfer(
      sourceSpace,
      resources,
      targetSpace,
      unref(currentFolder),
      currentFolder,
      clientService,
      $gettext,
      $ngettext
    )

    const transferData = await resourceTransfer.getTransferData(TransferType.DUPLICATE)

    if (!transferData.length) {
      return
    }

    startWorker(transferData, async ({ successful, failed }) => {
      resourceTransfer.showResultMessage(failed, successful, TransferType.DUPLICATE)

      if (!successful.length) {
        return
      }

      // handle store update, fetch resources first
      const loadingResources: Promise<void>[] = []
      const fetchedResources: Resource[] = []

      for (const resource of successful) {
        loadingResources.push(
          (async () => {
            const movedResource = await clientService.webdav.getFileInfo(targetSpace, resource)
            fetchedResources.push(movedResource)
          })()
        )
      }

      await Promise.allSettled(loadingResources)
      resourcesStore.upsertResources(fetchedResources)
    })
  }

  const handler = async ({ space: targetSpace, resources }: FileActionOptions) => {
    await duplicateResources({ targetSpace, sourceSpace: targetSpace, resources })
  }

  const actions = computed((): FileAction[] => {
    return [
      {
        name: 'duplicate',
        icon: 'folders',
        handler,
        label: () =>
          $pgettext(
            'Action to duplicate resources. Displayed in the resource context menu, app bar actions, and in the actions sidebar.',
            'Duplicate'
          ),
        isVisible: ({ resources }) => {
          if (
            !isLocationSpacesActive(router, 'files-spaces-generic') &&
            !isLocationPublicActive(router, 'files-public-link') &&
            !isLocationCommonActive(router, 'files-common-favorites') &&
            !isLocationCommonActive(router, 'files-common-search')
          ) {
            return false
          }

          if (isLocationSpacesActive(router, 'files-spaces-projects')) {
            return false
          }

          if (resources.length === 0) {
            return false
          }

          if (isLocationPublicActive(router, 'files-public-link')) {
            return unref(currentFolder)?.canCreate()
          }

          if (
            isLocationCommonActive(router, 'files-common-search') &&
            resources.every((r) => isProjectSpaceResource(r))
          ) {
            return false
          }

          if (!unref(resources)[0].canDownload()) {
            return false
          }

          return true
        },
        class: 'oc-files-actions-duplicate-trigger'
      }
    ]
  })

  return {
    actions
  }
}
