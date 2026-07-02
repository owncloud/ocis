<template>
  <oc-table
    v-bind="$attrs"
    id="files-space-table"
    :class="[
      {
        condensed: viewMode === FolderViewModeConstants.name.condensedTable,
        'files-table': resourceType === 'file',
        'files-table-squashed': resourceType === 'file' && isSideBarOpen,

        'spaces-table': resourceType === 'space',
        'spaces-table-squashed': resourceType === 'space' && isSideBarOpen
      }
    ]"
    :data="resources"
    :fields="fields"
    :highlighted="selectedIds"
    :disabled="disabledResources"
    :sticky="isSticky"
    :header-position="headerPosition"
    :drag-drop="dragDrop"
    :hover="hover"
    :item-dom-selector="resourceDomSelector"
    :selection="selectedResources"
    :sort-by="sortBy"
    :sort-dir="sortDir"
    :lazy="lazy"
    :grouping-settings="groupingSettings"
    padding-x="medium"
    @highlight="fileClicked"
    @row-mounted="rowMounted"
    @contextmenu-clicked="showContextMenu"
    @item-dropped="fileDropped"
    @item-dragged="fileDragged"
    @drop-row-styling="dropRowStyling"
    @sort="sort"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <template v-if="!isLocationPicker && !isFilePicker" #selectHeader>
      <div class="resource-table-select-all">
        <oc-checkbox
          id="resource-table-select-all"
          v-oc-tooltip="{ content: selectAllCheckboxLabel, placement: 'bottom' }"
          size="large"
          :label="selectAllCheckboxLabel"
          :disabled="resources.length === disabledResources.length"
          :label-hidden="false"
          :label-classes="['oc-invisible-sr']"
          :model-value="areAllResourcesSelected"
          @click.stop="toggleSelectionAll"
        />
      </div>
    </template>
    <template v-if="!isLocationPicker && !isFilePicker" #select="{ item }">
      <oc-spinner
        v-if="isResourceInDeleteQueue(item.id)"
        class="resource-table-activity-indicator"
        size="medium"
      />

      <oc-checkbox
        v-else
        :id="`resource-table-select-${resourceDomSelector(item)}`"
        :label="getResourceCheckboxLabel(item)"
        :label-hidden="true"
        size="large"
        :disabled="isResourceDisabled(item)"
        :model-value="isResourceSelected(item)"
        :outline="isLatestSelectedItem(item)"
        @click.stop="toggleSelection(item.id)"
      />
    </template>
    <template #name="{ item }">
      <div
        class="resource-table-resource-wrapper"
        :class="[{ 'resource-table-resource-wrapper-limit-max-width': hasRenameAction(item) }]"
      >
        <slot name="image" :resource="item" />
        <resource-list-item
          :key="`${item.path}-${resourceDomSelector(item)}-${item.thumbnail}`"
          :resource="item"
          :path-prefix="getPathPrefix(item)"
          :is-path-displayed="arePathsDisplayed"
          :parent-folder-name="getParentFolderName(item)"
          :is-icon-displayed="!$slots['image']"
          :is-extension-displayed="areFileExtensionsShown"
          :is-resource-clickable="isResourceClickable(item)"
          :link="getResourceLink(item)"
          :parent-folder-link="getParentFolderLink(item)"
          :parent-folder-link-icon-additional-attributes="
            getParentFolderLinkIconAdditionalAttributes(item)
          "
          :class="{ 'resource-table-resource-cut': isResourceCut(item) }"
          @click="emitFileClick(item)"
        />
        <oc-button
          v-if="hasRenameAction(item)"
          :aria-label="getRenameButtonAriaLabel(item)"
          class="resource-table-edit-name"
          appearance="raw"
          @click="openRenameDialog(item)"
        >
          <oc-icon name="edit-2" fill-type="line" size="small" />
        </oc-button>
      </div>
      <slot name="additionalResourceContent" :resource="item" />
    </template>
    <template #syncEnabled="{ item }">
      <!-- @slot syncEnabled column -->
      <slot name="syncEnabled" :resource="item" />
    </template>
    <template #size="{ item }">
      <resource-size :size="item.size || Number.NaN" />
    </template>
    <template #tags="{ item }">
      <component
        :is="userContextReady ? 'router-link' : 'span'"
        v-for="tag in item.tags.slice(0, 2)"
        :key="tag"
        v-bind="getTagComponentAttrs(tag)"
        class="resource-table-tag-wrapper"
      >
        <oc-tag
          v-oc-tooltip="getTagToolTip(tag)"
          class="resource-table-tag oc-ml-xs"
          :rounded="true"
          size="small"
        >
          <oc-icon name="price-tag-3" size="small" />
          <span class="oc-text-truncate">{{ tag }}</span>
        </oc-tag>
      </component>
      <oc-tag
        v-if="item.tags.length > 2"
        size="small"
        class="resource-table-tag-more"
        @click="openTagsSidebar"
      >
        + {{ item.tags.length - 2 }}
      </oc-tag>
    </template>
    <template #manager="{ item }">
      <slot name="manager" :resource="item" />
    </template>
    <template #members="{ item }">
      <slot name="members" :resource="item" />
    </template>
    <template #totalQuota="{ item }">
      <slot name="totalQuota" :resource="item" />
    </template>
    <template #usedQuota="{ item }">
      <slot name="usedQuota" :resource="item" />
    </template>
    <template #remainingQuota="{ item }">
      <slot name="remainingQuota" :resource="item" />
    </template>
    <template #mdate="{ item }">
      <span
        v-oc-tooltip="formatDate(item.mdate)"
        tabindex="0"
        v-text="formatDateRelative(item.mdate)"
      />
    </template>
    <template #indicators="{ item }">
      <resource-status-indicators
        :space="space"
        :resource="item"
        :disable-handler="isResourceDisabled(item)"
      />
    </template>
    <template #status="{ item }">
      <oc-icon
        v-oc-tooltip="item.disabled ? $gettext('Disabled') : $gettext('Enabled')"
        :name="item.disabled ? 'stop-circle' : 'play-circle'"
        :accessible-label="`${item.disabled ? $pgettext(`used by screen reader to announce spaces's status`, 'Space is disabled') : $pgettext(`used by reader to announce spaces's status`, 'Space is enabled')}`"
        size="small"
        fill-type="line"
      />
    </template>
    <template #sdate="{ item }">
      <span
        v-oc-tooltip="formatDate(item.sdate)"
        tabindex="0"
        v-text="formatDateRelative(item.sdate)"
      />
    </template>
    <template #ddate="{ item }">
      <p
        v-oc-tooltip="formatDate(item.ddate)"
        tabindex="0"
        class="oc-m-rm"
        v-text="formatDateRelative(item.ddate)"
      />
    </template>
    <template #sharedBy="{ item }">
      <oc-button
        appearance="raw-inverse"
        variation="passive"
        class="resource-table-shared-by"
        @click="openSharingSidebar(item)"
      >
        <oc-avatars
          class="resource-table-people"
          :items="getSharedByAvatarItems(item)"
          :is-tooltip-displayed="true"
          :accessible-description="getSharedByAvatarDescription(item)"
        />
      </oc-button>
    </template>
    <template #sharedWith="{ item }">
      <oc-button
        appearance="raw-inverse"
        variation="passive"
        class="resource-table-shared-with"
        @click="openSharingSidebar(item)"
      >
        <oc-avatars
          class="resource-table-people"
          :items="getSharedWithAvatarItems(item)"
          :stacked="true"
          :max-displayed="3"
          :is-tooltip-displayed="true"
          :accessible-description="getSharedWithAvatarDescription(item)"
        />
      </oc-button>
    </template>
    <template #actions="{ item }">
      <div v-if="!isResourceDisabled(item)" class="resource-table-actions">
        <!-- @slot Add quick actions before the `context-menu / three dot` button in the actions column -->
        <slot name="quickActions" :resource="item" />
        <context-menu-quick-action
          ref="contextMenuButton"
          :item="item"
          :resource-dom-selector="resourceDomSelector"
          class="resource-table-btn-action-dropdown"
          @quick-action-clicked="showContextMenuOnBtnClick($event, item)"
        >
          <template #contextMenu>
            <slot name="contextMenu" :resource="item" />
          </template>
        </context-menu-quick-action>
      </div>
    </template>
    <template v-if="$slots.footer" #footer>
      <!-- @slot Footer of the files table -->
      <slot name="footer" />
    </template>
  </oc-table>
  <Teleport v-if="dragItem" to="body">
    <resource-ghost-element ref="ghostElement" :preview-items="[dragItem, ...dragSelection]" />
  </Teleport>
