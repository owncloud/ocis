import { unref, computed } from 'vue'
import { useGettext } from 'vue3-gettext'
import { useSpaceHelpers } from '../../spaces'
import { useClientService } from '../../clientService'
import { useAbility } from '../../ability'
import { useRoute } from '../../router'
import { SpaceAction, SpaceActionOptions } from '../types'
import { SpaceResource } from '@ownclouders/web-client'
import {
  useMessages,
  useModals,
  useSharesStore,
  useSpacesStore,
  useUserStore
} from '../../piniaStores'

export const useSpaceActionsRename = () => {
  const { showMessage, showErrorMessage } = useMessages()
  const userStore = useUserStore()
  const { $gettext } = useGettext()
  const ability = useAbility()
  const clientService = useClientService()
  const route = useRoute()
  const { checkSpaceNameModalInput } = useSpaceHelpers()
  const { dispatchModal } = useModals()
  const spacesStore = useSpacesStore()
  const sharesStore = useSharesStore()

  const renameSpace = (space: SpaceResource, name: string) => {
    const graphClient = clientService.graphAuthenticated
    return graphClient.drives
      .updateDrive(space.id, { name }, sharesStore.graphRoles)
      .then(() => {
        if (unref(route).name === 'admin-settings-spaces') {
          space.name = name
        }
        spacesStore.updateSpaceField({ id: space.id, field: 'name', value: name })
        showMessage({ title: $gettext('Space name was changed successfully') })
      })
      .catch((error) => {
        console.error(error)
        showErrorMessage({
          title: $gettext('Failed to rename space'),
          errors: [error]
        })
      })
  }

  const handler = ({ resources }: SpaceActionOptions) => {
    if (resources.length !== 1) {
      return
    }

    dispatchModal({
      title: $gettext('Rename space') + ' ' + resources[0].name,
      confirmText: $gettext('Rename'),
      hasInput: true,
      inputLabel: $gettext('Space name'),
      inputValue: resources[0].name,
      onConfirm: (name: string) => renameSpace(resources[0], name),
      onInput: checkSpaceNameModalInput
    })
  }

  const actions = computed((): SpaceAction[] => [
    {
      name: 'rename',
      icon: 'pencil',
      label: () => {
        return $gettext('Rename')
      },
      handler,
      isVisible: ({ resources }) => {
        if (resources.length !== 1) {
          return false
        }

        return resources[0].canRename({ user: userStore.user, ability })
      },
      class: 'oc-files-actions-rename-trigger'
    }
  ])

  return {
    actions,

    // HACK: exported for unit tests:
    renameSpace
  }
}
