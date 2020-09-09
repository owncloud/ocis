<template>
  <div>
    <div class="uk-container uk-padding">
      <h1 v-text="$gettext('Accounts')" />
      <oc-button
        v-if="numberOfSelectedAccounts < 1"
        id="accounts-new-account-trigger"
        key="create-accounts-button"
        v-text="$gettext('Create new user')"
        variation="primary"
        :disabled="isAccountCreationInProgress || !isInitialized"
        :uk-tooltip="disabledCreateAccountBtnTooltip"
        @click="setAccountCreationProgress(true)"
      />
      <oc-grid v-else key="selected-accounts-info" gutter="small" class="uk-flex-middle">
        <span v-text="selectionInfoText" />
        <span>|</span>
        <div>
          <oc-button v-text="$gettext('Clear selection')" variation="raw" @click="RESET_ACCOUNTS_SELECTION" />
        </div>
        <div>
          <oc-action-drop class="accounts-actions-dropdown">
            <template v-slot:button>
              <span class="uk-margin-xsmall-right" v-text="$gettext('Actions')" />
              <oc-icon name="expand_more" />
            </template>
            <template v-slot:actions>
                <oc-button
                  v-for="(action, index) in actions"
                  :key="action.label"
                  :id="action.id"
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
        <accounts-list
          :accounts="accounts"
          :is-create-new-row-displayed="isAccountCreationInProgress"
          @cancelAccountCreation="setAccountCreationProgress(false)"
        />
      </template>
      <oc-loader v-else />
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions, mapState, mapMutations } from 'vuex'
import AccountsList from './accounts/AccountsList.vue'
export default {
  name: 'App',
  components: { AccountsList },
  data: () => ({
    isAccountCreationInProgress: false
  }),
  computed: {
    ...mapGetters('Accounts', ['isInitialized', 'getAccountsSorted']),
    ...mapState('Accounts', ['selectedAccounts']),

    accounts () {
      return this.getAccountsSorted
    },

    numberOfSelectedAccounts () {
      return this.selectedAccounts.length
    },

    selectionInfoText () {
      const translated = this.$ngettext('%{ amount } selected user', '%{ amount } selected users', this.numberOfSelectedAccounts)

      return this.$gettextInterpolate(translated, { amount: this.numberOfSelectedAccounts })
    },

    actions () {
      let actions = []
      const permanentActions = [
        {
          id: 'accounts-actions-dropdown-action-delete',
          label: this.$gettext('Delete'),
          handler: this.deleteAccounts
        }
      ]
      const numberOfDisabledAccounts = this.selectedAccounts.filter(account => !account.accountEnabled).length
      const isAnyAccountDisabled = numberOfDisabledAccounts > 0
      const isAnyAccountEnabled = numberOfDisabledAccounts < this.numberOfSelectedAccounts

      if (isAnyAccountDisabled) {
        actions.push({
          id: 'accounts-actions-dropdown-action-enable',
          label: this.$gettext('Enable'),
          handler: () => this.toggleAccountStatus(true)
        })
      }

      if (isAnyAccountEnabled) {
        actions.push({
          id: 'accounts-actions-dropdown-action-disable',
          label: this.$gettext('Disable'),
          handler: () => this.toggleAccountStatus(false)
        })
      }

      actions = actions.concat(permanentActions)

      return actions
    },

    disabledCreateAccountBtnTooltip () {
      if (!this.isInitialized) {
        return this.$gettext('Loading users')
      }

      if (this.isAccountCreationInProgress) {
        return this.$gettext('User creation is already in progress')
      }

      return null
    }
  },
  methods: {
    ...mapActions('Accounts', ['initialize', 'toggleAccountStatus', 'deleteAccounts']),
    ...mapMutations('Accounts', ['RESET_ACCOUNTS_SELECTION']),

    setAccountCreationProgress (isInProgress) {
      this.isAccountCreationInProgress = isInProgress
    }
  },
  created () {
    this.initialize()
  },
  beforeDestroy () {
    this.RESET_ACCOUNTS_SELECTION()
    this.setAccountCreationProgress(false)
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
