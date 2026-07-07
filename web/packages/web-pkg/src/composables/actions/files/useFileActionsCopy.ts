import {
  isLocationCommonActive,
  isLocationPublicActive,
  isLocationSpacesActive
} from '../../../router'
import { computed, unref } from 'vue'

import { useGettext } from 'vue3-gettext'
import { FileAction, FileActionOptions } from '../types'
import { isProjectSpaceResource } from '@ownclouders/web-client'
import { useRouter } from '../../router'
import { useConfigStore, useClipboardStore, useResourcesStore } from '../../piniaStores'
import { storeToRefs } from 'pinia'

export const useFileActionsCopy = () => {
  const configStore = useConfigStore()
  const router = useRouter()
  const { copyResources } = useClipboardStore()
  const language = useGettext()
  const { $gettext } = language

  const resourcesStore = useResourcesStore()
  const { currentFolder } = storeToRefs(resourcesStore)

  const isMacOs = computed(() => {
    return window.navigator.platform.match('Mac')
  })

  const copyShortcutString = computed(() => {
    if (unref(isMacOs)) {
      return $gettext('⌘ + C')
    }
    return $gettext('Ctrl + C')
  })

  const handler = ({ resources }: FileActionOptions) => {
    if (isLocationCommonActive(router, 'files-common-search')) {
      resources = resources.filter((r) => !isProjectSpaceResource(r))
    }

    copyResources(resources)
  }

  const actions = computed((): FileAction[] => {
    return [
      {
        name: 'copy',
        icon: 'file-copy-2',
        handler,
        shortcut: unref(copyShortcutString),
        label: () => $gettext('Copy'),
        isVisible: ({ resources }) => {
          if (
            !isLocationSpacesActive(router, 'files-spaces-generic') &&
            !isLocationPublicActive(router, 'files-public-link') &&
            !isLocationCommonActive(router, 'files-common-favorites') &&
            !isLocationCommonActive(router, 'files-common-search')
          ) {
            return false
          }
          if (isLocationSpacesActive(router, 'files-spaces-projects')) {
            return false
          }
          if (resources.length === 0) {
            return false
          }

          if (isLocationPublicActive(router, 'files-public-link')) {
            return unref(currentFolder)?.canCreate()
          }

          if (
            isLocationCommonActive(router, 'files-common-search') &&
            resources.every((r) => isProjectSpaceResource(r))
          ) {
            return false
          }

          if (unref(configStore.options.runningOnEos)) {
            // CERNBox does not allow actions above home/project root
            const elems = resources[0].path?.split('/').filter(Boolean) || [] //"/eos/project/c/cernbox"
            if (isLocationSpacesActive(router, 'files-spaces-generic') && elems.length < 5) {
              return false
            }
          }

          if (!unref(resources)[0].canDownload()) {
            return false
          }
          // copy can't be restricted in authenticated context, because
          // a user always has their home dir with write access
          return true
        },
        class: 'oc-files-actions-copy-trigger'
      }
    ]
  })

  return {
    actions
  }
}
