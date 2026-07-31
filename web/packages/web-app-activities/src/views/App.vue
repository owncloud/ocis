<template>
  <h1 v-text="$gettext('Activities')" />
  <div class="oc-width-1-1 oc-mb-m">
    <item-filter
      ref="mediaTypeFilter"
      :allow-multiple="false"
      :filter-label="$gettext('Location')"
      :filterable-attributes="['name']"
      :option-filter-label="$gettext('Filter location')"
      :show-option-filter="true"
      :items="filterableSpaces"
      :close-on-click="true"
      class="files-search-filter-file-type oc-mr-s"
      display-name-attribute="name"
      filter-name="location"
    >
      <template #image="{ item }">
        <oc-icon :name="getLocationFilterIcon(item)" />
      </template>
      <template #item="{ item }">
        <div v-text="item.name" />
      </template>
    </item-filter>
  </div>
  <app-loading-spinner v-if="isLoading" />
  <template v-else>
    <no-content-message v-if="!activities.length" icon="pulse">
      <template #message>
        <span v-text="$gettext('No activities found')" />
      </template>
    </no-content-message>
    <ActivityList v-else :activities="activities" />
  </template>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, unref, watch } from 'vue'
import {
  AppLoadingSpinner,
  ItemFilter,
  NoContentMessage,
  useClientService,
  useRouteQuery,
  useSpacesStore
} from '@ownclouders/web-pkg'
import { storeToRefs } from 'pinia'
import {
  call,
  isPersonalSpaceResource,
  isProjectSpaceResource,
  SpaceResource
} from '@ownclouders/web-client'
import { useTask } from 'vue-concurrency'
import { Activity } from '@ownclouders/web-client/graph/generated'
import ActivityList from './ActivityList.vue'

const spacesStore = useSpacesStore()
const { spaces } = storeToRefs(spacesStore)
const clientService = useClientService()
const activities = ref<Activity[]>([])

const locationQuery = useRouteQuery('q_location')

const filterableSpaces = computed(() => {
  return [...unref(spaces)]
    .filter(
      (space) =>
        !space.disabled && (isProjectSpaceResource(space) || isPersonalSpaceResource(space))
    )
    .sort((a, b) => {
      if (isPersonalSpaceResource(a) === isPersonalSpaceResource(b)) {
        return a.name.localeCompare(b.name)
      }
      return isPersonalSpaceResource(a) ? -1 : 1
    })
})

const loadActivitiesTask = useTask(function* (signal) {
  const filters = ['sort:desc', 'limit:100']

  if (unref(locationQuery)) {
    filters.push(`itemid:${unref(locationQuery)}`)
  }

  activities.value = yield* call(
    clientService.graphAuthenticated.activities.listActivities(filters.join(' AND '), {
      signal
    })
  )
})

const isLoading = computed(() => loadActivitiesTask.isRunning || !loadActivitiesTask.last)

const getLocationFilterIcon = (space: SpaceResource) => {
  if (isPersonalSpaceResource(space)) {
    return 'folder'
  }

  return 'layout-grid'
}

onMounted(() => {
  loadActivitiesTask.perform()
})

watch(locationQuery, () => {
  loadActivitiesTask.perform()
})
</script>
