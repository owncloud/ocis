<template>
  <div
    class="oc-text-editor oc-width-1-1 oc-height-1-1"
    :class="{ 'oc-text-editor-readonly': isReadOnly }"
  >
    <text-editor-component
      :resource="resource"
      :application-config="applicationConfig"
      :current-content="currentContent"
      :is-read-only="isReadOnly"
      @update:current-content="$emit('update:currentContent', $event)"
    />
  </div>
</template>

<script lang="ts" setup>
import { Resource } from '@ownclouders/web-client'
import { AppConfigObject, TextEditor as TextEditorComponent } from '@ownclouders/web-pkg'

interface Props {
  applicationConfig: AppConfigObject
  currentContent: string
  isReadOnly?: boolean
  resource: Resource
}
interface Emits {
  (e: 'update:currentContent', value: Event): void
}
const { applicationConfig, currentContent, isReadOnly, resource } = defineProps<Props>()

defineEmits<Emits>()
</script>

<style lang="scss">
.oc-text-editor-readonly {
  //Toastui Editor doesn't have margins in view mode, adjusted for uniformity
  padding: 18px 25px;

  //Fixes in readonly mode vertical scrolling is not available
  height: calc(100dvh - 52px);
  overflow: auto;
}

// Make url links limited in height and scrollable since base64 encoded images might be very long
.toastui-editor-md-link-url {
  display: block;
  max-height: 100px;
  overflow: auto;
}

.toastui-editor-defaultUI {
  // Adjustments to match our theming
  border: none;
  // Adjustments to not overlay with the app switcher
  z-index: 0;
}

.toastui-editor-defaultUI-toolbar,
.toastui-editor-md-tab-container {
  background-color: var(--oc-color-background-default) !important;
}
</style>
