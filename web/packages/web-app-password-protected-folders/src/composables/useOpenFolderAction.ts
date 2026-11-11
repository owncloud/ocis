import {
  FileAction,
  renameResource,
  useClientService,
  useConfigStore,
  useModals,
  useResourcesStore
} from '@ownclouders/web-pkg'
import {
  PASSWORD_PROTECTED_FOLDER_FILE_EXTENSION,
  PASSWORD_PROTECTED_FOLDER_RENAMED_MESSAGE
} from '@ownclouders/web-client'
import { computed } from 'vue'
import { dirname, join } from 'path'
import { useGettext } from 'vue3-gettext'
import FolderViewModal from '../components/FolderViewModal.vue'

export const useOpenFolderAction = () => {
  const { $gettext, $pgettext } = useGettext()
  const { dispatchModal } = useModals()
  const clientService = useClientService()
  const configStore = useConfigStore()
  const { upsertResource } = useResourcesStore()

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
      const serverOrigin = new URL(configStore.serverUrl).origin
      if (publicLinkUrl.origin !== serverOrigin) {
        throw new Error(
          $pgettext(
            'Error shown when opening a password-protected folder fails because the stored link points to a different server.',
            'This folder cannot be opened because the link it contains does not point to this server.'
          )
        )
      }

      // The real folder is only reachable through the public link session running inside the
      // modal, so it may be renamed there. The `.psec` pointer file cannot be re-resolved from
      // the folder afterwards (the two are coupled only by name), so the framed app posts the
      // new folder name to us at rename time and we keep the `.psec` file in sync here.
      const onFolderRenamed = async (event: MessageEvent) => {
        if (event.origin !== serverOrigin) {
          return
        }
        if (event.data?.name !== PASSWORD_PROTECTED_FOLDER_RENAMED_MESSAGE) {
          return
        }

        const newName = event.data?.data?.newName
        const currentName = file.name.replace(
          new RegExp(`\\.${PASSWORD_PROTECTED_FOLDER_FILE_EXTENSION}$`),
          ''
        )
        if (!newName || newName === currentName) {
          return
        }

        const newPsecName = `${newName}.${PASSWORD_PROTECTED_FOLDER_FILE_EXTENSION}`
        const newPsecPath = join(dirname(file.path), newPsecName)

        await clientService.webdav.moveFiles(space, file, space, { path: newPsecPath })

        const updatedFile = { ...file }
        renameResource(space, updatedFile, newPsecPath)
        upsertResource(updatedFile)
      }

      window.addEventListener('message', onFolderRenamed)

      dispatchModal({
        title: resources.at(0).name,
        elementClass: 'folder-view-modal',
        customComponent: FolderViewModal,
        customComponentAttrs: () => ({
          publicLink,
          serverUrl: configStore.serverUrl
        }),
        hideConfirmButton: true,
        cancelText: $gettext('Close folder'),
        onCancel: () => {
          window.removeEventListener('message', onFolderRenamed)
        }
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
