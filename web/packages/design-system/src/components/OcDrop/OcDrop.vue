<!-- eslint-disable vuejs-accessibility/no-static-element-interactions -->
<!-- eslint-disable vuejs-accessibility/click-events-have-key-events -->
<template>
  <div :id="dropId" ref="drop" class="oc-drop oc-box-shadow-medium oc-rounded" @click="onClick">
    <!-- eslint-enable vuejs-accessibility/no-static-element-interactions -->
    <!-- eslint-enable vuejs-accessibility/click-events-have-key-events -->
    <div
      v-if="$slots.default"
      :class="['oc-card oc-card-body oc-background-secondary', paddingClass]"
    >
      <slot />
    </div>
    <slot v-else name="special" />
  </div>
</template>

<script lang="ts" setup>
import tippy, { ReferenceElement, hideAll, Props as TippyProps } from 'tippy.js'
import { Modifier, detectOverflow } from '@popperjs/core'
import { destroy, hideOnEsc } from '../../directives/OcTooltip'
import { getSizeClass, uniqueId, AvailableSizeType } from '../../helpers'
import { onBeforeUnmount, onMounted, ref, unref, useTemplateRef, computed, watch } from 'vue'

/**
 * @component OcDrop
 * @description Position any element in relation to another element with a dropdown/popover interface.
 * Uses tippy.js for positioning and popper.js for overflow handling.
 *
 * @prop {String} [dropId] - Unique ID for the drop element. Defaults to auto-generated ID.
 * @prop {Object} [popperOptions] - Specifies custom Popper options
 * @prop {String} [toggle] - CSS selector for the element that triggers the drop. Defaults to previous sibling.
 * @prop {String} [position='bottom-start'] - Position of the drop relative to the target.
 *   Options: 'top-start', 'right-start', 'bottom-start', 'left-start', 'auto-start',
 *            'top-end', 'right-end', 'bottom-end', 'left-end', 'auto-end'
 * @prop {String} [mode='click'] - How the drop is triggered. Options: 'click', 'hover', 'manual'
 * @prop {Boolean} [closeOnClick=false] - Whether the drop closes when clicked inside.
 * @prop {Boolean} [isNested=false] - Whether this drop is nested inside another drop.
 * @prop {String} [target=null] - CSS selector for target element. Replaces default target selection.
 * @prop {String} [paddingSize='medium'] - Padding size applied to the drop content.
 *   Options: 'xsmall', 'small', 'medium', 'large', 'xlarge', 'xxlarge', 'xxxlarge', 'remove'
 * @prop {String} [offset] - Offset distance in the format "x, y" (e.g. "10, 20").
 *
 * @event showDrop - Emitted when the drop is shown
 * @event hideDrop - Emitted when the drop is hidden
 *
 * @slot default - Default content of the drop. Wrapped in a card with the specified padding.
 * @slot special - Special content slot used when default slot is not provided (no automatic styling).
 *
 * @exposes {Function} show - Show the drop with an optional duration
 * @exposes {Function} hide - Hide the drop with an optional duration
 * @exposes {Object} tippy - Reference to the internal tippy.js instance
 */

interface Props {
  dropId?: string
  popperOptions?: TippyProps['popperOptions']
  toggle?: string
  position?:
    | 'top-start'
    | 'right-start'
    | 'bottom-start'
    | 'left-start'
    | 'auto-start'
    | 'top-end'
    | 'right-end'
    | 'bottom-end'
    | 'left-end'
    | 'auto-end'
  mode?: 'click' | 'hover' | 'manual'
  closeOnClick?: boolean
  isNested?: boolean
  target?: string
  paddingSize?: AvailableSizeType | 'remove'
  offset?: string
  sameWidthAsTarget?: boolean
  focusOnOpen?: boolean
}
interface Emits {
  (e: 'hideDrop'): void
  (e: 'showDrop'): void
}

