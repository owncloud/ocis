import PQueue from 'p-queue'
import { isLocationSharesActive } from '../../../router'
import { useClientService } from '../../clientService'
import { useLoadingService } from '../../loadingService'
import { useRouter } from '../../router'
import { computed } from 'vue'
import { useGettext } from 'vue3-gettext'
import { FileAction, FileActionOptions } from '../../actions'
import { useMessages, useConfigStore, useResourcesStore } from '../../piniaStores'
import { IncomingShareResource } from '@ownclouders/web-client'

export const useFileActionsToggleHideShare = () => {
  const { showMessage, showErrorMessage } = useMessages()
  const router = useRouter()
  const { $gettext } = useGettext()

  const clientService = useClientService()
  const loadingService = useLoadingService()
  const configStore = useConfigStore()
  const { updateResourceField, resetSelection } = useResourcesStore()

  const handler = async ({ resources }: FileActionOptions<IncomingShareResource>) => {
    const errors: Error[] = []
    const triggerPromises: Promise<void>[] = []
    const triggerQueue = new PQueue({
      concurrency: configStore.options.concurrentRequests.resourceBatchActions
    })
    const hidden = !resources[0].hidden

    resources.forEach((resource) => {
      triggerPromises.push(
        triggerQueue.add(async () => {
          try {
            await clientService.graphAuthenticated.driveItems.updateDriveItem(
              resource.driveId,
              resource.id,
              { '@UI.Hidden': hidden }
            )

            updateResourceField<IncomingShareResource>({
              id: resource.id,
              field: 'hidden',
              value: hidden
            })
          } catch (error) {
            console.error(error)
            errors.push(error)
          }
        })
      )
    })

    await Promise.all(triggerPromises)

    if (errors.length === 0) {
      resetSelection()
      showMessage({
        title: hidden
          ? $gettext('The share was hidden successfully')
          : $gettext('The share was unhidden successfully')
      })

      return
    }

    showErrorMessage({
      title: hidden
        ? $gettext('Failed to hide the share')
        : $gettext('Failed to unhide share share'),
      errors
    })
  }

  const actions = computed((): FileAction<IncomingShareResource>[] => [
    {
      name: 'toggle-hide-share',
      icon: 'eye-off', // FIXME: change icon based on hidden status
      handler: (args) => loadingService.addTask(() => handler(args)),
      label: ({ resources }) => (resources[0].hidden ? $gettext('Unhide') : $gettext('Hide')),
      isVisible: ({ resources }) => {
        if (resources.length === 0) {
          return false
        }

        return isLocationSharesActive(router, 'files-shares-with-me')
      },
      class: 'oc-files-actions-hide-share-trigger'
    }
  ])

  return {
    actions
  }
}
