<template>
  <oc-recipient
    :data-testid="`recipient-container-${formattedRecipient.name}`"
    class="files-share-invite-recipient"
    :recipient="formattedRecipient"
  >
    <template #append>
      <oc-button
        class="files-share-invite-recipient-btn-remove"
        appearance="raw"
        :aria-label="btnDeselectRecipientLabel"
        @click.stop="deselect(recipient)"
      >
        <oc-icon name="close" size="small" />
      </oc-button>
    </template>
  </oc-recipient>
</template>

<script lang="ts" setup>
import { useGettext } from 'vue3-gettext'
import { CollaboratorAutoCompleteItem, ShareTypes } from '@ownclouders/web-client'
import { computed, unref, ref } from 'vue'
import { Recipient } from '@ownclouders/design-system/helpers'
import { useCapabilityStore, useConfigStore, useClientService } from '@ownclouders/web-pkg'
import { storeToRefs } from 'pinia'
import { avatarUrl } from '../../../../../helpers/user'

interface Props {
  recipient: CollaboratorAutoCompleteItem
  deselect?: (recipient: CollaboratorAutoCompleteItem) => void
}
const { recipient, deselect = null } = defineProps<Props>()
const capabilityStore = useCapabilityStore()
const capabilityRefs = storeToRefs(capabilityStore)
const clientService = useClientService()

const configStore = useConfigStore()
const { serverUrl } = storeToRefs(configStore)

const { $gettext } = useGettext()
const externalIssuer = computed(() => {
  if (recipient.shareType === ShareTypes.remote.value) {
    return recipient.identities?.[0]?.issuer
  }
  return ''
})

const userProfilePicture = unref(capabilityRefs.sharingUserProfilePicture)
const btnDeselectRecipientLabel = computed(() => {
  return $gettext('Deselect %{name}', { name: recipient.displayName })
})

const recipientName = computed(() => {
  let name = recipient.displayName
  if (unref(externalIssuer)) {
    name += ` (${unref(externalIssuer)})`
  }

  return name
})

const formattedRecipient = ref<Recipient>({
  name: unref(recipientName),
  icon: getRecipientIcon(),
  hasAvatar: recipient.shareType === ShareTypes.user.value,
  isLoadingAvatar: true
})

function getRecipientIcon() {
  switch (recipient.shareType) {
    case ShareTypes.group.value:
      return {
        name: ShareTypes.group.icon,
        label: $gettext('Group')
      }

    case ShareTypes.guest.value:
      return {
        name: ShareTypes.guest.icon,
        label: $gettext('Guest user')
      }

    case ShareTypes.remote.value:
      return {
        name: ShareTypes.remote.icon,
        label: $gettext('External user')
      }

    default:
      return {
        name: ShareTypes.user.icon,
        label: $gettext('User')
      }
  }
}

async function onCreated() {
  if (userProfilePicture && unref(formattedRecipient).hasAvatar) {
    try {
      unref(formattedRecipient).avatar = await avatarUrl({
        clientService,
        server: unref(serverUrl),
        username: recipient.displayName
      })
    } catch (error) {
      console.error(error)
    }
  }

  unref(formattedRecipient).isLoadingAvatar = false
}

onCreated()
</script>

<style lang="scss">
.files-share-invite-recipient {
  margin: 4px 2px 0;
  padding: 0 0.25em;
  overflow-wrap: anywhere;

  .oc-recipient-avatar {
    min-width: 16.8px;
  }
}
</style>
