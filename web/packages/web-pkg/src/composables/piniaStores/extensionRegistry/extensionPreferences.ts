import { defineStore } from 'pinia'
import { useLocalStorage } from '@vueuse/core'
import { unref } from 'vue'
import { Extension, ExtensionPoint } from './types'

export interface ExtensionPreferenceItem {
  extensionPointId: string
  selectedExtensionIds: string[]
}

export const useExtensionPreferencesStore = defineStore('extensionPreferences', () => {
  // maps extension point ids to selected extensions ids
  const extensionPreferences = useLocalStorage<Record<string, ExtensionPreferenceItem>>(
    'extensionPreferences',
    {}
  )

  const getExtensionPreference = (
    extensionPointId: string,
    defaultExtensionIds: string[]
  ): ExtensionPreferenceItem => {
    const extensionPreference = extensionPreferences.value[extensionPointId]
    if (extensionPreference) {
      return extensionPreference
    }
    return {
      extensionPointId,
      selectedExtensionIds: defaultExtensionIds
    }
  }
  const extractDefaultExtensionIds = (
    extensionPoint: ExtensionPoint<Extension>,
    extensions: Extension[]
  ): string[] => {
    if (extensionPoint.multiple) {
      return extensions.map((extension) => extension.id)
    }
    if (extensionPoint.defaultExtensionId) {
      return [extensionPoint.defaultExtensionId]
    }
    return []
  }
  const setSelectedExtensionIds = (extensionPointId: string, extensionIds: string[]) => {
    if (!Object.hasOwn(unref(extensionPreferences), extensionPointId)) {
      extensionPreferences.value[extensionPointId] = {
        extensionPointId,
        selectedExtensionIds: extensionIds
      }
      return
    }
    extensionPreferences.value[extensionPointId].selectedExtensionIds = extensionIds
  }

  return {
    extensionPreferences,
    extractDefaultExtensionIds,
    getExtensionPreference,
    setSelectedExtensionIds
  }
})

export type ExtensionPreferencesStore = ReturnType<typeof useExtensionPreferencesStore>
