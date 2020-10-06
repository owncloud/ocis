<template>
  <div>
    <oc-table middle divider>
      <oc-table-group>
        <oc-table-row class="fix-table-header">
          <oc-table-cell shrink type="head" class="uk-text-center">
            <oc-checkbox
                class="oc-ml-s"
                :value="areAllAccountsSelected"
                @input="toggleSelectionAll"
                :label="$gettext('Select all users')"
                hide-label
            />
          </oc-table-cell>
          <oc-table-cell shrink type="head" />
          <oc-table-cell type="head" v-text="$gettext('Username')" />
          <oc-table-cell type="head" v-text="$gettext('Display name')" />
          <oc-table-cell type="head" v-text="$gettext('Email')" />
          <oc-table-cell type="head" v-text="$gettext('Role')" />
          <oc-table-cell shrink type="head" v-text="$gettext('Activated')" />
        </oc-table-row>
      </oc-table-group>
      <oc-table-group>
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
import { mapActions, mapGetters, mapMutations } from 'vuex'
import AccountsListRow from './AccountsListRow.vue'

export default {
  name: 'AccountsList',
  components: {
    AccountsListRow
  },
  props: {
    accounts: {
      type: Array,
      required: true
    }
  },
  computed: {
    ...mapGetters('Accounts', ['areAllAccountsSelected'])
  },
  methods: {
    ...mapActions('Accounts', ['toggleSelectionAll']),
    ...mapMutations('Accounts', ['RESET_ACCOUNTS_SELECTION'])
  },
  beforeDestroy () {
    this.RESET_ACCOUNTS_SELECTION()
  }
}
</script>

<style scoped>
.fix-table-header > th {
  text-transform: none;
}
</style>
