import { computed } from 'vue'
import { SpaceAction } from '../types'
import { useGettext } from 'vue3-gettext'
import { useRouter } from '../../router'
import { SpaceResource } from '@ownclouders/web-client'
import { createLocationTrash } from '../../../router'
import { createFileRouteOptions } from '../../../helpers'
import { isProjectSpaceResource } from '@ownclouders/web-client'

export const useSpaceActionsNavigateToTrash = () => {
  const router = useRouter()
  const { $gettext } = useGettext()

  const getTrashLink = (space: SpaceResource) => {
    return createLocationTrash('files-trash-generic', {
      ...createFileRouteOptions(space, { fileId: space.fileId })
    })
  }

  const actions = computed((): SpaceAction[] => [
    {
      name: 'navigateToTrash',
      icon: 'delete-bin',
      label: () => $gettext('Show deleted files'),
      handler: ({ resources }) => {
        router.push(getTrashLink(resources[0]))
      },
      isVisible: ({ resources }) => {
        if (resources.length !== 1) {
          return false
        }

        if (!isProjectSpaceResource(resources[0])) {
          return false
        }

        if (resources[0].disabled) {
          return false
        }

        return true
      },
      class: 'oc-files-actions-navigate-to-trash-trigger'
    }
  ])

  return {
    actions
  }
}
