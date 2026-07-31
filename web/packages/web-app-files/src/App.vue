<template>
  <main id="files" class="oc-flex oc-height-1-1">
    <div v-if="dragareaEnabled" class="dragarea" />
    <router-view tabindex="0" class="files-wrapper oc-width-expand oc-height-1-1 oc-flex-wrap" />
  </main>
</template>
<script lang="ts" setup>
import { onBeforeUnmount, watch, ref, computed, unref } from 'vue'
import { useRoute, eventBus, useResourcesStore } from '@ownclouders/web-pkg'

const route = useRoute()
const { resetSelection } = useResourcesStore()

const dragareaEnabled = ref(false)

const routePath = computed(() => unref(route).path)

watch(routePath, () => {
  resetSelection()
})

const hideDropzone = () => {
  dragareaEnabled.value = false
}
const onDragOver = (event: DragEvent) => {
  dragareaEnabled.value = (event.dataTransfer.types || []).some((e) => e === 'Files')
}

const dragOver = eventBus.subscribe('drag-over', onDragOver)
const dragOut = eventBus.subscribe('drag-out', hideDropzone)
const drop = eventBus.subscribe('drop', hideDropzone)

onBeforeUnmount(() => {
  eventBus.unsubscribe('drag-over', dragOver)
  eventBus.unsubscribe('drag-out', dragOut)
  eventBus.unsubscribe('drop', drop)
})
</script>

<style lang="scss" scoped>
main {
  max-height: 100%;
}

.dragarea {
  background-color: rgba(60, 130, 225, 0.21);
  pointer-events: none;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  position: absolute;
  z-index: 9;
  border-radius: 14px;
  border: 2px dashed var(--oc-color-swatch-primary-muted);
}

#files {
  position: relative;
}
</style>
