<template>
  <div
    class="oc-filter-chip oc-flex"
    :class="{ 'oc-filter-chip-toggle': isToggle, 'oc-filter-chip-raw': raw }"
  >
    <oc-hidden-announcer :announcement="filterToggleAnnouncement" level="polite" />
    <oc-button
      :id="id"
      class="oc-filter-chip-button oc-pill"
      :class="{ 'oc-filter-chip-button-selected': filterActive }"
      :aria-pressed="filterActive"
      appearance="raw"
      @click="isToggle ? handleToggleFilter() : false"
    >
      <oc-icon
        :class="filterActive ? 'oc-filter-check-icon-active' : 'oc-filter-check-icon-inactive'"
        name="check"
        size="small"
        color="var(--oc-color-text-inverse)"
      />
      <slot name="active" :selected-item-names="selectedItemNames">
        <span
          class="oc-text-truncate oc-filter-chip-label"
          v-text="!!selectedItemNames.length ? selectedItemNames[0] : filterLabel"
        />
      </slot>
      <span v-if="selectedItemNames.length > 1" v-text="` +${selectedItemNames.length - 1}`" />
      <oc-icon v-if="!filterActive && !isToggle" name="arrow-down-s" size="small" />
    </oc-button>
    <oc-drop
      v-if="!isToggle"
      ref="dropRef"
      :toggle="'#' + id"
      class="oc-filter-chip-drop"
      mode="click"
      padding-size="small"
      :close-on-click="closeOnClick"
      @hide-drop="$emit('hideDrop')"
      @show-drop="$emit('showDrop')"
    >
      <slot />
    </oc-drop>
    <oc-button
      v-if="filterActive"
      v-oc-tooltip="$gettext('Clear filter')"
      class="oc-filter-chip-clear oc-px-xs"
      appearance="raw"
      :aria-label="$gettext('Clear filter')"
      @click="$emit('clearFilter')"
    >
      <oc-icon name="close" size="small" color="var(--oc-color-text-inverse)" />
    </oc-button>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, unref } from 'vue'
import { useGettext } from 'vue3-gettext'
import { uniqueId } from '../../helpers'
import OcDrop from '../OcDrop/OcDrop.vue'

/**
 * @component OcFilterChip
 * @description A filter chip component used for filtering data with optional toggle and dropdown functionality.
 *
 * @props {string} [id] - Unique identifier for the filter chip.
 * @props {string} filterLabel - The label displayed on the filter chip.
 * @props {string[]} [selectedItemNames] - List of selected item names displayed on the chip.
 * @props {boolean} [isToggle=false] - Determines if the chip acts as a toggle button.
 * @props {boolean} [isToggleActive=false] - Indicates if the toggle chip is active.
 * @props {boolean} [raw=false] - If true, applies raw styling to the chip.
 * @props {boolean} [closeOnClick=false] - If true, closes the dropdown when an item is clicked.
 *
 * @emits {void} clearFilter - Emitted when the filter is cleared.
 * @emits {void} hideDrop - Emitted when the dropdown is hidden.
 * @emits {void} showDrop - Emitted when the dropdown is shown.
 * @emits {void} toggleFilter - Emitted when the toggle filter is activated.
 *
 * @slots default - Slot for custom dropdown content.
 *
 * @methods hideDrop - Exposes a method to programmatically hide the dropdown.
 *
 * @example
 *  <oc-filter-chip
 *  id="my-filter-chip"
 *  filterLabel="Filter by category"
 *  selectedItemNames="['Category 1', 'Category 2']"
 *  :isToggle="true"
 *  :isToggleActive="true"
 *  :closeOnClick="true"
 *  :raw="false"
 *  @clearFilter="handleClearFilter"
 *  @hideDrop="handleHideDrop"
 *  @showDrop="handleShowDrop"
 *  @toggleFilter="handleToggleFilter"
 *  >
 */

interface Props {
  id?: string
  filterLabel: string
  selectedItemNames?: string[]
  isToggle?: boolean
  isToggleActive?: boolean
  closeOnClick?: boolean
  raw?: boolean
}
interface Emits {
  (e: 'clearFilter'): void
  (e: 'hideDrop'): void
  (e: 'showDrop'): void
  (e: 'toggleFilter'): void
}

defineOptions({
  name: 'OcFilterChip',
  status: 'ready',
  release: '15.0.0'
})

const { $pgettext } = useGettext()
const filterToggleAnnouncement = ref('')

const {
  id = uniqueId('oc-filter-chip-'),
  filterLabel,
  selectedItemNames = [],
  isToggle = false,
  isToggleActive = false,
  raw = false,
  closeOnClick = false
} = defineProps<Props>()

const emit = defineEmits<Emits>()
const dropRef = ref<typeof OcDrop>()

const handleToggleFilter = () => {
  filterToggleAnnouncement.value = isToggleActive
    ? $pgettext(
        'Accessibility announcement when a toggle filter is deactivated',
        `${filterLabel} filter removed`
      )
    : $pgettext(
        'Accessibility announcement when a toggle filter is activated',
        `${filterLabel} filter applied`
      )
  emit('toggleFilter')
}

const filterActive = computed(() => {
  if (isToggle) {
    return isToggleActive
  }
  return !!selectedItemNames.length
})

const hideDrop = () => {
  unref(dropRef)?.hide()
}

defineExpose({ hideDrop })
</script>

<style lang="scss">
.oc-filter-chip {
  &-button.oc-pill {
    align-items: center;
    background-color: var(--oc-color-background-default) !important;
    color: var(--oc-color-text-muted) !important;
    border: 1px solid var(--oc-color-text-muted);
    box-sizing: border-box;
    display: inline-flex;
    gap: var(--oc-space-xsmall);
    text-decoration: none;
    font-size: var(--oc-font-size-xsmall);
    line-height: 1rem;
    max-width: 180px;
    padding: var(--oc-space-xsmall) var(--oc-space-small) !important;
    height: 100%;
  }
  &-button-selected.oc-pill,
  &-button-selected.oc-pill:hover {
    background-color: var(--oc-color-swatch-primary-default) !important;
    color: var(--oc-color-text-inverse) !important;
    border-top-left-radius: 99px !important;
    border-bottom-left-radius: 99px !important;
    border-top-right-radius: 0px !important;
    border-bottom-right-radius: 0px !important;
    border: 0;
  }
  &-clear,
  &-clear:hover {
    background-color: var(--oc-color-swatch-primary-default) !important;
    color: var(--oc-color-text-inverse) !important;
    border-top-left-radius: 0px !important;
    border-bottom-left-radius: 0px !important;
    border-top-right-radius: 99px !important;
    border-bottom-right-radius: 99px !important;
  }
  &-clear:not(.oc-filter-chip-toggle .oc-filter-chip-clear),
  &-clear:hover:not(.oc-filter-chip-toggle .oc-filter-chip-clear) {
    margin-left: 1px;
  }
}
.oc-filter-chip-raw {
  .oc-filter-chip-button {
    background-color: transparent !important;
    border: none !important;
  }
}
.oc-filter-check-icon-active {
  transition: all 0.25s ease-in;
  transform: scale(1) !important;
}
.oc-filter-check-icon-inactive {
  transition: all 0.25 ease-in;
  transform: scale(0) !important;
  width: 0px !important;
}

// the focussed button needs to stay above the other to correctly display the focus outline
.oc-filter-chip-button,
.oc-filter-chip-clear {
  &:focus {
    z-index: 9;
  }
}
</style>
