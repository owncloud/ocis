import { eventBus } from '../../../services/eventBus'
import { SideBarEventTopics } from '../../sideBar'
import { computed } from 'vue'
import { SpaceAction, SpaceActionOptions } from '../types'
import { useGettext } from 'vue3-gettext'
import { useResourcesStore } from '../../piniaStores'

export const useSpaceActionsShowMembers = () => {
  const { $gettext } = useGettext()
  const resourcesStore = useResourcesStore()

  const handler = ({ resources }: SpaceActionOptions) => {
    resourcesStore.setSelection(resources.map(({ id }) => id))
    eventBus.publish(SideBarEventTopics.openWithPanel, 'space-share')
  }

  const actions = computed((): SpaceAction[] => [
    {
      name: 'show-members',
      icon: 'group',
      label: () => $gettext('Members'),
      handler,
      isVisible: ({ resources }) => resources.length === 1 && !resources[0].disabled,
      class: 'oc-files-actions-show-details-trigger'
    }
  ])

  return {
    actions
  }
}
