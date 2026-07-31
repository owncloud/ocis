<template>
  <iframe
    v-if="appUrl && method === 'GET'"
    ref="appIframe"
    :src="appUrl"
    class="oc-width-1-1 oc-height-1-1"
    :title="iFrameTitle"
    allowfullscreen
    allow="camera; clipboard-read; clipboard-write"
    @load="onIframeLoad"
  />
  <div v-if="appUrl && method === 'POST' && formParameters" class="oc-height-1-1 oc-width-1-1">
    <form :action="appUrl" target="app-iframe" method="post">
      <input ref="subm" type="submit" :value="formParameters" class="oc-hidden" />
      <div v-for="(item, key, index) in formParameters" :key="index">
        <input :name="key" :value="item" type="hidden" />
      </div>
    </form>
    <iframe
      ref="appIframe"
      name="app-iframe"
      :src="appUrl"
      class="oc-width-1-1 oc-height-1-1"
      :title="iFrameTitle"
      allowfullscreen
      allow="camera; clipboard-read; clipboard-write"
      @load="onIframeLoad"
    />
  </div>
</template>

<script lang="ts" setup>
import { stringify } from 'qs'
import {
  computed,
  inject,
  unref,
  nextTick,
  ref,
  watch,
  VNodeRef,
  onMounted,
  onBeforeUnmount,
  type Ref
} from 'vue'
import { useTask } from 'vue-concurrency'
import { useGettext } from 'vue3-gettext'

import { Resource, SpaceResource } from '@ownclouders/web-client'
import { urlJoin } from '@ownclouders/web-client'
import {
  isSameResource,
  useCapabilityStore,
  useConfigStore,
  useMessages,
  useModals,
  useRequest,
  useAppProviderService,
  useRoute,
  queryItemAsString,
  useRouteQuery
} from '@ownclouders/web-pkg'
import {
  isProjectSpaceResource,
  isPublicSpaceResource,
  isShareSpaceResource
} from '@ownclouders/web-client'

type ExtendedNavigator = Navigator & {
  userAgentData?: {
    mobile: boolean
    platform: string
    brands: { brand: string; version: string }[]
  }
}

interface Props {
  space: SpaceResource
  resource: Resource
  isReadOnly: boolean
}
const props = defineProps<Props>()
const language = useGettext()
const { $gettext } = language
const { showErrorMessage } = useMessages()
const { dispatchModal } = useModals()
const capabilityStore = useCapabilityStore()
const configStore = useConfigStore()
const route = useRoute()
const appProviderService = useAppProviderService()
const { makeRequest } = useRequest()

const viewModeQuery = useRouteQuery('view_mode')
const isMobileWidth =
  inject<Ref<boolean>>('isMobileWidth') || (navigator as ExtendedNavigator).userAgentData?.mobile
const viewModeQueryValue = computed(() => {
  return queryItemAsString(unref(viewModeQuery))
})

const templateIdQuery = useRouteQuery('templateId')
const templateIdQueryValue = computed(() => {
  return queryItemAsString(unref(templateIdQuery))
})

const appName = computed(() => {
  const lowerCaseAppName = unref(route)
    .name.toString()
    .replace('external-', '')
    .replace('-apps', '')
  return appProviderService.appNames.find((appName) => appName.toLowerCase() === lowerCaseAppName)
})

// The grouped Save-As control and the Host_PostmessageReady handshake are
// Collabora-specific WOPI extensions; other app providers don't implement them.
const isCollabora = computed(() => unref(appName) === 'Collabora')

const appUrl = ref()
const formParameters = ref({})
const method = ref()
const subm: VNodeRef = ref()
const appIframe = ref<HTMLIFrameElement>()

// Origin of the editor (e.g. the Collabora host) used as the target origin when
// posting messages into the iframe. Resolved against the current origin so a
// relative `app_url` (WOPI server co-hosted with the web app) still yields a
// usable origin instead of permanently failing the postMessage origin check.
const appOrigin = computed(() => {
  try {
    return new URL(unref(appUrl), window.location.origin).origin
  } catch {
    return ''
  }
})

