import { RouteLocation, RouteParams, Router, RouteRecordNormalized } from 'vue-router'
import {
  AuthContext,
  authContextValues,
  contextQueryToFileContextProps,
  queryItemAsString,
  WebRouteMeta
} from '@ownclouders/web-pkg'

/**
 * Checks if the `to` route or the route it was reached from (i.e. the `contextRoute`) needs authentication from the IdP and a successfully fetched ownCloud user.
 *
 * @param router {Router}
 * @param to {Route}
 * @returns {boolean}
 */
export const isUserContextRequired = (router: Router, to: RouteLocation): boolean => {
  const meta = getRouteMeta(to)
  if (meta.authContext === 'user') {
    return true
  }
  if (meta.authContext !== 'hybrid') {
    return false
  }

  const contextRoute = getContextRoute(router, to)
  return (
    !contextRoute ||
    getRouteMeta({ meta: contextRoute.meta } as RouteLocation).authContext === 'user'
  )
}

/**
 * Checks if the `to` route or the route it was reached from (i.e. the `contextRoute`) needs authentication from the IdP but should not try to fetch an ownCloud user.
 *
 * @param router {Router}
 * @param to {Route}
 * @returns {boolean}
 */
export const isIdpContextRequired = (router: Router, to: RouteLocation): boolean => {
  const meta = getRouteMeta(to)
  if (meta.authContext === 'idp') {
    return true
  }

  const contextRoute = getContextRoute(router, to)
  return (
    contextRoute && getRouteMeta({ meta: contextRoute.meta } as RouteLocation).authContext === 'idp'
  )
}

/**
 * Checks if the `to` route or the route it was reached from (i.e. the `contextRoute`) needs a resolved public link context (with or without password).
 *
 * @param router {Router}
 * @param to {Route}
 * @returns {boolean}
 */
export const isPublicLinkContextRequired = (router: Router, to: RouteLocation): boolean => {
  if (
    (to.params.driveAliasAndItem as string)?.startsWith('public/') ||
    (to.params.driveAliasAndItem as string)?.startsWith('ocm/')
  ) {
    return true
  }

  const meta = getRouteMeta(to)
  if (meta.authContext === 'publicLink') {
    return true
  }
  if (meta.authContext !== 'hybrid') {
    return false
  }

  const contextRoute = getContextRoute(router, to)
  return (
    contextRoute &&
    getRouteMeta({ meta: contextRoute.meta } as RouteLocation).authContext === 'publicLink'
  )
}

/**
 * Extracts the public link token from the various possible route params.
 *
 * @param to {Route}
 * @returns {string}
 */
export const extractPublicLinkToken = (to: RouteLocation): string => {
  const contextRouteParams = contextQueryToFileContextProps(to.query)?.routeParams
  if (contextRouteParams) {
    return extractPublicLinkTokenFromRouteParams(contextRouteParams)
  }
  return extractPublicLinkTokenFromRouteParams(to.params)
}

/**
 * Extracts the public link token from known possible occurrences in params of a route.
 *
 * @param params {LocationParams}
 */
const extractPublicLinkTokenFromRouteParams = (params: RouteParams): string => {
  if (Object.prototype.hasOwnProperty.call(params, 'driveAliasAndItem')) {
    const driveAliasAndItem = queryItemAsString(params.driveAliasAndItem)
    if (!driveAliasAndItem.startsWith('public/') && !driveAliasAndItem.startsWith('ocm/')) {
      return ''
    }
    return (params.driveAliasAndItem as string).split('/')[1]
  }
  return ((params.item || params.filePath || params.token || '') as string).split('/')[0]
}

/**
 * Asserts that no form of authentication is required.
 *
 * @param router {Router}
 * @param to {Route}
 * @returns {boolean}
 */
export const isAnonymousContext = (router: Router, to: RouteLocation): boolean => {
  return getRouteMeta(to).authContext === 'anonymous'
}

/**
 * The contextRoute in URLs is used to give applications additional context where the application route was triggered from
 * (e.g. from a project space, a public link file listing, a personal space, etc).
 * Application routes need to fulfill both their own auth requirements and the auth requirements from the context route.
 *
 * Example: the `preview` app and its routes don't explicitly require authentication (`meta.auth` is set to `false`), because
 * the app can be used from both an authenticated context or from a public link context. The information which endpoint
 * the preview app is supposed to load files from is transported via the contextRouteName, contextRouteParams and contextRouteQuery
 * in the URL (provided by the context that opens the preview app in the first place).
 */
const getContextRoute = (router: Router, to: RouteLocation): RouteRecordNormalized | null => {
  const contextRouteNameKey = 'contextRouteName'
  if (!to.query || !to.query[contextRouteNameKey]) {
    return null
  }

  return router.getRoutes().find((r) => r.name === to.query[contextRouteNameKey])
}

const getRouteMeta = (to: RouteLocation): WebRouteMeta => {
  if (!to.meta) {
    return {
      authContext: 'user'
    }
  }

  // rewrite deprecated `auth` property to the respective `authContext` value
  if (!to.meta.authContext && Object.prototype.hasOwnProperty.call(to.meta, 'auth')) {
    to.meta.authContext = to.meta.auth ? 'user' : 'hybrid'
    console.warn(
      `route key meta.auth is deprecated. Please switch to meta.authContext="${
        to.meta.authContext
      }" in route "${String(to.name)}".`
    )
  }

  if (to?.meta?.authContext) {
    if (authContextValues.includes(to.meta.authContext as AuthContext)) {
      return to.meta
    }
    console.warn(
      `invalid authContext "${to.meta.authContext}" in route "${String(
        to.name
      )}". must be one of [${authContextValues.join(', ')}].`
    )
  }
  if (to?.meta) {
    return {
      ...to.meta,
      authContext: 'user'
    }
  }
  return {
    authContext: 'user'
  }
}
