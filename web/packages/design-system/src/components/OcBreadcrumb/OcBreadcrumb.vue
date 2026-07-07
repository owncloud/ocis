<template>
  <nav
    :id="id"
    :class="`oc-breadcrumb oc-breadcrumb-${variation}`"
    :aria-label="$gettext('Breadcrumbs')"
    v-bind="attrs"
  >
    <ol class="oc-breadcrumb-list oc-flex oc-m-rm oc-p-rm">
      <li
        v-for="(item, index) in displayItems"
        :key="index"
        :data-key="index"
        :data-item-id="item.id"
        :aria-hidden="item.isTruncationPlaceholder"
        :tabindex="item.isTruncationPlaceholder ? -1 : 0"
        :class="[
          'oc-breadcrumb-list-item',
          'oc-flex',
          'oc-flex-middle',
          {
            'oc-invisible-sr':
              hiddenItems.indexOf(item) !== -1 ||
              (item.isTruncationPlaceholder && hiddenItems.length === 0)
          }
        ]"
        @dragover.prevent
        @dragenter.prevent="dropItemStyling(item as BreadcrumbItem, index, false, $event)"
        @dragleave.prevent="dropItemStyling(item as BreadcrumbItem, index, true, $event)"
        @mouseleave="dropItemStyling(item as BreadcrumbItem, index, true, $event as DragEvent)"
        @blur="dropItemStyling(item as BreadcrumbItem, index, true, $event as DragEvent)"
        @drop="dropItemEvent(item as BreadcrumbItem, index)"
      >
        <router-link
          v-if="item.to"
          :aria-current="getAriaCurrent(index)"
          :to="item.isTruncationPlaceholder ? lastHiddenItem.to : item.to"
        >
          <span class="oc-breadcrumb-item-text oc-breadcrumb-item-navigable">{{ item.text }}</span>
        </router-link>
        <oc-button
          v-else-if="item.onClick"
          :aria-current="getAriaCurrent(index)"
          appearance="raw"
          class="oc-flex"
          @click="item.onClick"
        >
          <span
            :class="[
              'oc-breadcrumb-item-text',
              'oc-breadcrumb-item-navigable',
              {
                'oc-breadcrumb-item-text-last': index === displayItems.length - 1
              }
            ]"
            v-text="item.text"
          />
        </oc-button>
        <span
          v-else
          class="oc-breadcrumb-item-text"
          :aria-current="getAriaCurrent(index)"
          tabindex="-1"
          v-text="item.text"
        />
        <oc-icon
          v-if="index !== displayItems.length - 1"
          color="var(--oc-color-text-default)"
          name="arrow-right-s"
          class="oc-mx-xs"
          fill-type="line"
        />
        <template v-if="showContextActions && index === displayItems.length - 1">
          <oc-button
            id="oc-breadcrumb-contextmenu-trigger"
            v-oc-tooltip="contextMenuLabel"
            :aria-label="contextMenuLabel"
            appearance="raw"
          >
            <oc-icon name="more-2" color="var(--oc-color-text-default)" />
          </oc-button>
          <oc-drop
            drop-id="oc-breadcrumb-contextmenu"
            toggle="#oc-breadcrumb-contextmenu-trigger"
            mode="click"
            close-on-click
            :padding-size="contextMenuPadding"
          >
            <!-- @slot Add context actions that open in a dropdown when clicking on the "three dots" button -->
            <slot name="contextMenu" />
          </oc-drop>
        </template>
      </li>
    </ol>
    <oc-button
      v-if="parentFolderTo && displayItems.length > 1"
      appearance="raw"
      type="router-link"
      :aria-label="$gettext('Navigate one level up')"
      :to="parentFolderTo"
      class="oc-breadcrumb-mobile-navigation"
    >
      <oc-icon name="arrow-left-s" fill-type="line" size="large" class="oc-mr-m" />
    </oc-button>
  </nav>
  <div v-if="displayItems.length > 1" class="oc-breadcrumb-mobile-current" v-bind="attrs">
    <span class="oc-text-truncate" aria-current="page" v-text="currentFolder.text" />
  </div>
</template>

<script lang="ts" setup>
import { computed, nextTick, ref, unref, watch, useAttrs } from 'vue'
import { useGettext } from 'vue3-gettext'
import {
  AvailableSizeType,
  EVENT_ITEM_DROPPED_BREADCRUMB,
  uniqueId,
  BreadcrumbItem
} from '../../helpers'
import OcButton from '../OcButton/OcButton.vue'
import OcDrop from '../OcDrop/OcDrop.vue'
import OcIcon from '../OcIcon/OcIcon.vue'
import { RouteLocationPathRaw } from 'vue-router'

