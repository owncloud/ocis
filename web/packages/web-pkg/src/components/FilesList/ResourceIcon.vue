<template>
  <oc-icon
    :key="`resource-icon-${icon.name}`"
    :name="icon.name"
    :color="icon.color"
    :size="size"
    :class="['oc-resource-icon', iconTypeClass]"
  />
</template>

<script lang="ts" setup>
import { computed, inject, unref } from 'vue'
import { Resource } from '@ownclouders/web-client'
import {
  IconType,
  createDefaultFileIconMapping,
  ResourceIconMapping,
  resourceIconMappingInjectionKey
} from '../../helpers/resource/icon'

interface Props {
  /**
   * The resource to be displayed
   */
  resource: Resource
  /**
   * The size of the icon. Defaults to large.
   * `xsmall, small, medium, large, xlarge, xxlarge`
   */
  size?: 'xsmall' | 'small' | 'medium' | 'large' | 'xlarge' | 'xxlarge' | 'xxxlarge' | string
}

const defaultFolderIcon: IconType = {
  name: 'resource-type-folder',
  color: 'var(--oc-color-icon-folder)'
}

const defaultSpaceIcon: IconType = {
  name: 'layout-grid',
  color: 'var(--oc-color-swatch-passive-default)'
}
const defaultFallbackIcon: IconType = {
  name: 'resource-type-file',
  color: 'var(--oc-color-text-default)'
}

const defaultFileIconMapping = createDefaultFileIconMapping()
const { resource, size = 'large' } = defineProps<Props>()

const iconMappingInjection = inject<ResourceIconMapping>(resourceIconMappingInjectionKey)

const isFolder = computed(() => {
  // fallback is necessary since
  // sometimes resources without a type
  // but with `isFolder` are being passed
  return resource.type === 'folder' || resource.isFolder
})

const isSpace = computed(() => {
  return resource.type === 'space'
})
const extension = computed(() => {
  return resource.extension?.toLowerCase()
})
const mimeType = computed(() => {
  return resource.mimeType?.toLowerCase()
})

const icon = computed((): IconType => {
  if (unref(isSpace)) {
    return defaultSpaceIcon
  }
  if (unref(isFolder)) {
    return defaultFolderIcon
  }

  const icon =
    defaultFileIconMapping[unref(extension)] ||
    iconMappingInjection?.mimeType[unref(mimeType)] ||
    iconMappingInjection?.extension[unref(extension)]

  return {
    ...defaultFallbackIcon,
    ...icon
  }
})

const iconTypeClass = computed(() => {
  if (unref(isSpace)) {
    return 'oc-resource-icon-space'
  }
  if (unref(isFolder)) {
    return 'oc-resource-icon-folder'
  }
  return 'oc-resource-icon-file'
})
</script>

<style lang="scss">
span.oc-resource-icon {
  display: inline-flex;
  align-items: center;
  vertical-align: middle;

  &-file svg {
    height: 70%;
  }
}
</style>
