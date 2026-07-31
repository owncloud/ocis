<template>
  <!-- eslint-disable-next-line vuejs-accessibility/no-static-element-interactions -->
  <div
    ref="observerTarget"
    class="oc-tile-card oc-card oc-card-default oc-rounded"
    :data-item-id="resource.id"
    :class="{
      'oc-tile-card-selected': isResourceSelected,
      'oc-tile-card-disabled': isResourceDisabled && !isProjectSpaceResource(resource),
      'state-trashed': isResourceDisabled && isProjectSpaceResource(resource)
    }"
    @contextmenu="$emit('contextmenu', $event)"
  >
    <div v-if="isHidden" class="oc-tile-card-lazy-shimmer"></div>
    <template v-else>
      <resource-link
        class="oc-card-media-top oc-flex oc-flex-center oc-flex-middle oc-m-rm"
        :resource="resource"
        :link="resourceRoute"
        :is-resource-clickable="isResourceClickable"
        tabindex="-1"
        @click="$emit('click')"
      >
        <div class="oc-tile-card-selection">
          <slot name="selection" :item="resource" />
        </div>
        <oc-tag
          v-if="isResourceDisabled && isProjectSpaceResource(resource)"
          class="resource-disabled-indicator oc-position-absolute"
          type="span"
        >
          <span v-text="$gettext('Disabled')" />
        </oc-tag>
        <div
          v-oc-tooltip="tooltipLabelIcon"
          class="oc-tile-card-preview oc-flex oc-flex-middle oc-flex-center"
          :aria-label="tooltipLabelIcon"
        >
          <div class="oc-tile-card-hover"></div>
          <slot name="imageField" :item="resource">
            <oc-img
              v-if="shouldDisplayThumbnails(resource)"
              class="tile-preview"
              :src="resource.thumbnail"
            />
            <resource-icon
              v-else
              :resource="resource"
              :size="resourceIconSize"
              class="tile-default-image oc-pt-xs"
            >
              <template v-if="showStatusIcon" #status>
                <oc-icon v-bind="statusIconAttrs" size="xsmall" />
              </template>
            </resource-icon>
          </slot>
        </div>
      </resource-link>
      <div class="oc-card-body oc-p-s">
        <div class="oc-flex oc-flex-between oc-flex-middle">
          <div class="oc-flex oc-flex-middle oc-text-truncate resource-name-wrapper">
            <resource-list-item
              :resource="resource"
              :is-icon-displayed="false"
              :is-extension-displayed="isExtensionDisplayed"
              :is-resource-clickable="isResourceClickable"
              :link="resourceRoute"
              @click="$emit('click')"
            />
          </div>
          <div class="oc-flex oc-flex-middle">
            <!-- Slot for indicators !-->
            <slot name="indicators" :item="resource" class="resource-indicators" />
            <!-- Slot for individual actions -->
            <slot name="actions" :item="resource" />
            <!-- Slot for contextmenu -->
            <slot name="contextMenu" :item="resource" />
          </div>
        </div>
        <p v-if="resourceDescription" class="oc-text-left oc-my-rm oc-text-truncate">
          <small v-text="resourceDescription" />
        </p>
      </div>
    </template>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import ResourceIcon from './ResourceIcon.vue'
import ResourceListItem from './ResourceListItem.vue'
import ResourceLink from './ResourceLink.vue'
import { isProjectSpaceResource, SpaceResource, Resource } from '@ownclouders/web-client'
import { useGettext } from 'vue3-gettext'
import { isSpaceResource } from '@ownclouders/web-client'
import { RouteLocationRaw } from 'vue-router'
import { useIsVisible } from '@ownclouders/design-system/composables'
import { customRef, ref, unref } from 'vue'

interface Props {
  resource: SpaceResource | Resource
  resourceRoute?: RouteLocationRaw | null
  isResourceSelected?: boolean
  isResourceClickable?: boolean
  isResourceDisabled?: boolean
  isExtensionDisplayed?: boolean
  resourceIconSize?: 'large' | 'xlarge' | 'xxlarge' | 'xxxlarge'
  lazy?: boolean
}
interface Emits {
  (e: 'contextmenu', event: MouseEvent): void
  (e: 'click'): void
  (e: 'itemVisible'): void
}

const {
  resource,
  resourceRoute = null,
  isResourceSelected = false,
  isResourceClickable = true,
  isResourceDisabled = false,
  isExtensionDisplayed = true,
  resourceIconSize = 'xlarge',
  lazy = false
} = defineProps<Props>()
const emit = defineEmits<Emits>()
const { $gettext } = useGettext()

