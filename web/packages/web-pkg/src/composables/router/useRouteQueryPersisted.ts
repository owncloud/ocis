import { Ref, watch, unref } from 'vue'
import { useRouteQuery } from './useRouteQuery'
import { useLocalStorage } from '../localStorage/useLocalStorage'
import { QueryValue } from './types'
import { queryItemAsString } from '../appDefaults'

export interface RouteQueryPersistedOptions {
  name: string
  defaultValue: QueryValue
  storagePrefix?: string
}

interface WatcherValue {
  value: QueryValue
  source: string
}

export const useRouteQueryPersisted = (options: RouteQueryPersistedOptions): Ref<QueryValue> => {
  const routeQueryVariable = useRouteQuery(options.name)
  const localStorageVariable = useLocalStorage<QueryValue>(localStorageKey(options))
  watch(
    (): WatcherValue => {
      if (unref(routeQueryVariable)) {
        return {
          value: unref(routeQueryVariable),
          source: 'route'
        }
      }
      if (unref(localStorageVariable)) {
        return {
          value: unref(localStorageVariable),
          source: 'storage'
        }
      }
      return {
        value: options.defaultValue,
        source: 'default'
      }
    },
    (val) => {
      if (['route', 'default'].includes(val.source)) {
        localStorageVariable.value =
          val.value === options.defaultValue ? undefined : queryItemAsString(val.value)
      }
      if (['storage', 'default'].includes(val.source)) {
        routeQueryVariable.value = val.value
      }
    },
    { immediate: true }
  )
  return routeQueryVariable
}

const localStorageKey = (options: RouteQueryPersistedOptions): string => {
  return ['oc-options', options.storagePrefix, options.name].filter(Boolean).join('_')
}
