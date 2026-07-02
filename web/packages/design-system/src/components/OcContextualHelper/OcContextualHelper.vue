<template>
  <div class="oc-contextual-helper">
    <oc-button :id="buttonId" :aria-label="$gettext('Show more information')" appearance="raw">
      <oc-icon name="question" fill-type="line" size="small" :accessible-label="title" />
    </oc-button>
    <oc-info-drop :drop-id="dropId" :toggle="toggleId" v-bind="props" />
  </div>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import { uniqueId } from '../../helpers'
import OcButton from '../OcButton/OcButton.vue'
import OcIcon from '../OcIcon/OcIcon.vue'
import OcInfoDrop from '../OcInfoDrop/OcInfoDrop.vue'
import { ContextualHelperDataListItem } from '../../helpers'

/**
 * OcContextualHelper - A component that displays contextual help information in a dropdown when clicked.
 * It shows a question mark icon that opens additional information in a popup.
 *
 * @prop {string} title - The title text displayed in the header of the contextual helper popup.
 * @prop {string} [text=''] - Optional main descriptive text content.
 * @prop {ContextualHelperDataListItem[]} [list=[]] - Optional array of list items to display in a definition list format.
 * @prop {string} [endText=''] - Optional text displayed after the list and before the "Read more" link.
 * @prop {string} [readMoreLink=''] - Optional URL for the "Read more" link that opens in a new tab.
 *
 * @example
 * ```vue
 * <!-- Basic usage -->
 * <oc-contextual-helper
 *   title="title?"
 *   text="explanation text"
 * />
 *
 * <!-- With list and read more link -->
 * <oc-contextual-helper
 *   title="title?"
 *   text="explanation text"
 *   :list="[
 *     { text: 'Personal storage: 5GB', headline: true },
 *     { text: 'Project storage: 50GB', headline: true }
 *   ]"
 *   end-text="Contact your administrator for more space."
 *   read-more-link="https://docs.example.com/storage-quotas"
 * />
 * ```
 */

interface Props {
  title: string
  text?: string
  list?: ContextualHelperDataListItem[]
  endText?: string
  readMoreLink?: string
}
defineOptions({
  name: 'OcContextualHelper',
  status: 'unreleased'
})

const { title, text = '', list = [], endText = '', readMoreLink = '' } = defineProps<Props>()
const props = computed(() => ({ title, text, list, endText, readMoreLink }))

const dropId = computed(() => uniqueId('oc-contextual-helper-'))
const buttonId = computed(() => `${unref(dropId)}-button`)
const toggleId = computed(() => `#${unref(buttonId)}`)
</script>

<style lang="scss">
.oc-contextual-helper {
  display: inline-block;
  .oc-button {
    vertical-align: middle;
  }
}
</style>
