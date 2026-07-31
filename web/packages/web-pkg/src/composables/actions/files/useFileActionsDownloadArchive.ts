import {
  isLocationCommonActive,
  isLocationPublicActive,
  isLocationSharesActive,
  isLocationSpacesActive
} from '../../../router'
import { useIsFilesAppActive } from '../helpers'
import path from 'path'
import first from 'lodash-es/first'
import { isProjectSpaceResource, isPublicSpaceResource, Resource } from '@ownclouders/web-client'
import { computed, unref } from 'vue'
import { useLoadingService } from '../../loadingService'
import { useRouter } from '../../router'

import { FileAction, FileActionOptions } from '../types'
import { useGettext } from 'vue3-gettext'
import { useArchiverService } from '../../archiverService'
import { formatFileSize } from '../../../helpers/filesize'
import { useAuthStore, useMessages } from '../../piniaStores'

export const useFileActionsDownloadArchive = () => {
  const { showErrorMessage } = useMessages()
  const router = useRouter()
  const loadingService = useLoadingService()
  const archiverService = useArchiverService()
  const { $ngettext, $gettext, current } = useGettext()
  const authStore = useAuthStore()
  const isFilesAppActive = useIsFilesAppActive()

  const handler = ({ space, resources }: FileActionOptions) => {
    if (resources.length > 1) {
      // the handler can be triggered successfully if project spaces are selected along with other files.
      // but we must filter out the project spaces in such a case (only the other selected files are allowed for download).
      resources = resources.filter((r) => r.canDownload() && !isProjectSpaceResource(r))
    }

    const fileOptions = unref(archiverService.fileIdsSupported)
      ? {
          fileIds: resources.map((resource) => resource.fileId)
        }
      : {
          dir: path.dirname(first<Resource>(resources).path) || '/',
          files: resources.map((resource) => resource.name)
        }

    return archiverService
      .triggerDownload({
        ...fileOptions,
        ...(space &&
          isPublicSpaceResource(space) && {
            publicToken: space.id as string,
            publicLinkPassword: authStore.publicLinkPassword,
            publicLinkShareOwner: space.publicLinkShareOwner,
            signatureAuth: resources[0].signatureAuth
          })
      })
      .catch((e) => {
        console.error(e)
        showErrorMessage({
          title: $ngettext(
            'Failed to download the selected folder.', // on single selection only available for folders
            'Failed to download the selected files.', // on multi selection available for files+folders
            resources.length
          ),
          errors: [e]
        })
      })
  }

  const areArchiverLimitsExceeded = (resources: Resource[]) => {
    const archiverCapabilities = unref(archiverService.capability)
    if (!archiverCapabilities) {
      return
    }

    const selectedFilesSize = resources.reduce(
      (accumulator, currentValue) => accumulator + parseInt(`${currentValue.size}`),
      0
    )

    return selectedFilesSize > parseInt(archiverCapabilities.max_size)
  }

  const actions = computed((): FileAction[] => {
    return [
      {
        name: 'download-archive',
        icon: 'inbox-archive',
        handler: async (args) => {
          await loadingService.addTask(() => handler(args))
        },
        label: () => $gettext('Download'),
        disabledTooltip: ({ resources }) => {
          return areArchiverLimitsExceeded(resources)
            ? $gettext('The selection exceeds the allowed archive size (max. %{maxSize})', {
                maxSize: formatFileSize(unref(archiverService.capability).max_size, current)
              })
            : ''
        },
        isDisabled: ({ resources }) => areArchiverLimitsExceeded(resources),
        isVisible: ({ resources }) => {
          if (
            unref(isFilesAppActive) &&
            !isLocationSpacesActive(router, 'files-spaces-generic') &&
            !isLocationPublicActive(router, 'files-public-link') &&
            !isLocationCommonActive(router, 'files-common-favorites') &&
            !isLocationCommonActive(router, 'files-common-search') &&
            !isLocationSharesActive(router, 'files-shares-with-me') &&
            !isLocationSharesActive(router, 'files-shares-with-others') &&
            !isLocationSharesActive(router, 'files-shares-via-link') &&
            !isLocationCommonActive(router, 'files-common-search')
          ) {
            return false
          }
          if (!unref(archiverService.available)) {
            return false
          }

          if (resources.length === 0) {
            return false
          }
          if (resources.length === 1 && !resources[0].isFolder) {
            return false
          }
          if (resources.length > 1 && resources.every((r) => isProjectSpaceResource(r))) {
            return false
          }
          if (isProjectSpaceResource(resources[0]) && resources[0].disabled) {
            return false
          }
          if (
            !unref(archiverService.fileIdsSupported) &&
            isLocationCommonActive(router, 'files-common-favorites')
          ) {
            return false
          }

          const downloadDisabled = resources.some((resource) => {
            return !resource.canDownload()
          })
          return !downloadDisabled
        },
        class: 'oc-files-actions-download-archive-trigger'
      }
    ]
  })

  return {
    actions
  }
}
