import { computed, unref } from 'vue'
import { SpaceAction, SpaceActionOptions } from '../types'
import { useRoute } from '../../router'
import { useAbility } from '../../ability'
import { useClientService } from '../../clientService'
import { useGettext } from 'vue3-gettext'
import { SpaceResource } from '@ownclouders/web-client'
import {
  useMessages,
  useModals,
  useSharesStore,
  useSpacesStore,
  useUserStore
} from '../../piniaStores'

export const useSpaceActionsEditDescription = () => {
  const { showMessage, showErrorMessage } = useMessages()
  const userStore = useUserStore()
  const { $gettext } = useGettext()
  const ability = useAbility()
  const clientService = useClientService()
  const route = useRoute()
  const { dispatchModal } = useModals()
  const spacesStore = useSpacesStore()
  const sharesStore = useSharesStore()

  const editDescriptionSpace = (space: SpaceResource, description: string) => {
    const graphClient = clientService.graphAuthenticated
    return graphClient.drives
      .updateDrive(space.id, { name: space.name, description }, sharesStore.graphRoles)
      .then(() => {
        spacesStore.updateSpaceField({ id: space.id, field: 'description', value: description })
        if (unref(route).name === 'admin-settings-spaces') {
          space.description = description
        }
        showMessage({ title: $gettext('Space subtitle was changed successfully') })
      })
      .catch((error) => {
        console.error(error)
        showErrorMessage({
          title: $gettext('Failed to change space subtitle'),
          errors: [error]
        })
      })
  }

  const handler = ({ resources }: SpaceActionOptions) => {
    if (resources.length !== 1) {
      return
    }

    dispatchModal({
      title: $gettext('Change subtitle for space') + ' ' + resources[0].name,
      confirmText: $gettext('Confirm'),
      hasInput: true,
      inputLabel: $gettext('Space subtitle'),
      inputValue: resources[0].description,
      onConfirm: (description: string) => editDescriptionSpace(resources[0], description)
    })
  }

  const actions = computed((): SpaceAction[] => [
    {
      name: 'editDescription',
      icon: 'h-2',
      iconFillType: 'none',
      label: () => {
        return $gettext('Edit subtitle')
      },
      handler,
      isVisible: ({ resources }) => {
        if (resources.length !== 1) {
          return false
        }

        return resources[0].canEditDescription({ user: userStore.user, ability })
      },
      class: 'oc-files-actions-edit-description-trigger'
    }
  ])

  return {
    actions,

    // HACK: exported for unit tests:
    editDescriptionSpace
  }
}
