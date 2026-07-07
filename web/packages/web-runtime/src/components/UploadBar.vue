<template>
  <div
    v-if="showInfo"
    id="upload-bar"
    class="oc-background-muted upload-bar"
    :class="{ 'oc-rounded oc-box-shadow-medium': headless === false }"
  >
    <div
      class="upload-bar-title oc-flex oc-flex-between oc-flex-middle oc-px-m oc-py-s oc-rounded-top"
    >
      <p v-oc-tooltip="uploadDetails" class="oc-my-xs" v-text="uploadInfoTitle" />
      <oc-button
        v-if="!filesInProgressCount"
        id="close-upload-bar-btn"
        v-oc-tooltip="$gettext('Close')"
        :aria-label="$gettext('Close')"
        appearance="raw-inverse"
        variation="brand"
        @click="closeInfo"
      >
        <oc-icon name="close" />
      </oc-button>
    </div>
    <div
      class="upload-bar-status oc-px-m oc-pt-m oc-flex oc-flex-between oc-flex-middle"
      :class="{
        'oc-pb-m': !runningUploads
      }"
    >
      <div v-if="runningUploads" class="oc-flex oc-flex-middle">
        <oc-icon v-if="uploadsPaused" name="pause" size="small" class="oc-mr-s" />
        <oc-spinner v-else size="small" class="oc-mr-s" :aria-label="$gettext('Uploading')" />
        <span class="oc-text-small oc-text-muted" v-text="remainingTime" />
      </div>
      <div
        v-else
        class="upload-bar-label"
        :class="{
          'upload-bar-danger': Object.keys(errors).length && !uploadsCancelled,
          'upload-bar-success': !Object.keys(errors).length && !uploadsCancelled
        }"
      >
        {{ uploadingLabel }}
      </div>
      <div class="oc-flex">
        <oc-button
          appearance="raw"
          class="oc-text-muted oc-text-small upload-bar-toggle-details-btn"
          @click="toggleInfo"
        >
          {{ infoExpanded ? $gettext('Hide details') : $gettext('Show details') }}
        </oc-button>
        <oc-button
          v-if="!runningUploads && Object.keys(errors).length && !disableActions"
          v-oc-tooltip="$gettext('Retry all failed uploads')"
          class="oc-ml-s"
          appearance="raw"
          :aria-label="$gettext('Retry all failed uploads')"
          @click="retryUploads"
        >
          <oc-icon name="restart" fill-type="line" />
        </oc-button>

        <oc-button
          v-if="
            runningUploads &&
            uploadsPausable &&
            !inPreparation &&
            !inFinalization &&
            !disableActions
          "
          id="pause-upload-bar-btn"
          v-oc-tooltip="uploadsPaused ? $gettext('Resume upload') : $gettext('Pause upload')"
          class="oc-ml-s"
          appearance="raw"
          :aria-label="uploadsPaused ? $gettext('Resume upload') : $gettext('Pause upload')"
          @click="togglePauseUploads"
        >
          <oc-icon :name="uploadsPaused ? 'play-circle' : 'pause-circle'" fill-type="line" />
        </oc-button>
        <oc-button
          v-if="runningUploads && !inPreparation && !inFinalization && !disableActions"
          id="cancel-upload-bar-btn"
          v-oc-tooltip="$gettext('Cancel upload')"
          class="oc-ml-s"
          appearance="raw"
          :aria-label="$gettext('Cancel upload')"
          @click="cancelAllUploads"
        >
          <oc-icon name="close-circle" fill-type="line" />
        </oc-button>
      </div>
    </div>
    <div v-if="runningUploads" class="upload-bar-progress oc-mx-m oc-pb-m oc-mt-s oc-text">
      <oc-progress
        :value="totalProgress"
        :max="100"
        size="small"
        :indeterminate="!filesInProgressCount"
        :aria-label="$gettext('Upload progress')"
      />
    </div>
    <div
      v-if="infoExpanded"
      class="upload-bar-items oc-px-m oc-pb-m"
      :class="{ 'has-errors': showErrorLog }"
    >
      <ul class="oc-list">
        <li v-for="(item, idx) in uploads" :key="idx">
          <span class="oc-flex oc-flex-middle">
            <oc-icon v-if="item.status === 'error'" name="close" variation="danger" size="small" />
            <oc-icon
              v-else-if="item.status === 'success'"
              name="check"
              variation="success"
              size="small"
            />
            <oc-icon v-else-if="item.status === 'cancelled'" name="close" size="small" />
            <oc-icon v-else-if="uploadsPaused" name="pause" size="small" />
            <div v-else class="oc-flex"><oc-spinner size="small" /></div>
            <resource-list-item
              v-if="displayFileAsResource(item)"
              :key="item.path"
              class="oc-ml-s"
              :resource="item as Resource"
              :is-path-displayed="true"
              :is-resource-clickable="isResourceClickable(item)"
              :parent-folder-name="parentFolderName(item)"
              :link="resourceLink(item)"
              :parent-folder-link="parentFolderLink(item)"
            />
            <span v-else class="oc-flex oc-flex-middle oc-text-truncate">
              <resource-icon
                :resource="item as Resource"
                size="large"
                class="file_info__icon oc-mx-s"
              />
              <resource-name
                :name="item.name"
                :extension="item.extension"
                :type="item.type"
                full-path=""
                :is-path-displayed="false"
              />
            </span>
          </span>
          <span
            v-if="getUploadItemMessage(item)"
            class="upload-bar-message oc-ml-xs oc-text-small"
            :class="getUploadItemClass(item)"
            v-text="getUploadItemMessage(item)"
          />
        </li>
      </ul>
    </div>
    <oc-error-log
      v-if="showErrorLog"
      class="upload-bar-error-log oc-pt-m oc-pb-m oc-px-m"
      :content="uploadErrorLogContent"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, unref, computed, onMounted, onWatcherCleanup } from 'vue'
