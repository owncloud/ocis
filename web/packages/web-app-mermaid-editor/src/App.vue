<template>
  <div class="mermaid-editor oc-width-1-1 oc-height-1-1">
    <mermaid-toolbar :view-mode="viewMode" @change-mode="viewMode = $event" />
    <div class="mermaid-editor-body" :class="bodyClass">
      <div class="mermaid-editor-body-editor">
        <mermaid-editor-pane
          :model-value="currentContent"
          :is-read-only="isReadOnly"
          @update:model-value="onInput"
        />
      </div>
      <div class="mermaid-editor-body-preview">
        <div v-if="previewPaused" class="mermaid-editor-preview-paused">
          <p class="mermaid-editor-preview-paused__text">
            {{
              $gettext(
                'This file is large, so the live preview is paused to keep the editor responsive.'
              )
            }}
          </p>
          <oc-button
            class="mermaid-editor-preview-render"
            appearance="filled"
            variation="primary"
            @click="showPreviewAnyway"
          >
            {{ $gettext('Show preview anyway') }}
          </oc-button>
        </div>
        <mermaid-preview-pane v-else :content="previewContent" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { Resource } from '@ownclouders/web-client'
import { AppConfigObject } from '@ownclouders/web-pkg'
import MermaidEditorPane from './components/MermaidEditorPane.vue'
import MermaidPreviewPane from './components/MermaidPreviewPane.vue'
import MermaidToolbar, { MermaidEditorViewMode } from './components/MermaidToolbar.vue'
import { isPreviewTooLarge } from './helpers/preview'

interface Props {
  applicationConfig: AppConfigObject
  currentContent: string
  isReadOnly?: boolean
  resource: Resource
}
interface Emits {
  (e: 'update:currentContent', value: string): void
}
// `currentContent` and `update:currentContent` are the contract that turns on the
// AppWrapper's WebDAV load/save, dirty tracking, Ctrl+S and unsaved-changes guard,
// same as every other AppWrapperRoute-based editor in this workspace (see
// web-app-html-editor for the reference implementation of that contract).
// `applicationConfig` and `resource` are declared so the wrapper binds them as
// props rather than as fallthrough attributes.
const { currentContent, isReadOnly = false } = defineProps<Props>()
const emit = defineEmits<Emits>()

const viewMode = ref<MermaidEditorViewMode>('split')

const bodyClass = computed(() => ({
  'mermaid-editor-body-split': viewMode.value === 'split',
  'mermaid-editor-body-editor-only': viewMode.value === 'editor',
  'mermaid-editor-body-preview-only': viewMode.value === 'preview'
}))

// Large diagrams are not auto-previewed: re-parsing the whole source on every
// settled keystroke can make mermaid's layout pass hang the tab. The user can opt
// in to render a large diagram once.
const renderLargeAnyway = ref(false)
const previewPaused = computed(() => isPreviewTooLarge(currentContent) && !renderLargeAnyway.value)

// The preview is debounced so typing does not re-render the diagram on every
// keystroke. Unlike the HTML editor, the raw source is passed straight through —
// there is no markup to wrap, mermaid.render() parses the diagram text as-is.
const previewContent = ref(previewPaused.value ? '' : (currentContent ?? ''))
let previewTimer: ReturnType<typeof setTimeout> | undefined
const schedulePreview = (value: string) => {
  if (previewTimer) {
    clearTimeout(previewTimer)
  }
  previewTimer = setTimeout(() => {
    previewContent.value = value ?? ''
  }, 250)
}

const showPreviewAnyway = () => {
  renderLargeAnyway.value = true
  // Explicit user action: render the current content immediately rather than
  // through the debounce, so the preview pane does not mount empty for 250 ms.
  if (previewTimer) {
    clearTimeout(previewTimer)
    previewTimer = undefined
  }
  previewContent.value = currentContent ?? ''
}

const onInput = (value: string) => {
  emit('update:currentContent', value)
}

// Drives the preview for both user edits (which round-trip back through the prop)
// and external content changes such as the initial WebDAV load. While paused, the
// debounced render is skipped entirely.
watch(
  () => currentContent,
  (value) => {
    // Re-arm the large-file guard on every content change. The "show anyway"
    // opt-in is scoped to the document the user explicitly approved; a later
    // change (notably an external conflict-reload) must re-pause rather than
    // silently auto-render a different, possibly hostile, large document.
    renderLargeAnyway.value = false
    if (isPreviewTooLarge(value)) {
      // Drop any queued render so a previously scheduled small-content preview
      // cannot fire after the content has grown past the limit.
      if (previewTimer) {
        clearTimeout(previewTimer)
        previewTimer = undefined
      }
      previewContent.value = ''
      return
    }
    schedulePreview(value ?? '')
  }
)

onBeforeUnmount(() => {
  if (previewTimer) {
    clearTimeout(previewTimer)
  }
})
</script>

<style lang="scss" scoped>
.mermaid-editor {
  display: flex;
  flex-direction: column;
}

.mermaid-editor-body {
  display: grid;
  grid-template-columns: 1fr 1fr;
  flex: 1 1 auto;
  min-height: 0;
  overflow: hidden;
}

.mermaid-editor-body-editor,
.mermaid-editor-body-preview {
  min-width: 0;
  height: 100%;
  overflow: hidden;
}

.mermaid-editor-body-preview {
  border-left: 1px solid var(--oc-color-border);
}

.mermaid-editor-preview-paused {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--oc-space-medium);
  height: 100%;
  padding: var(--oc-space-large);
  text-align: center;
  color: var(--oc-color-text-muted);
}

.mermaid-editor-body-split {
  grid-template-columns: 1fr 1fr;
}

.mermaid-editor-body-editor-only {
  grid-template-columns: 1fr 0;

  .mermaid-editor-body-preview {
    border-left: none;
  }
}

.mermaid-editor-body-preview-only {
  grid-template-columns: 0 1fr;

  .mermaid-editor-body-preview {
    border-left: none;
  }
}
</style>
