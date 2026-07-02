<template>
  <div id="space-list">
    <div class="space-filters oc-flex oc-flex-right oc-flex-wrap oc-flex-bottom oc-mx-m oc-mb-m">
      <oc-text-input
        id="spaces-filter"
        v-model="filterTerm"
        :label="$gettext('Search')"
        autocomplete="off"
      />
    </div>
    <oc-table
      class="spaces-table"
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
      @highlight="fileClicked"
    >
      <template #selectHeader>
        <oc-checkbox
          size="large"
          class="oc-ml-s"
          :label="$gettext('Select all spaces')"
          :model-value="allSpacesSelected"
          :label-hidden="true"
          @update:model-value="
            allSpacesSelected ? unselectAllSpaces() : selectSpaces(paginatedItems)
          "
        />
      </template>
      <template #select="{ item }">
        <oc-checkbox
          class="oc-ml-s"
          size="large"
          :model-value="isSpaceSelected(item)"
          :option="item"
          :label="getSelectSpaceLabel(item)"
          :label-hidden="true"
          @update:model-value="selectSpace(item)"
          @click.stop="fileClicked([item, $event])"
        />
      </template>
      <template #icon>
        <oc-icon name="layout-grid" />
      </template>
      <template #name="{ item }">
        <span
          class="spaces-table-space-name"
          :data-test-space-name="item.name"
          v-text="item.name"
        />
      </template>
      <template #manager="{ item }">
        {{ getManagerNames(item) }}
      </template>
      <template #members="{ item }">
        {{ getMemberCount(item) }}
      </template>
      <template #totalQuota="{ item }"> {{ getTotalQuota(item) }}</template>
      <template #usedQuota="{ item }"> {{ getUsedQuota(item) }}</template>
      <template #remainingQuota="{ item }"> {{ getRemainingQuota(item) }}</template>
      <template #mdate="{ item }">
        <span
          v-oc-tooltip="formatDate(item.mdate)"
          tabindex="0"
          v-text="formatDateRelative(item.mdate)"
        />
      </template>
      <template #status="{ item }">
        <oc-icon
          v-oc-tooltip="item.disabled ? $gettext('Disabled') : $gettext('Enabled')"
          :name="item.disabled ? 'stop-circle' : 'play-circle'"
          :accessible-label="
            item.disabled ? $gettext('Space is disabled') : $gettext('Space is enabled')
          "
          size="small"
          fill-type="line"
        />
      </template>
      <template #actions="{ item }">
        <div class="spaces-list-actions">
          <oc-button
            v-oc-tooltip="spaceDetailsLabel"
            :aria-label="spaceDetailsLabel"
            appearance="raw"
            class="oc-mr-xs quick-action-button spaces-table-btn-details oc-p-xs"
            @click.stop.prevent="showDetailsForSpace(item)"
          >
            <oc-icon name="information" fill-type="line" />
          </oc-button>
          <context-menu-quick-action
            ref="contextMenuButtonRef"
            :item="item"
            class="spaces-table-btn-action-dropdown"
            @quick-action-clicked="showContextMenuOnBtnClick($event, item)"
          >
            <template #contextMenu>
              <slot name="contextMenu" :space="item" />
            </template>
          </context-menu-quick-action>
        </div>
      </template>
      <template #footer>
        <pagination :pages="totalPages" :current-page="currentPage" />
        <div class="oc-text-center oc-width-1-1 oc-my-s">
          <p class="oc-text-muted">{{ footerTextTotal }}</p>
          <p v-if="filterTerm" class="oc-text-muted">{{ footerTextFilter }}</p>
        </div>
      </template>
    </oc-table>
  </div>
</template>

