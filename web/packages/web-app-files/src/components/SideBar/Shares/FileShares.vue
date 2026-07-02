<template>
  <div id="oc-files-sharing-sidebar" class="oc-position-relative">
    <div class="oc-flex oc-flex-between oc-flex-middle">
      <div class="oc-flex">
        <h4 v-translate class="oc-text-bold oc-text-medium oc-m-rm">Share with people</h4>
        <oc-contextual-helper
          v-if="helpersEnabled"
          class="oc-pl-xs"
          v-bind="inviteCollaboratorHelp"
        />
      </div>
      <copy-private-link :resource="resource" />
    </div>
    <invite-collaborator-form
      v-if="canShare({ resource, space })"
      key="new-collaborator"
      class="oc-my-s"
    />
    <p
      v-else
      key="no-share-permissions-message"
      data-testid="files-collaborators-no-share-permissions-message"
      v-text="noSharePermsMessage"
    />
    <template v-if="hasSharees">
      <div id="files-collaborators-headline" class="oc-flex oc-flex-middle oc-flex-between">
        <h5 class="oc-text-bold oc-text-medium oc-my-rm" v-text="sharedWithLabel" />
      </div>
      <portal-target
        name="app.files.sidebar.sharing.shared-with.top"
        :slot-props="{ space, resource }"
        :multiple="true"
      />
      <ul
        id="files-collaborators-list"
        class="oc-list oc-list-divider"
        :class="{ 'oc-mb-l': showSpaceMembers, 'oc-m-rm': !showSpaceMembers }"
        :aria-label="$gettext('Share receivers')"
      >
        <li v-for="collaborator in displayCollaborators" :key="collaborator.id">
          <collaborator-list-item
            :share="collaborator"
            :resource-name="resource.name"
            :deniable="isShareDeniable(collaborator)"
            :modifiable="isShareModifiable(collaborator)"
            :is-share-denied="isShareDenied(collaborator)"
            :shared-parent-route="getSharedParentRoute(collaborator)"
            :is-locked="resource.locked"
            @on-delete="deleteShareConfirmation"
            @on-set-deny="setDenyShare"
          />
        </li>
        <portal-target
          name="app.files.sidebar.sharing.shared-with.bottom"
          :slot-props="{ space, resource }"
          :multiple="true"
        />
      </ul>
      <div v-if="showShareToggle" class="oc-flex oc-flex-center">
        <oc-button
          appearance="raw"
          class="toggle-shares-list-btn"
          @click="toggleShareListCollapsed"
        >
          {{ collapseButtonTitle }}
        </oc-button>
      </div>
    </template>
    <template v-if="showSpaceMembers">
      <div class="oc-flex oc-flex-middle oc-flex-between">
        <h5 class="oc-text-bold oc-text-medium oc-my-s" v-text="spaceMemberLabel" />
      </div>
      <ul
        id="space-collaborators-list"
        class="oc-list oc-list-divider oc-overflow-hidden oc-m-rm"
        :aria-label="spaceMemberLabel"
      >
        <li v-for="(collaborator, i) in displaySpaceMembers" :key="i">
          <collaborator-list-item
            :share="collaborator"
            :resource-name="resource.name"
            :deniable="isSpaceMemberDeniable(collaborator)"
            :modifiable="false"
            :is-share-denied="isSpaceMemberDenied(collaborator)"
            :is-space-share="true"
            @on-set-deny="setDenyShare"
          />
        </li>
      </ul>
      <div v-if="showMemberToggle" class="oc-flex oc-flex-center">
        <oc-button appearance="raw" @click="toggleMemberListCollapsed">
          {{ collapseMemberButtonTitle }}
        </oc-button>
      </div>
    </template>
  </div>
</template>

<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import {
  useGetMatchingSpace,
  useModals,
  useUserStore,
  useMessages,
  useSpacesStore,
  useConfigStore,
  useSharesStore,
  useResourcesStore,
  useCanShare,
  useClientService,
  useRouter
} from '@ownclouders/web-pkg'
import { isLocationSharesActive } from '@ownclouders/web-pkg'
import { textUtils } from '../../../helpers/textUtils'
import { isShareSpaceResource, ShareTypes } from '@ownclouders/web-client'
import InviteCollaboratorForm from './Collaborators/InviteCollaborator/InviteCollaboratorForm.vue'
import CollaboratorListItem from './Collaborators/ListItem.vue'
import { useContextualHelpers } from '../../../composables/contextualHelpers/useContextualHelpers'
import { computed, inject, ref, Ref, unref } from 'vue'
import {
  isProjectSpaceResource,
  Resource,
  SpaceResource,
  CollaboratorShare,
  isSpaceResource
} from '@ownclouders/web-client'
import { getSharedAncestorRoute } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import CopyPrivateLink from '../../Shares/CopyPrivateLink.vue'

