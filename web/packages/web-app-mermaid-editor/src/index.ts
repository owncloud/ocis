import { useGettext } from 'vue3-gettext'
import translations from '../l10n/translations.json'
import MermaidEditor from './App.vue'
import {
  AppWrapperRoute,
  ApplicationFileExtension,
  ApplicationInformation,
  defineWebApplication
} from '@ownclouders/web-pkg'

export default defineWebApplication({
  setup() {
    const { $gettext } = useGettext()

    const appId = 'mermaid-editor'

    // `newFileMenu` on the `mmd` entry adds a "New > Mermaid diagram" item to the
    // Files create menu. The extensions also declare which file types open in this app.
    const mermaidFileLabel = () => $gettext('Mermaid diagram')
    const fileExtensions: ApplicationFileExtension[] = [
      {
        extension: 'mmd',
        label: mermaidFileLabel,
        newFileMenu: {
          menuTitle: mermaidFileLabel
        }
      },
      {
        extension: 'mermaid',
        label: mermaidFileLabel
      }
    ]

    const routes = [
      {
        path: '/:driveAliasAndItem(.*)?',
        component: AppWrapperRoute(MermaidEditor, {
          applicationId: appId
        }),
        name: 'mermaid-editor',
        meta: {
          authContext: 'hybrid',
          title: $gettext('Mermaid Editor'),
          patchCleanPath: true
        }
      }
    ]

    const appInfo: ApplicationInformation = {
      name: $gettext('Mermaid Editor'),
      id: appId,
      icon: 'flow-chart',
      color: '#2E8B57',
      defaultExtension: 'mmd',
      meta: {
        fileSizeLimit: 2000000
      },
      extensions: fileExtensions
    }

    return {
      appInfo,
      routes,
      translations
    }
  }
})
