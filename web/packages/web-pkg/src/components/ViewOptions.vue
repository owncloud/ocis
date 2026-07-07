<template>
  <div class="oc-flex oc-flex-middle">
    <oc-hidden-announcer :announcement="viewOptionsAnnouncement" level="polite" />
    <div
      v-if="viewModes.length > 1"
      class="viewmode-switch-buttons oc-button-group oc-visible@s oc-mr-s"
    >
      <oc-button
        v-for="viewMode in viewModes"
        :key="viewMode.name"
        v-oc-tooltip="$gettext(viewMode.label)"
        :class="viewMode.name"
        :appearance="viewModeCurrent === viewMode.name ? 'filled' : 'outline'"
        :aria-label="$gettext(viewMode.label)"
        :aria-pressed="viewModeCurrent === viewMode.name"
        variation="primary"
        @click="setViewMode(viewMode)"
      >
        <oc-icon
          :name="viewMode.icon.name"
          :fill-type="viewMode.icon.fillType"
          size="small"
          variation="inherit"
        />
      </oc-button>
    </div>
    <oc-button
      id="files-view-options-btn"
      key="files-view-options-btn"
      v-oc-tooltip="viewOptionsButtonLabel"
      data-testid="files-view-options-btn"
      :aria-label="viewOptionsButtonLabel"
      appearance="raw"
      variation="primary"
      class="oc-my-s oc-p-xs"
    >
      <oc-icon name="settings-3" fill-type="line" />
    </oc-button>
    <oc-drop
      drop-id="files-view-options-drop"
      toggle="#files-view-options-btn"
      mode="click"
      class="oc-width-auto"
      padding-size="medium"
    >
      <oc-list>
        <li v-if="shouldShowFlatListToggle" class="files-view-options-list-item">
          <oc-switch
            v-model:checked="showFlatList"
            data-testid="files-switch-flat-list"
            :label="$gettext('Flat List (A-Z)')"
            @update:checked="toggleFlatList"
          />
        </li>
        <li v-if="hasHiddenFiles" class="files-view-options-list-item">
          <oc-switch
            v-model:checked="hiddenFilesShownModel"
            data-testid="files-switch-hidden-files"
            :label="$gettext('Show hidden files')"
            @update:checked="updateHiddenFilesShownModel"
          />
        </li>
        <li v-if="hasFileExtensions" class="files-view-options-list-item">
          <oc-switch
            v-model:checked="fileExtensionsShownModel"
            data-testid="files-switch-files-extensions-files"
            :label="$gettext('Show file extensions')"
            @update:checked="updateFileExtensionsShownModel"
          />
        </li>
        <li v-if="hasPagination" class="files-view-options-list-item">
          <oc-page-size
            v-if="!queryParamsLoading"
            :selected="queryItemAsString(itemsPerPageCurrent)"
            data-testid="files-pagination-size"
            :label="$gettext('Items per page')"
            :options="paginationOptions"
            class="files-pagination-size"
            @change="setItemsPerPage"
          />
        </li>
        <li
          v-if="viewModes.find((v) => v.name === FolderViewModeConstants.name.tiles)"
          class="files-view-options-list-item oc-flex oc-flex-between oc-flex-middle"
        >
          <label for="tiles-size-slider" v-text="$gettext('Tile size')" />
          <input
            id="tiles-size-slider"
            v-model="viewSizeCurrent"
            type="range"
            :min="1"
            :max="viewSizeMax"
            class="oc-range"
            data-testid="files-tiles-size-slider"
          />
        </li>
      </oc-list>
    </oc-drop>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, unref, watch } from 'vue'
import { useGettext } from 'vue3-gettext'
import {
  queryItemAsString,
  useRoute,
  useRouteQuery,
  useRouteQueryPersisted,
  useRouter,
  PaginationConstants,
  FolderViewModeConstants,
  useRouteName,
  useResourcesStore,
  useViewSizeMax
} from '../composables'
import { FolderView } from '../ui/types'
import { storeToRefs } from 'pinia'

interface Props {
  shouldShowFlatListToggle?: boolean
  hasHiddenFiles?: boolean
  hasFileExtensions?: boolean
  hasPagination?: boolean
  paginationOptions?: string[]
  perPageQueryName?: string
  perPageDefault?: string
  perPageStoragePrefix: string
  viewModeDefault?: string
  viewModes?: FolderView[]
}
const {
  shouldShowFlatListToggle = true,
  hasHiddenFiles = true,
  hasFileExtensions = true,
  hasPagination = true,
  paginationOptions = PaginationConstants.options,
  perPageQueryName = PaginationConstants.perPageQueryName,
  perPageDefault = PaginationConstants.perPageDefault,
  perPageStoragePrefix,
  viewModeDefault = FolderViewModeConstants.defaultModeName,
  viewModes = []
} = defineProps<Props>()
const router = useRouter()
const currentRoute = useRoute()
const { $gettext, $pgettext } = useGettext()

const resourcesStore = useResourcesStore()
const { setAreHiddenFilesShown, setAreFileExtensionsShown, setShouldShowFlatList } = resourcesStore
const { areHiddenFilesShown, areFileExtensionsShown, shouldShowFlatList } =
  storeToRefs(resourcesStore)

