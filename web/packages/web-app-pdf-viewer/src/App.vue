<template>
  <object
    v-if="!isIos"
    class="pdf-viewer oc-width-1-1 oc-height-1-1"
    :data="url"
    :type="objectType"
    :aria-label="
      $pgettext('Accessible label for the object holding the PDF file content', 'PDF document')
    "
  />
  <div
    v-else
    class="oc-flex oc-flex-column oc-flex-center oc-flex-middle oc-width-1-1 oc-height-1-1 oc-gap-s"
  >
    <oc-button type="a" :href="url" target="_blank" appearance="filled" variation="primary">
      {{ $pgettext('Button to open a PDF file in the browser on iOS/iPadOS', 'Open PDF') }}
    </oc-button>
  </div>
</template>

<script lang="ts" setup>
interface Props {
  url: string
}
const { url } = defineProps<Props>()

const isSafari = navigator.userAgent?.includes('Safari') && !navigator.userAgent?.includes('Chrome')
const objectType = isSafari ? undefined : 'application/pdf'

// iOS/iPadOS does not support scrollable multi-page PDFs inside iframes.
const isIos = /iPad|iPhone|iPod/.test(navigator.userAgent)
</script>

<style scoped>
.pdf-viewer {
  border: none;
  margin: 0;
  padding: 0;
  overflow: hidden;
}
</style>