<script lang="ts" setup>
import {
  formatDateFromJSDate,
  formatRelativeDateFromJSDate,
  displayPositionedDropdown,
  formatFileSize,
  defaultFuseOptions,
  useKeyboardActions,
  ContextMenuBtnClickEventData,
  useIsTopBarSticky
} from '@ownclouders/web-pkg'
import { ComponentPublicInstance, computed, nextTick, onMounted, ref, unref, watch } from 'vue'
import { getSpaceManagers, SpaceResource } from '@ownclouders/web-client'
import Mark from 'mark.js'
import Fuse from 'fuse.js'
import { useGettext } from 'vue3-gettext'
import { eventBus, SortDir } from '@ownclouders/web-pkg'
import { SideBarEventTopics } from '@ownclouders/web-pkg'
import { ContextMenuQuickAction } from '@ownclouders/web-pkg'
import { useFileListHeaderPosition, useRoute, useRouter, usePagination } from '@ownclouders/web-pkg'
import { Pagination } from '@ownclouders/web-pkg'
import { perPageDefault, perPageStoragePrefix } from '../../defaults'
import { findIndex } from 'lodash-es'
import {
  useKeyboardTableMouseActions,
  useKeyboardTableNavigation
} from '../../composables/keyboardActions'
import { useSpaceSettingsStore } from '../../composables'
import { storeToRefs } from 'pinia'

const router = useRouter()
const route = useRoute()
const language = useGettext()
const { $gettext } = language
const { isSticky } = useIsTopBarSticky()

const { y: fileListHeaderY } = useFileListHeaderPosition('#admin-settings-app-bar')
const contextMenuButtonRef = ref(undefined)
const sortBy = ref('name')
const sortDir = ref(SortDir.Asc)
const filterTerm = ref('')
const markInstance = ref(undefined)

const lastSelectedSpaceIndex = ref(0)
const lastSelectedSpaceId = ref(null)

const spaceSettingsStore = useSpaceSettingsStore()
const { spaces, selectedSpaces } = storeToRefs(spaceSettingsStore)

const highlighted = computed(() => unref(selectedSpaces).map((s) => s.id))
const footerTextTotal = computed(() => {
  return $gettext('%{spaceCount} spaces in total', {
    spaceCount: unref(spaces).length.toString()
  })
})
const footerTextFilter = computed(() => {
  return $gettext('%{spaceCount} matching spaces', {
    spaceCount: unref(items).length.toString()
  })
})

const orderBy = (list: SpaceResource[], prop: string, desc: boolean) => {
  return [...list].sort((s1, s2) => {
    let a: string, b: string
    const numeric = ['totalQuota', 'usedQuota', 'remainingQuota'].includes(prop)

    switch (prop) {
      case 'members':
        a = getMemberCount(s1).toString()
        b = getMemberCount(s2).toString()
        break
      case 'totalQuota':
        a = (s1.spaceQuota.total || 0).toString()
        b = (s2.spaceQuota.total || 0).toString()
        break
      case 'usedQuota':
        a = (s1.spaceQuota.used || 0).toString()
        b = (s2.spaceQuota.used || 0).toString()
        break
      case 'remainingQuota':
        a = (s1.spaceQuota.remaining || 0).toString()
        b = (s2.spaceQuota.remaining || 0).toString()
        break
      case 'status':
        a = s1.disabled.toString()
        b = s2.disabled.toString()
        break
      default:
        a = s1[prop as keyof SpaceResource].toString() || ''
        b = s2[prop as keyof SpaceResource].toString() || ''
    }

    return desc
      ? b.localeCompare(a, undefined, { numeric })
      : a.localeCompare(b, undefined, { numeric })
  })
}
const items = computed(() =>
  orderBy(filter(unref(spaces), unref(filterTerm)), unref(sortBy), unref(sortDir) === SortDir.Desc)
)
const {
  items: paginatedItems,
  page: currentPage,
  total: totalPages
} = usePagination({ items, perPageDefault, perPageStoragePrefix })

