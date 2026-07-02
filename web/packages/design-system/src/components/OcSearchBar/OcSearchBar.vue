<template>
  <oc-grid
    flex
    :role="isFilter ? undefined : 'search'"
    class="oc-search oc-flex-middle"
    :class="{ 'oc-search-small': small }"
  >
    <div class="oc-width-expand oc-position-relative">
      <input
        ref="searchInput"
        role="combobox"
        aria-autocomplete="list"
        :aria-expanded="ariaExpanded"
        :aria-controls="ariaControls"
        :class="inputClass"
        :aria-label="label"
        :value="searchQuery"
        :disabled="loading"
        :placeholder="placeholder"
        @input="onType(($event.target as HTMLInputElement).value)"
        @keydown.enter="onSearch"
        @keyup="$emit('keyup', $event)"
      />
      <slot name="locationFilter" />
      <oc-button
        v-if="icon"
        v-oc-tooltip="$gettext('Search')"
        :aria-label="$gettext('Search')"
        class="oc-position-small oc-position-center-right oc-mt-rm"
        appearance="raw"
        @click.prevent.stop="$emit('advancedSearch', $event)"
      >
        <oc-icon v-show="!loading" :name="icon" size="small" fill-type="line" variation="passive" />
        <oc-spinner
          v-show="loading"
          :size="spinnerSize"
          :aria-label="loadingAccessibleLabelValue"
        />
      </oc-button>
    </div>
    <div class="oc-search-button-wrapper" :class="{ 'oc-invisible-sr': buttonHidden }">
      <oc-button
        tabindex="-1"
        class="oc-search-button"
        variation="primary"
        appearance="filled"
        :size="small ? 'small' : 'medium'"
        :disabled="loading || searchQuery.length < 1"
        @click="onSearch"
      >
        {{ buttonLabel }}
      </oc-button>
    </div>
    <oc-button
      v-if="showCancelButton"
      :variation="cancelButtonVariation"
      :appearance="cancelButtonAppearance"
      class="oc-ml-m"
      @click="onCancel"
    >
      <span v-text="$gettext('Cancel')" />
    </oc-button>
  </oc-grid>
</template>

<script lang="ts" setup>
import { computed, unref, ref, useSlots, watch } from 'vue'
import OcButton from '../OcButton/OcButton.vue'
import OcGrid from '../OcGrid/OcGrid.vue'
import OcIcon from '../OcIcon/OcIcon.vue'
import OcSpinner from '../OcSpinner/OcSpinner.vue'
import { useGettext } from 'vue3-gettext'

/**
 * @component OcSearchBar
 * @description
 * The OcSearchBar component is an input element used for searching server-side resources or filtering local results.
 * It supports features like type-ahead search, trimming input, and customizable buttons for search and cancel actions.
 *
 * @Accessibility
 * Landmark role=search**: Communicates its purpose as the main site search to screen readers. Use `isFilter="true"` to disable the landmark role if used as a filter form.
 * Submit button**: Ensures a submit button exists, even if visually hidden (`buttonHidden="true"`).
 * Loading spinner aria-label**: Set via `loadingAccessibleLabel` or defaults to "Loading results".
 *
 * @props
 * @prop {string|null} [value=null] - The search query value.
 * @prop {string} [icon='search'] - The icon to display in the search bar.
 * @prop {string} [placeholder=''] - Placeholder text for the input field.
 * @prop {string} [label=''] - Aria-label for the input field.
 * @prop {boolean} [small=false] - Whether the search bar should be smaller in size.
 * @prop {string} [ariaExpanded='false'] - Indicates if the dropdown is expanded (for combobox).
 * @prop {string} [ariaControls=''] - ID of the dropdown element controlled by this input.
 * @prop {string} [buttonLabel='Search'] - Label for the search button.
 * @prop {boolean} [buttonHidden=false] - Whether to hide the search button visually.
 * @prop {boolean} [typeAhead=false] - If true, triggers the search event on each character input.
 * @prop {boolean} [trimQuery=true] - Automatically trims whitespaces around the search term.
 * @prop {boolean} [loading=false] - If true, disables input and shows a loading spinner.
 * @prop {boolean} [isFilter=false] - If true, removes the search landmark role.
 * @prop {string} [loadingAccessibleLabel=''] - Aria-label for the loading spinner.
 * @prop {boolean} [showCancelButton=false] - Whether to show a cancel button.
 * @prop {'passive'|'primary'|'danger'|'success'|'warning'|'brand'} [cancelButtonVariation='primary'] - Variation of the cancel button.
 * @prop {'outline'|'filled'|'raw'|'raw-inverse'} [cancelButtonAppearance='raw'] - Appearance of the cancel button.
 * @prop {Function} [cancelHandler=() => {}] - Handler function for the cancel button click.
 *
 * @emits
 * @event advancedSearch {MouseEvent} - Emitted when the advanced search button is clicked.
 * @event clear {Event} - Emitted when the search input is cleared.
 * @event input {string} - Emitted on input change.
 * @event keyup {KeyboardEvent} - Emitted on keyup event in the input field.
 * @event search {string} - Emitted when the search button is clicked or enter is pressed.
 *
 * @example
 *   <OcSearchBar
 *     :value="searchQuery"
 *     icon="search"
 *     placeholder="Search for items"
 *     label="Search"
 *     :small="false"
 *     buttonLabel="Search"
 *     :buttonHidden="false"
 *     :typeAhead="true"
 *     :trimQuery="true"
 *     :loading="isLoading"
 *     :isFilter="false"
 *     loadingAccessibleLabel="Loading search results"
 *     :showCancelButton="true"
 *     cancelButtonVariation="primary"
 *     cancelButtonAppearance="outline"
 *     :cancelHandler="onCancel"
 *     @search="onSearch"
 *     @input="onInput"
 *   />
 */

