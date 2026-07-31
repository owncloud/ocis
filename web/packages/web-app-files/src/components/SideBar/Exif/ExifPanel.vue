<template>
  <div id="files-sidebar-panel-exif">
    <dl class="exif-data-list">
      <dt v-text="$gettext('Dimensions')" />
      <dd data-testid="exif-panel-dimensions" v-text="dimensions" />
      <dt v-text="$gettext('Device make')" />
      <dd data-testid="exif-panel-cameraMake" v-text="cameraMake" />
      <dt v-text="$gettext('Device model')" />
      <dd data-testid="exif-panel-cameraModel" v-text="cameraModel" />
      <dt v-text="$gettext('Focal length')" />
      <dd data-testid="exif-panel-focalLength" v-text="focalLength" />
      <dt v-text="$gettext('F number')" />
      <dd data-testid="exif-panel-fNumber" v-text="fNumber" />
      <dt v-text="$gettext('Exposure time')" />
      <dd data-testid="exif-panel-exposureTime" v-text="exposureTime" />
      <dt v-text="$gettext('ISO')" />
      <dd data-testid="exif-panel-iso" v-text="iso" />
      <dt v-text="$gettext('Orientation')" />
      <dd data-testid="exif-panel-orientation" v-text="orientation" />
      <dt v-text="$gettext('Taken time')" />
      <dd data-testid="exif-panel-takenDateTime" v-text="takenDateTime" />
      <dt v-text="$gettext('Location')" />
      <dd v-if="isCopyToClipboardAvailable" data-testid="exif-panel-location">
        <span>{{ location }}</span>
        <oc-button
          v-if="location"
          v-oc-tooltip="copyLocationToClipboardLabel"
          size="small"
          appearance="raw"
          class="oc-ml-s"
          :aria-label="copyLocationToClipboardLabel"
          @click="copyLocationToClipboard"
        >
          <oc-icon size="small" :name="isCopiedToClipboard ? 'checkbox-circle' : 'file-copy'" />
        </oc-button>
      </dd>
      <dd v-else data-testid="exif-panel-location" v-text="location" />
    </dl>
  </div>
</template>

<script lang="ts" setup>
import { computed, inject, Ref, unref } from 'vue'
import { Resource } from '@ownclouders/web-client'
import { formatDateFromISO, useMessages } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { useClipboard } from '@vueuse/core'

const resource = inject<Ref<Resource>>('resource')
const language = useGettext()
const { $gettext } = language
const { showMessage } = useMessages()
const {
  copy: copyToClipboard,
  copied: isCopiedToClipboard,
  isSupported: isCopyToClipboardSupported
} = useClipboard({ legacy: true, copiedDuring: 550 })

const dimensions = computed(() => {
  const image = unref(resource).image
  const width = image?.width
  const height = image?.height
  if (!width || !height) {
    return '-'
  }
  if ([5, 6, 7, 8].includes(unref(resource).photo?.orientation)) {
    // these orientations indicate portrait mode. tika normalizes width and height according to orientation.
    return `${height}x${width}`
  }
  return `${width}x${height}`
})

const cameraMake = computed(() => {
  return unref(resource).photo?.cameraMake || '-'
})

const cameraModel = computed(() => {
  return unref(resource).photo?.cameraModel || '-'
})

const focalLength = computed(() => {
  const photo = unref(resource).photo
  return photo?.focalLength ? `${photo.focalLength} mm` : '-'
})

const fNumber = computed(() => {
  const photo = unref(resource).photo
  return photo?.fNumber ? `f/${photo.fNumber}` : '-'
})

const exposureTime = computed(() => {
  const photo = unref(resource).photo
  return photo?.exposureDenominator
    ? `${photo.exposureNumerator}/${photo.exposureDenominator}`
    : '-'
})

const iso = computed(() => {
  return unref(resource).photo?.iso || '-'
})

const orientation = computed(() => {
  return unref(resource).photo?.orientation || '-'
})

const takenDateTime = computed(() => {
  const photo = unref(resource).photo
  return photo?.takenDateTime ? formatDateFromISO(photo.takenDateTime, language.current) : '-'
})

const location = computed(() => {
  const l = unref(resource).location
  if (!l?.latitude || !l?.longitude) {
    return '-'
  }
  return `${l.latitude}, ${l.longitude}`
})

const isCopyToClipboardAvailable = computed(() => {
  if (!unref(isCopyToClipboardSupported)) {
    return false
  }
  const l = unref(resource).location
  return l?.latitude && l?.longitude
})
const copyLocationToClipboard = () => {
  copyToClipboard(unref(location))
  showMessage({
    title: $gettext('The location has been copied to your clipboard.')
  })
}
const copyLocationToClipboardLabel = computed(() => {
  return $gettext('Copy location to clipboard')
})
</script>

<style lang="scss">
.exif-data-list {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);

  dt,
  dd {
    margin-bottom: var(--oc-space-small);
  }

  dt {
    font-weight: bold;
    white-space: nowrap;
  }

  dd {
    margin-inline-start: var(--oc-space-medium);
  }
}
</style>