const keyActions = useKeyboardActions()
useKeyboardTableNavigation(
  keyActions,
  paginatedItems,
  selectedSpaces,
  lastSelectedSpaceIndex,
  lastSelectedSpaceId
)
useKeyboardTableMouseActions(
  keyActions,
  paginatedItems,
  selectedSpaces,
  lastSelectedSpaceIndex,
  lastSelectedSpaceId
)

watch(currentPage, () => {
  unselectAllSpaces()
})

const allSpacesSelected = computed(() => {
  return unref(paginatedItems).length === unref(selectedSpaces).length
})

const handleSort = (event: { sortBy: string; sortDir: SortDir }) => {
  sortBy.value = event.sortBy
  sortDir.value = event.sortDir
}
const filter = (spaces: SpaceResource[], filterTerm: string) => {
  if (!(filterTerm || '').trim()) {
    return spaces
  }
  const searchEngine = new Fuse(spaces, { ...defaultFuseOptions, keys: ['name'] })
  return searchEngine.search(filterTerm).map((r) => r.item)
}
const isSpaceSelected = (space: SpaceResource) => {
  return unref(selectedSpaces).some((s) => s.id === space.id)
}

const fields = computed(() => [
  {
    name: 'select',
    title: '',
    type: 'slot',
    width: 'shrink',
    headerType: 'slot'
  },
  {
    name: 'icon',
    title: '',
    type: 'slot',
    width: 'shrink'
  },
  {
    name: 'name',
    title: $gettext('Name'),
    type: 'slot',
    sortable: true,
    tdClass: 'mark-element'
  },
  {
    name: 'status',
    title: $gettext('Status'),
    type: 'slot',
    sortable: true
  },
  {
    name: 'manager',
    title: $gettext('Manager'),
    type: 'slot'
  },
  {
    name: 'members',
    title: $gettext('Members'),
    type: 'slot',
    sortable: true
  },
  {
    name: 'totalQuota',
    title: $gettext('Total quota'),
    type: 'slot',
    sortable: true
  },
  {
    name: 'usedQuota',
    title: $gettext('Used quota'),
    type: 'slot',
    sortable: true
  },
  {
    name: 'remainingQuota',
    title: $gettext('Remaining quota'),
    type: 'slot',
    sortable: true
  },
  {
    name: 'mdate',
    title: $gettext('Modified'),
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
])

const getManagerNames = (space: SpaceResource) => {
  const allManagers = getSpaceManagers(space)
  const managers = allManagers.length > 2 ? allManagers.slice(0, 2) : allManagers
  let managerStr = managers
    .map(({ grantedTo }) => (grantedTo.user || grantedTo.group).displayName)
    .join(', ')
  if (allManagers.length > 2) {
    managerStr += `... +${allManagers.length - 2}`
  }
  return managerStr
}
const formatDate = (date: string) => {
  return formatDateFromJSDate(new Date(date), language.current)
}
const formatDateRelative = (date: string) => {
  return formatRelativeDateFromJSDate(new Date(date), language.current)
}
const getTotalQuota = (space: SpaceResource) => {
  if (space.spaceQuota.total === 0) {
    return $gettext('Unrestricted')
  }

  return formatFileSize(space.spaceQuota.total, language.current)
}
const getUsedQuota = (space: SpaceResource) => {
  if (space.spaceQuota.used === undefined) {
    return '-'
  }
  return formatFileSize(space.spaceQuota.used, language.current)
}
const getRemainingQuota = (space: SpaceResource) => {
  if (space.spaceQuota.remaining === undefined) {
    return '-'
  }
  return formatFileSize(space.spaceQuota.remaining, language.current)
}
const getMemberCount = (space: SpaceResource) => {
  return Object.keys(space.members).length
}

const getSelectSpaceLabel = (space: SpaceResource) => {
  return $gettext('Select %{ space }', { space: space.name }, true)
}

onMounted(() => {
  nextTick(() => {
    markInstance.value = new Mark('.mark-element')
  })
})

watch(filterTerm, async () => {
  await unref(router).push({ ...unref(route), query: { ...unref(route).query, page: '1' } })
})

watch([filterTerm, paginatedItems], () => {
  unref(markInstance)?.unmark()
  unref(markInstance)?.mark(unref(filterTerm), {
    element: 'span',
    className: 'mark-highlight'
  })
})

const fileClicked = (data: [SpaceResource, MouseEvent]) => {
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

  unselectAllSpaces()
  selectSpace(resource)
}

const showContextMenuOnBtnClick = (data: ContextMenuBtnClickEventData, space: SpaceResource) => {
  const { dropdown, event } = data
  if (dropdown?.tippy === undefined) {
    return
  }
  if (!isSpaceSelected(space)) {
    spaceSettingsStore.setSelectedSpaces([space])
  }
  displayPositionedDropdown(dropdown.tippy, event, unref(contextMenuButtonRef))
}
const showContextMenuOnRightClick = (
  row: ComponentPublicInstance<unknown>,
  event: MouseEvent,
  space: SpaceResource
) => {
  event.preventDefault()
  const dropdown = row.$el.getElementsByClassName('spaces-table-btn-action-dropdown')[0]
  if (dropdown === undefined) {
    return
  }
  if (!isSpaceSelected(space)) {
    spaceSettingsStore.setSelectedSpaces([space])
  }
  displayPositionedDropdown(dropdown._tippy, event, unref(contextMenuButtonRef))
}

const spaceDetailsLabel = computed(() => {
  return $gettext('Show details')
})
const showDetailsForSpace = (space: SpaceResource) => {
  selectSpace(space)
  eventBus.publish(SideBarEventTopics.open)
}

const selectSpace = (selectedSpace: SpaceResource) => {
  lastSelectedSpaceIndex.value = findIndex(unref(spaces), (g) => g.id === selectedSpace.id)
  lastSelectedSpaceId.value = selectedSpace.id
  keyActions.resetSelectionCursor()

  const isSpaceSelected = unref(selectedSpaces).find((space) => space.id === selectedSpace.id)
  if (!isSpaceSelected) {
    return spaceSettingsStore.addSelectedSpace(selectedSpace)
  }

  spaceSettingsStore.setSelectedSpaces(
    unref(selectedSpaces).filter((space) => space.id !== selectedSpace.id)
  )
}

const unselectAllSpaces = () => {
  spaceSettingsStore.setSelectedSpaces([])
}

const selectSpaces = (spaces: SpaceResource[]) => {
  spaceSettingsStore.setSelectedSpaces(spaces)
}
</script>

<style lang="scss">
#spaces-filter {
  width: 16rem;
}

