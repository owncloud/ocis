import { computed } from 'vue'
import { useGettext } from 'vue3-gettext'
import { eventBus } from '../../services'
import { SideBarEventTopics } from '../sideBar'
import { Action } from './types'

export const useActionsShowDetails = () => {
  const { $gettext } = useGettext()

  const actions = computed((): Action[] => [
    {
      name: 'show-details',
      icon: 'information',
      label: () => $gettext('Details'),
      handler: () => eventBus.publish(SideBarEventTopics.open),
      isVisible: ({ resources }) => {
        return (resources as unknown[]).length > 0
      },
      class: 'oc-admin-settings-show-details-trigger'
    }
  ])

  return {
    actions
  }
}
