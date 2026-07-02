import { RouteComponents } from './router'
import { RouteLocationNamedRaw, RouteRecordRaw } from 'vue-router'
import { createLocation, isLocationActiveDirector, $gettext } from './utils'

export type RouteShareTypes =
  | 'files-shares'
  | 'files-shares-with-me'
  | 'files-shares-with-others'
  | 'files-shares-via-link'

export const createLocationShares = (name: RouteShareTypes, location = {}): RouteLocationNamedRaw =>
  createLocation(name, location)

export const locationShares = createLocationShares('files-shares')
export const locationSharesWithMe = createLocationShares('files-shares-with-me')
export const locationSharesWithOthers = createLocationShares('files-shares-with-others')
export const locationSharesViaLink = createLocationShares('files-shares-via-link')

export const isLocationSharesActive = isLocationActiveDirector<RouteShareTypes>(
  locationSharesWithMe,
  locationSharesWithOthers,
  locationSharesViaLink
)

export const buildRoutes = (components: RouteComponents): RouteRecordRaw[] => [
  {
    name: locationShares.name,
    path: '/shares',
    component: components.App,
    redirect: locationSharesWithMe,
    children: [
      {
        name: locationSharesWithMe.name,
        path: 'with-me',
        component: components.Shares.SharedWithMe,
        meta: {
          authContext: 'user',
          title: $gettext('Files shared with me')
        }
      },
      {
        name: locationSharesWithOthers.name,
        path: 'with-others',
        component: components.Shares.SharedWithOthers,
        meta: {
          authContext: 'user',
          title: $gettext('Files shared with others')
        }
      },
      {
        name: locationSharesViaLink.name,
        path: 'via-link',
        component: components.Shares.SharedViaLink,
        meta: {
          authContext: 'user',
          title: $gettext('Files shared via link')
        }
      }
    ]
  }
]
