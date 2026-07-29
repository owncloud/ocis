<template>
  <div class="mermaid-preview-pane">
    <div v-if="hasError" class="mermaid-preview-pane-error">
      <oc-icon name="error-warning" fill-type="line" variation="danger" size="large" />
      <p class="mermaid-preview-pane-error__text">
        {{ $gettext('This diagram could not be rendered. Check the syntax.') }}
      </p>
      <pre v-if="errorDetail" class="mermaid-preview-pane-error__detail">{{ errorDetail }}</pre>
    </div>
    <!-- eslint-disable-next-line vue/no-v-html -->
    <div v-else class="mermaid-preview-pane-diagram" v-html="renderedSvg" />
  </div>
</template>

<script lang="ts" setup>
import { onBeforeUnmount, unref, watch, ref } from 'vue'
import mermaid from 'mermaid'
import DOMPurify from 'dompurify'
import { useThemeStore } from '@ownclouders/web-pkg'

interface Props {
  content: string
}
const { content } = defineProps<Props>()

const themeStore = useThemeStore()
const isDark = () => Boolean(unref(themeStore.currentTheme)?.isDark)

const renderedSvg = ref('')
// `hasError` drives the translated, static heading (rendered via the template's
// global `$gettext`, matching every other component in this workspace — leaf
// components here never call `useGettext()` from script). `errorDetail` carries
// only the raw, untranslated diagnostic (mermaid's own thrown message), and is
// empty for a plain failed parse where mermaid gives no message at all.
const hasError = ref(false)
const errorDetail = ref('')

let renderCounter = 0
// Bumped on every render call and checked after each async step, so a slow,
// now-superseded render (stale parse/render promise from a previous keystroke)
// can never clobber the result of a render started after it.
let renderToken = 0

const configureMermaid = () => {
  mermaid.initialize({
    startOnLoad: false,
    securityLevel: 'strict',
    theme: isDark() ? 'dark' : 'default'
  })
}

// Same defense-in-depth sanitize pass web-pkg's TextEditor.vue applies to
// markdown-embedded mermaid output: mermaid's `securityLevel: 'strict'` already
// sanitizes label-derived HTML internally, but some diagram types embed real HTML
// inside <foreignObject>, so the rendered SVG is run through DOMPurify again
// before it reaches `v-html`.
const sanitize = (svg: string) => DOMPurify.sanitize(svg, { ADD_ATTR: ['target'], ADD_TAGS: ['foreignObject'] })

const renderDiagram = async (source: string) => {
  const token = ++renderToken
  if (!source?.trim()) {
    renderedSvg.value = ''
    hasError.value = false
    errorDetail.value = ''
    return
  }

  configureMermaid()

  // `mermaid.parse` with `suppressErrors` validates the syntax without throwing
  // and, importantly, without mermaid's render path side effect of injecting an
  // error placeholder node straight into `document.body` on a failed parse. Only
  // syntax that parses cleanly is ever handed to `render()`.
  const isValid = await mermaid.parse(source, { suppressErrors: true })
  if (token !== renderToken) {
    return
  }
  if (!isValid) {
    renderedSvg.value = ''
    hasError.value = true
    errorDetail.value = ''
    return
  }

  try {
    const { svg } = await mermaid.render(`mermaid-preview-${renderCounter++}`, source)
    if (token !== renderToken) {
      return
    }
    renderedSvg.value = sanitize(svg)
    hasError.value = false
    errorDetail.value = ''
  } catch (error) {
    if (token !== renderToken) {
      return
    }
    renderedSvg.value = ''
    hasError.value = true
    errorDetail.value = error instanceof Error ? error.message : String(error)
  }
}

watch(
  () => content,
  (value) => {
    renderDiagram(value ?? '')
  },
  { immediate: true }
)
watch(
  () => isDark(),
  () => {
    renderDiagram(content ?? '')
  }
)

onBeforeUnmount(() => {
  // Invalidate any render still in flight so its promise continuation is a no-op
  // after the component (and its refs) are gone.
  renderToken++
})
</script>

<style scoped>
.mermaid-preview-pane {
  height: 100%;
  width: 100%;
  overflow: auto;
  padding: var(--oc-space-medium);
}

.mermaid-preview-pane-diagram :deep(svg) {
  max-width: 100%;
  height: auto;
}

.mermaid-preview-pane-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--oc-space-small);
  height: 100%;
  text-align: center;
  color: var(--oc-color-text-muted);
}

.mermaid-preview-pane-error__detail {
  max-width: 100%;
  overflow: auto;
  white-space: pre-wrap;
  font-family: Consolas, 'Liberation Mono', Menlo, monospace;
  color: var(--oc-color-swatch-danger-default);
}
</style>
