<template>
  <div id="user-list">
    <div class="user-filters oc-flex oc-flex-between oc-flex-wrap oc-flex-bottom oc-mx-m oc-mb-m">
      <slot name="filter" />
    </div>
    <app-loading-spinner v-if="isLoading" />
    <div v-else>
      <slot v-if="!users.length" name="noResults" />
      <oc-table
        v-else
        class="users-table"
        :sort-by="sortBy"
        :sort-dir="sortDir"
        :fields="fields"
        :data="paginatedItems"
        :highlighted="highlighted"
        :sticky="isSticky"
        :header-position="fileListHeaderY"
        :hover="true"
        @sort="handleSort"
        @contextmenu-clicked="showContextMenuOnRightClick"
        @highlight="rowClicked"
      >
        <template #selectHeader>
          <oc-checkbox
            size="large"
            class="oc-ml-s"
            :label="$gettext('Select all users')"
            :model-value="allUsersSelected"
            :label-hidden="false"
            :label-classes="['oc-invisible-sr']"
            @update:model-value="
              allUsersSelected ? unselectAllUsers() : selectUsers(paginatedItems)
            "
          />
        </template>
        <template #select="{ item }">
          <oc-checkbox
            class="oc-ml-s"
            size="large"
            :model-value="isUserSelected(item)"
            :option="item"
            :label="getSelectUserLabel(item)"
            :label-hidden="true"
            @update:model-value="selectUser(item)"
            @click.stop="rowClicked([item, $event])"
          />
        </template>
        <template #avatar="{ item }">
          <avatar-image :width="32" :userid="item.id" :user-name="item.displayName" />
        </template>
        <template #role="{ item }">
          <template v-if="item.appRoleAssignments">{{ getRoleDisplayNameByUser(item) }}</template>
        </template>
        <template #accountEnabled="{ item }">
          <span v-if="item.accountEnabled === false" class="oc-flex oc-flex-middle">
            <oc-icon name="stop-circle" fill-type="line" class="oc-mr-s" /><span
              v-text="$gettext('Forbidden')"
            />
          </span>
          <span v-else class="oc-flex oc-flex-middle">
            <oc-icon name="play-circle" fill-type="line" class="oc-mr-s" /><span
              v-text="$gettext('Allowed')"
            />
          </span>
        </template>
        <template #actions="{ item }">
          <oc-button
            v-oc-tooltip="$gettext('Show details')"
            :aria-label="$gettext('Show details')"
            appearance="raw"
            class="oc-mr-xs quick-action-button oc-p-xs users-table-btn-details"
            @click="showDetails(item)"
          >
            <oc-icon name="information" fill-type="line" />
          </oc-button>
          <oc-button
            v-oc-tooltip="$gettext('Edit')"
            :aria-label="$gettext('Edit')"
            appearance="raw"
            class="oc-mr-xs quick-action-button oc-p-xs users-table-btn-edit"
            @click="showEditPanel(item)"
          >
            <oc-icon name="pencil" fill-type="line" />
          </oc-button>
          <context-menu-quick-action
            ref="contextMenuButtonRef"
            :item="item"
            class="users-table-btn-action-dropdown"
            @quick-action-clicked="showContextMenuOnBtnClick($event, item)"
          >
            <template #contextMenu>
              <slot name="contextMenu" :user="item" />
            </template>
          </context-menu-quick-action>
        </template>
        <template #footer>
          <pagination :pages="totalPages" :current-page="currentPage" />
          <div class="oc-text-center oc-width-1-1 oc-my-s">
            <p class="oc-text-muted">{{ footerTextTotal }}</p>
          </div>
        </template>
      </oc-table>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useGettext } from 'vue3-gettext'
import { ComponentPublicInstance, computed, nextTick, onMounted, ref, unref, watch } from 'vue'
import {
  AppLoadingSpinner,
  ContextMenuBtnClickEventData,
  displayPositionedDropdown,
  eventBus,
  queryItemAsString,
  SortDir,
  useIsTopBarSticky,
  useKeyboardActions,
  useRouteQuery
} from '@ownclouders/web-pkg'
import { SideBarEventTopics } from '@ownclouders/web-pkg'
import { AppRole, User } from '@ownclouders/web-client/graph/generated'
import { ContextMenuQuickAction } from '@ownclouders/web-pkg'
import { useFileListHeaderPosition, usePagination } from '@ownclouders/web-pkg'
import { Pagination } from '@ownclouders/web-pkg'
import { perPageDefault, perPageStoragePrefix } from '../../defaults'
import { storeToRefs } from 'pinia'
import { useUserSettingsStore } from '../../composables/stores/userSettings'
import {
  useKeyboardTableMouseActions,
  useKeyboardTableNavigation
} from '../../composables/keyboardActions'
import { findIndex } from 'lodash-es'
import Mark from 'mark.js'

