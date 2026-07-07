<template>
  <div class="oc-error-log">
    <oc-textarea
      class="oc-error-log-textarea oc-mt-s oc-text-small"
      :label="contentLabel"
      :model-value="content"
      rows="4"
      readonly
    />
    <div class="oc-flex oc-flex-between oc-mt-s">
      <div class="oc-flex">
        <div v-if="showCopied" class="oc-flex oc-flex-middle">
          <oc-icon variation="success" name="checkbox-circle" />
          <p class="oc-error-log-content-copied oc-ml-s oc-my-rm" v-text="$gettext('Copied')" />
        </div>
      </div>
      <oc-button
        size="small"
        variation="primary"
        appearance="filled"
        @click="copyContentToClipboard"
      >
        {{ $gettext('Copy') }}
      </oc-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import { useGettext } from 'vue3-gettext'

/**
 * OcErrorLog Component
 *
 * This component displays an error log with a textarea containing the error content
 * and a button to copy the content to the clipboard. It also provides feedback
 * when the content is successfully copied.
 *
 * @name OcErrorLog
 * @status ready
 * @release 2.0.0
 *
 * @props {string} content - The error log content to be displayed in the textarea.
 *
 * @example
 * <OcErrorLog content="Error details" />
 *
 */

interface Props {
  content: string
}

defineOptions({
  name: 'OcErrorLog',
  status: 'ready',
  release: '2.0.0'
})
const { content } = defineProps<Props>()
const { $gettext } = useGettext()
const showCopied = ref(false)

const contentLabel = computed(() => {
  return $gettext(
    'Copy the following information and pass them to technical support to troubleshoot the problem:'
  )
})

const copyContentToClipboard = () => {
  navigator.clipboard.writeText(content)
  showCopied.value = true
  setTimeout(() => (showCopied.value = false), 500)
}
</script>

<style lang="scss">
.oc-error-log {
  &-textarea {
    resize: none;

    label {
      color: var(--oc-color-text-muted);
    }
  }

  &-content-copied {
    color: var(--oc-color-swatch-success-default);
  }
}
</style>
