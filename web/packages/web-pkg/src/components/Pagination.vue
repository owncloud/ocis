<template>
  <oc-pagination
    v-if="pages > 1"
    :pages="pages"
    :current-page="props.currentPage"
    :max-displayed="3"
    :current-route="$route"
    class="files-pagination oc-flex oc-flex-center oc-my-s"
  />
</template>

<script lang="ts" setup>
import { watch } from 'vue'

interface Props {
  pages: number
  currentPage: number
}

const props = defineProps<Props>()
watch(
  () => props.currentPage,
  () => {
    document.getElementsByClassName('files-view-wrapper')[0]?.scrollTo(0, 0)
  }
)
</script>
