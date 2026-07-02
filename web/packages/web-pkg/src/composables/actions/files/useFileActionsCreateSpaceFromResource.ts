import { computed, unref } from 'vue'
import { useGettext } from 'vue3-gettext'
import { FileAction, FileActionOptions } from '../../actions'

import { useAbility } from '../../ability'
import { useClientService } from '../../clientService'
import { useRouter } from '../../router'
import {
  HIDDEN_FILE_EXTENSIONS,
  Resource,
  SpaceResource,
  isPersonalSpaceResource
} from '@ownclouders/web-client'
import { isLocationSpacesActive } from '../../../router'
import { useCreateSpace } from '../../spaces'
import { useSpaceHelpers } from '../../spaces'
import PQueue from 'p-queue'
import {
  useConfigStore,
  useMessages,
  useModals,
  useResourcesStore,
  useSpacesStore
} from '../../piniaStores'

export const useFileActionsCreateSpaceFromResource = () => {
  const { showMessage, showErrorMessage } = useMessages()
  const { can } = useAbility()
  const { $gettext, $ngettext } = useGettext()
  const { createSpace } = useCreateSpace()
  const { checkSpaceNameModalInput } = useSpaceHelpers()
  const clientService = useClientService()
  const router = useRouter()
  const hasCreatePermission = computed(() => can('create-all', 'Drive'))
  const { dispatchModal } = useModals()
  const configStore = useConfigStore()
  const spacesStore = useSpacesStore()
  const resourcesStore = useResourcesStore()

  const confirmAction = async ({
    spaceName,
    resources,
    space
  }: {
    spaceName: string
    resources: Resource[]
    space: SpaceResource
  }) => {
    const { webdav } = clientService
    const queue = new PQueue({
      concurrency: configStore.options.concurrentRequests.resourceBatchActions
    })
    const copyOps = []

    try {
      const createdSpace = await createSpace(spaceName)
      spacesStore.upsertSpace(createdSpace)

      if (resources.length === 1 && resources[0].isFolder) {
        //If a single folder is selected we copy it's content to the Space's root folder
        resources = (await webdav.listFiles(space, { path: resources[0].path })).children
      }

      for (const resource of resources) {
        copyOps.push(
          queue.add(() => webdav.copyFiles(space, resource, createdSpace, { path: resource.name }))
        )
      }

      await Promise.all(copyOps)
      resourcesStore.resetSelection()
      showMessage({ title: $gettext('Space was created successfully') })
    } catch (error) {
      console.error(error)
      showErrorMessage({
        title: $gettext('Creating space failed…'),
        errors: [error]
      })
    }
  }
  const handler = ({ resources, space }: FileActionOptions) => {
    dispatchModal({
      title: $ngettext(
        'Create Space from "%{resourceName}"',
        'Create Space from selection',
        resources.length,
        {
          resourceName: resources[0].name
        }
      ),
      message: $ngettext(
        'Create Space with the content of "%{resourceName}".',
        'Create Space with the selected files.',
        resources.length,
        {
          resourceName: resources[0].name
        }
      ),
      contextualHelperLabel: $gettext('The marked elements will be copied.'),
      contextualHelperData: {
        title: $gettext('Restrictions'),
        text: $gettext('Shares, versions and tags will not be copied.')
      },
      confirmText: $gettext('Create'),
      hasInput: true,
      inputLabel: $gettext('Space name'),
      onInput: checkSpaceNameModalInput,
      onConfirm: (spaceName: string) => confirmAction({ spaceName, space, resources })
    })
  }

  const actions = computed((): FileAction[] => {
    return [
      {
        name: 'create-space-from-resource',
        icon: 'function',
        handler,
        label: () => {
          return $gettext('Create Space from selection')
        },
        isVisible: ({ resources, space }) => {
          if (!resources.length) {
            return false
          }

          if (!unref(hasCreatePermission)) {
            return false
          }

          if (
            !isLocationSpacesActive(router, 'files-spaces-generic') ||
            !isPersonalSpaceResource(space)
          ) {
            return false
          }

          if (HIDDEN_FILE_EXTENSIONS.includes(resources[0].extension)) {
            return false
          }

          return true
        },
        class: 'oc-files-actions-create-space-from-resource-trigger'
      }
    ]
  })

  return {
    actions
  }
}
