<template>
  <div>
    <oc-table middle divider>
      <oc-table-group>
        <oc-table-row>
          <oc-table-cell shrink type="head">
            <oc-checkbox
              :value="areAllAccountsSelected"
              :label="$gettext('Select all users')"
              hide-label
              @change="toggleSelectionAll"
            />
          </oc-table-cell>
          <oc-table-cell shrink type="head" />
          <oc-table-cell type="head" v-text="$gettext('Username')" />
          <oc-table-cell type="head" v-text="$gettext('Display name')" />
          <oc-table-cell type="head" v-text="$gettext('Email')" />
          <oc-table-cell type="head" v-text="$gettext('Role')" />
          <oc-table-cell shrink type="head" class="uk-text-nowrap" v-text="$gettext('Uid number')" />
          <oc-table-cell shrink type="head" class="uk-text-nowrap" v-text="$gettext('Gid number')" />
          <oc-table-cell shrink type="head" v-text="$gettext('Enabled')" />
        </oc-table-row>
      </oc-table-group>
      <oc-table-group>
        <accounts-list-new-account-row v-if="isCreateNewRowDisplayed" @cancel="emitCreationCancel" />
        <accounts-list-row
          v-for="account in accounts"
          :key="`account-list-row-${account.id}`"
          :account="account"
        />
      </oc-table-group>
    </oc-table>
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex'
import AccountsListRow from './AccountsListRow.vue'
import AccountsListNewAccountRow from './AccountsListNewAccountRow.vue'

export default {
  name: 'AccountsList',
  components: {
    AccountsListRow,
    AccountsListNewAccountRow
  },
  props: {
    accounts: {
      type: Array,
      required: true
    },
    isCreateNewRowDisplayed: {
      type: Boolean,
      required: false,
      default: false
    }
  },
  computed: {
    ...mapGetters('Accounts', ['areAllAccountsSelected'])
  },
  methods: {
    ...mapActions('Accounts', ['toggleSelectionAll']),

    emitCreationCancel () {
      this.$emit('cancelAccountCreation')
    }
  }
}
</script>
