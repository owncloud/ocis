import {
  isLocationCommonActive,
  isLocationPublicActive,
  isLocationSpacesActive
} from '../../../router'
import { computed, unref } from 'vue'
import { useGettext } from 'vue3-gettext'
import { useGetMatchingSpace } from '../../spaces'
import { useClientService } from '../../clientService'
import { useRouter } from '../../router'
import { FileAction, FileActionOptions } from '../types'
import { Resource, SpaceResource, isShareSpaceResource } from '@ownclouders/web-client'
import {
  ClipboardMode,
  useClipboardStore,
  useConfigStore,
  useMessages,
  useResourcesStore
} from '../../piniaStores'
import { ClipboardActions, ResourceTransfer, TransferType } from '../../../helpers'
import { storeToRefs } from 'pinia'
import { usePasteWorker } from '../../webWorkers/pasteWorker'

export const useFileActionsPaste = () => {
  const router = useRouter()
  const clientService = useClientService()
  const { getMatchingSpace } = useGetMatchingSpace()
  const { $gettext, $ngettext } = useGettext()
  const { showErrorMessage } = useMessages()
  const clipboardStore = useClipboardStore()
  const configStore = useConfigStore()
  const { startWorker } = usePasteWorker()

  const resourcesStore = useResourcesStore()
  const { currentFolder } = storeToRefs(resourcesStore)
  const { resources: clipboardResources } = storeToRefs(clipboardStore)

  const isCrossModePasteAllowed = computed(() => {
    return (
      !clipboardStore.sourceMode ||
      clipboardStore.sourceMode !== ClipboardMode.Vault ||
      configStore.isInVault
    )
  })

  const getSourceSpace = (resource: Resource): SpaceResource => {
    const sourceBucketId = clipboardStore.getClipboardSourceSpaceKey(resource)
    const persistedSpace = clipboardStore.sourceSpaces[sourceBucketId]

    return (persistedSpace as SpaceResource) || getMatchingSpace(resource)
  }

  const isMacOs = computed(() => {
    return window.navigator.platform.match('Mac')
  })

  const pasteShortcutString = computed(() => {
    if (unref(isMacOs)) {
      return $gettext('⌘ + V')
    }
    return $gettext('Ctrl + V')
  })

  const transferType = computed(() => {
    if (clipboardStore.action === ClipboardActions.Cut) {
      return TransferType.MOVE
    }

    return TransferType.COPY
  })

  const pasteSelectedFiles = async ({
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

    const transferData = await resourceTransfer.getTransferData(unref(transferType))
    if (!transferData.length) {
      return
    }

    const originalCurrentFolderId = unref(currentFolder)?.id

    startWorker(transferData, async ({ successful, failed }) => {
      resourceTransfer.showResultMessage(failed, successful, unref(transferType))

      if (!successful.length) {
        return
      }

      // user has navigated to another location meanwhile -> no need to update store
      if (unref(currentFolder) && originalCurrentFolderId !== unref(currentFolder).id) {
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

      // FIXME: move to buildResource as soon as it has space context
      if (isShareSpaceResource(targetSpace)) {
        fetchedResources.forEach((r) => {
          r.remoteItemId = targetSpace.id
        })
      }

      resourcesStore.upsertResources(fetchedResources)
    })
  }

  /**
   * Computed property that determines if the user is attempting to cut and paste files into the same folder.
   *
   * This property checks if the current folder is the same as the parent folder of any resource in the clipboard.
   * If the current folder is of type 'space', the folder ID is appended with the owner ID. Otherwise, only the folder ID is used.
   * It also verifies if the clipboard action is a cut operation.
   * User is allowed to copy and paste into the same folder, but not cut and paste into the same folder.
   *
   * @returns {boolean} - Returns `true` if the user is trying to cut and paste files into the same folder, otherwise `false`.
   */
  const isCuttingAndPastingIntoSameFolder = computed(() => {
    const folderId =
      unref(currentFolder)?.type === 'space'
        ? `${unref(currentFolder).id}!${unref(currentFolder).owner.id}`
        : unref(currentFolder)?.id

    const isPastingIntoSameFolder =
      unref(folderId) &&
      unref(clipboardResources).some((resource) => resource.parentFolderId === unref(folderId))

    return clipboardStore.action === ClipboardActions.Cut && unref(isPastingIntoSameFolder)
  })

  const handler = async ({ space: targetSpace }: FileActionOptions) => {
    if (!unref(isCrossModePasteAllowed)) {
      showErrorMessage({
        title: $gettext('Pasting from Vault into the default mode is not supported.')
      })
      return
    }

    if (unref(isCuttingAndPastingIntoSameFolder)) {
      return
    }

    const resourceSpaceMapping = clipboardStore.resources.reduce<
      Record<string, { space: SpaceResource; resources: Resource[] }>
    >((acc, resource) => {
      const sourceBucketId = clipboardStore.getClipboardSourceSpaceKey(resource)
      const sourceSpace = getSourceSpace(resource)

      if (sourceBucketId in acc) {
        acc[sourceBucketId].resources.push(resource)
        return acc
      }

      if (!(sourceSpace.id in acc)) {
        acc[sourceSpace.id] = { space: sourceSpace, resources: [] }
      }

      acc[sourceSpace.id].resources.push(resource)
      return acc
    }, {})

    const promises = Object.values(resourceSpaceMapping).map(
      ({ space: sourceSpace, resources: resourcesToCopy }) => {
        return pasteSelectedFiles({ targetSpace, sourceSpace, resources: resourcesToCopy })
      }
    )
    await Promise.all(promises)
    clipboardStore.clearClipboard()
  }

  const actions = computed((): FileAction[] => [
    {
      name: 'paste',
      icon: 'clipboard',
      handler,
      label: () => $gettext('Paste'),
      shortcut: unref(pasteShortcutString),
      isVisible: ({ resources }) => {
        if (clipboardStore.resources.length === 0) {
          return false
        }
        if (!unref(isCrossModePasteAllowed)) {
          return false
        }
        if (
          !isLocationSpacesActive(router, 'files-spaces-generic') &&
          !isLocationPublicActive(router, 'files-public-link') &&
          !isLocationCommonActive(router, 'files-common-favorites')
        ) {
          return false
        }
        if (resources.length === 0) {
          return false
        }

        if (isLocationPublicActive(router, 'files-public-link') && unref(currentFolder)) {
          return unref(currentFolder)?.canCreate()
        }

        // copy can't be restricted in authenticated context, because
        // a user always has their home dir with write access
        return true
      },
      class: 'oc-files-actions-copy-trigger'
    }
  ])

  return {
    actions,
    isCuttingAndPastingIntoSameFolder
  }
}
