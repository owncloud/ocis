import { RouteComponents } from './router'
import { RouteLocationNamedRaw, RouteRecordRaw } from 'vue-router'
import { createLocation, isLocationActiveDirector, $gettext } from './utils'

type shareTypes = 'files-public-link' | 'files-public-upload'

export const createLocationPublic = (name: shareTypes, location = {}): RouteLocationNamedRaw =>
  createLocation(name, location)

export const locationPublicLink = createLocationPublic('files-public-link')
export const locationPublicUpload = createLocationPublic('files-public-upload')

export const isLocationPublicActive = isLocationActiveDirector<shareTypes>(
  locationPublicLink,
  locationPublicUpload
)

export const buildRoutes = (components: RouteComponents): RouteRecordRaw[] => [
  {
    path: '/link',
    component: components.App,
    meta: {
      auth: false
    },
    children: [
      {
        name: locationPublicLink.name,
        path: ':driveAliasAndItem(.*)?',
        component: components.Spaces.DriveResolver,
        meta: {
          authContext: 'publicLink',
          patchCleanPath: true
        }
      }
    ]
  },
  {
    path: '/upload',
    component: components.App,
    meta: {
      auth: false,
      isUploadSnackbarHidden: true
    },
    children: [
      {
        name: locationPublicUpload.name,
        path: ':token?',
        component: components.FilesDrop,
        meta: {
          authContext: 'publicLink',
          title: $gettext('Public file upload')
        }
      }
    ]
  }
]
