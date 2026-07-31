<template>
  <div class="preview-details" :class="{ lightbox: isFullScreenModeActivated }">
    <div
      class="oc-background-brand oc-p-s oc-width-large oc-flex oc-flex-middle oc-flex-center oc-flex-around preview-controls-action-bar"
    >
      <oc-button
        v-oc-tooltip="$gettext('Show previous media file in folder')"
        class="preview-controls-previous"
        appearance="raw-inverse"
        variation="brand"
        :aria-label="$gettext('Show previous media file in folder')"
        @click="$emit('togglePrevious')"
      >
        <oc-icon size="large" name="arrow-drop-left" variation="inherit" />
      </oc-button>
      <p v-if="!isFolderLoading" class="oc-m-rm preview-controls-action-count">
        <span aria-hidden="true" v-text="ariaHiddenFileCount" />
        <span class="oc-invisible-sr" v-text="screenreaderFileCount" />
      </p>
      <oc-button
        v-oc-tooltip="$gettext('Show next media file in folder')"
        class="preview-controls-next"
        appearance="raw-inverse"
        variation="brand"
        :aria-label="$gettext('Show next media file in folder')"
        @click="$emit('toggleNext')"
      >
        <oc-icon size="large" name="arrow-drop-right" variation="inherit" />
      </oc-button>
      <div class="oc-flex">
        <oc-button
          v-oc-tooltip="
            isFullScreenModeActivated
              ? $gettext('Exit full screen mode')
              : $gettext('enter full screen mode')
          "
          class="preview-controls-fullscreen"
          appearance="raw-inverse"
          variation="brand"
          :aria-label="
            isFullScreenModeActivated
              ? $gettext('Exit full screen mode')
              : $gettext('enter full screen mode')
          "
          @click="$emit('toggleFullScreen')"
        >
          <oc-icon
            fill-type="line"
            :name="isFullScreenModeActivated ? 'fullscreen-exit' : 'fullscreen'"
            variation="inherit"
          />
        </oc-button>
      </div>
      <div v-if="showImageControls" class="oc-flex oc-flex-middle">
        <div class="oc-flex">
          <oc-button
            v-oc-tooltip="$gettext('Shrink the image')"
            class="preview-controls-image-shrink"
            appearance="raw-inverse"
            variation="brand"
            :aria-label="$gettext('Shrink the image')"
            @click="imageShrink"
          >
            <oc-icon fill-type="line" name="zoom-out" variation="inherit" />
          </oc-button>
          <oc-button
            v-oc-tooltip="$gettext('Show the image at its normal size')"
            class="preview-controls-image-original-size oc-ml-s oc-mr-s"
            appearance="raw-inverse"
            variation="brand"
            :aria-label="$gettext('Show the image at its normal size')"
            @click="$emit('setZoom', 1)"
          >
            <span v-text="currentZoomDisplayValue" />
          </oc-button>
          <oc-button
            v-oc-tooltip="$gettext('Enlarge the image')"
            class="preview-controls-image-zoom"
            appearance="raw-inverse"
            variation="brand"
            :aria-label="$gettext('Enlarge the image')"
            @click="imageZoom"
          >
            <oc-icon fill-type="line" name="zoom-in" variation="inherit" />
          </oc-button>
        </div>
        <div class="oc-ml-m">
          <oc-button
            v-oc-tooltip="$gettext('Rotate the image 90 degrees to the left')"
            class="preview-controls-rotate-left"
            appearance="raw-inverse"
            variation="brand"
            :aria-label="$gettext('Rotate the image 90 degrees to the left')"
            @click="imageRotateLeft"
          >
            <oc-icon fill-type="line" name="anticlockwise" variation="inherit" />
          </oc-button>
          <oc-button
            v-oc-tooltip="$gettext('Rotate the image 90 degrees to the right')"
            class="preview-controls-rotate-right"
            appearance="raw-inverse"
            variation="brand"
            :aria-label="$gettext('Rotate the image 90 degrees to the right')"
            @click="imageRotateRight"
          >
            <oc-icon fill-type="line" name="clockwise" variation="inherit" />
          </oc-button>
        </div>
        <div class="oc-ml-m">
          <oc-button
            v-oc-tooltip="$gettext('Reset')"
            class="preview-controls-image-reset"
            appearance="raw-inverse"
            variation="brand"
            :aria-label="$gettext('Reset')"
            @click="$emit('resetImage')"
          >
            <oc-icon fill-type="line" name="refresh" variation="inherit" />
          </oc-button>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { computed } from 'vue'
import { useGettext } from 'vue3-gettext'
import { Resource } from '@ownclouders/web-client'

interface Props {
  files: Resource[]
  activeIndex: number
  isFullScreenModeActivated?: boolean
  isFolderLoading?: boolean
  showImageControls?: boolean
  currentImageZoom?: number
  currentImageRotation?: number
}
interface Emits {
  (e: 'setRotation', rotation: number): void
  (e: 'setZoom', zoom: number): void
  (e: 'toggleFullScreen'): void
  (e: 'toggleNext'): void
  (e: 'togglePrevious'): void
  (e: 'resetImage'): void
}
const {
  files,
  activeIndex,
  isFullScreenModeActivated = false,
  isFolderLoading = false,
  showImageControls = false,
  currentImageZoom = 1,
  currentImageRotation = 0
} = defineProps<Props>()
const emit = defineEmits<Emits>()
const { $gettext } = useGettext()

const currentZoomDisplayValue = computed(() => {
  return `${(currentImageZoom * 100).toFixed(0)}%`
})

const ariaHiddenFileCount = computed(() => {
  return $gettext('%{ displayIndex } of %{ availableMediaFiles }', {
    displayIndex: (activeIndex + 1).toString(),
    availableMediaFiles: files.length.toString()
  })
})
const screenreaderFileCount = computed(() => {
  return $gettext('Media file %{ displayIndex } of %{ availableMediaFiles }', {
    displayIndex: (activeIndex + 1).toString(),
    availableMediaFiles: files.length.toString()
  })
})

const calculateZoom = (zoom: number, factor: number) => {
  return Math.round(zoom * factor * 20) / 20
}
const imageShrink = () => {
  emit('setZoom', Math.max(0.1, calculateZoom(currentImageZoom, 0.8)))
}
const imageZoom = () => {
  const maxZoomValue = calculateZoom(9, 1.25)
  emit('setZoom', Math.min(calculateZoom(currentImageZoom, 1.25), maxZoomValue))
}
const imageRotateLeft = () => {
  emit('setRotation', currentImageRotation === -270 ? 0 : currentImageRotation - 90)
}
const imageRotateRight = () => {
  emit('setRotation', currentImageRotation === 270 ? 0 : currentImageRotation + 90)
}
</script>

<style lang="scss" scoped>
.preview-details.lightbox {
  z-index: 1000;
  opacity: 0.9;
}

.preview-controls-action-count {
  color: var(--oc-color-swatch-brand-contrast);
}

.preview-controls-image-original-size {
  width: 42px;
}
</style>
