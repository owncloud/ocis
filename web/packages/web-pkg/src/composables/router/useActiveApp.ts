import { computed, ComputedRef, unref } from 'vue'
import { useRoute } from './useRoute'
import { RouteLocationNormalizedLoaded } from 'vue-router'

export const activeApp = (route: RouteLocationNormalizedLoaded): string => {
  const removeVaultPrefix = route.path.replace(/^\/vault/, '')
  return removeVaultPrefix.split('/')[1]
}

export const useActiveApp = (): ComputedRef<string> => {
  const route = useRoute()
  return computed(() => {
    return activeApp(unref(route))
  })
}
