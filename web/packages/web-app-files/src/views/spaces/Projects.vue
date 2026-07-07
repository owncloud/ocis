<template>
  <div class="oc-flex oc-width-1-1">
    <files-view-wrapper>
      <app-bar
        :breadcrumbs="breadcrumbs"
        :has-bulk-actions="true"
        :has-hidden-files="false"
        :has-file-extensions="false"
        :has-pagination="false"
        :is-side-bar-open="isSideBarOpen"
        :view-modes="viewModes"
        :view-mode-default="FolderViewModeConstants.name.tiles"
      >
        <template #actions>
          <create-space v-if="hasCreatePermission" class="oc-mr-s" />
          <div v-if="!selectedResourcesIds?.length" class="oc-flex oc-flex-middle oc-pl-s">
            <span v-text="$gettext('Learn about spaces')" />
            <oc-contextual-helper
              :list="spacesHelpList"
              :title="$gettext('Spaces')"
              class="oc-ml-xs"
            />
          </div>
        </template>
      </app-bar>
      <app-loading-spinner v-if="areResourcesLoading" />
      <template v-else>
        <no-content-message
          v-if="!spaces.length"
          id="files-spaces-empty"
          class="files-empty"
          icon="layout-grid"
        >
          <template #message>
            <span v-text="$gettext('You don\'t have access to any spaces')" />
          </template>
        </no-content-message>
        <div v-else class="spaces-list">
          <div
            class="spaces-list-filters oc-flex oc-flex-between oc-flex-wrap oc-flex-bottom oc-mx-m oc-mb-m"
          >
            <fieldset class="oc-flex">
              <div>
                <legend class="oc-mr-m oc-flex oc-flex-middle">
                  <oc-icon name="filter-2" class="oc-mr-xs" />
                  <span v-text="$gettext('Filter:')" />
                </legend>
              </div>
              <item-filter-toggle
                :filter-label="$gettext('Include disabled')"
                filter-name="includeDisabled"
                class="spaces-list-filter-include-disabled oc-mr-s"
              />
            </fieldset>
            <oc-text-input
              id="spaces-filter"
              v-model="filterTerm"
              :label="$gettext('Search')"
              autocomplete="off"
            />
          </div>
          <component
            :is="folderView.component"
            v-model:selected-ids="selectedResourcesIds"
            resource-type="space"
            :resources="paginatedItems"
            :fields-displayed="tableDisplayFields"
            :sort-fields="sortFields"
            :sort-by="sortBy"
            :sort-dir="sortDir"
            :header-position="fileListHeaderY"
            :is-side-bar-open="isSideBarOpen"
            :view-size="viewSize"
            v-bind="folderView.componentAttrs?.()"
            @sort="handleSort"
            @item-visible="loadPreview({ space: $event, resource: $event })"
          >
            <template #image="{ resource }">
              <template v-if="viewMode === FolderViewModeConstants.name.tiles">
                <img
                  v-if="resource.thumbnail"
                  class="tile-preview"
                  :src="resource.thumbnail"
                  :alt="resource.name"
                />
              </template>
              <template v-else>
                <img
                  v-if="resource.thumbnail"
                  class="table-preview oc-mr-s"
                  :src="resource.thumbnail"
                  alt=""
                  width="33"
                  height="33"
                />
                <resource-icon v-else class="oc-mr-s" :resource="resource" />
              </template>
            </template>
            <template #actions="{ resource }">
              <oc-button
                v-if="!resource.disabled"
                v-oc-tooltip="showSpaceMemberLabel"
                class="spaces-list-show-members-button"
                :aria-label="showSpaceMemberLabel"
                appearance="raw"
                @click="openSidebarSharePanel(resource as SpaceResource)"
              >
                <oc-icon name="group" fill-type="line" />
              </oc-button>
            </template>
            <template #contextMenu="{ resource }">
              <space-context-actions
                :action-options="{ resources: [resource] as SpaceResource[] }"
              />
            </template>
            <template #footer>
              <pagination :pages="totalPages" :current-page="currentPage" />
              <div class="oc-text-center oc-width-1-1 oc-my-s">
                <p class="oc-text-muted">{{ footerTextTotal }}</p>
                <p v-if="filterTerm" class="oc-text-muted">{{ footerTextFilter }}</p>
              </div>
            </template>
            <!--- table -->
            <template #status="{ resource }">
              <span v-if="resource.disabled" class="oc-flex oc-flex-middle">
                <oc-icon name="stop-circle" fill-type="line" class="oc-mr-s" /><span
                  v-text="$gettext('Disabled')"
                />
              </span>
              <span v-else class="oc-flex oc-flex-middle">
                <oc-icon name="play-circle" fill-type="line" class="oc-mr-s" /><span
                  v-text="$gettext('Enabled')"
                />
              </span>
            </template>
            <template #manager="{ resource }">
              {{ getManagerNames(resource) }}
            </template>
            <template #members="{ resource }">
              {{ getMemberCount(resource) }}
            </template>
            <template #totalQuota="{ resource }">
              {{ getTotalQuota(resource) }}
            </template>
            <template #usedQuota="{ resource }"> {{ getUsedQuota(resource) }}</template>
            <template #remainingQuota="{ resource }"> {{ getRemainingQuota(resource) }}</template>
          </component>
        </div>
      </template>
    </files-view-wrapper>
    <file-side-bar
      :is-open="isSideBarOpen"
      :active-panel="sideBarActivePanel"
      :space="selectedSpace"
    />
  </div>