import { isUndefined } from 'lodash-es'
import { getSpeed } from '@uppy/utils'

import { extractParentFolderName, HttpError, Resource, urlJoin } from '@ownclouders/web-client'
import {
  OcUppyFile,
  queryItemAsString,
  UppyService,
  useConfigStore,
  useService
} from '@ownclouders/web-pkg'
import { formatFileSize, ResourceListItem, ResourceIcon, ResourceName } from '@ownclouders/web-pkg'
import { storeToRefs } from 'pinia'
import { RouteLocationNamedRaw } from 'vue-router'
import { useGettext } from 'vue3-gettext'

type UploadResult = OcUppyFile & {
  path?: string
  targetRoute?: RouteLocationNamedRaw
  status?: string
  filesCount?: number
  successCount?: number
  errorCount?: number
}

const { infoExpandedInitial = false, headless = false } = defineProps<{
  /**
   * show the info including all uploads?
   * Prop only works intially, state gets copied to local var infoExpanded
   */
  infoExpandedInitial?: boolean
  /**
   * Render as headless component?
   * Renders the header and the close button
   */
  headless?: boolean
}>()

const { current, $pgettext, $npgettext } = useGettext()
const uppyService = useService<UppyService>('$uppyService')

const configStore = useConfigStore()
const { options: configOptions } = storeToRefs(configStore)

const showInfo = ref(false)
const infoExpanded = ref(false)
const uploads = ref<Record<string, UploadResult>>({})
const errors = ref<Record<string, HttpError>>({})
const successful = ref<string[]>([])
const filesInProgressCount = ref(0)
const totalProgress = ref(0)
const uploadsPaused = ref(false)
const uploadsCancelled = ref(false)
const inFinalization = ref(false)
const inPreparation = ref(true)
const runningUploads = ref(0)
const bytesTotal = ref(0)
const bytesUploaded = ref(0)
const uploadSpeed = ref(0)
const filesInEstimation = ref<Record<string, number>>({})
const timeStarted = ref<Date>(null)
const remainingTime = ref<string>(undefined)
const disableActions = ref(false)

const uploadDetails = computed(() => {
  if (!unref(uploadSpeed) || !unref(runningUploads)) {
    return ''
  }
  const uploadedBytes = formatFileSize(unref(bytesUploaded), unref(current))
  const totalBytes = formatFileSize(unref(bytesTotal), unref(current))
  const currentUploadSpeed = formatFileSize(unref(uploadSpeed), unref(current))

  return $pgettext(
    'The upload details displayed in the upload overlay.',
    '%{uploadedBytes} of %{totalBytes} (%{currentUploadSpeed}/s)',
    {
      uploadedBytes,
      totalBytes,
      currentUploadSpeed
    }
  )
})

