<template>
  <div>
    <label v-if="!labelHidden" :aria-hidden="true" :for="id" class="oc-label" v-text="label" />
    <oc-contextual-helper
      v-if="contextualHelper?.isEnabled"
      v-bind="contextualHelper?.data"
      class="oc-pl-xs"
    ></oc-contextual-helper>
    <vue-select
      ref="select"
      :disabled="disabled || readOnly"
      :filter="filter"
      :loading="loading"
      :searchable="searchable"
      :clearable="clearable"
      :multiple="multiple"
      class="oc-select"
      :class="{ 'oc-select-position-fixed': positionFixed }"
      style="background: transparent"
      :dropdown-should-open="selectDropdownShouldOpen"
      :map-keydown="selectMapKeydown"
      v-bind="additionalAttributes"
      @update:model-value="$emit('update:modelValue', $event)"
      @click="onSelectClick()"
      @search:blur="onSelectBlur()"
      @keydown="onSelectKeyDown($event)"
    >
      <template #search="{ attributes, events }">
        <input
          class="vs__search"
          v-bind="attributes"
          role="combobox"
          :aria-expanded="dropdownOpen ? 'true' : 'false'"
          aria-haspopup="listbox"
          :aria-controls="listboxId"
          :aria-label="label"
          @input="userInput"
          v-on="events"
        />
      </template>
      <template v-for="(index, name) in $slots" #[name]="data">
        <slot v-if="name !== 'search'" :name="name" v-bind="data" />
      </template>
      <template #no-options>
        <div v-text="$gettext('No options available.')" />
      </template>
      <template #spinner="{ loading: loadingSpinner }">
        <oc-spinner v-if="loadingSpinner" />
      </template>
      <template #selected-option-container="{ option, deselect }">
        <span class="vs__selected" :class="{ 'vs__selected-readonly': option.readonly }">
          <slot name="selected-option" v-bind="option">
            <oc-icon v-if="readOnly" name="lock" class="oc-mr-xs" size="small" />
            {{ getOptionLabel(option) }}
          </slot>
          <span v-if="multiple" class="oc-flex oc-flex-middle oc-ml-s oc-mr-xs">
            <oc-icon v-if="option.readonly" class="vs__deselect-lock" name="lock" size="small" />
            <oc-button
              v-else
              appearance="raw"
              :title="$gettext('Deselect %{label}', { label: getOptionLabel(option) })"
              :aria-label="$gettext('Deselect %{label}', { label: getOptionLabel(option) })"
              class="vs__deselect oc-mx-rm"
              @mousedown.stop.prevent
              @click="deselect(option)"
            >
              <oc-icon name="close" size="small" />
            </oc-button>
          </span>
        </span>
      </template>
    </vue-select>

    <div
      v-if="showMessageLine"
      class="oc-text-input-message"
      :class="{
        'oc-text-input-description': !!descriptionMessage,
        'oc-text-input-warning': !!warningMessage,
        'oc-text-input-danger': !!errorMessage
      }"
    >
      <oc-icon
        v-if="messageText !== null && !!descriptionMessage"
        name="information"
        size="small"
        fill-type="line"
      />

      <span
        :id="messageId"
        :class="{
          'oc-text-input-description': !!descriptionMessage,
          'oc-text-input-warning': !!warningMessage,
          'oc-text-input-danger': !!errorMessage
        }"
        v-text="messageText"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import Fuse from 'fuse.js'
import { uniqueId } from '../../helpers'
import {
  onMounted,
  ref,
  unref,
  VNodeRef,
  useAttrs,
  nextTick,
  onBeforeUnmount,
  watch,
  computed
} from 'vue'
import { useGettext } from 'vue3-gettext'
import 'vue-select/dist/vue-select.css'
import { ContextualHelper } from '../../helpers'
// @ts-ignore
import VueSelect from 'vue-select'

