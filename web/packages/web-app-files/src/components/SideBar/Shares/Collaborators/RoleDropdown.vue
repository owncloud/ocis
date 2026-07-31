<template>
  <span v-if="selectedRole" class="oc-flex oc-flex-middle">
    <span v-if="availableRoles.length === 1">
      <oc-icon v-if="showIcon" :name="selectedRole.icon" class="oc-mr-s" />
      <span v-text="inviteLabel" />
    </span>
    <div v-else v-oc-tooltip="dropButtonTooltip">
      <oc-button
        :id="roleButtonId"
        class="files-recipient-role-select-btn"
        appearance="raw"
        gap-size="none"
        :disabled="isLocked"
        :aria-labelledby="roleLabelId"
      >
        <oc-icon v-if="showIcon" :name="selectedRole.icon" class="oc-mr-s" />
        <span :id="roleLabelId" class="oc-text-truncate" v-text="inviteLabel"></span>
        <oc-icon name="arrow-down-s" />
      </oc-button>
      <oc-contextual-helper
        v-if="isDisabledRole"
        class="oc-ml-xs files-permission-actions-list"
        :text="customPermissionsText"
        :title="$gettext('Custom permissions')"
      />
    </div>
    <oc-drop
      v-if="availableRoles.length > 1"
      ref="rolesDrop"
      :toggle="'#' + roleButtonId"
      mode="click"
      padding-size="small"
      class="files-recipient-role-drop"
      offset="0"
      close-on-click
    >
      <oc-list
        class="files-recipient-role-drop-list"
        :aria-label="$gettext('Select role for the invitation')"
      >
        <li v-for="role in availableRoles" :key="role.id">
          <oc-button
            :id="`files-recipient-role-drop-btn-${role.id}`"
            ref="roleSelect"
            justify-content="space-between"
            class="files-recipient-role-drop-btn oc-p-s"
            :class="{
              'oc-background-primary-gradient': isSelectedRole(role),
              selected: isSelectedRole(role)
            }"
            :appearance="isSelectedRole(role) ? 'raw-inverse' : 'raw'"
            :variation="isSelectedRole(role) ? 'primary' : 'passive'"
            @click="selectRole(role)"
          >
            <span class="oc-flex oc-flex-middle">
              <oc-icon :name="role.icon" class="oc-pl-s oc-pr-m" variation="inherit" />
              <role-item :role="role" />
            </span>
            <span class="oc-flex">
              <oc-icon v-if="isSelectedRole(role)" name="check" variation="inherit" />
            </span>
          </oc-button>
        </li>
      </oc-list>
    </oc-drop>
  </span>
</template>

<script lang="ts" setup>
import get from 'lodash-es/get'
import RoleItem from '../Shared/RoleItem.vue'
import { v4 as uuidV4 } from 'uuid'
import {
  onBeforeUnmount,
  onMounted,
  useTemplateRef,
  inject,
  ComponentPublicInstance,
  computed,
  ref,
  unref,
  Ref,
  watch
} from 'vue'
import { useGettext } from 'vue3-gettext'
import { ShareRole } from '@ownclouders/web-client'

interface Props {
  existingShareRole?: ShareRole
  existingSharePermissions?: string[]
  domSelector?: string
  mode?: 'create' | 'edit'
  showIcon?: boolean
  isLocked?: boolean
  isExternal?: boolean
}
interface Emits {
  (event: 'optionChange', role: ShareRole): void
}

const {
  existingShareRole = undefined,
  existingSharePermissions = [],
  domSelector = undefined,
  mode = 'create',
  showIcon = false,
  isLocked = false,
  isExternal = false
} = defineProps<Props>()

const emit = defineEmits<Emits>()
const rolesDropTemplateRef = useTemplateRef('rolesDrop')
const roleSelectTemplateRef = useTemplateRef('roleSelect')
const { $gettext } = useGettext()

const dropButtonTooltip = computed(() => {
  if (isLocked) {
    return $gettext('Resource is temporarily locked, unable to manage share')
  }

  return ''
})
const customPermissionsText = computed(() =>
  $gettext('Dear user, please replace this legacy role with one of the currently available roles')
)

const availableInternalRoles = inject<Ref<ShareRole[]>>('availableInternalShareRoles')
const availableExternalRoles = inject<Ref<ShareRole[]>>('availableExternalShareRoles')
// const resource = inject<Resource>('resource')

const availableRoles = computed(() => {
  let roles = availableInternalRoles
  if (isExternal) {
    roles = availableExternalRoles
  }

  return unref(roles)
})

