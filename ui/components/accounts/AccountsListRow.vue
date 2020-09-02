<template>
  <oc-table-row>
    <oc-table-cell>
      <oc-checkbox :value="isAccountSelected" @change="TOGGLE_SELECTION_ACCOUNT(account)" />
    </oc-table-cell>
    <oc-table-cell>
      <avatar :user-name="account.displayName || account.onPremisesSamAccountName" :userid="account.id" :width="35" />
    </oc-table-cell>
    <oc-table-cell v-text="account.onPremisesSamAccountName" />
    <oc-table-cell v-text="account.displayName || '-'" />
    <oc-table-cell v-text="account.mail" />
    <oc-table-cell>
      <oc-button :id="`accounts-roles-select-trigger-${account.id}`" class="accounts-roles-select-trigger" variation="raw">
        <span class="uk-flex uk-flex-middle accounts-roles-current-role">
          {{ currentRole ? currentRole.displayName : $gettext('Select role') }}
          <oc-icon name="expand_more" aria-hidden="true" />
        </span>
      </oc-button>
      <oc-drop
        :drop-id="`accounts-roles-select-dropdown-${account.id}`"
        :toggle="`#accounts-roles-select-trigger-${account.id}`"
        mode="click"
        close-on-click
        :options="{ delayHide: 0 }"
      >
        <ul class="uk-list">
          <li v-for="role in roles" :key="role.id">
            <label class="accounts-roles-dropdown-role">
              <input
                type="radio"
                class="oc-radiobutton"
                v-model="currentRole"
                :value="role"
                @change="changeRole(role.id)"
              />
              {{ role.displayName }}
            </label>
          </li>
        </ul>
      </oc-drop>
    </oc-table-cell>
    <oc-table-cell v-text="account.uidNumber || '-'" />
    <oc-table-cell v-text="account.gidNumber || '-'" />
    <oc-table-cell class="uk-text-center">
      <oc-icon v-if="account.accountEnabled" name="ready" variation="success" :aria-label="$gettext('Account is enabled')" />
      <oc-icon v-else name="deprecated" variation="danger" :aria-label="$gettext('Account is disabled')" />
    </oc-table-cell>
  </oc-table-row>
</template>

<script>
import { mapGetters, mapState, mapActions, mapMutations } from 'vuex'
import { isObjectEmpty } from '../../helpers/utils'
import { injectAuthToken } from '../../helpers/auth'
// eslint-disable-next-line camelcase
import { RoleService_AssignRoleToUser, RoleService_ListRoleAssignments } from '../../client/settings'
import Avatar from './Avatar.vue'

export default {
  name: 'AccountsListRow',

  components: { Avatar },

  props: {
    account: {
      type: Object,
      required: true
    }
  },

  data () {
    return {
      currentRole: null
    }
  },

  computed: {
    ...mapGetters(['user', 'configuration']),
    ...mapState('Accounts', ['roles', 'selectedAccounts']),

    isAccountSelected () {
      return this.selectedAccounts.indexOf(this.account) > -1
    }
  },

  created () {
    this.getUsersCurrentRole()
  },

  methods: {
    ...mapActions(['showMessage']),
    ...mapMutations('Accounts', ['TOGGLE_SELECTION_ACCOUNT']),

    async changeRole (roleId) {
      injectAuthToken(this.user.token)

      const response = await RoleService_AssignRoleToUser({
        $domain: this.configuration.server,
        body: {
          account_uuid: this.account.id,
          role_id: roleId
        }
      })

      if (response.status === 201) {
        const roleId = response.data.assignment.roleId
        this.currentRole = this.roles.find(role => {
          return role.id === roleId
        })
      } else {
        this.showMessage({
          title: this.$gettext('Failed to change role.'),
          desc: response.statusText,
          status: 'danger'
        })
      }
    },

    async getUsersCurrentRole () {
      injectAuthToken(this.user.token)

      const response = await RoleService_ListRoleAssignments({
        $domain: this.configuration.server,
        body: {
          account_uuid: this.account.id
        }
      })

      if (response.status === 201) {
        const assignedRole = response.data

        if (isObjectEmpty(assignedRole)) {
          return
        }

        this.currentRole = this.roles.find(role => {
          return role.id === assignedRole.assignments[0].roleId
        })
      }
    }
  }
}
</script>