interface Props {
  roles: AppRole[]
  isLoading?: boolean
}
const { roles, isLoading = false } = defineProps<Props>()

const { $gettext } = useGettext()
const { isSticky } = useIsTopBarSticky()

const contextMenuButtonRef = ref(undefined)
const sortBy = ref('onPremisesSamAccountName')
const sortDir = ref<SortDir>(SortDir.Asc)
const { y: fileListHeaderY } = useFileListHeaderPosition('#admin-settings-app-bar')

const lastSelectedUserIndex = ref(0)
const lastSelectedUserId = ref(null)

const userSettingsStore = useUserSettingsStore()
const { users, selectedUsers } = storeToRefs(userSettingsStore)

const isUserSelected = (user: User) => {
  return unref(selectedUsers).some((s) => s.id === user.id)
}
const selectUser = (selectedUser: User) => {
  lastSelectedUserIndex.value = findIndex(unref(users), (u) => u.id === selectedUser.id)
  lastSelectedUserId.value = selectedUser.id
  keyActions.resetSelectionCursor()

  const isUserSelected = unref(selectedUsers).find((user) => user.id === selectedUser.id)
  if (!isUserSelected) {
    return userSettingsStore.addSelectedUser(selectedUser)
  }

  userSettingsStore.setSelectedUsers(
    unref(selectedUsers).filter((user) => user.id !== selectedUser.id)
  )
}

const unselectAllUsers = () => {
  userSettingsStore.setSelectedUsers([])
}

const selectUsers = (users: User[]) => {
  userSettingsStore.setSelectedUsers(users)
}

const showDetails = (user: User) => {
  if (!isUserSelected(user)) {
    selectUser(user)
  }
  eventBus.publish(SideBarEventTopics.open)
}

const showEditPanel = (user: User) => {
  if (!isUserSelected(user)) {
    selectUser(user)
  }
  eventBus.publish(SideBarEventTopics.openWithPanel, 'EditPanel')
}

const rowClicked = (data: [User, MouseEvent]) => {
  const resource = data[0]
  const eventData = data[1]
  const isCheckboxClicked = (eventData?.target as HTMLElement).getAttribute('type') === 'checkbox'

  const contextActionClicked =
    (eventData?.target as HTMLElement)?.closest('div')?.id === 'oc-files-context-menu'
  if (contextActionClicked) {
    return
  }

  if (eventData?.metaKey) {
    return eventBus.publish('app.resources.list.clicked.meta', resource)
  }
  if (eventData?.shiftKey) {
    return eventBus.publish('app.resources.list.clicked.shift', {
      resource,
      skipTargetSelection: isCheckboxClicked
    })
  }
  if (isCheckboxClicked) {
    return
  }
  unselectAllUsers()
  selectUser(resource)
}
const showContextMenuOnBtnClick = (data: ContextMenuBtnClickEventData, user: User) => {
  const { dropdown, event } = data
  if (dropdown?.tippy === undefined) {
    return
  }
  if (!isUserSelected(user)) {
    userSettingsStore.setSelectedUsers([user])
  }
  displayPositionedDropdown(dropdown.tippy, event, unref(contextMenuButtonRef))
}
const showContextMenuOnRightClick = (
  row: ComponentPublicInstance<unknown>,
  event: MouseEvent,
  user: User
) => {
  event.preventDefault()
  const dropdown = row.$el.getElementsByClassName('users-table-btn-action-dropdown')[0]
  if (dropdown === undefined) {
    return
  }
  if (!isUserSelected(user)) {
    userSettingsStore.setSelectedUsers([user])
  }
  displayPositionedDropdown(dropdown._tippy, event, unref(contextMenuButtonRef))
}

const getRoleDisplayNameByUser = (user: User) => {
  const assignedRole = user.appRoleAssignments[0]

  return (
    $gettext(roles.find((role) => role.id === assignedRole?.appRoleId)?.displayName || '') || '-'
  )
}

const orderBy = (list: User[], prop: string, desc: boolean) => {
  return [...list].sort((user1, user2) => {
    let a: string, b: string

    switch (prop) {
      case 'role':
        a = getRoleDisplayNameByUser(user1)
        b = getRoleDisplayNameByUser(user2)
        break
      case 'accountEnabled':
        a = ('accountEnabled' in user1 ? user1.accountEnabled : true).toString()
        b = ('accountEnabled' in user2 ? user2.accountEnabled : true).toString()
        break
      default:
        a = user1[prop as keyof User].toString() || ''
        b = user2[prop as keyof User].toString() || ''
    }

    return desc ? b.localeCompare(a) : a.localeCompare(b)
  })
}

