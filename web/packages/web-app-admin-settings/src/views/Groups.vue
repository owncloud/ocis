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
      :show-batch-actions="!!selectedGroups.length"
      :batch-actions="batchActions"
      :batch-action-items="selectedGroups"
      :show-view-options="true"
    >
      <template #topbarActions="{ limitedScreenSpace }">
        <div>
          <oc-button
            id="create-group-btn"
            v-oc-tooltip="limitedScreenSpace ? createGroupAction.label() : undefined"
            class="oc-mr-s"
            variation="primary"
            appearance="filled"
            @click="createGroupAction.handler()"
          >
            <oc-icon :name="createGroupAction.icon" />
            <span v-if="!limitedScreenSpace" v-text="createGroupAction.label()" />
          </oc-button>
        </div>
      </template>
      <template #mainContent>
        <app-loading-spinner v-if="isLoading" />
        <div v-else>
          <groups-list>
            <template #contextMenu>
              <context-actions :action-options="{ resources: selectedGroups }" />
            </template>
            <template #filter>
              <div class="oc-flex oc-flex-middle">
                <oc-text-input
                  id="groups-filter"
                  v-model.trim="filterTerm"
                  :label="$gettext('Search')"
                  autocomplete="off"
                  @enter-key-down="filterGroups"
                />
                <oc-button
                  id="groups-filter-confirm"
                  class="oc-ml-xs"
                  appearance="raw"
                  @click="filterGroups"
                >
                  <oc-icon name="search" fill-type="line" />
                </oc-button>
              </div>
            </template>
            <template #noResults>
              <no-content-message
                v-if="!groups.length"
                id="admin-settings-groups-empty"
                class="files-empty"
                icon="user"
              >
                <template #message>
                  {{
                    $pgettext(
                      'A message displayed when no groups are found in the groups list in admin settings when there is no filter applied.',
                      'No groups in here'
                    )
                  }}
                </template>
              </no-content-message>
            </template>
          </groups-list>
        </div>
      </template>
    </app-template>
  </div>
</template>

<script lang="ts" setup>
import AppTemplate from '../components/AppTemplate.vue'
import ContextActions from '../components/Groups/ContextActions.vue'
import DetailsPanel from '../components/Groups/SideBar/DetailsPanel.vue'
import EditPanel from '../components/Groups/SideBar/EditPanel.vue'
import GroupsList from '../components/Groups/GroupsList.vue'
import MembersPanel from '../components/Groups/SideBar/MembersPanel.vue'
import { useGroupSettingsStore } from '../composables'
import { useGroupActionsCreateGroup, useGroupActionsDelete } from '../composables/actions/groups'
import {
  AppLoadingSpinner,
  NoContentMessage,
  queryItemAsString,
  SideBarPanel,
  SideBarPanelContext,
  useClientService,
  useSideBar
} from '@ownclouders/web-pkg'
import { Group } from '@ownclouders/web-client/graph/generated'
import { computed, provide, ref, unref, onBeforeUnmount, onMounted } from 'vue'
import { useTask } from 'vue-concurrency'
import { useGettext } from 'vue3-gettext'
import { storeToRefs } from 'pinia'
import { call } from '@ownclouders/web-client'
import { useRoute, useRouter } from 'vue-router'
import { omit } from 'lodash-es'

const template = ref()
const groupSettingsStore = useGroupSettingsStore()
const { selectedGroups, groups } = storeToRefs(groupSettingsStore)
const clientService = useClientService()
const { $gettext } = useGettext()
const { sideBarActivePanel, isSideBarOpen } = useSideBar()
const router = useRouter()
const route = useRoute()

provide(
  'group',
  computed(() => unref(selectedGroups)[0])
)

const filterTerm = ref(queryItemAsString(unref(route).query.q_displayName))

const { actions: createGroupActions } = useGroupActionsCreateGroup()
const createGroupAction = computed(() => unref(createGroupActions)[0])

const loadResourcesTask = useTask(function* (signal) {
  const loadedGroups = yield* call(
    clientService.graphAuthenticated.groups.listGroups(
      {
        orderBy: ['displayName'],
        expand: ['members'],
        search: queryItemAsString(unref(route).query.q_displayName)
      },
      { signal }
    )
  )
  groupSettingsStore.setGroups(loadedGroups || [])
}).restartable()

const { actions: deleteActions } = useGroupActionsDelete()
const batchActions = computed(() => {
  return [...unref(deleteActions)].filter((item) =>
    item.isVisible({ resources: unref(selectedGroups) })
  )
})

const isLoading = computed(() => {
  return loadResourcesTask.isRunning || !loadResourcesTask.last
})

const sideBarPanelContext = computed<SideBarPanelContext<unknown, unknown, Group>>(() => {
  return {
    parent: null,
    items: unref(selectedGroups)
  }
})

const sideBarAvailablePanels = [
  {
    name: 'DetailsPanel',
    icon: 'group-2',
    title: () => $gettext('Details'),
    component: DetailsPanel,
    componentAttrs: () => ({ groups: unref(selectedGroups) }),
    isRoot: () => true,
    isVisible: () => true
  },
  {
    name: 'EditPanel',
    icon: 'pencil',
    title: () => $gettext('Edit group'),
    component: EditPanel,
    componentAttrs: ({ items }) => {
      return {
        group: items.length === 1 ? items[0] : null
      }
    },
    isVisible: ({ items }) => {
      return items.length === 1 && !items[0].groupTypes?.includes('ReadOnly')
    }
  },
  {
    name: 'GroupMembers',
    icon: 'group',
    title: () => $gettext('Members'),
    component: MembersPanel,
    isVisible: ({ items }) => items.length === 1
  }
] satisfies SideBarPanel<unknown, unknown, Group>[]

async function filterGroups() {
  await router.push({
    ...unref(route),
    query: {
      ...omit(unref(route).query, 'q_displayName'),
      ...(unref(filterTerm) ? { q_displayName: unref(filterTerm) } : {}),
      page: '1'
    }
  })
  loadResourcesTask.perform()
  groupSettingsStore.setSelectedGroups([])
}

onMounted(async () => {
  await loadResourcesTask.perform()
})

onBeforeUnmount(() => {
  groupSettingsStore.reset()
})

const breadcrumbs = computed(() => {
  return [
    { text: $gettext('Administration Settings'), to: { path: '/admin-settings' } },
    {
      text: $gettext('Groups'),
      onClick: () => {
        groupSettingsStore.setSelectedGroups([])
        loadResourcesTask.perform()
      }
    }
  ]
})
</script>

<style lang="scss" scoped>
#groups-filter-confirm {
  margin-top: calc(0.2rem + var(--oc-font-size-default));
}
</style>
