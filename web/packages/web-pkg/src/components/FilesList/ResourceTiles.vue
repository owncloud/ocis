<template>
  <div id="tiles-view" class="oc-px-m oc-pt-l">
    <div class="oc-flex oc-flex-middle oc-mb-m oc-pb-s oc-tiles-controls">
      <oc-checkbox
        id="tiles-view-select-all"
        v-oc-tooltip="selectAllCheckboxLabel"
        class="oc-ml-s"
        size="large"
        :label="selectAllCheckboxLabel"
        :label-hidden="true"
        :disabled="resources.length === disabledResourceIds.length"
        :model-value="areAllResourcesSelected"
        @click.stop="toggleSelectionAll"
      />
      <div v-if="sortFields.length" class="oc-tile-sorting oc-ml-m">
        <oc-select
          class="oc-tiles-sort-select oc-flex oc-flex-middle"
          :model-value="currentSortField"
          :label="$gettext('Sort by')"
          :options="sortFields"
          :clearable="false"
          :searchable="false"
          :position-fixed="true"
          @update:model-value="selectSorting"
        />
      </div>
    </div>
    <oc-list class="oc-tiles">
      <li
        v-for="resource in resources"
        :key="resource.id"
        class="oc-tiles-item has-item-context-menu"
      >
        <resource-tile
          :ref="(el) => (tileRefs.tiles[resource.id] = el as ResourceTileRef)"
          :resource="resource"
          :resource-route="getRoute(resource)"
          :is-resource-selected="isResourceSelected(resource)"
          :is-resource-clickable="isResourceClickable(resource)"
          :is-resource-disabled="isResourceDisabled(resource)"
          :is-extension-displayed="areFileExtensionsShown"
          :resource-icon-size="resourceIconSize"
          :draggable="dragDrop"
          :lazy="lazy"
          @vue:mounted="
            $emit('rowMounted', resource, tileRefs.tiles[resource.id], ImageDimension.Tile)
          "
          @contextmenu="showContextMenu($event, resource, tileRefs.tiles[resource.id])"
          @click="emitTileClick(resource)"
          @dragstart="dragStart(resource, $event)"
          @dragenter.prevent="setDropStyling(resource, false, $event)"
          @dragleave.prevent="setDropStyling(resource, true, $event)"
          @drop="fileDropped(resource, $event)"
          @dragover="$event.preventDefault()"
          @item-visible="$emit('itemVisible', resource)"
        >
          <template #selection>
            <oc-checkbox
              v-if="!isLocationPicker && !isFilePicker"
              :label="getResourceCheckboxLabel(resource)"
              :label-hidden="true"
              size="large"
              class="oc-flex-inline oc-p-s"
              :disabled="!isSpaceResource(resource) && isResourceDisabled(resource)"
              :model-value="isResourceSelected(resource)"
              @click.stop.prevent="toggleTile([resource, $event])"
            />
          </template>
          <template #imageField>
            <slot name="image" :resource="resource" />
          </template>
          <template #indicators>
            <resource-status-indicators
              :space="space"
              :resource="resource"
              :filter="(indicator) => indicator.category === 'system'"
              :disable-handler="!isSpaceResource(resource) && isResourceDisabled(resource)"
            />
          </template>
          <template #actions>
            <slot name="actions" :resource="resource" />
          </template>
          <template #contextMenu>
            <context-menu-quick-action
              v-if="isSpaceResource(resource) || !isResourceDisabled(resource)"
              :ref="(el) => (tileRefs.dropBtns[resource.id] = el as ContextMenuQuickActionRef)"
              :item="resource"
              class="resource-tiles-btn-action-dropdown"
              @quick-action-clicked="showContextMenuOnBtnClick($event, resource, resource.id)"
            >
              <template #contextMenu>
                <slot name="contextMenu" :resource="resource" />
              </template>
            </context-menu-quick-action>
          </template>
        </resource-tile>
      </li>
      <li
        v-for="index in ghostTilesCount"
        :key="`ghost-tile-${index}`"
        class="ghost-tile"
        :aria-hidden="true"
      />
    </oc-list>
    <Teleport v-if="dragItem" to="body">
      <resource-ghost-element ref="ghostElementRef" :preview-items="[dragItem, ...dragSelection]" />
    </Teleport>
    <div class="oc-tiles-footer">
      <slot name="footer" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  computed,
  ComponentPublicInstance,
  nextTick,
  onBeforeUnmount,
  onBeforeUpdate,
  onMounted,
  ref,
  unref,
  watch,
  type MaybeRef
} from 'vue'
import { useGettext } from 'vue3-gettext'
import {
  isPasswordProtectedFolderFileResource,
  isSpaceResource,
  Resource,
  SpaceResource
} from '@ownclouders/web-client'

