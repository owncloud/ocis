import { defineStore } from 'pinia'
import { computed, ref, unref } from 'vue'
import merge from 'lodash-es/merge'
import { OptionsConfig, RawConfig } from './types'
import { urlJoin } from '@ownclouders/web-client'
import { useAppsStore } from '../apps'

const defaultOptions = {
  cernFeatures: false,
  concurrentRequests: {
    resourceBatchActions: 4,
    shares: {
      create: 4,
      list: 2
    },
    sse: 4
  },
  contextHelpers: true,
  contextHelpersReadMore: true,
  defaultExtension: 'files',
  disabledExtensions: [] as string[],
  editor: {
    autosaveEnabled: true,
    autosaveInterval: 120
  },
  embed: {
    enabled: false,
    target: 'resources'
  },
  ocm: {
    openRemotely: false
  },
  routing: {
    idBased: true,
    fullShareOwnerPaths: false
  },
  runningOnEos: false,
  tokenStorageLocal: true,
  userListRequiresFilter: false,
  hideLogo: false
} satisfies Partial<OptionsConfig>

export const useConfigStore = defineStore('config', () => {
  const appsStore = useAppsStore()

  const server = ref<RawConfig['server']>('')
  const theme = ref<RawConfig['theme']>('')
  const options = ref<RawConfig['options']>({ ...defaultOptions })

  const apps = ref<RawConfig['apps']>([])
  const externalApps = ref<RawConfig['external_apps']>([])
  const customTranslations = ref<RawConfig['customTranslations']>([])
  const oAuth2 = ref<RawConfig['auth']>()
  const openIdConnect = ref<RawConfig['openIdConnect']>()
  const sentry = ref<RawConfig['sentry']>()
  const scripts = ref<RawConfig['scripts']>([])
  const styles = ref<RawConfig['styles']>([])

  const maintenanceMode = ref(false)

  const isInVault = ref(false)

  const serverUrl = computed(() =>
    urlJoin(unref(server) || window.location.origin, { trailingSlash: true })
  )

  const isOAuth2 = computed(() => !!unref(oAuth2))
  const isOIDC = computed(() => !!unref(openIdConnect))

  const loadConfig = (data: RawConfig) => {
    if (data.server) {
      server.value = data.server?.endsWith('/') ? data.server : data.server + '/'
    }

    apps.value = data.apps || []
    customTranslations.value = data.customTranslations || []
    oAuth2.value = data.auth
    openIdConnect.value = data.openIdConnect
    sentry.value = data.sentry
    scripts.value = data.scripts || []
    styles.value = data.styles || []
    theme.value = data.theme

    if (data.options) {
      options.value = merge({ ...defaultOptions }, data.options)
      // ocm.openRemotely will not be loaded from config, but set based on cernFeatures option
      unref(options).ocm.openRemotely = unref(options).cernFeatures
      // routing will not be loaded from config, but set based on cernFeatures option
      unref(options).routing.idBased = !unref(options).cernFeatures
      unref(options).routing.fullShareOwnerPaths = unref(options).cernFeatures
    }

    if (data.external_apps) {
      externalApps.value = data.external_apps
      data.external_apps.filter(Boolean).forEach((externalApp) => {
        appsStore.loadExternalAppConfig({ appId: externalApp.id, config: externalApp.config })
      })
    }
  }

  const setMaintenanceMode = (value: boolean) => {
    maintenanceMode.value = value
  }

  const setIsInVault = (value: boolean) => {
    isInVault.value = value
  }

  return {
    options,
    oAuth2,
    openIdConnect,
    isOAuth2,
    isOIDC,
    customTranslations,
    apps,
    externalApps,
    sentry,
    theme,
    scripts,
    styles,
    serverUrl,
    maintenanceMode,
    setIsInVault,
    isInVault,
    loadConfig,
    setMaintenanceMode
  }
})

export type ConfigStore = ReturnType<typeof useConfigStore>
