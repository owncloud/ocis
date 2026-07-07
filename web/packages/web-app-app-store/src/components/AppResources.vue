<template>
  <ul class="oc-mb-rm oc-p-rm">
    <li v-for="resource in resources" :key="resource.label" class="app-resource-item">
      <a
        :href="resource.url"
        data-testid="resource-link"
        target="_blank"
        class="oc-flex-inline oc-flex-middle"
      >
        <oc-icon
          v-if="resource.icon"
          data-testid="resource-icon"
          :name="resource.icon"
          size="medium"
          class="oc-mr-xs"
        />
        <span data-testid="resource-label">{{ resource.label }}</span>
      </a>
    </li>
  </ul>
</template>
<script lang="ts" setup>
import { computed } from 'vue'
import { App } from '../types'
import { isEmpty } from 'lodash-es'

interface Props {
  app?: App
}
const { app = undefined } = defineProps<Props>()

const resources = computed(() => {
  return (app.resources || []).filter((resource) => {
    if (isEmpty(resource.url) || isEmpty(resource.label)) {
      return false
    }
    try {
      new URL(resource.url)
    } catch {
      return false
    }
    return true
  })
})
</script>

<style lang="scss">
.app-resource-item {
  list-style: none;
}
</style>
