<template>
  <div class="app-details oc-card oc-card-default oc-card-rounded">
    <div class="oc-p-xs">
      <router-link :to="{ name: `${APPID}-list` }" class="oc-flex oc-flex-middle app-details-back">
        <oc-icon name="arrow-left-s" fill-type="line" />
        <span v-text="$gettext('Back to list')" />
      </router-link>
    </div>
    <app-image-gallery :app="app" :show-pagination="true" />
    <div class="app-content oc-card-body oc-p">
      <div class="oc-flex oc-flex-middle">
        <h2 class="oc-my-s oc-text-truncate app-details-title">{{ app.name }}</h2>
        <span class="oc-ml-s oc-text-muted oc-text-small oc-mt-s">
          v{{ app.mostRecentVersion.version }}
        </span>
      </div>
      <p class="oc-my-rm">{{ app.subtitle }}</p>
      <div v-if="app.description">
        <h3>{{ $gettext('Details') }}</h3>
        <text-editor
          class="oc-my-s"
          :is-read-only="true"
          :markdown-mode="true"
          :current-content="app.description"
        />
      </div>
      <div v-if="app.tags">
        <h3>{{ $gettext('Tags') }}</h3>
        <app-tags :app="app" @click="onTagClicked" />
      </div>
      <div v-if="app.authors">
        <h3>{{ $gettext('Author') }}</h3>
        <app-authors :app="app" />
      </div>
      <div v-if="app.resources">
        <h3>{{ $gettext('Resources') }}</h3>
        <app-resources :app="app" />
      </div>
      <div v-if="app.versions">
        <h3>
          {{ $gettext('Releases') }}
          <app-contextual-helper />
        </h3>
        <app-versions :app="app" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import { App } from '../types'
import { APPID } from '../appid'
import { TextEditor, useRouteParam, useRouter } from '@ownclouders/web-pkg'
import { useAppsStore } from '../piniaStores'
import AppResources from '../components/AppResources.vue'
import AppTags from '../components/AppTags.vue'
import AppVersions from '../components/AppVersions.vue'
import AppAuthors from '../components/AppAuthors.vue'
import AppImageGallery from '../components/AppImageGallery.vue'
import AppContextualHelper from '../components/AppContextualHelper.vue'

const appIdRouteParam = useRouteParam('appId')
const appId = computed(() => {
  return decodeURIComponent(unref(appIdRouteParam))
})
const appsStore = useAppsStore()
const router = useRouter()

const app = computed<App>(() => {
  return appsStore.getById(unref(appId))
})

const onTagClicked = (tag: string) => {
  router.push({ name: `${APPID}-list`, query: { filter: tag } })
}
</script>

<style lang="scss">
// .oc-my-s > .md-editor.md-editor-dark.md-editor-previewOnly {
//   max-width: 100%;
// }
.app-details {
  background-color: var(--oc-color-background-highlight);
  box-shadow: none;
  max-width: 600px;
  margin: 0 auto;
  outline: 1px solid var(--oc-color-border);

  .app-content {
    display: flex;
    flex-flow: column;
    gap: 1rem;
  }
}
</style>