// Constants should match what is being used in OcTable/ResourceTable
// Alignment regarding naming would be an API-breaking change and can
// Be done at a later point in time?
import { ContextMenuQuickAction } from '../ContextActions'
import { createLocationSpaces } from '../../router'
import {
  ContextMenuBtnClickEventData,
  createFileRouteOptions,
  CreateTargetRouteOptions,
  displayPositionedDropdown
} from '../../helpers'
import { eventBus } from '../../services'
import { ImageDimension } from '../../constants'
import ResourceTile from './ResourceTile.vue'
import ResourceGhostElement from './ResourceGhostElement.vue'
import {
  FolderViewModeConstants,
  SortDir,
  SortField,
  useResourceRouteResolver,
  useTileSize,
  useResourcesStore,
  useViewSizeMax,
  useEmbedMode,
  useCanBeOpenedWithSecureView,
  useFileActions,
  useGetMatchingSpace,
  embedModeFilePickMessageData,
  routeToContextQuery,
  useRouter
} from '../../composables'
import ResourceStatusIndicators from './ResourceStatusIndicators.vue'

type ResourceTileRef = ComponentPublicInstance<typeof ResourceTile>
type ContextMenuQuickActionRef = ComponentPublicInstance<typeof ContextMenuQuickAction>
type IconSize = 'large' | 'xlarge' | 'xxlarge' | 'xxxlarge'

interface Props {
  resources?: Resource[]
  selectedIds?: string[]
  targetRouteCallback?: (arg: CreateTargetRouteOptions) => unknown
  space?: SpaceResource
  sortFields?: SortField[]
  sortBy?: string
  sortDir?: SortDir
  viewSize?: number
  dragDrop?: boolean
  lazy?: boolean
}
interface Emits {
  (e: 'fileClick', data: { space: SpaceResource; resource: Resource }): void
  (e: 'fileDropped', resourceId: string): void
  (e: 'rowMounted', resource: Resource, tileRef: ResourceTileRef, dimension: ImageDimension): void
  (e: 'sort', sortData: { sortBy: string; sortDir: MaybeRef<SortDir> }): void
  (e: 'itemVisible', resource: Resource): void
  (e: 'update:selectedIds', selectedIds: string[]): void
}
const {
  resources = [],
  selectedIds = [],
  targetRouteCallback = undefined,
  space = null,
  sortFields = [],
  sortBy = undefined,
  sortDir = undefined,
  viewSize = FolderViewModeConstants.tilesSizeDefault,
  dragDrop = false,
  lazy = true
} = defineProps<Props>()
const emit = defineEmits<Emits>()
const { $gettext } = useGettext()
const router = useRouter()
const resourcesStore = useResourcesStore()
const { getDefaultAction } = useFileActions()
const { getMatchingSpace } = useGetMatchingSpace()
const { canBeOpenedWithSecureView } = useCanBeOpenedWithSecureView()
const {
  isEnabled: isEmbedModeEnabled,
  fileTypes: embedModeFileTypes,
  isLocationPicker,
  isFilePicker,
  postMessage
} = useEmbedMode()
const viewSizeMax = useViewSizeMax()
const viewSizeCurrent = computed(() => {
  return Math.min(unref(viewSizeMax), viewSize)
})