let initialSelectedRole: ShareRole
const hasExistingShareRole = computed(() => !!existingShareRole)
const hasExistingSharePermissions = computed(() => !!existingSharePermissions.length)
const isDisabledRole = computed(
  () => !unref(hasExistingShareRole) && unref(hasExistingSharePermissions)
)
switch (true) {
  // if no role is set and no permissions are set, we use the first available role as the default
  case !unref(hasExistingShareRole) && !unref(hasExistingSharePermissions):
    initialSelectedRole = unref(availableRoles)[0]
    break
  // in the rare case that a role is disabled and permissions are set aka a disabled unified role ...
  case unref(isDisabledRole):
    // ... we need to create a fake role as an indicator that the permissions are custom
    initialSelectedRole = {
      displayName: $gettext('Custom permissions')
    }
    break
  default:
    initialSelectedRole = existingShareRole
    break
}

const selectedRole = ref<ShareRole>(initialSelectedRole)
const isSelectedRole = (role: ShareRole) => {
  return unref(selectedRole).id === role.id
}

const selectRole = (role: ShareRole) => {
  selectedRole.value = role
  emit('optionChange', unref(selectedRole))
}
const roleButtonId = computed(() => {
  if (domSelector) {
    return `files-collaborators-role-button-${domSelector}-${uuidV4()}`
  }
  return 'files-collaborators-role-button-new'
})

const inviteLabel = computed(() => {
  return $gettext(unref(selectedRole)?.displayName || '')
})

const roleLabelId = computed(() => `${unref(roleButtonId)}-label`)

watch(
  () => isExternal,
  () => {
    if (!unref(hasExistingShareRole)) {
      // when no role exists and the external flag changes, we need to reset the selected role
      selectedRole.value = unref(availableRoles)[0]
    }
  }
)

function cycleRoles(event: KeyboardEvent) {
  // events only need to be captured if the roleSelect element is visible
  if (!get(rolesDropTemplateRef.value, 'tippy.state.isShown', false)) {
    return
  }

  const { code } = event
  const isArrowUp = code === 'ArrowUp'
  const isArrowDown = code === 'ArrowDown'

  // to cycle through the list of roles only up and down keyboard events are allowed
  // if this is not the case we can return early and stop the script execution from here
  if (!isArrowUp && !isArrowDown) {
    return
  }

  // if there is only 1 or no roleSelect we can early return
  // it does not make sense to cycle through it if value is less than 1
  const roleSelect = (roleSelectTemplateRef.value as ComponentPublicInstance[]) || []

  if (roleSelect.length <= 1) {
    return
  }

  // obtain active role select in following priority chain:
  // first try to get the focused select
  // then try to get the selected select
  // and if none of those applies we fall back to the first role select
  const activeRoleSelect =
    roleSelect.find((rs) => rs.$el === document.activeElement) ||
    roleSelect.find((rs) => rs.$el.classList.contains('selected')) ||
    roleSelect[0]
  const activeRoleSelectIndex = roleSelect.indexOf(activeRoleSelect)
  const activateRoleSelect = (idx: number) => roleSelect[idx].$el.focus()

  // if the event key is arrow up
  // and the next active role select index would be less than 0
  // then activate the last available role select
  if (isArrowUp && activeRoleSelectIndex - 1 < 0) {
    activateRoleSelect(roleSelect.length - 1)

    return
  }

  // if the event key is arrow down
  // and the next active role select index would be greater or even to the available amount of role selects
  // then activate the first available role select
  if (isArrowDown && activeRoleSelectIndex + 1 >= roleSelect.length) {
    activateRoleSelect(0)

    return
  }

  // the only missing part is to navigate up or down, this only happens if:
  // the next active role index is greater than 0
  // the next active role index is less than the amount of all available role selects
  activateRoleSelect(activeRoleSelectIndex + (isArrowUp ? -1 : 1))
}
onMounted(() => {
  window.addEventListener('keydown', cycleRoles)
})
onBeforeUnmount(() => {
  window.removeEventListener('keydown', cycleRoles)
})
</script>

<style lang="scss" scoped>
.files-recipient {
  &-role-drop {
    @media (max-width: $oc-breakpoint-medium-default) {
      width: 100%;
    }
    @media (min-width: $oc-breakpoint-medium-default) {
      width: 400px;
    }

    &-list {
      li {
        margin: var(--oc-space-xsmall) 0;

        &:first-child {
          margin-top: 0;
        }

        &:last-child {
          margin-bottom: 0;
        }
      }
    }

    &-btn {
      width: 100%;
      gap: var(--oc-space-medium);

      &:hover,
      &:focus {
        background-color: var(--oc-color-background-hover);
        text-decoration: none;
      }
    }
  }

  &-role-select-btn {
    max-width: 100%;
  }
}
</style>
