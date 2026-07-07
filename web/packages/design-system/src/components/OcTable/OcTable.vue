<template>
  <table v-bind="extractTableProps()" class="has-item-context-menu">
    <caption v-if="caption || isTableSortable" :class="{ 'oc-invisible-sr': !captionVisible }">
      {{
        caption
      }}
      <span
        v-if="isTableSortable"
        class="oc-invisible-sr"
        v-text="
          $pgettext(
            'Table component caption sorting explanation available only to screen readers.',
            'Column headers with buttons are sortable.'
          )
        "
      />
    </caption>
    <oc-thead v-if="hasHeader">
      <oc-tr class="oc-table-header-row">
        <oc-th
          v-for="(field, index) in fields"
          :key="`oc-thead-${field.name}`"
          v-bind="extractThProps(field, index)"
          :aria-sort="getAriaSortValue(field.name)"
        >
          <oc-button
            v-if="field.sortable"
            appearance="raw"
            class="oc-button-sort oc-width-1-1"
            @click="handleSort(field)"
          >
            <span v-if="field.headerType === 'slot'" class="oc-table-thead-content">
              <slot :name="field.name + 'Header'" />
            </span>
            <span
              v-else
              class="oc-table-thead-content header-text"
              v-text="extractFieldTitle(field)"
            />
            <oc-icon
              :name="sortDir === 'asc' ? 'arrow-down' : 'arrow-up'"
              fill-type="line"
              :class="{ 'oc-invisible-sr': sortBy !== field.name }"
              size="small"
              variation="passive"
            />
          </oc-button>
          <div v-else>
            <span v-if="field.headerType === 'slot'" class="oc-table-thead-content">
              <slot :name="field.name + 'Header'" />
            </span>
            <span
              v-else
              class="oc-table-thead-content header-text"
              :class="{ 'oc-invisible-sr': !extractFieldTitle(field) }"
              :aria-hidden="extractFieldTitle(field) ? 'false' : 'true'"
              v-text="
                hasIconsColumn
                  ? $pgettext('Sets a hidden table header text for screen readers', 'Icons')
                  : extractFieldTitle(field)
              "
            />
          </div>
        </oc-th>
      </oc-tr>
    </oc-thead>
    <oc-tbody class="has-item-context-menu">
      <oc-tr
        v-for="(item, trIndex) in data"
        :key="`oc-tbody-tr-${domElementSelector(item) || trIndex}`"
        :ref="`row-${trIndex}`"
        v-bind="extractTbodyTrProps(item, trIndex)"
        :data-item-id="item[idKey as keyof Item]"
        :draggable="dragDrop"
        @click="$emit(constants.EVENT_TROW_CLICKED, [item, $event])"
        @contextmenu="
          $emit(
            constants.EVENT_TROW_CONTEXTMENU,
            ($refs[`row-${trIndex}`] as HTMLElement[])[0],
            $event,
            item
          )
        "
        @vue:mounted="
          $emit(constants.EVENT_TROW_MOUNTED, item, ($refs[`row-${trIndex}`] as HTMLElement[])[0])
        "
        @dragstart="dragStart(item, $event)"
        @drop="dropRowEvent(domElementSelector(item), $event)"
        @dragenter.prevent="dropRowStyling(domElementSelector(item), false, $event)"
        @dragleave.prevent="dropRowStyling(domElementSelector(item), true, $event)"
        @mouseleave="dropRowStyling(domElementSelector(item), true, $event)"
        @blur="dropRowStyling(domElementSelector(item), true, $event)"
        @dragover="dragOver($event)"
        @item-visible="$emit('itemVisible', item)"
      >
        <oc-td
          v-for="(field, tdIndex) in fields"
          :key="'oc-tbody-td-' + cellKey(field, tdIndex, item)"
          v-bind="extractTdProps(field, tdIndex, item)"
        >
          <slot v-if="isFieldTypeSlot(field)" :name="field.name" :item="item" />
          <template v-else-if="isFieldTypeCallback(field)">
            {{ field.callback(item[field.name as keyof Item]) }}
          </template>
          <template v-else>
            {{ item[field.name as keyof Item] }}
          </template>
        </oc-td>
      </oc-tr>
    </oc-tbody>
    <tfoot v-if="$slots.footer" class="oc-table-footer">
      <tr class="oc-table-footer-row">
        <td :colspan="fullColspan" class="oc-table-footer-cell">
          <!-- @slot Footer of the table -->
          <slot name="footer" />
        </td>
      </tr>
    </tfoot>
  </table>
