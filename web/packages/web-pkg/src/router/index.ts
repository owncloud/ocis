import { RouteRecordRaw } from 'vue-router'

import {
  buildRoutes as buildCommonRoutes,
  isLocationCommonActive,
  createLocationCommon
} from './common'
import { buildRoutes as buildDeprecatedRoutes, isLocationActive } from './deprecated'
import {
  buildRoutes as buildPublicRoutes,
  createLocationPublic,
  isLocationPublicActive,
  locationPublicLink,
  locationPublicUpload
} from './public'
import { RouteComponents } from './router'
import {
  buildRoutes as buildSharesRoutes,
  isLocationSharesActive,
  createLocationShares,
  locationSharesViaLink,
  locationSharesWithMe,
  locationSharesWithOthers,
  RouteShareTypes
} from './shares'
import {
  buildRoutes as buildSpacesRoutes,
  isLocationSpacesActive,
  createLocationSpaces,
  locationSpacesGeneric
} from './spaces'
import {
  buildRoutes as buildTrashRoutes,
  isLocationTrashActive,
  createLocationTrash
} from './trash'
import type { ActiveRouteDirectorFunc } from './utils'

const ROOT_ROUTE = {
  name: 'root',
  path: '/',
  redirect: (to) => createLocationSpaces('files-spaces-generic', to)
} as RouteRecordRaw

const buildRoutes = (components: RouteComponents): RouteRecordRaw[] => [
  ROOT_ROUTE,
  ...buildCommonRoutes(components),
  ...buildSharesRoutes(components),
  ...buildPublicRoutes(components),
  ...buildSpacesRoutes(components),
  ...buildTrashRoutes(components),
  ...buildDeprecatedRoutes()
]

export {
  createLocationCommon,
  createLocationShares,
  createLocationSpaces,
  createLocationPublic,
  isLocationCommonActive,
  isLocationSharesActive,
  isLocationSpacesActive,
  isLocationPublicActive,
  isLocationActive,
  isLocationTrashActive,
  createLocationTrash,
  locationPublicLink,
  locationPublicUpload,
  locationSpacesGeneric,
  locationSharesViaLink,
  locationSharesWithMe,
  locationSharesWithOthers,
  buildRoutes,
  ActiveRouteDirectorFunc
}

export type { RouteShareTypes }
