<template>
  <div v-if="isFolderLoading" class="oc-width-1-1">
    <div class="oc-position-center">
      <oc-spinner :aria-label="$gettext('Loading media file')" size="xlarge" />
    </div>
  </div>
  <!-- eslint-disable-next-line vuejs-accessibility/no-static-element-interactions -->
  <div
    v-else
    ref="preview"
    class="oc-flex oc-width-1-1 oc-height-1-1"
    tabindex="-1"
    @keydown.left="goToPrev"
    @keydown.right="goToNext"
  >
    <div class="stage" :class="{ lightbox: isFullScreenModeActivated }">
      <div class="stage_media">
        <div v-if="!activeMediaFileCached || activeMediaFileCached.isLoading" class="oc-width-1-1">
          <div class="oc-position-center">
            <oc-spinner :aria-label="$gettext('Loading media file')" size="xlarge" />
          </div>
        </div>
        <div
          v-else-if="activeMediaFileCached.isError"
          class="oc-width-1-1 oc-flex oc-flex-column oc-flex-middle oc-flex-center"
        >
          <oc-icon name="file-damage" variation="danger" size="xlarge" />
          <p>
            {{ $gettext('Failed to load "%{filename}"', { filename: activeMediaFileCached.name }) }}
          </p>
        </div>
        <media-image
          v-else-if="activeMediaFileCached.isImage"
          :file="activeMediaFileCached"
          :current-image-rotation="currentImageRotation"
          :current-image-zoom="currentImageZoom"
          :current-image-position-x="currentImagePositionX"
          :current-image-position-y="currentImagePositionY"
          @pan-zoom-change="onPanZoomChanged"
        />
        <media-video
          v-else-if="activeMediaFileCached.isVideo"
          :file="activeMediaFileCached"
          :is-auto-play-enabled="isAutoPlayEnabled"
        />
        <media-audio
          v-else-if="activeMediaFileCached.isAudio"
          :file="activeMediaFileCached"
          :resource="activeFilteredFile"
          :is-auto-play-enabled="isAutoPlayEnabled"
        />
      </div>
      <media-controls
        class="stage_controls"
        :files="filteredFiles"
        :active-index="activeIndex"
        :is-full-screen-mode-activated="isFullScreenModeActivated"
        :is-folder-loading="isFolderLoading"
        :show-image-controls="activeMediaFileCached?.isImage && !activeMediaFileCached?.isError"
        :current-image-rotation="currentImageRotation"
        :current-image-zoom="currentImageZoom"
        @set-rotation="currentImageRotation = $event"
        @set-zoom="currentImageZoom = $event"
        @reset-image="resetImage"
        @toggle-full-screen="toggleFullScreenMode"
        @toggle-previous="goToPrev"
        @toggle-next="goToNext"
      />
    </div>
  </div>
</template>
<script lang="ts" setup>
import { computed, ref, unref, nextTick, watch, Ref, onMounted, onBeforeMount } from 'vue'
import { IncomingShareResource, Resource } from '@ownclouders/web-client'
import {
  AppFolderHandlingResult,
  FileContext,
  ProcessorType,
  SortDir,
  createFileRouteOptions,
  queryItemAsString,
  sortHelper,
  useRoute,
  useRouteQuery,
  useRouter,
  usePreviewService,
  useAppStore,
  useGetMatchingSpace,
  isLocationSharesActive
} from '@ownclouders/web-pkg'
import MediaControls from './components/MediaControls.vue'
import MediaAudio from './components/Sources/MediaAudio.vue'
import MediaImage from './components/Sources/MediaImage.vue'
import MediaVideo from './components/Sources/MediaVideo.vue'
import { CachedFile } from './helpers/types'
import {
  useFileTypes,
  useFullScreenMode,
  useImageControls,
  usePreviewDimensions
} from './composables'
import { mimeTypes } from './mimeTypes'
import { RouteLocationRaw } from 'vue-router'

const PRELOAD_COUNT = 5