const uploadInfoTitle = computed(() => {
  if (unref(inFinalization)) {
    return $pgettext(
      'The title of the upload overlay when upload has been transferred but still needs to be finalized.',
      'Finalizing upload...'
    )
  }

  if (unref(filesInProgressCount) && !unref(inPreparation)) {
    return $npgettext(
      'The title of the upload overlay when upload is being transferred.',
      '%{ filesInProgressCount } item uploading...',
      '%{ filesInProgressCount } items uploading...',
      unref(filesInProgressCount),
      { filesInProgressCount: (unref(filesInProgressCount) as number).toString() }
    )
  }

  if (unref(uploadsCancelled)) {
    return $pgettext(
      'The title of the upload overlay when upload has been cancelled.',
      'Upload cancelled'
    )
  }

  if (Object.keys(unref(errors)).length) {
    return $pgettext('The title of the upload overlay when upload has failed.', 'Upload failed')
  }

  if (!unref(runningUploads)) {
    return $pgettext(
      'The title of the upload overlay when upload has been completed.',
      'Upload complete'
    )
  }

  return $pgettext(
    'The title of the upload overlay when upload is being prepared.',
    'Preparing upload...'
  )
})

const uploadingLabel = computed(() => {
  if (Object.keys(unref(errors)).length) {
    const count = unref(successful).length + Object.keys(unref(errors)).length

    return $npgettext(
      'The label of the upload overlay displayed when at least one item has failed.',
      '%{ errors } of %{ uploads } item failed',
      '%{ errors } of %{ uploads } items failed',
      count,
      { uploads: count.toString(), errors: Object.keys(unref(errors)).length.toString() }
    )
  }

  return $npgettext(
    'The label of the upload overlay displayed when all items have been uploaded.',
    '%{ successfulUploads } item uploaded',
    '%{ successfulUploads } items uploaded',
    unref(successful).length,
    { successfulUploads: unref(successful).length.toString() }
  )
})

const uploadsPausable = computed(() => {
  return uppyService.tusActive()
})

const showErrorLog = computed(() => {
  return unref(infoExpanded) && unref(uploadErrorLogContent)
})

const uploadErrorLogContent = computed(() => {
  const requestIds = Object.values(unref(errors)).reduce((acc: Array<string>, error: any) => {
    const requestId = error.originalRequest?._headers?.['X-Request-ID']

    if (requestId) {
      acc.push(requestId)
    }

    return acc
  }, []) as Array<any>

  return requestIds.map((item) => `X-Request-Id: ${item}`).join('\r\n')
})

const onBeforeUnload = (e: BeforeUnloadEvent) => {
  if (unref(runningUploads)) {
    e.preventDefault()
  }
}

const getRemainingTime = (remainingMilliseconds: number): string => {
  const roundedRemainingMinutes = Math.round(remainingMilliseconds / 1000 / 60)

  if (roundedRemainingMinutes >= 1 && roundedRemainingMinutes < 60) {
    return $npgettext(
      'The remaining upload time in minutes displayed in the upload overlay.',
      '%{ roundedRemainingMinutes } minute left',
      '%{ roundedRemainingMinutes } minutes left',
      roundedRemainingMinutes,
      { roundedRemainingMinutes: roundedRemainingMinutes.toString() }
    )
  }

  const roundedRemainingHours = Math.round(remainingMilliseconds / 1000 / 60 / 60)

  if (roundedRemainingHours > 0) {
    return $npgettext(
      'The remaining upload time in hours displayed in the upload overlay.',
      '%{ roundedRemainingHours } hour left',
      '%{ roundedRemainingHours } hours left',
      roundedRemainingHours,
      { roundedRemainingHours: roundedRemainingHours.toString() }
    )
  }

  return $pgettext(
    'The remaining upload time in seconds displayed in the upload overlay.',
    'Few seconds left'
  )
}

const handleTopLevelFolderUpdate = (file: OcUppyFile, status: string) => {
  const topLevelFolder = unref(uploads)[file.meta.topLevelFolderId]

  if (status === 'success') {
    topLevelFolder.successCount += 1
  } else {
    topLevelFolder.errorCount += 1
  }

  // all files for this top level folder are finished
  if (topLevelFolder.successCount + topLevelFolder.errorCount === topLevelFolder.filesCount) {
    topLevelFolder.status = topLevelFolder.errorCount ? 'error' : 'success'
  }
}

const closeInfo = () => {
  showInfo.value = false
  infoExpanded.value = false

  cleanOverlay()
  resetProgress()
}

