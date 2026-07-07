<template>
  <span class="oc-resource-size" v-text="formattedSize" />
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { formatFileSize } from '../../helpers'
import { useGettext } from 'vue3-gettext'

/**
 * Displays a formatted resource size
 */

interface Props {
  size: string | number
}
const { size } = defineProps<Props>()
const { current: currentLanguage } = useGettext()
const formattedSize = computed(() => formatFileSize(size, currentLanguage))
</script>