interface Props {
  activeFiles: Resource[]
  currentFileContext: FileContext
  loadFolderForFileContext: (fileContext: FileContext) => Promise<AppFolderHandlingResult>
  getUrlForResource: (space: Resource, resource: Resource) => Promise<string>
  revokeUrl: (url: string | undefined) => void
  isFolderLoading: boolean
}
interface Emits {
  (e: 'update:resource', resource: Resource | undefined): void
}
const {
  activeFiles,
  currentFileContext,
  loadFolderForFileContext,
  getUrlForResource,
  revokeUrl,
  isFolderLoading
} = defineProps<Props>()
const emit = defineEmits<Emits>()
const router = useRouter()
const route = useRoute()
const contextRouteQuery = useRouteQuery('contextRouteQuery') as unknown as Ref<
  Record<string, string>
>

const { isFileTypeAudio, isFileTypeImage, isFileTypeVideo } = useFileTypes()
const previewService = usePreviewService()
const { dimensions } = usePreviewDimensions()
const { getMatchingSpace } = useGetMatchingSpace()
const {
  currentImageZoom,
  currentImageRotation,
  currentImagePositionX,
  currentImagePositionY,
  resetImage,
  onPanZoomChanged
} = useImageControls()

const { isFullScreenModeActivated, toggleFullScreenMode } = useFullScreenMode()

const activeIndex = ref<number>()
const cachedFiles = ref<Record<string, CachedFile>>({})
const folderLoaded = ref(false)
const isAutoPlayEnabled = ref(true)
const preview = ref<HTMLElement>()
const appStore = useAppStore()

const space = computed(() => {
  return getMatchingSpace(unref(activeFilteredFile))
})

const sortBy = computed(() => {
  if (!unref(contextRouteQuery)) {
    return 'name'
  }
  return unref(contextRouteQuery)['sort-by'] ?? 'name'
})
const sortDir = computed<SortDir>(() => {
  if (!unref(contextRouteQuery)) {
    return SortDir.Desc
  }
  return (unref(contextRouteQuery)['sort-dir'] as SortDir) ?? SortDir.Asc
})

const fileIdQuery = useRouteQuery('fileId')
const fileId = computed(() => queryItemAsString(unref(fileIdQuery)))

const filteredFiles = computed(() => {
  if (!activeFiles) {
    return []
  }

  const files = activeFiles.filter((file) => {
    if (
      unref(currentFileContext.routeQuery)?.['q_share-visibility'] === 'hidden' &&
      !(file as IncomingShareResource).hidden
    ) {
      return false
    }

    if (
      unref(currentFileContext.routeQuery)?.['q_share-visibility'] !== 'hidden' &&
      (file as IncomingShareResource).hidden
    ) {
      return false
    }

    return mimeTypes.includes(file.mimeType?.toLowerCase()) && file.canDownload()
  })

  return sortHelper(files, [{ name: unref(sortBy) }], unref(sortBy), unref(sortDir))
})
const activeFilteredFile = computed(() => {
  return unref(filteredFiles)[unref(activeIndex)]
})
const activeMediaFileCached = computed(() => {
  return unref(cachedFiles)[unref(activeFilteredFile)?.id]
})

const loadFileIntoCache = async (file: Resource) => {
  if (Object.hasOwn(unref(cachedFiles), file.id)) {
    return
  }

  const cachedFile: CachedFile = {
    id: file.id,
    name: file.name,
    url: undefined,
    ext: file.extension,
    mimeType: file.mimeType,
    isVideo: isFileTypeVideo(file),
    isImage: isFileTypeImage(file),
    isAudio: isFileTypeAudio(file),
    isLoading: ref(true),
    isError: ref(false)
  }
  cachedFiles.value[file.id] = cachedFile

  try {
    if (cachedFile.isImage) {
      cachedFile.url = await previewService.loadPreview(
        {
          space: unref(space),
          resource: file,
          dimensions: unref(dimensions),
          processor: ProcessorType.enum.fit
        },
        false,
        false
      )
      return
    }
    cachedFile.url = await getUrlForResource(unref(space), file)
  } catch (e) {
    console.error(e)
    cachedFile.isError.value = true
  } finally {
    cachedFile.isLoading.value = false
  }
}

