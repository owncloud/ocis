import { useModals } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import CreateFolderModal from '../components/CreateFolderModal.vue'

export const useCustomHandler = () => {
  const { dispatchModal } = useModals()
  const { $gettext } = useGettext()

  const customHandler = () => {
    dispatchModal({
      title: $gettext('Create a new password protected folder'),
      customComponent: CreateFolderModal,
      confirmText: $gettext('Create')
    })
  }

  return { customHandler }
}
