<template>
  <div id="oc-no-selection" class="oc-text-center oc-mt-xl">
    <oc-icon size="xxlarge" name="drag-drop" fill-type="line" />
    <p data-testid="selectedFilesText" v-text="selectedFilesString" />
  </div>
</template>
<script lang="ts" setup>
import { computed } from 'vue'
import { useGettext } from 'vue3-gettext'

const { $gettext } = useGettext()
const selectedFilesString = computed(() => $gettext('Select a file or folder to view details'))
</script>
