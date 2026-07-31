<template>
  <oc-status-indicators
    v-if="indicators.length > 0"
    v-bind="attrs"
    :resource="resource"
    :indicators="indicators"
  />
</template>

<script setup lang="ts">
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { ResourceIndicator, getIndicators } from '../../helpers'
import { computed, useAttrs } from 'vue'
import { useResourcesStore, useUserStore } from '../../composables/piniaStores'

const attrs = useAttrs()
const {
  space = null,
  resource,
  filter = null
} = defineProps<{
  space?: SpaceResource
  resource: Resource
  filter?: (indicator: ResourceIndicator) => boolean
}>()

const userStore = useUserStore()
const resourcesStore = useResourcesStore()
const indicators = computed(() => {
  const list = getIndicators({
    space,
    resource,
    ancestorMetaData: resourcesStore.ancestorMetaData,
    user: userStore.user
  })

  if (filter) {
    return list.filter(filter)
  }

  return list
})
</script>
