<template>
  <div>
    <main class="uk-flex uk-flex-column uk-height-1-1" id="accounts-app">
      <template v-if="isInitialized">
        <h1 class="oc-invisible-sr">
          <translate>Accounts</translate>
        </h1>
        <div class="oc-app-bar">
          <accounts-batch-actions
            v-if="isAnyAccountSelected"
            :number-of-selected-accounts="numberOfSelectedAccounts"
            :selected-accounts="selectedAccounts"
          />
          <accounts-create v-else />
        </div>
        <oc-grid class="uk-flex-1 uk-overflow-auto">
          <div class="uk-width-expand">
            <accounts-list :accounts="accounts" />
          </div>
        </oc-grid>
      </template>
      <template v-else-if="hasFailed">
        <oc-alert variation="warning" no-close class="oc-m" id="accounts-list-loading-failed">
          <oc-icon name="warning" variation="warning" class="uk-float-left oc-mr-s" />
          <translate>You don't have permissions to manage accounts.</translate>
        </oc-alert>
      </template>
      <oc-loader id="accounts-list-loader" v-else />
    </main>
  </div>
</template>

<script>
import { mapGetters, mapActions, mapState } from 'vuex'
import AccountsList from './accounts/AccountsList.vue'
import AccountsCreate from './accounts/AccountsCreate.vue'
import AccountsBatchActions from './accounts/AccountsBatchActions.vue'

export default {
  name: 'App',
  components: { AccountsBatchActions, AccountsList, AccountsCreate },
  computed: {
    ...mapGetters('Accounts', ['isInitialized', 'hasFailed', 'getAccountsSorted', 'isAnyAccountSelected']),
    ...mapState('Accounts', ['selectedAccounts']),

    accounts () {
      return this.getAccountsSorted
    },
    numberOfSelectedAccounts () {
      return this.selectedAccounts.length
    }
  },
  methods: {
    ...mapActions('Accounts', ['initialize'])
  },
  created () {
    this.initialize()
  }
}
</script>
