<template>
  <oc-drop
    ref="drop"
    class="oc-width-1-1 oc-info-drop"
    :drop-id="dropId"
    :toggle="toggle"
    :mode="mode"
    close-on-click
    @hide-drop="() => (dropOpen = false)"
    @show-drop="() => (dropOpen = true)"
  >
    <focus-trap :active="dropOpen">
      <div class="info-drop-content">
        <div class="oc-flex oc-flex-between info-header oc-border-b oc-pb-s">
          <h4 class="oc-m-rm info-title" v-text="$gettext(title)" />
          <oc-button
            v-oc-tooltip="$gettext('Close')"
            appearance="raw"
            :aria-label="$gettext('Close')"
          >
            <oc-icon
              name="close"
              fill-type="line"
              size="medium"
              variation="inherit"
              :accessible-label="$gettext('Close')"
            />
          </oc-button>
        </div>
        <p v-if="text" class="info-text" v-text="$gettext(text)" />
        <dl v-if="listItems.length" class="info-list">
          <component
            :is="item.headline ? 'dt' : 'dd'"
            v-for="(item, index) in listItems"
            :key="index"
          >
            {{ $gettext(item.text) }}
          </component>
        </dl>
        <p v-if="endText" class="info-text-end" v-text="$gettext(endText)" />
        <oc-button
          v-if="readMoreLink"
          type="a"
          appearance="raw"
          size="small"
          class="info-more-link"
          :href="readMoreLink"
          target="_blank"
        >
          {{ $gettext('Read more') }}
        </oc-button>
      </div>
    </focus-trap>
  </oc-drop>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import OcButton from '../OcButton/OcButton.vue'
import OcIcon from '../OcIcon/OcIcon.vue'
import OcDrop from '../OcDrop/OcDrop.vue'
import { uniqueId } from '../../helpers'
import { FocusTrap } from 'focus-trap-vue'
import { ContextualHelperDataListItem } from '../../helpers'

/**
 * OcInfoDrop - A dropdown component that displays contextual help information
 *
 * @prop {string} [dropId] - Optional unique identifier for the dropdown. If not provided, an auto-generated ID will be used.
 * @prop {string} [toggle=''] - CSS selector for the element that triggers the dropdown.
 * @prop {'click' | 'hover' | 'manual'} [mode='click'] - Interaction mode to trigger the dropdown.
 * @prop {string} title - Required title text displayed in the header of the dropdown.
 * @prop {string} [text] - Optional main descriptive text content.
 * @prop {ContextualHelperDataListItem[]} [list] - Optional array of list items to display in a definition list format.
 * @prop {string} [endText] - Optional text displayed after the list and before the "Read more" link.
 * @prop {string} [readMoreLink] - Optional URL for the "Read more" link that opens in a new tab.
 *
 * @event {void} hide-drop - Emitted when the dropdown is hidden.
 * @event {void} show-drop - Emitted when the dropdown is shown.
 *
 * @example
 * ```vue
 * <template>
 *   <oc-info-drop
 *     title="title"
 *     text="text"
 *     :list="[
 *       {text: 'text', headline: true},
 *       {text: 'text', headline: true}
 *     ]"
 *     end-text="endText"
 *     read-more-link="https://example.com"
 *   />
 * </template>
 * ```
 */

interface Props {
  dropId?: string
  toggle?: string
  mode?: 'click' | 'hover' | 'manual'
  title: string
  text?: string
  list?: ContextualHelperDataListItem[]
  endText?: string
  readMoreLink?: string
}

defineOptions({
  name: 'OcInfoDrop',
  status: 'unreleased'
})

const {
  dropId = uniqueId('oc-info-drop-'),
  toggle = '',
  mode = 'click',
  title,
  text = '',
  list = [],
  endText = '',
  readMoreLink = ''
} = defineProps<Props>()

const dropOpen = ref(false)
const listItems = computed(() => {
  return (list || []).filter((item) => !!item.text)
})
</script>

<style lang="scss">
.oc-info-drop {
  display: inline-block;
  .oc-button {
    vertical-align: middle;
  }
  .info-drop-content {
    font-size: var(--oc-font-size-medium);
    color: var(--oc-color-text-default);
  }
  .info-more-link {
    font-size: var(--oc-font-size-medium) !important;
  }
  .info-header {
    align-items: center;
  }
  .info-title {
    font-size: 1.125rem;
    font-weight: normal;
  }
  .info-list:first-child,
  .info-text:first-child {
    margin-top: 0;
  }
  .info-list:last-child,
  .info-text:last-child {
    margin-bottom: 0;
  }
  .info-list {
    font-weight: bold;
    margin-bottom: var(--oc-space-xsmall);
    margin-top: var(--oc-space-small);
    dt {
      &:first-child {
        margin-top: 0;
      }
    }
    dd {
      margin-left: 0;
      font-weight: normal;
    }
  }
}
</style>

<docs>
## Examples
A simple example, using only text.
```js
<template>
  <div>
    <oc-info-drop v-bind="helperContent"/>
  </div>
</template>
<script>
export default {
  computed: {
    helperContent() {
      return {
        text: "Invite persons or groups to access this file or folder.",
      }
    }
  },
}
</script>
```

An example using Title, Text, List, End-Text and Read-More-Link properties.
```js
<template>
  <div>
    <oc-info-drop v-bind="helperContent"/>
  </div>
</template>
<script>
export default {
  computed: {
    helperContent() {
      return {
        title: 'Choose how access is granted ',
        text: "Share a file or folder by link",
        list: [
          {text: "Only invited people can access", headline: true},
          {text: "Only people from the list \"Invited people\" can access. If there is no list, no people are invited yet."},
          {text: "Everyone with the link", headline: true },
          {text: "Everyone with the link can access. Note: If you share this link with people from the list \"Invited people\", they need to login-in so that their individual assigned permissions can take effect. If they are not logged-in, the permissions of the link take effect." }
        ],
        endText: "Invited persons can not see who else has access",
        readMoreLink: "https://owncloud.design"
      }
    }
  },
}
</script>
```
</docs>
