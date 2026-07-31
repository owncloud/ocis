<template>
  <div
    :data-testid="`collaborator-${isAnyUserShareType ? 'user' : 'group'}-item-${
      share.sharedWith.displayName
    }`"
    class="files-collaborators-collaborator oc-py-xs"
  >
    <div class="oc-width-1-1 oc-flex oc-flex-middle files-collaborators-collaborator-details">
      <div class="oc-width-2-3 oc-flex oc-flex-middle">
        <div>
          <template v-if="isShareDenied">
            <oc-avatar-item
              :width="36"
              icon-size="medium"
              icon="stop"
              :name="$gettext('Access denied')"
              class="files-collaborators-collaborator-indicator"
            />
          </template>
          <template v-else>
            <avatar-image
              v-if="isAnyUserShareType"
              :userid="share.sharedWith.id"
              :user-name="share.sharedWith.displayName"
              :width="36"
              class="files-collaborators-collaborator-indicator"
            />
            <oc-avatar-item
              v-else
              :width="36"
              icon-size="medium"
              :icon="shareTypeIcon"
              :name="shareTypeKey"
              class="files-collaborators-collaborator-indicator"
            />
          </template>
        </div>
        <div class="files-collaborators-collaborator-name-wrapper oc-pl-s">
          <div class="oc-text-truncate">
            <span
              aria-hidden="true"
              class="files-collaborators-collaborator-name"
              v-text="shareDisplayName"
            />
            <span class="oc-invisible-sr" v-text="screenreaderShareDisplayName" />
            <oc-contextual-helper
              v-if="isExternalShare"
              :text="
                $gettext(
                  'External user, registered with another organization’s account but granted access to your resources. External users can only have “view” or “edit” permission.'
                )
              "
              :title="$gettext('External user')"
            />
          </div>
          <div v-if="isExternalShare" class="oc-text-small" data-testid="external-share-domain">
            {{ externalShareDomainName }}
          </div>
          <div>
            <div
              v-if="isShareDenied"
              v-oc-tooltip="shareDeniedTooltip"
              class="oc-flex oc-flex-nowrap oc-flex-middle"
              v-text="$gettext('Access denied')"
            />
            <template v-else>
              <div v-if="modifiable" class="oc-flex oc-flex-nowrap oc-flex-middle">
                <role-dropdown
                  :dom-selector="shareDomSelector"
                  :existing-share-role="share.role"
                  :existing-share-permissions="share.permissions"
                  :is-locked="isLocked"
                  :is-external="isExternalShare"
                  class="files-collaborators-collaborator-role"
                  mode="edit"
                  @option-change="shareRoleChanged"
                />
              </div>
              <div v-else-if="share.role">
                <span
                  v-oc-tooltip="$gettext(share.role.description)"
                  class="oc-mr-xs"
                  v-text="$gettext(share.role.displayName)"
                />
              </div>
            </template>
          </div>
        </div>
      </div>
      <div class="oc-flex oc-flex-middle oc-width-1-3 files-collaborators-collaborator-navigation">
        <expiration-date-indicator
          v-if="hasExpirationDate"
          class="files-collaborators-collaborator-expiration oc-mr-xs"
          data-testid="recipient-info-expiration-date"
          :expiration-date="DateTime.fromISO(share.expirationDateTime)"
        />
        <oc-icon
          v-if="!isShareDenied && sharedParentRoute"
          v-oc-tooltip="sharedViaTooltip"
          name="folder-shared"
          fill-type="line"
          class="files-collaborators-collaborator-shared-via oc-mx-xs"
        />
        <edit-dropdown
          class="files-collaborators-collaborator-edit oc-ml-xs"
          data-testid="collaborator-edit"
          :expiration-date="share.expirationDateTime ? share.expirationDateTime : null"
          :share-category="shareCategory"
          :can-edit="modifiable"
          :is-share-denied="isShareDenied"
          :is-locked="isLocked"
          :deniable="deniable"
          :shared-parent-route="!isShareDenied ? sharedParentRoute : undefined"
          :access-details="accessDetails"
          @expiration-date-changed="shareExpirationChanged"
          @remove-share="removeShare"
          @set-deny-share="setDenyShare"
          @notify-share="showNotifyShareModal"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { DateTime } from 'luxon'

import EditDropdown from './EditDropdown.vue'
import RoleDropdown from './RoleDropdown.vue'
import { CollaboratorShare, ShareRole, ShareTypes } from '@ownclouders/web-client'
import {
  queryItemAsString,
  useMessages,
  useModals,
  useSpacesStore,
  useUserStore,
  useSharesStore
} from '@ownclouders/web-pkg'
import { Resource, extractDomSelector } from '@ownclouders/web-client'
import { computed, inject, Ref, unref } from 'vue'
import { formatDateFromDateTime } from '@ownclouders/web-pkg'
import { useClientService } from '@ownclouders/web-pkg'
import { RouteLocationNamedRaw } from 'vue-router'
import { useGettext } from 'vue3-gettext'
import { SpaceResource } from '@ownclouders/web-client'
import { isProjectSpaceResource } from '@ownclouders/web-client'
import { ContextualHelperDataListItem } from '@ownclouders/design-system/helpers'
import ExpirationDateIndicator from '../ExpirationDateIndicator.vue'

interface Props {
  share: CollaboratorShare
  isShareDenied?: boolean
  modifiable?: boolean
  sharedParentRoute?: RouteLocationNamedRaw | null
  resourceName?: string
  deniable?: boolean
  isLocked?: boolean
  isSpaceShare?: boolean
}
interface Emits {
  (e: 'onDelete', share: CollaboratorShare): void
  (e: 'onSetDeny', payload: { share: CollaboratorShare; value: boolean }): void
}

const {
  share,
  isShareDenied = false,
  modifiable = false,
  sharedParentRoute = null,
  resourceName = '',
  deniable = false,
  isLocked = false,
  isSpaceShare = false
} = defineProps<Props>()

const emit = defineEmits<Emits>()
const space = inject<Ref<SpaceResource>>('space')
const resource = inject<Ref<Resource>>('resource')
const { showMessage, showErrorMessage } = useMessages()
const userStore = useUserStore()
const clientService = useClientService()
const language = useGettext()
const { $gettext } = language
const { dispatchModal } = useModals()

const sharesStore = useSharesStore()
const { graphRoles } = storeToRefs(sharesStore)
const { updateShare } = sharesStore
const { upsertSpace } = useSpacesStore()

const { user } = storeToRefs(userStore)

const sharedParentDir = computed(() => {
  return queryItemAsString(sharedParentRoute?.params?.driveAliasAndItem).split('/').pop()
})

const shareDate = computed(() => {
  return formatDateFromDateTime(DateTime.fromISO(share.createdDateTime), language.current)
})

const isExternalShare = computed(() => share.shareType === ShareTypes.remote.value)

const setDenyShare = (value: boolean) => {
  emit('onSetDeny', { share: share, value })
}

const showNotifyShareModal = () => {
  dispatchModal({
    variation: 'warning',
    icon: 'mail-send',
    title: $gettext('Send a reminder'),
    confirmText: $gettext('Send'),
    message: $gettext('Are you sure you want to send a reminder about this share?'),
    onConfirm: notifyShare
  })
}
const notifyShare = async () => {
  // FIXME: cern code
  // const response = await clientService.owncloudSdk.shares.notifyShare(props.share.id)
}

const sharedViaTooltip = computed(() =>
  $gettext('Shared via the parent folder "%{sharedParentDir}"', {
    sharedParentDir: unref(sharedParentDir)
  })
)

const shareType = computed(() => ShareTypes.getByValue(share.shareType))

const shareTypeIcon = computed(() => unref(shareType).icon)

const shareTypeKey = computed(() => unref(shareType).key)

const shareDomSelector = computed(() => {
  if (!share.id) {
    return undefined
  }
  return extractDomSelector(share.id)
})

const isAnyUserShareType = computed(() => ShareTypes.user === unref(shareType))

const shareTypeText = computed(() => $gettext(unref(shareType).label))

const shareCategory = computed(() => (ShareTypes.isIndividual(unref(shareType)) ? 'user' : 'group'))

const shareDeniedTooltip = computed(() => {
  return $gettext('%{shareType} cannot access %{resourceName}', {
    shareType: unref(shareTypeText),
    resourceName: resourceName
  })
})

const shareDisplayName = computed(() => {
  if (unref(user).id === share.sharedWith.id) {
    return $gettext('%{collaboratorName} (me)', {
      collaboratorName: share.sharedWith.displayName
    })
  }
  return share.sharedWith.displayName
})

const screenreaderShareDisplayName = computed(() => {
  const context = {
    displayName: share.sharedWith.displayName
  }

  return $gettext('Share receiver name: %{ displayName }', context)
})

const hasExpirationDate = computed(() => !!share.expirationDateTime)

const expirationDate = computed(() => {
  return formatDateFromDateTime(
    DateTime.fromISO(share.expirationDateTime).endOf('day'),
    language.current
  )
})

const shareOwnerDisplayName = computed(() => share.sharedBy.displayName)

const externalShareDomainName = computed(() => {
  if (unref(isExternalShare)) {
    const [, serverUrl] = share.sharedWith.id.split('@')

    return serverUrl
  }

  return null
})

const accessDetails = computed(() => {
  const list: ContextualHelperDataListItem[] = []

  list.push({ text: $gettext('Name'), headline: true }, { text: unref(shareDisplayName) })
  unref(isExternalShare) &&
    list.push(
      { text: $gettext('Domain'), headline: true },
      { text: unref(externalShareDomainName) }
    )

  list.push({ text: $gettext('Type'), headline: true }, { text: unref(shareTypeText) })
  list.push(
    { text: $gettext('Access expires'), headline: true },
    { text: unref(hasExpirationDate) ? unref(expirationDate) : $gettext('no') }
  )
  list.push({ text: $gettext('Shared on'), headline: true }, { text: unref(shareDate) })

  if (!isSpaceShare) {
    list.push(
      { text: $gettext('Invited by'), headline: true },
      { text: unref(shareOwnerDisplayName) }
    )
  }

  return list
})

function removeShare() {
  emit('onDelete', share)
}

async function shareRoleChanged(role: ShareRole) {
  const expirationDateTime = share.expirationDateTime
  try {
    await saveShareChanges({ role, expirationDateTime })
  } catch (e) {
    console.error(e)
    showErrorMessage({
      title: $gettext('Failed to apply new permissions'),
      errors: [e]
    })
  }
}

async function shareExpirationChanged({ expirationDateTime }: { expirationDateTime: string }) {
  const role = share.role
  try {
    await saveShareChanges({ role, expirationDateTime })
  } catch (e) {
    console.error(e)
    showErrorMessage({
      title: $gettext('Failed to apply expiration date'),
      errors: [e]
    })
  }
}

async function saveShareChanges({
  role,
  expirationDateTime
}: {
  role: ShareRole
  expirationDateTime?: string
}) {
  try {
    await updateShare({
      clientService,
      space: unref(space),
      resource: unref(resource),
      collaboratorShare: share,
      options: { roles: [role.id], expirationDateTime }
    })

    if (isProjectSpaceResource(unref(resource))) {
      const client = clientService.graphAuthenticated
      const space = await client.drives.getDrive(unref(resource).id, unref(graphRoles))

      upsertSpace(space)
    }

    showMessage({ title: $gettext('Share successfully changed') })
  } catch (e) {
    console.error(e)
    showErrorMessage({
      title: $gettext('Error while editing the share.'),
      errors: [e]
    })
  }
}
</script>

<style lang="scss" scoped>
.sharee-avatar {
  min-width: 36px;
}

.files-collaborators-collaborator-navigation {
  align-items: center;
  justify-content: end;
}

.files-collaborators-collaborator-role {
  max-width: 100%;
}

.files-collaborators-collaborator-name-wrapper {
  max-width: 100%;
}
</style>
