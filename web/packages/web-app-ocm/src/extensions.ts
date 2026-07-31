import {
  ApplicationInformation,
  FileActionOptions,
  useClientService,
  useConfigStore,
  useMessages,
  useUserStore,
  useWindowOpen
} from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { computed } from 'vue'
import { Extension } from '@ownclouders/web-pkg'
import { OCM_PROVIDER_ID, urlJoin } from '@ownclouders/web-client'

export const extensions = (appInfo: ApplicationInformation) => {
  const { showErrorMessage } = useMessages()
  const clientService = useClientService()
  const configStore = useConfigStore()
  const userStore = useUserStore()
  const { $gettext } = useGettext()
  const { openUrl } = useWindowOpen()

  const handler = async ({ resources }: FileActionOptions) => {
    const resource = resources[0]

    try {
      const params = new URLSearchParams()
      params.append('file', resource.id.toString())

      const { data } = await clientService.httpAuthenticated.post(
        '/sciencemesh/open-in-app',
        params
      )
      if (data.app_url) {
        openUrl(data.app_url)
      } else {
        showErrorMessage({
          title: $gettext('An error occurred'),
          desc: $gettext("Couldn't open remotely")
        })
      }
    } catch (error) {
      console.log(error)
      showErrorMessage({
        title: $gettext('An error occurred'),
        desc: $gettext("Couldn't open remotely"),
        errors: [error]
      })
    }
  }

  return computed<Extension[]>(() => [
    {
      id: 'com.github.owncloud.web.open-file-remote',
      type: 'action',
      extensionPointIds: ['global.files.context-actions'],
      action: {
        name: 'open-file-remote',
        category: 'actions',
        icon: 'remote-control',
        handler,
        label: () => $gettext('Open remotely'),
        isVisible: ({ resources }: FileActionOptions) => {
          if (!resources?.length) {
            return false
          }
          return (
            configStore.options.ocm.openRemotely &&
            resources[0]?.storageId?.startsWith(OCM_PROVIDER_ID)
          )
        },
        class: 'oc-files-actions-open-file-remote'
      }
    },
    ...((userStore.user && [
      {
        id: `app.${appInfo.id}.menuItem`,
        type: 'appMenuItem',
        label: () => appInfo.name,
        color: appInfo.color,
        icon: appInfo.icon,
        path: urlJoin(appInfo.id)
      }
    ]) ||
      [])
  ])
}
