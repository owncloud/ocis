import {
  ApplicationInformation,
  Extension,
  useCapabilityStore,
  useConfigStore,
  useRouter,
  useSearch,
  useUserStore
} from '@ownclouders/web-pkg'
import { computed } from 'vue'
import { SDKSearch } from './search'
import { useSideBarPanels } from './composables/extensions/useFileSideBars'
import { useFolderViews } from './composables/extensions/useFolderViews'
import { useFileActions } from './composables/extensions/useFileActions'
import { urlJoin } from '@ownclouders/web-client'

export const extensions = (appInfo: ApplicationInformation) => {
  const capabilityStore = useCapabilityStore()
  const configStore = useConfigStore()
  const userStore = useUserStore()
  const router = useRouter()
  const { search: searchFunction } = useSearch()

  const fileActionExtensions = useFileActions()
  const folderViewExtensions = useFolderViews()
  const sideBarPanelExtensions = useSideBarPanels()

  return computed<Extension[]>(() => [
    ...fileActionExtensions,
    ...folderViewExtensions,
    ...sideBarPanelExtensions,
    {
      id: 'com.github.owncloud.web.files.search',
      extensionPointIds: ['app.search.provider'],
      type: 'search',
      searchProvider: new SDKSearch(capabilityStore, router, searchFunction, configStore)
    },
    ...((userStore.user && [
      {
        id: `app.${appInfo.id}.menuItem`,
        type: 'appMenuItem',
        label: () => appInfo.name,
        color: appInfo.color,
        icon: appInfo.icon,
        priority: 10,
        path: urlJoin(appInfo.id)
      }
    ]) ||
      [])
  ])
}
