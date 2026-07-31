<template>
  <div class="oc-height-1-1" tabindex="0">
    <app-loading-spinner v-if="isLoading" />
    <iframe
      v-show="!isLoading"
      ref="iframeRef"
      class="oc-width-1-1 oc-height-1-1"
      :title="iframeTitle"
      :src="iframeSrc"
      tabindex="0"
      @load="onLoad"
    ></iframe>
  </div>
</template>

<script lang="ts" setup>
import { onBeforeUnmount, onMounted, ref } from 'vue'
import {
  EDITOR_MODE_EDIT,
  Modal,
  useEmbedMode,
  useGetMatchingSpace,
  useModals,
  useRouter,
  useThemeStore,
  useFileActions,
  embedModeFilePickMessageData
} from '../../composables'
import { ApplicationInformation } from '../../apps'
import { RouteLocationRaw } from 'vue-router'
import AppLoadingSpinner from '../AppLoadingSpinner.vue'
import { isShareSpaceResource } from '@ownclouders/web-client'
import { unref } from 'vue'

interface Props {
  modal: Modal
  app: ApplicationInformation
  parentFolderLink: RouteLocationRaw
}
const { modal, app, parentFolderLink } = defineProps<Props>()
const iframeRef = ref<HTMLIFrameElement>()
const isLoading = ref(true)
const router = useRouter()
const { removeModal } = useModals()
const { getMatchingSpace } = useGetMatchingSpace()
const themeStore = useThemeStore()
const { getEditorRouteOpts } = useFileActions()
const { verifyMessageOrigin } = useEmbedMode()
const parentFolderRoute = router.resolve(parentFolderLink)

const availableFileTypes = (app as ApplicationInformation).extensions.map((e) =>
  e.extension ? e.extension : e.mimeType
)

const iframeTitle = themeStore.currentTheme.common?.name
const iframeUrl = new URL(parentFolderRoute.href, window.location.origin)
iframeUrl.searchParams.append('hide-logo', 'true')
iframeUrl.searchParams.append('embed', 'true')
iframeUrl.searchParams.append('embed-target', 'file')
iframeUrl.searchParams.append('embed-delegate-authentication', 'false')
iframeUrl.searchParams.append('embed-file-types', availableFileTypes.join(','))

const iframeSrc = iframeUrl.href

const onLoad = () => {
  isLoading.value = false
  unref(iframeRef).contentWindow.focus()
}

const onFilePick = ({ data, origin }: MessageEvent) => {
  if (!verifyMessageOrigin(origin)) {
    return
  }

  if (data.name !== 'owncloud-embed:file-pick') {
    return
  }

  const { resource, locationQuery }: embedModeFilePickMessageData = data.data

  const space = getMatchingSpace(resource)
  const remoteItemId = isShareSpaceResource(space) ? space.id : undefined

  const routeOpts = getEditorRouteOpts(
    unref(router.currentRoute).name,
    space,
    resource,
    EDITOR_MODE_EDIT,
    remoteItemId
  )
  routeOpts.query = { ...routeOpts.query, ...locationQuery }

  const editorRoute = router.resolve(routeOpts)
  const editorRouteUrl = new URL(editorRoute.href, window.location.origin)

  removeModal(modal.id)
  window.open(editorRouteUrl.href, '_blank')
}

const onCancel = ({ data, origin }: MessageEvent) => {
  if (!verifyMessageOrigin(origin)) {
    return
  }

  if (data.name !== 'owncloud-embed:cancel') {
    return
  }

  removeModal(modal.id)
}

onMounted(() => {
  window.addEventListener('message', onFilePick)
  window.addEventListener('message', onCancel)
})

onBeforeUnmount(() => {
  window.removeEventListener('message', onFilePick)
  window.removeEventListener('message', onCancel)
})
</script>

<style lang="scss">
.oc-modal.open-with-app-modal {
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
}
</style>
