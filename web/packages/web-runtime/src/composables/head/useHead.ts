import { computed, unref } from 'vue'
import { useHead as _useHead } from '@vueuse/head'
import {
  useCapabilityStore,
  useThemeStore,
  getBackendVersion,
  getWebVersion
} from '@ownclouders/web-pkg'
import { storeToRefs } from 'pinia'

export const useHead = () => {
  const themeStore = useThemeStore()
  const capabilityStore = useCapabilityStore()
  const { currentTheme } = storeToRefs(themeStore)

  const favicon = computed(() => currentTheme.value.logo.favicon)

  // Get brand color directly from theme store (reactive to theme changes)
  const themeColor = computed(
    () => currentTheme.value.designTokens.colorPalette['swatch-brand-default']
  )

  _useHead(
    computed(() => {
      return {
        meta: [
          {
            name: 'generator',
            content: [getWebVersion(), getBackendVersion({ capabilityStore })]
              .filter(Boolean)
              .join(', ')
          },
          {
            name: 'theme-color',
            content: unref(themeColor)
          }
        ],
        ...(unref(favicon) && { link: [{ rel: 'icon', href: unref(favicon) }] })
      }
    })
  )
}
