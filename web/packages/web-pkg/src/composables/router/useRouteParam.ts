import { computed, Ref, unref } from 'vue'
import { useRouter } from './useRouter'
import { ParamValue } from './types'
import { queryItemAsString } from '../appDefaults'

export const useRouteParam = (name: string, defaultValue?: ParamValue): Ref<ParamValue> => {
  const router = useRouter()

  return computed({
    get() {
      return queryItemAsString(unref(router.currentRoute).params[name]) || defaultValue
    },
    async set(v) {
      if (unref(router.currentRoute).params[name] === v) {
        return
      }
      await router.replace({
        params: {
          ...unref(router.currentRoute).params,
          [name]: v
        }
      })
    }
  })
}
