<template>
  <div id="oc-files-file-link" class="oc-position-relative">
    <div class="oc-flex oc-flex-middle">
      <h4 class="oc-text-bold oc-text-medium oc-m-rm" v-text="$gettext('Public links')" />
      <oc-contextual-helper v-if="helpersEnabled" class="oc-pl-xs" v-bind="shareViaLinkHelp" />
    </div>
    <p v-if="!directLinks.length" class="files-links-empty oc-mt-m" v-text="noLinksLabel" />
    <ul
      v-else
      id="files-links-list"
      class="oc-list oc-list-divider oc-mt-m"
      :aria-label="$gettext('Public links')"
    >
      <li v-for="link in displayLinks" :key="link.id">
        <list-item
          :can-rename="true"
          :is-folder-share="resource.isFolder"
          :is-modifiable="canEditLink(link)"
          :is-password-enforced="isPasswordEnforcedForLinkType(link.type)"
          :is-password-removable="canDeletePublicLinkPassword(link)"
          :link-share="link"
          @update-link="handleLinkUpdate"
          @remove-public-link="deleteLinkConfirmation"
        />
      </li>
    </ul>
    <div v-if="directLinks.length > 3" class="oc-flex oc-flex-center">
      <oc-button class="indirect-link-list-toggle" appearance="raw" @click="toggleLinkListCollapsed"
        ><span v-text="collapseButtonTitle"
      /></oc-button>
    </div>
    <div class="oc-mt-m">
      <oc-button
        v-if="canCreateLinks"
        id="files-file-link-add"
        variation="primary"
        appearance="raw"
        data-testid="files-link-add-btn"
        @click="addNewLink"
      >
        <span v-text="$gettext('Add link')"
      /></oc-button>
      <p
        v-else
        data-testid="files-links-no-share-permissions-message"
        class="oc-mt-m"
        v-text="$gettext('You do not have permission to create public links.')"
      />
    </div>
    <div v-if="indirectLinks.length" class="files-links-indirect oc-mt-m">
      <hr class="oc-my-m" />
      <h5 class="oc-text-bold oc-text-medium oc-m-rm">
        {{ indirectLinksHeading }}
        <oc-contextual-helper
          v-if="helpersEnabled"
          class="oc-pl-xs"
          v-bind="shareViaIndirectLinkHelp"
        />
      </h5>
      <div
        class="files-links-indirect-list"
        :class="{ 'files-links-indirect-list-open': !indirectLinkListCollapsed }"
      >
        <ul class="oc-list oc-list-divider" :aria-label="$gettext('Public links')">
          <li v-for="link in indirectLinks" :key="link.id">
            <list-item
              :is-folder-share="resource.isFolder"
              :is-modifiable="false"
              :link-share="link"
            />
          </li>
        </ul>
      </div>
      <div class="oc-flex oc-flex-center">
        <oc-button
          class="indirect-link-list-toggle"
          appearance="raw"
          @click="toggleIndirectLinkListCollapsed"
        >
          <span v-text="indirectCollapseButtonTitle" />
        </oc-button>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { computed, inject, ref, Ref, unref } from 'vue'
import {
  useAbility,
  useFileActionsCreateLink,
  FileAction,
  useClientService,
  useModals,
  useMessages,
  useConfigStore,
  useResourcesStore,
  useLinkTypes,
  useCanShare,
  UpdateLinkOptions,
  useRouter
} from '@ownclouders/web-pkg'
import { useContextualHelpers } from '../../../composables/contextualHelpers/useContextualHelpers'
import { isSpaceResource, LinkShare } from '@ownclouders/web-client'
import ListItem from './Links/ListItem.vue'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { isLocationSharesActive, useSharesStore } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { storeToRefs } from 'pinia'
import { SharingLinkType } from '@ownclouders/web-client/graph/generated'

const { showMessage, showErrorMessage } = useMessages()
const { $gettext } = useGettext()
const ability = useAbility()
const clientService = useClientService()
const { can } = ability
const { dispatchModal } = useModals()
const { removeResources } = useResourcesStore()
const { isPasswordEnforcedForLinkType } = useLinkTypes()
const { canShare } = useCanShare()

const canCreateLinks = computed(() => {
  if (!ability.can('create-all', 'PublicLink')) {
    return false
  }
  return canShare({ space: unref(space), resource: unref(resource) })
})

const sharesStore = useSharesStore()
const router = useRouter()
const { updateLink, deleteLink } = sharesStore
const { linkShares } = storeToRefs(sharesStore)

const configStore = useConfigStore()
const { options: configOptions } = storeToRefs(configStore)

const { shareViaLinkHelp, shareViaIndirectLinkHelp } = useContextualHelpers()

const { actions: createLinkActions } = useFileActionsCreateLink()
const createLinkAction = computed<FileAction>(() =>
  unref(createLinkActions).find(({ name }) => name === 'create-links')
)