/**
 * OcBreadcrumb - component is responsible for showing breadcrumbs
 *
 * @prop {string} [id] - Optional ID for the breadcrumbs. If it's empty, a generated one will be used.
 * @prop {BreadcrumbItem[]} items - Array of breadcrumb items.
 * @prop {string} [variation='default'] - Variation of breadcrumbs. Can be `default` or `lead`.
 * @prop {string} [contextMenuPadding='medium'] - Defines the padding size around the drop content. Defaults to `medium`.
 * Valid values: xsmall, small, medium, large, xlarge, xxlarge, xxxlarge, remove
 * @prop {number} [maxWidth=-1] - Defines the maximum width of the breadcrumb. If the breadcrumb is wider than the given value, it will be reduced from the left side. If the value is -1, the breadcrumb will not be reduced.
 * @prop {number} [truncationOffset=2] - Defines the number of items that should always be displayed at the beginning of the breadcrumb. The default value is 2.
 * @prop {boolean} [showContextActions=false] - Determines if the last breadcrumb item should have context menu actions.
 *
 * @event {RouteLocationPathRaw} itemDroppedBreadcrumb - Event emitted when an item is dropped on the breadcrumb.
 *
 * @example
 * ```vue
 * <template>
 *   <OcBreadcrumb :items="items" variation="lead" @item-dropped-breadcrumb="handleItemDropped" />
 * </template>
 * ```
 */

interface Props {
  id?: string
  items: BreadcrumbItem[]
  variation?: 'default' | 'lead'
  contextMenuPadding?: AvailableSizeType | 'remove'
  maxWidth?: number
  truncationOffset?: number
  showContextActions?: boolean
}
interface Emits {
  (e: typeof EVENT_ITEM_DROPPED_BREADCRUMB, payload: RouteLocationPathRaw): void
}

const attrs = useAttrs()
const emits = defineEmits<Emits>()

defineOptions({
  name: 'OcBreadcrumb',
  status: 'ready',
  release: '1.0.0'
})

const {
  id = uniqueId('oc-breadcrumbs-'),
  items,
  variation = 'default',
  contextMenuPadding = 'medium',
  maxWidth = -1,
  truncationOffset = 2,
  showContextActions = false
} = defineProps<Props>()

const { $gettext } = useGettext()
const visibleItems = ref<BreadcrumbItem[]>([])
const hiddenItems = ref<BreadcrumbItem[]>([])
const displayItems = ref<BreadcrumbItem[]>(items.slice())

const getBreadcrumbElement = (id: string): HTMLElement => {
  return document.querySelector(`.oc-breadcrumb-list [data-item-id="${id}"]`)
}

const isDropAllowed = (item: BreadcrumbItem, index: number): boolean => {
  return !(
    !item.id ||
    index === unref(displayItems).length - 1 ||
    item.isTruncationPlaceholder ||
    item.isStaticNav
  )
}
const dropItemEvent = (item: BreadcrumbItem, index: number) => {
  if (!isDropAllowed(item, index)) {
    return
  }

  if (typeof item.to === 'object') {
    const itemTo = item.to as RouteLocationPathRaw
    itemTo.path = itemTo.path || '/'
    emits(EVENT_ITEM_DROPPED_BREADCRUMB, itemTo)
  }
}

const calculateTotalBreadcrumbWidth = () => {
  let totalBreadcrumbWidth = 100 // 100px margin to the right to avoid breadcrumb from getting too close to the controls
  visibleItems.value.forEach((item) => {
    const breadcrumbElement = getBreadcrumbElement(item.id)
    const itemClientWidth = breadcrumbElement?.getBoundingClientRect()?.width || 0
    totalBreadcrumbWidth += itemClientWidth
  })
  return totalBreadcrumbWidth
}

const reduceBreadcrumb = (offsetIndex: number) => {
  const breadcrumbMaxWidth = maxWidth
  if (!breadcrumbMaxWidth) {
    return
  }
  const totalBreadcrumbWidth = calculateTotalBreadcrumbWidth()

  const isOverflowing = breadcrumbMaxWidth < totalBreadcrumbWidth
  if (!isOverflowing || visibleItems.value.length <= truncationOffset + 1) {
    return
  }
  // Remove from the left side
  const removed = visibleItems.value.splice(offsetIndex, 1)

  hiddenItems.value.push(removed[0])
  reduceBreadcrumb(offsetIndex)
}