const items = computed(() => {
  return orderBy(unref(users), unref(sortBy), unref(sortDir) === SortDir.Desc)
})

const {
  items: paginatedItems,
  page: currentPage,
  total: totalPages
} = usePagination({ items, perPageDefault, perPageStoragePrefix })

const keyActions = useKeyboardActions()
useKeyboardTableNavigation(
  keyActions,
  paginatedItems,
  selectedUsers,
  lastSelectedUserIndex,
  lastSelectedUserId
)
useKeyboardTableMouseActions(
  keyActions,
  paginatedItems,
  selectedUsers,
  lastSelectedUserIndex,
  lastSelectedUserId
)

const markInstance = ref<Mark>(null)
onMounted(async () => {
  await nextTick()
  markInstance.value = new Mark('.mark-element')
})
const displayNameQuery = useRouteQuery('q_displayName')
watch([displayNameQuery, paginatedItems], () => {
  unref(markInstance)?.unmark()
  const filterTerm = queryItemAsString(unref(displayNameQuery))
  if (filterTerm) {
    unref(markInstance)?.mark(filterTerm, {
      element: 'span',
      className: 'mark-highlight'
    })
  }
})
function handleSort(event: { sortBy: string; sortDir: SortDir }) {
  sortBy.value = event.sortBy
  sortDir.value = event.sortDir
}
function getSelectUserLabel(user: User) {
  return $gettext('Select %{ user }', { user: user.displayName }, true)
}
const allUsersSelected = computed(() => {
  return unref(paginatedItems).length === unref(selectedUsers).length
})
const footerTextTotal = computed(() => {
  return $gettext('%{userCount} users in total', {
    userCount: unref(users).length.toString()
  })
})
const fields = computed(() => {
  return [
    {
      name: 'select',
      title: '',
      type: 'slot',
      width: 'shrink',
      headerType: 'slot'
    },
    {
      name: 'avatar',
      title: '',
      type: 'slot',
      width: 'shrink'
    },
    {
      name: 'onPremisesSamAccountName',
      title: $gettext('User name'),
      sortable: true
    },
    {
      name: 'displayName',
      title: $gettext('First and last name'),
      sortable: true,
      tdClass: 'mark-element'
    },
    {
      name: 'mail',
      title: $gettext('Email'),
      sortable: true
    },
    {
      name: 'role',
      title: $gettext('Role'),
      type: 'slot',
      sortable: true
    },
    {
      name: 'accountEnabled',
      title: $gettext('Login'),
      type: 'slot',
      sortable: true
    },
    {
      name: 'actions',
      title: $gettext('Actions'),
      sortable: false,
      type: 'slot',
      alignH: 'right'
    }
  ]
})
const highlighted = computed(() => {
  return unref(selectedUsers).map((user: User) => user.id)
})
</script>

<style lang="scss">
.users-table {
  .oc-table-header-cell-actions,
  .oc-table-data-cell-actions {
    white-space: nowrap;
  }

  .oc-table-header-cell-role,
  .oc-table-data-cell-role,
  .oc-table-header-cell-accountEnabled,
  .oc-table-data-cell-accountEnabled,
  .oc-table-header-cell-mail,
  .oc-table-data-cell-mail {
    display: none;

    @media only screen and (min-width: 1200px) {
      display: table-cell;
    }
  }

  .oc-table-header-cell-displayName,
  .oc-table-data-cell-displayName {
    display: none;

    @media only screen and (min-width: 1000px) {
      display: table-cell;
    }
  }

  &-squashed {
    .oc-table-header-cell-role,
    .oc-table-data-cell-role,
    .oc-table-header-cell-accountEnabled,
    .oc-table-data-cell-accountEnabled,
    .oc-table-header-cell-mail,
    .oc-table-data-cell-mail {
      display: none;

      @media only screen and (min-width: 1600px) {
        display: table-cell;
      }
    }

    .oc-table-header-cell-displayName,
    .oc-table-data-cell-displayName {
      display: none;

      @media only screen and (min-width: 1400px) {
        display: table-cell;
      }
    }

    .oc-table-header-cell-mail,
    .oc-table-data-cell-mail {
      display: none;

      @media only screen and (min-width: 1200px) {
        display: table-cell;
      }
    }
  }
}
</style>
