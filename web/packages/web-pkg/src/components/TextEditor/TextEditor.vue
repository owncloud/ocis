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
      no-upload-img
      @on-change="(value) => $emit('update:currentContent', value)"
    >
      <template #defFooters>
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
      </template>
    </md-editor>
  </div>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
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
import { useThemeStore } from '../../composables'

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

const { current: currentLanguage } = useGettext()
const { currentTheme } = useThemeStore()

// Should not be a ref, otherwise functions like setMarkdown won't work
const editorConfig = computed(() => {
  const { showPreviewOnlyMd = true }: AppConfigObject = applicationConfig
  return { showPreviewOnlyMd }
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
  }
}
</style>
