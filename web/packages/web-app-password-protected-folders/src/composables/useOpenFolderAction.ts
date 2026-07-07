import { FileAction, useClientService, useConfigStore, useModals } from '@ownclouders/web-pkg'
import { computed } from 'vue'
import { useGettext } from 'vue3-gettext'
import FolderViewModal from '../components/FolderViewModal.vue'

export const useOpenFolderAction = () => {
  const { $gettext, $pgettext } = useGettext()
  const { dispatchModal } = useModals()
  const clientService = useClientService()
  const configStore = useConfigStore()

  const action = computed<FileAction>(() => ({
    name: 'open-password-protected-folder',
    category: 'context',
    icon: 'external-link',
    async handler({ resources, space }) {
      const [file] = resources
      const { body } = await clientService.webdav.getFileContents(space, file)
      const publicLink = atob(body)
      const publicLinkUrl = new URL(publicLink)
      if (!['https:', 'http:'].includes(publicLinkUrl.protocol)) {
        throw new Error(
          $pgettext(
            'Error shown when opening a password-protected folder fails because the stored link has an unexpected format.',
            'This folder cannot be opened because the link it contains is invalid.'
          )
        )
      }
      if (publicLinkUrl.origin !== new URL(configStore.serverUrl).origin) {
        throw new Error(
          $pgettext(
            'Error shown when opening a password-protected folder fails because the stored link points to a different server.',
            'This folder cannot be opened because the link it contains does not point to this server.'
          )
        )
      }

      dispatchModal({
        title: resources.at(0).name,
        elementClass: 'folder-view-modal',
        customComponent: FolderViewModal,
        customComponentAttrs: () => ({
          publicLink,
          serverUrl: configStore.serverUrl
        }),
        hideConfirmButton: true,
        cancelText: $gettext('Close folder')
      })
    },
    label: () => $gettext('Open folder'),
    isDisabled: () => false,
    isVisible: ({ resources }) => {
      if (resources.length !== 1) {
        return false
      }

      return resources[0].extension === 'psec'
    },
    componentType: 'button',
    class: 'oc-files-actions-open-password-protected-folder'
  }))

  return action
}
