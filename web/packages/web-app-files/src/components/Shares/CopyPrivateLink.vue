<template>
  <div class="oc-flex oc-flex-middle copy-private-link">
    <oc-button v-oc-tooltip="tooltip" gap-size="none" appearance="raw" @click="copyLinkToClipboard">
      <oc-icon size="small" :name="copied ? 'checkbox-circle' : 'file-copy'" fill-type="line" />
      <span class="oc-ml-xs" v-text="$gettext('Permanent link')"
    /></oc-button>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { Resource } from '@ownclouders/web-client'
import { useGettext } from 'vue3-gettext'
import { useClipboard } from '@vueuse/core'
import { useMessages } from '@ownclouders/web-pkg'

interface Props {
  resource: Resource
}
const props = defineProps<Props>()

const { $gettext } = useGettext()
const { showMessage } = useMessages()
const { copy, copied } = useClipboard({ legacy: true, copiedDuring: 550 })

const copyLinkToClipboard = () => {
  copy(props.resource.privateLink)
  showMessage({
    title: $gettext('Permanent link copied'),
    desc: $gettext('The permanent link has been copied to your clipboard.')
  })
}

const tooltip = computed(() => {
  return $gettext(
    'Copy the link to point your team to this item. Works only for people with existing access.'
  )
})
</script>
