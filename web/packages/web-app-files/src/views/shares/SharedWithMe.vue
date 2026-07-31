<template>
  <div class="oc-flex">
    <files-view-wrapper class="oc-flex-column">
      <app-bar :has-bulk-actions="true" :is-side-bar-open="isSideBarOpen">
        <template #navigation>
          <SharesNavigation />
        </template>
      </app-bar>
      <app-loading-spinner v-if="areResourcesLoading" />
      <template v-else>
        <oc-hidden-announcer :announcement="visibilityFilterAnnouncement" level="polite" />
        <div
          class="shared-with-me-filters oc-flex oc-flex-between oc-flex-wrap oc-flex-bottom oc-mx-m oc-mb-m"
        >
          <fieldset class="oc-flex oc-flex-wrap">
            <div>
              <legend class="oc-mr-m oc-flex oc-flex-middle">
                <oc-icon name="filter-2" class="oc-mr-xs" />
                <span v-text="$gettext('Filter:')" />
              </legend>
            </div>
            <item-filter-inline
              class="share-visibility-filter"
              filter-name="share-visibility"
              :filter-options="visibilityOptions"
              @toggle-filter="setAreHiddenFilesShown"
            />
            <item-filter
              :allow-multiple="true"
              :filter-label="$gettext('Share Type')"
              :filterable-attributes="['label']"
              :items="shareTypes"
              :option-filter-label="$gettext('Filter share types')"
              :show-option-filter="true"
              id-attribute="key"
              class="share-type-filter oc-ml-s"
              display-name-attribute="label"
              filter-name="shareType"
            >
              <template #item="{ item }">
                <span class="oc-ml-s" v-text="item.label" />
              </template>
            </item-filter>
            <item-filter
              :allow-multiple="true"
              :filter-label="$gettext('Shared By')"
              :filterable-attributes="['displayName']"
              :items="fileOwners"
              :option-filter-label="$gettext('Filter shared by')"
              :show-option-filter="true"
              id-attribute="id"
              class="shared-by-filter oc-ml-s"
              display-name-attribute="displayName"
              filter-name="sharedBy"
            >
              <template #image="{ item }">
                <avatar-image :width="32" :userid="item.id" :user-name="item.displayName" />
              </template>
              <template #item="{ item }">
                <span class="oc-ml-s" v-text="item.displayName" />
              </template>
            </item-filter>
          </fieldset>
          <div>
            <oc-text-input
              v-model="filterTerm"
              class="search-filter"
              :label="$gettext('Search')"
              autocomplete="off"
            />
          </div>
        </div>
        <shared-with-me-section
          id="files-shared-with-me-view"
          :file-list-header-y="fileListHeaderY"
          :items="items"
          :is-side-bar-open="isSideBarOpen"
          :sort-by="sortBy"
          :sort-dir="sortDir"
          :sort-handler="handleSort"
          :title="shareSectionTitle"
          :empty-message="
            areHiddenFilesShown ? $gettext('No hidden shares') : $gettext('No shares')
          "
          :grouping-settings="groupingSettings"
        />
      </template>
    </files-view-wrapper>
    <file-side-bar
      :is-open="isSideBarOpen"
      :active-panel="sideBarActivePanel"
      :space="selectedShareSpace"
    />
  </div>
</template>

<script lang="ts" setup>
import Fuse from 'fuse.js'
import Mark from 'mark.js'
import { useResourcesViewDefaults } from '../../composables'

import {
  AppLoadingSpinner,
  FileSideBar,
  InlineFilterOption,
  ItemFilter,
  useAppsStore,
  useResourcesStore
} from '@ownclouders/web-pkg'
import { AppBar, ItemFilterInline } from '@ownclouders/web-pkg'
import { queryItemAsString, useRouteQuery } from '@ownclouders/web-pkg'
import SharedWithMeSection from '../../components/Shares/SharedWithMeSection.vue'
import { computed, onMounted, ref, unref, watch } from 'vue'
import FilesViewWrapper from '../../components/FilesViewWrapper.vue'
import { useGetMatchingSpace, useSort } from '@ownclouders/web-pkg'
import { useGroupingSettings } from '@ownclouders/web-pkg'
import SharesNavigation from '../../components/AppBar/SharesNavigation.vue'
import { useGettext } from 'vue3-gettext'
import { useOpenWithDefaultApp, defaultFuseOptions } from '@ownclouders/web-pkg'
import { IncomingShareResource, ShareTypes } from '@ownclouders/web-client'
import { uniq } from 'lodash-es'

const { openWithDefaultApp } = useOpenWithDefaultApp()
const appsStore = useAppsStore()
const resourcesStore = useResourcesStore()

const {
  areResourcesLoading,
  sortFields,
  fileListHeaderY,
  loadResourcesTask,
  selectedResources,
  sideBarActivePanel,
  isSideBarOpen,
  paginatedResources,
  scrollToResourceFromRoute
} = useResourcesViewDefaults<IncomingShareResource, any, any>()

