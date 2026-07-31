<template>
  <div class="app-tags">
    <oc-tag
      v-for="tag in app.tags"
      :key="`app-tag-${app.id}-${tag}`"
      data-testid="tag-button"
      size="medium"
      class="oc-text-nowrap"
      type="button"
      @click="emitClick(tag)"
    >
      <span class="mark-element">{{ tag }}</span>
    </oc-tag>
  </div>
</template>
<script lang="ts" setup>
import { App } from '../types'

interface Props {
  app?: App
}
interface Emits {
  (e: 'click', value: string): void
}
const { app = undefined } = defineProps<Props>()
const emit = defineEmits<Emits>()
const emitClick = (tagName: string) => {
  emit('click', tagName)
}
</script>

<style lang="scss">
.app-tags {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}
</style>
