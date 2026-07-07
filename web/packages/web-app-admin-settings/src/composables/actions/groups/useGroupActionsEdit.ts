import { eventBus } from '@ownclouders/web-pkg'
import { SideBarEventTopics } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { computed } from 'vue'
import { GroupAction } from '@ownclouders/web-pkg'

export const useGroupActionsEdit = () => {
  const { $gettext } = useGettext()

  const actions = computed((): GroupAction[] => [
    {
      name: 'edit',
      icon: 'pencil',
      label: () => $gettext('Edit'),
      handler: () => eventBus.publish(SideBarEventTopics.openWithPanel, 'EditPanel'),
      isVisible: ({ resources }) => {
        return resources.length === 1 && !resources[0].groupTypes?.includes('ReadOnly')
      },
      class: 'oc-groups-actions-edit-trigger'
    }
  ])

  return {
    actions
  }
}
