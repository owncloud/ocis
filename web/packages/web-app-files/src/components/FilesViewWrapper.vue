<template>
  <div class="files-view-wrapper oc-width-expand">
    <div id="files-view" v-bind="$attrs">
      <slot />
    </div>
  </div>

  <portal v-if="isEmbedModeEnabled" to="app.runtime.footer">
    <embed-actions />
  </portal>
</template>

<script lang="ts" setup>
import { useEmbedMode } from '@ownclouders/web-pkg'
import EmbedActions from './EmbedActions/EmbedActions.vue'

defineOptions({
  inheritAttrs: false
})

const { isEnabled: isEmbedModeEnabled } = useEmbedMode()
</script>

<style lang="scss" scoped>
.files-view-wrapper {
  height: 100%;
  position: relative;
  overflow-y: auto;
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: max-content max-content 1fr;
  gap: 0 0;
  grid-template-areas:
    'header'
    'upload'
    'main';

  &:focus {
    outline: none;
  }
}

#files-view {
  grid-area: main;
  z-index: 0;
  outline: none;
}
</style>
