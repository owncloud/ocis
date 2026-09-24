<template>
  <div id="text-editor-container" class="oc-height-1-1">
    <md-preview
      v-if="isReadOnly"
      id="space-description-preview"
      :model-value="currentContent"
      :language="languages[currentLanguage] || 'en-US'"
      :theme="theme"
      read-only
      :toolbars="[]"
      :sanitize="sanitize"
    />
    <md-editor
      v-else
      id="text-editor-component"
      :model-value="currentContent"
      :language="languages[currentLanguage] || 'en-US'"
      :theme="theme"
      :preview="isMarkdown"
      :toolbars="isMarkdown ? undefined : []"
      :footers="['markdownTotal', 0, '=', 'scrollSwitch']"
      :read-only="isReadOnly"
      :auto-focus="autoFocus"
      :sanitize="sanitize"
      :toolbars-exclude="['save', 'github']"
      @on-change="(value) => $emit('update:currentContent', value)"
      @on-upload-img="onUploadImg"
    >
      <template #defFooters>
        <span class="footer-slot">
          <span v-if="isAddingImages" class="footer-image-progress" role="status">
            {{ $gettext('Adding image…') }}
          </span>

          <span class="footer-links">
            <a
              href="https://imzbf.github.io/md-editor-v3/en-US/api#%F0%9F%AA%A1%20Shortcut%20keys"
              target="_blank"
              rel="noopener noreferrer"
              >{{
                $pgettext(
                  'A link to a list of keyboard shortcuts that can be used in the markdown editor.',
                  'Keyboard shortcuts'
                )
              }}</a
            >

            <a
              href="https://highlightjs.readthedocs.io/en/latest/supported-languages.html"
              target="_blank"
              rel="noopener noreferrer"
              >{{
                $pgettext(
                  'A link to a list of supported programming languages that can be used in the markdown editor.',
                  'Supported programming languages'
                )
              }}</a
            >
          </span>
        </span>
      </template>
    </md-editor>
    <span v-if="!isReadOnly" id="text-editor-focus-out-hint" class="oc-invisible-sr">
      {{ $gettext('Press Control+M to move focus out of the text area to the toolbar.') }}
    </span>
  </div>
</template>

<script lang="ts" setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, unref, watch } from 'vue'
import { Resource } from '@ownclouders/web-client'
import dompurify from 'dompurify'

import { config, MdEditor, MdPreview } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

import screenfull from 'screenfull'

import katex from 'katex'
import 'katex/dist/katex.min.css'

import mermaid from 'mermaid'

import highlight from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'

import * as prettier from 'prettier'
import parserMarkdown from 'prettier/plugins/markdown'

import { languageUserDefined, languages } from './l18n'

import { useGettext } from 'vue3-gettext'
import { AppConfigObject } from '../../apps'
import { useMessages, useThemeStore } from '../../composables'
import { DEFAULT_MAX_DOCUMENT_IMAGE_SIZE, DEFAULT_MAX_IMAGE_SIZE } from '../../constants'
import { formatFileSize } from '../../helpers'

interface TextEditorProps {
  applicationConfig?: AppConfigObject
  currentContent: string
  markdownMode?: boolean
  isReadOnly?: boolean
  resource?: Resource
  autoFocus?: boolean
}
interface TextEditorEmits {
  (e: 'update:currentContent', value: string): void
}
const {
  markdownMode = false,
  isReadOnly = false,
  applicationConfig,
  currentContent,
  resource,
  autoFocus = true
} = defineProps<TextEditorProps>()

defineEmits<TextEditorEmits>()

const { current: currentLanguage, $gettext, $ngettext } = useGettext()
const { currentTheme } = useThemeStore()
const { showErrorMessage } = useMessages()

// Should not be a ref, otherwise functions like setMarkdown won't work
const editorConfig = computed(() => {
  const {
    showPreviewOnlyMd = true,
    maxImageSize = DEFAULT_MAX_IMAGE_SIZE,
    maxDocumentImageSize = DEFAULT_MAX_DOCUMENT_IMAGE_SIZE
  }: AppConfigObject = applicationConfig
  return { showPreviewOnlyMd, maxImageSize, maxDocumentImageSize }
})

const isMarkdown = computed(() => {
  return (
    markdownMode ||
    ['md', 'markdown'].includes(resource?.extension) ||
    !unref(editorConfig).showPreviewOnlyMd
  )
})

const theme = computed(() => (unref(currentTheme).isDark ? 'dark' : 'light'))

const sanitize = (html) =>
  dompurify.sanitize(html, { ADD_ATTR: ['target'], ADD_TAGS: ['foreignObject'] })

const dataUriRegex = /data:image\/[a-z0-9.+-]+;base64,[A-Za-z0-9+/=]+/gi