</template>

<script lang="ts" setup>
import { onMounted, computed, unref, ref, watch, nextTick, onBeforeUnmount } from 'vue'
import { useTask } from 'vue-concurrency'
import Mark from 'mark.js'
import Fuse from 'fuse.js'

import {
  AppLoadingSpinner,
  useResourcesStore,
  useSpacesStore,
  useExtensionRegistry,
  ItemFilterToggle,
  useRouteQuery,
  queryItemAsString,
  useLoadPreview
} from '@ownclouders/web-pkg'

import { AppBar, CreateSpace } from '@ownclouders/web-pkg'
import {
  useAbility,
  useClientService,
  FolderViewModeConstants,
  useRouteQueryPersisted,
  useSort,
  useRouteName,
  usePagination,
  useRouter,
  useRoute,
  Pagination,
  FileSideBar,
  NoContentMessage,
  useCapabilityStore
} from '@ownclouders/web-pkg'
import SpaceContextActions from '../../components/Spaces/SpaceContextActions.vue'
import {
  getSpaceManagers,
  isProjectSpaceResource,
  ProjectSpaceResource,
  SpaceResource
} from '@ownclouders/web-client'
import FilesViewWrapper from '../../components/FilesViewWrapper.vue'
import { eventBus } from '@ownclouders/web-pkg'
import { SideBarEventTopics, useSideBar } from '@ownclouders/web-pkg'
import { sortFields as availableSortFields, translateSortFields } from '@ownclouders/web-pkg'
import { defaultFuseOptions, formatFileSize, ResourceIcon } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { useKeyboardActions } from '@ownclouders/web-pkg'
import {
  useKeyboardTableNavigation,
  useKeyboardTableMouseActions,
  useKeyboardTableActions
} from '../../composables/keyboardActions'
import { orderBy } from 'lodash-es'
import { useResourcesViewDefaults } from '../../composables'
import { folderViewsProjectSpacesExtensionPoint } from '../../extensionPoints'

const spacesStore = useSpacesStore()
const router = useRouter()
const route = useRoute()
const clientService = useClientService()
const { can } = useAbility()
const language = useGettext()
const { $gettext } = language
const filterTerm = ref('')
const markInstance = ref(undefined)
const includeDisabledParam = useRouteQuery('q_includeDisabled')
const capabilityStore = useCapabilityStore()

const { setSelection, initResourceList, clearResourceList, setAncestorMetaData } =
  useResourcesStore()

const loadResourcesTask = useTask(function* (signal) {
  clearResourceList()
  setAncestorMetaData({})
  yield spacesStore.reloadProjectSpaces({
    graphClient: clientService.graphAuthenticated,
    signal,
    isInVault: unref(route)?.params?.scope === 'vault'
  })
  initResourceList({ currentFolder: null, resources: unref(spaces) })
})

const {
  viewSize,
  fileListHeaderY,
  scrollToResourceFromRoute,
  areResourcesLoading,
  selectedResourcesIds,
  selectedResources
} = useResourcesViewDefaults({ loadResourcesTask })

let loadPreviewToken: string = null
const { isSideBarOpen, sideBarActivePanel } = useSideBar()

const runtimeSpaces = computed(() => {
  return spacesStore.spaces.filter((space) => isProjectSpaceResource(space)) || []
})
const selectedSpace = computed(() => {
  if (
    unref(selectedResources).length === 1 &&
    isProjectSpaceResource(unref(selectedResources)[0])
  ) {
    return unref(selectedResources)[0] as ProjectSpaceResource
  }
  return null
})

const tableDisplayFields = [
  'image',
  'name',
  'manager',
  'members',
  'totalQuota',
  'usedQuota',
  'remainingQuota',
  'status',
  'mdate'
]

const sortFields = translateSortFields(availableSortFields, language)
const {
  sortBy,
  sortDir,
  items: spaces,
  handleSort
} = useSort<SpaceResource>({
  items: runtimeSpaces,
  fields: sortFields
})
const filter = (spaces: Array<ProjectSpaceResource>, filterTerm: string) => {
  const includeDisabled = queryItemAsString(unref(includeDisabledParam)) === 'true'

  if (!includeDisabled) {
    spaces = spaces.filter((space) => space.disabled !== true)
  }

  if (!(filterTerm || '').trim()) {
    return spaces
  }

  const searchEngine = new Fuse(spaces, { ...defaultFuseOptions, keys: ['name'] })
  return searchEngine.search(filterTerm).map((r) => r.item)
}
const items = computed(() =>
  orderBy(
    filter(unref(spaces), unref(filterTerm)),
    [
      (item: SpaceResource) =>
        typeof item[unref(sortBy)] === 'string'
          ? item[unref(sortBy)].toLowerCase()
          : item[unref(sortBy)]
    ],
    unref(sortDir)
  )
)