/**
 * OcSelect Component
 *
 * A select component with a trigger and dropdown based on [Vue Select](https://vue-select.org/).
 *
 * @component
 * @name OcSelect
 * @status ready
 * @release 4.3.0
 *
 * @props
 * @property {string} [id] - The ID of the element. Defaults to a unique ID.
 * @property {Function} [filter] - Function to filter items when searching. Defaults to a Fuse.js-based search.
 * @property {boolean} [disabled=false] - Disable the select component.
 * @property {string} label - Label of the select component, required for accessibility.
 * @property {boolean} [labelHidden=false] - Hide the label visually but keep it for accessibility.
 * @property {ContextualHelper|null} [contextualHelper=null] - Contextual helper data.
 * @property {string} [optionLabel='label'] - Key to use as label when option is an object.
 * @property {boolean} [searchable=true] - Determines if the select field is searchable.
 * @property {boolean} [clearable=false] - Determines if the select field is clearable.
 * @property {boolean} [loading=false] - Determines if the select field is in a loading state.
 * @property {string|null} [warningMessage=null] - A warning message shown below the select.
 * @property {string|null} [errorMessage=null] - An error message shown below the select.
 * @property {boolean} [fixMessageLine=false] - Reserve vertical space for a one-line message.
 * @property {string|null} [descriptionMessage=null] - A description text shown below the select field.
 * @property {boolean} [multiple=false] - Determines if multiple options can be selected.
 * @property {boolean} [readOnly=false] - Determines if the select field is read-only.
 * @property {boolean} [positionFixed=false] - Sets the dropdown menu to `position: fixed`.
 *
 * @emits
 * @event search:input - Triggered when the search input value changes.
 * @property {string} query - The search query.
 * @event update:modelValue - Triggered when the model value is updated.
 * @property {unknown} value - The updated value.
 *
 * @example
 *   <OcSelect
 *     :id="'example-select'"
 *     :label="'Example Select'"
 *     :options="options"
 *     :multiple="true"
 *     :searchable="true"
 *     :clearable="true"
 *     update:modelValue="onValueUpdate"
 *   />
 */

interface Props {
  id?: string
  filter?: (items: unknown[], search: string, label: string) => unknown[]
  disabled?: boolean
  label: string
  labelHidden?: boolean
  contextualHelper?: ContextualHelper | null
  optionLabel?: string
  getOptionLabel?: (option: unknown) => string
  searchable?: boolean
  clearable?: boolean
  loading?: boolean
  warningMessage?: string | null
  errorMessage?: string | null
  fixMessageLine?: boolean
  descriptionMessage?: string | null
  multiple?: boolean
  readOnly?: boolean
  positionFixed?: boolean
}
interface Emits {
  (event: 'search:input', query: string): void
  (event: 'update:modelValue', value: unknown): void
}
// the keycode property is deprecated in the JS event API, vue-select still works with it though
enum KeyCode {
  Enter = 13,
  ArrowDown = 40,
  ArrowUp = 38
}

defineOptions({
  name: 'OcSelect',
  status: 'ready',
  release: '4.3.0',
  // needed for unit testing, otherwise, it will be rendered as anonymous component -->
  components: { VueSelect }
})
const {
  id = uniqueId('oc-select-'),
  filter = (items: unknown[], search: string, { label }: { label?: string }) => {
    if (items.length < 1) {
      return []
    }

    const fuse = new Fuse(items, {
      ...(label && { keys: [label] }),
      shouldSort: true,
      threshold: 0,
      ignoreLocation: true,
      distance: 100,
      minMatchCharLength: 1
    })

    return search.length ? fuse.search(search).map(({ item }) => item) : items
  },
  disabled = false,
  label,
  labelHidden = false,
  contextualHelper = null,
  optionLabel = 'label',
  getOptionLabel: getOptionLabelProp = null,
  searchable = true,
  clearable = false,
  loading = false,
  warningMessage = null,
  errorMessage = null,
  fixMessageLine = false,
  descriptionMessage = null,
  multiple = false,
  readOnly = false,
  positionFixed = false
} = defineProps<Props>()

const emit = defineEmits<Emits>()
const { $gettext } = useGettext()
const select: VNodeRef = ref()
const attrs = useAttrs()
const listboxId = computed(() => `${id}-listbox`)

const userInput = (event: Event) => {
  /**
   * Triggers when a value of search input is changed
   *
   * @property {string} query search query
   */
  emit('search:input', (event.target as HTMLInputElement).value)
}

const dropdownEnabled = ref(false)
const setDropdownEnabled = (enabled: boolean) => {
  dropdownEnabled.value = enabled
}

