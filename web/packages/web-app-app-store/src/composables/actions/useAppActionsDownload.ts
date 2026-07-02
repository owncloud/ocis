import { Action, triggerDownloadWithFilename } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { App, AppVersion } from '../../types'

export type AppActionOptions = {
  app: App
  version?: AppVersion
}

export const useAppActionsDownload = () => {
  const { $gettext } = useGettext()

  const downloadAppAction: Action<AppActionOptions> = {
    name: 'download-app',
    icon: 'download',
    label: () => {
      return $gettext('Download')
    },
    handler: (options?) => {
      const version = options.version || options.app.mostRecentVersion
      const filename = version.filename || version.url.split('/').pop()
      triggerDownloadWithFilename(version.url, filename)
    },
    isVisible: () => {
      return true
    },
    appearance: 'outline'
  }

  return {
    downloadAppAction
  }
}