const lastHiddenItem = computed(() =>
  hiddenItems.value.length >= 1 ? unref(hiddenItems)[unref(hiddenItems).length - 1] : { to: {} }
)

const renderBreadcrumb = () => {
  displayItems.value = [...items]
  if (displayItems.value.length > truncationOffset - 1) {
    displayItems.value.splice(truncationOffset - 1, 0, {
      text: '...',
      allowContextActions: false,
      to: {} as BreadcrumbItem['to'],
      isTruncationPlaceholder: true
    })
  }
  visibleItems.value = [...displayItems.value]
  hiddenItems.value = []
  nextTick(() => {
    reduceBreadcrumb(truncationOffset)
  })
}

watch([() => maxWidth, () => items], renderBreadcrumb, { immediate: true })

const currentFolder = computed<BreadcrumbItem>(() => {
  if (items.length === 0 || !items) {
    return undefined
  }
  return [...items].reverse()[0]
})
const parentFolderTo = computed(() => {
  return [...items].reverse()[1]?.to
})

const contextMenuLabel = computed(() => {
  return $gettext('Show actions for current folder')
})

const getAriaCurrent = (index: number): 'page' | null => {
  return items.length - 1 === index ? 'page' : null
}

const dropItemStyling = (
  item: BreadcrumbItem,
  index: number,
  leaving: boolean,
  event: DragEvent
) => {
  if (!isDropAllowed(item, index)) {
    return
  }
  const hasFilePayload = (event.dataTransfer?.types || []).some((e) => e === 'Files')
  if (hasFilePayload) return
  if ((event.currentTarget as HTMLElement)?.contains(event.relatedTarget as HTMLElement)) {
    return
  }

  const classList = getBreadcrumbElement(item.id).children[0].classList
  const className = 'oc-breadcrumb-item-dragover'
  leaving ? classList.remove(className) : classList.add(className)
}
</script>

<style lang="scss">
.oc-breadcrumb {
  overflow: visible;
  &-item-dragover {
    transition:
      background 0.06s,
      border 0s 0.08s,
      border-color 0s,
      border-width 0.06s;
    background-color: var(--oc-color-background-highlight);
    box-shadow: 0 0 0 5px var(--oc-color-background-highlight);
    border-radius: 5px;
  }
  &-item-text {
    max-width: 200px;
    white-space: nowrap;
    text-overflow: ellipsis;
    overflow: hidden;

    &-last {
      vertical-align: text-bottom;
    }
  }

  &-item-navigable:hover {
    text-decoration: underline;
  }

  &-mobile-current,
  &-mobile-navigation {
    @media (min-width: $oc-breakpoint-small-default) {
      display: none !important;
    }
  }

  &-list {
    list-style: none;
    align-items: baseline;
    flex-wrap: nowrap;

    @media (max-width: $oc-breakpoint-xsmall-max) {
      display: none !important;
    }

    #oc-breadcrumb-contextmenu-trigger > span {
      vertical-align: middle;
      border: 3px solid transparent;
    }

    #oc-breadcrumb-contextmenu li button {
      display: inline-flex;
    }

    > li button {
      display: inline;
    }

    > :nth-child(n + 2)::before {
      color: var(--oc-color-text-default);
      display: inline-block;
    }

    > :last-child > span {
      color: var(--oc-color-text-default);
    }
  }

  /* stylelint-disable */
  &-list-item {
    a:first-of-type,
    button:first-of-type,
    span:first-of-type {
      font-size: var(--oc-font-size-medium);
      color: var(--oc-color-text-default);
      display: inline-block;
      vertical-align: sub;
      line-height: normal;
    }
  }

  &-lead &-list-item {
    a,
    button,
    span {
      font-size: var(--oc-font-size-large);
    }
  }
}
</style>

<docs>
```js
<template>
<section>
  <div>
    <oc-breadcrumb :items="items" />
  </div>
  <div>
    <oc-breadcrumb :items="items" variation="lead" />
    <oc-breadcrumb :items="items" >
      <template v-slot:contextMenu>
        <p class="oc-my-rm">I'm an example item</p>
      </template>
    </oc-breadcrumb>
  </div>
</section>
</template>
<script>
  export default {
    data: () => {
      return {
        items: [
          {text:'First folder',to:{path:'folder'}},
          {text:'Subfolder', to: {path: 'subfolder'}},
          {text:'Deep',to:{path:'deep'}},
          {text:'Deeper ellipsize in responsive mode'},
        ]
      }
    }
  }
</script>
```
</docs>
