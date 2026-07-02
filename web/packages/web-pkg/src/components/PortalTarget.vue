<template>
  <portal-target-vue v-bind="properties" />
</template>

<script lang="ts" setup>
import { computed, onMounted } from 'vue'
import { eventBus } from '../services'
import { PortalTargetEventTopics } from '../composables/portalTarget'
import { PortalTarget as PortalTargetVue } from 'portal-vue'

const props = defineProps({
  ...PortalTargetVue.props
})

const properties = computed<typeof PortalTargetVue.props>(() => props)

onMounted(() => {
  eventBus.publish(PortalTargetEventTopics.mounted, props)
})
</script>
