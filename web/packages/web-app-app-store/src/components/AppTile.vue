<template>
  <li class="app-tile oc-card oc-card-default oc-card-rounded">
    <router-link :to="{ name: `${APPID}-details`, params: { appId: encodeURIComponent(app.id) } }">
      <app-image-gallery :app="app" />
    </router-link>
    <div class="app-tile-body oc-card-body oc-p">
      <div class="app-content">
        <div class="oc-flex oc-flex-middle">
          <h3 class="oc-my-s oc-text-truncate mark-element app-tile-title">
            <router-link
              :to="{ name: `${APPID}-details`, params: { appId: encodeURIComponent(app.id) } }"
            >
              {{ app.name }}
            </router-link>
          </h3>
          <span class="oc-ml-s oc-text-muted oc-text-small oc-mt-xs">
            v{{ app.mostRecentVersion.version }}
          </span>
        </div>
        <p class="oc-my-s mark-element">{{ app.subtitle }}</p>
      </div>
      <app-tags :app="app" @click="emitSearchTerm" />
      <app-actions :app="app" />
    </div>
  </li>
</template>

<script lang="ts" setup>
import { App } from '../types'
import { APPID } from '../appid'
import AppTags from './AppTags.vue'
import AppActions from './AppActions.vue'
import AppImageGallery from './AppImageGallery.vue'

interface Props {
  app?: App
}
interface Emits {
  (e: 'search', value: string): void
}
const { app = undefined } = defineProps<Props>()
const emit = defineEmits<Emits>()
const emitSearchTerm = (term: string) => {
  emit('search', term)
}
</script>

<style lang="scss">
.app-tile {
  overflow: hidden;
  background-color: var(--oc-color-background-highlight) !important;
  box-shadow: none;
  height: 100%;
  display: flex;
  flex-flow: column;
  outline: 1px solid var(--oc-color-border);

  .app-tile-body {
    display: flex;
    flex-flow: column;
    justify-content: space-between;
    height: 100%;
  }

  .app-tile-title {
    .mark-highlight {
      font-weight: unset !important;
      color: var(--oc-color-text-default);
    }
  }
}
</style>
