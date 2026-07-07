<template>
  <div class="app-image-wrapper">
    <div v-if="app.badge" class="app-image-ribbon" :class="[`app-image-ribbon-${app.badge.color}`]">
      <span>{{ app.badge.label }}</span>
    </div>
    <div class="app-image">
      <oc-img v-if="currentImage?.url" :src="currentImage?.url" :alt="app.name" />
      <div v-else class="fallback-icon">
        <oc-icon name="computer" size="xxlarge" />
      </div>
    </div>
    <ul v-if="hasPagination" class="app-image-navigation">
      <li>
        <oc-button
          data-testid="prev-image"
          class="oc-p-xs"
          appearance="raw"
          variation="primary"
          @click="previousImage"
        >
          <oc-icon name="arrow-left-s" />
        </oc-button>
      </li>
      <li v-for="(image, index) in images" :key="`gallery-page-${index}`">
        <oc-button
          data-testid="set-image"
          class="oc-py-xs"
          appearance="raw"
          variation="primary"
          @click="setImageIndex(index)"
        >
          <oc-icon
            name="circle"
            size="small"
            :fill-type="index === currentImageIndex ? 'fill' : 'line'"
          />
        </oc-button>
      </li>
      <li>
        <oc-button
          data-testid="next-image"
          class="oc-p-xs"
          appearance="raw"
          variation="primary"
          @click="nextImage"
        >
          <oc-icon name="arrow-right-s" />
        </oc-button>
      </li>
    </ul>
  </div>
</template>
<script lang="ts" setup>
import { computed, ref, unref } from 'vue'
import { App, AppImage } from '../types'

interface Props {
  app?: App
  showPagination?: boolean
}
const { app = undefined, showPagination = false } = defineProps<Props>()
const images = computed(() => {
  return [app.coverImage, ...app.screenshots]
})

const currentImageIndex = ref<number>(0)
const currentImage = computed<AppImage>(() => unref(images)[unref(currentImageIndex)])
const hasPagination = computed(() => showPagination && unref(images).length > 1)
const nextImage = () => {
  currentImageIndex.value = (unref(currentImageIndex) + 1) % unref(images).length
}
const previousImage = () => {
  currentImageIndex.value =
    (unref(currentImageIndex) - 1 + unref(images).length) % unref(images).length
}
const setImageIndex = (index: number) => {
  currentImageIndex.value = index
}
</script>

<style lang="scss">
.app-image-wrapper {
  position: relative;

  .app-image-ribbon {
    position: absolute;
    top: 0;
    right: 0;
    z-index: 1;
    overflow: hidden;
    width: 7rem;
    height: 7rem;
    text-align: right;

    &-primary {
      span {
        color: var(--oc-color-swatch-primary-contrast);
        background-color: var(--oc-color-swatch-primary-default);
      }
    }
    &-success {
      span {
        color: var(--oc-color-swatch-success-contrast);
        background-color: var(--oc-color-swatch-success-default);
      }
    }
    &-danger {
      span {
        color: var(--oc-color-swatch-danger-contrast);
        background-color: var(--oc-color-swatch-danger-default);
      }
    }

    span {
      position: absolute;
      top: 1.8rem;
      right: -2.2rem;
      font-size: 0.7rem;
      font-weight: bold;
      text-align: center;
      line-height: 2rem;
      transform: rotate(45deg);
      -webkit-transform: rotate(45deg);
      width: 10rem;
      display: block;
    }
  }

  .app-image {
    width: 100%;

    img {
      width: 100%;
      max-width: 100%;
      aspect-ratio: 3/2;
      object-fit: cover;
    }

    .fallback-icon {
      width: 100%;
      aspect-ratio: 3/2;
      background-color: white;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }

  .app-image-navigation {
    list-style: none;
    width: 100%;
    position: absolute;
    bottom: 0;
    display: flex;
    flex-direction: row;
    gap: 0.5rem;
    align-items: center;
    justify-content: center;
    padding: var(--oc-space-small) 0;
    margin: 0;
    background-color: rgba(255, 255, 255, 0.8);
  }
}
</style>