const selectDropdownShouldOpen = ({
  noDrop,
  open,
  mutableLoading
}: {
  noDrop?: boolean
  open?: boolean
  mutableLoading?: boolean
}) => {
  return !noDrop && open && !mutableLoading && unref(dropdownEnabled)
}

const onSelectClick = () => {
  setDropdownEnabled(true)
}

const onSelectBlur = () => {
  setDropdownEnabled(false)
}

/**
 * Sets the outline for the highlighted option. This needs to be applied when
 * navigating via keyboard because of a11y.
 */
const setKeyboardOutline = async () => {
  const optionEls = unref(select).$refs.dropdownMenu.querySelectorAll('li')
  const highlightedOption = optionEls[unref(select).typeAheadPointer]
  if (highlightedOption) {
    await nextTick()
    highlightedOption.classList.add('keyboard-outline')
  }
}

const selectMapKeydown = (map: Record<number, (e: KeyboardEvent) => void>) => {
  return {
    ...map,
    [KeyCode.Enter]: (e: KeyboardEvent) => {
      if (!unref(dropdownEnabled)) {
        setDropdownEnabled(true)
        return
      }
      map[KeyCode.Enter](e)
      unref(select).searchEl.focus()
    },
    [KeyCode.ArrowDown]: async (e: KeyboardEvent) => {
      e.preventDefault()
      unref(select).typeAheadDown()

      if (unref(dropdownOpen)) {
        await setKeyboardOutline()
      }
    },
    [KeyCode.ArrowUp]: async (e: KeyboardEvent) => {
      e.preventDefault()
      unref(select).typeAheadUp()

      if (unref(dropdownOpen)) {
        await setKeyboardOutline()
      }
    }
  }
}

const onSelectKeyDown = async (e: KeyboardEvent) => {
  if (e.key === 'Enter' || e.key === 'Tab') {
    if (unref(dropdownOpen)) {
      // set initial outline for highlighted option
      await setKeyboardOutline()
    }

    return
  }

  setDropdownEnabled(true)
}

const setDropdownPosition = () => {
  const dropdownMenu = unref(select).$refs.dropdownMenu
  if (!dropdownMenu) {
    return
  }

  const toggleClientRect = unref(select).$refs.toggle.getBoundingClientRect()
  const dropdownMenuBottomOffset = 25
  const dropdownMenuMaxHeight = Math.min(
    window.innerHeight - toggleClientRect.bottom - dropdownMenuBottomOffset,
    window.innerHeight
  )

  dropdownMenu.style.maxHeight = `${dropdownMenuMaxHeight}px`
  dropdownMenu.style.width = `${toggleClientRect.width}px`
  dropdownMenu.style.top = `${toggleClientRect.top + toggleClientRect.height + 1}px`
  dropdownMenu.style.left = `${toggleClientRect.left}px`
}

const dropdownOpen = computed(() => {
  return unref(select)?.dropdownOpen
})

const getOptionLabel = computed(() => {
  return (
    getOptionLabelProp ||
    ((option: string | Record<string, unknown>): string => {
      if (typeof option === 'object') {
        const key = optionLabel || label
        if (!Object.hasOwn(option, key)) {
          console.warn(
            `[vue-select warn]: Label key "option.${key}" does not` +
              ` exist in options object ${JSON.stringify(option)}.\n` +
              'https://vue-select.org/api/html#getoptionlabel'
          )
          return ''
        }
        return option[key] as string
      }
      return option
    })
  )
})

const additionalAttributes = computed(() => {
  const additionalAttrs: Record<string, unknown> = {}
  additionalAttrs['input-id'] = id
  additionalAttrs['getOptionLabel'] = unref(getOptionLabel)
  additionalAttrs['label'] = optionLabel

  return { ...attrs, ...additionalAttrs }
})
const showMessageLine = computed(() => {
  return fixMessageLine || !!warningMessage || !!errorMessage || !!descriptionMessage
})
const messageText = computed(() => {
  if (errorMessage) {
    return errorMessage
  }

  if (warningMessage) {
    return warningMessage
  }

  return descriptionMessage
})
const messageId = computed(() => {
  return `${id}-message`
})
watch(dropdownOpen, async () => {
  if (positionFixed && unref(dropdownOpen)) {
    await nextTick()
    setDropdownPosition()
  }
})

onMounted(() => {
  if (positionFixed) {
    window.addEventListener('resize', setDropdownPosition)
  }
})