const space = inject<Ref<SpaceResource>>('space')
const resource = inject<Ref<Resource>>('resource')

const linkListCollapsed = ref(true)
const indirectLinkListCollapsed = ref(true)
const directLinks = computed(() =>
  unref(linkShares)
    .filter((l) => !l.indirect)
    .sort((a, b) => b.createdDateTime.localeCompare(a.createdDateTime))
    .map((share) => {
      return { ...share, key: 'direct-link-' + share.id }
    })
)
const indirectLinks = computed(() =>
  unref(linkShares)
    .filter((l) => l.indirect)
    .sort((a, b) => b.createdDateTime.localeCompare(a.createdDateTime))
    .map((share) => {
      return { ...share, key: 'indirect-link-' + share.id }
    })
)

const canDeleteReadOnlyPublicLinkPassword = computed(() =>
  can('delete-all', 'ReadOnlyPublicLinkPassword')
)

const canEditLink = (linkShare: LinkShare) => {
  return (
    unref(canCreateLinks) &&
    (can('create-all', 'PublicLink') || linkShare.type === SharingLinkType.Internal)
  )
}

const addNewLink = () => {
  const handlerArgs = { space: unref(space), resources: [unref(resource)] }
  return unref(createLinkAction)?.handler(handlerArgs)
}

const canDeletePublicLinkPassword = (linkShare: LinkShare) => {
  const isPasswordEnforced = isPasswordEnforcedForLinkType(linkShare.type)

  if (!isPasswordEnforced) {
    return true
  }

  return linkShare.type === SharingLinkType.View && unref(canDeleteReadOnlyPublicLinkPassword)
}

const handleLinkUpdate = async ({
  linkShare,
  options
}: {
  linkShare: LinkShare
  options: UpdateLinkOptions['options']
}) => {
  try {
    await updateLink({
      clientService,
      space: unref(space),
      resource: unref(resource),
      linkShare,
      options
    })
    showMessage({ title: $gettext('Link was updated successfully') })
  } catch (e) {
    console.error(e)
    showErrorMessage({
      title: $gettext('Failed to update link'),
      errors: [e]
    })
  }
}

const toggleLinkListCollapsed = () => {
  linkListCollapsed.value = !unref(linkListCollapsed)
}

const toggleIndirectLinkListCollapsed = () => {
  indirectLinkListCollapsed.value = !unref(indirectLinkListCollapsed)
}

const noLinksLabel = computed(() => {
  if (isSpaceResource(unref(resource))) {
    return $gettext('This space has no public links.')
  }
  if (unref(resource).isFolder) {
    return $gettext('This folder has no public link.')
  }
  return $gettext('This file has no public link.')
})

function deleteLinkConfirmation({ link }: { link: LinkShare }) {
  dispatchModal({
    variation: 'danger',
    title: $gettext('Delete link'),
    message: $gettext(
      'Are you sure you want to delete this link? Recreating the same link again is not possible.'
    ),
    confirmText: $gettext('Delete'),
    onConfirm: async () => {
      let lastLinkId = unref(linkShares).length === 1 ? unref(linkShares)[0].id : undefined

      try {
        await deleteLink({
          clientService: clientService,
          space: unref(space),
          resource: unref(resource),
          linkShare: link
        })

        showMessage({ title: $gettext('Link was deleted successfully') })

        if (lastLinkId && isLocationSharesActive(router, 'files-shares-via-link')) {
          if (isSpaceResource(unref(resource))) {
            // spaces need their actual id instead of their share id to be removed from the file list
            lastLinkId = unref(resource).id.toString()
          }
          removeResources([{ id: lastLinkId }] as Resource[])
        }
      } catch (e) {
        console.error(e)
        showErrorMessage({
          title: $gettext('Failed to delete link'),
          errors: [e]
        })
      }
    }
  })
}
const collapseButtonTitle = computed(() => {
  return unref(linkListCollapsed) ? $gettext('Show more') : $gettext('Show less')
})
const indirectCollapseButtonTitle = computed(() => {
  return unref(indirectLinkListCollapsed) ? $gettext('Show') : $gettext('Hide')
})

const helpersEnabled = computed(() => {
  return unref(configOptions).contextHelpers
})

const indirectLinksHeading = computed(() => {
  return $gettext('Indirect (%{ count })', {
    count: unref(indirectLinks).length.toString()
  })
})

const displayLinks = computed(() => {
  if (unref(directLinks).length > 3 && unref(linkListCollapsed)) {
    return unref(directLinks).slice(0, 3)
  }
  return unref(directLinks)
})
</script>
<style lang="scss">
#oc-files-file-link,
#oc-files-sharing-sidebar {
  border-radius: 5px;
}

.files-links-indirect-list {
  display: grid;
  grid-template-rows: 0fr;
  transition: all 0.25s ease-out;

  ul {
    overflow: hidden;
  }

  &-open {
    grid-template-rows: 1fr;
    margin-top: var(--oc-space-medium);
  }
}
</style>
