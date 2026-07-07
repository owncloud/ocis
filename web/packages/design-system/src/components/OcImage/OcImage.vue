<template>
  <img :src="src" :alt="alt" :aria-hidden="`${ariaHidden}`" :title="title" :loading="loadingType" />
</template>
<script lang="ts" setup>
import { computed } from 'vue'

/**
 * OcImg Component
 *
 * Displaying images with support for lazy loading, accessibility, and additional metadata.
 *
 * @component
 * @name OcImg
 * @status ready
 * @release 1.0.0
 *
 * @props
 * @prop {string} [src=''] - The source URL of the image.
 * @prop {string} [alt=''] - The alt text for the image, used for accessibility.
 * @prop {string} [title=''] - The title of the image, displayed on hover.
 * @prop {'lazy' | 'eager'} [loadingType='eager'] - Specifies whether the image should be loaded lazily or eagerly.
 *
 * @example
 * <h3>Example image, with width set explicitly:</h3>
 * <oc-img width="200" height="300" src="https://example.com/image" title="I am random" alt="example image" />
 *
 * <h3>Example grayscale image:</h3>
 * <oc-img src="https://example.com/image" title="Grayscale" alt="example grayscale image" />
 *
 * <h3>Example using an SVG file:</h3>
 * <oc-img width="100" src="icons/folder.svg" title="I am a folder" alt="folder icon" />
 */

interface Props {
  src?: string
  alt?: string
  title?: string
  loadingType?: 'lazy' | 'eager'
}
defineOptions({
  name: 'OcImg',
  status: 'ready',
  release: '1.0.0'
})
const { src = null, alt = '', title = null, loadingType = 'eager' } = defineProps<Props>()

const ariaHidden = computed(() => alt.length === 0)
</script>
