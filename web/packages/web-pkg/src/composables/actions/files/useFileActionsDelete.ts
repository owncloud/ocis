import { useFileActionsDeleteResources } from '../helpers'
import {
  isLocationPublicActive,
  isLocationSpacesActive,
  isLocationTrashActive,
  isLocationCommonActive
} from '../../../router'
import { isProjectSpaceResource } from '@ownclouders/web-client'
import { useRouter } from '../../router'
import { useGettext } from 'vue3-gettext'
import { FileAction, FileActionOptions } from '../types'
import { computed } from 'vue'
import { useUserStore, useCapabilityStore } from '../../piniaStores'

export const useFileActionsDelete = () => {
  const userStore = useUserStore()
  const capabilityStore = useCapabilityStore()
  const router = useRouter()
  const { displayDialog, filesList_delete } = useFileActionsDeleteResources()

  const { $gettext } = useGettext()

  const handler = ({
    space,
    resources,
    deletePermanent
  }: FileActionOptions & { deletePermanent: boolean }) => {
    if (isLocationCommonActive(router, 'files-common-search')) {
      resources = resources.filter(
        (r) => r.canBeDeleted() && !r.isShareRoot() && !isProjectSpaceResource(r)
      )
    }
    if (deletePermanent) {
      displayDialog(space, resources)
      return
    }

    filesList_delete(resources)
  }

  const actions = computed((): FileAction[] => [
    {
      name: 'delete',
      icon: 'delete-bin-5',
      label: () => $gettext('Delete'),
      handler: ({ space, resources }) => handler({ space, resources, deletePermanent: false }),
      isVisible: ({ space, resources }) => {
        if (
          !isLocationSpacesActive(router, 'files-spaces-generic') &&
          !isLocationPublicActive(router, 'files-public-link') &&
          !isLocationCommonActive(router, 'files-common-search')
        ) {
          return false
        }

        if (resources.length === 0) {
          return false
        }

        if (
          isLocationSpacesActive(router, 'files-spaces-generic') &&
          space?.driveType === 'share' &&
          resources[0].path === '/'
        ) {
          return false
        }

        if (resources.length === 1 && resources[0].locked) {
          return false
        }

        if (isLocationCommonActive(router, 'files-common-search')) {
          return resources.some(
            (r) => r.canBeDeleted() && !r.isShareRoot() && !isProjectSpaceResource(r)
          )
        }

        const deleteDisabled = resources.some((resource) => {
          return !resource.canBeDeleted()
        })
        return !deleteDisabled
      },
      class: 'oc-files-actions-delete-trigger'
    },
    {
      // this menu item is ONLY for the trashbin (permanently delete a file/folder)
      name: 'delete-permanent',
      icon: 'delete-bin-5',
      label: () => $gettext('Delete'),
      handler: ({ space, resources }) => handler({ space, resources, deletePermanent: true }),
      isVisible: ({ space, resources }) => {
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

        return resources.length > 0
      },
      class: 'oc-files-actions-delete-permanent-trigger'
    }
  ])

  return {
    actions
  }
}
