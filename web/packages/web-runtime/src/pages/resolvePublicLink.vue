<template>
  <div
    class="oc-link-resolve oc-height-viewport oc-flex oc-flex-column oc-flex-center oc-flex-middle"
  >
    <div class="oc-card oc-text-center oc-width-large">
      <template v-if="errorMessage">
        <div class="oc-card-header oc-link-resolve-error-title">
          <h2 key="public-link-error">
            <span v-text="$gettext('An error occurred while loading the public link')" />
          </h2>
        </div>
        <div class="oc-card-body oc-link-resolve-error-message">
          <p class="oc-text-xlarge">{{ errorMessage }}</p>
        </div>
      </template>
      <template v-else-if="isPasswordRequired">
        <form @submit.prevent="resolvePublicLinkTask.perform(true)">
          <div class="oc-card-header">
            <h2>
              <span v-text="$gettext('This resource is password-protected')" />
            </h2>
          </div>
          <div class="oc-card-body">
            <oc-text-input
              ref="passwordInput"
              v-model="password"
              :error-message="wrongPasswordMessage"
              :label="passwordFieldLabel"
              type="password"
              class="oc-mb-s"
            />
            <oc-button
              variation="primary"
              appearance="filled"
              class="oc-login-authorize-button"
              :disabled="!password || resolvePublicLinkTask.isRunning"
              :show-spinner="resolvePublicLinkTask.isRunning"
              submit="submit"
            >
              <span v-text="$gettext('Continue')" />
            </oc-button>
          </div>
        </form>
      </template>
      <template v-else>
        <div class="oc-card-header">
          <h2 key="public-link-loading">
            <span v-text="$gettext('Loading public link…')" />
          </h2>
        </div>
        <div class="oc-card-body">
          <oc-spinner :aria-hidden="true" />
        </div>
      </template>
      <div class="oc-card-footer oc-pt-rm">
        <p>{{ footerSlogan }}</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { DavHttpError, SharePermissionBit } from '@ownclouders/web-client'
import { authService } from '../services/auth'

import {
  queryItemAsString,
  useAuthStore,
  useClientService,
  useConfigStore,
  useRoute,
  useRouteParam,
  useRouteQuery,
  useRouter,
  useSpacesStore,
  useThemeStore
} from '@ownclouders/web-pkg'
import { useTask } from 'vue-concurrency'
import { ref, unref, computed, onMounted } from 'vue'
import {
  buildPublicSpaceResource,
  isPublicSpaceResource,
  PublicSpaceResource
} from '@ownclouders/web-client'
import { useGettext } from 'vue3-gettext'
import { urlJoin } from '@ownclouders/web-client'
import { RouteLocationNamedRaw } from 'vue-router'
import { dirname } from 'path'
import { storeToRefs } from 'pinia'

const configStore = useConfigStore()
const clientService = useClientService()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const { $gettext } = useGettext()
const token = useRouteParam('token')
const redirectUrl = useRouteQuery('redirectUrl')
const themeStore = useThemeStore()
const spacesStore = useSpacesStore()

const { currentTheme } = storeToRefs(themeStore)
const password = ref('')

const isOcmLink = computed(() => {
  const split = unref(route).path.split('/')?.[1]
  return split === 'o'
})

const publicLinkType = computed(() => (unref(isOcmLink) ? 'ocm' : 'public-link'))

const publicLinkSpace = computed(() =>
  buildPublicSpaceResource({
    id: unref(token),
    driveType: 'public',
    publicLinkType: unref(publicLinkType),
    ...(unref(password) && { publicLinkPassword: unref(password) })
  })
)

const item = computed(() => queryItemAsString(unref(route)?.params?.driveAliasAndItem))

const detailsQuery = useRouteQuery('details')
const details = computed(() => queryItemAsString(unref(detailsQuery)))

const loadedSpace = ref<PublicSpaceResource>()
const isPasswordRequired = ref(false)
const isInternalLink = ref(false)

const loadPublicSpaceTask = useTask(function* (signal) {
  try {
    loadedSpace.value = yield clientService.webdav.getFileInfo(
      unref(publicLinkSpace),
      {},
      { signal }
    )
  } catch (error) {
    const err = error as DavHttpError

    if (err.statusCode === 401) {
      if (err.errorCode === 'ERR_MISSING_BASIC_AUTH') {
        isPasswordRequired.value = true
      }

      if (err.errorCode === 'ERR_MISSING_BEARER_AUTH') {
        isInternalLink.value = true
      }

      return
    }
    if (err.statusCode === 404) {
      throw new Error($gettext('The resource could not be located, it may not exist anymore.'))
    }
    throw err
  }
})

