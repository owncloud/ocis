<template>
  <div id="oc-trash-no-selection" class="oc-text-center oc-mt-xl">
    <oc-icon size="xxlarge" name="delete-bin" fill-type="line" />
    <p data-testid="selectTrashText" v-text="$gettext('Select a trash to view details')" />
  </div>
</template>
<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'TrashNoSelection'
})
</script>
