<template>
  <div>
    <app-template
      ref="template"
      :loading="loadResourcesTask.isRunning || !loadResourcesTask.last"
      :breadcrumbs="breadcrumbs"
      :side-bar-active-panel="sideBarActivePanel"
      :side-bar-available-panels="sideBarAvailablePanels"
      :side-bar-panel-context="sideBarPanelContext"
      :is-side-bar-open="isSideBarOpen"
      :show-batch-actions="!!selectedSpaces.length"
      :batch-actions="batchActions"
      :batch-action-items="selectedSpaces"
      :show-view-options="true"
    >
      <template #topbarActions="{ limitedScreenSpace }">
        <create-space
          v-if="hasCreatePermission"
          :show-label="!limitedScreenSpace"
          @space-created="(space) => spaceSettingsStore.upsertSpace(space)"
        />
      </template>
      <template #sideBarHeader>
        <space-info
          v-if="selectedSpaces.length === 1"
          :space-resource="selectedSpaces[0]"
          class="sidebar-panel__space_info"
        />
      </template>
      <template #mainContent>
        <app-loading-spinner v-if="isLoading" />
        <template v-else>
          <no-content-message
            v-if="!spaces.length"
            id="admin-settings-spaces-empty"
            class="spaces-empty"
            icon="layout-grid"
          >
            <template #message>
              <span v-translate>No spaces in here</span>
            </template>
          </no-content-message>
          <div v-else>
            <spaces-list :class="{ 'spaces-table-squashed': isSideBarOpen }">
              <template #contextMenu>
                <context-actions :items="selectedSpaces" />
              </template>
            </spaces-list>
          </div>
        </template>
      </template>
    </app-template>
  </div>
</template>

<script lang="ts" setup>
import AppTemplate from '../components/AppTemplate.vue'
import SpacesList from '../components/Spaces/SpacesList.vue'
import ContextActions from '../components/Spaces/ContextActions.vue'
import MembersPanel from '../components/Spaces/SideBar/MembersPanel.vue'
import ActionsPanel from '../components/Spaces/SideBar/ActionsPanel.vue'
import {
  NoContentMessage,
  SideBarPanel,
  SideBarPanelContext,
  SpaceAction,
  SpaceDetails,
  SpaceDetailsMultiple,
  SpaceInfo,
  SpaceNoSelection,
  eventBus,
  queryItemAsString,
  useClientService,
  useRouteQuery,
  useSideBar,
  useSpaceActionsDelete,
  useSpaceActionsDisable,
  useSpaceActionsRestore,
  useSpaceActionsEditQuota,
  AppLoadingSpinner,
  useSharesStore,
  useAbility,
  CreateSpace
} from '@ownclouders/web-pkg'
import { call, SpaceResource } from '@ownclouders/web-client'
import { computed, provide, onBeforeUnmount, onMounted, ref, unref } from 'vue'
import { useTask } from 'vue-concurrency'
import { useGettext } from 'vue3-gettext'

import { useSpaceSettingsStore } from '../composables'
import { storeToRefs } from 'pinia'
import { Quota } from '@ownclouders/web-client/graph/generated'

const clientService = useClientService()
const { $gettext } = useGettext()
const { isSideBarOpen, sideBarActivePanel } = useSideBar()
const sharesStore = useSharesStore()
const { can } = useAbility()

const loadResourcesEventToken = ref(null)
let updateQuotaForSpaceEventToken: string
const template = ref(null)
const spaceSettingsStore = useSpaceSettingsStore()
const { spaces, selectedSpaces } = storeToRefs(spaceSettingsStore)

provide(
  'resource',
  computed(() => selectedSpaces.value[0])
)

const currentPageQuery = useRouteQuery('page', '1')
const currentPage = computed(() => {
  return parseInt(queryItemAsString(unref(currentPageQuery)))
})

const itemsPerPageQuery = useRouteQuery('items-per-page', '1')
const itemsPerPage = computed(() => {
  return parseInt(queryItemAsString(unref(itemsPerPageQuery)))
})

const hasCreatePermission = computed(() => can('create-all', 'Drive'))

