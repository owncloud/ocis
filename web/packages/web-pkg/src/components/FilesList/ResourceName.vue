<template>
  <span
    class="oc-resource-name"
    :class="[{ 'oc-display-inline-block': !truncateName }]"
    :data-test-resource-path="fullPath"
    :data-test-resource-name="fullName"
    :data-test-resource-type="type"
    :title="htmlTitle"
    :aria-label="htmlTitle"
  >
    <span v-if="truncateName" class="oc-text-truncate">
      <span class="oc-resource-basename" v-text="displayName" />
    </span>
    <span v-else class="oc-resource-basename oc-text-break" v-text="displayName" /><span
      v-if="extension && isExtensionDisplayed"
      class="oc-resource-extension"
      v-text="displayExtension"
    />
  </span>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'

/**
 * Props for the ResourceName component.
 * @property {string} name - The name of the resource.
 * @property {string} [extension] - The file extension, if any.
 * @property {string} type - The type of the resource.
 * @property {string} fullPath - The full path to the resource.
 * @property {boolean} [isPathDisplayed] - Whether to display the path.
 * @property {boolean} [isExtensionDisplayed] - Whether to display the extension.
 * @property {boolean} [truncateName] - Whether to truncate the name.
 */

interface Props {
  name: string
  extension?: string
  type: string
  fullPath: string
  isPathDisplayed?: boolean
  isExtensionDisplayed?: boolean
  truncateName?: boolean
}
const {
  name,
  extension = '',
  type,
  fullPath,
  isPathDisplayed = false,
  isExtensionDisplayed = true,
  truncateName = true
} = defineProps<Props>()

const fullName = computed(() => {
  return (unref(displayPath) || '') + name
})

const displayName = computed(() => {
  if (extension) {
    return name.slice(0, -extension.length - 1)
  }
  return name
})

const displayExtension = computed(() => {
  return extension ? '.' + extension : ''
})

const displayPath = computed(() => {
  if (!isPathDisplayed) {
    return null
  }
  const pathSplit = fullPath.replace(/^\//, '').split('/')
  if (pathSplit.length < 2) {
    return null
  }
  if (pathSplit.length === 2) {
    return pathSplit[0] + '/'
  }
  return `…/${pathSplit[pathSplit.length - 2]}/`
})

const htmlTitle = computed(() => {
  if (isExtensionDisplayed) {
    return `${unref(displayName)}${unref(displayExtension)}`
  }
  return unref(displayName)
})
</script>

<style lang="scss">
.oc-resource {
  &-name {
    display: flex;
    min-width: 0;

    &:hover {
      text-decoration: underline;
      text-decoration-color: var(--oc-color-text-default);
    }
  }

  &-basename,
  &-extension {
    color: var(--oc-color-text-default);
    white-space: pre;
  }

  &-path {
    color: var(--oc-color-text-muted);
    white-space: pre;
  }
}
</style>
