<template>
  <div id="ghost-element" class="ghost-element">
    <div class="ghost-element-layer1 oc-rounded">
      <resource-icon class="oc-p-xs" :resource="previewItems[0]" />
      <div v-if="showSecondLayer" class="ghost-element-layer2 oc-rounded" />
      <div v-if="showThirdLayer" class="ghost-element-layer3 oc-rounded" />
    </div>
    <span class="badge">{{ itemCount }}</span>
  </div>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import { Resource } from '@ownclouders/web-client'
import ResourceIcon from './ResourceIcon.vue'

/**
 * Please head to the ownCloud web ResourceTable component (https://github.com/owncloud/web/blob/master/packages/web-app-files/src/components/FilesList/ResourceTable.vue) for a demo of the Ghost Element.
 */

interface Props {
  previewItems: Resource[]
}
const { previewItems } = defineProps<Props>()
const layerCount = computed(() => {
  return Math.min(previewItems.length, 3)
})
const showSecondLayer = computed(() => {
  return unref(layerCount) > 1
})
const showThirdLayer = computed(() => {
  return unref(layerCount) > 2
})
const itemCount = computed(() => {
  return previewItems.length
})
</script>

<style lang="scss">
.ghost-element-layer1 {
  position: relative;
  background-color: var(--oc-color-background-hover);

  .ghost-element-layer2 {
    position: absolute;
    background-color: var(--oc-color-background-hover);
    filter: brightness(0.82);
    top: 3px;
    left: 3px;
    right: -3px;
    bottom: -3px;
    z-index: -1;
  }
  .ghost-element-layer3 {
    position: absolute;
    background-color: var(--oc-color-background-hover);
    filter: brightness(0.72);
    top: 6px;
    left: 6px;
    right: -6px;
    bottom: -6px;
    z-index: -2;
  }
}
.ghost-element {
  background-color: transparent;
  padding-top: var(--oc-space-xsmall);
  padding-left: 5px;
  z-index: var(--oc-z-index-modal);
  position: absolute;
  .icon-wrapper {
    position: relative;
  }
  .badge {
    position: absolute;
    top: -2px;
    right: -8px;
    padding: var(--oc-space-xsmall);
    line-height: var(--oc-space-small);
    -webkit-border-radius: 30px;
    -moz-border-radius: 30px;
    border-radius: 30px;
    min-width: var(--oc-space-small);
    height: var(--oc-space-small);
    text-align: center;

    font-size: 12px;
    background: red;
    color: white;
  }
}
</style>
