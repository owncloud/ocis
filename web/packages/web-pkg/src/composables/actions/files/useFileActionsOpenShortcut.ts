import {
  isLocationCommonActive,
  isLocationPublicActive,
  isLocationSharesActive,
  isLocationSpacesActive
} from '../../../router'
import { useIsFilesAppActive } from '../helpers'
import { useRouter } from '../../router'
import { FileAction, FileActionOptions } from '../types'
import { useIsSearchActive } from '../helpers'
import { computed, unref } from 'vue'
import { useGettext } from 'vue3-gettext'
import { useClientService } from '../../clientService'
import DOMPurify from 'dompurify'
import { useMessages } from '../../piniaStores'

export const useFileActionsOpenShortcut = () => {
  const { showErrorMessage } = useMessages()
  const router = useRouter()
  const { $gettext } = useGettext()
  const isFilesAppActive = useIsFilesAppActive()
  const isSearchActive = useIsSearchActive()
  const clientService = useClientService()

  const extractUrl = (fileContents: string) => {
    const regex = /URL=(.+)/
    const match = fileContents.match(regex)

    if (match && match[1]) {
      return match[1]
    } else {
      throw new Error('unable to extract url')
    }
  }
  const handler = async ({ resources, space }: FileActionOptions) => {
    try {
      const webURL = new URL(window.location.href)
      const fileContents = (await clientService.webdav.getFileContents(space, resources[0])).body
      let url = extractUrl(fileContents)

      // Add protocol if missing
      url = url.match(/^http[s]?:\/\//) ? url : `https://${url}`

      // Omit possible xss code
      url = DOMPurify.sanitize(url, { USE_PROFILES: { html: true } })

      if (url.startsWith(webURL.origin)) {
        window.location.href = url
        return
      }

      window.open(url)
    } catch (e) {
      console.error(e)
      showErrorMessage({
        title: $gettext('Failed to open shortcut'),
        errors: [e]
      })
    }
  }

  const actions = computed((): FileAction[] => [
    {
      name: 'open-shortcut',
      icon: 'external-link',
      category: 'context',
      handler,
      label: () => {
        return $gettext('Open shortcut')
      },
      isVisible: ({ resources }) => {
        if (
          unref(isFilesAppActive) &&
          !unref(isSearchActive) &&
          !isLocationSpacesActive(router, 'files-spaces-generic') &&
          !isLocationPublicActive(router, 'files-public-link') &&
          !isLocationCommonActive(router, 'files-common-favorites') &&
          !isLocationCommonActive(router, 'files-common-search') &&
          !isLocationSharesActive(router, 'files-shares-with-me') &&
          !isLocationSharesActive(router, 'files-shares-with-others') &&
          !isLocationSharesActive(router, 'files-shares-via-link')
        ) {
          return false
        }
        if (resources.length !== 1) {
          return false
        }
        if (resources[0].extension !== 'url') {
          return false
        }
        return resources[0].canDownload()
      },
      class: 'oc-files-actions-open-short-cut-trigger'
    }
  ])

  return {
    actions,

    // Hack for unit tests
    extractUrl
  }
}
