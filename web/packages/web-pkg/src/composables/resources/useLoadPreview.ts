import PQueue from 'p-queue'
import { computed, onUnmounted, Ref, unref } from 'vue'
import { useTask } from 'vue-concurrency'
import {
  buildSpaceImageResource,
  isProjectSpaceResource,
  Resource,
  SpaceResource
} from '@ownclouders/web-client'
import { FolderViewModeConstants } from '../viewMode'
import { usePreviewService } from '../previewService'
import { ProcessorType } from '../../services'
import { useResourcesStore } from '../piniaStores'
import { ImageDimension } from '../../constants'

type LoadPreviewOptions = {
  space: SpaceResource
  resource: Resource
  dimensions?: [number, number]
  processor?: ProcessorType

  /**
   * Cancel potential running tasks before loading.
   * Recommended when loading one preview at a time (hence not in a file list).
   * @default false
   */
  cancelRunning?: boolean

  /**
   * Update resource store after loading.
   * Recommended when loading previews in a file list.
   * @default true
   */
  updateStore?: boolean
}

export const useLoadPreview = (viewMode?: Ref<string>) => {
  const previewService = usePreviewService()
  const { updateResourceField } = useResourcesStore()
  const previewQueue = new PQueue({ concurrency: 4 })

  const isTilesView = computed(() => unref(viewMode) === FolderViewModeConstants.name.tiles)
  const defaultProcessor = computed(() =>
    unref(isTilesView) ? ProcessorType.enum.fit : ProcessorType.enum.thumbnail
  )
  const defaultDimensions = computed(() =>
    unref(isTilesView) ? ImageDimension.Tile : ImageDimension.Thumbnail
  )

  const loadPreviewTask = useTask<string, LoadPreviewOptions[]>(function* (
    signal,
    { space, resource, dimensions, processor, updateStore = true }
  ) {
    const item = isProjectSpaceResource(resource) ? buildSpaceImageResource(resource) : resource
    const preview = yield previewQueue.add(() =>
      previewService.loadPreview(
        {
          space,
          resource: item,
          processor: processor || unref(defaultProcessor),
          dimensions: dimensions || unref(defaultDimensions)
        },
        true,
        true,
        signal
      )
    )

    if (preview && updateStore) {
      updateResourceField({ id: resource.id, field: 'thumbnail', value: preview })
    }

    return preview
  })

  const loadPreview = async (options: LoadPreviewOptions) => {
    const { resource, cancelRunning } = options
    if (cancelRunning) {
      cancelTasks()
    }

    if (isProjectSpaceResource(resource) && (!resource.spaceImageData || resource.disabled)) {
      return null
    }

    try {
      return await loadPreviewTask.perform(options)
    } catch (e) {
      // ignore errors on cancel
      if (e !== 'cancel') {
        console.error(e)
      }
    }
  }

  const previewsLoading = computed(() => loadPreviewTask.isRunning)

  const cancelTasks = () => {
    loadPreviewTask.cancelAll()
    previewQueue.clear()
  }

  onUnmounted(cancelTasks)

  return { loadPreview, previewsLoading }
}