const observerTarget = customRef((track, trigger) => {
  let $el: HTMLElement
  return {
    get() {
      track()
      return $el
    },
    set(value) {
      $el = value
      trigger()
    }
  }
})

const showStatusIcon = computed(() => {
  return resource.locked || resource.processing
})

const statusIconAttrs = computed(() => {
  if (resource.locked) {
    return {
      name: 'lock',
      fillType: 'fill'
    }
  }

  if (resource.processing) {
    return {
      name: 'loop-right',
      fillType: 'line'
    }
  }

  return {}
})

const tooltipLabelIcon = computed(() => {
  if (resource.locked) {
    return $gettext('This item is locked')
  }
  return null
})
const resourceDescription = computed(() => {
  if (isSpaceResource(resource)) {
    return resource.description
  }
  return ''
})

const shouldDisplayThumbnails = (resource: SpaceResource | Resource) => {
  return resource.thumbnail
}

const { isVisible } = lazy
  ? useIsVisible({
      target: observerTarget,
      onVisibleCallback: () => emit('itemVisible')
    })
  : { isVisible: ref(true) }

const isHidden = computed(() => !unref(isVisible))

if (!lazy) {
  emit('itemVisible')
}
</script>

<style lang="scss">
.oc-tile-card {
  background-color: var(--oc-color-background-highlight) !important;
  box-shadow: none;
  height: 100%;
  display: flex;
  flex-flow: column;
  outline: 1px solid var(--oc-color-border);

  &-disabled {
    pointer-events: none;
    background-color: var(--oc-color-background-muted) !important;
    opacity: 0.7;
    filter: grayscale(0.6);

    // Show tooltip on status indicators without handler
    span.oc-status-indicators-indicator {
      pointer-events: all;
    }
  }

  &.state-trashed {
    .tile-image,
    .tile-default-image > svg {
      filter: grayscale(100%);
      opacity: 80%;
    }
  }

  .tile-default-image {
    position: relative;
  }

  .oc-card-media-top {
    position: relative;
    aspect-ratio: 16/9;
    justify-content: center;
    width: 100%;

    .oc-tag {
      color: var(--oc-color-text-default);

      &.resource-disabled-indicator {
        z-index: 1;
      }
    }

    .tile-preview {
      aspect-ratio: 16/9;
      height: 100%;
      object-fit: cover;
      width: 100%;
      border-top-left-radius: 5px;
      border-top-right-radius: 5px;
    }

    &:hover {
      .oc-tile-card-hover {
        opacity: 15%;
      }
    }
  }

  &-selected {
    outline: 2px solid var(--oc-color-swatch-primary-hover);

    .oc-tile-card-preview {
      width: calc(100% - var(--oc-space-medium));
      height: calc(100% - var(--oc-space-medium));

      .tile-preview,
      .oc-tile-card-hover {
        border-radius: 5px !important;
      }

      .oc-tile-card-hover {
        opacity: 10%;
      }
    }
  }

  &-selection {
    z-index: 1;
    position: absolute;
    top: 0;
    left: 0;

    input {
      background-color: var(--oc-color-background-muted);
    }

    input.oc-checkbox-checked {
      background-color: var(--oc-color-swatch-inverse-default);
    }
  }

  &-preview {
    position: absolute;
    height: 100%;
    width: 100%;
    text-align: center;
  }

  &-hover {
    position: absolute;
    width: 100%;
    height: 100%;
    background: #000;
    opacity: 0;
    border-top-left-radius: 5px;
    border-top-right-radius: 5px;
  }

  .resource-name-wrapper {
    color: var(--oc-color-text-default);
    max-width: 70%;
    overflow: hidden;
  }

  &-lazy-shimmer {
    height: 120px;
    opacity: 0.2;
    position: relative;
    overflow: hidden;
  }

  &-lazy-shimmer::after {
    animation: shimmer 2s infinite;
    background-image: linear-gradient(
      90deg,
      rgba(#4c5f79, 0) 0,
      rgba(#4c5f79, 0.2) 20%,
      rgba(#4c5f79, 0.5) 60%,
      rgba(#4c5f79, 0)
    );
    bottom: 0;
    content: '';
    left: 0;
    position: absolute;
    right: 0;
    top: 0;
    transform: translateX(-100%);
  }

  @keyframes shimmer {
    100% {
      transform: translateX(100%);
    }
  }
}
</style>
