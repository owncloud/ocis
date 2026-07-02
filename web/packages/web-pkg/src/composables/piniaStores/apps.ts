import { defineStore } from 'pinia'
import { computed, ref, unref } from 'vue'
import { AppConfigObject, ApplicationInformation, ApplicationFileExtension } from '../../apps'
import { Translations } from 'vue3-gettext'

export const useAppsStore = defineStore('apps', () => {
  const apps = ref<Record<string, ApplicationInformation>>({})
  const externalAppConfig = ref<Record<string, AppConfigObject>>({})
  const fileExtensions = ref<ApplicationFileExtension[]>([])

  const appIds = computed(() => Object.keys(unref(apps)))

  const registerApp = (appInfo: ApplicationInformation, translations?: Translations) => {
    if (!appInfo.id) {
      return
    }

    if (appInfo.extensions) {
      appInfo.extensions.forEach((extension) => {
        registerFileExtension({ appId: appInfo.id, data: extension })
      })
    }

    unref(apps)[appInfo.id] = {
      defaultExtension: appInfo.defaultExtension || '',
      icon: 'check_box_outline_blank',
      name: appInfo.name || appInfo.id,
      translations,
      ...appInfo
    }
  }

  const registerFileExtension = ({
    appId,
    data
  }: {
    appId: string
    data: ApplicationFileExtension
  }) => {
    unref(fileExtensions).push({
      app: appId,
      extension: data.extension,
      createFileHandler: data.createFileHandler,
      label: data.label,
      mimeType: data.mimeType,
      routeName: data.routeName,
      newFileMenu: data.newFileMenu,
      icon: data.icon,
      name: data.name,
      hasPriority:
        data.hasPriority ||
        unref(externalAppConfig)?.[appId]?.priorityExtensions?.includes(data.extension) ||
        false,
      secureView: data.secureView || false,
      customHandler: data.customHandler || null
    })
  }

  const loadExternalAppConfig = ({ appId, config }: { appId: string; config: AppConfigObject }) => {
    externalAppConfig.value = { ...unref(externalAppConfig), [appId]: config }
  }

  const isAppEnabled = (appId: string) => {
    return unref(appIds).includes(appId)
  }

  return {
    apps,
    externalAppConfig,
    appIds,
    fileExtensions,

    registerApp,
    registerFileExtension,
    loadExternalAppConfig,
    isAppEnabled
  }
})

export type AppsStore = ReturnType<typeof useAppsStore>
