<template>
  <div ref="editorRef" class="html-editor-pane" />
</template>

<script lang="ts" setup>
import { onBeforeUnmount, onMounted, ref, unref, watch } from 'vue'
import {
  EditorView,
  drawSelection,
  highlightActiveLine,
  highlightActiveLineGutter,
  keymap,
  lineNumbers
} from '@codemirror/view'
import { Compartment, EditorState } from '@codemirror/state'
import { defaultKeymap, history, historyKeymap, indentWithTab } from '@codemirror/commands'
import { html } from '@codemirror/lang-html'
import {
  bracketMatching,
  defaultHighlightStyle,
  indentOnInput,
  syntaxHighlighting
} from '@codemirror/language'
import { useThemeStore } from '@ownclouders/web-pkg'

interface Props {
  modelValue: string
  isReadOnly?: boolean
}
interface Emits {
  (e: 'update:modelValue', value: string): void
}
const { modelValue, isReadOnly = false } = defineProps<Props>()
const emit = defineEmits<Emits>()

// The active theme drives the editor's dark/light chrome. `currentTheme` is a ref
// the dark-mode watch below tracks, so the editor reconfigures on a theme switch.
const themeStore = useThemeStore()
const isDark = () => Boolean(unref(themeStore.currentTheme)?.isDark)

const editorRef = ref<HTMLElement>()
let view: EditorView | undefined

const themeCompartment = new Compartment()
const readOnlyCompartment = new Compartment()

// Editor chrome is expressed in ODS tokens so it follows the active theme; the
// `dark` flag lets CodeMirror pick sensible defaults for selection/cursor.
const buildTheme = () =>
  EditorView.theme(
    {
      '&': {
        height: '100%',
        color: 'var(--oc-color-text-default)',
        backgroundColor: 'var(--oc-color-background-default)'
      },
      '.cm-scroller': {
        fontFamily: 'Consolas, "Liberation Mono", Menlo, monospace',
        lineHeight: '1.5',
        overflow: 'auto'
      },
      '.cm-gutters': {
        backgroundColor: 'var(--oc-color-background-muted)',
        color: 'var(--oc-color-text-muted)',
        border: 'none'
      },
      '.cm-content': { caretColor: 'var(--oc-color-text-default)' },
      '&.cm-focused': { outline: 'none' }
    },
    { dark: isDark() }
  )

const buildReadOnly = () => [
  EditorState.readOnly.of(isReadOnly),
  EditorView.editable.of(!isReadOnly)
]

onMounted(() => {
  view = new EditorView({
    parent: unref(editorRef),
    state: EditorState.create({
      doc: modelValue ?? '',
      extensions: [
        lineNumbers(),
        highlightActiveLine(),
        highlightActiveLineGutter(),
        drawSelection(),
        history(),
        indentOnInput(),
        bracketMatching(),
        syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
        html(),
        // No Ctrl+S binding on purpose: AppWrapper binds save at the document
        // level (see DECISIONS.md D6), and CodeMirror's default keymap leaves
        // Ctrl+S free so the keystroke propagates to it.
        keymap.of([...defaultKeymap, ...historyKeymap, indentWithTab]),
        themeCompartment.of(buildTheme()),
        readOnlyCompartment.of(buildReadOnly()),
        EditorView.updateListener.of((update) => {
          if (update.docChanged) {
            emit('update:modelValue', update.state.doc.toString())
          }
        })
      ]
    })
  })
})

// External content changes (initial WebDAV load, refresh-after-conflict) flow in
// via the prop. Only dispatch when the value actually differs to avoid a loop
// with our own updateListener emit.
watch(
  () => modelValue,
  (value) => {
    if (!view) {
      return
    }
    const current = view.state.doc.toString()
    if ((value ?? '') !== current) {
      view.dispatch({ changes: { from: 0, to: current.length, insert: value ?? '' } })
    }
  }
)

watch(
  () => isReadOnly,
  () => {
    view?.dispatch({ effects: readOnlyCompartment.reconfigure(buildReadOnly()) })
  }
)

watch(
  () => isDark(),
  () => {
    view?.dispatch({ effects: themeCompartment.reconfigure(buildTheme()) })
  }
)

onBeforeUnmount(() => {
  view?.destroy()
  view = undefined
})

// Exposed for unit tests to drive the editor without depending on layout.
defineExpose({ getView: () => view })
</script>

<style scoped>
.html-editor-pane {
  height: 100%;
  width: 100%;
  overflow: hidden;
}
</style>