const areFileExtensionsShown = computed(() => resourcesStore.areFileExtensionsShown)

const selectAllCheckboxLabel = computed(() => {
  return unref(areAllResourcesSelected) ? $gettext('Clear selection') : $gettext('Select all')
})

const dragItem = ref()
const ghostElementRef = ref()

const tileRefs = ref({
  tiles: {} as Record<string, ResourceTileRef>,
  dropBtns: {} as Record<string, ContextMenuQuickActionRef>
})

const resourceRouteResolver = useResourceRouteResolver(
  {
    space: ref(space),
    targetRouteCallback: computed(() => targetRouteCallback)
  },
  emit
)

const getRoute = (resource: Resource) => {
  if (isSpaceResource(resource)) {
    return resource.disabled
      ? null
      : createLocationSpaces(
          'files-spaces-generic',
          createFileRouteOptions(resource as SpaceResource, {
            path: '',
            fileId: resource.fileId
          })
        )
  }

  if (resource.isFolder) {
    return resourceRouteResolver.createFolderLink({
      path: resource.path,
      fileId: resource.fileId,
      resource: resource
    })
  }

  let currentSpace = space
  if (!currentSpace) {
    currentSpace = getMatchingSpace(resource)
  }

  const action = getDefaultAction({ resources: [resource], space })
  if (!action?.route) {
    return null
  }

  return action.route({ space, resources: [resource] })
}
const emitTileClick = (resource: Resource) => {
  if (unref(isEmbedModeEnabled) && unref(isFilePicker)) {
    return postMessage<embedModeFilePickMessageData>('owncloud-embed:file-pick', {
      resource: JSON.parse(JSON.stringify(resource)),
      locationQuery: JSON.parse(JSON.stringify(routeToContextQuery(unref(router.currentRoute))))
    })
  }

  if (resource.type !== 'space' && resource.type !== 'folder') {
    resourceRouteResolver.createFileAction(resource)
  }
}

const showContextMenuOnBtnClick = (
  data: ContextMenuBtnClickEventData,
  item: Resource,
  index: string
) => {
  const { dropdown, event } = data
  if (dropdown?.tippy === undefined) {
    return
  }
  resourcesStore.setSelection([item.id])
  displayPositionedDropdown(dropdown.tippy, event, unref(tileRefs).dropBtns[index])
}

const isResourceSelected = (resource: Resource) => {
  return selectedIds.includes(resource.id)
}

const selectedResources = computed(() => {
  return resources.filter((resource) => selectedIds.includes(resource.id))
})

const isResourceClickable = (resource: Resource) => {
  if (isResourceDisabled(resource)) {
    return false
  }

  if (resource.isFolder || isPasswordProtectedFolderFileResource(resource.name)) {
    return true
  }

  if (!resource.canDownload() && !canBeOpenedWithSecureView(resource)) {
    return false
  }

  if (unref(isEmbedModeEnabled) && !unref(isFilePicker)) {
    return false
  }

  return true
}

const isResourceDisabled = (resource: Resource) => {
  if (unref(isEmbedModeEnabled) && unref(embedModeFileTypes)?.length) {
    return (
      !unref(embedModeFileTypes).includes(resource.extension) &&
      !unref(embedModeFileTypes).includes(resource.mimeType) &&
      !resource.isFolder
    )
  }

  if (isSpaceResource(resource) && resource.disabled) {
    return true
  }

  return resource.processing === true
}

const disabledResourceIds = computed(() => {
  return (
    resources
      ?.filter((resource) => isResourceDisabled(resource) === true)
      ?.map((resource) => resource.id) || []
  )
})

const emitSelect = (selectedIds: string[]) => {
  emit('update:selectedIds', selectedIds)
}

