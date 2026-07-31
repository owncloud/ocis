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
import { useDownloadFile } from '../../download'

export const useFileActionsDownloadFile = () => {
  const router = useRouter()
  const { $gettext } = useGettext()
  const isFilesAppActive = useIsFilesAppActive()
  const isSearchActive = useIsSearchActive()
  const { downloadFile } = useDownloadFile()
  const handler = ({ space, resources }: FileActionOptions) => {
    downloadFile(space, resources[0])
  }

  const actions = computed((): FileAction[] => [
    {
      name: 'download-file',
      icon: 'file-download',
      handler,
      label: () => {
        return $gettext('Download')
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
        if (resources[0].isFolder) {
          return false
        }
        return resources[0].canDownload()
      },
      class: 'oc-files-actions-download-file-trigger'
    }
  ])

  return {
    actions
  }
}
