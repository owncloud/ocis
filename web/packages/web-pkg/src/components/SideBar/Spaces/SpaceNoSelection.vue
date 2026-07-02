<template>
  <div class="oc-mt-xl">
    <div class="oc-flex space-info oc-text-center">
      <oc-icon name="layout-grid" size="xxlarge" />
      <p v-translate>Select a space to view details</p>
    </div>
  </div>
</template>
<script lang="ts">
import { defineComponent } from 'vue'
export default defineComponent({
  name: 'SpaceNoSelection'
})
</script>
<style lang="scss">
.space-info {
  align-items: center;
  flex-direction: column;
}
</style>