</template>

<script lang="ts" setup>
import { computed, unref, ref, ComputedRef, ComponentPublicInstance, nextTick } from 'vue'
import { useWindowSize } from '@vueuse/core'
import {
  IncomingShareResource,
  isPasswordProtectedFolderFileResource,
  isProjectSpaceResource,
  isSpaceResource,
  Resource,
  TrashResource
} from '@ownclouders/web-client'
import { extractDomSelector, SpaceResource } from '@ownclouders/web-client'
import { ShareTypes, isShareResource } from '@ownclouders/web-client'

import {
  SortDir,
  FolderViewModeConstants,
  useGetMatchingSpace,
  useFolderLink,
  useEmbedMode,
  useAuthStore,
  useCapabilityStore,
  useClipboardStore,
  useResourcesStore,
  useRouter,
  useCanBeOpenedWithSecureView,
  useFileActions,
  useIsTopBarSticky,
  embedModeFilePickMessageData,
  routeToContextQuery,
  useSpaceActionsRename
} from '../../composables'
import ResourceListItem from './ResourceListItem.vue'
import ResourceGhostElement from './ResourceGhostElement.vue'
import ResourceSize from './ResourceSize.vue'
import { EVENT_TROW_MOUNTED, EVENT_FILE_DROPPED, ImageDimension } from '../../constants'
import { eventBus } from '../../services'
import {
  ContextMenuBtnClickEventData,
  CreateTargetRouteOptions,
  displayPositionedDropdown,
  formatDateFromJSDate,
  formatRelativeDateFromJSDate
} from '../../helpers'
import { SideBarEventTopics } from '../../composables/sideBar'
import ContextMenuQuickAction from '../ContextActions/ContextMenuQuickAction.vue'