defineOptions({
  name: 'OcDrop',
  status: 'ready',
  release: '1.0.0'
})

const {
  dropId = uniqueId('oc-drop-'),
  popperOptions = {},
  toggle = '',
  position = 'bottom-start',
  mode = 'click',
  closeOnClick = false,
  isNested = false,
  target = null,
  paddingSize = 'medium',
  offset = '',
  sameWidthAsTarget = false,
  focusOnOpen = false
} = defineProps<Props>()

const emit = defineEmits<Emits>()

const drop = useTemplateRef<HTMLElement>('drop')
const tippyInstance = ref(null)
const triggerEl = ref<Element | null>(null)

const getFocusableItems = () =>
  Array.from(
    unref(drop)?.querySelectorAll<HTMLElement>(
      'a, button:not([disabled]), [tabindex]:not([tabindex="-1"])'
    ) ?? []
  )

const onTriggerKeydown = (e: KeyboardEvent) => {
  if (!unref(tippyInstance)?.state.isVisible) {
    return
  }
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    getFocusableItems()[0]?.focus()
  } else if (e.key === 'Tab') {
    e.preventDefault()
    hide()
  }
}

const onDropKeydown = (e: KeyboardEvent) => {
  const items = getFocusableItems()
  const currentIdx = items.indexOf(document.activeElement as HTMLElement)

  if (e.key === 'ArrowDown') {
    e.preventDefault()
    e.stopPropagation()
    const next = items[currentIdx + 1] ?? items[0]
    next?.focus()
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    e.stopPropagation()
    const prev = items[currentIdx - 1] ?? items[items.length - 1]
    prev?.focus()
  } else if (e.key === 'Escape') {
    e.preventDefault()
    e.stopPropagation()
    hide()
    ;(unref(triggerEl) as HTMLElement)?.focus()
  } else if (e.key === 'Tab') {
    e.preventDefault()
    hide()
    ;(unref(triggerEl) as HTMLElement)?.focus()
  }
}

const show = (duration?: number) => {
  unref(tippyInstance)?.show(duration)
}
const hide = (duration?: number) => {
  unref(tippyInstance)?.hide(duration)
}

const onClick = () => {
  if (closeOnClick) {
    hide()
  }
}

const onFocusOut = (event: FocusEvent) => {
  const tippyBox = unref(drop)?.closest('.tippy-box')
  if (!tippyBox) {
    return
  }
  const focusLeft = event.relatedTarget && !tippyBox.contains(event.relatedTarget as Node)
  if (focusLeft) {
    hide()
  }
}

const triggerMapping = computed(() => {
  return (
    {
      hover: 'mouseenter focus'
    }[mode] || mode
  )
})

const paddingClass = computed(() => {
  return `oc-p-${getSizeClass(paddingSize)}`
})

watch(
  () => [position, mode],
  () => {
    if (tippyInstance.value) {
      tippyInstance.value.setProps({ placement: position, trigger: triggerMapping.value })
    }
  },
  { immediate: true }
)

