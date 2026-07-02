<template>
  <div>
    <app-template
      ref="template"
      :breadcrumbs="breadcrumbs"
      :side-bar-active-panel="sideBarActivePanel"
      :side-bar-available-panels="sideBarAvailablePanels"
      :side-bar-panel-context="sideBarPanelContext"
      :is-side-bar-open="isSideBarOpen"
      :side-bar-loading="sideBarLoading"
      :show-batch-actions="!!selectedUsers.length"
      :batch-actions="batchActions"
      :batch-action-items="selectedUsers"
      :show-view-options="true"
    >
      <template #topbarActions="{ limitedScreenSpace }">
        <div>
          <oc-button
            v-if="createUserAction.isVisible()"
            id="create-user-btn"
            v-oc-tooltip="limitedScreenSpace ? createUserAction.label() : undefined"
            class="oc-mr-s"
            variation="primary"
            appearance="filled"
            @click="createUserAction.handler()"
          >
            <oc-icon :name="createUserAction.icon" />
            <span v-if="!limitedScreenSpace" v-text="createUserAction.label()" />
          </oc-button>
        </div>
      </template>
      <template #mainContent>
        <users-list
          :is-loading="isLoading"
          :roles="roles"
          :class="{ 'users-table-squashed': isSideBarOpen }"
        >
          <template #contextMenu>
            <context-actions :items="selectedUsers" />
          </template>
          <template #filter>
            <fieldset class="oc-flex oc-flex-middle">
              <div>
                <legend class="oc-mr-m oc-flex oc-flex-middle">
                  <oc-icon name="filter-2" class="oc-mr-xs" />
                  <span v-text="$gettext('Filter:')" />
                </legend>
              </div>
              <item-filter
                v-if="groups.length"
                :allow-multiple="true"
                :filter-label="$gettext('Groups')"
                :filterable-attributes="['displayName']"
                :items="groups"
                :option-filter-label="$gettext('Filter groups')"
                :show-option-filter="true"
                class="oc-mr-s"
                display-name-attribute="displayName"
                filter-name="groups"
                @selection-change="filterGroups"
              >
                <template #image="{ item }">
                  <avatar-image :width="32" :userid="item.id" :user-name="item.displayName" />
                </template>
                <template #item="{ item }">
                  <div v-text="item.displayName" />
                </template>
              </item-filter>
              <item-filter
                v-if="roles.length"
                :allow-multiple="true"
                :filter-label="$gettext('Roles')"
                :filterable-attributes="['displayName']"
                :items="roles"
                :option-filter-label="$gettext('Filter roles')"
                :show-option-filter="true"
                display-name-attribute="displayName"
                filter-name="roles"
                @selection-change="filterRoles"
              >
                <template #image="{ item }">
                  <avatar-image
                    :width="32"
                    :userid="item.id"
                    :user-name="$gettext(item.displayName)"
                  />
                </template>
                <template #item="{ item }">
                  <div v-text="$gettext(item.displayName)" />
                </template>
              </item-filter>
            </fieldset>
            <div class="oc-flex oc-flex-middle">
              <oc-text-input
                id="users-filter"
                v-model.trim="filterTermDisplayName"
                :label="$gettext('Search')"
                autocomplete="off"
                @keypress.enter="filterDisplayName"
              />
              <oc-button
                id="users-filter-confirm"
                class="oc-ml-xs"
                appearance="raw"
                @click="filterDisplayName"
              >
                <oc-icon name="search" fill-type="line" />
              </oc-button>
            </div>
          </template>
          <template #noResults>
            <no-content-message
              v-if="isFilteringMandatory && !isFilteringActive"
              icon="error-warning"
            >
              <template #message>
                <span v-text="$gettext('Please specify a filter to see results')" />
              </template>
            </no-content-message>
            <no-content-message v-else icon="user">
              <template #message>
                <span v-text="$gettext('No users in here')" />
              </template>
            </no-content-message>
          </template>
        </users-list>
      </template>
    </app-template>
  </div>
</template>

<script lang="ts" setup>
import AppTemplate from '../components/AppTemplate.vue'
import UsersList from '../components/Users/UsersList.vue'
import ContextActions from '../components/Users/ContextActions.vue'
import DetailsPanel from '../components/Users/SideBar/DetailsPanel.vue'
import EditPanel from '../components/Users/SideBar/EditPanel.vue'
import {
  useUserActionsDelete,
  useUserActionsRemoveFromGroups,
  useUserActionsAddToGroups,
  useUserActionsEditLogin,
  useUserActionsEditQuota,
  useUserActionsCreateUser
} from '../composables'
import { User, Group, AppRole, Quota } from '@ownclouders/web-client/graph/generated'
import {
  ItemFilter,
  NoContentMessage,
  eventBus,
  queryItemAsString,
  useClientService,
  useRoute,
  useRouteQuery,
  useRouter,
  useSideBar,
  SideBarPanel,
  SideBarPanelContext,
  useConfigStore,
  QueryValue
} from '@ownclouders/web-pkg'
import { computed, ref, onBeforeUnmount, onMounted, unref, watch, Ref } from 'vue'
import { useTask } from 'vue-concurrency'
import { useGettext } from 'vue3-gettext'
import { format } from 'util'
import { omit } from 'lodash-es'
import { storeToRefs } from 'pinia'

