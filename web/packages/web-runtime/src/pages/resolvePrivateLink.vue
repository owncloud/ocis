<template>
  <main
    class="oc-link-resolve oc-height-viewport oc-flex oc-flex-column oc-flex-center oc-flex-middle"
  >
    <div class="oc-card oc-text-center oc-width-large">
      <template v-if="loading">
        <div class="oc-card-header">
          <h2 key="private-link-loading" class="oc-link-resolve-loading">
            <span v-text="$gettext('Resolving private link…')" />
          </h2>
        </div>
        <div class="oc-card-body">
          <oc-spinner :aria-hidden="true" />
        </div>
      </template>
      <template v-else-if="errorMessage">
        <div class="oc-card-header oc-link-resolve-error-title">
          <h2 key="private-link-error">
            <span v-text="$gettext('An error occurred while resolving the private link')" />
          </h2>
        </div>
        <div class="oc-card-body oc-link-resolve-error-message">
          <p class="oc-text-xlarge">{{ errorMessage }}</p>
        </div>
      </template>
    </div>
    <oc-button
      v-if="isUnacceptedShareError"
      type="router-link"
      variation="primary"
      appearance="filled"
      target="_blank"
      class="oc-mt-m oc-text-center oc-width-medium"
      :to="sharedWithMeRoute"
    >
      <span class="text" v-text="openSharedWithMeLabel" />
    </oc-button>
  </main>
</template>

<script lang="ts">
import {
  useRouteParam,
  useRouter,
  queryItemAsString,
  useRouteQuery,
  createLocationSpaces,
  createLocationShares,
  useClientService,
  useCapabilityStore
} from '@ownclouders/web-pkg'
import { unref, defineComponent, computed, onMounted, ref, Ref } from 'vue'
import { dirname } from 'path'
import { createFileRouteOptions, useGetResourceContext } from '@ownclouders/web-pkg'
import { useTask } from 'vue-concurrency'
import { isShareSpaceResource, Resource, SHARE_JAIL_ID } from '@ownclouders/web-client'
import { RouteLocationNamedRaw } from 'vue-router'
import { useGettext } from 'vue3-gettext'

export default defineComponent({
  name: 'ResolvePrivateLink',
  setup() {
    const router = useRouter()
    const id = useRouteParam('fileId')
    const scope = useRouteParam('scope')
    const { $gettext } = useGettext()
    const clientService = useClientService()
    const capabilityStore = useCapabilityStore()

    const resource: Ref<Resource> = ref()
    const sharedParentResource: Ref<Resource> = ref()
    const isUnacceptedShareError = ref(false)

    const { getResourceContext } = useGetResourceContext()

    const openWithDefaultAppQuery = useRouteQuery('openWithDefaultApp')
    const openWithDefaultApp = computed(() => queryItemAsString(unref(openWithDefaultAppQuery)))

    const detailsQuery = useRouteQuery('details')
    const details = computed(() => {
      return queryItemAsString(unref(detailsQuery))
    })

    onMounted(() => {
      resolvePrivateLinkTask.perform(queryItemAsString(unref(id)))
    })

    const resolvePrivateLinkTask = useTask(function* (signal, id) {
      // vault file but non-vault /f/ URL: hard-reload into /vault so the vault client
      // base URLs are picked up (they are derived from the URL at boot). The MFA
      // step-up is enforced by the authService guard on the vault scope, which also
      // fires on a router.push, so it is not the reason for the full reload.
      if (
        capabilityStore.vaultEnabled &&
        id?.split('$')[0] === capabilityStore.vaultStorageProvider &&
        unref(scope) !== 'vault'
      ) {
        window.location.replace(
          router.resolve({
            name: 'resolvePrivateLink',
            params: { scope: 'vault', fileId: id },
            query: router.currentRoute.value.query
          }).href
        )
        return
      }

      if (
        [
          `${SHARE_JAIL_ID}$${SHARE_JAIL_ID}!${SHARE_JAIL_ID}`,
          `${SHARE_JAIL_ID}$${SHARE_JAIL_ID}`
        ].includes(id)
      ) {
        return router.push(
          createLocationShares('files-shares-with-me', { params: { scope: unref(scope) } })
        )
      }

      let result: Awaited<ReturnType<typeof getResourceContext>>
      try {
        result = yield getResourceContext(id)
      } catch (e) {
        // error means the resurce is an unaccepted/unsynced share
        isUnacceptedShareError.value = true
        throw Error(e)
      }

      const { space, resource } = result
      let { path } = result

      if (!path) {
        // empty path means the user has no access to the resource or it doesn't exist
        throw new Error('The file or folder does not exist')
      }

      let resourceIsNestedInShare = false
      let isHiddenShare = false
      if (isShareSpaceResource(space)) {
        sharedParentResource.value = resource
        resourceIsNestedInShare = path !== '/'
        if (!resourceIsNestedInShare) {
          // FIXME: get drive item by id as soon as server supports it
          const driveItems = yield clientService.graphAuthenticated.driveItems.listSharedWithMe()
          const share = driveItems.find(({ remoteItem }) => remoteItem.id === resource.id)

          isHiddenShare = share?.['@UI.Hidden']
        }
      }

      let fileId: string
      let targetLocation: RouteLocationNamedRaw
      if (unref(resource).type === 'folder') {
        fileId = unref(resource).fileId
        targetLocation = createLocationSpaces('files-spaces-generic')
      } else {
        fileId = unref(resource).parentFolderId
        path = dirname(path)
        targetLocation =
          space.driveType === 'share' && !resourceIsNestedInShare
            ? createLocationShares('files-shares-with-me')
            : createLocationSpaces('files-spaces-generic')
      }

      const { params, query } = createFileRouteOptions(space, { fileId, path })
      const openWithDefault = unref(openWithDefaultApp) !== 'false' && !unref(details)

      targetLocation.params = { ...params, scope: unref(scope) }
      targetLocation.query = {
        ...query,
        scrollTo: unref(resource).fileId,
        ...(unref(details) && { details: unref(details) }),
        ...(isHiddenShare && { 'q_share-visibility': 'hidden' }),
        ...(openWithDefault && { openWithDefaultApp: 'true' })
      }

      router.push(targetLocation)
    })

    const loading = computed(() => {
      return !resolvePrivateLinkTask.last || resolvePrivateLinkTask.isRunning
    })

    const sharedWithMeRoute = computed(() => ({
      name: 'files-shares-with-me',
      params: { scope: unref(scope) }
    }))

    const openSharedWithMeLabel = computed(() => {
      return $gettext('Open "Shared with me"')
    })

    const errorMessage = computed(() => {
      if (unref(isUnacceptedShareError)) {
        return $gettext(
          'The link you are trying to access is invalid or you do not have permission to view the content. Please check the link for any errors or contact the person who shared it for assistance.'
        )
      }

      if (resolvePrivateLinkTask.isError) {
        return resolvePrivateLinkTask.last.error.message
      }
      return null
    })

    return {
      errorMessage,
      loading,
      resource,
      isUnacceptedShareError,
      sharedWithMeRoute,
      openSharedWithMeLabel,

      // HACK: for unit tests
      resolvePrivateLinkTask
    }
  }
})
</script>

<style lang="scss">
.oc-link-resolve {
  .oc-card {
    background: var(--oc-color-background-highlight);
    border-radius: 15px;
  }

  .oc-card-header h2 {
    margin: 0;
  }
}
</style>