const userStore = useUserStore()
const clientService = useClientService()
const { getMatchingSpace } = useGetMatchingSpace()
const { dispatchModal } = useModals()
const { canShare } = useCanShare()
const { $gettext } = useGettext()
const router = useRouter()
const { showMessage, showErrorMessage } = useMessages()

const resourcesStore = useResourcesStore()
const { removeResources, getAncestorById } = resourcesStore

const { getSpaceMembers } = useSpacesStore()

const configStore = useConfigStore()
const { options: configOptions } = storeToRefs(configStore)

const { shareInviteCollaboratorHelp, shareInviteCollaboratorHelpCern } = useContextualHelpers()

const sharesStore = useSharesStore()
const { addShare, deleteShare } = sharesStore

const { user } = storeToRefs(userStore)

const resource = inject<Ref<Resource>>('resource')
const space = inject<Ref<SpaceResource>>('space')

const collaboratorShares = computed(() => {
  if (isProjectSpaceResource(unref(space))) {
    // filter out project space members, they are listed separately (see down below)
    return sharesStore.collaboratorShares.filter((c) => c.resourceId !== unref(space).id)
  }
  return sharesStore.collaboratorShares
})

const spaceMembers = computed(() => getSpaceMembers(unref(space)))

const sharesListCollapsed = ref(true)
const toggleShareListCollapsed = () => {
  sharesListCollapsed.value = !unref(sharesListCollapsed)
}
const memberListCollapsed = ref(true)
const toggleMemberListCollapsed = () => {
  memberListCollapsed.value = !unref(memberListCollapsed)
}

const matchingSpace = computed(() => {
  return getMatchingSpace(unref(resource))
})

const collaborators = computed(() => {
  const collaboratorsComparator = (c1: CollaboratorShare, c2: CollaboratorShare) => {
    // Sorted by: type, direct, display name, creation date
    const name1 = c1.sharedWith.displayName.toLowerCase().trim()
    const name2 = c2.sharedWith.displayName.toLowerCase().trim()
    const c1UserShare = ShareTypes.containsAnyValue(ShareTypes.individuals, [c1.shareType])
    const c2UserShare = ShareTypes.containsAnyValue(ShareTypes.individuals, [c2.shareType])
    const c1DirectShare = !c1.indirect
    const c2DirectShare = !c2.indirect

    if (c1UserShare === c2UserShare) {
      if (c1DirectShare === c2DirectShare) {
        return textUtils.naturalSortCompare(name1, name2)
      }

      return c1DirectShare ? -1 : 1
    }

    return c1UserShare ? -1 : 1
  }

  return unref(collaboratorShares).sort(collaboratorsComparator)
})

const inviteCollaboratorHelp = computed(() => {
  const cernFeatures = configOptions.value.cernFeatures

  if (cernFeatures) {
    const mergedHelp = { ...unref(shareInviteCollaboratorHelp) }
    mergedHelp.list = [...unref(shareInviteCollaboratorHelpCern).list, ...mergedHelp.list]
    return mergedHelp
  }

  return unref(shareInviteCollaboratorHelp)
})

const helpersEnabled = computed(() => {
  return configOptions.value.contextHelpers
})

const sharedWithLabel = computed(() => {
  return $gettext('Shared with')
})

const spaceMemberLabel = computed(() => {
  return $gettext('Space members')
})

const collapseButtonTitle = computed(() => {
  return unref(sharesListCollapsed) ? $gettext('Show more') : $gettext('Show less')
})

const collapseMemberButtonTitle = computed(() => {
  return memberListCollapsed.value ? $gettext('Show more') : $gettext('Show less')
})

const hasSharees = computed(() => {
  return unref(displayCollaborators).length > 0
})

const displayCollaborators = computed(() => {
  if (unref(collaborators).length > 3 && unref(sharesListCollapsed)) {
    return unref(collaborators).slice(0, 3)
  }

  return unref(collaborators)
})

