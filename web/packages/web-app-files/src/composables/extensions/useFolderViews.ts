import { FolderViewExtension, ResourceTable, ResourceTiles } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import {
  folderViewsFavoritesExtensionPoint,
  folderViewsFolderExtensionPoint,
  folderViewsProjectSpacesExtensionPoint
} from '../../extensionPoints'

export const useFolderViews = (): FolderViewExtension[] => {
  const { $gettext } = useGettext()

  return [
    {
      id: 'com.github.owncloud.web.files.folder-view.resource-table',
      type: 'folderView',
      extensionPointIds: [
        folderViewsFolderExtensionPoint.id,
        folderViewsProjectSpacesExtensionPoint.id,
        folderViewsFavoritesExtensionPoint.id
      ],
      folderView: {
        name: 'resource-table',
        label: $gettext('Switch to default table view'),
        icon: {
          name: 'menu-line',
          fillType: 'none'
        },
        component: ResourceTable
      }
    },
    {
      id: 'com.github.owncloud.web.files.folder-view.resource-table-condensed',
      type: 'folderView',
      extensionPointIds: [folderViewsFolderExtensionPoint.id],
      folderView: {
        name: 'resource-table-condensed',
        label: $gettext('Switch to condensed table view'),
        icon: {
          name: 'menu-line-condensed',
          fillType: 'none'
        },
        component: ResourceTable
      }
    },
    {
      id: 'com.github.owncloud.web.files.folder-view.resource-tiles',
      type: 'folderView',
      extensionPointIds: [
        folderViewsFolderExtensionPoint.id,
        folderViewsProjectSpacesExtensionPoint.id,
        folderViewsFavoritesExtensionPoint.id
      ],
      folderView: {
        name: 'resource-tiles',
        label: $gettext('Switch to tiles view'),
        icon: {
          name: 'apps-2',
          fillType: 'line'
        },
        component: ResourceTiles
      }
    }
  ]
}