import { ClipboardActions } from '../../helpers/clipboardActions'
import { determineResourceTableSortFields } from '../../helpers/ui/resourceTable'
import { useFileActionsRename } from '../../composables/actions'
import { createLocationCommon } from '../../router'
import get from 'lodash-es/get'
import { storeToRefs } from 'pinia'
import { OcButton, OcTable } from '@ownclouders/design-system/components'
import { FieldType } from '@ownclouders/design-system/helpers'
import { OcSpinner } from '@ownclouders/design-system/components'
import ResourceStatusIndicators from './ResourceStatusIndicators.vue'
import { useGettext } from 'vue3-gettext'

const TAGS_MINIMUM_SCREEN_WIDTH = 850

interface GroupingSettings {
  groupingBy: string
  showGroupingOptions: boolean
  groupingFunctions: {
    [key: string]: (row: IncomingShareResource) => string | void
  }
  sortGroups: {
    [key: string]: (groups: { name: string }[]) => { name: string }[]
  }
}
/**
 * Resources to be displayed in the table.
 * Required fields:
 * - name: The name of the resource containing the file extension in case of a file
 * - path: The full path of the resource
 * - type: The type of the resource. Can be `file` or `folder`
 * Optional fields:
 * - thumbnail
 * - size: The size of the resource
 * - modificationDate: The date of the last modification of the resource
 * - shareDate: The date when the share was created
 * - deletionDate: The date when the resource has been deleted
 * - syncEnabled: The sync status of the share
 */
interface Props {
  resources: Resource[]
  resourceDomSelector?: (resource: Resource) => string
  arePathsDisplayed?: boolean
  selectedIds?: string[]
  hasActions?: boolean
  targetRouteCallback?: (arg: CreateTargetRouteOptions) => unknown
  areResourcesClickable?: boolean
  headerPosition?: number
  isSelectable?: boolean
  isSideBarOpen?: boolean
  dragDrop?: boolean
  viewMode?:
    | typeof FolderViewModeConstants.name.condensedTable
    | typeof FolderViewModeConstants.name.table
  hover?: boolean
  sortBy?: string
  sortDir?: SortDir
  fieldsDisplayed?: string[] | null
  space?: SpaceResource | null
  resourceType?: 'file' | 'space'
  lazy?: boolean
  groupingSettings?: GroupingSettings
}
interface Emits {
  (e: 'fileClick', data: { space: SpaceResource; resources: Resource[] }): void
  (e: 'sort', data: { sortBy: string; sortDir: SortDir }): void
  (
    e: 'rowMounted',
    resource: Resource,
    component: ComponentPublicInstance<unknown>,
    thumbnailDimension: ImageDimension
  ): void
  (e: typeof EVENT_FILE_DROPPED, data: string): void
  (e: 'update:selectedIds', selectedIds: string[]): void
  (e: 'update:modelValue', value: string[]): void
}
const {
  resources,
  resourceDomSelector = (resource: Resource) => extractDomSelector(resource.id),
  arePathsDisplayed = false,
  selectedIds = [],
  hasActions = true,
  targetRouteCallback = undefined,
  areResourcesClickable = true,
  headerPosition = 0,
  isSelectable = true,
  isSideBarOpen = false,
  dragDrop = false,
  viewMode = FolderViewModeConstants.defaultModeName,
  hover = true,
  sortBy = undefined,
  sortDir = undefined,
  fieldsDisplayed = null,
  space = null,
  resourceType = 'file',
  lazy = true,
  groupingSettings = null
} = defineProps<Props>()

const emit = defineEmits<Emits>()
const router = useRouter()
const capabilityStore = useCapabilityStore()
const { getMatchingSpace } = useGetMatchingSpace()
const { canBeOpenedWithSecureView } = useCanBeOpenedWithSecureView()
const folderLinkUtils = useFolderLink({
  space: ref(space),
  targetRouteCallback: computed(() => targetRouteCallback)
})
const {
  getPathPrefix,
  getParentFolderName,
  getParentFolderLink,
  getParentFolderLinkIconAdditionalAttributes
} = folderLinkUtils
const { isSticky } = useIsTopBarSticky()
const {
  isLocationPicker,
  isFilePicker,
  postMessage,
  isEnabled: isEmbedModeEnabled,
  fileTypes: embedModeFileTypes
} = useEmbedMode()
const { getDefaultAction } = useFileActions()
const language = useGettext()
const { $pgettext, $gettext, $ngettext } = language

const clipboardStore = useClipboardStore()
const { resources: clipboardResources, action: clipboardAction } = storeToRefs(clipboardStore)