interface Props {
  value?: string | null
  icon?: string
  placeholder?: string
  label?: string
  small?: boolean
  buttonLabel?: string
  buttonHidden?: boolean
  typeAhead?: boolean
  trimQuery?: boolean
  loading?: boolean
  isFilter?: boolean
  loadingAccessibleLabel?: string
  showCancelButton?: boolean
  cancelButtonVariation?: 'passive' | 'primary' | 'danger' | 'success' | 'warning' | 'brand'
  cancelButtonAppearance?: 'outline' | 'filled' | 'raw' | 'raw-inverse'
  cancelHandler?: () => void
  ariaExpanded?: string
  ariaControls?: string
}

interface Emits {
  (e: 'advancedSearch', event: MouseEvent): void
  (e: 'clear', event: Event): void
  (e: 'input', event: string): void
  (e: 'keyup', event: KeyboardEvent): void
  (e: 'search', event: string): void
}
defineOptions({
  name: 'OcSearchBar',
  status: 'ready',
  release: '1.0.0'
})
const {
  value = null,
  icon = 'search',
  placeholder = '',
  label = '',
  small = false,
  buttonLabel = 'Search',
  buttonHidden = false,
  typeAhead = false,
  trimQuery = true,
  loading = false,
  isFilter = false,
  loadingAccessibleLabel = '',
  showCancelButton = false,
  cancelButtonVariation = 'primary',
  cancelButtonAppearance = 'raw',
  cancelHandler = () => {},
  ariaExpanded = 'false',
  ariaControls = ''
} = defineProps<Props>()
const emit = defineEmits<Emits>()
const slots = useSlots()
const query = ref<string>('')
const { $gettext } = useGettext()

function onSearch() {
  /**
   * Search event on filter or search user input
   * @event search
   * @type {event}
   */
  emit('search', unref(query))
}
function onType(currentQuery: string) {
  query.value = trimQuery ? currentQuery.trim() : currentQuery
  /**
   * Input event to support model directive
   * @event Input
   * @type {event}
   */
  emit('input', unref(query))
  if (typeAhead) onSearch()
}

function onCancel() {
  query.value = ''
  onType('')
  onSearch()
  cancelHandler()
}

watch(
  () => value,
  () => {
    if (!value) {
      query.value = ''
    }
  }
)
const inputIconRightPadding = computed(() => {
  if (slots.locationFilter?.().length > 0) {
    return '125px'
  }
  return '48px'
})

const searchQuery = computed(() => {
  // please don't treat empty string the same as null...
  return value === null ? unref(query) : value
})
const spinnerSize = computed(() => {
  if (small) {
    return 'xsmall'
  }
  return 'medium'
})
const inputClass = computed(() => {
  const classes = ['oc-search-input', 'oc-input']

  !buttonHidden && classes.push('oc-search-input-button')

  if (icon || slots.locationFilter) {
    classes.push('oc-search-input-icon')
  }

  return classes
})
const loadingAccessibleLabelValue = computed(() => {
  return loadingAccessibleLabel || $gettext('Loading results')
})
</script>

<style lang="scss">
.oc-search {
  min-width: $form-width-medium;

  &-button {
    border-bottom-left-radius: 0;
    border-top-left-radius: 0;
    // Prevent double borders
    // from input and button
    transform: translateX(-1px);
    z-index: 0;
  }

  &-icon {
    align-items: center;
    bottom: 0;
    color: var(--oc-color-text-muted);
    display: inline-flex;
    justify-content: center;
    left: 0;
    position: absolute;
    top: 0;
    width: 40px;
  }

  &-input {
    border-radius: 25px !important;
    border: none;
    padding: var(--oc-space-medium);
    color: var(--oc-color-input-text-muted) !important;

    &:focus {
      background-color: var(--oc-color-input-bg);
      border-color: var(--oc-color-input-text-default);
      color: var(--oc-color-input-text-default);
      background-image: none;
    }

    &::-ms-clear,
    &::-ms-reveal {
      display: none;
    }
  }

  &-input-icon {
    padding-left: var(--oc-space-xlarge) !important;
    padding-right: v-bind(inputIconRightPadding) !important;
  }

  &-input-button {
    border-bottom-right-radius: 0;
    border-top-right-radius: 0;
  }

  &-clear {
    right: var(--oc-space-large) !important;
  }

  &-small {
    .oc-search-input {
      height: 30px;
      line-height: 28px;
      padding-left: var(--oc-space-xlarge);
    }

    .oc-icon {
      &,
      svg {
        height: 18px;
        width: 18px;
      }
    }
  }
}
</style>