onBeforeUnmount(() => {
  if (positionFixed) {
    window.removeEventListener('resize', setDropdownPosition)
  }
})
</script>
<style lang="scss">
.vs--disabled {
  cursor: not-allowed;

  .vs__clear,
  .vs__dropdown-toggle,
  .vs__open-indicator,
  .vs__search,
  .vs__selected {
    background-color: var(--oc-color-background-muted) !important;
    color: var(--oc-color-input-text-muted) !important;
    pointer-events: none;
  }

  .vs__actions {
    opacity: 0.3;
  }
}

.oc-select {
  line-height: normal;
  padding: 1px 0;
  color: var(--oc-color-input-text-default);

  &-position-fixed {
    .vs__dropdown-menu {
      position: fixed;
      overflow-y: auto;
    }
  }

  .vs {
    &__search {
      color: var(--oc-color-input-text-default);
    }

    &__search::placeholder,
    &__dropdown-toggle,
    &__dropdown-menu {
      -webkit-appearance: none;
      appearance: none;
      background-color: var(--oc-color-background-highlight);
      border-radius: 0;
      border-radius: 5px;
      border: 1px solid var(--oc-color-input-border);
      box-sizing: border-box;
      color: var(--oc-color-input-text-default);
      line-height: inherit;
      margin: 0;
      max-width: 100%;
      outline: none;
      padding: 2px;
      transition-duration: 0.2s;
      transition-timing-function: ease-in-out;
      transition-property: color, background-color;
      width: 100%;
    }

    &__selected-readonly {
      background-color: var(--oc-color-background-muted) !important;
    }

    &__search,
    &__search:focus {
      padding: 0 5px;
    }

    &__dropdown-menu {
      padding: 0;
      background-color: var(--oc-color-background-default);
      margin-top: -1px;
    }

    &__clear,
    &__open-indicator,
    &__deselect {
      fill: var(--oc-color-input-text-default);
    }

    &__deselect {
      margin: 0 var(--oc-space-small);
    }

    &__dropdown-option,
    &__no-options {
      color: var(--oc-color-input-text-default);
      white-space: normal;
      padding: 6px 0.6rem;
      border-radius: 5px;
      line-height: var(--vs-line-height);

      &--highlight,
      &--selected {
        background-color: var(--oc-color-background-hover);
        color: var(--oc-color-swatch-brand-hover);
      }
    }

    &__actions {
      flex-flow: row wrap;
      gap: var(--oc-space-xsmall);
      cursor: pointer;

      svg {
        overflow: visible;
      }
    }

    &__clear svg {
      max-width: var(--oc-space-small);
    }

    &__selected-options {
      flex: auto;
      padding: 0;

      > * {
        padding: 0px 2px;
        margin: 2px 2px 2px 1px;
        color: var(--oc-color-input-text-default);
      }

      > *:not(input) {
        padding-left: 3px;
        background-color: var(--oc-color-background-default);
        fill: var(--oc-color-text-default);
      }
    }
  }

  &.vs--multiple {
    .vs {
      &__selected-options > *:not(input) {
        color: var(--oc-color-input-text-default);
        background-color: var(--oc-color-background-default);
      }
    }
  }

  &:focus-within {
    .vs__dropdown-menu,
    .vs__dropdown-toggle {
      border-color: var(--oc-color-swatch-passive-default);
    }
  }

  .keyboard-outline {
    outline: 2px var(--oc-color-swatch-passive-default) solid !important;
    outline-offset: -2px;
  }
}

.oc-background-highlight {
  .oc-select {
    .vs {
      &__search {
        color: var(--oc-color-input-text-default);
      }

      &__search::placeholder,
      &__dropdown-toggle,
      &__dropdown-menu {
        background-color: var(--oc-color-input-bg);
      }
    }

    &.vs--multiple {
      .vs__selected-options > *:not(input) {
        color: var(--oc-color-input-text-default);
        background-color: var(--oc-color-background-highlight);
      }
    }

    &:focus-within {
      .vs__dropdown-menu,
      .vs__dropdown-toggle {
        background-color: var(--oc-color-background-default);
      }
    }
  }
}

.vs--single {
  &.vs--open .vs__selected {
    opacity: 0.8 !important;
  }

  .vs__selected-options > *:not(input) {
    background-color: transparent !important;
  }
}
</style>
