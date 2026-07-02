<template>
  <main id="activities">
    <router-view />
  </main>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'LayoutContainer'
})
</script>

<style lang="scss">
#activities {
  overflow: auto;
  padding: var(--oc-space-medium) !important;
}
</style>
