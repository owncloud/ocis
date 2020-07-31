<template>
  <oc-table-row>
    <oc-table-cell>
      <avatar :user-name="account.displayName || account.onPremisesSamAccountName" :userid="account.id" :width="35" />
    </oc-table-cell>
    <oc-table-cell v-text="account.onPremisesSamAccountName" />
    <oc-table-cell v-text="account.displayName || '-'" />
    <oc-table-cell v-text="account.mail" />
    <oc-table-cell>
      <oc-button :class="`accounts-roles-select-trigger-${account.id}`" variation="raw">
        <span class="uk-flex uk-flex-middle">
          {{ currentRole ? currentRole.displayName : $gettext('Select role') }}
          <oc-icon name="expand_more" aria-hidden="true" />
        </span>
      </oc-button>
      <oc-drop
        :drop-id="`accounts-roles-select-dropdown-${account.id}`"
        :toggle="`.accounts-roles-select-trigger-${account.id}`"
        mode="click"
        close-on-click
        :options="{ delayHide: 0 }"
      >
        <ul class="uk-nav">
          <li v-for="role in roles" :key="role.id" class="uk-margin-small">
            <oc-button variation="raw" @click="changeRole(role.id)">
              <span v-text="role.displayName" :class="{ 'uk-text-bold': role === currentRole }" />
            </oc-button>
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
import { mapGetters, mapState, mapActions } from 'vuex'
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
    ...mapState('Accounts', ['roles'])
  },

  created () {
    this.getUsersCurrentRole()
  },

  methods: {
    ...mapActions(['showMessage']),

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
        this.getUsersCurrentRole()
      } else {
        this.showMessage({
          title: 'Failed to change role.',
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