.spaces-table {
  .oc-table-header-cell-actions,
  .oc-table-data-cell-actions {
    white-space: nowrap;
  }

  .oc-table-header-cell-manager,
  .oc-table-data-cell-manager,
  .oc-table-header-cell-remainingQuota,
  .oc-table-data-cell-remainingQuota {
    display: none;

    @media only screen and (min-width: 1200px) {
      display: table-cell;
    }
  }

  .oc-table-header-cell-totalQuota,
  .oc-table-data-cell-totalQuota,
  .oc-table-header-cell-usedQuota,
  .oc-table-data-cell-usedQuota {
    display: none;

    @media only screen and (min-width: 1600px) {
      display: table-cell;
    }
  }

  &-squashed {
    .oc-table-header-cell-manager,
    .oc-table-data-cell-manager,
    .oc-table-header-cell-totalQuota,
    .oc-table-data-cell-totalQuota,
    .oc-table-header-cell-usedQuota,
    .oc-table-data-cell-usedQuota {
      display: none;
    }

    .oc-table-header-cell-remainingQuota,
    .oc-table-data-cell-remainingQuota,
    .oc-table-header-cell-mdate,
    .oc-table-data-cell-mdate {
      display: none;
      @media only screen and (min-width: 1400px) {
        display: table-cell;
      }
    }
  }
}
</style>
