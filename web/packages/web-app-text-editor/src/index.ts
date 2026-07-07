import { useGettext } from 'vue3-gettext'
import translations from '../l10n/translations.json'
import TextEditor from './App.vue'
import {
  AppMenuItemExtension,
  AppWrapperRoute,
  ApplicationFileExtension,
  ApplicationInformation,
  defineWebApplication,
  useOpenEmptyEditor,
  useUserStore
} from '@ownclouders/web-pkg'
import { computed } from 'vue'
import { urlJoin } from '@ownclouders/web-client'

export default defineWebApplication({
  setup({ applicationConfig }) {
    const { $gettext } = useGettext()
    const userStore = useUserStore()
    const { openEmptyEditor } = useOpenEmptyEditor()

    const appId = 'text-editor'

    const fileExtensions = () => {
      const extensions: ApplicationFileExtension[] = [
        {
          extension: 'txt',
          label: () => $gettext('Plain text file')
        },
        {
          extension: 'md',
          label: () => $gettext('Markdown file')
        },
        {
          extension: 'markdown',
          label: () => $gettext('Markdown file')
        },
        {
          extension: 'js',
          label: () => $gettext('JavaScript file')
        },
        {
          extension: 'json',
          label: () => $gettext('JSON file')
        },
        {
          extension: 'xml',
          label: () => $gettext('XML file')
        },
        {
          extension: 'py',
          label: () => $gettext('Python file')
        },
        {
          extension: 'php',
          label: () => $gettext('PHP file')
        },
        {
          extension: 'yaml',
          label: () => $gettext('YAML file')
        },
        {
          extension: 'log',
          label: () => $gettext('Log file')
        },
        {
          extension: 'conf',
          label: () => $gettext('Configuration file')
        }
      ]

      const config = applicationConfig || {}
      extensions.push(...(config.extraExtensions || []).map((ext: string) => ({ extension: ext })))

      let primaryExtensions: string[] = config.primaryExtensions || ['txt', 'md']

      if (typeof primaryExtensions === 'string') {
        primaryExtensions = [primaryExtensions]
      }

      return extensions.reduce<ApplicationFileExtension[]>((acc, extensionItem) => {
        const isPrimary = primaryExtensions.includes(extensionItem.extension)
        if (isPrimary) {
          extensionItem.newFileMenu = {
            menuTitle() {
              if (typeof extensionItem.label === 'function') {
                return extensionItem.label()
              }
              return extensionItem.label
            }
          }
        }
        acc.push(extensionItem)
        return acc
      }, [])
    }

    const routes = [
      {
        path: '/:driveAliasAndItem(.*)?',
        component: AppWrapperRoute(TextEditor, {
          applicationId: appId
        }),
        name: 'text-editor',
        meta: {
          authContext: 'hybrid',
          title: $gettext('Text Editor'),
          patchCleanPath: true
        }
      }
    ]

    const appInfo: ApplicationInformation = {
      name: $gettext('Text Editor'),
      id: appId,
      icon: 'file-text',
      color: '#0D856F',
      defaultExtension: 'txt',
      meta: {
        fileSizeLimit: 2000000
      },
      extensions: fileExtensions().map((extensionItem) => {
        return {
          extension: extensionItem.extension,
          ...(Object.prototype.hasOwnProperty.call(extensionItem, 'newFileMenu') && {
            newFileMenu: extensionItem.newFileMenu
          })
        }
      })
    }

    const menuItems = computed<AppMenuItemExtension[]>(() => {
      const items: AppMenuItemExtension[] = []

      if (userStore.user) {
        items.push({
          id: `app.${appInfo.id}.menuItem`,
          type: 'appMenuItem',
          label: () => appInfo.name,
          color: appInfo.color,
          icon: appInfo.icon,
          priority: 20,
          path: urlJoin(appInfo.id),
          handler: () => openEmptyEditor(appInfo.id, appInfo.defaultExtension)
        })
      }

      return items
    })

    return {
      appInfo,
      routes,
      translations,
      extensions: menuItems
    }
  }
})
