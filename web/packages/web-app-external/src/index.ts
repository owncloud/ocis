import {
  AppWrapperRoute,
  defineWebApplication,
  ApplicationInformation,
  Extension
} from '@ownclouders/web-pkg'
import translations from '../l10n/translations.json'
import App from './App.vue'
import { useGettext } from 'vue3-gettext'
import { useAppProviderService } from '@ownclouders/web-pkg/src/composables/appProviderService'
import Redirect from './Redirect.vue'
import { useApplicationReadyStore } from './piniaStores'
import { computed } from 'vue'
import { useActionExtensionCreateFromTemplate } from './extensions'
import { useCreateFileHandler } from './composables'

export default defineWebApplication({
  setup(options: any) {
    const { $gettext } = useGettext()
    const appProviderService = useAppProviderService()
    const { createFileHandler } = useCreateFileHandler()

    if (!Object.hasOwn(options, 'appName')) {
      const appInfo: ApplicationInformation = {
        name: $gettext('External'),
        id: 'external'
      }
      const routes = [
        {
          // fallback route for old external-app URLs, in case someone made a bookmark. Can be removed with the next major release.
          name: 'apps',
          path: '/:driveAliasAndItem(.*)?',
          component: Redirect,
          meta: {
            authContext: 'hybrid',
            title: $gettext('Redirecting to external app'),
            patchCleanPath: true
          }
        }
      ]
      return {
        appInfo,
        routes,
        ready: () => {
          const applicationReadyStore = useApplicationReadyStore()
          applicationReadyStore.setReady()
        }
      }
    }

    const { appName } = options
    const appId = `external-${appName.toLowerCase()}`
    const mimeTypes = appProviderService.getMimeTypesByAppName(appName)
    const appInfo: ApplicationInformation = {
      name: appName,
      id: appId,
      extensions: mimeTypes.map((mimeType) => {
        const provider = mimeType.app_providers.find((provider) => provider.name === appName)
        return {
          extension: mimeType.ext,
          label: () => $gettext('Open in %{app}', { app: provider.name }),
          icon: provider.icon,
          name: provider.name,
          mimeType: mimeType.mime_type,
          secureView: provider.secure_view,
          routeName: `${appId}-apps`,
          hasPriority: mimeType.default_application === provider.name,
          ...(mimeType.allow_creation && { newFileMenu: { menuTitle: () => mimeType.name } }),
          createFileHandler
        }
      })
    }

    const routes = [
      {
        name: 'apps',
        path: '/:driveAliasAndItem(.*)?',
        component: AppWrapperRoute(App, {
          applicationId: appInfo.id
        }),
        meta: {
          authContext: 'hybrid',
          title: appName,
          patchCleanPath: true
        }
      }
    ]

    const actionCreateFromTemplate = useActionExtensionCreateFromTemplate(appInfo)
    const extensions = computed<Extension[]>(() => {
      return [actionCreateFromTemplate]
    })

    return {
      appInfo,
      routes,
      translations,
      extensions
    }
  }
})
