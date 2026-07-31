import { canBeMoved } from '../../../helpers/permissions'
import {
  isLocationCommonActive,
  isLocationPublicActive,
  isLocationSpacesActive
} from '../../../router'
import { useGettext } from 'vue3-gettext'
import { ActionOptions, FileAction } from '../types'
import { computed, unref } from 'vue'
import { useRouter } from '../../router'
import { useClipboardStore, useResourcesStore } from '../../piniaStores'
import { Resource } from '@ownclouders/web-client'
import { storeToRefs } from 'pinia'

export const useFileActionsMove = () => {
  const router = useRouter()
  const { cutResources } = useClipboardStore()
  const language = useGettext()
  const { $gettext } = language

  const resourcesStore = useResourcesStore()
  const { currentFolder } = storeToRefs(resourcesStore)

  const isMacOs = computed(() => {
    return window.navigator.platform.match('Mac')
  })

  const cutShortcutString = computed(() => {
    if (unref(isMacOs)) {
      return $gettext('⌘ + X')
    }
    return $gettext('Ctrl + X')
  })

  const handler = ({ resources }: ActionOptions) => {
    cutResources(resources as Resource[])
  }
  const actions = computed((): FileAction[] => [
    {
      name: 'cut',
      icon: 'scissors',
      handler,
      shortcut: unref(cutShortcutString),
      label: () => $gettext('Cut'),
      isVisible: ({ resources }) => {
        if (
          !isLocationSpacesActive(router, 'files-spaces-generic') &&
          !isLocationPublicActive(router, 'files-public-link') &&
          !isLocationCommonActive(router, 'files-common-favorites')
        ) {
          return false
        }
        if (resources.length === 0) {
          return false
        }

        if (!unref(currentFolder)) {
          return false
        }

        if (resources.length === 1 && resources[0].locked) {
          return false
        }

        const moveDisabled = resources.some((resource) => {
          return canBeMoved(resource, unref(currentFolder).path) === false
        })
        return !moveDisabled
      },
      class: 'oc-files-actions-move-trigger'
    }
  ])

  return {
    actions
  }
}
