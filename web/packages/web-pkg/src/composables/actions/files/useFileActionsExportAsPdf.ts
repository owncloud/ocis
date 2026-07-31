import { FileAction, FileActionOptions } from '../types'
import { computed, Ref, unref } from 'vue'
import { useGettext } from 'vue3-gettext'
import { useModals } from '../../piniaStores'
import { useFolderLink } from '../../folderLink'
import ExportAsPdfModal from '../../../components/Modals/ExportAsPdfModal.vue'
import { useIsFilesAppActive } from '../helpers'

export const useFileActionsExportAsPdf = ({ content }: { content: Ref<unknown> }) => {
  const { $pgettext } = useGettext()
  const { dispatchModal } = useModals()
  const { getParentFolderLink } = useFolderLink()
  const isFilesAppActive = useIsFilesAppActive()

  function handler({ resources }: FileActionOptions) {
    const parentFolderLink = getParentFolderLink(resources[0])

    dispatchModal({
      elementClass: 'export-as-pdf-modal',
      title: $pgettext(
        'Title of the modal allowing users to select a location where they would like to store a file as a PDF.',
        'Export as PDF'
      ),
      customComponent: ExportAsPdfModal,
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
      name: 'export-as-pdf',
      icon: 'file-pdf',
      handler,
      label: () =>
        $pgettext(
          'File action available in editors that shows a modal allowing users to select a location where they would like to store the current file as a PDF.',
          'Export as PDF'
        ),
      isVisible: ({ resources }) =>
        !unref(isFilesAppActive) && resources.length === 1 && resources.at(0).extension === 'md',
      class: 'oc-files-actions-export-as-pdf-trigger'
    }
  ])

  return { actions }
}