function initializeTippy() {
  destroy(unref(tippyInstance))
  const to = unref(target)
    ? document.querySelector(unref(target))
    : unref(toggle)
      ? document.querySelector(unref(toggle))
      : drop.value.previousElementSibling
  const content = drop.value

  if (!to || !content) {
    return
  }

  const sameWidthModifier = {
    name: 'sameWidth',
    enabled: true,
    phase: 'beforeWrite',
    requires: ['computeStyles'],
    fn: ({ state }) => {
      state.styles.popper.width = `${state.rects.reference.width}px`
    },
    effect: ({ state }) => {
      state.elements.popper.style.width = `${state.elements.reference.offsetWidth}px`
    }
  }

  const config: any = {
    trigger: unref(triggerMapping),
    placement: unref(position),
    arrow: false,
    hideOnClick: true,
    interactive: true,
    plugins: [hideOnEsc],
    theme: 'none',
    maxWidth: 416,
    offset: unref(offset) ?? 0,
    ...(!unref(isNested) && {
      onShow: (instance: ReferenceElement) => {
        emit('showDrop')
        hideAll({ exclude: instance })
      },
      onShown: () => {
        if (focusOnOpen) {
          getFocusableItems()[0]?.focus()
        }
      },
      onHide: () => {
        emit('hideDrop')
      }
    }),
    popperOptions: {
      ...unref(popperOptions),
      modifiers: [
        ...(unref(popperOptions)?.modifiers ? unref(popperOptions).modifiers : []),
        ...(unref(sameWidthAsTarget) ? [sameWidthModifier] : []),
        {
          name: 'fixVerticalPosition',
          enabled: true,
          phase: 'beforeWrite',
          requiresIfExists: ['offset', 'preventOverflow', 'flip'],
          fn({ state }) {
            const overflow = detectOverflow(state)
            const dropHeight = state.modifiersData.fullHeight || state.elements.popper.offsetHeight
            const dropYPos = overflow.top * -1 - 10
            const availableHeight = dropYPos + dropHeight + overflow.bottom * -1
            const spaceBelow = availableHeight - dropYPos
            const spaceAbove = availableHeight - spaceBelow

            if (dropHeight > spaceBelow && dropHeight > spaceAbove) {
              /*
                  if context menu placement from the top 'dropYPos' is the same or less than space above
                  and placement is right-start or left-start
                  then subtract the dropYPos from spaceAbove and set the drop on top of the screen
                */
              if (
                dropYPos <= spaceAbove &&
                ['right-start', 'left-start'].includes(state.placement)
              ) {
                state.styles.popper.top = `${spaceAbove - dropYPos}px`
                state.modifiersData.fullHeight = dropHeight
              } else {
                // place drop on top of screen because of limited screen estate above and below
                state.styles.popper.top = `-${dropYPos}px`
                state.modifiersData.fullHeight = dropHeight
              }
            }

            if (dropHeight > availableHeight) {
              // drop is bigger than total available height
              state.styles.popper.maxHeight = `${availableHeight - 10}px`
              state.styles.popper.overflowY = `auto`
              state.styles.popper.overflowX = `hidden`
            }
          }
        } as Modifier<'fixVerticalPosition', unknown>
      ]
    },
    content,
    role: 'listbox'
  }

  if (unref(target)) {
    config.triggerTarget = unref(toggle)
      ? document.querySelector(unref(toggle))
      : drop.value.previousElementSibling
  }

  tippyInstance.value = tippy(to, config)
  triggerEl.value = to
  to.addEventListener('keydown', onTriggerKeydown)
}

onMounted(() => {
  unref(drop).addEventListener('focusout', onFocusOut)
  unref(drop).addEventListener('keydown', onDropKeydown)
  initializeTippy()
})

onBeforeUnmount(() => {
  unref(drop).removeEventListener('focusout', onFocusOut)
  unref(drop).removeEventListener('keydown', onDropKeydown)
  unref(triggerEl)?.removeEventListener('keydown', onTriggerKeydown)
  destroy(unref(tippyInstance))
})

defineExpose({ show, hide, tippy: tippyInstance })
</script>

<style lang="scss">
.tippy-box[data-theme~='none'] {
  background-color: transparent;
  font-size: inherit;
  line-height: inherit;

  .tippy-content {
    // note: needed so that the box shadow from `oc-box-shadow-medium` doesn't get suppressed
    padding: 8px;
  }

  li.oc-menu-item-hover {
    a,
    .item-has-switch,
    button:not([role='switch']) {
      box-sizing: border-box;
      padding: var(--oc-space-small);
      color: var(--oc-color-swatch-passive-default);

      &:focus:not([disabled]),
      &:hover:not([disabled]) {
        background-color: var(--oc-color-background-hover);

        text-decoration: none !important;
        border-radius: 5px;
      }

      &:hover span {
        color: var(--oc-color-swatch-brand-hover) !important;
      }

      span {
        text-decoration: none !important;
      }
    }
  }
}

