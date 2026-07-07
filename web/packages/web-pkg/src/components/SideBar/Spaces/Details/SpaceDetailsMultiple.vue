<template>
  <div id="oc-spaces-details-multiple-sidebar">
    <div class="spaces-preview oc-mb">
      <div class="spaces-preview-body">
        <oc-icon class="preview-icon" size="xxlarge" variation="passive" name="layout-grid" />
        <p class="preview-text" v-text="selectedSpacesString" />
      </div>
    </div>
    <oc-definition-list :aria-label="detailsTableLabel" :items="items" />
  </div>
</template>
<script lang="ts" setup>
import { formatFileSize } from '../../../../helpers'
import { computed, unref } from 'vue'
import { SpaceResource } from '@ownclouders/web-client'
import { useGettext } from 'vue3-gettext'

interface Props {
  selectedSpaces: SpaceResource[]
}
const { selectedSpaces } = defineProps<Props>()

const language = useGettext()
const { $gettext, $ngettext } = language
const totalSelectedSpaceQuotaTotal = computed(() => {
  let total = 0
  selectedSpaces.forEach((space) => {
    total += space.spaceQuota.total
  })
  return formatFileSize(total, language.current)
})
const totalSelectedSpaceQuotaRemaining = computed(() => {
  let remaining = 0
  selectedSpaces.forEach((space) => {
    if (space.disabled) {
      return
    }
    remaining += space.spaceQuota.remaining
  })
  return formatFileSize(remaining, language.current)
})
const totalSelectedSpaceQuotaUsed = computed(() => {
  let used = 0
  selectedSpaces.forEach((space) => {
    if (space.disabled) {
      return
    }
    used += space.spaceQuota.used
  })
  return formatFileSize(used, language.current)
})
const totalEnabledSpaces = computed(() => {
  return selectedSpaces.filter((s) => !s.disabled).length
})
const totalDisabledSpaces = computed(() => {
  return selectedSpaces.filter((s) => s.disabled).length
})
const detailsTableLabel = computed(() => {
  return $gettext('Overview of the information about the selected spaces')
})
const selectedSpacesString = computed(() => {
  return $ngettext(
    '%{ itemCount } space selected',
    '%{ itemCount } spaces selected',
    selectedSpaces.length,
    {
      itemCount: selectedSpaces.length.toString()
    }
  )
})

const items = [
  { term: $gettext('Total quota:'), definition: unref(totalSelectedSpaceQuotaTotal) },
  { term: $gettext('Remaining quota:'), definition: unref(totalSelectedSpaceQuotaRemaining) },
  { term: $gettext('Used quota:'), definition: unref(totalSelectedSpaceQuotaUsed) },
  { term: $gettext('Enabled:'), definition: unref(totalEnabledSpaces) },
  { term: $gettext('Disabled:'), definition: unref(totalDisabledSpaces) }
]
</script>
<style lang="scss">
.spaces-preview {
  position: relative;
  background-color: var(--oc-color-background-muted);
  border: 10px solid var(--oc-color-background-muted);
  height: 230px;
  text-align: center;
  color: var(--oc-color-swatch-passive-muted);

  &-body {
    margin: 0;
    position: absolute;
    top: 50%;
    left: 50%;
    -ms-transform: translate(-50%, -50%);
    transform: translate(-50%, -50%);

    .preview-icon {
      display: inline-block;
    }
    .preview-text {
      display: block;
    }
  }
}
</style>