const toggleSelectionAll = () => {
  if (unref(areAllResourcesSelected)) {
    return emit('update:selectedIds', [])
  }

  emit(
    'update:selectedIds',
    resources
      .filter((resource) => !unref(disabledResourceIds).includes(resource.id))
      .map((resource) => resource.id)
  )
}

const showContextMenu = (
  event: MouseEvent,
  item: Resource,
  reference: ComponentPublicInstance<unknown>
) => {
  event.preventDefault()
  const drop = unref(tileRefs).tiles[item.id]?.$el.getElementsByClassName(
    'resource-tiles-btn-action-dropdown'
  )[0]

  if (drop === undefined) {
    return
  }
  if (!isResourceSelected(item)) {
    emitSelect([item.id])
  }
  displayPositionedDropdown(drop._tippy, event, reference)
}

const toggleTile = (data: [Resource, MouseEvent]) => {
  const resource = data[0]
  const eventData = data[1]

  if (eventData && eventData.metaKey) {
    return eventBus.publish('app.files.list.clicked.meta', resource)
  }
  if (eventData && eventData.shiftKey) {
    return eventBus.publish('app.files.list.clicked.shift', {
      resource,
      skipTargetSelection: false
    })
  }
  toggleSelection(resource)
}

const toggleSelection = (resource: Resource) => {
  const currentSelectedIds = !isResourceSelected(resource)
    ? [...selectedIds, resource.id]
    : selectedIds.filter((id) => id !== resource.id)
  emit('update:selectedIds', currentSelectedIds)
}

const getResourceCheckboxLabel = (resource: Resource) => {
  switch (resource.type) {
    case 'folder':
      return $gettext('Select folder')
    case 'space':
      return $gettext('Select space')
    default:
      return $gettext('Select file')
  }
}

const currentSortField = computed(() => {
  return sortFields.find((o) => o.name === sortBy && o.sortDir === sortDir) || sortFields[0]
})
const selectSorting = (field: SortField) => {
  emit('sort', { sortBy: field.name, sortDir: field.sortDir })
}

const resourceIconSize = computed<IconSize>(() => {
  const sizeMap: Record<number, string> = {
    1: 'xlarge',
    2: 'xlarge',
    3: 'xxlarge',
    4: 'xxlarge',
    5: 'xxxlarge',
    6: 'xxxlarge'
  }
  const size = unref(viewSizeCurrent)
  return (sizeMap[size] as IconSize) ?? ('xxlarge' as IconSize)
})
onBeforeUpdate(() => {
  tileRefs.value = {
    tiles: {},
    dropBtns: {}
  }
})

const setDropStyling = (resource: Resource, leaving: boolean, event: DragEvent) => {
  const hasFilePayload = (event.dataTransfer?.types || []).some((e) => e === 'Files')
  if (
    hasFilePayload ||
    (event.currentTarget as HTMLElement)?.contains(event.relatedTarget as HTMLElement) ||
    selectedIds.includes(resource.id) ||
    resource.type !== 'folder'
  ) {
    return
  }
  const el = unref(tileRefs).tiles[resource.id]
  if (leaving) {
    el.$el.classList.remove('oc-tiles-item-drop-highlight')
    return
  }
  el.$el.classList.add('oc-tiles-item-drop-highlight')
}
const dragSelection = computed(() => {
  return selectedIds.filter((id) => id !== unref(dragItem).id)
})
const setDragItem = async (item: Resource, event: DragEvent) => {
  dragItem.value = item
  await nextTick()
  unref(ghostElementRef).$el.ariaHidden = 'true'
  unref(ghostElementRef).$el.style.left = '-99999px'
  unref(ghostElementRef).$el.style.top = '-99999px'
  event.dataTransfer.setDragImage(unref(ghostElementRef).$el, 0, 0)
  event.dataTransfer.dropEffect = 'move'
  event.dataTransfer.effectAllowed = 'move'
}
const dragStart = async (resource: Resource, event: DragEvent) => {
  if (!isResourceSelected(resource)) {
    toggleSelection(resource)
  }
  await setDragItem(resource, event)
}

