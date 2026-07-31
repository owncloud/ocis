import { computed } from 'vue'
import { useGettext } from 'vue3-gettext'
import translations from '../l10n/translations.json'
import { AppMenuItemExtension, AppWrapperRoute, defineWebApplication } from '@ownclouders/web-pkg'

export default defineWebApplication({
  setup() {
    const { $gettext } = useGettext()

    const appId = 'epub-reader'

    const routes = [
      {
        path: '/',
        component: () => import('./views/LibraryView.vue'),
        name: 'epub-library',
        meta: {
          authContext: 'user',
          title: $gettext('Library')
        }
      },
      {
        path: '/read/:driveAliasAndItem(.*)?',
        component: async () => {
          // lazy loading to avoid loading the epubjs package on page load
          const EpubReader = (await import('./App.vue')).default
          return AppWrapperRoute(EpubReader, {
            applicationId: appId,
            fileContentOptions: {
              responseType: 'blob'
            }
          })
        },
        name: 'epub-reader',
        meta: {
          authContext: 'hybrid',
          title: $gettext('Epub Reader'),
          patchCleanPath: true
        }
      }
    ]

    const appInfo = {
      name: $gettext('Epub Reader'),
      id: appId,
      icon: 'book-read',
      color: '#4f46e5',
      extensions: [
        {
          extension: 'epub',
          routeName: 'epub-reader'
        }
      ]
    }

    const extensions = computed<AppMenuItemExtension[]>(() => [
      {
        id: `app.${appId}.menuItem`,
        type: 'appMenuItem',
        label: () => $gettext('Library'),
        color: appInfo.color,
        icon: 'book',
        priority: 50,
        path: `/${appId}`
      }
    ])

    return {
      appInfo,
      routes,
      translations,
      extensions
    }
  }
})
