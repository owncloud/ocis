<template>
  <div class="app-list oc-mb-m">
    <h2 class="oc-mt-rm app-list-headline">
      {{ $gettext('App Store') }}
      <app-contextual-helper />
    </h2>
    <div class="oc-flex oc-flex-middle">
      <oc-text-input
        id="apps-filter"
        :model-value="filterTermInput"
        :label="$gettext('Search')"
        :clear-button-enabled="true"
        autocomplete="off"
        @update:model-value="setFilterTerm"
      />
    </div>
    <no-content-message v-if="!filteredApps.length" icon="store">
      <template #message>
        <span v-text="$gettext('No apps found matching your search')" />
      </template>
    </no-content-message>
    <oc-list v-else class="app-tiles">
      <app-tile
        v-for="app in filteredApps"
        :key="`app-${app.repository.name}-${app.id}`"
        :app="app"
        class="oc-my-m"
        @search="setFilterTerm"
      />
    </oc-list>
  </div>
</template>

<script lang="ts" setup>
import { computed, nextTick, onMounted, ref, unref, watch } from 'vue'
import Mark from 'mark.js'
import Fuse from 'fuse.js'
import { useAppsStore } from '../piniaStores'
import AppTile from '../components/AppTile.vue'
import { storeToRefs } from 'pinia'
import { App } from '../types'
import {
  defaultFuseOptions,
  NoContentMessage,
  queryItemAsString,
  useRouteQuery,
  useRouter
} from '@ownclouders/web-pkg'
import AppContextualHelper from '../components/AppContextualHelper.vue'

const appsStore = useAppsStore()
const { apps } = storeToRefs(appsStore)
const router = useRouter()

const filterTermQueryItem = useRouteQuery('filter', '')
const filterTerm = computed(() => {
  return queryItemAsString(unref(filterTermQueryItem))
})
const filterTermInput = ref('')

const setFilterTerm = (term: string) => {
  return router.replace({
    query: {
      ...(term && { filter: term.trim() })
    }
  })
}
const filter = (apps: App[], filterTerm: string) => {
  if (!(filterTerm || '').trim()) {
    return apps
  }
  const searchEngine = new Fuse(apps, {
    ...defaultFuseOptions,
    keys: ['name', 'subtitle', 'tags']
  })
  return searchEngine.search(filterTerm).map((r) => r.item)
}
const filteredApps = computed(() => {
  // TODO: debounce the filtering by 100-300ms
  return filter(unref(apps), unref(filterTerm))
})

const markInstance = ref<Mark>(null)
onMounted(async () => {
  await nextTick()
  markInstance.value = new Mark('.mark-element')
})
watch([filterTerm, markInstance], () => {
  filterTermInput.value = unref(filterTerm)
  unref(markInstance)?.unmark()
  if (unref(filterTerm)) {
    unref(markInstance)?.mark(unref(filterTerm), {
      element: 'span',
      className: 'mark-highlight'
    })
  }
})
</script>

<style lang="scss">
.app-list {
  .app-tiles {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 2rem;
  }
}
</style>