const displaySpaceMembers = computed(() => {
  if (unref(spaceMembers).length > 3 && unref(memberListCollapsed)) {
    return unref(spaceMembers).slice(0, 3)
  }
  return unref(spaceMembers)
})

const showShareToggle = computed(() => {
  return unref(collaborators).length > 3
})

const showMemberToggle = computed(() => {
  return unref(spaceMembers).length > 3
})

const noSharePermsMessage = computed(() => {
  const translatedFile = $gettext("You don't have permission to share this file.")
  const translatedFolder = $gettext("You don't have permission to share this folder.")
  return unref(resource).type === 'file' ? translatedFile : translatedFolder
})

const showSpaceMembers = computed(() => {
  return (
    unref(space)?.driveType === 'project' &&
    unref(resource).type !== 'space' &&
    unref(space)?.isMember(unref(user))
  )
})

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function getDeniedShare(_collaborator: CollaboratorShare) {
  // FIXME: currently not supported by sharing NG
  return undefined
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function isShareDenied(_collaborator: CollaboratorShare) {
  // FIXME: currently not supported by sharing NG
  return false
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function getDeniedSpaceMember(_collaborator: CollaboratorShare) {
  // FIXME: currently not supported by sharing NG
  return undefined
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function isSpaceMemberDenied(_collaborator: CollaboratorShare) {
  // FIXME: currently not supported by sharing NG
  return false
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function isSpaceMemberDeniable(_collaborator: CollaboratorShare) {
  // FIXME: currently not supported by sharing NG
  return false
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
function isShareDeniable(_collaborator: CollaboratorShare) {
  // FIXME: currently not supported by sharing NG
  return false
}

async function setDenyShare({ value, share }: { value: boolean; share: CollaboratorShare }) {
  if (value === true) {
    try {
      await addShare({
        clientService,
        space: unref(space),
        resource: unref(resource),
        options: {}
      })
      showMessage({
        title: $gettext('Access was denied successfully')
      })
    } catch (e) {
      console.error(e)
      showErrorMessage({
        title: $gettext('Failed to deny access'),
        errors: [e]
      })
    }
  } else {
    try {
      await deleteShare({
        clientService,
        space: unref(space),
        resource: unref(resource),
        collaboratorShare: isSpaceResource(unref(resource))
          ? getDeniedSpaceMember(share)
          : getDeniedShare(share)
      })
      showMessage({
        title: $gettext('Access was granted successfully')
      })
    } catch (e) {
      console.error(e)
      showErrorMessage({
        title: $gettext('Failed to grant access'),
        errors: [e]
      })
    }
  }
}

function deleteShareConfirmation(collaboratorShare: CollaboratorShare) {
  dispatchModal({
    variation: 'danger',
    title: $gettext('Remove share'),
    confirmText: $gettext('Remove'),
    message: $gettext('Are you sure you want to remove this share?'),
    hasInput: false,
    onConfirm: async () => {
      const lastShareId = unref(collaborators).length === 1 ? unref(collaborators)[0].id : undefined

      try {
        await deleteShare({
          clientService,
          space: unref(space),
          resource: unref(resource),
          collaboratorShare
        })

        showMessage({
          title: $gettext('Share was removed successfully')
        })
        if (lastShareId && isLocationSharesActive(router, 'files-shares-with-others')) {
          removeResources([{ id: lastShareId }] as Resource[])
        }
      } catch (error) {
        console.error(error)
        showErrorMessage({
          title: $gettext('Failed to remove share'),
          errors: [error]
        })
      }
    }
  })
}

function getSharedParentRoute(collaborator: CollaboratorShare) {
  if (!collaborator.indirect) {
    return null
  }
  const sharedAncestor = getAncestorById(collaborator.resourceId)
  if (!sharedAncestor) {
    return null
  }

  return getSharedAncestorRoute({
    sharedAncestor,
    matchingSpace: unref(space) || unref(matchingSpace)
  })
}

function isShareModifiable(collaborator: CollaboratorShare) {
  if (collaborator.indirect) {
    return false
  }

  if (isProjectSpaceResource(unref(space)) || isShareSpaceResource(unref(space))) {
    return unref(space).canShare({ user: unref(user) })
  }

  return true
}
</script>

<style lang="scss" scoped>
#files-collaborators-headline {
  height: 40px;
}
</style>
