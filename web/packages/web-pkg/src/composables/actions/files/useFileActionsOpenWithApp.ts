import { FileAction, FileActionOptions } from '../types'
import { computed, unref } from 'vue'
import { useGettext } from 'vue3-gettext'
import { useAppsStore, useModals } from '../../piniaStores'
import { storeToRefs } from 'pinia'
import FilePickerModal from '../../../components/Modals/FilePickerModal.vue'
import { useFolderLink } from '../../folderLink'
import { useIsFilesAppActive } from '../helpers'

export const useFileActionsOpenWithApp = ({ appId }: { appId: string }) => {
  const { $gettext } = useGettext()
  const isFilesAppActive = useIsFilesAppActive()
  const { dispatchModal } = useModals()
  const appsStore = useAppsStore()
  const { apps } = storeToRefs(appsStore)
  const { getParentFolderLink } = useFolderLink()

  const handler = ({ resources }: FileActionOptions) => {
    const app = unref(apps)[appId]
    const parentFolderLink = getParentFolderLink(resources[0])

    dispatchModal({
      elementClass: 'open-with-app-modal',
      title: $gettext('Open file in %{app}', { app: app.name }),
      customComponent: FilePickerModal,
      hideActions: true,
      customComponentAttrs: () => ({
        app,
        parentFolderLink
      }),
      focusTrapInitial: false
    })
  }

  const actions = computed((): FileAction[] => [
    {
      name: 'open-with-app',
      icon: 'folder-open',
      handler,
      label: () => {
        return $gettext('Open')
      },
      isVisible: () => {
        return !unref(isFilesAppActive)
      },
      class: 'oc-files-actions-open-with-app-trigger'
    }
  ])

  return {
    actions
  }
}
