<template>
  <picture>
    <source v-if="src.xl" key="responsive-image-xl" media="(min-width: 1600px)" :srcset="src.xl" />
    <source v-if="src.lg" key="responsive-image-lg" media="(min-width: 1200px)" :srcset="src.lg" />
    <source v-if="src.md" key="responsive-image-md" media="(min-width: 960px)" :srcset="src.md" />
    <source v-if="src.sm" key="responsive-image-sm" media="(min-width: 640px)" :srcset="src.sm" />
    <source v-if="src.xs" key="responsive-image-xs" :srcset="src.xs" />

    <img :src="defaultSrc" :alt="alt" :loading="loading" />
  </picture>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineOptions({ name: 'OcResponsiveImage' })

type ResponsiveImageSources = {
  xs?: string
  sm?: string
  md?: string
  lg?: string
  xl?: string
}

const {
  src,
  alt,
  loading = 'lazy'
} = defineProps<{
  src: ResponsiveImageSources
  alt: string
  loading?: HTMLImageElement['loading']
}>()

const defaultSrc = computed(() => src.xl || src.lg || src.md || src.sm || src.xs || '')
</script>
