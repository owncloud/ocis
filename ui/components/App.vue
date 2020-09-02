<template>
  <div>
    <div class="uk-container uk-padding">
      <h1 v-text="$gettext('Accounts')" />
      <oc-grid v-if="selectedAccountsAmount > 0" key="selected-accounts-info" gutter="small" class="uk-flex-middle">
        <span v-text="selectionInfoText" />
        <div>
          <oc-action-drop>
            <template v-slot:button>
              <span class="uk-margin-xsmall-right" v-text="$gettext('Actions')" />
              <oc-icon name="expand_more" />
            </template>
            <template v-slot:actions>
                <oc-button
                  v-for="(action, index) in actions"
                  :key="action.label"
                  variation="raw"
                  role="menuitem"
                  :class="{ 'uk-margin-small-bottom': index + 1 !== actions.length }"
                  class="uk-width-1-1 uk-flex-left"
                  @click="action.handler"
                >
                  {{ action.label }}
                </oc-button>
            </template>
          </oc-action-drop>
        </div>
      </oc-grid>
      <template v-if="isInitialized">
        <accounts-list :accounts="accounts" />
      </template>
      <oc-loader v-else />
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions, mapState } from 'vuex'
import AccountsList from './accounts/AccountsList.vue'
export default {
  name: 'App',
  components: { AccountsList },
  computed: {
    ...mapGetters('Accounts', ['isInitialized', 'getAccountsSorted']),
    ...mapState('Accounts', ['selectedAccounts']),

    accounts () {
      return this.getAccountsSorted
    },

    selectedAccountsAmount () {
      return this.selectedAccounts.length
    },

    selectionInfoText () {
      const translated = this.$ngettext('%{ amount } selected user', '%{ amount } selected users', this.selectedAccountsAmount)

      return this.$gettextInterpolate(translated, { amount: this.selectedAccountsAmount })
    },

    actions () {
      const actions = []
      const isAnyAccountDisabled = this.selectedAccounts.some(account => !account.accountEnabled)
      const isAnyAccountEnabled = this.selectedAccounts.some(account => account.accountEnabled)

      if (isAnyAccountDisabled) {
        actions.push({
          label: this.$gettext('Enable'),
          handler: this.enableAccounts
        })
      }

      if (isAnyAccountEnabled) {
        actions.push({
          label: this.$gettext('Disable'),
          handler: this.disableAccounts
        })
      }

      return actions
    }
  },
  methods: {
    ...mapActions('Accounts', ['initialize', 'enableAccounts', 'disableAccounts'])
  },
  created () {
    this.initialize()
  }
}
</script>

<style>
/* TODO: After https://github.com/owncloud/owncloud-design-system/pull/418 gets merged
there won't be an extra span and this won't be needed anymore */
.accounts-selection-actions-btn > span {
  display: flex;
  align-items: center;
}

/* TODO: Adjust in ODS */
.oc-dropdown-menu {
  width: 150px;
}
</style>