const authStore = useAuthStore()
const { userContextReady } = storeToRefs(authStore)

const resourcesStore = useResourcesStore()
const { areFileExtensionsShown, latestSelectedId, deleteQueue } = storeToRefs(resourcesStore)

const dragItem = ref<Resource>()
const ghostElement = ref()
const contextMenuButton = ref<ComponentPublicInstance<typeof OcButton>>()

const { width } = useWindowSize()
const constants = ref({
  ImageDimension,
  EVENT_TROW_MOUNTED
})
const hasTags = computed(
  () => capabilityStore.filesTags && width.value >= TAGS_MINIMUM_SCREEN_WIDTH
)

const { actions: renameActions } = useFileActionsRename()
const { actions: renameActionsSpace } = useSpaceActionsRename()
const renameHandler = computed(() => unref(renameActions)[0].handler)
const renameHandlerSpace = computed(() => unref(renameActionsSpace)[0].handler)

const getTagToolTip = (text: string) => (text.length > 7 ? text : '')

const isResourceDisabled = (resource: Resource) => {
  if (unref(isEmbedModeEnabled) && unref(embedModeFileTypes)?.length) {
    return (
      !unref(embedModeFileTypes).includes(resource.extension) &&
      !unref(embedModeFileTypes).includes(resource.mimeType) &&
      !resource.isFolder
    )
  }
  return resource.processing === true || isResourceInDeleteQueue(resource.id)
}

const disabledResources: ComputedRef<Array<Resource['id']>> = computed(() => {
  return (
    resources
      ?.filter((resource) => isResourceDisabled(resource) === true)
      ?.map((resource) => resource.id) || []
  )
})

const isResourceClickable = (resource: Resource) => {
  if (!areResourcesClickable) {
    return false
  }

  if (isProjectSpaceResource(resource) && resource.disabled) {
    return false
  }

  if (!resource.isFolder && !isPasswordProtectedFolderFileResource(resource.name)) {
    if (!resource.canDownload() && !canBeOpenedWithSecureView(resource)) {
      return false
    }

    if (unref(isEmbedModeEnabled) && !unref(isFilePicker)) {
      return false
    }
  }

  return !unref(disabledResources).includes(resource.id)
}

const emitSelect = (selectedIds: string[]) => {
  eventBus.publish('app.files.list.clicked')
  emit('update:selectedIds', selectedIds)
}

const toggleSelection = (resourceId: string) => {
  resourcesStore.toggleSelection(resourceId)
  emitSelect(resourcesStore.selectedIds)
}

const getResourceLink = (resource: Resource) => {
  if (resource.isFolder) {
    return folderLinkUtils.getFolderLink(resource)
  }

  let currentSpace = space
  if (!currentSpace) {
    currentSpace = getMatchingSpace(resource)
  }

  const action = getDefaultAction({ resources: [resource], space: currentSpace })

  if (!action?.route) {
    return
  }

  return action.route({ space: currentSpace, resources: [resource] })
}

const isResourceInDeleteQueue = (id: string): boolean => {
  return unref(deleteQueue).includes(id)
}

