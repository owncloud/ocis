import {
  Key,
  KeyboardActions,
  Modifier,
  useClipboardStore,
  useResourcesStore
} from '@ownclouders/web-pkg'

export const useKeyboardTableActions = (keyActions: KeyboardActions) => {
  const resourcesStore = useResourcesStore()
  const { copyResources } = useClipboardStore()

  keyActions.bindKeyAction({ modifier: Modifier.Ctrl, primary: Key.C }, () => {
    copyResources(resourcesStore.selectedResources)
  })
}