const {
  items: paginatedItems,
  page: currentPage,
  total: totalPages
} = usePagination({
  items,
  perPageDefault: '50',
  perPageStoragePrefix: 'spaces-list'
})

watch(filterTerm, async () => {
  const instance = unref(markInstance)
  if (!instance) {
    return
  }
  await router.push({ ...unref(route), query: { ...unref(route).query, page: '1' } })
  instance.unmark()
  instance.mark(unref(filterTerm), {
    element: 'span',
    className: 'mark-highlight',
    exclude: ['th *', 'tfoot *']
  })
})

const hasCreatePermission = computed(() => can('create-all', 'Drive'))
const canAccessVault = computed(() => capabilityStore.vaultEnabled && can('read-all', 'Vault'))

const extensionRegistry = useExtensionRegistry()
const viewModes = computed(() => {
  return [
    ...extensionRegistry
      .requestExtensions(folderViewsProjectSpacesExtensionPoint)
      .map((e) => e.folderView)
  ]
})

const routeName = useRouteName()

const viewMode = useRouteQueryPersisted({
  name: `${unref(routeName)}-${FolderViewModeConstants.queryName}`,
  defaultValue: FolderViewModeConstants.name.tiles
})
const viewModeStr = computed(() => queryItemAsString(unref(viewMode)))

const { loadPreview } = useLoadPreview(viewModeStr)

const keyActions = useKeyboardActions()
useKeyboardTableNavigation(keyActions, runtimeSpaces, viewMode)
useKeyboardTableMouseActions(keyActions, viewMode)
useKeyboardTableActions(keyActions)

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

onMounted(async () => {
  await loadResourcesTask.perform()

  loadPreviewToken = eventBus.subscribe(
    'app.files.spaces.uploaded-image',
    (space: SpaceResource) => {
      loadPreview({ space, resource: space })
    }
  )
  scrollToResourceFromRoute(unref(spaces), 'files-app-bar')
  nextTick(() => {
    markInstance.value = new Mark('.spaces-table')
  })
})

onBeforeUnmount(() => {
  eventBus.unsubscribe('app.files.spaces.uploaded-image', loadPreviewToken)
})

const footerTextTotal = computed(() => {
  const disabledSpaces = unref(spaces).filter((space) => space.disabled === true)

  if (!disabledSpaces.length) {
    return $gettext('%{spaceCount} spaces in total', {
      spaceCount: unref(spaces).length.toString()
    })
  }

  return $gettext('%{spaceCount} spaces in total (including %{disabledSpaceCount} disabled)', {
    spaceCount: unref(spaces).length.toString(),
    disabledSpaceCount: disabledSpaces.length.toString()
  })
})
const footerTextFilter = computed(() => {
  return $gettext('%{spaceCount} matching spaces', {
    spaceCount: unref(items).length.toString()
  })
})

const folderView = computed(() => {
  const viewModeName = unref(viewMode) || FolderViewModeConstants.name.tiles
  return unref(viewModes).find((v) => v.name === viewModeName)
})

const spacesHelpList = computed(() => {
  return [
    {
      text: $gettext('Spaces are special folders for making files accessible to a team.')
    },
    {
      text: $gettext(
        'Spaces belong to a team and not to a single person. Even if members are removed, the files remain in the Space so that the team can continue to work on the files.'
      )
    },
    {
      text: $gettext(
        'Members with the Manager role can add or remove other members from the Space.'
      )
    },
    {
      text: $gettext('A Space can have multiple Managers. Each Space has at least one Manager.')
    }
  ]
})
const getBreadcrumbText = () => {
  if (!unref(canAccessVault)) {
    return $gettext('Spaces')
  }

  if (unref(route).params.scope === 'vault') {
    return $gettext('Vault')
  }

  return $gettext('Drive')
}

const breadcrumbs = computed(() => {
  return [
    {
      text: getBreadcrumbText(),
      onClick: () => loadResourcesTask.perform(),
      isStativNav: true
    }
  ]
})

const showSpaceMemberLabel = computed(() => $gettext('Show members'))
function openSidebarSharePanel(space: SpaceResource) {
  setSelection([space.id])
  eventBus.publish(SideBarEventTopics.openWithPanel, 'space-share')
}
</script>

<style lang="scss">
#files-spaces-empty {
  height: 75dvh;
}

.table-preview {
  border-radius: 3px;
}

.state-trashed {
  .tile-preview,
  .tile-default-image > svg {
    filter: grayscale(100%);
    opacity: 80%;
  }
}
</style>
