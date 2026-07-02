<template>
  <div class="file_info oc-flex oc-flex-between oc-p-s">
    <div class="oc-flex oc-flex-middle">
      <resource-icon
        v-if="isSubPanelActive"
        :resource="resource"
        size="large"
        class="file_info__icon oc-mr-s oc-position-relative"
      />
      <div class="file_info__body oc-text-overflow">
        <!-- The following rule is false-negative as the resource-name component includes the necessary content -->
        <!-- eslint-disable-next-line vuejs-accessibility/heading-has-content -->
        <h3 data-testid="files-info-name" class="oc-font-semibold">
          <resource-name
            :name="resource.name"
            :extension="resource.extension"
            :type="resource.type"
            :full-path="resource.webDavPath"
            :is-extension-displayed="areFileExtensionsShown"
            :is-path-displayed="false"
            :truncate-name="false"
          />
        </h3>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, inject, unref } from 'vue'
import { HIDDEN_FILE_EXTENSIONS, Resource } from '@ownclouders/web-client'
import { useResourcesStore } from '../../../composables'
import ResourceIcon from '../../FilesList/ResourceIcon.vue'
import ResourceName from '../../FilesList/ResourceName.vue'

interface Props {
  isSubPanelActive?: boolean
}
const { isSubPanelActive = true } = defineProps<Props>()
const resourcesStore = useResourcesStore()

const resource = inject<Resource>('resource')
const areFileExtensionsShown = computed(
  () =>
    resourcesStore.areFileExtensionsShown &&
    !HIDDEN_FILE_EXTENSIONS.includes(unref(resource).extension)
)
</script>

<style lang="scss">
.file_info {
  &.sidebar-panel__file_info {
    border-bottom: 1px solid var(--oc-color-border);
  }

  button {
    white-space: nowrap;
  }

  &__body {
    text-align: left;

    h3 {
      font-size: var(--oc-font-size-medium);
      margin: 0;
      word-break: break-all;
    }
  }

  &__favorite {
    .oc-star {
      display: inline-block;

      &-shining svg {
        fill: #ffba0a !important;

        path:not([fill='none']) {
          stroke: var(--oc-color-swatch-passive-default);
        }
      }
    }
  }
}
</style>
