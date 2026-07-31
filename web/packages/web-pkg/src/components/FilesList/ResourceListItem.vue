<template>
  <div
    ref="resourceListItem"
    class="oc-resource oc-text-overflow"
    :class="{ 'oc-resource-no-interaction': !isResourceClickable }"
  >
    <span v-if="isIconDisplayed" class="oc-resource-icon-wrapper">
      <oc-img
        v-if="hasThumbnail"
        :key="thumbnail"
        v-oc-tooltip="tooltipLabelIcon"
        :src="thumbnail"
        class="oc-resource-thumbnail"
        width="40"
        height="40"
        :aria-label="tooltipLabelIcon"
      />
      <resource-icon
        v-else
        v-oc-tooltip="tooltipLabelIcon"
        :aria-label="tooltipLabelIcon"
        :aria-hidden="tooltipLabelIcon === null"
        :resource="resource"
        :role="tooltipLabelIcon ? 'img' : 'presentation'"
      />
    </span>
    <div class="oc-resource-details oc-text-overflow" :class="{ 'oc-pl-s': isIconDisplayed }">
      <resource-link
        :resource="resource"
        :is-resource-clickable="isResourceClickable"
        :link="link"
        class="oc-text-overflow"
        @click="emitClick"
      >
        <resource-name
          :key="resource.name"
          :name="resource.name"
          :extension="resource.extension"
          :type="resource.type"
          :full-path="resource.path"
          :is-path-displayed="isPathDisplayed"
          :is-extension-displayed="
            isExtensionDisplayed && !HIDDEN_EXTENSIONS.includes(resource.extension)
          "
        />
      </resource-link>
      <div class="oc-resource-indicators">
        <component
          :is="parentFolderComponentType"
          v-if="isPathDisplayed"
          v-oc-tooltip="parentFolderPathTooltip"
          :to="parentFolderLink"
          :style="parentFolderStyle"
          class="parent-folder oc-text-truncate"
          :aria-current="isSearchResult ? null : 'page'"
        >
          <oc-icon v-bind="parentFolderLinkIconAttrs" />
          <span class="text oc-text-truncate" v-text="parentFolderName" />
        </component>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { computed, useTemplateRef } from 'vue'
import { useGettext } from 'vue3-gettext'
import { HIDDEN_FILE_EXTENSIONS, SpaceResource, Resource } from '@ownclouders/web-client'
import ResourceIcon from './ResourceIcon.vue'
import ResourceLink from './ResourceLink.vue'
import ResourceName from './ResourceName.vue'
import { RouteLocationRaw } from 'vue-router'
import { dirname, join } from 'node:path'

interface Props {
  resource: SpaceResource | Resource
  pathPrefix?: string
  link?: RouteLocationRaw | null
  isPathDisplayed?: boolean
  parentFolderLink?: RouteLocationRaw | null
  parentFolderName?: string
  parentFolderLinkIconAdditionalAttributes?: Record<string, any>
  isExtensionDisplayed?: boolean
  isThumbnailDisplayed?: boolean
  isIconDisplayed?: boolean
  isResourceClickable?: boolean
  isSearchResult?: boolean
}

interface Emits {
  (e: 'click'): void
}
const {
  resource,
  pathPrefix = '',
  link = null,
  isPathDisplayed = false,
  parentFolderLink = null,
  parentFolderName = '',
  parentFolderLinkIconAdditionalAttributes = {},
  isExtensionDisplayed = true,
  isThumbnailDisplayed = true,
  isIconDisplayed = true,
  isResourceClickable = true,
  isSearchResult = false
} = defineProps<Props>()

const resourceListItem = useTemplateRef('resourceListItem')
const emit = defineEmits<Emits>()
const { $gettext } = useGettext()
const HIDDEN_EXTENSIONS = HIDDEN_FILE_EXTENSIONS

defineExpose({ resourceListItem })

/**
 * Checks if a path has a valid second segment after splitting by '/'.
 * Used for tooltip display - if true, shows "segment1 > segment2 > ...",
 * otherwise just shows "segment1" without trailing ">".
 *
 * @param {string} value - The path to check
 * @returns {boolean} - True if path has a valid second segment, false otherwise
 */

function splitPathHasSecondSegment(value: string): boolean {
  const result = value.split('/')
  return result?.length > 1 && result[1].length > 0
}

function emitClick() {
  /**
   * Triggered when the resource is a file and the name is clicked
   */
  emit('click')
}

const parentFolderPathTooltip = computed(() => {
  if (!isPathDisplayed) {
    return null
  }

  const parentFolderPath = dirname(resource.path)

  if (pathPrefix) {
    if (splitPathHasSecondSegment(parentFolderPath)) {
      return join(pathPrefix, parentFolderPath).replaceAll('/', ' > ')
    }
    return pathPrefix.replace('/', ' > ')
  }

  return parentFolderPath.replaceAll('/', ' > ')
})
const parentFolderComponentType = computed(() => {
  return parentFolderLink ? 'router-link' : 'span'
})

const parentFolderStyle = computed(() => {
  return {
    cursor: parentFolderLink ? 'pointer' : 'default'
  }
})

const parentFolderLinkIconAttrs = computed(() => {
  return {
    'fill-type': 'line',
    name: 'folder-2',
    size: 'small',
    ...parentFolderLinkIconAdditionalAttributes
  }
})

const hasThumbnail = computed(() => {
  return isThumbnailDisplayed && Object.prototype.hasOwnProperty.call(resource, 'thumbnail')
})

const thumbnail = computed(() => {
  return resource.thumbnail
})

const tooltipLabelIcon = computed(() => {
  if (resource.locked) {
    return $gettext('This item is locked')
  }
  return null
})
</script>

<style lang="scss">
.oc-resource {
  align-items: center;
  display: inline-flex;
  justify-content: flex-start;
  overflow: visible !important;

  &-no-interaction {
    pointer-events: none;
  }

  &-icon-link {
    position: relative;
  }

  &-thumbnail {
    border-radius: 2px;
    object-fit: cover;
    height: $oc-size-icon-default * 1.5;
    max-height: $oc-size-icon-default * 1.5;
    width: $oc-size-icon-default * 1.5;
    max-width: $oc-size-icon-default * 1.5;
  }

  &-details {
    display: block;

    a {
      text-decoration: none;
    }

    a:hover,
    a:focus {
      outline-offset: 0;
    }
  }

  &-indicators {
    display: flex;

    a {
      &:hover {
        background-color: var(--oc-color-input-bg);
        border-radius: 2px;
      }

      .text {
        &:hover {
          color: var(--oc-color-text-default);
          text-decoration: underline;
        }
      }
    }

    .parent-folder {
      display: flex;
      align-items: center;

      padding: 0 2px 0 2px;
      margin: 0 8px 0 -2px;

      .oc-icon {
        padding-right: 3px;
      }

      .text {
        font-size: 0.8125rem;
        color: var(--oc-color-text-muted);
      }
    }
  }
}
</style>