</template>
<script lang="ts" setup>
import OcThead from '../OcTableHeader/OcTableHeader.vue'
import OcTbody from '../OcTableBody/OcTableBody.vue'
import OcTr from '../OcTableRow/OcTableRow.vue'
import OcTd from '../OcTableCellData/OcTableCellData.vue'
import OcTh from '../OcTableCellHead/OcTableCellHead.vue'
import OcButton from '../OcButton/OcButton.vue'
import { getSizeClass, Item, FieldType } from '../../helpers'
import { computed, unref } from 'vue'
import { useGettext } from 'vue3-gettext'

/**
 * @component OcTable
 * @description A flexible table component that supports various customization options including
 * sorting, row highlighting, custom cell rendering, sticky headers, and drag & drop functionality.
 *
 * @prop {Array<Object>} items - Array of data items to display in the table rows.
 * @prop {String} idKey - Key to use as unique identifier for items, defaults to 'id'.
 * @prop {Array<Object>} fields - Column configuration for the table.
 *        Each field object can include:
 *        - name {String} - Key of the data item to display (required)
 *        - title {String} - Column header text (optional, defaults to name)
 *        - headerType {String} - Type of header ('slot' or default)
 *        - type {String} - Cell content type ('slot', 'callback', or default)
 *        - callback {Function} - Function to render content if type is 'callback'
 *        - alignH {String} - Horizontal alignment ('left', 'center', 'right')
 *        - alignV {String} - Vertical alignment ('top', 'middle', 'bottom')
 *        - width {String} - Column width ('auto', 'shrink', 'expand')
 *        - wrap {String} - Text wrapping ('truncate', 'overflow', 'nowrap', 'break')
 *        - thClass {String} - Additional classes for header cells
 *        - tdClass {String} - Additional classes for data cells
 *        - sortable {Boolean} - Whether column is sortable
 *        - sortDir {String} - Default sort direction if sortable
 *        - accessibleLabelCallback {Function} - Function to generate accessible labels
 * @prop {Boolean} hasHeader - Whether to display the table header, defaults to true.
 * @prop {Boolean} sticky - Whether the header should be sticky, defaults to false.
 * @prop {Boolean} hover - Whether to highlight rows on hover, defaults to false.
 * @prop {String|Array} highlighted - IDs of rows to highlight.
 * @prop {Array<String|Number>} disabled - IDs of rows to disable.
 * @prop {Number} headerPosition - Top position of sticky header in pixels, defaults to 0.
 * @prop {String} paddingX - Horizontal padding size ('xsmall', 'small', 'medium', 'large', 'xlarge'), defaults to 'small'.
 * @prop {Boolean} dragDrop - Enable drag and drop functionality, defaults to false.
 * @prop {Array} selection - Array of pre-selected items.
 * @prop {Boolean} lazy - Whether table content should be loaded lazily, defaults to false.
 * @prop {String} sortDir - Current sort direction ('asc' or 'desc').
 * @prop {String} sortBy - Current sort column name.
 * @prop {Object} groupingSettings - Grouping configuration (CERN-specific).
 *
 * @event item-dropped - Emitted when an item is dropped during drag and drop.
 * @event item-dragged - Emitted when an item starts being dragged.
 * @event thead-clicked - Emitted when a table header is clicked.
 * @event trow-clicked - Emitted when a table row is clicked.
 * @event trow-mounted - Emitted when a table row is mounted.
 * @event trow-contextmenu - Emitted when right-click is performed on a row.
 * @event sort - Emitted when table sorting is requested with sortBy and sortDir.
 * @event dropRowStyling - Emitted to customize styling during drag operations.
 * @event itemVisible - Emitted when an item becomes visible (with lazy loading).
 *
 * @example
 * <template>
 *   <oc-table
 *     :items="users"
 *     :fields="fields"
 *     :has-header="true"
 *     :hover="true"
 *     :sticky="true"
 *     :sort-by="sortBy"
 *     :sort-dir="sortDir"
 *     @sort="handleSort"
 *     @trow-clicked="handleRowClick"
 *   >
 *     <template #actions="{ item }">
 *       <oc-button @click="editUser(item)">Edit</oc-button>
 *     </template>
 *   </oc-table>
 * </template>
 *
 */

import {
  EVENT_THEAD_CLICKED,
  EVENT_TROW_CLICKED,
  EVENT_TROW_MOUNTED,
  EVENT_TROW_CONTEXTMENU,
  EVENT_ITEM_DROPPED,
  EVENT_ITEM_DRAGGED
} from '../../helpers/constants'

const SORT_DIRECTION_ASC = 'asc' as const
const SORT_DIRECTION_DESC = 'desc' as const

