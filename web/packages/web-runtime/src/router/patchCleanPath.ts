import { RouteLocation, RouteLocationNormalizedLoaded, RouteLocationRaw, Router } from 'vue-router'

// type: patch
// temporary patch till we have upgraded web to the latest vue router which make this obsolete
// this takes care that routes like 'foo/bar/baz' which by default would be converted to 'foo%2Fbar%2Fbaz' stay as they are
// should immediately go away and be removed after finalizing the update
// to apply the patch to a route add meta.patchCleanPath = true to it
// to patch needs to be enabled on a route level, to do so add meta.patchCleanPath = true property to the route
// c.f. https://github.com/vuejs/router/issues/1638
export const patchRouter = (router: Router) => {
  const cleanPath = (route: string) =>
    [
      ['%2F', '/'],
      ['//', '/']
    ].reduce((path, rule) => path.replaceAll(rule[0], rule[1]), route || '')

  const bindResolve = router.resolve.bind(router)
  router.resolve = (
    raw: RouteLocationRaw,
    currentLocation?: RouteLocationNormalizedLoaded
  ): RouteLocation & {
    href: string
  } => {
    const resolved = bindResolve(raw, currentLocation)
    if (resolved.meta?.patchCleanPath !== true) {
      return resolved
    }

    return {
      ...resolved,
      href: cleanPath(resolved.href),
      path: cleanPath(resolved.path),
      fullPath: cleanPath(resolved.fullPath)
    }
  }

  const routerMethodFactory =
    (method: (arg: RouteLocationRaw) => ReturnType<Router['push']>) => (to: RouteLocationRaw) => {
      const resolved = router.resolve(to)
      if (resolved.meta?.patchCleanPath !== true) {
        return method(to)
      }

      return method({
        path: cleanPath(resolved.fullPath),
        query: resolved.query
      })
    }

  router.push = routerMethodFactory(router.push.bind(router))
  router.replace = routerMethodFactory(router.replace.bind(router))

  return router
}