// Collabora exposes "Save As" (export to another format, saved back to storage
// via WOPI PutRelativeFile) as a host-delegated operation: it shows a grouped
// Save-As control and posts a `UI_SaveAs` message, expecting the host to reply
// with the target filename. Request that grouped control via `ui_defaults`,
// which is a semicolon-delimited list, so append rather than overwrite it in
// case the server already decorates the URL with its own defaults.
const withSaveAsUiDefaults = (rawUrl: string) => {
  try {
    const url = new URL(rawUrl)
    const existing = url.searchParams.get('ui_defaults')
    url.searchParams.set(
      'ui_defaults',
      existing ? `${existing};SaveAsMode=group` : 'SaveAsMode=group'
    )
    return url.href
  } catch {
    return rawUrl
  }
}

// Post a WOPI postMessage into the editor iframe (messages are JSON strings).
// Accepts an explicit target/origin so callers that captured them at the time
// a user action started (e.g. opening the Save-As modal) don't re-read the
// live refs, which may have moved on to a different resource by the time the
// user confirms. Returns whether the message was actually sent.
const postToApp = (
  message: Record<string, unknown>,
  target: Window | null | undefined = unref(appIframe)?.contentWindow,
  origin: string = unref(appOrigin)
): boolean => {
  if (!target || !origin) {
    return false
  }
  target.postMessage(JSON.stringify({ ...message, SendTime: Date.now() }), origin)
  return true
}

// The editor only emits its rich postMessage API (UI_SaveAs, App_LoadingStatus,
// ...) once the host has announced itself with `Host_PostmessageReady`. Without
// this handshake "Save As" silently does nothing. Send it as soon as the iframe
// has loaded. This handshake is Collabora-specific; other providers don't expect it.
const onIframeLoad = () => {
  if (!unref(isCollabora)) {
    return
  }
  postToApp({ MessageId: 'Host_PostmessageReady', Values: {} })
}

// Reject empty names, path separators and the reserved "." / ".." names before
// the name is sent on. The server-side `PutRelativeFile` sanitizes the name as
// well, but validating here gives immediate feedback and avoids a pointless
// round-trip on invalid input.
const validateSaveAsFilename = (filename: string): string => {
  const trimmed = filename?.trim()
  if (!trimmed) {
    return $gettext('The file name cannot be empty.')
  }
  if (/[/\\]/.test(trimmed)) {
    return $gettext('The file name cannot contain "/" or "\\".')
  }
  if (trimmed === '.' || trimmed === '..') {
    return $gettext('The file name cannot be "." or "..".')
  }
  return ''
}

// Guards against a second Save-As modal being dispatched while one is already
// open (e.g. the editor's grouped Save-As control emitting `UI_SaveAs` twice
// in quick succession).
const isSaveAsModalOpen = ref(false)

// Handle the editor's "Save As" request: ask the user for the copy's name (the
// extension selects the export format) and reply with `Action_SaveAs`, which
// makes the editor render to that format and PutRelativeFile it into the space.
const onSaveAs = () => {
  if (props.isReadOnly) {
    showErrorMessage({ title: $gettext('Cannot save a copy: file is read-only') })
    return
  }
  if (unref(isSaveAsModalOpen)) {
    return
  }

  // Capture the target iframe/origin now, at dispatch time, rather than
  // re-reading the (reactive) refs in `onConfirm`: if the component gets
  // reused for a different resource while the modal is open, the live refs
  // would point at the new document by the time the user presses Save.
  const target = unref(appIframe)?.contentWindow
  const origin = unref(appOrigin)

  isSaveAsModalOpen.value = true
  dispatchModal({
    variation: 'passive',
    title: $gettext('Save a copy'),
    confirmText: $gettext('Save'),
    hasInput: true,
    inputValue: props.resource.name,
    inputLabel: $gettext('File name'),
    onInput: (filename: string, setError: (error: string) => void) => {
      setError(validateSaveAsFilename(filename))
    },
    onCancel: () => {
      isSaveAsModalOpen.value = false
    },
    onConfirm: (filename: string) => {
      isSaveAsModalOpen.value = false
      // Guard again on confirm (defense-in-depth); onInput normally prevents
      // reaching here with an invalid name.
      if (validateSaveAsFilename(filename)) {
        return
      }
      const sent = postToApp(
        { MessageId: 'Action_SaveAs', Values: { Filename: filename.trim(), Notify: true } },
        target,
        origin
      )
      if (!sent) {
        showErrorMessage({ title: $gettext('Cannot save a copy: the editor is not available') })
      }
    }
  })
}

