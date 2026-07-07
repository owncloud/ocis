import { eventBus } from '@ownclouders/web-pkg'
import { SideBarEventTopics } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { computed } from 'vue'
import { UserAction } from '@ownclouders/web-pkg'

export const useUserActionsEdit = () => {
  const { $gettext } = useGettext()

  const actions = computed((): UserAction[] => [
    {
      name: 'edit',
      icon: 'pencil',
      label: () => $gettext('Edit'),
      handler: () => eventBus.publish(SideBarEventTopics.openWithPanel, 'EditPanel'),
      isVisible: ({ resources }) => {
        return resources.length === 1
      },
      class: 'oc-users-actions-edit-trigger'
    }
  ])

  return {
    actions
  }
}
