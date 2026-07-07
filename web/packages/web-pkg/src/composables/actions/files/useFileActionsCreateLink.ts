import { computed, unref } from 'vue'
import { useGettext } from 'vue3-gettext'
import { FileAction, FileActionOptions } from '../../actions'
import CreateLinkModal from '../../../components/CreateLinkModal.vue'
import { useAbility } from '../../ability'
import { LinkShare, isProjectSpaceResource } from '@ownclouders/web-client'
import { useLinkTypes } from '../../links'
import { useLoadingService } from '../../loadingService'
import {
  useMessages,
  useModals,
  useUserStore,
  useCapabilityStore,
  useSharesStore
} from '../../piniaStores'
import { useClipboard } from '../../clipboard'
import { useClientService } from '../../clientService'

export const useFileActionsCreateLink = ({
  enforceModal = false
}: {
  enforceModal?: boolean
} = {}) => {
  const clientService = useClientService()
  const userStore = useUserStore()
  const { showMessage, showErrorMessage } = useMessages()
  const { $gettext, $ngettext } = useGettext()
  const capabilityStore = useCapabilityStore()
  const ability = useAbility()
  const loadingService = useLoadingService()
  const { defaultLinkType } = useLinkTypes()
  const { addLink } = useSharesStore()
  const { dispatchModal } = useModals()
  const { copyToClipboard } = useClipboard()

  const proceedResult = async ({
    result,
    password,
    options = {}
  }: {
    result: PromiseSettledResult<LinkShare>[]
    password?: string
    options?: { copyPassword?: boolean }
  }) => {
    const succeeded = result.filter(
      (val): val is PromiseFulfilledResult<LinkShare> => val.status === 'fulfilled'
    )

    if (succeeded.length) {
      let successMessage = $gettext('Link has been created successfully')

      if (result.length === 1) {
        // Only copy to clipboard if the user tries to create one single link
        try {
          const copyToClipboardText = options.copyPassword
            ? $gettext(
                '%{link} Password:%{password}',
                {
                  link: succeeded[0].value.webUrl,
                  password
                },
                true
              )
            : succeeded[0].value.webUrl

          await copyToClipboard(copyToClipboardText)
          successMessage = $gettext('The link has been copied to your clipboard.')
        } catch (e) {
          console.warn('Unable to copy link to clipboard', e)
        }
      }

      showMessage({
        title: $ngettext(successMessage, 'Links have been created successfully.', succeeded.length)
      })
    }

    const failed = result.filter(({ status }) => status === 'rejected')
    if (failed.length) {
      showErrorMessage({
        errors: (failed as PromiseRejectedResult[]).map(({ reason }) => reason),
        title: $ngettext('Failed to create link', 'Failed to create links', failed.length)
      })
    }
  }

  const handler = async ({ space, resources }: FileActionOptions) => {
    const passwordEnforced = capabilityStore.sharingPublicPasswordEnforcedFor.read_only === true
    if (enforceModal || passwordEnforced) {
      dispatchModal({
        title: $ngettext(
          'Copy link for "%{resourceName}"',
          'Copy links for the selected items',
          resources.length,
          { resourceName: resources[0].name }
        ),
        customComponent: CreateLinkModal,
        customComponentAttrs: () => ({
          space,
          resources,
          callbackFn: proceedResult
        }),
        hideActions: true
      })
      return
    }

    const promises = resources.map((resource) =>
      addLink({
        clientService,
        space,
        resource,
        options: {
          '@libre.graph.quickLink': false,
          displayName: $gettext('Unnamed link'),
          type: unref(defaultLinkType)
        }
      })
    )
    const result = await loadingService.addTask(() => Promise.allSettled<LinkShare>(promises))

    proceedResult({ result })
  }

  const isVisible = ({ resources }: FileActionOptions) => {
    if (!resources.length) {
      return false
    }

    for (const resource of resources) {
      if (!resource.canShare({ user: userStore.user, ability })) {
        return false
      }

      if (isProjectSpaceResource(resource) && resource.disabled) {
        return false
      }
    }

    return true
  }

  const actions = computed((): FileAction[] => {
    return [
      {
        name: 'create-links',
        icon: 'link',
        handler,
        label: () => {
          return $gettext('Create links')
        },
        isVisible,
        class: 'oc-files-actions-create-links'
      }
    ]
  })

  return {
    actions
  }
}
