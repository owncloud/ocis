<template>
  <div>
    <span v-if="loading">
      <oc-spinner />
    </span>
    <span
      v-else-if="exportInProgress"
      class="oc-flex oc-flex-middle"
      data-testid="export-in-process"
    >
      <div class="oc-flex oc-flex-middle">
        <oc-icon name="time" fill-type="line" size="small" />
        <span
          class="oc-ml-xs"
          v-text="$gettext('Export is being processed. This can take up to 24 hours.')"
        />
      </div>
    </span>
    <div v-else class="oc-flex">
      <oc-button
        appearance="raw"
        variation="primary"
        data-testid="request-export-btn"
        class="oc-mr-s"
        @click="requestExport"
      >
        <div class="oc-flex oc-flex-middle">
          <oc-icon name="question-answer" fill-type="line" size="small" />
          <span class="oc-ml-xs" v-text="$gettext('Request new export')" />
        </div>
      </oc-button>
      <oc-button
        v-if="exportFile"
        v-oc-tooltip="$gettext('Latest export from: %{date}', { date: exportDate })"
        appearance="raw"
        variation="primary"
        data-testid="download-export-btn"
        @click="downloadExport"
      >
        <div class="oc-flex oc-flex-middle">
          <oc-icon name="download" fill-type="line" size="small" />
          <span class="oc-ml-xs" v-text="$gettext('Download')" />
        </div>
      </oc-button>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onMounted, onUnmounted, ref, unref } from 'vue'
import { useTask } from 'vue-concurrency'
import { useGettext } from 'vue3-gettext'
import { Resource } from '@ownclouders/web-client'
import { useClientService, useMessages, useSpacesStore, useUserStore } from '@ownclouders/web-pkg'
import { useDownloadFile } from '@ownclouders/web-pkg'
import { formatDateFromJSDate } from '@ownclouders/web-pkg'

const GDPR_EXPORT_FILE_NAME = '.personal_data_export.json'
const POLLING_INTERVAL = 30000

export default defineComponent({
  name: 'GdprExport',
  setup() {
    const { showMessage, showErrorMessage } = useMessages()
    const userStore = useUserStore()
    const spacesStore = useSpacesStore()
    const language = useGettext()
    const { $gettext } = language
    const clientService = useClientService()
    const { downloadFile } = useDownloadFile()

    const loading = ref(true)
    const checkInterval = ref<ReturnType<typeof setInterval>>()
    const exportFile = ref<Resource>()
    const exportInProgress = ref(false)

    const loadExportTask = useTask(function* (signal) {
      try {
        const resource = yield clientService.webdav.getFileInfo(
          spacesStore.personalSpace,
          { path: `/${GDPR_EXPORT_FILE_NAME}` },
          { signal }
        )

        if (resource.processing) {
          exportInProgress.value = true
          if (!unref(checkInterval)) {
            checkInterval.value = setInterval(() => {
              loadExportTask.perform()
            }, POLLING_INTERVAL)
          }
          return
        }

        exportFile.value = resource
        exportInProgress.value = false
        if (unref(checkInterval)) {
          clearInterval(unref(checkInterval))
          checkInterval.value = undefined
        }
      } catch (e) {
        if (e.statusCode !== 404) {
          // resource seems to exist, but something else went wrong
          console.error(e)
        }
      } finally {
        loading.value = false
      }
    }).restartable()

    const requestExport = async () => {
      try {
        await clientService.graphAuthenticated.users.exportPersonalData(userStore.user.id, {
          storageLocation: `/${GDPR_EXPORT_FILE_NAME}`
        })
        await loadExportTask.perform()
        showMessage({ title: $gettext('GDPR export has been requested') })
      } catch (e) {
        showErrorMessage({
          title: $gettext('GDPR export could not be requested. Please contact an administrator.'),
          errors: [e]
        })
      }
    }
    const downloadExport = () => {
      return downloadFile(spacesStore.personalSpace, unref(exportFile))
    }

    const exportDate = computed(() => {
      return formatDateFromJSDate(new Date(unref(exportFile).mdate), language.current)
    })

    onMounted(() => {
      loadExportTask.perform()
    })

    onUnmounted(() => {
      if (unref(checkInterval)) {
        clearInterval(unref(checkInterval))
      }
    })

    return {
      loading,
      loadExportTask,
      exportFile,
      exportInProgress,
      requestExport,
      downloadExport,
      exportDate
    }
  }
})
</script>
