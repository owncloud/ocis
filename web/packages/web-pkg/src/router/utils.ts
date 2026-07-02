import { Router, RouteLocationNamedRaw } from 'vue-router'
import merge from 'lodash-es/merge'
import { unref } from 'vue'

export interface ActiveRouteDirectorFunc<T extends string> {
  (router: Router, ...comparatives: T[]): boolean
}

/**
 * helper function to find out if comparative route location is active or not.
 * it uses vue router resolve to do so.
 *
 * @param router
 * @param comparatives
 */
export const isLocationActive = (
  router: Router,
  ...comparatives: [RouteLocationNamedRaw, ...RouteLocationNamedRaw[]]
): boolean => {
  // FIXME: router.resolve cleans the path. we don't need it, if we can rely on
  // router.currentRoute to not have slashs encoded for paths
  const { href: currentHref } = router.resolve(unref(router.currentRoute))
  return comparatives
    .map((comparative) => {
      const { href: comparativeHref } = router.resolve({
        ...comparative
        // ...(comparative.name && { name: resolveRouteName(comparative.name) })
      })

      /**
       * Href might be '/' or '#/' if router is not able to resolve the proper path.
       * This happens if the we don't pass a param which is defined in the route configuration, for example:
       * path: user/:id
       *
       * This implies that the comparative route is not active
       **/
      if (comparativeHref === '/' || comparativeHref === '#/') {
        return false
      }
      return currentHref.startsWith(comparativeHref)
    })
    .some(Boolean)
}

/**
 * wraps isLocationActive to be used as a closure,
 * the resulting closure then can be used to check a location against the defined set of director locations
 *
 * @param defaultComparatives
 */
export const isLocationActiveDirector = <T extends string>(
  ...defaultComparatives: [RouteLocationNamedRaw, ...RouteLocationNamedRaw[]]
): ActiveRouteDirectorFunc<T> => {
  return (router: Router, ...comparatives: T[]): boolean => {
    if (!comparatives.length) {
      return isLocationActive(router, ...defaultComparatives)
    }

    const [first, ...rest] = comparatives.map((name) => {
      const match = defaultComparatives.find((c) => c.name === name)

      if (!match) {
        throw new Error(`unknown comparative '${name}'`)
      }

      return match
    })

    return isLocationActive(router, first, ...rest)
  }
}

/**
 * just a dummy function to trick gettext tools
 *
 * @param msg
 */
export function $gettext(msg: string): string {
  return msg
}

/**
 * create a location with attached default values
 *
 * @param name
 * @param locations
 */
export const createLocation = (
  name: string,
  ...locations: RouteLocationNamedRaw[]
): RouteLocationNamedRaw =>
  merge(
    {},
    {
      name
    },
    ...locations.map((location) => ({
      ...(location.params && { params: location.params }),
      ...(location.query && { query: location.query })
    }))
  )
