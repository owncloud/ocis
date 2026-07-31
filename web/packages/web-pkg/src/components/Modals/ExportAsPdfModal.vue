<template>
  <div class="oc-height-1-1" tabindex="0">
    <app-loading-spinner v-if="isLoading" />
    <iframe
      v-show="!isLoading"
      ref="iframeRef"
      class="oc-width-1-1 oc-height-1-1"
      :title="iframeTitle"
      :src="iframeUrl.href"
      tabindex="0"
      @load="onLoad"
    />
  </div>
</template>

<script lang="ts" setup>
import { onBeforeUnmount, onMounted, ref } from 'vue'
import {
  embedModeLocationPickMessageData,
  Modal,
  useEmbedMode,
  useGetMatchingSpace,
  useMessages,
  useModals,
  useRouter,
  useThemeStore
} from '../../composables'
import { RouteLocationRaw } from 'vue-router'
import AppLoadingSpinner from '../AppLoadingSpinner.vue'
import { Resource } from '@ownclouders/web-client'
import { unref } from 'vue'
import { useGettext } from 'vue3-gettext'
import { useExportAsPdfWorker } from '../../composables/webWorkers/exportAsPdfWorker'

const { modal, parentFolderLink, originalResource, content } = defineProps<{
  modal: Modal
  parentFolderLink: RouteLocationRaw
  originalResource: Resource
  content: string
}>()

const themeStore = useThemeStore()
const { $pgettext } = useGettext()
const router = useRouter()
const { removeModal } = useModals()
const { showMessage, showErrorMessage } = useMessages()
const { getMatchingSpace } = useGetMatchingSpace()
const { startWorker } = useExportAsPdfWorker()
const { verifyMessageOrigin } = useEmbedMode()

const parentFolderRoute = router.resolve(parentFolderLink)
const iframeTitle = themeStore.currentTheme.common?.name
const iframeUrl = new URL(parentFolderRoute.href, window.location.origin)
iframeUrl.searchParams.append('hide-logo', 'true')
iframeUrl.searchParams.append('embed', 'true')
iframeUrl.searchParams.append('embed-target', 'location')
iframeUrl.searchParams.append('embed-choose-file-name', 'true')
iframeUrl.searchParams.append('embed-delegate-authentication', 'false')
iframeUrl.searchParams.append(
  'embed-choose-file-name-suggestion',
  originalResource.name.replace('.md', '.pdf')
)

const iframeRef = ref<HTMLIFrameElement>()
const isLoading = ref(true)

function onLoad() {
  isLoading.value = false
  unref(iframeRef).contentWindow.focus()
}

function onLocationPick({ data }: MessageEvent) {
  const { resources, fileName }: embedModeLocationPickMessageData = data.data

  const destinationFolder: Resource = resources[0]
  const space = getMatchingSpace(destinationFolder)

  startWorker(destinationFolder, space, fileName, content, (result) => {
    if (result.failed.length > 0) {
      console.error(result.failed)
      showErrorMessage({
        title: $pgettext(
          'Error toast message title shown to a user when exporting a file as PDF via the export as PDF modal failed.',
          'Unable to export "%{fileName}"',
          { fileName }
        ),
        errors: [result.failed[0].error]
      })
    }

    if (result.successful.length > 0) {
      showMessage({
        title: $pgettext(
          'Success toast message title shown to a user when a file is exported as PDF via the export as PDF modal.',
          '"%{fileName}" was exported successfully',
          { fileName }
        )
      })
    }
  })

  removeModal(modal.id)
}

function handleMessage(event: MessageEvent) {
  if (!verifyMessageOrigin(event.origin)) {
    return
  }

  if (event.data.name === 'owncloud-embed:select') {
    onLocationPick(event)
    return
  }

  if (event.data.name === 'owncloud-embed:cancel') {
    removeModal(modal.id)
  }
}

onMounted(() => {
  window.addEventListener('message', handleMessage)
})

onBeforeUnmount(() => {
  window.removeEventListener('message', handleMessage)
})
</script>

<style lang="scss">
.oc-modal.export-as-pdf-modal {
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

  .export-pdf-preview {
    h1,
    h2,
    h3,
    h4,
    h5,
    h6,
    p {
      color: var(--oc-color-text-default);
    }

    .md-editor-preview {
      hyphens: auto;
      overflow-wrap: break-word;
      word-break: normal;
      word-break: auto-phrase;
    }
  }
}
</style>