const queryParamsLoading = ref(false)
const viewOptionsAnnouncement = ref('')

const currentPageQuery = useRouteQuery('page')
const currentPage = computed(() => {
  if (!unref(currentPageQuery)) {
    return 1
  }
  return parseInt(queryItemAsString(unref(currentPageQuery)))
})
const itemsPerPageQuery = useRouteQueryPersisted({
  name: perPageQueryName,
  defaultValue: perPageDefault,
  storagePrefix: perPageStoragePrefix
})

const routeName = useRouteName()
const viewModeQuery = useRouteQueryPersisted({
  name: `${unref(routeName)}-${FolderViewModeConstants.queryName}`,
  defaultValue: viewModeDefault
})

const viewSizeQuery = useRouteQueryPersisted({
  name: FolderViewModeConstants.tilesSizeQueryName,
  defaultValue: FolderViewModeConstants.tilesSizeDefault.toString()
})

function toggleFlatList(event: boolean) {
  setShouldShowFlatList(event)
  viewOptionsAnnouncement.value = event
    ? $pgettext('Accessibility announcement when flat list view is enabled', 'Flat list enabled')
    : $pgettext('Accessibility announcement when folder tree view is enabled', 'Flat list disabled')
}
function updateHiddenFilesShownModel(event: boolean) {
  setAreHiddenFilesShown(event)
  viewOptionsAnnouncement.value = event
    ? $pgettext(
        'Accessibility announcement when hidden files are shown',
        'Hidden files are now visible'
      )
    : $pgettext(
        'Accessibility announcement when hidden files are hidden',
        'Hidden files are now hidden'
      )
}
function updateFileExtensionsShownModel(event: boolean) {
  setAreFileExtensionsShown(event)
  viewOptionsAnnouncement.value = event
    ? $pgettext(
        'Accessibility announcement when file extensions are shown',
        'File extensions are now visible'
      )
    : $pgettext(
        'Accessibility announcement when file extensions are hidden',
        'File extensions are now hidden'
      )
}

const setItemsPerPage = (itemsPerPage: string) => {
  return router.replace({
    query: {
      ...unref(currentRoute).query,
      [perPageQueryName]: itemsPerPage,
      ...(unref(currentPage) > 1 && { page: '1' })
    }
  })
}

const setViewMode = (mode: FolderView) => {
  viewModeQuery.value = mode.name
  if (mode.name === 'resource-table-condensed') {
    viewOptionsAnnouncement.value = $pgettext(
      'Accessibility announcement when condensed table view is selected',
      'Condensed table view selected'
    )
  } else if (mode.name === 'resource-tiles') {
    viewOptionsAnnouncement.value = $pgettext(
      'Accessibility announcement when tiles view is selected',
      'Tiles view selected'
    )
  } else {
    viewOptionsAnnouncement.value = $pgettext(
      'Accessibility announcement when table view is selected',
      'Table view selected'
    )
  }
}

watch(
  [itemsPerPageQuery, viewModeQuery, viewSizeQuery],
  (params) => {
    queryParamsLoading.value = params.some((p) => !p)
  },
  { immediate: true, deep: true }
)

const viewSizeMax = useViewSizeMax()
const viewModeCurrent = viewModeQuery
const viewSizeCurrent = viewSizeQuery
const itemsPerPageCurrent = itemsPerPageQuery

const viewOptionsButtonLabel = $gettext('Display customization options of the files list')
const showFlatList = computed(() => unref(shouldShowFlatList))
const hiddenFilesShownModel = computed(() => unref(areHiddenFilesShown))
const fileExtensionsShownModel = computed(() => unref(areFileExtensionsShown))
</script>

<style lang="scss" scoped>
.viewmode-switch-buttons {
  flex-flow: initial;
}

.viewmode-switch-buttons.oc-button-group {
  outline: 1px solid var(--oc-color-swatch-primary-default);
}

#files-view-options-btn {
  vertical-align: middle;
  border: 3px solid transparent;
  &:hover {
    background-color: var(--oc-color-background-hover);
    border-radius: 3px;
  }
}

.files-view-options-list-item {
  &:not(:last-child) {
    margin-bottom: var(--oc-space-medium);
  }

  & > * {
    display: flex;
    justify-content: space-between;
  }

  & + & {
    margin-top: var(--oc-space-small);
  }
}

.oc-range {
  -webkit-appearance: none;
  -webkit-transition: 0.2s;
  border-radius: 0.3rem;
  background: var(--oc-color-border);
  height: 0.5rem;
  opacity: 0.7;
  outline: none;
  transition: opacity 0.2s;
  width: 100%;
  max-width: 50%;

  &:hover {
    opacity: 1;
  }

  &::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    background: var(--oc-color-swatch-primary-default);
    border-radius: 50%;
    cursor: pointer;
    height: 1rem;
    width: 1rem;
  }

  &::-moz-range-thumb {
    background: var(--oc-color-swatch-primary-default);
    border-radius: 50%;
    cursor: pointer;
    height: 1rem;
    width: 1rem;
  }
}
</style>
