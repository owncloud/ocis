import { computed, ComputedRef, unref } from 'vue'
import { useRoute } from './useRoute'

export const useRouteName = (): ComputedRef<string> => {
  const route = useRoute()
  return computed(() => {
    return unref(route).name as string
  })
}
