import { computed, unref } from 'vue'
import { isLocationCommonActive, isLocationSpacesActive } from '../../../router'
import { useGettext } from 'vue3-gettext'
import { FileAction, FileActionOptions, useIsFilesAppActive } from '../../actions'
import { useRouter } from '../../router'
import { useClientService } from '../../clientService'
import { useAbility } from '../../ability'
import { useMessages, useCapabilityStore, useResourcesStore } from '../../piniaStores'
import { useEventBus } from '../../eventBus'

export const useFileActionsFavorite = () => {
  const { showErrorMessage } = useMessages()
  const capabilityStore = useCapabilityStore()
  const router = useRouter()
  const { $gettext } = useGettext()
  const clientService = useClientService()
  const isFilesAppActive = useIsFilesAppActive()
  const ability = useAbility()
  const resourcesStore = useResourcesStore()
  const eventBus = useEventBus()

  const handler = async ({ space, resources }: FileActionOptions) => {
    try {
      const newValue = !resources[0].starred
      await clientService.webdav.setFavorite(space, resources[0], newValue)

      resourcesStore.updateResourceField({ id: resources[0].id, field: 'starred', value: newValue })
      if (!newValue) {
        eventBus.publish('app.files.list.removeFromFavorites', resources[0].id)
      }
    } catch (error) {
      const title = $gettext(
        'Failed to change favorite state of "%{file}"',
        { file: resources[0].name },
        true
      )
      showErrorMessage({ title, errors: [error] })
    }
  }

  const actions = computed((): FileAction[] => [
    {
      name: 'favorite',
      icon: 'star',
      handler,
      label: ({ resources }) => {
        if (resources[0].starred) {
          return $gettext('Remove from favorites')
        }
        return $gettext('Add to favorites')
      },
      isVisible: ({ resources }) => {
        if (
          unref(isFilesAppActive) &&
          !isLocationSpacesActive(router, 'files-spaces-generic') &&
          !isLocationCommonActive(router, 'files-common-favorites')
        ) {
          return false
        }
        if (resources.length !== 1) {
          return false
        }

        return capabilityStore.filesFavorites && ability.can('create', 'Favorite')
      },
      class: 'oc-files-actions-favorite-trigger'
    }
  ])

  return {
    actions
  }
}
