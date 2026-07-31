import { ActionExtension } from '@ownclouders/web-pkg'
import { computed, unref } from 'vue'
import { useOpenFolderAction } from './useOpenFolderAction'

export const useExtensions = () => {
  const action = useOpenFolderAction()

  const actionExtension = computed<ActionExtension>(() => ({
    id: 'com.github.owncloud.web-extensions.password-protected-folders',
    type: 'action',
    extensionPointIds: ['global.files.context-actions', 'global.files.default-actions'],
    action: unref(action)
  }))

  return computed(() => [unref(actionExtension)])
}