const cleanOverlay = () => {
  uploadsCancelled.value = false
  uploads.value = {}
  errors.value = {}
  successful.value = []
  filesInProgressCount.value = 0
  runningUploads.value = 0
  disableActions.value = false
}

const resetProgress = () => {
  bytesTotal.value = 0
  bytesUploaded.value = 0
  filesInEstimation.value = {}
  timeStarted.value = null
  remainingTime.value = undefined
  inPreparation.value = true
  inFinalization.value = false
  uploadsPaused.value = false
}

const displayFileAsResource = (file: UploadResult): boolean => {
  return !!file.targetRoute
}

const isResourceClickable = (file: UploadResult): boolean => {
  return file.meta.isFolder === true
}

const resourceLink = (file: UploadResult): RouteLocationNamedRaw => {
  if (!file.meta.isFolder) {
    return {}
  }
  return {
    ...file.targetRoute,
    params: {
      ...file.targetRoute.params,
      driveAliasAndItem: urlJoin(
        queryItemAsString(file.targetRoute.params.driveAliasAndItem),
        file.name,
        {
          leadingSlash: false
        }
      )
    },
    query: {
      ...file.targetRoute.query,
      ...(unref(configOptions).routing.idBased &&
        !isUndefined(file.meta.fileId) && { fileId: file.meta.fileId })
    }
  }
}

const parentFolderLink = (file: UploadResult): RouteLocationNamedRaw => {
  return {
    ...file.targetRoute,
    query: {
      ...file.targetRoute.query,
      ...(unref(configOptions).routing.idBased &&
        !isUndefined(file.meta.currentFolderId) && { fileId: file.meta.currentFolderId })
    }
  }
}

const buildRouteFromUppyResource = (resource: OcUppyFile): RouteLocationNamedRaw => {
  if (!resource.meta.routeName) {
    return null
  }

  return {
    name: resource.meta.routeName,
    params: {
      driveAliasAndItem: resource.meta.routeDriveAliasAndItem
    },
    query: {
      ...(resource.meta.routeShareId && { shareId: resource.meta.routeShareId })
    }
  }
}

const parentFolderName = (file: UploadResult): string => {
  const {
    meta: { spaceName, driveType }
  } = file

  const parentFolder = extractParentFolderName(file as Resource)
  if (parentFolder) {
    return parentFolder
  }

  if (driveType === 'personal') {
    return $pgettext(
      'The root parent folder name displayed in the upload overlay in personal space.',
      'Personal'
    )
  }

  if (driveType === 'public') {
    return $pgettext(
      'The root parent folder name displayed in the upload overlay in public links.',
      'Public link'
    )
  }

  return spaceName
}

const toggleInfo = () => {
  infoExpanded.value = !unref(infoExpanded)
}

const retryUploads = () => {
  filesInProgressCount.value += Object.keys(unref(errors)).length
  runningUploads.value += 1

  for (const fileID of Object.keys(unref(errors))) {
    uploads.value[fileID].status = undefined

    const topLevelFolderId = unref(uploads)[fileID].meta.topLevelFolderId

    if (topLevelFolderId) {
      uploads.value[topLevelFolderId].status = undefined
      uploads.value[topLevelFolderId].errorCount = 0
    }
  }

  errors.value = {}
  uppyService.retryAllUploads()
}

const togglePauseUploads = () => {
  if (unref(uploadsPaused)) {
    uppyService.resumeAllUploads()
    timeStarted.value = null
  } else {
    uppyService.pauseAllUploads()
  }

  uploadsPaused.value = !unref(uploadsPaused)
}

const cancelAllUploads = () => {
  uploadsCancelled.value = true
  filesInProgressCount.value = 0
  runningUploads.value = 0

  resetProgress()
  uppyService.cancelAllUploads()

  const _runningUploads = Object.values(unref(uploads)).filter(
    (u) => u.status !== 'success' && u.status !== 'error'
  )

  for (const item of _runningUploads) {
    uploads.value[item.meta.uploadId].status = 'cancelled'
  }
}