const getRenameButtonAriaLabel = (resource: Resource): string => {
  if (isSpaceResource(resource)) {
    return $pgettext(
      'The label of the rename button in the resource table for spaces',
      'Rename space'
    )
  } else if (resource.isFolder) {
    return $pgettext(
      'The label of the rename button in the resource table for folders',
      'Rename folder'
    )
  }

  return $pgettext('The label of the rename button in the resource table for files', 'Rename file')
}
const fields = computed(() => {
  if (resources.length === 0) {
    return []
  }
  const firstResource = resources[0]
  const fields: FieldType[] = []
  if (isSelectable) {
    fields.push({
      name: 'select',
      title: '',
      type: 'slot',
      headerType: 'slot',
      width: 'shrink'
    })
  }

  const sortFields = determineResourceTableSortFields(firstResource)
  fields.push(
    ...(
      [
        {
          name: 'name',
          title: $gettext('Name'),
          type: 'slot',
          width: 'expand',
          wrap: 'truncate'
        },

        {
          name: 'manager',
          prop: 'members',
          title: $gettext('Manager'),
          type: 'slot'
        },
        {
          name: 'members',
          title: $gettext('Members'),
          prop: 'members',
          type: 'slot'
        },
        {
          name: 'totalQuota',
          prop: 'spaceQuota.total',
          title: $gettext('Total quota'),
          type: 'slot',
          sortable: true
        },
        {
          name: 'usedQuota',
          prop: 'spaceQuota.used',
          title: $gettext('Used quota'),
          type: 'slot',
          sortable: true
        },
        {
          name: 'remainingQuota',
          prop: 'spaceQuota.remaining',
          title: $gettext('Remaining quota'),
          type: 'slot',
          sortable: true
        },
        {
          name: 'indicators',
          title: $gettext('Status'),
          type: 'slot',
          alignH: 'right',
          wrap: 'nowrap',
          width: 'shrink'
        },
        {
          name: 'size',
          title: $gettext('Size'),
          type: 'slot',
          alignH: 'right',
          wrap: 'nowrap',
          width: 'shrink'
        },
        {
          name: 'syncEnabled',
          title: $gettext('Info'),
          type: 'slot',
          alignH: 'right',
          wrap: 'nowrap',
          width: 'shrink'
        },
        {
          name: 'status',
          prop: 'disabled',
          title: $gettext('Status'),
          type: 'slot',
          alignH: 'right',
          wrap: 'nowrap',
          width: 'shrink'
        },
        {
          name: 'tags',
          title: $gettext('Tags'),
          type: 'slot',
          alignH: 'right',
          wrap: 'nowrap',
          width: 'shrink'
        },
        {
          name: 'sharedBy',
          title: $gettext('Shared by'),
          type: 'slot',
          alignH: 'right',
          wrap: 'nowrap',
          width: 'shrink'
        },
        {
          name: 'sharedWith',
          title: $gettext('Shared with'),
          type: 'slot',
          alignH: 'right',
          wrap: 'nowrap',
          width: 'shrink'
        },
        {
          name: 'mdate',
          title: $gettext('Modified'),
          type: 'slot',
          alignH: 'right',
          wrap: 'nowrap',
          width: 'shrink',
          accessibleLabelCallback: (item) =>
            formatDateRelative((item as Resource).mdate) +
            ' (' +
            formatDate((item as Resource).mdate) +
            ')'
        },
        {
          name: 'sdate',
          title: $gettext('Shared on'),
          type: 'slot',
          alignH: 'right',
          wrap: 'nowrap',
          width: 'shrink',
          accessibleLabelCallback: (item) =>
            formatDateRelative((item as IncomingShareResource).sdate) +
            ' (' +
            formatDate((item as IncomingShareResource).sdate) +
            ')'
        },
        {
          name: 'ddate',
          title: $gettext('Deleted'),
          type: 'slot',
          alignH: 'right',
          wrap: 'nowrap',
          width: 'shrink',
          accessibleLabelCallback: (item) =>
            formatDateRelative((item as TrashResource).ddate) +
            ' (' +
            formatDate((item as TrashResource).ddate) +
            ')'
        }
      ] as FieldType[]
    )
      .filter((field) => {
        if (field.name === 'tags' && !unref(hasTags)) {
          return false
        }

        if (field.name === 'indicators') {
          return true
        }

        let hasField: boolean
        if (field.prop) {
          hasField = get(firstResource, field.prop) !== undefined
        } else {
          hasField = Object.prototype.hasOwnProperty.call(firstResource, field.name)
        }
        if (!fieldsDisplayed) {
          return hasField
        }

        return hasField && fieldsDisplayed.includes(field.name)
      })
      .map((field) => {
        const sortField = sortFields.find((f) => f.name === field.name)
        if (sortField) {
          Object.assign(field, {
            sortable: sortField.sortable,
            sortDir: sortField.sortDir
          })
        }
        return field
      })
  )
  if (hasActions) {
    fields.push({
      name: 'actions',
      title: $gettext('Actions'),
      type: 'slot',
      alignH: 'right',
      wrap: 'nowrap',
      width: 'shrink'
    })
  }

  return fields
})
const areAllResourcesSelected = computed(() => {
  const allResourcesDisabled = unref(disabledResources).length === resources.length
  const allSelected =
    unref(selectedResources).length === resources.length - unref(disabledResources).length

  return !allResourcesDisabled && allSelected
})
const selectAllCheckboxLabel = computed(() => {
  return unref(areAllResourcesSelected) ? $gettext('Clear selection') : $gettext('Select all')
})
const selectedResources = computed(() => {
  return resources.filter((resource) => selectedIds.includes(resource.id))
})
const dragSelection = computed(() => {
  const selection = [...unref(selectedResources)]
  selection.splice(
    selection.findIndex((i) => i.id === unref(dragItem).id),
    1
  )
  return selection
})
function isResourceSelected(item: Resource) {
  return selectedIds.includes(item.id)
}
function isResourceCut(resource: Resource) {
  if (unref(clipboardAction) !== ClipboardActions.Cut) {
    return false
  }
  return unref(clipboardResources).some((r) => r.id === resource.id)
}
function getTagLink(tag: string) {
  const currentTerm = unref(router.currentRoute).query?.term
  return createLocationCommon('files-common-search', {
    query: { provider: 'files.sdk', q_tags: tag, ...(currentTerm && { term: currentTerm }) }
  })
}
function getTagComponentAttrs(tag: string) {
  if (!userContextReady) {
    return {}
  }

  return {
    to: getTagLink(tag)
  }
}
function isLatestSelectedItem(item: Resource) {
  return item.id === unref(latestSelectedId)
}
function hasRenameAction(item: Resource) {
  if (isProjectSpaceResource(item)) {
    return unref(renameActionsSpace).filter((menuItem) => menuItem.isVisible({ resources: [item] }))
      .length
  }

  return unref(renameActions).filter((menuItem) => menuItem.isVisible({ space, resources: [item] }))
    .length
}
function openRenameDialog(item: Resource) {
  if (isProjectSpaceResource(item)) {
    return unref(renameHandlerSpace)({
      resources: [item]
    })
  }
  unref(renameHandler)({
    space: getMatchingSpace(item),
    resources: [item]
  })
}
function openTagsSidebar() {
  eventBus.publish(SideBarEventTopics.open)
}
function openSharingSidebar(file: Resource) {
  let panelToOpen
  if (file.type === 'space') {
    panelToOpen = 'space-share'
  } else {
    panelToOpen = 'sharing'
  }
  eventBus.publish(SideBarEventTopics.openWithPanel, panelToOpen)
}
async function fileDragged(file: Resource, event: DragEvent) {
  if (!unref(dragDrop)) {
    return
  }

  await setDragItem(file, event)

  addSelectedResource(file)
}
function fileDropped(selector: HTMLElement, event: DragEvent) {
  if (!unref(dragDrop)) {
    return
  }
  const hasFilePayload = (event.dataTransfer.types || []).some((e) => e === 'Files')
  if (hasFilePayload) {
    return
  }
  dragItem.value = null
  const dropTarget = event.target as HTMLElement
  const dropTargetTr = dropTarget.closest('tr')
  const dropItemId = dropTargetTr.dataset.itemId
  dropRowStyling(selector, true, event)

  emit(EVENT_FILE_DROPPED, dropItemId)
}
async function setDragItem(item: Resource, event: DragEvent) {
  dragItem.value = item
  await nextTick()
  unref(ghostElement).$el.ariaHidden = 'true'
  unref(ghostElement).$el.style.left = '-99999px'
  unref(ghostElement).$el.style.top = '-99999px'
  event.dataTransfer.setDragImage(unref(ghostElement).$el, 0, 0)
  event.dataTransfer.dropEffect = 'move'
  event.dataTransfer.effectAllowed = 'move'
}
function dropRowStyling(selector: HTMLElement, leaving: boolean, event: DragEvent) {
  const hasFilePayload = (event.dataTransfer?.types || []).some((e) => e === 'Files')
  if (hasFilePayload) {
    return
  }
  if ((event.currentTarget as HTMLElement)?.contains(event.relatedTarget as HTMLElement)) {
    return
  }

  const classList = document.getElementsByClassName(`oc-tbody-tr-${selector}`)[0].classList
  const className = 'highlightedDropTarget'
  leaving ? classList.remove(className) : classList.add(className)
}
function sort(opts: { sortBy: string; sortDir: SortDir }) {
  emit('sort', opts)
}
function addSelectedResource(file: Resource) {
  const isSelected = isResourceSelected(file)
  if (isSelected) {
    return
  }
  toggleSelection(file.id)
}
function showContextMenuOnBtnClick(data: ContextMenuBtnClickEventData, item: Resource) {
  if (unref(isResourceDisabled)(item)) {
    return false
  }

  const { dropdown, event } = data
  if (dropdown?.tippy === undefined) {
    return
  }
  if (!isResourceSelected(item)) {
    emitSelect([item.id])
  }
  displayPositionedDropdown(dropdown.tippy, event, unref(contextMenuButton))
}
function showContextMenu(row: ComponentPublicInstance<unknown>, event: MouseEvent, item: Resource) {
  event.preventDefault()

  if (isResourceDisabled(item)) {
    return false
  }

  const instance = row.$el.getElementsByClassName('resource-table-btn-action-dropdown')[0]
  if (instance === undefined) {
    return
  }
  if (!isResourceSelected(item)) {
    emitSelect([item.id])
  }
  displayPositionedDropdown(instance._tippy, event, unref(contextMenuButton))
}
function rowMounted(resource: Resource, component: ComponentPublicInstance<unknown>) {
  /**
   * Triggered whenever a row is mounted
   * @property {object} resource The resource which was mounted as table row
   * @property {object} component The table row component
   */
  emit('rowMounted', resource, component, unref(constants).ImageDimension.Thumbnail)
}
function fileClicked(data: [Resource, MouseEvent, boolean]) {
  /**
   * Triggered when the file row is clicked
   * @property {object} resource The resource for which the event is triggered
   */
  const resource = data[0]

  if (isResourceDisabled(resource)) {
    return
  }

  if (unref(isEmbedModeEnabled) && unref(isFilePicker) && !resource.isFolder) {
    return postMessage<embedModeFilePickMessageData>('owncloud-embed:file-pick', {
      resource: JSON.parse(JSON.stringify(resource)),
      locationQuery: JSON.parse(JSON.stringify(routeToContextQuery(unref(router.currentRoute))))
    })
  }

  const eventData = data[1]
  const skipTargetSelection = data[2] ?? false

  const isCheckboxClicked = (eventData?.target as HTMLElement).getAttribute('type') === 'checkbox'
  const contextActionClicked =
    (eventData?.target as HTMLElement)?.closest('div')?.id === 'oc-files-context-menu'
  if (contextActionClicked) {
    return
  }
  if (eventData && eventData.metaKey) {
    return eventBus.publish('app.files.list.clicked.meta', resource)
  }
  if (eventData && eventData.shiftKey) {
    return eventBus.publish('app.files.list.clicked.shift', { resource, skipTargetSelection })
  }
  if (isCheckboxClicked) {
    return
  }

  if (isResourceSelected(resource)) {
    return
  }

  resourcesStore.setLastSelectedId(resource.id)

  return emitSelect([resource.id])
}
function formatDate(date: string) {
  return formatDateFromJSDate(new Date(date), language.current)
}
function formatDateRelative(date: string) {
  return formatRelativeDateFromJSDate(new Date(date), language.current)
}
function toggleSelectionAll() {
  if (unref(areAllResourcesSelected)) {
    return emitSelect([])
  }
  emitSelect(
    resources
      .filter((resource) => !unref(disabledResources).includes(resource.id))
      .map((resource) => resource.id)
  )
}
function emitFileClick(resource: Resource) {
  const space = getMatchingSpace(resource)

  /**
   * Triggered when a default action is triggered on a file
   * @property {object} resource resource for which the event is triggered
   */
  emit('fileClick', { space, resources: [resource] })
}
function getResourceCheckboxLabel(resource: Resource) {
  if (resource.type === 'folder') {
    return $gettext('Select folder')
  }
  return $gettext('Select file')
}
function getSharedWithAvatarDescription(resource: Resource) {
  if (!isShareResource(resource)) {
    return
  }
  const resourceType = resource.type === 'folder' ? $gettext('folder') : $gettext('file')

  const shareCount = resource.sharedWith.filter(({ shareType }) =>
    ShareTypes.authenticated.includes(ShareTypes.getByValue(shareType))
  ).length

  if (!shareCount) {
    return ''
  }

  return $ngettext(
    'This %{ resourceType } is shared via %{ shareCount } invite',
    'This %{ resourceType } is shared via %{ shareCount } invites',
    shareCount,
    {
      resourceType,
      shareCount: shareCount.toString()
    }
  )
}
function getSharedByAvatarDescription(resource: Resource) {
  if (!isShareResource(resource)) {
    return ''
  }

  const resourceType = resource.type === 'folder' ? $gettext('folder') : $gettext('file')
  return $gettext('This %{ resourceType } is shared by %{ user }', {
    resourceType,
    user: resource.sharedBy.map(({ displayName }) => displayName).join(', ')
  })
}
function getSharedByAvatarItems(resource: Resource) {
  if (!isShareResource(resource)) {
    return []
  }

  return resource.sharedBy.map((s) => ({
    displayName: s.displayName,
    name: s.displayName,
    shareType: ShareTypes.user.value,
    username: s.id
  }))
}
function getSharedWithAvatarItems(resource: Resource) {
  if (!isShareResource(resource)) {
    return []
  }

  return resource.sharedWith
    .filter(({ shareType }) => ShareTypes.authenticated.includes(ShareTypes.getByValue(shareType)))
    .map((s) => ({
      displayName: s.displayName,
      name: s.displayName,
      shareType: s.shareType,
      username: s.id
    }))
}
</script>
<style lang="scss">
.oc-table.condensed > tbody > tr {
  height: 0 !important;
}

