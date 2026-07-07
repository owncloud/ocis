<template>
  <span>
    <span
      v-oc-tooltip="tooltip"
      class="oc-avatars"
      :class="{ 'oc-avatars-stacked': stacked }"
      aria-hidden="true"
    >
      <template v-if="avatars.length > 0">
        <oc-avatar
          v-for="avatar in avatars"
          :key="avatar.username"
          :src="avatar.avatar"
          :user-name="avatar.displayName"
          :width="30"
        />
      </template>
      <template v-if="otherItems.length > 0">
        <component
          :is="getAvatarComponentForItem(item)"
          v-for="(item, index) in otherItems"
          :key="item.name + index"
          :name="item.name"
        />
      </template>
      <oc-avatar-count v-if="isOverlapping" :count="items.length - maxDisplayed" />
    </span>
    <span v-if="accessibleDescription" class="oc-invisible-sr" v-text="accessibleDescription" />
  </span>
</template>

<script lang="ts" setup>
import { shareType } from '../../utils/shareType'
import OcAvatar from '../OcAvatar/OcAvatar.vue'
import OcAvatarCount from '../OcAvatarCount/OcAvatarCount.vue'
import OcAvatarLink from '../OcAvatarLink/OcAvatarLink.vue'
import OcAvatarGroup from '../OcAvatarGroup/OcAvatarGroup.vue'
import OcAvatarFederated from '../OcAvatarFederated/OcAvatarFederated.vue'
import OcAvatarGuest from '../OcAvatarGuest/OcAvatarGuest.vue'
import { computed, unref } from 'vue'

/**
 * OcAvatars - A component for displaying a group of different types of avatars with various display options.
 *
 * @prop {Array<Object>} items - Users, public links, groups, federated and guests to be displayed with avatar.
 * @prop {boolean} [stacked=false] - Asserts whether avatars should be stacked on each other.
 * @prop {boolean} [isTooltipDisplayed=false] - Asserts whether tooltip should be displayed on hover/focus.
 * @prop {number} [maxDisplayed=null] - Limits the number of avatars which will be displayed.
 * @prop {string} [accessibleDescription=null] - A description of the avatar group for screen readers. This is required as the avatar group element
 *   is hidden for screen readers.
 *
 * @example
 * ```vue
 * <!-- Basic usage with users and groups -->
 * <oc-avatars
 *   :items="[
 *     { username: 'john', displayName: 'doe', shareType: 0, avatar: 'avatar-url' },
 *     { name: 'Developers', shareType: 1 }
 *   ]"
 *   accessible-description="Marie and Developers group have access to this resource"
 * />
 *
 * <!-- Stacked avatars with tooltip and limited display -->
 * <oc-avatars
 *   :items="usersList"
 *   :stacked="true"
 *   :is-tooltip-displayed="true"
 *   :max-displayed="5"
 *   accessible-description="This resource is shared with multiple users and groups"
 * />
 * ```
 */

type Item = {
  displayName?: string
  name?: string
  shareType?: number
  username?: string
  avatar?: string
}
interface Props {
  items: Item[]
  stacked?: boolean
  isTooltipDisplayed?: boolean
  maxDisplayed?: number
  accessibleDescription?: string | null
}

defineOptions({
  name: 'OcAvatars',
  status: 'ready',
  release: '2.1.0'
})
const {
  items,
  stacked = false,
  isTooltipDisplayed = false,
  maxDisplayed = null,
  accessibleDescription = null
} = defineProps<Props>()

function getAvatarComponentForItem(item: Item) {
  switch (item.shareType) {
    case shareType.link:
      return OcAvatarLink
    case shareType.remote:
      return OcAvatarFederated
    case shareType.group:
      return OcAvatarGroup
    case shareType.guest:
      return OcAvatarGuest
  }
}

const isOverlapping = computed(() => {
  return maxDisplayed && maxDisplayed < items.length
})