const getUploadItemMessage = (item: UploadResult) => {
  const error = unref(errors)[item.meta.uploadId]

  if (!error) {
    return
  }

  //TODO: Remove extraction code as soon as https://github.com/tus/tus-js-client/issues/448 is solved
  const formatErrorMessageToObject = (
    errorMessage: string
  ): {
    responseCode: number | null
    errorCode: string | null
    errorMessage: string | null
  } => {
    const responseCode = errorMessage.match(/response code: (\d+)/)?.[1]
    const responseText = errorMessage.match(/response text: ([\s\S]+?), request id/)?.[1]
    const errorBody = JSON.parse(responseText?.startsWith('{') ? responseText : '{}')

    return {
      responseCode: responseCode ? parseInt(responseCode) : null,
      errorCode: errorBody?.error?.code,
      errorMessage: errorBody?.error?.message
    }
  }

  const errorObject = formatErrorMessageToObject(error.message)

  if (unref(errors)[item.meta.uploadId]?.statusCode === 423) {
    return $pgettext(
      'The message displayed in the upload overlay for a single upload item when the folder you are uploading to is locked',
      "The folder you're uploading to is locked"
    )
  }

  switch (errorObject.responseCode) {
    case 507:
      return $pgettext(
        'The message displayed in the upload overlay for a single upload item when the quota has been exceeded',
        'Quota exceeded'
      )
    case 412:
      return $pgettext(
        'The message displayed in the upload overlay for a single upload item when the parent folder does not exist',
        'Parent folder does not exist'
      )
    default:
      return errorObject.errorMessage
        ? $pgettext(
            'The message displayed in the upload overlay for a single upload item when the error message is returned from the API',
            errorObject.errorMessage
          )
        : $pgettext(
            'The message displayed in the upload overlay for a single upload item when the error message is unknown',
            'Unknown error'
          )
  }
}

const getUploadItemClass = (item: UploadResult) => {
  return unref(errors)[item.meta.uploadId] ? 'upload-bar-danger' : 'upload-bar-success'
}

watch(runningUploads, (val) => {
  if (val === 0) {
    window.removeEventListener('beforeunload', onBeforeUnload)

    return
  }

  window.addEventListener('beforeunload', onBeforeUnload)

  onWatcherCleanup(() => {
    window.removeEventListener('beforeunload', onBeforeUnload)
  })
})