import { useUserSettingsStore } from '../composables/stores/userSettings'
import { call } from '@ownclouders/web-client'

const { $gettext } = useGettext()
const router = useRouter()
const route = useRoute()
const clientService = useClientService()
const configStore = useConfigStore()
const { isSideBarOpen, sideBarActivePanel } = useSideBar()

const userSettingsStore = useUserSettingsStore()
const { users, selectedUsers } = storeToRefs(userSettingsStore)

const writableGroups = computed<Group[]>(() => {
  return unref(groups).filter((g) => !g.groupTypes?.includes('ReadOnly'))
})

const { actions: createUserActions } = useUserActionsCreateUser()
const createUserAction = computed(() => unref(createUserActions)[0])

const { actions: deleteActions } = useUserActionsDelete()
const { actions: removeFromGroupsActions } = useUserActionsRemoveFromGroups({
  groups: writableGroups
})
const { actions: addToGroupsActions } = useUserActionsAddToGroups({
  groups: writableGroups
})
const { actions: editLoginActions } = useUserActionsEditLogin()
const { actions: editQuotaActions } = useUserActionsEditQuota()

const groups = ref([])
const roles = ref([])
const additionalUserDataLoadedForUserIds = ref([])
const applicationId = ref()
const selectedUserIds = computed(() => unref(selectedUsers).map((selectedUser) => selectedUser.id))
const isFilteringMandatory = ref(configStore.options.userListRequiresFilter)

const sideBarLoading = ref(false)
const template = ref()
const displayNameQuery = useRouteQuery('q_displayName')
const filterTermDisplayName = ref(queryItemAsString(unref(displayNameQuery)))

let editQuotaActionEventToken: string

const loadGroupsTask = useTask(function* (signal) {
  groups.value = yield* call(
    clientService.graphAuthenticated.groups.listGroups(
      {
        orderBy: ['displayName'],
        expand: ['members']
      },
      { signal }
    )
  )
}).restartable()

const loadAppRolesTask = useTask(function* (signal) {
  const applications = yield* call(
    clientService.graphAuthenticated.applications.listApplications({ signal })
  )
  roles.value = applications[0].appRoles
  applicationId.value = applications[0].id
})

const loadUsersTask = useTask(function* (signal) {
  if (unref(isFilteringMandatory) && !unref(isFilteringActive)) {
    return userSettingsStore.setUsers([])
  }

  const filter = Object.values(filters)
    .reduce((acc, f) => {
      if ('value' in f) {
        if (unref(f.value)) {
          acc.push(format(f.query, unref(f.value)))
        }
        return acc
      }

      const str = unref(f.ids)
        .map((id) => format(f.query, id))
        .join(' or ')
      if (str) {
        acc.push(`(${str})`)
      }
      return acc
    }, [])
    .filter(Boolean)
    .join(' and ')

  const usersResponse = yield clientService.graphAuthenticated.users.listUsers(
    {
      orderBy: ['displayName'],
      filter,
      expand: ['appRoleAssignments']
    },
    { signal }
  )
  userSettingsStore.setUsers(usersResponse || [])
})

const isLoading = computed(() => {
  return (
    loadUsersTask.isRunning ||
    !loadUsersTask.last ||
    loadResourcesTask.isRunning ||
    !loadResourcesTask.last
  )
})

const loadResourcesTask = useTask(function* () {
  yield loadUsersTask.perform()
  yield loadGroupsTask.perform()
  yield loadAppRolesTask.perform()
})

/**
 * This function reloads the user with expanded attributes,
 * this is necessary as we don't load all the data while listing the users
 * for performance reasons
 */
const loadAdditionalUserDataTask = useTask(function* (signal, user, forceReload = false) {
  /**
   * Prevent load additional user data multiple times if not needed
   */
  if (!forceReload && unref(additionalUserDataLoadedForUserIds).includes(user.id)) {
    return
  }

  const data = yield clientService.graphAuthenticated.users.getUser(user.id, {}, { signal })
  unref(additionalUserDataLoadedForUserIds).push(user.id)

  Object.assign(user, data)
})

const resetPagination = () => {
  return router.push({ ...unref(route), query: { ...unref(route).query, page: '1' } })
}

const filters: Record<
  string,
  { param: Ref<QueryValue>; query: string; ids?: Ref<string[]>; value?: Ref<string> }
