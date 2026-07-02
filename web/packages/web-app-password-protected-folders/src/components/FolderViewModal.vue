<template>
  <div class="oc-height-1-1" tabindex="0">
    <app-loading-spinner v-if="isLoading" />
    <iframe
      v-show="!isLoading"
      id="iframe-folder-view"
      ref="iframeRef"
      class="oc-width-1-1 oc-height-1-1"
      :title="iframeTitle"
      :src="iframeUrl.href"
      sandbox="allow-scripts allow-forms allow-same-origin allow-popups"
      tabindex="0"
      @load="onLoad"
    ></iframe>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { Modal, useThemeStore } from '@ownclouders/web-pkg/src/composables'
import AppLoadingSpinner from '@ownclouders/web-pkg/src/components/AppLoadingSpinner.vue'
import { unref } from 'vue'
import { useGettext } from 'vue3-gettext'

const props = defineProps<{
  modal: Modal
  publicLink: string
  serverUrl: string
}>()

const iframeRef = ref<HTMLIFrameElement>()
const isLoading = ref(true)
const themeStore = useThemeStore()
const { current } = useGettext()

const iframeTitle = themeStore.currentTheme.common?.name
const iframeUrl = new URL(props.publicLink)
if (!['https:', 'http:'].includes(iframeUrl.protocol)) {
  throw new Error('Invalid URL scheme for iframe')
}
if (iframeUrl.origin !== new URL(props.serverUrl).origin) {
  throw new Error('URL does not belong to this server')
}
iframeUrl.searchParams.append('hide-logo', 'true')
iframeUrl.searchParams.append('hide-app-switcher', 'true')
iframeUrl.searchParams.append('hide-account-menu', 'true')
iframeUrl.searchParams.append('hide-navigation', 'true')
iframeUrl.searchParams.append('lang', current)

const onLoad = () => {
  isLoading.value = false
  unref(iframeRef).contentWindow.focus()
}
</script>

<style lang="scss">
.oc-modal.folder-view-modal {
  max-width: 80dvw;
  border: none;
  overflow: hidden;

  .oc-modal-title {
    display: none;
  }

  .oc-modal-body {
    padding: 0;

    &-message {
      height: 60dvh;
      margin: 0;
    }
  }

  .oc-modal-body-actions {
    background-color: var(--oc-color-swatch-brand-default);

    .oc-modal-body-actions-cancel {
      color: var(--oc-color-swatch-brand-contrast);
      outline-color: var(--oc-color-swatch-brand-contrast);

      &:hover:not([disabled]),
      &:focus:not([disabled]) {
        background-color: var(--oc-color-swatch-brand-contrast);
        color: var(--oc-color-swatch-brand-default);
        border-color: var(--oc-color-swatch-brand-contrast);
      }
    }
  }
}
</style>
