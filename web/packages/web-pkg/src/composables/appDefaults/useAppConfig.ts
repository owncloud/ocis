import { computed, Ref } from 'vue'
import type { AppConfigObject } from '../../apps'
import { AppsStore } from '../piniaStores'

export interface AppConfigOptions {
  appsStore: AppsStore
  applicationId: string
}

export interface AppConfigResult {
  applicationConfig: Ref<AppConfigObject>
}

export function useAppConfig(options: AppConfigOptions): AppConfigResult {
  const applicationConfig = computed(
    () => options.appsStore.externalAppConfig[options.applicationId] || {}
  )

  return {
    applicationConfig
  }
}
