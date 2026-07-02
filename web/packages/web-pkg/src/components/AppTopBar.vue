<template>
  <portal to="app.runtime.header.left">
    <div class="oc-app-top-bar oc-flex">
      <span class="oc-app-top-bar-inner oc-px-m oc-flex oc-flex-middle oc-flex-between">
        <div class="open-file-bar oc-flex">
          <resource-list-item
            v-if="resource"
            id="app-top-bar-resource"
            :is-thumbnail-displayed="false"
            :is-extension-displayed="areFileExtensionsShown"
            :path-prefix="getPathPrefix(resource)"
            :resource="resource"
            :parent-folder-name="getParentFolderName(resource)"
            :parent-folder-link-icon-additional-attributes="
              getParentFolderLinkIconAdditionalAttributes(resource)
            "
            :is-path-displayed="isPathDisplayed"
          />
        </div>
        <div class="oc-flex main-actions">
          <template v-if="dropDownMenuSections.length">
            <oc-button
              id="oc-openfile-contextmenu-trigger"
              v-oc-tooltip="contextMenuLabel"
              :aria-label="contextMenuLabel"
              appearance="raw-inverse"
              class="oc-p-xs"
              variation="brand"
            >
              <oc-icon name="more-2" />
            </oc-button>
            <oc-drop
              drop-id="oc-openfile-contextmenu"
              mode="click"
              padding-size="small"
              toggle="#oc-openfile-contextmenu-trigger"
              close-on-click
              focus-on-open
              @click.stop.prevent
            >
              <context-action-menu
                :menu-sections="dropDownMenuSections"
                :action-options="dropDownActionOptions"
              />
            </oc-drop>
          </template>
          <span v-if="hasAutosave" class="oc-flex oc-flex-middle">
            <oc-icon
              v-oc-tooltip="autoSaveTooltipText"
              :accessible-label="autoSaveTooltipText"
              name="refresh"
              color="white"
            />
          </span>
          <template v-if="mainActions.length && resource">
            <context-action-menu
              :menu-sections="[
                {
                  name: 'main-actions',
                  items: mainActions
                    .filter((action) => action.isVisible())
                    .map((action) => {
                      return { ...action, class: 'oc-p-xs', hideLabel: true }
                    })
                }
              ]"
              :action-options="{
                resources: [resource]
              }"
              appearance="raw-inverse"
              variation="brand"
            />
          </template>
          <oc-button
            id="app-top-bar-close"
            v-oc-tooltip="closeButtonLabel"
            appearance="raw-inverse"
            variation="brand"
            :aria-label="closeButtonLabel"
            @click="$emit('close')"
          >
            <oc-icon name="close" size="small" />
          </oc-button>
        </div>
      </span>
    </div>
  </portal>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import ContextActionMenu, { MenuSection } from './ContextActions/ContextActionMenu.vue'
import { useGettext } from 'vue3-gettext'
import {
  Action,
  FileActionOptions,
  useConfigStore,
  useFolderLink,
  useGetMatchingSpace,
  useResourcesStore
} from '../composables'
import ResourceListItem from './FilesList/ResourceListItem.vue'
import { isPublicSpaceResource, Resource } from '@ownclouders/web-client'
import { Duration } from 'luxon'

interface Props {
  dropDownMenuSections: MenuSection[]
  dropDownActionOptions: FileActionOptions
  mainActions: Action[]
  hasAutoSave: boolean
  isEditor: boolean
  resource: Resource
}
interface Emits {
  (e: 'close'): void
}

const {
  dropDownMenuSections = [],
  dropDownActionOptions = {
    space: null,
    resources: []
  },
  mainActions = [],
  hasAutoSave = true,
  isEditor = false,
  resource = null
} = defineProps<Partial<Props>>()

defineEmits<Emits>()

const { getPathPrefix, getParentFolderName, getParentFolderLinkIconAdditionalAttributes } =
  useFolderLink()
const { $gettext, current: currentLanguage } = useGettext()
const resourcesStore = useResourcesStore()
const configStore = useConfigStore()
const { getMatchingSpace } = useGetMatchingSpace()

const areFileExtensionsShown = computed(() => resourcesStore.areFileExtensionsShown)
const contextMenuLabel = computed(() => $gettext('Show context menu'))
const closeButtonLabel = computed(() => $gettext('Close'))
const hasAutosave = computed(
  () => isEditor && hasAutoSave && configStore.options.editor.autosaveEnabled
)
const autoSaveTooltipText = computed(() => {
  const duration = Duration.fromObject(
    { seconds: configStore.options.editor.autosaveInterval },
    { locale: currentLanguage }
  )
  return $gettext(`Autosave (every %{ duration })`, { duration: duration.toHuman() })
})

const space = computed(() => getMatchingSpace(resource))

const isPathDisplayed = computed(() => {
  return !isPublicSpaceResource(unref(space))
})
</script>

<style lang="scss">
.oc-app-top-bar {
  align-self: center;
  grid-column: 1 / 4;
  grid-row: secondRow;

  @media (min-width: $oc-breakpoint-small-default) {
    grid-column: 2;
    grid-row: 1;
  }
}

.oc-app-top-bar-inner {
  align-self: center;
  background-color: var(--oc-color-components-apptopbar-background);
  border-radius: 10px;
  border: 1px solid var(--oc-color-components-apptopbar-border);
  display: inline-flex;
  gap: 25px;
  height: 40px;
  margin: 10px auto;
  width: 100%;

  @media (min-width: $oc-breakpoint-small-default) {
    flex-basis: 250px;
    margin: 0;
  }

  .oc-resource-indicators {
    .text {
      color: var(--oc-color-swatch-brand-contrast);
    }
  }
}

.open-file-bar {
  #app-top-bar-resource {
    max-width: 360px;

    @media (max-width: $oc-breakpoint-medium-default) {
      max-width: 240px;
    }

    @media (min-width: $oc-breakpoint-small-default) {
      widows: initial;
    }

    svg,
    .oc-resource-name span {
      fill: var(--oc-color-swatch-inverse-default) !important;
      color: var(--oc-color-swatch-inverse-default) !important;
    }
  }

  .oc-resource-icon:hover,
  .oc-resource-name:hover {
    cursor: default;
    text-decoration: none;
  }
}
</style>
