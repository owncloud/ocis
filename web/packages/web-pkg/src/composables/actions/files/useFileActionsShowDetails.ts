import { isLocationTrashActive } from '../../../router'
import { eventBus } from '../../../services/eventBus'
import { SideBarEventTopics } from '../../sideBar'
import { computed } from 'vue'
import { useGettext } from 'vue3-gettext'
import { useRouter } from '../../router'
import { FileAction } from '../types'
import { useResourcesStore } from '../../piniaStores'

export const useFileActionsShowDetails = () => {
  const router = useRouter()
  const resourcesStore = useResourcesStore()

  const { $gettext } = useGettext()

  const actions = computed((): FileAction[] => [
    {
      name: 'show-details',
      icon: 'information',
      class: 'oc-files-actions-show-details-trigger',
      label: () => $gettext('Details'),
      // we don't have details in the trashbin, yet.
      // remove trashbin route rule once we have them.
      isVisible: ({ resources }) => {
        if (isLocationTrashActive(router, 'files-trash-generic')) {
          return false
        }
        return resources.length > 0
      },
      handler({ resources }) {
        resourcesStore.setSelection(resources.map(({ id }) => id))
        eventBus.publish(SideBarEventTopics.open)
      }
    }
  ])

  return {
    actions
  }
}
