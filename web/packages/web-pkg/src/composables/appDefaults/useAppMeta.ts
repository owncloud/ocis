import { computed } from 'vue'
import { AppsStore } from '../piniaStores'

interface AppMetaOptions {
  appsStore: AppsStore
  applicationId: string
}

export function useAppMeta({ appsStore, applicationId }: AppMetaOptions) {
  const applicationMeta = computed(() => {
    const appInfo = appsStore.apps[applicationId]
    if (!appInfo) {
      throw new Error(`useAppConfig: could not find config for applicationId: ${applicationId}`)
    }
    return appInfo || {}
  })

  return {
    applicationMeta
  }
}
