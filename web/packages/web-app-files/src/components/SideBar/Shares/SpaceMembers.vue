<template>
  <div id="oc-files-sharing-sidebar" class="oc-position-relative">
    <div class="oc-flex">
      <div v-if="canShare({ space: resource, resource })" class="oc-flex oc-py-s">
        <h4 class="oc-text-bold oc-text-medium oc-m-rm" v-text="$gettext('Add members')" />
        <oc-contextual-helper
          v-if="helpersEnabled"
          class="oc-pl-xs"
          v-bind="shareSpaceAddMemberHelp"
        />
      </div>
      <copy-private-link :resource="resource" />
    </div>
    <invite-collaborator-form
      v-if="canShare({ space: resource, resource })"
      key="new-collaborator"
      :save-button-label="$gettext('Add')"
      :invite-label="$gettext('Search')"
      class="oc-my-s"
    />
    <template v-if="hasCollaborators">
      <div
        id="files-collaborators-headline"
        class="oc-flex oc-flex-middle oc-flex-between oc-position-relative"
      >
        <div class="oc-flex">
          <h5 class="oc-text-bold oc-text-medium oc-my-rm" v-text="$gettext('Members')" />
          <oc-button
            v-oc-tooltip="$gettext('Filter members')"
            class="open-filter-btn oc-ml-s"
            :aria-label="$gettext('Filter members')"
            appearance="raw"
            :aria-expanded="isFilterOpen"
            @click="toggleFilter"
          >
            <oc-icon name="search" fill-type="line" size="small" />
          </oc-button>
        </div>
      </div>
      <div
        class="oc-flex oc-flex-between space-members-filter-container"
        :class="{ 'space-members-filter-container-expanded': isFilterOpen }"
      >
        <oc-text-input
          ref="filterInput"
          v-model="filterTerm"
          class="oc-text-truncate space-members-filter oc-mr-s oc-width-1-1"
          :label="$gettext('Filter members')"
          :clear-button-enabled="true"
        />
        <oc-button
          v-oc-tooltip="$gettext('Close filter')"
          class="close-filter-btn oc-mt-m"
          :aria-label="$gettext('Close filter')"
          appearance="raw"
          @click="toggleFilter"
        >
          <oc-icon name="arrow-up-s" fill-type="line" />
        </oc-button>
      </div>

      <ul
        id="files-collaborators-list"
        ref="collaboratorList"
        class="oc-list oc-list-divider oc-m-rm"
        :aria-label="$gettext('Space members')"
      >
        <li v-for="collaborator in filteredSpaceMembers" :key="collaborator.id">
          <collaborator-list-item
            :share="collaborator"
            :modifiable="isModifiable(collaborator)"
            :is-space-share="true"
            @on-delete="deleteMemberConfirm(collaborator)"
          />
        </li>
      </ul>
    </template>
  </div>
</template>

<script lang="ts" setup>
import Fuse from 'fuse.js'
import Mark from 'mark.js'
import { useGettext } from 'vue3-gettext'
import { storeToRefs } from 'pinia'
import { computed, inject, nextTick, ref, Ref, unref, useTemplateRef, watch } from 'vue'
import {
  createLocationSpaces,
  isLocationSpacesActive,
  useCanShare,
  useConfigStore,
  useMessages,
  useModals,
  useSharesStore,
  useSpacesStore,
  useUserStore,
  defaultFuseOptions,
  useClientService,
  useRouter
} from '@ownclouders/web-pkg'
import { useContextualHelpers } from '../../../composables/contextualHelpers/useContextualHelpers'
import {
  ProjectSpaceResource,
  CollaboratorShare,
  GraphSharePermission
} from '@ownclouders/web-client'
import { OcTextInput } from '@ownclouders/design-system/components'
import CopyPrivateLink from '../../Shares/CopyPrivateLink.vue'
import CollaboratorListItem from './Collaborators/ListItem.vue'
import InviteCollaboratorForm from './Collaborators/InviteCollaborator/InviteCollaboratorForm.vue'

const filterInput = useTemplateRef<typeof OcTextInput>('filterInput')
const collaboratorList = useTemplateRef('collaboratorList')

