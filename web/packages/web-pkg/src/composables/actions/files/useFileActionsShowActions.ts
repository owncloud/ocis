import { isLocationTrashActive } from '../../../router'
import { eventBus } from '../../../services/eventBus'
import { SideBarEventTopics } from '../../sideBar'
import { computed, unref } from 'vue'
import { useIsFilesAppActive } from '../helpers'
import { useRouter } from '../../router'
import { useGettext } from 'vue3-gettext'
import { FileAction } from '../types'

export const useFileActionsShowActions = () => {
  const router = useRouter()
  const { $gettext } = useGettext()
  const isFilesAppActive = useIsFilesAppActive()

  const handler = () => {
    // we don't have details in the trashbin, yet. the actions panel is the default
    // panel at the moment, so we need to use `null` as panel name for trashbins.
    // unconditionally return hardcoded `actions` once we have a dedicated
    // details panel in trashbins.
    const panelName = isLocationTrashActive(router, 'files-trash-generic') ? null : 'actions'
    eventBus.publish(SideBarEventTopics.openWithPanel, panelName)
  }

  const actions = computed((): FileAction[] => [
    {
      name: 'show-actions',
      icon: 'slideshow-3',
      label: () => $gettext('All Actions'),
      handler,
      isVisible: ({ resources }) => {
        // sidebar is currently only available inside files app
        if (!unref(isFilesAppActive)) {
          return false
        }

        // we don't have batch actions in the right sidebar, yet.
        // return hardcoded `true` in all cases once we have them.
        return resources.length === 1
      },
      class: 'oc-files-actions-show-actions-trigger'
    }
  ])

  return {
    actions
  }
}