const loadResourcesTask = useTask(function* (signal) {
  const drives = yield* call(
    clientService.graphAuthenticated.drives.listAllDrives(
      sharesStore.graphRoles,
      {
        orderBy: 'name asc',
        filter: 'driveType eq project'
      },
      { signal }
    )
  )
  spaceSettingsStore.setSpaces(drives)
})

const isLoading = computed(() => {
  return loadResourcesTask.isRunning || !loadResourcesTask.last
})

const breadcrumbs = computed(() => [
  { text: $gettext('Administration Settings'), to: { path: '/admin-settings' } },
  {
    text: $gettext('Spaces'),
    onClick: () => {
      spaceSettingsStore.setSelectedSpaces([])
      loadResourcesTask.perform()
    }
  }
])

const { actions: deleteActions } = useSpaceActionsDelete()
const { actions: disableActions } = useSpaceActionsDisable()
const { actions: editQuotaActions } = useSpaceActionsEditQuota()
const { actions: restoreActions } = useSpaceActionsRestore()

const batchActions = computed((): SpaceAction[] => {
  return [
    ...unref(editQuotaActions),
    ...unref(restoreActions),
    ...unref(deleteActions),
    ...unref(disableActions)
  ].filter((item) => item.isVisible({ resources: unref(selectedSpaces) }))
})

const sideBarPanelContext = computed<SideBarPanelContext<unknown, unknown, SpaceResource>>(() => {
  return {
    parent: null,
    items: unref(selectedSpaces)
  }
})
const sideBarAvailablePanels = [
  {
    name: 'SpaceNoSelection',
    icon: 'layout-grid',
    title: () => $gettext('Details'),
    component: SpaceNoSelection,
    isRoot: () => true,
    isVisible: ({ items }) => items.length === 0
  },
  {
    name: 'SpaceDetails',
    icon: 'layout-grid',
    title: () => $gettext('Details'),
    component: SpaceDetails,
    componentAttrs: () => ({
      showSpaceImage: false,
      showShareIndicators: false
    }),
    isRoot: () => true,
    isVisible: ({ items }) => items.length === 1
  },
  {
    name: 'SpaceDetailsMultiple',
    icon: 'layout-grid',
    title: () => $gettext('Details'),
    component: SpaceDetailsMultiple,
    componentAttrs: ({ items }) => ({
      selectedSpaces: items
    }),
    isRoot: () => true,
    isVisible: ({ items }) => items.length > 1
  },
  {
    name: 'SpaceActions',
    icon: 'play-circle',
    iconFillType: 'line',
    title: () => $gettext('Actions'),
    component: ActionsPanel,
    isVisible: ({ items }) => items.length === 1
  },
  {
    name: 'SpaceMembers',
    icon: 'group',
    title: () => $gettext('Members'),
    component: MembersPanel,
    isVisible: ({ items }) => items.length === 1
  }
] satisfies SideBarPanel<unknown, unknown, SpaceResource>[]

onMounted(async () => {
  await loadResourcesTask.perform()

  loadResourcesEventToken.value = eventBus.subscribe('app.admin-settings.list.load', async () => {
    await loadResourcesTask.perform()
    selectedSpaces.value = []

    const pageCount = Math.ceil(unref(spaces).length / unref(itemsPerPage))
    if (unref(currentPage) > 1 && unref(currentPage) > pageCount) {
      // reset pagination to avoid empty lists (happens when deleting all items on the last page)
      currentPageQuery.value = pageCount.toString()
    }
  })

  updateQuotaForSpaceEventToken = eventBus.subscribe(
    'app.admin-settings.spaces.space.quota.updated',
    ({ spaceId, quota }: { spaceId: string; quota: Quota }) => {
      const space = unref(spaces).find((s) => s.id === spaceId)
      if (space) {
        space.spaceQuota = quota
      }
    }
  )
})

onBeforeUnmount(() => {
  spaceSettingsStore.reset()
  eventBus.unsubscribe('app.admin-settings.list.load', unref(loadResourcesEventToken))
  eventBus.unsubscribe(
    'app.admin-settings.spaces.space.quota.updated',
    updateQuotaForSpaceEventToken
  )
})
</script>
