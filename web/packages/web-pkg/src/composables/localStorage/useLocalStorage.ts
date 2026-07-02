import { ref, watch, unref, Ref } from 'vue'

export const useLocalStorage = <T>(key: string, defaultValue: T = undefined) => {
  const existingValue = window.localStorage.getItem(key)
  const variable = ref(defaultValue)

  if (existingValue) {
    try {
      variable.value = JSON.parse(existingValue)
    } catch {
      ;(variable as Ref<string>).value = existingValue
    }
  }

  watch(
    () => unref(variable),
    (val, old) => {
      if (val === old) {
        return
      }
      if (val !== undefined) {
        window.localStorage.setItem(key, typeof val === 'string' ? val : JSON.stringify(val))
      } else {
        window.localStorage.removeItem(key)
      }
    },
    { deep: true }
  )
  return variable
}
