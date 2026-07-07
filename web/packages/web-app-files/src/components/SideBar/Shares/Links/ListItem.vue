<template>
  <div class="oc-width-1-1 oc-flex oc-flex-middle oc-flex-between files-links-details">
    <div class="oc-flex oc-flex-middle files-links-content">
      <oc-avatar-item :width="36" icon-size="medium" icon="link" name="df" />
      <div class="files-links-name-wrapper oc-pl-s">
        <div class="oc-flex oc-flex-middle">
          <div class="oc-text-truncate">
            <span aria-hidden="true" class="files-links-name" v-text="linkShare.displayName" />
          </div>
        </div>
        <div class="oc-flex oc-flex-nowrap oc-flex-middle">
          <link-role-dropdown
            v-if="isModifiable"
            :model-value="currentLinkType"
            :available-link-type-options="availableLinkTypeOptions"
            drop-offset="0"
            @update:model-value="updateSelectedType"
          />
          <span
            v-else
            v-oc-tooltip="$gettext(currentLinkRoleDescription)"
            class="link-current-role"
            v-text="$gettext(currentLinkRoleLabel)"
          />
        </div>
      </div>
    </div>
    <div class="oc-flex oc-flex-middle">
      <div class="oc-flex">
        <oc-icon
          v-if="linkShare.hasPassword"
          v-oc-tooltip="$gettext('This link is password-protected')"
          name="lock-password"
          class="oc-files-file-link-has-password oc-mr-xs"
          fill-type="line"
          :accessible-label="$gettext('This link is password-protected')"
        />
      </div>
      <expiration-date-indicator
        v-if="linkShare.expirationDateTime"
        :expiration-date="DateTime.fromISO(linkShare.expirationDateTime)"
        class="oc-mx-xs"
      />
      <copy-link :link-share="linkShare" class="oc-mx-xs" />
      <edit-dropdown
        :can-rename="canRename"
        :is-modifiable="isModifiable"
        :is-password-removable="isPasswordRemovable"
        :link-share="linkShare"
        class="oc-ml-xs"
        @remove-public-link="$emit('removePublicLink', $event)"
        @update-link="$emit('updateLink', $event)"
        @show-password-modal="showPasswordModal"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { DateTime } from 'luxon'
import { LinkRoleDropdown, useAbility, useLinkTypes, useModals } from '@ownclouders/web-pkg'
import { LinkShare } from '@ownclouders/web-client'
import { computed, inject, Ref, ref, unref } from 'vue'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { useGettext } from 'vue3-gettext'
import SetLinkPasswordModal from '../../../Modals/SetLinkPasswordModal.vue'
import { SharingLinkType } from '@ownclouders/web-client/graph/generated'
import ExpirationDateIndicator from '../ExpirationDateIndicator.vue'
import CopyLink from './CopyLink.vue'
import EditDropdown from './EditDropdown.vue'

interface Props {
  canRename?: boolean
  isFolderShare?: boolean
  isModifiable?: boolean
  isPasswordEnforced?: boolean
  isPasswordRemovable?: boolean
  linkShare: LinkShare
}
interface Emits {
  (e: 'removePublicLink', linkShare: LinkShare): void
  (e: 'updateLink', payload: { linkShare: LinkShare; options?: { type?: SharingLinkType } }): void
  (e: 'showPasswordModal'): void
}
const {
  canRename = false,
  isFolderShare = false,
  isModifiable = false,
  isPasswordEnforced = false,
  isPasswordRemovable = false,
  linkShare
} = defineProps<Props>()
const emit = defineEmits<Emits>()
const { dispatchModal } = useModals()
const { $gettext } = useGettext()
const { can } = useAbility()
const { getAvailableLinkTypes, getLinkRoleByType } = useLinkTypes()

const space = inject<Ref<SpaceResource>>('space')
const resource = inject<Ref<Resource>>('resource')

const currentLinkType = ref<SharingLinkType>(linkShare.type)

const canDeleteReadOnlyPublicLinkPassword = computed(() =>
  can('delete-all', 'ReadOnlyPublicLinkPassword')
)

const updateSelectedType = (type: SharingLinkType) => {
  currentLinkType.value = type

  const needsNoPw =
    type === SharingLinkType.Internal ||
    (unref(canDeleteReadOnlyPublicLinkPassword) && type === SharingLinkType.View)

  if (!linkShare.hasPassword && !needsNoPw && isPasswordEnforced) {
    showPasswordModal(() => emit('updateLink', { linkShare: { ...linkShare }, options: { type } }))
    return
  }

  emit('updateLink', { linkShare, options: { type } })
}

const showPasswordModal = (callbackFn: () => void = undefined) => {
  dispatchModal({
    title: linkShare.hasPassword ? $gettext('Edit password') : $gettext('Add password'),
    customComponent: SetLinkPasswordModal,
    customComponentAttrs: () => ({
      space: unref(space),
      resource: unref(resource),
      link: linkShare,
      ...(callbackFn && { callbackFn })
    })
  })
}

const availableLinkTypeOptions = computed(() => getAvailableLinkTypes({ isFolder: isFolderShare }))

const currentLinkRoleDescription = computed(() => {
  return getLinkRoleByType(unref(currentLinkType))?.description || ''
})

const currentLinkRoleLabel = computed(() => {
  return getLinkRoleByType(unref(currentLinkType))?.displayName || ''
})
</script>
<style lang="scss" scoped>
.files-links-content {
  min-width: 0;
}

.files-links-name-wrapper {
  min-width: 0;
}
</style>
