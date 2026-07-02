import FileDetails from '../../components/SideBar/Details/FileDetails.vue'
import FileDetailsMultiple from '../../components/SideBar/Details/FileDetailsMultiple.vue'
import ExifPanel from '../../components/SideBar/Exif/ExifPanel.vue'
import FileActions from '../../components/SideBar/Actions/FileActions.vue'
import FileVersions from '../../components/SideBar/Versions/FileVersions.vue'
import SharesPanel from '../../components/SideBar/Shares/SharesPanel.vue'
import NoSelection from '../../components/SideBar/NoSelection.vue'
import TrashNoSelection from '../../components/SideBar/TrashNoSelection.vue'
import SpaceActions from '../../components/SideBar/Actions/SpaceActions.vue'
import ActivitiesPanel from '../../components/SideBar/ActivitiesPanel.vue'
import {
  SpaceDetails,
  SpaceDetailsMultiple,
  SpaceNoSelection,
  isLocationTrashActive,
  isLocationSpacesActive,
  isLocationSharesActive,
  useRouter,
  SidebarPanelExtension,
  useIsFilesAppActive,
  useGetMatchingSpace,
  useCapabilityStore,
  useCanListShares,
  useCanListVersions
} from '@ownclouders/web-pkg'
import { isProjectSpaceResource, SpaceResource } from '@ownclouders/web-client'
import { Resource } from '@ownclouders/web-client'
import { useGettext } from 'vue3-gettext'
import { unref } from 'vue'
import { fileSideBarExtensionPoint } from '../../extensionPoints'
import AudioMetaPanel from '../../components/SideBar/Audio/AudioMetaPanel.vue'
import { isEmpty } from 'lodash-es'

