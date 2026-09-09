<template>
  <div id="user-edit-panel" class="oc-mt-xl">
    <UserInfoBox :user="user" />
    <form id="user-edit-form" class="oc-background-highlight oc-p-m" autocomplete="off">
      <div>
        <oc-text-input
          id="userName-input"
          v-model="editUser.onPremisesSamAccountName"
          class="oc-mb-s"
          :label="$gettext('User name')"
          :error-message="formData.userName.errorMessage"
          :fix-message-line="true"
          :read-only="isInputFieldReadOnly('user.onPremisesSamAccountName')"
          @update:model-value="validateUserName"
        />
        <oc-text-input
          id="displayName-input"
          v-model="editUser.displayName"
          class="oc-mb-s"
          :label="$gettext('First and last name')"
          :error-message="formData.displayName.errorMessage"
          :fix-message-line="true"
          :read-only="isInputFieldReadOnly('user.displayName')"
          @update:model-value="validateDisplayName"
        />
        <oc-text-input
          id="email-input"
          v-model="editUser.mail"
          class="oc-mb-s"
          :label="$gettext('Email')"
          :error-message="formData.email.errorMessage"
          type="email"
          :fix-message-line="true"
          :read-only="isInputFieldReadOnly('user.mail')"
          @change="validateEmail"
        />
        <oc-text-input
          id="password-input"
          :model-value="editUser.passwordProfile?.password"
          class="oc-mb-s"
          :label="$gettext('Password')"
          type="password"
          :fix-message-line="true"
          placeholder="●●●●●●●●"
          :read-only="isInputFieldReadOnly('user.passwordProfile')"
          @update:model-value="onUpdatePassword"
        />
        <div class="oc-mb-s">
          <oc-select
            id="role-input"
            :model-value="selectedRoleValue"
            :label="$gettext('Role')"
            option-label="displayName"
            :options="translatedRoleOptions"
            :clearable="false"
            :read-only="isInputFieldReadOnly('user.appRoleAssignments')"
            @update:model-value="onUpdateRole"
          />
          <div class="oc-text-input-message"></div>
        </div>
        <div class="oc-mb-s">
          <oc-select
            id="login-input"
            :disabled="isLoginInputDisabled"
            :model-value="selectedLoginValue"
            :label="$gettext('Login')"
            :options="loginOptions"
            :clearable="false"
            :read-only="isInputFieldReadOnly('user.accountEnabled')"
            @update:model-value="onUpdateLogin"
          />

          <div class="oc-text-input-message"></div>
        </div>
        <quota-select
          id="quota-select-form"
          :key="'quota-select-' + user.id"
          :disabled="isQuotaInputDisabled"
          class="oc-mb-s"
          :label="$gettext('Personal quota')"
          :total-quota="editUser.drive?.quota?.total || 0"
          :max-quota="maxQuota"
          :fix-message-line="true"
          :description-message="
            isQuotaInputDisabled && !isInputFieldReadOnly('drive.quota')
              ? $gettext('To set an individual quota, the user needs to have logged in once.')
              : ''
          "
          :read-only="isInputFieldReadOnly('drive.quota')"
          @selected-option-change="changeSelectedQuotaOption"
        />
        <group-select
          class="oc-mb-s"
          :read-only="isInputFieldReadOnly('user.memberOf')"
          :selected-groups="editUser.memberOf"
          :group-options="groupOptions"
          @selected-option-change="changeSelectedGroupOption"
        />
      </div>
      <compare-save-dialog
        class="edit-compare-save-dialog oc-mb-l"
        :original-object="compareSaveDialogOriginalObject"
        :compare-object="editUser"
        :confirm-button-disabled="invalidFormData"
        @revert="revertChanges"
        @confirm="onEditUser({ user, editUser })"
      ></compare-save-dialog>
    </form>
  </div>
</template>
<script lang="ts" setup>
import { computed, watch, ref, unref, isRef } from 'vue'
import * as EmailValidator from 'email-validator'
import UserInfoBox from './UserInfoBox.vue'
import {
  CompareSaveDialog,
  QuotaSelect,
  useUserStore,
  useCapabilityStore,
  useEventBus,
  useMessages,
  useSpacesStore,
  useSharesStore
} from '@ownclouders/web-pkg'
import GroupSelect from '../GroupSelect.vue'
import { cloneDeep, isEmpty, isEqual, omit } from 'lodash-es'
import { AppRole, AppRoleAssignment, Group, User } from '@ownclouders/web-client/graph/generated'
import { MaybeRef, useClientService } from '@ownclouders/web-pkg'
import { storeToRefs } from 'pinia'
import { diff } from 'deep-object-diff'
import { useUserSettingsStore } from '../../../composables/stores/userSettings'
import { useGettext } from 'vue3-gettext'