const isAddingImages = ref(false)
const pendingImageSize = ref(0)

const formatSize = (size: number) => formatFileSize(size, currentLanguage)

const MAX_LISTED_FILE_NAMES = 3

const formatFileNames = (names: Array<string>): string => {
  const listed = names.slice(0, MAX_LISTED_FILE_NAMES).join(', ')
  if (names.length <= MAX_LISTED_FILE_NAMES) {
    return listed
  }
  return $gettext('%{ files } and %{ count } more', {
    files: listed,
    count: (names.length - MAX_LISTED_FILE_NAMES).toString()
  })
}

const usedDocumentImageSize = (markdown: string): number => {
  const dataUris: string[] = markdown?.match(dataUriRegex) ?? []
  return dataUris.reduce((total, dataUri) => total + dataUri.length, 0)
}

const readAsDataUrl = (file: File): Promise<string> =>
  new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result as string)
    reader.onerror = () => reject(reader.error)
    reader.readAsDataURL(file)
  })

type UploadImgCallBack = (urls: Array<{ url: string; alt: string; title: string }>) => void

const onUploadImg = async (files: Array<File>, callBack: UploadImgCallBack) => {
  if (!unref(isMarkdown)) {
    callBack([])
    return
  }

  const { maxImageSize, maxDocumentImageSize } = unref(editorConfig)
  isAddingImages.value = true
  try {
    callBack(await pickImages(files, maxImageSize, maxDocumentImageSize))
  } finally {
    isAddingImages.value = false
  }
}

const pickImages = async (
  files: Array<File>,
  maxImageSize: number,
  maxDocumentImageSize: number
): Promise<Array<{ url: string; alt: string; title: string }>> => {
  let usedSize = usedDocumentImageSize(currentContent) + unref(pendingImageSize)
  const accepted: Array<{ url: string; alt: string; title: string }> = []
  const tooBig: Array<string> = []
  const unreadable: Array<string> = []
  const doesNotFit: Array<string> = []

  for (const file of files) {
    if (!file.type.startsWith('image/')) {
      continue
    }

    if (file.size > maxImageSize) {
      tooBig.push(file.name)
      continue
    }

    let dataUri: string
    try {
      dataUri = await readAsDataUrl(file)
    } catch {
      unreadable.push(file.name)
      continue
    }

    if (usedSize + dataUri.length > maxDocumentImageSize) {
      doesNotFit.push(file.name)
      continue
    }

    usedSize += dataUri.length
    pendingImageSize.value += dataUri.length
    accepted.push({ url: dataUri, alt: file.name, title: '' })
  }

  if (tooBig.length) {
    showErrorMessage({
      title: $ngettext('Image is too big', 'Images are too big', tooBig.length),
      desc: $gettext('%{ files }. Max image size: %{ limit }.', {
        files: formatFileNames(tooBig),
        limit: formatSize(maxImageSize)
      })
    })
  }

  if (doesNotFit.length) {
    showErrorMessage({
      title: $ngettext(
        'Image does not fit in this document',
        'Images do not fit in this document',
        doesNotFit.length
      ),
      desc: $ngettext(
        '%{ files }. Only %{ remaining } left - remove an image or link the file instead.',
        '%{ files }. Only %{ remaining } left - remove an image or link the files instead.',
        doesNotFit.length,
        {
          files: formatFileNames(doesNotFit),
          remaining: formatSize(Math.max(maxDocumentImageSize - usedSize, 0))
        }
      )
    })
  }

  if (unreadable.length) {
    showErrorMessage({
      title: $ngettext('Image could not be read', 'Images could not be read', unreadable.length),
      desc: formatFileNames(unreadable)
    })
  }

  return accepted
}

watch(
  () => currentContent,
  () => {
    pendingImageSize.value = 0
  }
)

// CodeMirror (used internally by md-editor-v3) binds Tab to indent instead of moving
// focus, which is correct for editing but leaves keyboard users with no way to reach
// the toolbar without leaving the editor entirely. CodeMirror's own escape hatch for
// this (Ctrl-M, or Shift-Alt-M on Mac) can't actually be triggered on a standard Mac
// keyboard layout, since Option+Shift+M is consumed by macOS as a dead-key/diacritic
// combo before it reaches the page. Ctrl-M (no Alt/Option involved) works identically
// on every platform, so it's handled here instead, ahead of CodeMirror's own listener.
const onFocusOutShortcut = (event: KeyboardEvent) => {
  if (event.key.toLowerCase() !== 'm' || !event.ctrlKey || event.altKey || event.metaKey) {
    return
  }
  const target = document.querySelector<HTMLElement>(
    '#text-editor-component .md-editor-toolbar-item:not([disabled])'
  )
  event.preventDefault()
  event.stopPropagation()
  if (target) {
    target.focus()
    return
  }
  document.querySelector<HTMLElement>('#text-editor-component .cm-content')?.blur()
}

