<template>
  <div class="html-editor oc-width-1-1 oc-height-1-1">
    <html-toolbar :view-mode="viewMode" @change-mode="viewMode = $event" />
    <div class="html-editor-body" :class="bodyClass">
      <div class="html-editor-body-editor">
        <html-editor-pane
          :model-value="currentContent"
          :is-read-only="isReadOnly"
          @update:model-value="onInput"
        />
      </div>
      <div class="html-editor-body-preview">
        <div v-if="previewPaused" class="html-editor-preview-paused">
          <p class="html-editor-preview-paused__text">
            {{
              $gettext(
                'This file is large, so the live preview is paused to keep the editor responsive.'
              )
            }}
          </p>
          <oc-button
            class="html-editor-preview-render"
            appearance="filled"
            variation="primary"
            @click="showPreviewAnyway"
          >
            {{ $gettext('Show preview anyway') }}
          </oc-button>
        </div>
        <html-preview-pane v-else :content="previewContent" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { Resource } from '@ownclouders/web-client'
import { AppConfigObject } from '@ownclouders/web-pkg'
import HtmlEditorPane from './components/HtmlEditorPane.vue'
import HtmlPreviewPane from './components/HtmlPreviewPane.vue'
import HtmlToolbar, { HtmlEditorViewMode } from './components/HtmlToolbar.vue'
import { isPreviewTooLarge, wrapWithPreviewCsp } from './helpers/preview'

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
// AppWrapper's WebDAV load/save, dirty tracking, Ctrl+S and unsaved-changes guard
// (see DECISIONS.md D5/D6/D7). `applicationConfig` and `resource` are declared so
// the wrapper binds them as props rather than as fallthrough attributes.
const { currentContent, isReadOnly = false } = defineProps<Props>()
const emit = defineEmits<Emits>()

const viewMode = ref<HtmlEditorViewMode>('split')

const bodyClass = computed(() => ({
  'html-editor-body-split': viewMode.value === 'split',
  'html-editor-body-editor-only': viewMode.value === 'editor',
  'html-editor-body-preview-only': viewMode.value === 'preview'
}))

// Large documents are not auto-previewed: re-parsing the whole document into the
// iframe on every settled edit (and a hostile script in it) can hang the tab.
// The user can opt in to render a large file once. See DECISIONS.md D3 and SECURITY-REVIEW.md.
const renderLargeAnyway = ref(false)
const previewPaused = computed(() => isPreviewTooLarge(currentContent) && !renderLargeAnyway.value)

// The preview is debounced so typing does not reload the iframe on every
// keystroke. The content is wrapped with a strict, iframe-scoped CSP before it
// reaches the (sandboxed, opaque-origin) preview. Both panes stay mounted across
// view-mode switches (CSS grid collapses the hidden column) so the editor keeps
// its cursor/undo.
const previewContent = ref(previewPaused.value ? '' : wrapWithPreviewCsp(currentContent ?? ''))
let previewTimer: ReturnType<typeof setTimeout> | undefined
const schedulePreview = (value: string) => {
  if (previewTimer) {
    clearTimeout(previewTimer)
  }
  previewTimer = setTimeout(() => {
    previewContent.value = wrapWithPreviewCsp(value ?? '')
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
  previewContent.value = wrapWithPreviewCsp(currentContent ?? '')
}

const onInput = (value: string) => {
  emit('update:currentContent', value)
}

// Drives the preview for both user edits (which round-trip back through the prop)
// and external content changes such as the initial WebDAV load. While paused, the
// expensive wrap/render is skipped entirely.
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
.html-editor {
  display: flex;
  flex-direction: column;
}

.html-editor-body {
  display: grid;
  grid-template-columns: 1fr 1fr;
  flex: 1 1 auto;
  min-height: 0;
  overflow: hidden;
}

.html-editor-body-editor,
.html-editor-body-preview {
  min-width: 0;
  height: 100%;
  overflow: hidden;
}

.html-editor-body-preview {
  border-left: 1px solid var(--oc-color-border);
}

.html-editor-preview-paused {
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

.html-editor-body-split {
  grid-template-columns: 1fr 1fr;
}

.html-editor-body-editor-only {
  grid-template-columns: 1fr 0;

  .html-editor-body-preview {
    border-left: none;
  }
}

.html-editor-body-preview-only {
  grid-template-columns: 0 1fr;

  .html-editor-body-preview {
    border-left: none;
  }
}
</style>