> = {
  groups: {
    param: useRouteQuery('q_groups'),
    query: `memberOf/any(m:m/id eq '%s')`,
    ids: ref([])
  },
  roles: {
    param: useRouteQuery('q_roles'),
    query: `appRoleAssignments/any(m:m/appRoleId eq '%s')`,
    ids: ref([])
  },
  displayName: {
    param: useRouteQuery('q_displayName'),
    query: `contains(displayName,'%s')`,
    value: ref('')
  }
}

const isFilteringActive = computed(() => {
  return (
    unref(filters.groups.ids)?.length ||
    unref(filters.roles.ids)?.length ||
    unref(filters.displayName.value)?.length
  )
})
const filterGroups = (groups: Group[]) => {
  filters.groups.ids.value = groups.map((g) => g.id)
  loadUsersTask.perform()
  if (userSettingsStore.selectedUsers.length) {
    // only reset selection if there are selected users because is messes with the focus otherwise
    userSettingsStore.setSelectedUsers([])
  }
  additionalUserDataLoadedForUserIds.value = []
  return resetPagination()
}
const filterRoles = (roles: AppRole[]) => {
  filters.roles.ids.value = roles.map((r) => r.id)
  loadUsersTask.perform()
  if (userSettingsStore.selectedUsers.length) {
    // only reset selection if there are selected users because is messes with the focus otherwise
    userSettingsStore.setSelectedUsers([])
  }
  additionalUserDataLoadedForUserIds.value = []
  return resetPagination()
}
const filterDisplayName = async () => {
  await router.push({
    ...unref(route),
    query: {
      ...omit(unref(route).query, 'q_displayName'),
      ...(unref(filterTermDisplayName) && { q_displayName: unref(filterTermDisplayName) })
    }
  })
  filters.displayName.value.value = unref(filterTermDisplayName)
  loadUsersTask.perform()
  userSettingsStore.setSelectedUsers([])
  additionalUserDataLoadedForUserIds.value = []
  return resetPagination()
}

watch(selectedUserIds, async () => {
  sideBarLoading.value = true
  await Promise.all(unref(selectedUsers).map((user) => loadAdditionalUserDataTask.perform(user)))
  sideBarLoading.value = false
})

const batchActions = computed(() => {
  return [
    ...unref(deleteActions),
    ...unref(editQuotaActions),
    ...unref(addToGroupsActions),
    ...unref(removeFromGroupsActions),
    ...unref(editLoginActions)
  ].filter((item) => item.isVisible({ resources: unref(selectedUsers) }))
})

const updateSpaceQuota = ({ spaceId, quota }: { spaceId: string; quota: Quota }) => {
  const user = unref(users).find((u) => u.drive?.id === spaceId)
  user.drive.quota = quota
  userSettingsStore.upsertUser(user)
}

onMounted(async () => {
  for (const f in filters) {
    if (unref(filters[f]).hasOwnProperty('ids')) {
      filters[f].ids.value = queryItemAsString(unref(filters[f].param))?.split('+') || []
    }
    if (unref(filters[f]).hasOwnProperty('value')) {
      filters[f].value.value = queryItemAsString(unref(filters[f].param))
    }
  }

  await loadResourcesTask.perform()

  editQuotaActionEventToken = eventBus.subscribe(
    'app.admin-settings.users.user.quota.updated',
    updateSpaceQuota
  )
})

onBeforeUnmount(() => {
  userSettingsStore.reset()

  eventBus.unsubscribe('app.admin-settings.users.user.quota.updated', editQuotaActionEventToken)
})

const sideBarPanelContext = computed<SideBarPanelContext<unknown, unknown, User>>(() => {
  return {
    parent: null,
    items: unref(selectedUsers)
  }
})
const sideBarAvailablePanels = [
  {
    name: 'DetailsPanel',
    icon: 'user',
    title: () => $gettext('Details'),
    component: DetailsPanel,
    componentAttrs: ({ items }) => ({
      user: items.length === 1 ? items[0] : null,
      users: items,
      roles: unref(roles)
    }),
    isRoot: () => true,
    isVisible: () => true
  },
  {
    name: 'EditPanel',
    icon: 'pencil',
    title: () => $gettext('Edit user'),
    component: EditPanel,
    isVisible: ({ items }) => items.length === 1,
    componentAttrs: ({ items }) => ({
      user: items.length === 1 ? items[0] : null,
      roles: unref(roles),
      groups: unref(groups),
      applicationId: unref(applicationId)
    })
  }
] satisfies SideBarPanel<unknown, unknown, User>[]

const breadcrumbs = computed(() => {
  return [
    { text: $gettext('Administration Settings'), to: { path: '/admin-settings' } },
    {
      text: $gettext('Users'),
      onClick: () => {
        userSettingsStore.setSelectedUsers([])
        loadResourcesTask.perform()
      }
    }
  ]
})
</script>
<style lang="scss" scoped>
#users-filter {
  width: 16rem;

  &-confirm {
    margin-top: calc(0.2rem + var(--oc-font-size-default));
  }
}
</style>