interface Props {
  user?: User | null
  roles: AppRole[]
  groups: Group[]
  applicationId: string
}
const { user = null, roles, groups, applicationId } = defineProps<Props>()
const capabilityStore = useCapabilityStore()
const capabilityRefs = storeToRefs(capabilityStore)
const maxQuota = capabilityRefs.spacesMaxQuota
const clientService = useClientService()
const userStore = useUserStore()
const userSettingsStore = useUserSettingsStore()
const spacesStore = useSpacesStore()
const sharesStore = useSharesStore()
const eventBus = useEventBus()
const { showErrorMessage } = useMessages()
const { $gettext } = useGettext()

let editUser: MaybeRef<User> = ref()
const formData = ref({
  displayName: {
    errorMessage: '',
    valid: true
  },
  userName: {
    errorMessage: '',
    valid: true
  },
  email: {
    errorMessage: '',
    valid: true
  }
})
function changeSelectedQuotaOption(option: { value: number; displayValue: string }) {
  unref(editUser).drive.quota.total = option.value
}
function changeSelectedGroupOption(option: Group[]) {
  unref(editUser).memberOf = option
}
async function validateUserName() {
  unref(formData).userName.valid = false

  if (unref(editUser).onPremisesSamAccountName.trim() === '') {
    unref(formData).userName.errorMessage = $gettext('User name cannot be empty')
    return false
  }

  if (unref(editUser).onPremisesSamAccountName.includes(' ')) {
    unref(formData).userName.errorMessage = $gettext('User name cannot contain white spaces')
    return false
  }

  if (
    unref(editUser).onPremisesSamAccountName.length &&
    !isNaN(parseInt(unref(editUser).onPremisesSamAccountName[0]))
  ) {
    unref(formData).userName.errorMessage = $gettext('User name cannot start with a number')
    return false
  }

  if (unref(editUser).onPremisesSamAccountName.length > 255) {
    unref(formData).userName.errorMessage = $gettext('User name cannot exceed 255 characters')
    return false
  }

  if (unref(user).onPremisesSamAccountName !== unref(editUser).onPremisesSamAccountName) {
    try {
      // Validate username by fetching the user. If the request succeeds, we throw a validation error
      const client = clientService.graphAuthenticated
      await client.users.getUser(unref(editUser).onPremisesSamAccountName)
      unref(formData).userName.errorMessage = $gettext('User "%{userName}" already exists', {
        userName: unref(editUser).onPremisesSamAccountName
      })
      return false
    } catch {}
  }

  unref(formData).userName.errorMessage = ''
  unref(formData).userName.valid = true
  return true
}
function validateDisplayName() {
  unref(formData).displayName.valid = false

  if (unref(editUser).displayName.trim() === '') {
    unref(formData).displayName.errorMessage = $gettext('First and last name cannot be empty')
    return false
  }

  if (unref(editUser).displayName.length > 255) {
    unref(formData).displayName.errorMessage = $gettext(
      'First and last name cannot exceed 255 characters'
    )
    return false
  }

  unref(formData).displayName.errorMessage = ''
  unref(formData).displayName.valid = true
  return true
}
function validateEmail() {
  unref(formData).email.valid = false

  if (!EmailValidator.validate(unref(editUser).mail)) {
    unref(formData).email.errorMessage = $gettext('Please enter a valid email')
    return false
  }

  unref(formData).email.errorMessage = ''
  unref(formData).email.valid = true
  return true
}
function revertChanges() {
  if (isRef(editUser)) {
    editUser.value = cloneDeep(unref(user))
  } else {
    editUser = ref(cloneDeep(unref(user)))
  }
  Object.values(unref(formData)).forEach((formDataValue: any) => {
    formDataValue.valid = true
    formDataValue.errorMessage = ''
  })
}
function onUpdateRole(role: AppRoleAssignment) {
  if (!unref(editUser).appRoleAssignments.length) {
    // FIXME: Add resourceId and principalId to be able to remove type cast
    unref(editUser).appRoleAssignments.push({
      appRoleId: role.id
    } as AppRoleAssignment)
    return
  }
  unref(editUser).appRoleAssignments[0].appRoleId = role.id
}
function onUpdatePassword(password: string) {
  unref(editUser).passwordProfile = {
    password
  }
}
function onUpdateLogin({ value }: { value: boolean }) {
  unref(editUser).accountEnabled = value
}
const groupOptions = computed(() => {
  const { memberOf: selectedGroups } = unref(editUser)
  return groups.filter(
    (g) => !selectedGroups.some((s) => s.id === g.id) && !g.groupTypes?.includes('ReadOnly')
  )
})
const isLoginInputDisabled = computed(() => userStore.user.id === user.id)
const isInputFieldReadOnly = (key: string) => {
  return capabilityStore.graphUsersReadOnlyAttributes.includes(key)
}