const iFrameTitle = computed(() => {
  return $gettext('"%{appName}" app content area', {
    appName: unref(appName)
  })
})

const errorPopup = (error: string) => {
  showErrorMessage({
    title: $gettext('An error occurred'),
    desc: error,
    errors: [new Error(error)]
  })
}

const loadAppUrl = useTask(function* (signal, viewMode: string) {
  try {
    if (props.isReadOnly && viewMode === 'write') {
      showErrorMessage({ title: $gettext('Cannot open file in edit mode as it is read-only') })
      return
    }

    const fileId = props.resource.fileId
    const baseUrl = urlJoin(configStore.serverUrl, capabilityStore.filesAppProviders[0].open_url)

    const query = stringify({
      file_id: fileId,
      lang: language.current,
      mobile: unref(isMobileWidth) ? 1 : 0,
      ...(unref(appName) && { app_name: encodeURIComponent(unref(appName)) }),
      ...(viewMode && { view_mode: viewMode }),
      ...(unref(templateIdQueryValue) && { template_id: unref(templateIdQueryValue) })
    })

    const url = `${baseUrl}?${query}`
    const response = yield makeRequest('POST', url, {
      validateStatus: () => true,
      signal
    })

    if (response.status !== 200) {
      switch (response.status) {
        case 425:
          errorPopup(
            $gettext(
              'This file is currently being processed and is not yet available for use. Please try again shortly.'
            )
          )
          break
        default:
          errorPopup(response.data?.message)
      }

      throw new Error('Error fetching app information')
    }

    if (!response.data.app_url || !response.data.method) {
      throw new Error('Error in app server response')
    }

    appUrl.value = unref(isCollabora)
      ? withSaveAsUiDefaults(response.data.app_url)
      : response.data.app_url
    method.value = response.data.method

    if (response.data.form_parameters) {
      formParameters.value = response.data.form_parameters
    }

    if (method.value === 'POST' && formParameters.value) {
      yield nextTick()
      unref(subm).click()
    }
  } catch (e) {
    console.error('web-app-external error', e)
    throw e
  }
}).restartable()

const determineOpenAsPreview = (appName: string) => {
  const openAsPreview = configStore.options.editor.openAsPreview
  return openAsPreview === true || (Array.isArray(openAsPreview) && openAsPreview.includes(appName))
}

// Single handler for the editor's postMessage events. Only messages coming from
// the editor's own origin are accepted.
const onAppMessage = (event: MessageEvent) => {
  // Fail closed: reject every message until the editor origin is known (i.e.
  // `appUrl` has resolved to an absolute URL) and accept only messages from that
  // exact origin. The listener is attached on mount, before the WOPI `open_url`
  // POST resolves, so an empty `appOrigin` must reject rather than wave messages
  // through. No legitimate editor message can arrive before `appUrl` is set (the
  // iframe cannot have loaded yet), so this is strictly correct and also covers a
  // backend returning a non-absolute `app_url`.
  if (!unref(appOrigin) || event.origin !== unref(appOrigin)) {
    return
  }
  let message: { MessageId?: string }
  try {
    message = JSON.parse(event.data)
  } catch {
    return
  }
  switch (message?.MessageId) {
    case 'UI_Edit':
      // switch to write mode when edit is clicked
      if (determineOpenAsPreview(unref(appName))) {
        loadAppUrl.perform('write')
      }
      break
    case 'UI_SaveAs':
      onSaveAs()
      break
  }
}
onMounted(() => {
  window.addEventListener('message', onAppMessage)
})
onBeforeUnmount(() => {
  window.removeEventListener('message', onAppMessage)
})

watch(
  [props.resource],
  ([newResource], [oldResource]) => {
    if (isSameResource(newResource, oldResource)) {
      return
    }

    let viewMode = 'view'

    if (!props.isReadOnly) {
      viewMode = unref(viewModeQueryValue) || 'write'
    }

    if (
      determineOpenAsPreview(unref(appName)) &&
      (isShareSpaceResource(props.space) ||
        isPublicSpaceResource(props.space) ||
        isProjectSpaceResource(props.space))
    ) {
      viewMode = 'view'
    }
    loadAppUrl.perform(viewMode)
  },
  { immediate: true, deep: true }
)
</script>