.remove-legend-styling {
  margin: 0px;
  padding: 0px;
  border: 0;
}

.resource-table {
  &-resource-cut {
    opacity: 0.7;
  }

  &-resource-wrapper {
    display: flex;
    align-items: center;

    &-limit-max-width {
      max-width: calc(100% - var(--oc-space-medium));
    }

    &:hover > .resource-table-edit-name {
      svg {
        fill: var(--oc-color-text-default);
      }
    }
  }

  &-tag {
    max-width: 80px;
  }

  &-tag-more {
    cursor: pointer;
    border: 0 !important;
    vertical-align: text-bottom;
  }

  &-edit-name,
  &-activity-indicator {
    display: inline-flex;
    margin-left: var(--oc-space-xsmall);

    svg {
      fill: var(--oc-color-text-muted);
    }
  }

  &-people {
    margin-right: -5px;
  }

  &-actions {
    align-items: center;
    display: flex;
    flex-flow: row nowrap;
    gap: var(--oc-space-xsmall);
    justify-content: flex-end;
  }

  &-select-all {
    align-items: center;
    display: flex;
    justify-content: center;
  }
}

.spaces-table {
  .oc-table-header-cell-mdate,
  .oc-table-data-cell-mdate,
  .oc-table-header-cell-manager,
  .oc-table-data-cell-manager,
  .oc-table-header-cell-remainingQuota,
  .oc-table-data-cell-remainingQuota,
  .oc-table-header-cell-members,
  .oc-table-data-cell-members,
  .oc-table-header-cell-status,
  .oc-table-data-cell-status {
    display: none;

    @media only screen and (min-width: 960px) {
      display: table-cell;
    }
  }

  .oc-table-header-cell-totalQuota,
  .oc-table-data-cell-totalQuota,
  .oc-table-header-cell-usedQuota,
  .oc-table-data-cell-usedQuota {
    display: none;

    @media only screen and (min-width: 1200px) {
      display: table-cell;
    }
  }

  &-squashed {
    /**
     * squashed = right sidebar is open.
     * same media queries as above but +440px width of the right sidebar
     * (because the right sidebar steals 440px from the file list)
     */
    .oc-table-header-cell-status,
    .oc-table-data-cell-status,
    .oc-table-header-cell-manager,
    .oc-table-data-cell-manager,
    .oc-table-header-cell-totalQuota,
    .oc-table-data-cell-totalQuota,
    .oc-table-header-cell-usedQuota,
    .oc-table-data-cell-usedQuota,
    .oc-table-header-cell-members,
    .oc-table-data-cell-members {
      display: none;

      @media only screen and (min-width: 1400px) {
        display: table-cell;
      }
    }

    .oc-table-header-cell-mdate,
    .oc-table-data-cell-mdate,
    .oc-table-header-cell-remainingQuota,
    .oc-table-data-cell-remainingQuota,
    .oc-table-header-cell-mdate,
    .oc-table-data-cell-mdate {
      display: none;

      @media only screen and (min-width: 1600px) {
        display: table-cell;
      }
    }
  }
}

