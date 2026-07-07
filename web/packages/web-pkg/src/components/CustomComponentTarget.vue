<template>
  <component
    :is="toRaw(extension.content)"
    v-for="extension in extensions"
    :key="`custom-component-${extension.id}`"
  />
</template>

<script lang="ts" setup>
import { computed, unref, toRaw } from 'vue'
import {
  CustomComponentExtension,
  ExtensionPoint,
  useExtensionPreferencesStore,
  useExtensionRegistry
} from '../composables'

interface Props {
  extensionPoint: ExtensionPoint<CustomComponentExtension>
}
const { extensionPoint } = defineProps<Props>()
const extensionRegistry = useExtensionRegistry()
const extensionPreferences = useExtensionPreferencesStore()

const allExtensions = computed(() => {
  return extensionRegistry.requestExtensions(extensionPoint)
})

const defaultExtensionIds = extensionPreferences.extractDefaultExtensionIds(
  extensionPoint,
  unref(allExtensions)
)

const extensions = computed<CustomComponentExtension[]>(() => {
  // TODO: for `multiple` we want to respect the selected extensions as well in the future.
  if (extensionPoint.multiple || unref(allExtensions).length <= 1) {
    return unref(allExtensions)
  }

  const preference = extensionPreferences.getExtensionPreference(
    extensionPoint.id,
    defaultExtensionIds
  )
  if (preference.selectedExtensionIds.length) {
    return [
      unref(allExtensions).find((extension) =>
        preference.selectedExtensionIds.includes(extension.id)
      ) || unref(allExtensions)[0]
    ]
  }

  // if no user preference and no default provided, return the first one.
  return [unref(allExtensions)[0]]
})
</script>