const updateLocalHistory = () => {
  // this is a rare edge case when browsing quickly through a lot of files
  // we workaround context being null, when useDriveResolver is in loading state
  if (!currentFileContext) {
    return
  }

  const { params, query } = createFileRouteOptions(unref(space), unref(activeFilteredFile))
  const { fullPath, ...routeWithoutFullPath } = unref(route)

  router.replace({
    ...routeWithoutFullPath,
    path: fullPath,
    params: { ...routeWithoutFullPath.params, ...params },
    query: { ...routeWithoutFullPath.query, ...query }
  } as RouteLocationRaw)
}

watch(
  () => currentFileContext,
  async () => {
    if (!currentFileContext) {
      return
    }

    if (!unref(folderLoaded)) {
      try {
        await loadFolderForFileContext(currentFileContext)
        folderLoaded.value = true
      } catch (e) {
        appStore.error = e
      }
    }

    setActiveFile()
  },
  { immediate: true, deep: true }
)

watch(
  () => activeFilteredFile.value,
  (file) => {
    emit('update:resource', unref(file))
  }
)

const loading = computed(() => {
  if (isFolderLoading) {
    return true
  }
  const file = unref(activeMediaFileCached)
  if (!file) {
    return true
  }
  return unref(file.isLoading)
})
watch(
  () => loading.value,
  async (loadingState) => {
    if (!loadingState) {
      await nextTick()
      unref(preview).focus()
    }
  },
  { immediate: true }
)

function setActiveFile() {
  for (let i = 0; i < unref(filteredFiles).length; i++) {
    const filterAttr = isLocationSharesActive(router, 'files-shares-with-me')
      ? 'remoteItemId'
      : 'fileId'

    // match the given file id with the filtered files to get the current index
    if (unref(filteredFiles)[i][filterAttr] === unref(fileId)) {
      activeIndex.value = i
      return
    }

    activeIndex.value = 0
  }
}
// react to PopStateEvent ()
function handleLocalHistoryEvent() {
  setActiveFile()
}
function goToNext() {
  if (unref(activeIndex) + 1 >= unref(filteredFiles).length) {
    activeIndex.value = 0
    updateLocalHistory()
    return
  }
  activeIndex.value++
  updateLocalHistory()
}
function goToPrev() {
  if (unref(activeIndex) === 0) {
    activeIndex.value = unref(filteredFiles).length - 1
    updateLocalHistory()
    return
  }
  activeIndex.value--
  updateLocalHistory()
}
function preloadImages() {
  const preloadFile = (preloadFileIndex: number) => {
    const cycleIndex =
      (((unref(activeIndex) + preloadFileIndex) % unref(filteredFiles).length) +
        unref(filteredFiles).length) %
      unref(filteredFiles).length

    const file = unref(filteredFiles)[cycleIndex]
    loadFileIntoCache(file)
  }

  for (let followingFileIndex = 1; followingFileIndex <= PRELOAD_COUNT; followingFileIndex++) {
    preloadFile(followingFileIndex)
  }

  for (let previousFileIndex = -1; previousFileIndex >= PRELOAD_COUNT * -1; previousFileIndex--) {
    preloadFile(previousFileIndex)
  }
}

watch(activeIndex, (newValue, oldValue) => {
  if (newValue !== oldValue) {
    loadFileIntoCache(unref(activeFilteredFile))
    preloadImages()
  }

  if (oldValue !== null) {
    isAutoPlayEnabled.value = false
  }

  currentImageZoom.value = 1
  currentImageRotation.value = 0
})
onMounted(() => {
  // keep a local history for this component
  window.addEventListener('popstate', handleLocalHistoryEvent)
})
onBeforeMount(() => {
  window.removeEventListener('popstate', handleLocalHistoryEvent)

  Object.values(unref(cachedFiles)).forEach((cachedFile) => {
    revokeUrl(unref(cachedFile).url)
  })
})
</script>

<style lang="scss" scoped>
.stage {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  text-align: center;

  &_media {
    flex-grow: 1;
    overflow: hidden;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  &_controls {
    height: auto;
    margin: 10px auto;
  }
}
</style>