const { $gettext, $pgettext } = useGettext()

const areHiddenFilesShown = ref(false)
const filterTerm = ref('')
const markInstance = ref<Mark>()
const visibilityFilterAnnouncement = ref('')

const shareSectionTitle = computed(() => {
  return unref(areHiddenFilesShown) ? $gettext('Hidden Shares') : $gettext('Shares')
})

const visibilityOptions = computed(() => [
  { name: 'visible', label: $gettext('Shares') },
  { name: 'hidden', label: $gettext('Hidden Shares') }
])

const setAreHiddenFilesShown = (value: InlineFilterOption) => {
  areHiddenFilesShown.value = value.name === 'hidden'
  resourcesStore.resetSelection()
  visibilityFilterAnnouncement.value = areHiddenFilesShown.value
    ? $pgettext(
        'Accessibility announcement for screen readers when the user switches to show hidden shares only',
        'Showing hidden shares now'
      )
    : $pgettext(
        'Accessibility announcement for screen readers when the user switches to show visible shares only',
        'Visible shares are now showing'
      )
}

const visibleShares = computed(() => unref(paginatedResources).filter((r) => !r.hidden))
const hiddenShares = computed(() => unref(paginatedResources).filter((r) => r.hidden))
const currentItems = computed(() => {
  return unref(areHiddenFilesShown) ? unref(hiddenShares) : unref(visibleShares)
})

const selectedShareTypesQuery = useRouteQuery('q_shareType')
const selectedSharedByQuery = useRouteQuery('q_sharedBy')
const scrollToTarget = useRouteQuery('scrollTo')

const filteredItems = computed(() => {
  let result = unref(currentItems)

  const selectedShareTypes = queryItemAsString(unref(selectedShareTypesQuery))?.split('+')
  if (selectedShareTypes?.length) {
    result = result.filter(({ shareTypes }) => {
      return ShareTypes.getByKeys(selectedShareTypes)
        .map(({ value }) => value)
        .some((t) => shareTypes.includes(t))
    })
  }

  const selectedSharedBy = queryItemAsString(unref(selectedSharedByQuery))?.split('+')
  if (selectedSharedBy?.length) {
    result = result.filter(({ sharedBy }) =>
      sharedBy.some(({ id }) => selectedSharedBy.includes(id))
    )
  }

  if (unref(filterTerm).trim()) {
    const usersSearchEngine = new Fuse(result, { ...defaultFuseOptions, keys: ['name'] })
    const fuseResult = usersSearchEngine.search(unref(filterTerm)).map((r) => r.item)
    result = fuseResult.filter((item) => result.includes(item))
  }

  return result
})

watch(filteredItems, () => {
  if (!unref(areResourcesLoading)) {
    if (!unref(markInstance)) {
      markInstance.value = new Mark('.oc-resource-details')
    }

    unref(markInstance).unmark()
    unref(markInstance).mark(unref(filterTerm), {
      element: 'span',
      className: 'mark-highlight'
    })
  }
})

const { sortBy, sortDir, items, handleSort } = useSort({
  items: filteredItems,
  fields: sortFields
})

const { groupingSettings } = useGroupingSettings({ sortBy, sortDir })

const { getMatchingSpace } = useGetMatchingSpace()

const selectedShareSpace = computed(() => {
  if (unref(selectedResources).length !== 1) {
    return null
  }
  const resource = unref(selectedResources)[0]
  return getMatchingSpace(resource)
})

const openWithDefaultAppQuery = useRouteQuery('openWithDefaultApp')
const performLoaderTask = async () => {
  await loadResourcesTask.perform()
  scrollToResourceFromRoute(unref(items), 'files-app-bar')
  if (queryItemAsString(unref(openWithDefaultAppQuery)) === 'true') {
    openWithDefaultApp({
      space: unref(selectedShareSpace),
      resource: unref(selectedResources)[0]
    })
  }
}

const shareTypes = computed(() => {
  const uniqueShareTypes = uniq(unref(paginatedResources).flatMap((i) => i.shareTypes))

  const ocmAvailable = appsStore.appIds.includes('open-cloud-mesh')
  if (ocmAvailable && !uniqueShareTypes.includes(ShareTypes.remote.value)) {
    uniqueShareTypes.push(ShareTypes.remote.value)
  }

  return ShareTypes.getByValues(uniqueShareTypes).map((shareType) => {
    return {
      key: shareType.key,
      value: shareType.value,
      label: $gettext(shareType.label)
    }
  })
})

const fileOwners = computed(() => {
  const flatList = unref(paginatedResources)
    .map((i) => i.sharedBy)
    .flat()
  return [...new Map(flatList.map((item) => [item.displayName, item])).values()]
})

onMounted(() => {
  performLoaderTask()
})

watch(scrollToTarget, (value) => {
  if (!value) {
    return
  }

  scrollToResourceFromRoute(unref(items), 'files-app-bar')
})
</script>

<style lang="scss" scoped>
.search-filter {
  width: 16rem;
}
</style>
