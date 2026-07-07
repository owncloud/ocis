import {
  isLocationSharesActive,
  isLocationSpacesActive,
  createLocationShares
} from '../../../router'
import PQueue from 'p-queue'
import { IncomingShareResource } from '@ownclouders/web-client'
import { useClientService } from '../../clientService'
import { useLoadingService } from '../../loadingService'
import { useRouter } from '../../router'
import { computed } from 'vue'
import { useGettext } from 'vue3-gettext'
import { FileAction, FileActionOptions } from '../types'
import { useMessages, useConfigStore, useResourcesStore } from '../../piniaStores'

export const useFileActionsDisableSync = () => {
  const { showMessage, showErrorMessage } = useMessages()
  const router = useRouter()
  const { $gettext, $ngettext } = useGettext()

  const clientService = useClientService()
  const loadingService = useLoadingService()
  const configStore = useConfigStore()
  const { updateResourceField } = useResourcesStore()

  const handler = async ({ resources }: FileActionOptions<IncomingShareResource>) => {
    const errors: Error[] = []
    const triggerPromises: Promise<void>[] = []
    const triggerQueue = new PQueue({
      concurrency: configStore.options.concurrentRequests.resourceBatchActions
    })
    resources.forEach((resource) => {
      triggerPromises.push(
        triggerQueue.add(async () => {
          try {
            const { graphAuthenticated } = clientService
            await graphAuthenticated.driveItems.deleteDriveItem(resource.driveId, resource.id)

            updateResourceField<IncomingShareResource>({
              id: resource.id,
              field: 'syncEnabled',
              value: false
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
      if (isLocationSpacesActive(router, 'files-spaces-generic')) {
        showMessage({
          title: $ngettext(
            'Sync for the selected share was disabled successfully',
            'Sync for the selected shares was disabled successfully',
            resources.length
          )
        })
        router.push(createLocationShares('files-shares-with-me'))
      }

      return
    }

    showErrorMessage({
      title: $ngettext(
        'Failed to disable sync for the the selected share',
        'Failed to disable sync for the selected shares',
        resources.length
      ),
      errors
    })
  }

  const actions = computed((): FileAction<IncomingShareResource>[] => [
    {
      name: 'disable-sync',
      icon: 'spam-3',
      handler: (args) => loadingService.addTask(() => handler(args)),
      label: () => $gettext('Disable sync'),
      isVisible: ({ space, resources }) => {
        if (
          !isLocationSharesActive(router, 'files-shares-with-me') &&
          !isLocationSpacesActive(router, 'files-spaces-generic')
        ) {
          return false
        }
        if (resources.length === 0) {
          return false
        }

        if (
          isLocationSpacesActive(router, 'files-spaces-generic') &&
          (space?.driveType !== 'share' || resources.length > 1 || resources[0].path !== '/')
        ) {
          return false
        }

        return resources.some((resource) => resource.syncEnabled)
      },
      class: 'oc-files-actions-disable-sync-trigger'
    }
  ])

  return {
    actions
  }
}