.oc-drop {
  max-width: 100%;
  width: 300px;

  .oc-card {
    border-radius: var(--oc-space-small) !important;
  }
}
</style>

<docs>
```js
<template>
  <div class="oc-button-group oc-mt-s">
    <oc-button id="my_menu" class="oc-mr-s">Menu</oc-button>
    <oc-drop toggle="#my_menu" mode="click">
      <ul>
        <li icon="create_new_folder" autofocus>Item with icon</li>
        <li>Item without icon</li>
        <li>Active item</li>
      </ul>
    </oc-drop>

    <oc-button id="my_filter" class="oc-mr-s">Filter</oc-button>
    <oc-drop toggle="#my_filter" mode="hover">
      <p>
        Lets filter:
      </p>
      <ul class="oc-list">
        <li>
          <oc-checkbox label=""/>
          <span class="oc-text-muted">Show Files</span>
        </li>
        <li>
          <oc-checkbox label=""/>
          <span class="oc-text-muted">Show Folders</span>
        </li>
        <li>
          <oc-search-bar small placeholder="Filter by name" :button="false" label=""/>
        </li>
      </ul>
    </oc-drop>

    <oc-button id="my_advanced" class="oc-mr-s">Advanced</oc-button>
    <oc-drop dropId="oc-drop" toggle="#my_advanced" mode="click" closeOnClick>
      <div slot="special" class="oc-card">
        <div class="oc-card-header">
          <h3 class="oc-card-title">
            Advanced
          </h3>
        </div>
        <div class="oc-card-body">
          <p>
            I'm a slightly more advanced drop down and I'll be closed as soon as you click on me.
          </p>
        </div>
      </div>
    </oc-drop>

    <oc-button id="my_submenu_parent"> With submenu</oc-button>
    <oc-drop
      id="drop"
      ref="submenu_parent"
      drop-id="oc-drop"
      toggle="#my_submenu_parent"
      mode="click"
      style="max-width: 200px"
    >
      <oc-list class="user-menu-list">
        <li>
          <oc-button appearance="raw"> Menu item 1</oc-button>
        </li>
        <li>
          <oc-button id="menu_with_submenu" appearance="raw">
            Menu item 2
            <oc-icon
              name="arrow-drop-right"
              fill-type="line"
              class="oc-p-xs"
            />
          </oc-button>
          <oc-drop
            id="submenu"
            toggle="#menu_with_submenu"
            mode="hover"
            position="right-start"
            isNested
            closeOnClick
            style="max-width: 200px"
          >
            <oc-list class="user-menu-list">
              <li>
                <oc-button appearance="raw"> Submenu item 1</oc-button>
              </li>
              <li>
                <oc-button appearance="raw"> Submenu item 2</oc-button>
              </li>
            </oc-list>
          </oc-drop>
        </li>
        <li>
          <oc-button appearance="raw"> Menu item 3</oc-button>
        </li>
      </oc-list>
    </oc-drop>
  </div>
</template>
```

### Custom target
```js
<div>
  <div>
    <p id="target">This is the target of the drop</p>
  </div>
  <oc-button id="custom-target-toggle">Trigger drop</oc-button>
  <oc-drop dropId="oc-drop-custom-target" toggle="#custom-target-toggle" target="#target" mode="click" closeOnClick>
    I am attached to a custom element
  </oc-drop>
</div>
```

### Open drop programatically
```js
<template>
  <div>
    <oc-button id="manual-target" @click="open">Open</oc-button>
    <oc-drop ref="drop" mode="manual" target="#manual-target">
      I am triggered manually
    </oc-drop>
  </div>
</template>
<script>
  export default {
    methods: {
      open() {
        this.$refs.drop.show()
      }
    }
  }
</script>
```
</docs>
