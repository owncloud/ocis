<template>
  <oc-emoji-picker :theme="theme" @emoji-select="onEmojiSelect" />
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import { useThemeStore } from '../../composables'
import { storeToRefs } from 'pinia'

interface Emits {
  (e: 'confirm', value: string): void
}
const emit = defineEmits<Emits>()
const themeStore = useThemeStore()
const { currentTheme } = storeToRefs(themeStore)

const theme = computed(() => {
  return unref(currentTheme).isDark ? 'dark' : 'light'
})

const onEmojiSelect = (emoji: string) => {
  emit('confirm', emoji)
}
</script>