interface Props {
  data: Item[]
  idKey?: string
  itemDomSelector?: (item: Item) => string
  fields: FieldType[]
  hasHeader?: boolean
  sticky?: boolean
  hover?: boolean
  highlighted?: string | string[]
  disabled?: Array<string | number>
  headerPosition?: number
  paddingX?: 'xsmall' | 'small' | 'medium' | 'large' | 'xlarge'
  dragDrop?: boolean
  selection?: Item[]
  lazy?: boolean
  sortDir?: 'asc' | 'desc'
  sortBy?: string
  caption?: string
  captionVisible?: boolean
  hasIconsColumn?: boolean
}

interface Emits {
  (e: 'itemDropped', selector: string, event: DragEvent): void
  (e: 'itemDragged', item: Item, event: DragEvent): void
  (e: 'rowMounted', item: Item, element: HTMLElement): void
  (e: 'theadClicked', event: MouseEvent): void
  (e: 'highlight', args: [Item, MouseEvent]): void
  (e: 'contextmenuClicked', element: HTMLElement, event: MouseEvent, item: Item): void
  (e: 'itemVisible', item: Item): void
  (e: 'sort', sort: { sortBy: string; sortDir: 'asc' | 'desc' }): void
  (e: 'dropRowStyling', selector: string, leaving: boolean, event: DragEvent): void
}
defineOptions({
  name: 'OcTable',
  status: 'ready',
  release: '2.1.0'
})

const {
  data,
  idKey = 'id',
  itemDomSelector,
  fields,
  hasHeader = true,
  sticky = false,
  hover = false,
  highlighted = null,
  disabled = [],
  headerPosition = 0,
  paddingX = 'small',
  dragDrop = false,
  lazy = false,
  sortDir = undefined,
  sortBy = undefined,
  caption = '',
  captionVisible = true,
  hasIconsColumn = false
} = defineProps<Props>()

const emit = defineEmits<Emits>()

const { $pgettext } = useGettext()
const domElementSelector = (item: Item) => {
  if (itemDomSelector) {
    return itemDomSelector(item)
  }
  return item[idKey as keyof Item]
}

const constants = {
  EVENT_THEAD_CLICKED,
  EVENT_TROW_CLICKED,
  EVENT_TROW_MOUNTED,
  EVENT_TROW_CONTEXTMENU
}

const isTableSortable = computed(() => fields.some((field) => field.sortable))

function dragOver(event: DragEvent) {
  event.preventDefault()
}
function dragStart(item: Item, event: DragEvent) {
  emit(EVENT_ITEM_DRAGGED, item, event)
}
function dropRowEvent(selector: string, event: DragEvent) {
  emit(EVENT_ITEM_DROPPED, selector, event)
}
function dropRowStyling(selector: string, leaving: boolean, event: DragEvent) {
  emit('dropRowStyling', selector, leaving, event)
}
function isFieldTypeSlot(field: FieldType) {
  return field.type === 'slot'
}
function isFieldTypeCallback(field: FieldType) {
  return ['callback', 'function'].indexOf(field.type) >= 0
}
function extractFieldTitle(field: FieldType) {
  if (Object.prototype.hasOwnProperty.call(field, 'title')) {
    return field.title
  }
  return field.name
}
function extractTableProps() {
  return {
    class: tableClasses.value
  }
}
function extractThProps(field: FieldType, index: number) {
  const props = extractCellProps(field)
  props.class = `oc-table-header-cell oc-table-header-cell-${field.name}`
  if (Object.prototype.hasOwnProperty.call(field, 'thClass')) {
    props.class += ` ${field.thClass}`
  }
  if (sticky) {
    props.style = `top: ${headerPosition}px;`
  }

  if (index === 0) {
    props.class += ` oc-pl-${getSizeClass(paddingX)} `
  }

  if (index === fields.length - 1) {
    props.class += ` oc-pr-${getSizeClass(paddingX)}`
  }

  return props
}
function extractTbodyTrProps(item: Item, index: number) {
  return {
    ...(lazy && { lazy: { colspan: fullColspan.value } }),
    class: [
      'oc-tbody-tr',
      `oc-tbody-tr-${domElementSelector(item) || index}`,
      isHighlighted(item) ? 'oc-table-highlighted' : undefined,
      isDisabled(item) ? 'oc-table-disabled' : undefined
    ].filter(Boolean)
  }
}
function extractTdProps(field: FieldType, index: number, item: Item) {
  const props = extractCellProps(field)
  props.class = `oc-table-data-cell oc-table-data-cell-${field.name}`
  if (Object.prototype.hasOwnProperty.call(field, 'tdClass')) {
    props.class += ` ${field.tdClass}`
  }
  if (Object.prototype.hasOwnProperty.call(field, 'wrap')) {
    props.wrap = field.wrap
  }

  if (index === 0) {
    props.class += ` oc-pl-${getSizeClass(paddingX)} `
  }

  if (index === fields.length - 1) {
    props.class += ` oc-pr-${getSizeClass(paddingX)}`
  }

  if (Object.prototype.hasOwnProperty.call(field, 'accessibleLabelCallback')) {
    props['aria-label'] = field.accessibleLabelCallback(item)
  }

  return props
}
function extractCellProps(field: FieldType): Record<string, string> {
  return {
    ...(field?.alignH && { alignH: field.alignH }),
    ...(field?.alignV && { alignV: field.alignV }),
    ...(field?.width && { width: field.width }),
    class: undefined,
    wrap: undefined,
    style: undefined
  }
}
function isHighlighted(item: Item) {
  if (!highlighted) {
    return false
  }

  if (Array.isArray(highlighted)) {
    return highlighted.indexOf(item[idKey as keyof Item]) > -1
  }

  return highlighted === item[idKey as keyof Item]
}
function isDisabled(item: Item) {
  if (!disabled.length) {
    return false
  }

  return disabled.indexOf(item[idKey as keyof Item]) > -1
}

