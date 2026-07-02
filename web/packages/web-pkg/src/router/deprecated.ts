import {
  RouteRecordRaw,
  RouteLocationNamedRaw,
  Router,
  RouteRecordName,
  RouteLocation,
  RouteLocationRaw,
  RouteLocationPathRaw
} from 'vue-router'
import { createLocationSpaces } from './spaces'
import { createLocationShares } from './shares'
import { createLocationCommon } from './common'
import { createLocationPublic } from './public'
import { isLocationActive as isLocationActiveNoCompat } from './utils'
import { createLocationTrash } from './trash'
import { urlJoin } from '@ownclouders/web-client'
import { queryItemAsString } from './../composables/appDefaults'

/**
 * all route configs created by buildRoutes are deprecated,
 * this helper wraps a route config and warns the user that it will go away and redirect to the new location.
 *
 * @param routeConfig
 */
const deprecatedRedirect = (routeConfig: RouteRecordRaw): RouteRecordRaw => {
  return {
    meta: { ...routeConfig.meta, authContext: 'anonymous' }, // authContext belongs to the redirect target, not to the redirect itself.
    path: routeConfig.path,
    redirect: (to) => {
      const location = (routeConfig.redirect as (to: RouteLocation) => RouteLocationRaw)(to)

      console.warn(
        `route "${routeConfig.path}" is deprecated, use "${
          String((location as RouteLocationPathRaw).path) ||
          String((location as RouteLocationNamedRaw).name)
        }" instead.`
      )

      return location
    }
  }
}

/**
 * listed routes only exist to keep backwards compatibility intact,
 * all routes written in  a flat syntax to keep them readable.
 */
export const buildRoutes = (): RouteRecordRaw[] =>
  (
    [
      {
        path: '/list',
        redirect: (to) =>
          createLocationSpaces('files-spaces-generic', {
            ...to,
            params: { ...to.params, driveAliasAndItem: 'personal/home' }
          })
      },
      {
        path: '/list/all/:item(.*)',
        redirect: (to) =>
          createLocationSpaces('files-spaces-generic', {
            ...to,
            params: {
              ...to.params,
              driveAliasAndItem: urlJoin('personal/home', queryItemAsString(to.params.item), {
                leadingSlash: false
              })
            }
          })
      },
      {
        path: '/list/favorites',
        redirect: (to) => createLocationCommon('files-common-favorites', to)
      },
      {
        path: '/list/shared-with-me',
        redirect: (to) => createLocationShares('files-shares-with-me', to)
      },
      {
        path: '/list/shared-with-others',
        redirect: (to) => createLocationShares('files-shares-with-others', to)
      },
      {
        path: '/list/shared-via-link',
        redirect: (to) => createLocationShares('files-shares-via-link', to)
      },
      {
        path: '/trash-bin',
        redirect: (to) => createLocationTrash('files-trash-generic', to)
      },
      {
        path: '/public/list/:item(.*)',
        redirect: (to) => createLocationPublic('files-public-link', to)
      },
      {
        path: '/private-link/:fileId',
        redirect: (to) => ({ name: 'resolvePrivateLink', params: { fileId: to.params.fileId } })
      },
      {
        path: '/public-link/:token',
        redirect: (to) => ({ name: 'resolvePublicLink', params: { token: to.params.token } })
      }
    ] as RouteRecordRaw[]
  ).map(deprecatedRedirect)

/**
 * same as utils.isLocationActive with the difference that it remaps old route names to new ones and warns
 * @param router
 * @param comparatives
 */
export const isLocationActive = (
  router: Router,
  ...comparatives: [RouteLocationNamedRaw, ...RouteLocationNamedRaw[]]
): boolean => {
  const [first, ...rest] = comparatives.map((c) => {
    const newName: RouteRecordName = {
      'files-personal': createLocationSpaces('files-spaces-generic').name,
      'files-favorites': createLocationCommon('files-common-favorites').name,
      'files-shared-with-others': createLocationShares('files-shares-with-others').name,
      'files-shared-with-me': createLocationShares('files-shares-with-me').name,
      'files-trashbin	': createLocationTrash('files-trash-generic').name,
      'files-public-list': createLocationPublic('files-public-link').name
    }[c.name.toString()]

    if (newName) {
      console.warn(`route name "${name}" is deprecated, use "${newName.toString()}" instead.`)
    }

    return {
      ...c,
      ...(!!newName && { name: newName })
    }
  })

  return isLocationActiveNoCompat(router, first, ...rest)
}
