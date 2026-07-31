<template>
  <div
    class="space-header oc-p-m"
    :class="{ 'oc-flex': !imageExpanded && !isMobileWidth, 'space-header-squashed': isSideBarOpen }"
  >
    <div
      class="space-header-image"
      :class="{ 'space-header-image-expanded': imageExpanded || isMobileWidth }"
    >
      <button v-if="imageContent" class="btn-toggle-image" @click="toggleImageExpanded">
        <label class="oc-invisible-sr" for="btn-toggle-image">
          {{
            imageExpanded
              ? $pgettext(
                  'Accessibility label to inform user when space image has been collapsed by being clicked on',
                  'Collapse space image'
                )
              : $pgettext(
                  'Accessibility label to inform user when space image has been expanded by being clicked on',
                  'Expand space image'
                )
          }}
        </label>
        <img alt="" :src="imageContent" />
      </button>
      <div v-else class="space-header-image-default oc-flex oc-flex-middle oc-flex-center">
        <oc-icon name="layout-grid" size="xxlarge" class="oc-px-m oc-py-m" />
      </div>
    </div>
    <div class="space-header-infos">
      <div class="oc-flex oc-mb-s oc-flex-middle oc-flex-between">
        <div class="oc-flex oc-flex-middle space-header-infos-heading">
          <h2 class="space-header-name">{{ space.name }}</h2>
          <oc-button
            :id="`space-context-btn`"
            v-oc-tooltip="$gettext('Show context menu')"
            :aria-label="$gettext('Show context menu')"
            appearance="raw"
            class="oc-ml-s"
          >
            <oc-icon name="more-2" />
          </oc-button>
          <oc-drop
            :drop-id="`space-context-drop`"
            :toggle="`#space-context-btn`"
            mode="click"
            close-on-click
            :options="{ delayHide: 0 }"
            padding-size="small"
            position="right-start"
          >
            <space-context-actions :action-options="{ resources: [space] }" />
          </oc-drop>
        </div>
        <oc-button
          v-if="memberCount"
          :aria-label="$gettext('Open context menu and show members')"
          appearance="raw"
          @click="openSideBarSharePanel"
        >
          <oc-icon name="group" fill-type="line" size="small" />
          <span class="space-header-people-count oc-text-small" v-text="memberCountString"></span>
        </oc-button>
      </div>
      <p v-if="space.description" class="oc-mt-rm oc-text-bold">{{ space.description }}</p>
      <div
        v-if="markdownResource && markdownContent"
        ref="markdownContainerRef"
        class="markdown-container oc-flex oc-flex-middle"
      >
        <text-editor is-read-only :current-content="markdownContent" />

        <div class="markdown-container-edit oc-ml-s">
          <router-link
            v-oc-tooltip="$gettext('Edit description')"
            size="small"
            appearance="raw"
            :to="editReadMeContentLink"
          >
            <oc-icon name="pencil" size="small" fill-type="line" />
            <span class="oc-invisible-sr">
              {{ $gettext('Edit description') }}
            </span>
          </router-link>
        </div>
      </div>
      <div
        v-if="showMarkdownCollapse && markdownContent"
        class="markdown-collapse oc-text-center oc-mt-s"
      >
        <oc-button appearance="raw" @click="toggleMarkdownCollapsed">
          <oc-icon :name="toggleMarkdownCollapsedIcon" />
          <span>{{ toggleMarkdownCollapsedText }}</span>
        </oc-button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, inject, onBeforeUnmount, onMounted, Ref, ref, unref, watch } from 'vue'
import { buildSpaceImageResource, SpaceResource } from '@ownclouders/web-client'
import {
  useClientService,
  ProcessorType,
  useResourcesStore,
  useFileActions,
  useLoadPreview,
  TextEditor
} from '@ownclouders/web-pkg'
import { ImageDimension } from '@ownclouders/web-pkg'
import SpaceContextActions from './SpaceContextActions.vue'
import { eventBus } from '@ownclouders/web-pkg'
import { SideBarEventTopics } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { DriveItem } from '@ownclouders/web-client/graph/generated'
interface Props {
  space: SpaceResource
  isSideBarOpen: boolean
}
const markdownContainerCollapsedClass = 'collapsed'
const { space, isSideBarOpen = false } = defineProps<Props>()

const language = useGettext()
const { $gettext, $ngettext } = language
const clientService = useClientService()
const { getFileContents, getFileInfo } = clientService.webdav
const resourcesStore = useResourcesStore()
const { getDefaultAction } = useFileActions()
const { loadPreview } = useLoadPreview()

