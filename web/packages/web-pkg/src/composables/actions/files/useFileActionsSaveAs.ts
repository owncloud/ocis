import { FileAction, FileActionOptions } from '../types'
import { computed, unref, Ref } from 'vue'
import { useGettext } from 'vue3-gettext'
import { useModals } from '../../piniaStores'
import SaveAsModal from '../../../components/Modals/SaveAsModal.vue'
import { useFolderLink } from '../../folderLink'
import { useIsFilesAppActive } from '../helpers'

export const useFileActionsSaveAs = ({ content }: { content: Ref<unknown> }) => {
  const { $gettext } = useGettext()
  const isFilesAppActive = useIsFilesAppActive()
  const { dispatchModal } = useModals()
  const { getParentFolderLink } = useFolderLink()

  const handler = ({ resources }: FileActionOptions) => {
    const parentFolderLink = getParentFolderLink(resources[0])

    dispatchModal({
      elementClass: 'save-as-modal',
      title: $gettext('Save as'),
      customComponent: SaveAsModal,
      hideActions: true,
      customComponentAttrs: () => ({
        content: unref(content),
        parentFolderLink,
        originalResource: resources[0]
      }),
      focusTrapInitial: false
    })
  }

  const actions = computed((): FileAction[] => [
    {
      name: 'save-as',
      icon: 'save-2',
      handler,
      label: () => {
        return $gettext('Save as')
      },
      isVisible: ({ resources }) => {
        return !unref(isFilesAppActive) || resources.length !== 1
      },
      class: 'oc-files-actions-save-as-trigger'
    }
  ])

  return {
    actions
  }
}