const verifyPasswordTask = useTask(function* (signal) {
  try {
    loadedSpace.value = yield clientService.webdav.getFileInfo(
      unref(publicLinkSpace),
      {},
      { signal }
    )
    if (!isPublicSpaceResource(unref(loadedSpace))) {
      const e: any = new Error($gettext('The resource is not a public link.'))
      e.resource = unref(loadedSpace)
      throw e
    }
  } catch (e) {
    if (e.statusCode === 401) {
      throw e
    }
    throw new Error($gettext('The resource could not be located, it may not exist anymore.'))
  }
})
const wrongPassword = computed(() => {
  if (verifyPasswordTask.isError) {
    return verifyPasswordTask.last.error.statusCode === 401
  }
  return false
})

const resolvePublicLinkTask = useTask(function* (signal, passwordRequired: boolean) {
  if (unref(isOcmLink) && !configStore.options.ocm.openRemotely) {
    throw new Error($gettext('Opening files from remote is disabled'))
  }

  if (unref(isInternalLink)) {
    router.push({ name: 'login', query: { redirectUrl: `/i/${unref(token)}` } })
    return
  }

  yield authService.resolvePublicLink(
    unref(token),
    passwordRequired,
    passwordRequired ? unref(password) : '',
    unref(publicLinkType)
  )

  if (passwordRequired) {
    try {
      yield verifyPasswordTask.perform()
    } catch (e) {
      authStore.clearPublicLinkContext()
      console.error(e, e.resource)
      throw e
    }
  }

  const url = queryItemAsString(unref(redirectUrl))
  if (url) {
    router.push({ path: url })
    return
  }

  if (unref(loadedSpace).publicLinkPermission === SharePermissionBit.Create) {
    router.push({
      name: 'files-public-upload',
      params: { token: unref(token) },
      query: { fileId: unref(publicLinkSpace).fileId }
    })
    return
  }

  let scrollTo: string
  let fileId: string
  let path: string

  if (['folder', 'space'].includes(unref(loadedSpace).type)) {
    fileId = unref(loadedSpace).fileId
    path = unref(item)
  } else {
    fileId = unref(loadedSpace).parentFolderId
    scrollTo = unref(loadedSpace).fileId
    path = dirname(unref(item))
  }

  spacesStore.upsertSpace(unref(loadedSpace))

  const driveAliasAndItem = urlJoin(unref(isOcmLink) ? `ocm/` : `public/`, unref(token), path)
  const targetLocation: RouteLocationNamedRaw = {
    name: 'files-public-link',
    query: {
      openWithDefaultApp: 'true',
      ...(!!fileId && { fileId }),
      ...(!!scrollTo && { scrollTo }),
      ...(unref(details) && { details: unref(details) })
    },
    params: {
      driveAliasAndItem
    }
  }

  router.push(targetLocation)
})

const errorMessage = computed<string>(() => {
  if (resolvePublicLinkTask.isError && resolvePublicLinkTask.last.error.statusCode !== 401) {
    return resolvePublicLinkTask.last.error.message
  }

  if (loadPublicSpaceTask.isError) {
    return loadPublicSpaceTask.last.error.message
  }
  return null
})

onMounted(async () => {
  try {
    if (unref(isOcmLink)) {
      await resolvePublicLinkTask.perform(false)
      return
    }

    await loadPublicSpaceTask.perform()

    if (!unref(isPasswordRequired)) {
      await resolvePublicLinkTask.perform(false)
    }
  } catch (e) {
    console.error(e)
  }
})

const footerSlogan = computed(() => currentTheme.value.common.slogan)
const passwordFieldLabel = computed(() => {
  return $gettext('Enter password for public link')
})
const wrongPasswordMessage = computed(() => {
  if (unref(wrongPassword)) {
    return $gettext('Incorrect password')
  }
  return null
})
</script>

<style lang="scss">
.oc-link-resolve {
  .oc-card {
    background: var(--oc-color-background-highlight);
    border-radius: 15px;
  }

  .oc-text-input-message {
    justify-content: center;
  }

  .oc-card-header h2,
  .oc-card-footer p {
    margin: 0;
  }
}
</style>