const onUpdateUserAppRoleAssignments = (user: User, editUser: User) => {
  const client = clientService.graphAuthenticated
  return client.users.createUserAppRoleAssignment(user.id, {
    appRoleId: editUser.appRoleAssignments[0].appRoleId,
    resourceId: applicationId,
    principalId: editUser.id
  })
}
const onUpdateUserGroupAssignments = (user: User, editUser: User) => {
  const client = clientService.graphAuthenticated
  const groupsToAdd = editUser.memberOf.filter(
    (editUserGroup) => !user.memberOf.some((g) => g.id === editUserGroup.id)
  )
  const groupsToDelete = user.memberOf.filter(
    (editUserGroup) => !editUser.memberOf.some((g) => g.id === editUserGroup.id)
  )
  const requests = []

  for (const groupToAdd of groupsToAdd) {
    requests.push(client.groups.addMember(groupToAdd.id, user.id))
  }
  for (const groupToDelete of groupsToDelete) {
    requests.push(client.groups.deleteMember(groupToDelete.id, user.id))
  }

  return Promise.all(requests)
}

const onUpdateUserDrive = async (editUser: User) => {
  const client = clientService.graphAuthenticated
  const updateSpace = await client.drives.updateDrive(
    editUser.drive.id,
    {
      name: editUser.drive.name,
      quota: { total: editUser.drive.quota.total }
    },
    sharesStore.graphRoles
  )

  if (editUser.id === userStore.user.id) {
    // Load current user quota
    spacesStore.updateSpaceField({
      id: editUser.drive.id,
      field: 'spaceQuota',
      value: updateSpace.spaceQuota
    })
  }
}

const onEditUser = async ({ user, editUser }: { user: User; editUser: User }) => {
  try {
    const client = clientService.graphAuthenticated
    const graphEditUserPayloadExtractor = (user: User) => {
      return omit(user, ['drive', 'appRoleAssignments', 'memberOf'])
    }
    const graphEditUserPayload = diff(
      graphEditUserPayloadExtractor(user),
      graphEditUserPayloadExtractor(editUser)
    ) as User

    if (!isEmpty(graphEditUserPayload)) {
      await client.users.editUser(editUser.id, graphEditUserPayload)
    }

    if (!isEqual(user.drive?.quota?.total, editUser.drive?.quota?.total)) {
      await onUpdateUserDrive(editUser)
    }

    if (!isEqual(user.memberOf, editUser.memberOf)) {
      await onUpdateUserGroupAssignments(user, editUser)
    }

    if (
      !isEqual(user.appRoleAssignments[0]?.appRoleId, editUser.appRoleAssignments[0]?.appRoleId)
    ) {
      await onUpdateUserAppRoleAssignments(user, editUser)
    }

    const updatedUser = await client.users.getUser(user.id)
    userSettingsStore.upsertUser(updatedUser)

    eventBus.publish('sidebar.entity.saved')

    if (userStore.user.id === updatedUser.id) {
      userStore.setUser(updatedUser)
    }

    return updatedUser
  } catch (error) {
    console.error(error)
    showErrorMessage({
      title: $gettext('Failed to edit user'),
      errors: [error]
    })
  }
}
const loginOptions = computed(() => {
  return [
    {
      label: $gettext('Allowed'),
      value: true
    },
    {
      label: $gettext('Forbidden'),
      value: false
    }
  ]
})
const selectedLoginValue = computed(() => {
  return unref(loginOptions).find((option) =>
    !('accountEnabled' in unref(editUser))
      ? option.value === true
      : unref(editUser).accountEnabled === option.value
  )
})
const translatedRoleOptions = computed(() => {
  return roles.map((role) => {
    return { ...role, displayName: $gettext(role.displayName) }
  })
})
const selectedRoleValue = computed(() => {
  const assignedRole = unref(editUser)?.appRoleAssignments?.[0]
  return unref(translatedRoleOptions).find((role) => role.id === assignedRole?.appRoleId)
})
const invalidFormData = computed(() => {
  return Object.values(unref(formData))
    .map((v: any) => !!v.valid)
    .includes(false)
})
const showQuota = computed(() => {
  return unref(editUser).drive?.quota
})
const isQuotaInputDisabled = computed(() => {
  return typeof unref(showQuota) === 'undefined'
})
const compareSaveDialogOriginalObject = computed(() => {
  return cloneDeep(user)
})

watch(
  () => user,
  () => {
    if (isRef(editUser)) {
      editUser.value = cloneDeep(user)
    } else {
      editUser = ref(cloneDeep(user))
    }
  },
  {
    deep: true,
    immediate: true
  }
)

watch(
  editUser,
  () => {
    /**
     * Property accountEnabled won't be always set, but this still means, that login is allowed.
     * So we actually don't need to change the property if missing and not set to forbidden in the UI.
     * This also avoids the compare save dialog from displaying that there are unsaved changes.
     */
    if (unref(editUser).accountEnabled === true && !('accountEnabled' in user)) {
      delete unref(editUser).accountEnabled
    }
  },
  {
    deep: true
  }
)
</script>

<style lang="scss">
#user-edit-panel {
  #user-edit-form {
    border-radius: 5px;
  }
}
</style>