function cellKey(field: FieldType, index: number, item: Item) {
  const prefix = [item[idKey as keyof Item], index + 1].filter(Boolean)

  if (isFieldTypeSlot(field)) {
    return [...prefix, field.name].join('-')
  }

  if (isFieldTypeCallback(field)) {
    return [...prefix, field.callback(item[field.name as keyof Item])].join('-')
  }

  return [...prefix, item[field.name as keyof Item]].join('-')
}

function fieldIsSortable({ sortable }: FieldType) {
  return !!sortable
}
function handleSort(field: FieldType) {
  if (!fieldIsSortable(field)) {
    return
  }

  let sortedDir = sortDir
  // toggle sortDir if already sorted by this column
  if (sortBy === field.name && sortDir !== undefined) {
    sortedDir = sortDir === SORT_DIRECTION_DESC ? SORT_DIRECTION_ASC : SORT_DIRECTION_DESC
  }
  // set default sortDir of the field when sortDir not set or sortBy changed
  if (sortBy !== field.name || sortDir === undefined) {
    sortedDir = (field.sortDir || SORT_DIRECTION_DESC) as 'asc' | 'desc'
  }

  /**
   * Triggers when table heads are clicked
   *
   * @property {string} sortBy requested column to sort by
   * @property {string} sortDir requested order to sort in (either asc or desc)
   */
  emit('sort', {
    sortBy: field.name,
    sortDir: sortedDir
  })
}
const tableClasses = computed(() => {
  const result = ['oc-table']

  if (hover) {
    result.push('oc-table-hover')
  }

  if (sticky) {
    result.push('oc-table-sticky')
  }

  return result
})

const fullColspan = computed(() => {
  return fields.length
})

function getAriaSortValue(field: string): string | null {
  if (unref(sortBy) !== field) {
    return null
  }

  return unref(sortDir) === 'asc' ? 'ascending' : 'descending'
}
</script>
<style lang="scss">
.oc-table {
  border-collapse: collapse;
  border-spacing: 0;
  color: var(--oc-color-text-default);
  width: 100%;

  &-hover tr {
    transition: background-color $transition-duration-short ease-in-out;
  }

  tr {
    outline: none;
    height: var(--oc-size-height-table-row);
  }

  tr + tr {
    border-top: 1px solid var(--oc-color-border);
  }

  &-hover tr:not(&-footer-row):hover {
    background-color: var(--oc-color-background-hover);
  }

  &-highlighted {
    background-color: var(--oc-color-background-highlight) !important;
  }

  &-accentuated {
    background-color: var(--oc-color-background-accentuate);
  }

  &-disabled {
    background-color: var(--oc-color-background-muted);
    opacity: 0.7;
    filter: grayscale(0.6);
    pointer-events: none;
  }

  &-sticky {
    position: relative;

    .oc-table-header-cell {
      background-color: var(--oc-color-background-default);
      position: sticky;
      z-index: 1;
    }
  }

  .highlightedDropTarget {
    background-color: var(--oc-color-input-border);
  }

  &-thead-content {
    vertical-align: middle;
    display: inline-table;
    color: var(--oc-color-swatch-passive-default);
    &:hover {
      text-decoration: underline;
    }
  }

  &-footer {
    border-top: 1px solid var(--oc-color-border);

    &-cell {
      color: var(--oc-color-text-muted);
      font-size: 0.875rem;
      line-height: 1.4;
      padding: var(--oc-space-xsmall);
    }
  }
}
.oc-button-sort {
  display: flex !important;
  justify-content: start;
  .oc-icon {
    &:hover {
      background-color: var(--oc-color-background-hover);
    }
  }
}
</style>
