<template>
  <app-loading-spinner v-if="isLoading" v-bind="attrs" />
  <template v-else>
    <app-banner v-if="fileId" :file-id="fileId" v-bind="attrs"></app-banner>
    <drive-redirect v-if="!space" :drive-alias-and-item="driveAliasAndItem" v-bind="attrs" />
    <generic-trash v-else-if="isTrashRoute" :space="space" :item-id="itemId" v-bind="attrs" />
    <generic-space v-else :space="space" :item="item" :item-id="itemId" v-bind="attrs" />
  </template>
</template>

<script lang="ts" setup>
import DriveRedirect from './DriveRedirect.vue'
import GenericSpace from './GenericSpace.vue'
import GenericTrash from './GenericTrash.vue'

import { computed, onMounted, ref, unref, useAttrs } from 'vue'
import {
  queryItemAsString,
  useAuthStore,
  useClientService,
  useConfigStore,
  useDriveResolver,
  useGetMatchingSpace,
  useRouteParam,
  useRouteQuery,
  useRouter,
  createLocationSpaces,
  isLocationTrashActive,
  useActiveLocation,
  locationPublicUpload,
  createFileRouteOptions,
  AppLoadingSpinner,
  AppBanner
} from '@ownclouders/web-pkg'
import {
  isPublicSpaceResource,
  PublicSpaceResource,
  SharePermissionBit,
  SpaceResource
} from '@ownclouders/web-client'
import { dirname } from 'path'

const authStore = useAuthStore()
const configStore = useConfigStore()
const clientService = useClientService()
const router = useRouter()
const driveAliasAndItem = useRouteParam('driveAliasAndItem')
const isTrashRoute = useActiveLocation(isLocationTrashActive, 'files-trash-generic')
const { space, item, itemId, loading } = useDriveResolver({ driveAliasAndItem })
const { getInternalSpace } = useGetMatchingSpace()

const attrs = useAttrs()
const isProcessing = ref(true)
const isLoading = computed(() => {
  return unref(isProcessing) || unref(loading)
})

const fileIdQueryItem = useRouteQuery('fileId')
const fileId = computed(() => {
  return queryItemAsString(unref(fileIdQueryItem))
})

const getSpaceResource = async (): Promise<SpaceResource> => {
  const spaceValue = unref(space)
  try {
    return (await clientService.webdav.getFileInfo(spaceValue)) as SpaceResource
  } catch (e) {
    console.error(e)
    return spaceValue
  }
}

const resolveToInternalLocation = async (path: string) => {
  const internalSpace = getInternalSpace(unref(fileId).split('!')[0])
  if (internalSpace) {
    const resource = await clientService.webdav.getFileInfo(
      internalSpace,
      { path },
      { headers: { Authorization: `Bearer ${authStore.accessToken}` } }
    )

    const resourceId = resource.type !== 'folder' ? resource.parentFolderId : resource.fileId
    const resourcePath = resource.type !== 'folder' ? dirname(path) : path
    space.value = internalSpace
    item.value = resourcePath

    const { params, query } = createFileRouteOptions(internalSpace, {
      fileId: resourceId,
      path: resourcePath
    })
    return router.push(
      createLocationSpaces('files-spaces-generic', {
        params,
        query: {
          ...query,
          scrollTo: unref(resource).fileId,
          openWithDefaultApp: 'true'
        }
      })
    )
  }

  // no internal space found -> share -> resolve via private link as it holds all the necessary logic
  return router.push({
    name: 'resolvePrivateLink',
    params: { fileId: unref(fileId) },
    query: {
      openWithDefaultApp: 'true'
    }
  })
}

onMounted(async () => {
  if (!unref(driveAliasAndItem) && unref(fileId)) {
    return router.push({
      name: 'resolvePrivateLink',
      params: { fileId: unref(fileId) },
      query: {
        openWithDefaultApp: 'true'
      }
    })
  }

  const spaceValue = unref(space)
  if (spaceValue && isPublicSpaceResource(spaceValue)) {
    const isRunningOnEos = configStore.options.runningOnEos
    if (authStore.userContextReady && unref(fileId) && !isRunningOnEos) {
      try {
        const path = await clientService.webdav.getPathForFileId(unref(fileId), {
          headers: { Authorization: `Bearer ${authStore.accessToken}` }
        })
        await resolveToInternalLocation(path)
        isProcessing.value = false
        return
      } catch {
        // getPathForFileId failed means the user doesn't have internal access to the resource
      }
    }

    /**
     * This is to make sure that an already resolved public link still resolves correctly
     * upon reload if the link type has been changed to "Uploader" meanwhile.
     * If the space ids differ, it means we're coming from the resolvePublicLink page
     * that already feetched the space. Hence the fileId and the id differ.
     * It also means the resolvePublicLink page already handled a link of type "Uploader".
     *
     * Ideally we would redirect the user via the resolvePublicLink page, but we didn't
     * find an easy way to do that.
     **/
    if (spaceValue.fileId === spaceValue.id) {
      const publicSpace = (await getSpaceResource()) as PublicSpaceResource

      // FIXME: check for type once https://github.com/owncloud/ocis/issues/8740 is resolved
      if (publicSpace.publicLinkPermission === SharePermissionBit.Create) {
        router.push({
          name: locationPublicUpload.name,
          params: { token: spaceValue.id.toString() }
        })
      }
    }
  }

  isProcessing.value = false
})
</script>
