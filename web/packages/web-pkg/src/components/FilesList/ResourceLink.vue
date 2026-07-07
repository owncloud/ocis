<template>
  <component
    :is="componentType"
    v-bind="componentProps"
    v-if="isResourceClickable"
    :target="linkTarget"
    :draggable="false"
    class="oc-resource-link"
    @dragstart.prevent.stop
    @click="emitClick"
  >
    <slot />
  </component>
  <span v-else>
    <slot />
  </span>
</template>

<script lang="ts" setup>
import { useConfigStore } from '../../composables'
import { storeToRefs } from 'pinia'
import { computed, unref } from 'vue'
import { SpaceResource, Resource } from '@ownclouders/web-client'
import { isSpaceResource } from '@ownclouders/web-client'
import { RouteLocationRaw } from 'vue-router'

/**
 * Wraps content in a resource link
 */

interface Props {
  /**
   * The resource folder link
   */
  link?: RouteLocationRaw
  /**
   * The resource to be displayed
   */
  resource: SpaceResource | Resource
  /**
   * Asserts whether clicking on the resource name triggers any action
   */
  isResourceClickable?: boolean
}
interface Emits {
  (event: 'click'): void
}
const { link = null, resource, isResourceClickable = true } = defineProps<Props>()
const emit = defineEmits<Emits>()
const configStore = useConfigStore()
const { options } = storeToRefs(configStore)

const linkTarget = computed(() => {
  return unref(options).cernFeatures && link && !resource.isFolder ? '_blank' : '_self'
})
const isNavigatable = computed(() => {
  if (isSpaceResource(resource)) {
    return (resource.isFolder || link) && !resource.disabled
  }

  return resource.isFolder || link
})
const componentType = computed(() => {
  return unref(isNavigatable) ? 'router-link' : 'oc-button'
})
const componentProps = computed(() => {
  if (!unref(isNavigatable)) {
    return {
      appearance: 'raw',
      gapSize: 'none',
      justifyContent: 'left'
    }
  }

  return {
    to: link
  }
})
function emitClick() {
  if (unref(isNavigatable)) {
    return
  }

  /**
   * Triggered when the resource is a file and the name is clicked
   */
  emit('click')
}
</script>
<style lang="scss">
.oc-resource-link {
  // necessary for the focus outline to show up
  display: inline-flex;
}
</style>