onMounted(() => {
  infoExpanded.value = infoExpandedInitial

  uppyService.subscribe('uploadStarted', () => {
    if (!unref(remainingTime)) {
      remainingTime.value = $pgettext(
        'The message displayed in the upload overlay when the estimated time is being calculated',
        'Calculating estimated time...'
      )
    }

    // No upload in progress -> clean overlay
    if (!unref(runningUploads) && unref(showInfo)) {
      cleanOverlay()
    }

    showInfo.value = true
    runningUploads.value += 1
    inFinalization.value = false
  })

  uppyService.subscribe('addedForUpload', (files: OcUppyFile[]) => {
    filesInProgressCount.value += files.filter((f) => !f.meta.isFolder).length

    for (const file of files) {
      if (!unref(disableActions) && file.isRemote) {
        disableActions.value = true
      }

      if (file.data?.size) {
        bytesTotal.value += file.data.size
      }

      const { relativeFolder, uploadId, topLevelFolderId } = file.meta
      const isTopLevelItem = !relativeFolder

      // only add top level items to this.uploads because we only show those
      if (isTopLevelItem) {
        uploads.value[uploadId] = file

        // top level folders get initialized with file counts about their files inside
        if (file.meta.isFolder && uploads.value[uploadId].filesCount === undefined) {
          uploads.value[uploadId].filesCount = 0
          uploads.value[uploadId].errorCount = 0
          uploads.value[uploadId].successCount = 0
        }
      }

      // count all files inside top level folders to mark them as successful or failed later
      if (!file.meta.isFolder && !isTopLevelItem && unref(uploads)[topLevelFolderId]) {
        uploads.value[topLevelFolderId].filesCount += 1
      }
    }
  })

  uppyService.subscribe('uploadCompleted', () => {
    runningUploads.value -= 1

    if (!unref(runningUploads)) {
      resetProgress()
    }
  })

  uppyService.subscribe('progress', (value: number) => {
    totalProgress.value = value
  })

  uppyService.subscribe(
    'upload-progress',
    ({ file, progress }: { file: OcUppyFile; progress: { bytesUploaded: number } }) => {
      if (!unref(timeStarted)) {
        timeStarted.value = new Date()
        inPreparation.value = false
      }

      if (unref(filesInEstimation)[file.meta.uploadId] === undefined) {
        filesInEstimation.value[file.meta.uploadId] = 0
      }

      const byteIncrease = progress.bytesUploaded - unref(filesInEstimation)[file.meta.uploadId]

      bytesUploaded.value += byteIncrease
      filesInEstimation.value[file.meta.uploadId] = progress.bytesUploaded

      const timeElapsed = +new Date().getTime() - unref(timeStarted).getTime()

      uploadSpeed.value = getSpeed({
        bytesUploaded: unref(bytesUploaded),
        uploadStarted: unref(timeStarted).getTime(),
        bytesTotal: unref(bytesTotal)
      })

      const progressPercent = (100 * unref(bytesUploaded)) / unref(bytesTotal)

      if (progressPercent === 0) {
        return
      }

      const totalTimeNeededInMilliseconds = (timeElapsed / progressPercent) * 100
      const remainingMilliseconds = totalTimeNeededInMilliseconds - timeElapsed

      remainingTime.value = getRemainingTime(remainingMilliseconds)

      if (progressPercent === 100) {
        inFinalization.value = true
      }
    }
  )

  uppyService.subscribe('uploadError', ({ file, error }: { file: OcUppyFile; error: Error }) => {
    if (unref(errors)[file.meta.uploadId]) {
      return
    }

    // file inside folder -> was not added to this.uploads, but must be now because of error
    if (!unref(uploads)[file.meta.uploadId]) {
      uploads.value[file.meta.uploadId] = file
    }

    if (file.meta.relativePath) {
      uploads.value[file.meta.uploadId].path = file.meta.relativePath
    } else {
      uploads.value[file.meta.uploadId].path = urlJoin(file.meta.currentFolder, file.name)
    }

    uploads.value[file.meta.uploadId].targetRoute = buildRouteFromUppyResource(file)
    uploads.value[file.meta.uploadId].status = 'error'
    errors.value[file.meta.uploadId] = error as HttpError
    filesInProgressCount.value -= 1
    runningUploads.value -= 1

    if (file.meta.topLevelFolderId) {
      handleTopLevelFolderUpdate(file, 'error')
    }
  })

  uppyService.subscribe('uploadSuccess', (file: OcUppyFile) => {
    // item inside folder
    if (!unref(uploads)[file.meta.uploadId]) {
      if (!file.meta.isFolder) {
        successful.value.push(file.meta.uploadId)
        filesInProgressCount.value -= 1

        if (file.meta.topLevelFolderId) {
          handleTopLevelFolderUpdate(file, 'success')
        }
      }

      return
    }

    // file inside folder that succeeded via retry can now be removed again from this.uploads
    if (file.meta.relativeFolder) {
      if (!file.meta.isFolder) {
        successful.value.push(file.meta.uploadId)
        filesInProgressCount.value -= 1

        if (file.meta.topLevelFolderId) {
          handleTopLevelFolderUpdate(file, 'success')
        }
      }

      delete uploads.value[file.meta.uploadId]
      return
    }

    uploads.value[file.meta.uploadId] = file
    uploads.value[file.meta.uploadId].path = urlJoin(file.meta.currentFolder, file.name)
    uploads.value[file.meta.uploadId].targetRoute = buildRouteFromUppyResource(file)

    if (!file.meta.isFolder) {
      uploads.value[file.meta.uploadId].status = 'success'
      successful.value.push(file.meta.uploadId)
      filesInProgressCount.value -= 1
    }
  })
})

// FIXME: drop once tests are adjusted not to interact with the vm directly
defineExpose({
  showInfo,
  inPreparation,
  filesInProgressCount,
  runningUploads,
  successful,
  errors,
  inFinalization,
  infoExpanded,
  uploads,
  uploadsCancelled,
  getRemainingTime
})
</script>

<style lang="scss">
.upload-bar {
  @media (max-width: 640px) {
    margin: 0 auto;
  }

  .oc-resource-details {
    padding-left: var(--oc-space-xsmall);
  }

  .upload-bar-title {
    background-color: var(--oc-color-swatch-brand-default);
  }

  .upload-bar-title p {
    color: var(--oc-color-swatch-brand-contrast);
  }

  .oc-resource-indicators .parent-folder .text {
    color: var(--oc-color-text-default);
  }

  .upload-bar-items {
    overflow-y: auto;
  }

  .upload-bar-danger {
    color: var(--oc-color-swatch-danger-default);
  }

  .upload-bar-success {
    color: var(--oc-color-swatch-success-default);
  }
}
</style>
