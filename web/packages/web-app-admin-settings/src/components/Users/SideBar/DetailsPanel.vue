<template>
  <div class="oc-mt-xl">
    <div v-if="noUsers" class="oc-flex user-info oc-text-center">
      <oc-icon name="user" size="xxlarge" />
      <p v-translate data-testid="no-user-selected">Select a user to view details</p>
    </div>
    <div v-if="multipleUsers" class="oc-flex group-info">
      <oc-icon name="group" size="xxlarge" />
      <p>{{ multipleUsersSelectedText }}</p>
    </div>
    <div v-if="user">
      <UserInfoBox :user="user" />
      <table
        class="details-table"
        :aria-label="$gettext('Overview of the information about the selected user')"
      >
        <tbody>
          <tr>
            <th scope="col" class="oc-pr-s" v-text="$gettext('User name')" />
            <td v-text="user.onPremisesSamAccountName" />
          </tr>
          <tr>
            <th scope="col" class="oc-pr-s" v-text="$gettext('First and last name')" />
            <td v-text="user.displayName" />
          </tr>
          <tr>
            <th scope="col" class="oc-pr-s" v-text="$gettext('Email')" />
            <td>
              <span v-text="user.mail" />
            </td>
          </tr>
          <tr>
            <th scope="col" class="oc-pr-s" v-text="$gettext('Role')" />
            <td>
              <span v-if="user.appRoleAssignments" v-text="roleDisplayName" />
              <span v-else>
                <span class="oc-mr-xs">-</span>
                <oc-contextual-helper
                  :text="
                    $gettext(
                      'User roles become available once the user has logged in for the first time.'
                    )
                  "
                  :title="$gettext('User role')"
                />
              </span>
            </td>
          </tr>
          <tr>
            <th scope="col" class="oc-pr-s" v-text="$gettext('Login')" />
            <td>
              <span v-text="loginDisplayValue" />
            </td>
          </tr>
          <tr>
            <th scope="col" class="oc-pr-s" v-text="$gettext('Quota')" />
            <td>
              <span v-if="showUserQuota" v-text="quotaDisplayValue" />
              <span v-else>
                <span class="oc-mr-xs">-</span>
                <oc-contextual-helper
                  :text="
                    $gettext(
                      'User quota becomes available once the user has logged in for the first time.'
                    )
                  "
                  :title="$gettext('Quota')"
                />
              </span>
            </td>
          </tr>
          <tr>
            <th scope="col" class="oc-pr-s" v-text="$gettext('Groups')" />
            <td>
              <span v-if="user.memberOf.length" v-text="groupsDisplayValue" />
              <span v-else>
                <span class="oc-mr-xs">-</span>
                <oc-contextual-helper
                  :text="$gettext('No groups assigned.')"
                  :title="$gettext('Groups')"
                />
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { computed } from 'vue'
import UserInfoBox from './UserInfoBox.vue'
import { AppRole, User } from '@ownclouders/web-client/graph/generated'
import { formatFileSize } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'

interface Props {
  user?: User | null
  users: User[]
  roles: AppRole[]
}
const { user = null, users, roles } = defineProps<Props>()
const language = useGettext()
const { $gettext } = language
const currentLanguage = computed(() => {
  return language.current
})
const noUsers = computed(() => !users.length)
const multipleUsers = computed(() => users.length > 1)
const multipleUsersSelectedText = computed(() => {
  return $gettext('%{count} users selected', {
    count: users.length.toString()
  })
})
const roleDisplayName = computed(() => {
  const assignedRole = user.appRoleAssignments[0]

  return (
    $gettext(roles.find((role) => role.id === assignedRole?.appRoleId)?.displayName || '') || '-'
  )
})
const groupsDisplayValue = computed(() => {
  return user.memberOf
    .map((group) => group.displayName)
    .sort()
    .join(', ')
})
const showUserQuota = computed(() => {
  return 'total' in (user.drive?.quota || {})
})
const quotaDisplayValue = computed(() => {
  return user.drive.quota.total === 0
    ? $gettext('No restriction')
    : formatFileSize(user.drive.quota.total, currentLanguage.value)
})
const loginDisplayValue = computed(() => {
  return user.accountEnabled === false ? $gettext('Forbidden') : $gettext('Allowed')
})
</script>
<style lang="scss">
.details-table {
  text-align: left;

  tr {
    height: 1.5rem;
  }

  th {
    font-weight: 600;
  }
}
</style>
