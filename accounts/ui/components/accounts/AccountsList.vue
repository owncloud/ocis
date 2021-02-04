<template>
  <div>
    <oc-table-simple>
      <oc-thead>
        <oc-tr class="fix-table-header">
          <oc-th shrink type="head" alignH="center">
            <oc-checkbox
              class="oc-ml-s"
              :value="areAllAccountsSelected"
              @input="toggleSelectionAll"
              :label="$gettext('Select all users')"
              hide-label
            />
          </oc-th>
          <oc-th shrink type="head" />
          <oc-th type="head" v-text="$gettext('Username')" />
          <oc-th type="head" v-text="$gettext('Display name')" />
          <oc-th type="head" v-text="$gettext('Email')" />
          <oc-th type="head" v-text="$gettext('Role')" />
          <oc-th shrink type="head" v-text="$gettext('Activated')" alignH="center"/>
        </oc-tr>
      </oc-thead>
      <oc-tbody>
        <accounts-list-row
          v-for="account in accounts"
          :key="`account-list-row-${account.id}`"
          :account="account"
        />
      </oc-tbody>
    </oc-table-simple>
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
