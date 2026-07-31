import { RouteComponents } from './router'
import { RouteLocationNamedRaw, RouteRecordRaw } from 'vue-router'
import { $gettext, createLocation, isLocationActiveDirector } from './utils'

type trashTypes = 'files-trash-generic' | 'files-trash-overview'

export const createLocationTrash = (name: trashTypes, location = {}): RouteLocationNamedRaw =>
  createLocation(name, location)

export const locationTrashGeneric = createLocationTrash('files-trash-generic')

export const locationTrashOverview = createLocationTrash('files-trash-overview')

export const isLocationTrashActive = isLocationActiveDirector<trashTypes>(
  locationTrashGeneric,
  locationTrashOverview
)

export const buildRoutes = (components: RouteComponents): RouteRecordRaw[] => [
  {
    path: '/trash',
    component: components.App,
    children: [
      {
        path: 'overview',
        name: locationTrashOverview.name,
        component: components.Trash.Overview,
        meta: {
          authContext: 'user',
          title: $gettext('Trash overview')
        }
      },
      {
        name: locationTrashGeneric.name,
        path: ':driveAliasAndItem(.*)?',
        component: components.Spaces.DriveResolver,
        meta: {
          authContext: 'user',
          patchCleanPath: true
        }
      }
    ]
  }
]
