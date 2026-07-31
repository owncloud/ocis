<template>
  <div class="html-editor-toolbar oc-flex oc-flex-middle oc-px-s oc-py-xs">
    <div class="oc-button-group" role="group" :aria-label="$gettext('View mode')">
      <oc-button
        v-for="mode in modes"
        :key="mode.name"
        v-oc-tooltip="$gettext(mode.label)"
        :class="`html-editor-viewmode-${mode.name}`"
        :appearance="viewMode === mode.name ? 'filled' : 'outline'"
        :aria-label="$gettext(mode.label)"
        :aria-pressed="viewMode === mode.name"
        variation="primary"
        size="small"
        @click="$emit('changeMode', mode.name)"
      >
        <oc-icon :name="mode.icon" fill-type="line" size="small" variation="inherit" />
      </oc-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
export type HtmlEditorViewMode = 'editor' | 'split' | 'preview'

interface Props {
  viewMode: HtmlEditorViewMode
}
interface Emits {
  (e: 'changeMode', value: HtmlEditorViewMode): void
}
defineProps<Props>()
defineEmits<Emits>()

const modes: { name: HtmlEditorViewMode; label: string; icon: string }[] = [
  { name: 'editor', label: 'Editor', icon: 'code-s-slash' },
  { name: 'split', label: 'Split', icon: 'layout-column' },
  { name: 'preview', label: 'Preview', icon: 'eye' }
]
</script>

<style scoped>
.html-editor-toolbar {
  flex: 0 0 auto;
  border-bottom: 1px solid var(--oc-color-border);
  background-color: var(--oc-color-background-default);
}
</style>