onMounted(async () => {
  if (isReadOnly) {
    return
  }
  // md-editor-v3's underlying CodeMirror instance renders its editable div without
  // an accessible name and doesn't expose a prop for one, so it's set imperatively here
  await nextTick()
  document
    .querySelector('#text-editor-component .cm-content')
    ?.setAttribute('aria-label', $gettext('Text editor'))
  document
    .querySelector('#text-editor-component .cm-content')
    ?.setAttribute('aria-describedby', 'text-editor-focus-out-hint')
  document
    .getElementById('text-editor-container')
    ?.addEventListener('keydown', onFocusOutShortcut, { capture: true })
})

onBeforeUnmount(() => {
  document
    .getElementById('text-editor-container')
    ?.removeEventListener('keydown', onFocusOutShortcut, { capture: true })
})

config({
  editorConfig: {
    languageUserDefined
  },
  editorExtensions: {
    prettier: {
      prettierInstance: prettier,
      parserMarkdownInstance: parserMarkdown
    },
    highlight: {
      instance: highlight
    },
    screenfull: {
      instance: screenfull
    },
    katex: {
      instance: katex
    },
    mermaid: {
      instance: mermaid
    },
    // The crop entry is hidden (see the CSS below), so Cropper is never constructed.
    // md-editor-v3 CDN-injects cropperjs from unpkg unless an instance is registered
    // (composition.ts: noCropperScript) — this stub keeps that request from being made.
    cropper: {
      instance: class {} as never
    }
  },
  markdownItConfig(md) {
    md.renderer.rules.link_open = function (tokens, idx, options, _, self) {
      const token = tokens[idx]
      const href = token.attrGet('href')

      if (!href) {
        return self.renderToken(tokens, idx, options)
      }

      token.attrSet('target', '_blank')
      token.attrSet('rel', 'noopener noreferrer')

      return self.renderToken(tokens, idx, options)
    }
  }
})
</script>
<style lang="scss">
#text-editor-component {
  height: 100%;

  .md-editor-mermaid {
    .messageText,
    .legend text,
    .titleText,
    .sectionTitle.sectionTitle0,
    .grid .tick text,
    text {
      fill: var(--oc-color-text-default);
      opacity: 0.8;
    }

    line {
      stroke: var(--oc-color-text-default);
      opacity: 0.8;
    }

    .slice {
      fill: #000;
    }

    .sectionTitle.sectionTitle1,
    .taskText.taskText1,
    .taskText.taskText0 {
      fill: #fff;
    }

    .messageLine1,
    .messageLine0,
    .flowchart-link,
    .transition,
    .relationshipLine {
      stroke: var(--oc-color-text-default);
      opacity: 0.8;
    }

    .nodeLabel p {
      fill: #000;
      color: #000;
    }
  }

  .footer-slot {
    display: inline-flex;
    align-items: center;
    gap: 0.625rem;
    vertical-align: middle;
    padding-inline-start: 10px;
  }

  .footer-links {
    display: inline-flex;
    gap: 0.625rem;
  }

  #text-editor-component-html-wrapper {
    margin-left: var(--oc-space-xsmall);
  }
  .md-editor-code-head {
    z-index: 0;
  }

  // Hides the "Crop And Upload" entry: cropping is dropped until users ask for it.
  // md-editor-v3 hardcodes the three image-dropdown entries with identical classes and
  // exposes no prop to hide one (still true in v7), so position is the only handle.
  // Requiring :nth-child(3) and :last-child together means an upstream reorder stops the
  // selector matching, showing all entries rather than hiding the wrong one.
  .md-editor-menu-item-image:nth-child(3):last-child {
    display: none;
  }
}

.toastui-editor-tabs {
  // Fix tab with for long i18n text
  .tab-item {
    width: auto;
    padding-left: var(--oc-space-small);
    padding-right: var(--oc-space-small);
  }
}

#space-description-preview {
  background-color: transparent;

  .md-editor-preview-wrapper {
    padding: 0;
  }

  .md-editor-preview {
    color: var(--oc-color-text-default);
    font-size: var(--oc-text-default);
    // override md-editor-v3's default link/inline-code colors, which don't meet
    // WCAG AA contrast against the app's card background
    --md-theme-link-color: var(--oc-color-swatch-primary-muted);
    --md-theme-code-inline-color: var(--oc-color-text-default);
    --md-theme-code-inline-bg-color: var(--oc-color-background-default);
  }
}
</style>