const markdownContainerRef = ref(null)
const markdownContent = ref('')
const markdownResource = ref(null)
const markdownCollapsed = ref(true)
const showMarkdownCollapse = ref(false)
const toggleMarkdownCollapsedIcon = computed(() => {
  return unref(markdownCollapsed) ? 'add' : 'subtract'
})
const toggleMarkdownCollapsedText = computed(() => {
  return unref(markdownCollapsed) ? $gettext('Show more') : $gettext('Show less')
})
const toggleMarkdownCollapsed = () => {
  markdownCollapsed.value = !unref(markdownCollapsed)
  unref(markdownContainerRef).classList.toggle(markdownContainerCollapsedClass)
}
const onMarkdownResize = () => {
  if (!unref(markdownContainerRef)) {
    return
  }

  unref(markdownContainerRef).classList.remove(markdownContainerCollapsedClass)
  const markdownContainerHeight = unref(markdownContainerRef).offsetHeight
  if (markdownContainerHeight < 150) {
    showMarkdownCollapse.value = false
    return
  }
  showMarkdownCollapse.value = true

  if (unref(markdownCollapsed)) {
    unref(markdownContainerRef).classList.add(markdownContainerCollapsedClass)
  }
}
const markdownResizeObserver = new ResizeObserver(onMarkdownResize)
const observeMarkdownContainerResize = () => {
  if (!markdownResizeObserver || !unref(markdownContainerRef)) {
    return
  }
  markdownResizeObserver.unobserve(unref(markdownContainerRef))
  markdownResizeObserver.observe(unref(markdownContainerRef))
}
const unobserveMarkdownContainerResize = () => {
  if (!markdownResizeObserver || !unref(markdownContainerRef)) {
    return
  }
  markdownResizeObserver.unobserve(unref(markdownContainerRef))
}
onMounted(observeMarkdownContainerResize)
onBeforeUnmount(() => {
  unobserveMarkdownContainerResize()
})
watch(
  computed(() => space.spaceReadmeData),
  async (data: DriveItem) => {
    if (!data) {
      return
    }

    const fileContentsResponse = await getFileContents(space, {
      path: `.space/${space.spaceReadmeData.name}`
    })

    const fileInfoResponse = await getFileInfo(space, {
      path: `.space/${space.spaceReadmeData.name}`
    })

    unobserveMarkdownContainerResize()
    markdownContent.value = fileContentsResponse.body
    markdownResource.value = fileInfoResponse

    if (unref(markdownContent)) {
      observeMarkdownContainerResize()
    }
  },
  { deep: true, immediate: true }
)

const imageContent = ref<string>(null)
const imageExpanded = ref(false)

const editReadMeContentLink = computed(() => {
  const action = getDefaultAction({ resources: [unref(markdownResource)], space })

  if (!action.route) {
    return null
  }

  return action.route({ space: space, resources: [unref(markdownResource)] })
})
const toggleImageExpanded = () => {
  imageExpanded.value = !unref(imageExpanded)
}

watch(
  computed(() => space.spaceImageData),
  async (data) => {
    if (!data) {
      return
    }
    const resource = buildSpaceImageResource(space)
    imageContent.value = await loadPreview({
      space,
      resource,
      dimensions: ImageDimension.Tile,
      processor: ProcessorType.enum.fit,
      cancelRunning: true,
      updateStore: false
    })
  },
  { immediate: true }
)

const memberCount = computed(() => {
  return Object.keys(space.members).length
})
const memberCountString = computed(() => {
  return $ngettext('%{count} member', '%{count} members', unref(memberCount), {
    count: unref(memberCount).toString()
  })
})

const openSideBarSharePanel = () => {
  resourcesStore.setSelection([])
  eventBus.publish(SideBarEventTopics.openWithPanel, 'space-share')
}

const isMobileWidth = inject<Ref<boolean>>('isMobileWidth')
</script>

<style lang="scss">
.space-header {
  &-squashed {
    .space-header-image {
      @media only screen and (max-width: 1200px) {
        display: none;
      }
    }
  }

  &-image {
    width: 280px;
    min-width: 280px;
    aspect-ratio: 16 / 9;
    margin-right: var(--oc-space-large);
    max-height: 158px;

    &-default {
      background-color: var(--oc-color-background-highlight);
      height: 100%;
      border-radius: 10px;
    }

    &-expanded {
      width: 100%;
      margin: 0;
      max-height: 100%;
      max-width: 100%;
    }

    img {
      border-radius: 10px;
      height: 100%;
      width: 100%;
      max-height: 100%;
      object-fit: cover;
    }

    .btn-toggle-image {
      background: transparent;
      border: none;
      cursor: pointer;
      outline: none;
      height: 100%;
      width: 100%;
    }
  }

  &-infos {
    flex: 1;

    &-heading {
      max-width: 100%;
    }
  }

  &-name {
    font-size: 1.5rem;
    word-break: break-all;
  }

  &-people-count {
    white-space: nowrap;
  }

  .markdown-container.collapsed {
    max-height: 100px;
    overflow: hidden;
    -webkit-mask-image: linear-gradient(180deg, #000 90%, transparent);
  }
}
</style>
