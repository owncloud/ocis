<template>
  <div class="oc-flex oc-width-1-1">
    <files-view-wrapper>
      <app-bar
        :breadcrumbs="breadcrumbs"
        :has-bulk-actions="true"
        :is-side-bar-open="isSideBarOpen"
        :space="space"
      />
      <app-loading-spinner v-if="areResourcesLoading" />
      <template v-else>
        <no-content-message
          v-if="isEmpty"
          id="files-trashbin-empty"
          class="files-empty"
          icon="delete-bin-7"
          icon-fill-type="line"
        >
          <template #message>
            <span>{{ emptyTrashMessage }}</span>
          </template>
        </no-content-message>
        <resource-table
          v-else
          v-model:selected-ids="selectedResourcesIds"
          :is-side-bar-open="isSideBarOpen"
          :fields-displayed="['name', 'ddate']"
          :are-paths-displayed="true"
          :resources="paginatedResources"
          :are-resources-clickable="false"
          :header-position="fileListHeaderY"
          :sort-by="sortBy"
          :sort-dir="sortDir"
          :space="space"
          :has-actions="showActions"
          @sort="handleSort"
        >
          <template #contextMenu="{ resource }">
            <context-actions
              v-if="isResourceInSelection(resource)"
              :action-options="{ space, resources: selectedResources }"
            />
          </template>
          <template #footer>
            <pagination :pages="paginationPages" :current-page="paginationPage" />
            <list-info v-if="paginatedResources.length > 0" class="oc-width-1-1 oc-my-s" />
          </template>
        </resource-table>
      </template>
    </files-view-wrapper>
    <file-side-bar :is-open="isSideBarOpen" :active-panel="sideBarActivePanel" :space="space" />
  </div>
</template>

<script lang="ts" setup>
import { storeToRefs } from 'pinia'

import { AppBar, ContextActions, FileSideBar, useUserStore } from '@ownclouders/web-pkg'
import FilesViewWrapper from '../../components/FilesViewWrapper.vue'
import ListInfo from '../../components/FilesList/ListInfo.vue'
import { ResourceTable } from '@ownclouders/web-pkg'
import { AppLoadingSpinner } from '@ownclouders/web-pkg'
import { NoContentMessage } from '@ownclouders/web-pkg'
import { Pagination } from '@ownclouders/web-pkg'

import { eventBus } from '@ownclouders/web-pkg'
import { useResourcesViewDefaults } from '../../composables'
import { computed, onMounted, onBeforeUnmount, unref } from 'vue'
import { Resource } from '@ownclouders/web-client'
import { createLocationTrash } from '@ownclouders/web-pkg'
import { isProjectSpaceResource, SpaceResource } from '@ownclouders/web-client'
import { useDocumentTitle } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'

interface Props {
  space?: SpaceResource
}
const { space = null } = defineProps<Props>()
const { $gettext } = useGettext()
const userStore = useUserStore()
const { user } = storeToRefs(userStore)

let loadResourcesEventToken: string
const emptyTrashMessage = computed(() => {
  return space.driveType === 'personal'
    ? $gettext('You have no deleted files')
    : $gettext('Space has no deleted files')
})

const titleSegments = computed(() => {
  const segments = [$gettext('Deleted files')]
  segments.unshift(space.name)

  return segments
})
useDocumentTitle({ titleSegments })

const {
  loadResourcesTask,
  refreshFileListHeaderPosition,
  scrollToResourceFromRoute,
  paginatedResources,
  isSideBarOpen,
  areResourcesLoading,
  sortBy,
  sortDir,
  selectedResourcesIds,
  fileListHeaderY,
  handleSort,
  isResourceInSelection,
  selectedResources,
  paginationPages,
  paginationPage,
  sideBarActivePanel
} = useResourcesViewDefaults<Resource, any, any[]>()

const performLoaderTask = async () => {
  await loadResourcesTask.perform(space)
  refreshFileListHeaderPosition()
  scrollToResourceFromRoute(unref(paginatedResources), 'files-app-bar')
}
const isEmpty = computed(() => {
  return unref(paginatedResources).length < 1
})

const breadcrumbs = computed(() => {
  let currentNodeName = space?.name
  if (space.driveType === 'personal') {
    currentNodeName = $gettext('Personal')
  }
  return [
    {
      text: $gettext('Deleted files'),
      to: createLocationTrash('files-trash-overview')
    },
    {
      text: currentNodeName,
      onClick: () => eventBus.publish('app.files.list.load')
    }
  ]
})

const showActions = computed(() => {
  return (
    !isProjectSpaceResource(space) ||
    space.canDeleteFromTrashBin({ user: unref(user) }) ||
    space.canRestoreFromTrashbin({ user: unref(user) })
  )
})
onMounted(() => {
  performLoaderTask()
  loadResourcesEventToken = eventBus.subscribe('app.files.list.load', () => {
    performLoaderTask()
  })
})

onBeforeUnmount(() => {
  eventBus.unsubscribe('app.files.list.load', loadResourcesEventToken)
})
</script>
