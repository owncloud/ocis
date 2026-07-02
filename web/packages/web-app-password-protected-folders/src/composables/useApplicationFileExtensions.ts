import { PASSWORD_PROTECTED_FOLDER_FILE_EXTENSION } from '@ownclouders/web-client'
import { useGettext } from 'vue3-gettext'
import { shareType } from '../../../design-system/src/utils/shareType'
import { useCustomHandler } from './useCustomHandler'

export function useApplicationFileExtensions() {
  const { $gettext } = useGettext()
  const { customHandler } = useCustomHandler()

  return [
    {
      newFileMenu: {
        menuTitle: () => $gettext('Password Protected Folder'),
        isVisible: ({ currentFolder }) => !currentFolder?.shareTypes?.includes(shareType.link)
      },
      extension: PASSWORD_PROTECTED_FOLDER_FILE_EXTENSION,
      customHandler
    }
  ]
}
