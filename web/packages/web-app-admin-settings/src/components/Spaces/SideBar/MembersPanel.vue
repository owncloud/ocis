<template>
  <div class="oc-ml-s">
    <oc-text-input
      v-model="filterTerm"
      class="oc-text-truncate oc-mr-s oc-mt-m"
      :label="$gettext('Filter members')"
    />
    <div ref="membersListRef" data-testid="space-members">
      <div v-if="!filteredSpaceMembers.length">
        <h3 class="oc-text-bold oc-text-medium" v-text="$gettext('No members found')" />
      </div>
      <div v-for="(role, i) in availableRoles" :key="i">
        <div
          v-if="getMembersForRole(role).length"
          class="oc-mb-m"
          :data-testid="`space-members-role-${role.displayName}`"
        >
          <h3 class="oc-text-bold oc-text-medium" v-text="role.displayName" />
          <members-role-section :members="getMembersForRole(role)" />
        </div>
      </div>
      <div v-if="membersWithoutRole.length" class="space-members-custom">
        <h3 class="oc-text-bold oc-text-medium" v-text="$gettext('Custom role')" />
        <members-role-section :members="membersWithoutRole" />
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { computed, inject, ref, watch, unref } from 'vue'
import { ShareRole, SpaceMember, SpaceResource } from '@ownclouders/web-client'
import MembersRoleSection from './MembersRoleSection.vue'
import Fuse from 'fuse.js'
import Mark from 'mark.js'
import { defaultFuseOptions, useSharesStore } from '@ownclouders/web-pkg'

const sharesStore = useSharesStore()

const resource = inject<SpaceResource>('resource')
const filterTerm = ref('')
const markInstance = ref(null)
const membersListRef = ref(null)

const filterMembers = (collection: SpaceMember[], term: string) => {
  if (!(term || '').trim()) {
    return collection
  }

  const searchEngine = new Fuse(collection, {
    ...defaultFuseOptions,
    keys: ['grantedTo.user.displayName', 'grantedTo.group.displayName']
  })
  return searchEngine.search(term).map((r) => r.item)
}

const spaceMembers = computed<SpaceMember[]>(() => {
  return Object.values(unref(resource).members)
})

const filteredSpaceMembers = computed<SpaceMember[]>(() => {
  return filterMembers(unref(spaceMembers), unref(filterTerm))
})

const availableRoles = computed<ShareRole[]>(() => {
  const permissionsWithRole = unref(spaceMembers).filter((p) => !!p.roleId)
  const roleIds = [...new Set(permissionsWithRole.map((p) => p.roleId))]
  return roleIds
    .map((r) => sharesStore.graphRoles[r])
    .filter(Boolean)
    .sort((a, b) => {
      // sort roles by amount of permissions (most likely translates to manager > editor > viewer)
      const permissionsA = a.rolePermissions.flatMap((r) => r.allowedResourceActions)
      const permissionsB = b.rolePermissions.flatMap((r) => r.allowedResourceActions)
      return permissionsB.length - permissionsA.length
    })
})

const membersWithoutRole = computed<SpaceMember[]>(() => {
  return unref(filteredSpaceMembers).filter(({ roleId }) => !roleId)
})

const getMembersForRole = (role: ShareRole): SpaceMember[] => {
  return unref(filteredSpaceMembers).filter(({ roleId }) => roleId === role.id)
}

watch(filterTerm, () => {
  if (unref(membersListRef)) {
    markInstance.value = new Mark(unref(membersListRef))
    unref(markInstance).unmark()
    unref(markInstance).mark(unref(filterTerm), {
      element: 'span',
      className: 'mark-highlight'
    })
  }
})
</script>
