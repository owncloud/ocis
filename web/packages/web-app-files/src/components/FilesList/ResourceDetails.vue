<template>
  <div class="resource-details oc-flex oc-flex-column oc-flex-middle">
    <div class="oc-width-1-3@l oc-width-1-2@m oc-width-3-4">
      <file-info />
      <file-details class="oc-mb" />
      <file-actions />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, provide, unref } from 'vue'
import { Resource, SpaceResource } from '@ownclouders/web-client'

import FileActions from '../SideBar/Actions/FileActions.vue'
import FileDetails from '../SideBar/Details/FileDetails.vue'
import { FileInfo, useOpenWithDefaultApp } from '@ownclouders/web-pkg'
import { useRouteQuery } from '@ownclouders/web-pkg'

interface Props {
  singleResource?: Resource
  space?: SpaceResource
}
const { singleResource = null, space = null } = defineProps<Props>()
provide(
  'resource',
  computed(() => singleResource)
)
provide(
  'space',
  computed(() => space)
)

const { openWithDefaultApp } = useOpenWithDefaultApp()
const openWithDefaultAppQuery = useRouteQuery('openWithDefaultApp')
if (unref(openWithDefaultAppQuery) === 'true') {
  openWithDefaultApp({
    space: space,
    resource: singleResource
  })
}
</script>
