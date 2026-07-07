import { isLocationTrashActive } from '../../../router'
import { SpaceResource } from '@ownclouders/web-client'
import { isProjectSpaceResource } from '@ownclouders/web-client'
import { computed } from 'vue'
import { useClientService } from '../../clientService'
import { useRouter } from '../../router'
import { useGettext } from 'vue3-gettext'
import { FileAction, FileActionOptions } from '../types'
import {
  useMessages,
  useModals,
  useUserStore,
  useCapabilityStore,
  useResourcesStore
} from '../../piniaStores'
import { useLoadingService } from '../../loadingService'

export const useFileActionsEmptyTrashBin = () => {
  const { showMessage, showErrorMessage } = useMessages()
  const userStore = useUserStore()
  const capabilityStore = useCapabilityStore()
  const router = useRouter()
  const { $gettext } = useGettext()
  const clientService = useClientService()
  const { dispatchModal } = useModals()
  const resourcesStore = useResourcesStore()
  const loadingService = useLoadingService()

  const emptyTrashBin = ({ space }: { space: SpaceResource }) => {
    return clientService.webdav
      .clearTrashBin(space)
      .then(() => {
        showMessage({ title: $gettext('All deleted files were removed') })
        resourcesStore.clearResources()
        resourcesStore.resetSelection()
      })
      .catch((error) => {
        console.error(error)
        showErrorMessage({
          title: $gettext('Failed to empty trash bin'),
          errors: [error]
        })
      })
  }

  const handler = ({ space }: FileActionOptions) => {
    dispatchModal({
      variation: 'danger',
      title: $gettext('Empty trash bin'),
      confirmText: $gettext('Delete'),
      message: $gettext(
        'Are you sure you want to permanently delete the listed items? You can’t undo this action.'
      ),
      hasInput: false,
      onConfirm: () => loadingService.addTask(() => emptyTrashBin({ space }))
    })
  }

  const actions = computed((): FileAction[] => [
    {
      name: 'empty-trash-bin',
      icon: 'delete-bin-5',
      label: () => $gettext('Empty trash bin'),
      handler,
      isVisible: ({ space }) => {
        if (!isLocationTrashActive(router, 'files-trash-generic')) {
          return false
        }
        if (!capabilityStore.filesPermanentDeletion) {
          return false
        }

        if (
          isProjectSpaceResource(space) &&
          !space.canDeleteFromTrashBin({ user: userStore.user })
        ) {
          return false
        }

        return true
      },
      isDisabled: () => {
        return resourcesStore.activeResources.length === 0
      },
      class: 'oc-files-actions-empty-trash-bin-trigger',
      variation: 'danger',
      appearance: 'filled'
    }
  ])

  return {
    actions,
    // HACK: exported for unit tests:
    emptyTrashBin
  }
}
