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
        <no-content-message
          v-if="isEmpty"
          id="files-shared-via-link-empty"
          class="files-empty"
          icon="link"
        >
          <template #message>
            <span v-translate>You have not shared any resource via link.</span>
          </template>
        </no-content-message>
        <resource-table
          v-else
          v-model:selected-ids="selectedResourcesIds"
          :is-side-bar-open="isSideBarOpen"
          :fields-displayed="['name', 'sdate']"
          :are-paths-displayed="true"
          :resources="paginatedResources"
          :header-position="fileListHeaderY"
          :sort-by="sortBy"
          :sort-dir="sortDir"
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
            <list-info v-if="paginatedResources.length > 0" class="oc-width-1-1 oc-my-s" />
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
  FileSideBar,
  useFileActions,
  useLoadPreview,
  useResourcesStore
} from '@ownclouders/web-pkg'
import { AppLoadingSpinner } from '@ownclouders/web-pkg'
import { NoContentMessage } from '@ownclouders/web-pkg'
import { AppBar } from '@ownclouders/web-pkg'
import ListInfo from '../../components/FilesList/ListInfo.vue'
import { ContextActions } from '@ownclouders/web-pkg'
import FilesViewWrapper from '../../components/FilesViewWrapper.vue'
import { ResourceTable } from '@ownclouders/web-pkg'
import { Pagination } from '@ownclouders/web-pkg'

import { useResourcesViewDefaults } from '../../composables'
import { computed, unref } from 'vue'
import { useGetMatchingSpace } from '@ownclouders/web-pkg'
import SharesNavigation from '../../../src/components/AppBar/SharesNavigation.vue'
import { OutgoingShareResource } from '@ownclouders/web-client'

const { getMatchingSpace } = useGetMatchingSpace()

const resourcesStore = useResourcesStore()
const { triggerDefaultAction } = useFileActions()

const {
  loadResourcesTask,
  selectedResourcesIds,
  paginatedResources,
  viewMode,
  scrollToResourceFromRoute,
  isSideBarOpen,
  sideBarActivePanel,
  selectedResourceSpace,
  sortBy,
  sortDir,
  isResourceInSelection,
  handleSort,
  areResourcesLoading,
  fileListHeaderY,
  paginationPages,
  paginationPage,
  selectedResources
} = useResourcesViewDefaults<OutgoingShareResource, any, any[]>()
const { loadPreview } = useLoadPreview(viewMode)

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

const isEmpty = computed(() => {
  return unref(paginatedResources).length < 1
})
async function created() {
  await unref(loadResourcesTask).perform()
  scrollToResourceFromRoute(unref(paginatedResources), 'files-app-bar')
}

created()
</script>