const userStore = useUserStore()
const router = useRouter()
const clientService = useClientService()
const { canShare } = useCanShare()
const { dispatchModal } = useModals()
const sharesStore = useSharesStore()
const { $gettext } = useGettext()
const { deleteShare } = sharesStore
const { graphRoles } = storeToRefs(sharesStore)
const spacesStore = useSpacesStore()
const { showMessage, showErrorMessage } = useMessages()
const { upsertSpace, getSpaceMembers } = spacesStore

const configStore = useConfigStore()
// const { options: configOptions } = storeToRefs(configStore)

const { shareSpaceAddMemberHelp } = useContextualHelpers()

const { user } = storeToRefs(userStore)

const markInstance = ref<Mark>()
const filterTerm = ref('')
const isFilterOpen = ref(false)

const resource = inject<Ref<ProjectSpaceResource>>('resource')

const spaceMembers = computed(() => getSpaceMembers(unref(resource)))

const filteredSpaceMembers = computed(() => {
  return filter(unref(spaceMembers), unref(filterTerm))
})
const helpersEnabled = computed(() => {
  return unref(configStore.options.contextHelpers)
})
const hasCollaborators = computed(() => {
  return unref(spaceMembers).length > 0
})
function filter(collection: CollaboratorShare[], term: string) {
  if (!(term || '').trim()) {
    return collection
  }
  const searchEngine = new Fuse(collection, {
    ...defaultFuseOptions,
    keys: ['sharedWith.displayName', 'sharedWith.name']
  })

  return searchEngine.search(term).map((r) => r.item)
}
watch(isFilterOpen, () => {
  filterTerm.value = ''
})
watch(filterTerm, () => {
  nextTick(() => {
    if (unref(collaboratorList)) {
      markInstance.value = new Mark(unref(collaboratorList) as HTMLElement)
      markInstance.value.unmark()
      markInstance.value.mark(unref(filterTerm), {
        element: 'span',
        className: 'mark-highlight'
      })
    }
  })
})

async function toggleFilter() {
  isFilterOpen.value = !unref(isFilterOpen)
  if (unref(isFilterOpen)) {
    await nextTick()
    unref(filterInput)?.focus()
  }
}
function isModifiable(share: CollaboratorShare) {
  if (!canShare({ space: unref(resource), resource: unref(resource) })) {
    return false
  }

  const memberCanUpdateMembers = share.permissions.includes(GraphSharePermission.updatePermissions)
  if (!memberCanUpdateMembers) {
    return true
  }

  // make sure at least one member can edit other members
  const managers = unref(spaceMembers).filter(({ permissions }) =>
    permissions.includes(GraphSharePermission.updatePermissions)
  )
  return managers.length > 1
}

function deleteMemberConfirm(share: CollaboratorShare) {
  dispatchModal({
    variation: 'danger',
    title: $gettext('Remove member'),
    confirmText: $gettext('Remove'),
    message: $gettext('Are you sure you want to remove this member?'),
    hasInput: false,
    onConfirm: async () => {
      try {
        const currentUserRemoved = share.sharedWith.id === unref(user).id
        await deleteShare({
          clientService,
          space: unref(resource),
          resource: unref(resource),
          collaboratorShare: share
        })

        if (!currentUserRemoved) {
          const client = clientService.graphAuthenticated
          const space = await client.drives.getDrive(share.resourceId, unref(graphRoles))
          upsertSpace(space)
        }

        showMessage({
          title: $gettext('Share was removed successfully')
        })

        if (currentUserRemoved) {
          if (isLocationSpacesActive(router, 'files-spaces-projects')) {
            await router.go(0)
            return
          }
          await router.push(createLocationSpaces('files-spaces-projects'))
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
</script>

<style lang="scss">
#oc-files-sharing-sidebar {
  .copy-private-link {
    margin-left: auto;
  }
}

.space-members-filter {
  label {
    font-size: var(--oc-font-size-small);
  }

  &-container {
    max-height: 0px;
    visibility: hidden;
    transition:
      max-height 0.25s ease-in-out,
      margin-bottom 0.25s ease-in-out,
      visibility 0.25s ease-in-out;

    &-expanded {
      max-height: 60px;
      visibility: visible;
      transition:
        max-height 0.25s ease-in-out,
        margin-bottom 0.25s ease-in-out,
        visibility 0s;
      margin-bottom: var(--oc-space-medium);
    }
  }
}

#files-collaborators-headline {
  height: 40px;
}
</style>