const fileDropped = (resource: Resource, event: DragEvent) => {
  const hasFilePayload = (event.dataTransfer.types || []).some((e) => e === 'Files')
  if (hasFilePayload) {
    return
  }
  dragItem.value = null
  setDropStyling(resource, true, event)
  emit('fileDropped', resource.id)
}

const viewWidth = ref(0)
const updateViewWidth = () => {
  const element = document.getElementById('tiles-view')
  const style = getComputedStyle(element)
  const paddingLeft = parseInt(style.getPropertyValue('padding-left'), 10) | 0
  const paddingRight = parseInt(style.getPropertyValue('padding-right'), 10) | 0
  viewWidth.value = element.clientWidth - paddingLeft - paddingRight
}
const gapSizePixels = computed(() => {
  return parseFloat(getComputedStyle(document.documentElement).fontSize)
})
const { calculateTileSizePixels } = useTileSize()
const maxTilesAll = computed<number[]>(() => {
  const viewSizes = [...Array(FolderViewModeConstants.tilesSizeMax).keys()].map((i) => i + 1)
  return [
    ...new Set<number>(
      viewSizes.map((viewSize) => {
        const pixels = calculateTileSizePixels(viewSize)
        return pixels ? Math.round(unref(viewWidth) / (pixels + unref(gapSizePixels))) : 0
      })
    )
  ]
})
const maxTilesCurrent = computed(() => {
  const maxTiles = unref(maxTilesAll)
  return maxTiles.length < unref(viewSizeCurrent)
    ? maxTiles[maxTiles.length - 1]
    : maxTiles[unref(viewSizeCurrent) - 1]
})
const ghostTilesCount = computed(() => {
  const remainder = unref(maxTilesCurrent) ? resources.length % unref(maxTilesCurrent) : 0
  if (!remainder) {
    return 0
  }
  return unref(maxTilesCurrent) - remainder
})

const tileSizePixels = computed(() => {
  return unref(viewWidth) / unref(maxTilesCurrent) - unref(gapSizePixels)
})

const areAllResourcesSelected = computed(() => {
  const allResourcesDisabled = unref(disabledResourceIds).length === resources.length
  const allSelected =
    unref(selectedResources).length === resources.length - unref(disabledResourceIds).length

  return !allResourcesDisabled && allSelected
})

watch(
  tileSizePixels,
  (px: number) => {
    document.documentElement.style.setProperty(`--oc-size-tiles-actual`, `${px}px`)
  },
  { immediate: true }
)
watch(maxTilesAll, (all) => {
  viewSizeMax.value = Math.max(all.length, 1)
})

onMounted(() => {
  window.addEventListener('resize', updateViewWidth)
  updateViewWidth()
})
onBeforeUnmount(() => {
  window.removeEventListener('resize', updateViewWidth)
})
</script>

<style lang="scss">
.oc-tiles {
  column-gap: 1rem;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(var(--oc-size-tiles-actual), 1fr));
  justify-content: flex-start;
  row-gap: 1rem;

  &-controls {
    border-bottom: 1px solid var(--oc-color-border);
  }

  &-item-drop-highlight {
    background-color: var(--oc-color-input-border) !important;
  }

  &-footer {
    color: var(--oc-color-text-muted);
    font-size: var(--oc-font-size-default);
    line-height: 1.4;
    padding: var(--oc-space-xsmall);
  }

  &-sort-select {
    min-width: var(--oc-size-width-xsmall);

    .v-select {
      flex: 1;
      margin-left: var(--oc-space-small);
    }
    .vs__dropdown-menu {
      min-width: var(--oc-size-width-small);
    }
  }
}

.ghost-tile {
  display: list-item;

  div {
    opacity: 0;
    box-shadow: none;
    height: 100%;
    display: flex;
    flex-flow: column;
    outline: 1px solid var(--oc-color-border);
  }
}
</style>
