<template>
  <div id="oc-trash-no-selection" class="oc-text-center oc-mt-xl">
    <oc-icon size="xxlarge" name="settings-4" fill-type="fill" />
    <p v-text="$gettext('Select a resource from the left sidebar to manage it')" />
  </div>
</template>
<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'DetailsPanel'
})
</script>
