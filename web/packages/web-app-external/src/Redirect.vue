<template>
  <main
    class="external-redirect oc-height-viewport oc-flex oc-flex-column oc-flex-center oc-flex-middle"
  >
    <h1 v-if="pageTitle" class="oc-invisible-sr" v-text="pageTitle" />
    <div class="oc-card oc-card-body oc-text-center oc-width-large">
      <template v-if="hasError">
        <h2 key="external-redirect-error">
          <span v-text="$gettext('We could not open this file')" />
        </h2>
        <p v-text="errorMessage" />
      </template>
      <template v-else>
        <h2 key="external-redirect-loading">
          <span v-text="$gettext('One moment please…')" />
        </h2>
        <p v-text="$gettext('You are being redirected.')" />
      </template>
    </div>
  </main>
</template>

<script lang="ts" setup>
import { computed, ref, unref, watch } from 'vue'
import {
  queryItemAsString,
  useAppProviderService,
  useClientService,
  useRouteMeta,
  useRouteQuery,
  useSharesStore
} from '@ownclouders/web-pkg'
import { buildSpace } from '@ownclouders/web-client'
import { DavProperty } from '@ownclouders/web-client/webdav'
import { useRouter } from 'vue-router'
import { omit } from 'lodash-es'
import { useGettext } from 'vue3-gettext'
import { useApplicationReadyStore } from './piniaStores'
import { storeToRefs } from 'pinia'

const { $gettext } = useGettext()
const appProviderService = useAppProviderService()
const clientService = useClientService()
const sharesStore = useSharesStore()
const router = useRouter()
const { isReady } = storeToRefs(useApplicationReadyStore())

const appQuery = useRouteQuery('app')
const appNameQuery = useRouteQuery('appName')
const fileIdQuery = useRouteQuery('fileId')

const hasError = ref(false)
const errorMessage = ref('')

// An explicit app/appName query parameter always wins. It is set by callers that
// already know which app to open the file with.
const explicitAppName = computed(() => {
  if (unref(appQuery)) {
    return queryItemAsString(unref(appQuery))
  }
  if (unref(appNameQuery)) {
    return queryItemAsString(unref(appNameQuery))
  }
  return ''
})

const fail = (message: string) => {
  errorMessage.value = message
  hasError.value = true
}

const redirectToApp = (appName: string) => {
  hasError.value = false
  router.replace({
    name: `external-${appName.toLowerCase()}-apps`,
    query: omit(unref(router.currentRoute).query, ['app', 'appName'])
  })
}

// Resolve the target app from the file's mime type. A single PROPFIND by file id
// is enough and is safe on a cold deep link (no space from the store required),
// mirroring useGetResourceContext's loadFileInfoById.
const resolveAppNameByFileId = async (fileId: string): Promise<string | undefined> => {
  const space = buildSpace({ id: fileId, name: '' }, sharesStore.graphRoles)
  const resource = await clientService.webdav.getFileInfo(
    space,
    { fileId },
    { davProperties: [DavProperty.MimeType, DavProperty.FileId, DavProperty.Name] }
  )

  if (!resource?.mimeType) {
    return undefined
  }

  return appProviderService.getDefaultAppNameForMimeType(resource.mimeType)
}

const redirect = async () => {
  const explicit = unref(explicitAppName)
  if (explicit) {
    redirectToApp(explicit)
    return
  }

  const fileId = queryItemAsString(unref(fileIdQuery))
  if (!fileId) {
    fail($gettext('No file was specified to open.'))
    return
  }

  let appName: string | undefined
  try {
    appName = await resolveAppNameByFileId(fileId)
  } catch (e) {
    console.error(e)
    fail(
      $gettext(
        'The file could not be loaded. It may have been deleted or you might not have access to it.'
      )
    )
    return
  }

  if (!appName) {
    fail($gettext('No application is available to open this file type.'))
    return
  }

  redirectToApp(appName)
}

watch(
  isReady,
  (ready) => {
    if (!ready) {
      return
    }
    void redirect()
  },
  { immediate: true }
)

const title = useRouteMeta('title')
const pageTitle = computed(() => {
  return $gettext(unref(title))
})
</script>

<style lang="scss">
.external-redirect {
  .oc-card {
    background: var(--oc-color-background-highlight);
    border-radius: 15px;

    &-body {
      h2 {
        margin-top: 0;
      }
      p {
        font-size: var(--oc-font-size-large);
      }
    }
  }
}
</style>
