import { isLocationSpacesActive } from '../../../router'
import { usePreviewService } from '../../previewService'
import { useClientService } from '../../clientService'
import { useLoadingService } from '../../loadingService'
import { useRouter } from '../../router'
import { useGettext } from 'vue3-gettext'
import { computed } from 'vue'
import { FileAction, FileActionOptions } from '../types'
import { useMessages, useSharesStore, useSpacesStore, useUserStore } from '../../piniaStores'
import { useCreateSpace, useSpaceHelpers } from '../../spaces'

export const useFileActionsSetImage = () => {
  const { showMessage, showErrorMessage } = useMessages()
  const userStore = useUserStore()
  const router = useRouter()
  const { $gettext } = useGettext()
  const clientService = useClientService()
  const loadingService = useLoadingService()
  const previewService = usePreviewService()
  const spacesStore = useSpacesStore()
  const sharesStore = useSharesStore()
  const { createDefaultMetaFolder } = useCreateSpace()
  const { getDefaultMetaFolder } = useSpaceHelpers()

  const handler = async ({ space, resources }: FileActionOptions) => {
    const graphClient = clientService.graphAuthenticated
    const storageId = space?.id
    const { copyFiles, getFileInfo } = clientService.webdav

    try {
      let metaFolder = await getDefaultMetaFolder(space)
      if (!metaFolder) {
        metaFolder = await createDefaultMetaFolder(space)
      }

      if (resources[0].id !== space.spaceImageData?.id) {
        await copyFiles(
          space,
          { fileId: resources[0].id },
          space,
          { parentFolderId: metaFolder.id, name: resources[0].name },
          { overwrite: true }
        )
      }

      const { fileId } = await getFileInfo(space, { fileId: resources[0].id })
      const updatedSpace = await graphClient.drives.updateDrive(
        storageId,
        {
          name: space.name,
          special: [{ specialFolder: { name: 'image' }, id: fileId }]
        },
        sharesStore.graphRoles
      )

      spacesStore.updateSpaceField({
        id: storageId,
        field: 'spaceImageData',
        value: updatedSpace.spaceImageData
      })

      showMessage({ title: $gettext('Space image was set successfully') })
    } catch (error) {
      console.error(error)
      showErrorMessage({
        title: $gettext('Failed to set space image'),
        errors: [error]
      })
    }
  }

  const actions = computed((): FileAction[] => [
    {
      name: 'set-space-image',
      icon: 'image-edit',
      handler: (args) => loadingService.addTask(() => handler(args)),
      label: () => {
        return $gettext('Set as space image')
      },
      isVisible: ({ space, resources }) => {
        if (resources.length !== 1) {
          return false
        }
        if (!resources[0].mimeType) {
          return false
        }
        if (!previewService.isMimetypeSupported(resources[0].mimeType, true)) {
          return false
        }

        if (!isLocationSpacesActive(router, 'files-spaces-generic')) {
          return false
        }
        if (!space) {
          return false
        }

        return space.canEditImage({ user: userStore.user })
      },
      class: 'oc-files-actions-set-space-image-trigger'
    }
  ])

  return {
    actions
  }
}
