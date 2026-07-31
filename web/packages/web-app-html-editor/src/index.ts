import { useGettext } from 'vue3-gettext'
import translations from '../l10n/translations.json'
import HtmlEditor from './App.vue'
import {
  AppWrapperRoute,
  ApplicationFileExtension,
  ApplicationInformation,
  defineWebApplication
} from '@ownclouders/web-pkg'
import { PREVIEW_SIZE_LIMIT } from './helpers/preview'

// AppWrapper refuses to open files larger than this (bytes). It must stay strictly
// above PREVIEW_SIZE_LIMIT so the editor still opens "large-but-openable" files
// whose live preview is merely paused — otherwise the two limits would silently
// disagree and such files would be rejected before the pause logic ever runs.
const FILE_SIZE_LIMIT = 2_000_000
if (PREVIEW_SIZE_LIMIT >= FILE_SIZE_LIMIT) {
  throw new Error(
    `html-editor: PREVIEW_SIZE_LIMIT (${PREVIEW_SIZE_LIMIT}) must be smaller than ` +
      `FILE_SIZE_LIMIT (${FILE_SIZE_LIMIT})`
  )
}

export default defineWebApplication({
  setup() {
    const { $gettext } = useGettext()

    const appId = 'html-editor'

    // `newFileMenu` on the html entry adds a "New > HTML file" item to the Files
    // create menu. The extensions also declare which file types open in this app.
    const htmlFileLabel = () => $gettext('HTML file')
    const fileExtensions: ApplicationFileExtension[] = [
      {
        extension: 'html',
        label: htmlFileLabel,
        newFileMenu: {
          menuTitle: htmlFileLabel
        }
      },
      {
        extension: 'htm',
        label: htmlFileLabel
      },
      {
        extension: 'xhtml',
        label: () => $gettext('XHTML file')
      }
    ]

    const routes = [
      {
        path: '/:driveAliasAndItem(.*)?',
        component: AppWrapperRoute(HtmlEditor, {
          applicationId: appId
        }),
        name: 'html-editor',
        meta: {
          authContext: 'hybrid',
          title: $gettext('HTML Editor'),
          patchCleanPath: true
        }
      }
    ]

    const appInfo: ApplicationInformation = {
      name: $gettext('HTML Editor'),
      id: appId,
      icon: 'file-code',
      color: '#e34c26',
      defaultExtension: 'html',
      meta: {
        fileSizeLimit: FILE_SIZE_LIMIT
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
