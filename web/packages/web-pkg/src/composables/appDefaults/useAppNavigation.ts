import { Ref, ref, unref } from 'vue'
import { RouteLocation, Router, RouteParams } from 'vue-router'

import { MaybeRef } from '../../utils'
import { FileContext } from './types'
import { LocationQuery, QueryValue } from '../router'
import { Resource } from '@ownclouders/web-client'
import { useFileRouteReplace } from '../router/useFileRouteReplace'

interface AppNavigationOptions {
  router: Router
  currentFileContext: MaybeRef<FileContext>
}

export interface AppNavigationResult {
  closeApp(): void
  replaceInvalidFileRoute(context: MaybeRef<FileContext>, resource: Resource): boolean
  closed: Ref<boolean>
}

export const contextRouteNameKey = 'contextRouteName'
export const contextRouteParamsKey = 'contextRouteParams'
export const contextRouteQueryKey = 'contextRouteQuery'

/*
  vue-router type bindings do not allow nested objects
  because they are not handled by default. We override
  parseQuery and stringifyQuery and handle it there.
  That's why we have types that match the router types
  and break them here once on purpose in encapsulated
  functions
*/
export const routeToContextQuery = (location: RouteLocation): LocationQuery => {
  const { params, query } = location

  const contextQuery: Record<string, QueryValue> = {}
  const contextQueryItems = ['fileId', 'shareId', 'q_share-visibility'].concat(
    (location as any).meta?.contextQueryItems || []
  ) as string[]
  for (const queryItem of contextQueryItems) {
    contextQuery[queryItem] = query[queryItem]
  }

  return {
    [contextRouteNameKey]: location.name,
    [contextRouteParamsKey]: params,
    [contextRouteQueryKey]: contextQuery
  } as any
}
export const contextQueryToFileContextProps = (
  query: LocationQuery
): { routeName: string; routeParams: RouteParams; routeQuery: LocationQuery } => {
  return {
    routeName: queryItemAsString(query[contextRouteNameKey]),
    routeParams: query[contextRouteParamsKey] as any,
    routeQuery: query[contextRouteQueryKey] as any
  }
}

export const queryItemAsString = (
  queryItem: string | number | Exclude<string | number, null | undefined>[]
): string => {
  if (Array.isArray(queryItem)) {
    return queryItem[0].toString()
  }

  return queryItem?.toString()
}

export function useAppNavigation({
  router,
  currentFileContext
}: AppNavigationOptions): AppNavigationResult {
  const navigateToContext = (context: MaybeRef<FileContext>) => {
    const { fileName, routeName, routeParams, routeQuery } = unref(context)

    if (!unref(routeName)) {
      return router.push({ path: '/' })
    }

    return router.push({
      name: unref(routeName),
      params: unref(routeParams),
      query: {
        ...unref(routeQuery),
        scrollTo: unref(fileName)
      }
    })
  }

  const { replaceInvalidFileRoute: replaceInvalidFileRouteGeneric } = useFileRouteReplace({
    router
  })
  const replaceInvalidFileRoute = (context: MaybeRef<FileContext>, resource: Resource) => {
    const ctx = unref(context)
    return replaceInvalidFileRouteGeneric({
      space: unref(ctx.space),
      resource,
      path: unref(ctx.item),
      fileId: unref(ctx.itemId)
    })
  }

  const closed = ref(false)
  const closeApp = () => {
    closed.value = true
    return navigateToContext(currentFileContext)
  }

  return {
    replaceInvalidFileRoute,
    closeApp,
    closed
  }
}
