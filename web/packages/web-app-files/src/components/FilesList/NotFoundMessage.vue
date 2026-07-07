<template>
  <div
    id="files-list-not-found-message"
    class="oc-text-center oc-flex-middle oc-flex oc-flex-center oc-flex-column"
  >
    <oc-icon name="cloud" type="div" size="xxlarge" />
    <div class="oc-text-muted oc-text-xlarge">
      <span v-translate>Resource not found</span>
    </div>
    <div class="oc-text-muted">
      <span v-translate>
        We went looking everywhere, but were unable to find the selected resource.
      </span>
    </div>
    <div class="oc-mt-s">
      <oc-button
        v-if="showSpacesButton"
        id="space-not-found-button-go-spaces"
        type="router-link"
        appearance="raw"
        :to="spacesRoute"
      >
        <span v-translate>Go to »Spaces Overview«</span>
      </oc-button>
      <oc-button
        v-if="showHomeButton"
        id="files-list-not-found-button-go-home"
        type="router-link"
        appearance="raw"
        :to="homeRoute"
      >
        <span v-translate>Go to »Personal« page</span>
      </oc-button>
      <oc-button
        v-if="showPublicLinkButton"
        id="files-list-not-found-button-reload-link"
        type="router-link"
        appearance="raw"
        :to="publicLinkRoute"
      >
        <span v-translate>Reload public link</span>
      </oc-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {
  createLocationPublic,
  createLocationSpaces,
  isLocationPublicActive,
  isLocationSpacesActive,
  useRouter,
  createFileRouteOptions
} from '@ownclouders/web-pkg'
import { SpaceResource } from '@ownclouders/web-client'

interface Props {
  space?: SpaceResource
}
const { space = null } = defineProps<Props>()
const router = useRouter()
const isProjectSpace = space?.driveType === 'project'

const showPublicLinkButton = isLocationPublicActive(router, 'files-public-link')
const showHomeButton = isLocationSpacesActive(router, 'files-spaces-generic') && !isProjectSpace
const showSpacesButton = isLocationSpacesActive(router, 'files-spaces-generic') && isProjectSpace
const homeRoute = createLocationSpaces('files-spaces-generic', {
  params: {
    driveAliasAndItem: 'personal'
  }
})
const publicLinkRoute = createLocationPublic('files-public-link', createFileRouteOptions(space, {}))
const spacesRoute = createLocationSpaces('files-spaces-projects')
</script>
<style>
#files-list-not-found-message {
  height: 75dvh;
}
</style>
