<template>
  <oc-button
    v-if="isClipboardCopySupported"
    v-oc-tooltip="$gettext('Copy link to clipboard')"
    :aria-label="$gettext('Copy link to clipboard')"
    appearance="raw"
    class="oc-files-public-link-copy-url"
    @click="copyLinkToClipboard"
  >
    <oc-icon :name="copied ? 'checkbox-circle' : 'file-copy'" fill-type="line" />
  </oc-button>
</template>

<script lang="ts" setup>
import { useMessages } from '@ownclouders/web-pkg'
import { useClipboard } from '@vueuse/core'
import { useGettext } from 'vue3-gettext'
import { LinkShare } from '@ownclouders/web-client'

interface Props {
  linkShare: LinkShare
}
const { linkShare } = defineProps<Props>()
const { $gettext } = useGettext()
const { showMessage } = useMessages()

const {
  copy,
  copied,
  isSupported: isClipboardCopySupported
} = useClipboard({ legacy: true, copiedDuring: 550 })

const copyLinkToClipboard = () => {
  copy(linkShare.webUrl)
  showMessage({
    title: linkShare.isQuickLink
      ? $gettext('The link has been copied to your clipboard.')
      : $gettext('The link "%{linkName}" has been copied to your clipboard.', {
          linkName: linkShare.displayName
        })
  })
}
</script>