// Hide files table columns
.files-table {
  .oc-table-header-cell-size,
  .oc-table-data-cell-size,
  .oc-table-header-cell-sharedWith,
  .oc-table-data-cell-sharedWith,
  .oc-table-header-cell-sharedBy,
  .oc-table-data-cell-sharedBy,
  .oc-table-header-cell-status,
  .oc-table-data-cell-status {
    display: none;

    @media only screen and (min-width: 640px) {
      display: table-cell;
    }
  }

  .oc-table-header-cell-mdate,
  .oc-table-data-cell-mdate,
  .oc-table-header-cell-sdate,
  .oc-table-data-cell-sdate,
  .oc-table-header-cell-ddate,
  .oc-table-data-cell-ddate {
    display: none;

    @media only screen and (min-width: 960px) {
      display: table-cell;
    }
  }

  .oc-table-header-cell-sharedBy,
  .oc-table-data-cell-sharedBy,
  .oc-table-header-cell-tags,
  .oc-table-data-cell-tags,
  .oc-table-header-cell-indicators,
  .oc-table-data-cell-indicators {
    display: none;

    @media only screen and (min-width: 1200px) {
      display: table-cell;
    }
  }

  &-squashed {
    /**
     * squashed = right sidebar is open.
     * same media queries as above but +440px width of the right sidebar
     * (because the right sidebar steals 440px from the file list)
     */
    .oc-table-header-cell-size,
    .oc-table-data-cell-size,
    .oc-table-header-cell-sharedWith,
    .oc-table-data-cell-sharedWith,
    .oc-table-header-cell-sharedBy,
    .oc-table-data-cell-sharedBy,
    .oc-table-header-cell-status,
    .oc-table-data-cell-status {
      display: none;

      @media only screen and (min-width: 1080px) {
        display: table-cell;
      }
    }

    .oc-table-header-cell-mdate,
    .oc-table-data-cell-mdate,
    .oc-table-header-cell-sdate,
    .oc-table-data-cell-sdate,
    .oc-table-header-cell-ddate,
    .oc-table-data-cell-ddate {
      display: none;

      @media only screen and (min-width: 1400px) {
        display: table-cell;
      }
    }

    .oc-table-header-cell-sharedBy,
    .oc-table-data-cell-sharedBy,
    .oc-table-header-cell-tags,
    .oc-table-data-cell-tags,
    .oc-table-header-cell-indicators,
    .oc-table-data-cell-indicators {
      display: none;

      @media only screen and (min-width: 1640px) {
        display: table-cell;
      }
    }
  }
}

// shared with me: on tablets hide shared with column and display sharedBy column instead
#files-shared-with-me-view .files-table .oc-table-header-cell-sharedBy,
#files-shared-with-me-view .files-table .oc-table-data-cell-sharedBy,
#files-shared-with-me-view .files-table .oc-table-header-cell-syncEnabled,
#files-shared-with-me-view .files-table .oc-table-data-cell-syncEnabled {
  @media only screen and (min-width: 640px) {
    display: table-cell;
  }
}

#files-shared-with-me-view .files-table .oc-table-header-cell-sharedWith,
#files-shared-with-me-view .files-table .oc-table-data-cell-sharedWith,
#files-shared-with-me-view .files-table .oc-table-header-cell-syncEnabled,
#files-shared-with-me-view .files-table .oc-table-data-cell-syncEnabled {
  @media only screen and (max-width: 1199px) {
    display: none;
  }
}

// Show tooltip on status indicators without handler
.oc-table-data-cell-indicators {
  span.oc-status-indicators-indicator {
    pointer-events: all;
  }
}
</style>
