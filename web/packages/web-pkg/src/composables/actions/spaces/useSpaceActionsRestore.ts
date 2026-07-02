import { SpaceResource } from '@ownclouders/web-client'
import { computed, unref } from 'vue'
import { SpaceAction, SpaceActionOptions } from '../types'
import { useRoute } from '../../router'
import { useAbility } from '../../ability'
import { useClientService } from '../../clientService'
import { useLoadingService } from '../../loadingService'
import { useGettext } from 'vue3-gettext'
import { isProjectSpaceResource } from '@ownclouders/web-client'
import {
  useMessages,
  useModals,
  useSharesStore,
  useSpacesStore,
  useUserStore
} from '../../piniaStores'

export const useSpaceActionsRestore = () => {
  const { showMessage, showErrorMessage } = useMessages()
  const userStore = useUserStore()
  const { $gettext, $ngettext } = useGettext()
  const ability = useAbility()
  const clientService = useClientService()
  const loadingService = useLoadingService()
  const route = useRoute()
  const { dispatchModal } = useModals()
  const spacesStore = useSpacesStore()
  const sharesStore = useSharesStore()

  const filterResourcesToRestore = (resources: SpaceResource[]): SpaceResource[] => {
    return resources.filter(
      (r) => isProjectSpaceResource(r) && r.canRestore({ user: userStore.user, ability })
    )
  }

  const restoreSpaces = async (spaces: SpaceResource[]) => {
    const client = clientService.graphAuthenticated
    const promises = spaces.map((space) =>
      client.drives
        .updateDrive(space.id, { name: space.name }, sharesStore.graphRoles, {
          headers: { Restore: 'true' }
        })
        .then((updatedSpace) => {
          if (unref(route).name === 'admin-settings-spaces') {
            space.disabled = false
            space.spaceQuota = updatedSpace.spaceQuota
          }
          spacesStore.updateSpaceField({ id: space.id, field: 'disabled', value: false })
          return true
        })
    )
    const results = await loadingService.addTask(() => {
      return Promise.allSettled(promises)
    })
    const succeeded = results.filter((r) => r.status === 'fulfilled')
    if (succeeded.length) {
      const title =
        succeeded.length === 1 && spaces.length === 1
          ? $gettext('Space "%{space}" was enabled successfully', { space: spaces[0].name })
          : $ngettext(
              '%{spaceCount} space was enabled successfully',
              '%{spaceCount} spaces were enabled successfully',
              succeeded.length,
              { spaceCount: succeeded.length.toString() },
              true
            )
      showMessage({ title })
    }

    const failed = results.filter((r) => r.status === 'rejected')
    if (failed.length) {
      failed.forEach(console.error)

      const title =
        failed.length === 1 && spaces.length === 1
          ? $gettext('Failed to enabled space "%{space}"', { space: spaces[0].name })
          : $ngettext(
              'Failed to enable %{spaceCount} space',
              'Failed to enable %{spaceCount} spaces',
              failed.length,
              { spaceCount: failed.length.toString() },
              true
            )
      showErrorMessage({
        title,
        errors: (failed as PromiseRejectedResult[]).map((f) => f.reason)
      })
    }
  }

  const handler = ({ resources }: SpaceActionOptions) => {
    const allowedResources = filterResourcesToRestore(resources)
    if (!allowedResources.length) {
      return
    }
    const message = $ngettext(
      'If you enable the selected space, it can be accessed again.',
      'If you enable the %{count} selected spaces, they can be accessed again.',
      allowedResources.length,
      { count: allowedResources.length.toString() }
    )
    const confirmText = $gettext('Enable')

    dispatchModal({
      title: $ngettext(
        'Enable Space "%{space}"?',
        'Enable %{spaceCount} Spaces?',
        allowedResources.length,
        {
          space: allowedResources[0].name,
          spaceCount: allowedResources.length.toString()
        }
      ),
      confirmText,
      icon: 'alert',
      message,
      hasInput: false,
      onConfirm: () => restoreSpaces(allowedResources)
    })
  }

  const actions = computed((): SpaceAction[] => [
    {
      name: 'restore',
      icon: 'play-circle',
      label: () => $gettext('Enable'),
      handler,
      isVisible: ({ resources }) => {
        return !!filterResourcesToRestore(resources).length
      },
      class: 'oc-files-actions-restore-trigger'
    }
  ])

  return {
    actions,

    // HACK: exported for unit tests:
    restoreSpaces
  }
}