const tooltip = computed(() => {
  if (isTooltipDisplayed) {
    const names = unref(avatars).map((user) => user.displayName)

    if (unref(otherItems).length > 0) {
      names.push(...unref(otherItems).map((item) => item.name))
    }

    let tooltip = names.join(', ')

    if (unref(isOverlapping)) {
      tooltip += ` +${items.length - maxDisplayed}`
    }

    return tooltip
  }

  return null
})

const avatars = computed(() => {
  const a = items.filter((u) => u.shareType === shareType.user)
  if (!unref(isOverlapping)) {
    return a
  }
  return a.slice(0, maxDisplayed)
})

const otherItems = computed(() => {
  const a = items.filter((u) => u.shareType !== shareType.user)
  if (!unref(isOverlapping)) {
    return a
  }
  if (maxDisplayed <= unref(avatars).length) {
    return []
  }
  return a.slice(0, maxDisplayed - unref(avatars).length)
})
</script>

<style lang="scss">
.oc-avatars {
  display: inline-flex;
  box-sizing: border-box;
  flex-flow: row nowrap;
  gap: var(--oc-space-xsmall);
  width: fit-content;

  &-stacked {
    .oc-avatar + .oc-avatar,
    .oc-avatar-count,
    .oc-avatar + .oc-avatar-item,
    .oc-avatar-item + .oc-avatar-item {
      border: 1px solid var(--oc-color-text-inverse);
      margin-left: -25px;
    }
  }
}
</style>

<docs>
```js
<template>
  <div>
    <h3>Default configuration</h3>
    <p>No stacking, no tooltip, no <b>:maxDisplayed</b> configured</p>
    <oc-avatars :items="items" accessible-description="This resource is shared with many users." class="oc-mb" />
    <h3>Stacked, tooltip, maxDisplayed</h3>
    <p>Using <b>:stacked="true"</b>, <b>:isTooltipDisplayed="true"</b> and <b>:maxDisplayed="5"</b></p>
    <oc-avatars :items="items" accessible-description="This resource is shared with many users." :stacked="true" :maxDisplayed="5" :isTooltipDisplayed="true" />
    <h3>Unstacked, tooltip, maxDisplayed</h3>
    <p>Using <b>:isTooltipDisplayed="true"</b> and <b>:maxDisplayed="2"</b></p>
    <oc-avatars :items="items" accessible-description="This resource is shared with many users." :maxDisplayed="2" :isTooltipDisplayed="true" />
  </div>
</template>
<script>
import { shareType } from "../../utils/shareType"
export default {
  data: () => ({
    items: [
      {
        name: "bob",
        shareType: shareType.remote
      },
      {
        username: "marie",
        displayName: "Marie",
        avatar: "https://images.unsplash.com/photo-1584308972272-9e4e7685e80f?ixid=MXwxMjA3fDB8MHxzZWFyY2h8Mzh8fGZhY2V8ZW58MHwyfDB8&ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=60",
        shareType: shareType.user
      },
      {
        username: "peter",
        displayName: "Peter",
        avatar: "https://images.unsplash.com/photo-1610216705422-caa3fcb6d158?ixid=MXwxMjA3fDB8MHxzZWFyY2h8MTB8fGZhY2V8ZW58MHwyfDB8&ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=60",
        shareType: shareType.user
      },
      {
        username: "udo",
        displayName: "Udo",
        avatar: "https://images.unsplash.com/photo-1584308972272-9e4e7685e80f?ixid=MXwxMjA3fDB8MHxzZWFyY2h8Mzh8fGZhY2V8ZW58MHwyfDB8&ixlib=rb-1.2.1&auto=format&fit=crop&w=800&q=60",
        shareType: shareType.user
      },
      {
        name: "john",
        shareType: shareType.guest
      },
      {
        name: "Public link",
        shareType: shareType.link
      },
      {
        name: "Test",
        shareType: shareType.group
      }
    ]
  })
}
</script>
```
</docs>
