<template>
  <div
    v-oc-tooltip="
      isRemoteUploadInProgress ? $gettext('Please wait until all imports have finished') : null
    "
  >
    <slot :trigger-upload="triggerUpload" :upload-label-id="uploadLabelId">
      <oc-button
        :class="btnClass"
        justify-content="left"
        appearance="raw"
        :disabled="isRemoteUploadInProgress"
        @click="triggerUpload"
      >
        <resource-icon :resource="resource" size="medium" />
        <span :id="uploadLabelId">{{ buttonLabel }}</span>
      </oc-button>
    </slot>
    <input
      :id="inputId"
      ref="input"
      v-bind="inputAttrs"
      class="upload-input-target"
      type="file"
      :aria-labelledby="uploadLabelId"
      :name="isFolder ? 'file' : 'folder'"
      tabindex="-1"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, onBeforeUnmount, ref, useTemplateRef } from 'vue'
import { Resource } from '@ownclouders/web-client'
import { useService, ResourceIcon } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import type { UppyService } from '@ownclouders/web-pkg'

interface Props {
  btnLabel?: string
  btnClass?: string
  isFolder?: boolean
}
const { btnLabel = '', btnClass = '', isFolder = false } = defineProps<Props>()
const input = useTemplateRef('input')
const language = useGettext()
const { $gettext } = language

const uppyService = useService<UppyService>('$uppyService')
const isRemoteUploadInProgress = ref(uppyService.isRemoteUploadInProgress())

let uploadStartedSub: string
let uploadCompletedSub: string

const resource = computed(() => {
  return { extension: '', isFolder: isFolder } as Resource
})

const onUploadStarted = () =>
  (isRemoteUploadInProgress.value = uppyService.isRemoteUploadInProgress())
const onUploadCompleted = () => (isRemoteUploadInProgress.value = false)

onMounted(() => {
  uploadStartedSub = uppyService.subscribe('uploadStarted', onUploadStarted)
  uploadCompletedSub = uppyService.subscribe('uploadCompleted', onUploadCompleted)
  uppyService.registerUploadInput(input.value as HTMLInputElement)
})

onBeforeUnmount(() => {
  uppyService.unsubscribe('uploadStarted', uploadStartedSub)
  uppyService.unsubscribe('uploadCompleted', uploadCompletedSub)
  uppyService.removeUploadInput(input.value as HTMLInputElement)
})

function triggerUpload() {
  ;(input.value as HTMLInputElement).click()
}
const inputId = computed(() => {
  if (isFolder) {
    return 'files-folder-upload-input'
  }
  return 'files-file-upload-input'
})
const uploadLabelId = computed(() => {
  if (isFolder) {
    return 'files-folder-upload-button'
  }
  return 'files-file-upload-button'
})
const buttonLabel = computed(() => {
  if (btnLabel) {
    return btnLabel
  }
  if (isFolder) {
    return $gettext('Folder')
  }
  return $gettext('Files')
})
const inputAttrs = computed(() => {
  if (isFolder) {
    return {
      webkitdirectory: true,
      mozdirectory: true,
      allowdirs: true
    }
  }
  return { multiple: true }
})
</script>

<style scoped>
.upload-input-target {
  position: absolute;
  left: -99999px;
}
</style>
