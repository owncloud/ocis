<template>
  <div class="audio-container oc-flex oc-flex-column">
    <audio :key="`media-audio-${file.id}`" controls preload="preload" :autoplay="isAutoPlayEnabled">
      <source :src="file.url" :type="file.mimeType" />
    </audio>
    <p v-if="audioText" class="oc-text-muted oc-text-small" v-text="audioText"></p>
  </div>
</template>
<script lang="ts" setup>
import { computed } from 'vue'
import { CachedFile } from '../../helpers/types'
import { Resource } from '@ownclouders/web-client'

interface Props {
  file: CachedFile
  resource: Resource
  isAutoPlayEnabled?: boolean
}
const { file, resource, isAutoPlayEnabled = true } = defineProps<Props>()
const audioText = computed(() => {
  if (resource.audio?.artist && resource.audio?.title) {
    return `${resource.audio.artist} - ${resource.audio.title}`
  }
  return ''
})
</script>
<style lang="scss" scoped>
.audio-container {
  width: 300px;
}
</style>