export const useSideBarPanels = (): SidebarPanelExtension<SpaceResource, Resource, Resource>[] => {
  const router = useRouter()
  const capabilityStore = useCapabilityStore()
  const { $gettext } = useGettext()
  const isFilesAppActive = useIsFilesAppActive()
  const { isPersonalSpaceRoot } = useGetMatchingSpace()
  const { canListShares } = useCanListShares()
  const { canListVersions } = useCanListVersions()

  return [
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.no-selection',
      type: 'sidebarPanel',
      extensionPointIds: [fileSideBarExtensionPoint.id],
      panel: {
        name: 'no-selection',
        icon: 'questionnaire-line',
        title: () => $gettext('Details'),
        component: NoSelection,
        isRoot: () => true,
        isVisible: ({ parent, items }) => {
          if (isLocationTrashActive(router, 'files-trash-overview')) {
            // trash overview has its own "no selection" panel
            return false
          }
          if (isLocationSpacesActive(router, 'files-spaces-projects')) {
            // project spaces overview has its own "no selection" panel
            return false
          }
          if (items?.length > 0) {
            return false
          }
          // empty selection in a project space root shows a "details" panel for the space itself
          return !isProjectSpaceResource(parent)
        }
      }
    },
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.trash-no-selection',
      type: 'sidebarPanel',
      extensionPointIds: [fileSideBarExtensionPoint.id],
      panel: {
        name: 'no-selection',
        icon: 'questionnaire-line',
        title: () => $gettext('Details'),
        component: TrashNoSelection,
        isRoot: () => true,
        isVisible: () => {
          return isLocationTrashActive(router, 'files-trash-overview')
        }
      }
    },
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.details-single-selection',
      type: 'sidebarPanel',
      extensionPointIds: [fileSideBarExtensionPoint.id],
      panel: {
        name: 'details',
        icon: 'questionnaire-line',
        title: () => $gettext('Details'),
        component: FileDetails,
        componentAttrs: ({ items }) => ({
          previewEnabled: unref(isFilesAppActive),
          tagsEnabled:
            !isPersonalSpaceRoot(items[0]) && !isLocationTrashActive(router, 'files-trash-generic'),
          versionsEnabled: !isLocationTrashActive(router, 'files-trash-generic')
        }),
        isRoot: () => true,
        isVisible: ({ items }) => {
          if (items?.length !== 1) {
            return false
          }
          // project spaces have their own "details" panel
          return !isProjectSpaceResource(items[0])
        }
      }
    },
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.details-multi-selection',
      type: 'sidebarPanel',
      extensionPointIds: [fileSideBarExtensionPoint.id],
      panel: {
        name: 'details-multiple',
        icon: 'questionnaire-line',
        title: () => $gettext('Details'),
        component: FileDetailsMultiple,
        componentAttrs: () => ({
          get showSpaceCount() {
            return (
              !isLocationSpacesActive(router, 'files-spaces-generic') &&
              !isLocationSharesActive(router, 'files-shares-with-me') &&
              !isLocationTrashActive(router, 'files-trash-generic')
            )
          }
        }),
        isRoot: () => true,
        isVisible: ({ items }) => {
          if (isLocationSpacesActive(router, 'files-spaces-projects')) {
            // project spaces overview has its own "no selection" panel
            return false
          }
          return items?.length > 1
        }
      }
    },
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.exif',
      type: 'sidebarPanel',
      extensionPointIds: ['global.files.sidebar'],
      panel: {
        name: 'exif',
        icon: 'image',
        title: () => $gettext('Image Info'),
        component: ExifPanel,
        isVisible: ({ items }) => {
          if (items?.length !== 1) {
            return false
          }
          const item = items[0]
          if (item.type !== 'file') {
            return false
          }
          return !isEmpty(item.image) || !isEmpty(item.photo)
        }
      }
    },
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.audio-meta',
      type: 'sidebarPanel',
      extensionPointIds: ['global.files.sidebar'],
      panel: {
        name: 'audio-meta',
        icon: 'music',
        title: () => $gettext('Audio Info'),
        component: AudioMetaPanel,
        isVisible: ({ items }) => {
          if (items?.length !== 1) {
            return false
          }
          const item = items[0]
          if (item.type !== 'file') {
            return false
          }
          return !isEmpty(item.audio)
        }
      }
    },
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.actions',
      type: 'sidebarPanel',
      extensionPointIds: [fileSideBarExtensionPoint.id],
      panel: {
        name: 'actions',
        icon: 'play-circle',
        iconFillType: 'line',
        title: () => $gettext('Actions'),
        component: FileActions,
        isRoot: () => false,
        isVisible: ({ items }) => {
          if (items?.length !== 1) {
            return false
          }
          if (isPersonalSpaceRoot(items[0])) {
            // actions panel is not available on the personal space root for now ;-)
            return false
          }
          // project spaces have their own "actions" panel
          return !isProjectSpaceResource(items[0])
        }
      }
    },
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.sharing',
      type: 'sidebarPanel',
      extensionPointIds: [fileSideBarExtensionPoint.id],
      panel: {
        name: 'sharing',
        icon: 'user-add',
        iconFillType: 'line',
        title: () => $gettext('Shares'),
        component: SharesPanel,
        componentAttrs: () => ({
          showSpaceMembers: false,
          get showLinks() {
            return capabilityStore.sharingPublicEnabled
          }
        }),
        isVisible: ({ items, root }) => {
          if (items?.length !== 1) {
            return false
          }
          if (isProjectSpaceResource(items[0])) {
            // project space roots have their own "sharing" panel (= space members)
            return false
          }
          return canListShares({ space: root, resource: items[0] })
        }
      }
    },
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.versions',
      type: 'sidebarPanel',
      extensionPointIds: [fileSideBarExtensionPoint.id],
      panel: {
        name: 'versions',
        icon: 'git-branch',
        title: () => $gettext('Versions'),
        component: FileVersions,
        componentAttrs: () => ({
          isReadOnly: !unref(isFilesAppActive)
        }),
        isVisible: ({ items, root }) => {
          if (items?.length !== 1) {
            return false
          }
          return canListVersions({ space: root, resource: items[0] })
        }
      }
    },
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.projects.no-selection',
      type: 'sidebarPanel',
      extensionPointIds: [fileSideBarExtensionPoint.id],
      panel: {
        name: 'no-selection',
        icon: 'questionnaire-line',
        title: () => $gettext('Details'),
        component: SpaceNoSelection,
        isRoot: () => true,
        isVisible: ({ items }) => {
          if (!isLocationSpacesActive(router, 'files-spaces-projects')) {
            // only for project spaces overview
            return false
          }
          return items?.length === 0
        }
      }
    },
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.projects.details-single-selection',
      type: 'sidebarPanel',
      extensionPointIds: [fileSideBarExtensionPoint.id],
      panel: {
        name: 'details-space',
        icon: 'questionnaire-line',
        title: () => $gettext('Details'),
        component: SpaceDetails,
        isRoot: () => true,
        isVisible: ({ items }) => {
          return items?.length === 1 && isProjectSpaceResource(items[0])
        }
      }
    },
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.projects.details-multi-selection',
      type: 'sidebarPanel',
      extensionPointIds: [fileSideBarExtensionPoint.id],
      panel: {
        name: 'details-space-multiple',
        icon: 'questionnaire-line',
        title: () => $gettext('Details'),
        component: SpaceDetailsMultiple,
        componentAttrs: ({ items }) => ({
          selectedSpaces: items
        }),
        isRoot: () => true,
        isVisible: ({ items }) => {
          return items?.length > 1 && isLocationSpacesActive(router, 'files-spaces-projects')
        }
      }
    },
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.projects.actions',
      type: 'sidebarPanel',
      extensionPointIds: [fileSideBarExtensionPoint.id],
      panel: {
        name: 'space-actions',
        icon: 'play-circle',
        iconFillType: 'line',
        title: () => $gettext('Actions'),
        component: SpaceActions,
        isVisible: ({ items }) => {
          if (items?.length !== 1) {
            return false
          }
          if (!isProjectSpaceResource(items[0])) {
            return false
          }
          if (
            !isLocationSpacesActive(router, 'files-spaces-projects') &&
            !isLocationSpacesActive(router, 'files-spaces-generic')
          ) {
            return false
          }
          return true
        }
      }
    },
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.projects.sharing',
      type: 'sidebarPanel',
      extensionPointIds: [fileSideBarExtensionPoint.id],
      panel: {
        name: 'space-share',
        icon: 'group',
        title: () => $gettext('Members'),
        component: SharesPanel,
        componentAttrs: () => ({
          showSpaceMembers: true,
          get showLinks() {
            return capabilityStore.sharingPublicEnabled
          }
        }),
        isVisible: ({ items }) => {
          return items?.length === 1 && isProjectSpaceResource(items[0]) && !items[0].disabled
        }
      }
    },
    {
      id: 'com.github.owncloud.web.files.sidebar-panel.activities',
      type: 'sidebarPanel',
      extensionPointIds: [fileSideBarExtensionPoint.id],
      panel: {
        name: 'activities',
        icon: 'pulse',
        title: () => $gettext('Activities'),
        component: ActivitiesPanel,
        isVisible: ({ items }) => {
          if (items?.length !== 1) {
            return false
          }
          if (isLocationTrashActive(router, 'files-trash-generic')) {
            return false
          }
          return true
        }
      }
    }
  ]
}
