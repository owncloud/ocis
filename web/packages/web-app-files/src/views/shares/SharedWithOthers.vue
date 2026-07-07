<template>
  <div class="oc-flex">
    <files-view-wrapper>
      <app-bar :is-side-bar-open="isSideBarOpen">
        <template #navigation>
          <SharesNavigation />
        </template>
      </app-bar>
      <app-loading-spinner v-if="areResourcesLoading" />
      <template v-else>
        <fieldset v-if="shareTypes.length > 1" class="oc-flex oc-m-m">
          <div>
            <legend class="oc-mr-m oc-flex oc-flex-middle">
              <oc-icon name="filter-2" class="oc-mr-xs" />
              <span v-text="$gettext('Filter:')" />
            </legend>
          </div>
          <item-filter
            :allow-multiple="true"
            :filter-label="$gettext('Share Type')"
            :filterable-attributes="['label']"
            :items="shareTypes"
            :option-filter-label="$gettext('Filter share types')"
            :show-option-filter="true"
            id-attribute="key"
            class="share-type-filter oc-mx-s"
            display-name-attribute="label"
            filter-name="shareType"
          >
            <template #item="{ item }">
              <span class="oc-ml-s" v-text="item.label" />
            </template>
          </item-filter>
        </fieldset>
        <no-content-message
          v-if="isEmpty"
          id="files-shared-with-others-empty"
          class="files-empty"
          icon="reply"
        >
          <template #message>
            <span v-translate> You have not shared any resources with other people. </span>
          </template>
        </no-content-message>
        <resource-table
          v-else
          v-model:selected-ids="selectedResourcesIds"
          :is-side-bar-open="isSideBarOpen"
          :fields-displayed="['name', 'sharedWith', 'sdate']"
          :are-paths-displayed="true"
          :resources="filteredItems"
          :header-position="fileListHeaderY"
          :sort-by="sortBy"
          :sort-dir="sortDir"
          :grouping-settings="groupingSettings"
          @file-click="triggerDefaultAction"
          @item-visible="loadPreview({ space: getMatchingSpace($event), resource: $event })"
          @sort="handleSort"
        >
          <template #contextMenu="{ resource }">
            <context-actions
              v-if="isResourceInSelection(resource)"
              :action-options="{ space: getMatchingSpace(resource), resources: selectedResources }"
            />
          </template>
          <template #footer>
            <pagination :pages="paginationPages" :current-page="paginationPage" />
            <list-info v-if="filteredItems.length > 0" class="oc-width-1-1 oc-my-s" />
          </template>
        </resource-table>
      </template>
    </files-view-wrapper>
    <file-side-bar
      :is-open="isSideBarOpen"
      :active-panel="sideBarActivePanel"
      :space="selectedResourceSpace"
    />
  </div>
</template>

<script lang="ts" setup>
import {
  queryItemAsString,
  useAppsStore,
  useFileActions,
  useLoadPreview,
  useResourcesStore,
  useRouteQuery
} from '@ownclouders/web-pkg'
import { ItemFilter } from '@ownclouders/web-pkg'
import { uniq } from 'lodash-es'

import { FileSideBar, ResourceTable } from '@ownclouders/web-pkg'
import { AppLoadingSpinner } from '@ownclouders/web-pkg'
import { NoContentMessage } from '@ownclouders/web-pkg'
import { AppBar } from '@ownclouders/web-pkg'
import ListInfo from '../../components/FilesList/ListInfo.vue'
import { Pagination } from '@ownclouders/web-pkg'
import { ContextActions } from '@ownclouders/web-pkg'
import FilesViewWrapper from '../../components/FilesViewWrapper.vue'

import { useResourcesViewDefaults } from '../../composables'
import { computed, unref } from 'vue'
import { useGroupingSettings } from '@ownclouders/web-pkg'
import { useGetMatchingSpace } from '@ownclouders/web-pkg'
import SharesNavigation from '../../components/AppBar/SharesNavigation.vue'
import { OutgoingShareResource, ShareTypes } from '@ownclouders/web-client'
import { useGettext } from 'vue3-gettext'

const { getMatchingSpace } = useGetMatchingSpace()
const appsStore = useAppsStore()
const { $gettext } = useGettext()

const resourcesStore = useResourcesStore()

const { triggerDefaultAction } = useFileActions()
const resourcesViewDefaults = useResourcesViewDefaults<OutgoingShareResource, any, any[]>()
const {
  sortBy,
  sortDir,
  loadResourcesTask,
  selectedResourcesIds,
  paginatedResources,
  viewMode,
  scrollToResourceFromRoute,
  isSideBarOpen,
  areResourcesLoading,
  fileListHeaderY,
  handleSort,
  isResourceInSelection,
  paginationPages,
  paginationPage,
  sideBarActivePanel,
  selectedResourceSpace,
  selectedResources
} = resourcesViewDefaults
const { loadPreview } = useLoadPreview(viewMode)

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
const selectedShareTypesQuery = useRouteQuery('q_shareType')
const filteredItems = computed(() => {
  const selectedShareTypes = queryItemAsString(unref(selectedShareTypesQuery))?.split('+')
  if (!selectedShareTypes || selectedShareTypes.length === 0) {
    return unref(paginatedResources)
  }
  return unref(paginatedResources).filter((item) => {
    return ShareTypes.getByKeys(selectedShareTypes)
      .map(({ value }) => value)
      .some((t) => item.shareTypes.includes(t))
  })
})

resourcesStore.$onAction((action) => {
  if (action.name !== 'updateResourceField') {
    return
  }

  if (selectedResourcesIds.value.length !== 1) return
  const id = selectedResourcesIds.value[0]

  const match = unref(paginatedResources).find((r) => {
    return r.id === id
  })
  if (!match) return

  loadResourcesTask.perform()

  const matchedNewResource = unref(paginatedResources).find((r) => r.fileId === match.fileId)
  if (!matchedNewResource) return

  selectedResourcesIds.value = [matchedNewResource.id]
})
const isEmpty = computed(() => unref(filteredItems).length < 1)
const { groupingSettings } = useGroupingSettings({ sortBy, sortDir })

async function created() {
  await unref(loadResourcesTask).perform()
  scrollToResourceFromRoute(unref(filteredItems), 'files-app-bar')
}

created()
</script>
