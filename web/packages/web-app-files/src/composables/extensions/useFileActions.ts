import {
  ActionExtension,
  useFileActionsCopyPermanentLink,
  useFileActionsOpenShortcut,
  useFileActionsShowShares
} from '@ownclouders/web-pkg'
import {
  contextActionsExtensionPoint,
  defaultActionsExtensionPoint,
  quickActionsExtensionPoint
} from '../../extensionPoints'
import { unref } from 'vue'

export const useFileActions = (): ActionExtension[] => {
  const { actions: openShortcutActions } = useFileActionsOpenShortcut()
  const { actions: showSharesActions } = useFileActionsShowShares()
  const { actions: permanentLinkActions } = useFileActionsCopyPermanentLink()

  return [
    {
      id: 'com.github.owncloud.web.files.context-action.open-shortcut',
      extensionPointIds: [contextActionsExtensionPoint.id, defaultActionsExtensionPoint.id],
      type: 'action',
      action: unref(openShortcutActions)[0]
    },
    {
      id: 'com.github.owncloud.web.files.quick-action.collaborator',
      extensionPointIds: [quickActionsExtensionPoint.id],
      type: 'action',
      action: unref(showSharesActions)[0]
    },
    {
      id: 'com.github.owncloud.web.files.quick-action.quicklink',
      extensionPointIds: [quickActionsExtensionPoint.id],
      type: 'action',
      action: unref(permanentLinkActions)[0]
    }
  ]
}
