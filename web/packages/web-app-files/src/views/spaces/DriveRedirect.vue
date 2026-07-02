<template>
  <div class="oc-flex oc-width-1-1">
    <no-content-message
      v-if="spaceNotFound"
      id="files-space-not-found"
      class="oc-width-1-1"
      icon="layout-grid"
    >
      <template #message>
        <span v-text="$gettext('Space not found')" />
      </template>

      <template #callToAction>
        <oc-button
          id="space-not-found-button-go-spaces"
          type="router-link"
          appearance="raw"
          class="oc-mt-s"
          :to="spacesRoute"
        >
          <span v-translate>Go to »Spaces Overview«</span>
        </oc-button>
      </template>
    </no-content-message>

    <app-loading-spinner />
  </div>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import { NoContentMessage, useRoute, useRouter, useSpacesStore } from '@ownclouders/web-pkg'
import { AppLoadingSpinner } from '@ownclouders/web-pkg'
import { urlJoin } from '@ownclouders/web-client'
import { createFileRouteOptions } from '@ownclouders/web-pkg'
import { createLocationSpaces } from '@ownclouders/web-pkg'
import { RouteLocationRaw } from 'vue-router'

// 'personal/home' is used as personal drive alias from static contexts
// (i.e. places where we can't load the actual personal space)
const fakePersonalDriveAlias = 'personal/home'

const { driveAliasAndItem = '' } = defineProps<{
  driveAliasAndItem?: string
}>()

const router = useRouter()
const route = useRoute()
const spacesStore = useSpacesStore()

const personalSpace = computed(() => {
  return spacesStore.spaces.find((space) => space.driveType === 'personal')
})

const spacesRoute = computed(() => createLocationSpaces('files-spaces-projects'))

const spaceNotFound = computed(
  () =>
    driveAliasAndItem !== '' &&
    !driveAliasAndItem.startsWith(fakePersonalDriveAlias) &&
    !driveAliasAndItem.startsWith('personal')
)

if (!unref(spaceNotFound)) {
  if (!unref(personalSpace)) {
    router.replace(unref(spacesRoute))
  } else {
    const itemPath = driveAliasAndItem.startsWith(fakePersonalDriveAlias)
      ? urlJoin(driveAliasAndItem.slice(fakePersonalDriveAlias.length))
      : '/'

    const { params, query } = createFileRouteOptions(unref(personalSpace), {
      path: itemPath
    })

    const { fullPath, ...routeWithoutFullPath } = unref(route)

    router
      .replace({
        ...routeWithoutFullPath,
        path: fullPath,
        params: {
          ...routeWithoutFullPath.params,
          ...params
        },
        query
      } as RouteLocationRaw)
      // avoid NavigationDuplicated error in console
      .catch(() => {})
  }
}
</script>
