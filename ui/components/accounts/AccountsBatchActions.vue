<template>
  <oc-grid key="selected-accounts-info" gutter="small" class="uk-flex-middle">
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
</template>

<script>
import { mapActions, mapMutations } from 'vuex'

export default {
  name: 'AccountsBatchActions',
  props: {
    numberOfSelectedAccounts: {
      type: Number,
      required: true
    },
    selectedAccounts: {
      type: Array,
      required: true
    }
  },
  computed: {
    selectionInfoText () {
      const translated = this.$ngettext('%{ amount } selected user', '%{ amount } selected users', this.numberOfSelectedAccounts)
      return this.$gettextInterpolate(translated, { amount: this.numberOfSelectedAccounts })
    },
    actions () {
      const actions = []
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

      actions.push({
        id: 'accounts-actions-dropdown-action-delete',
        label: this.$gettext('Delete'),
        handler: this.deleteAccounts
      })

      return actions
    }
  },
  methods: {
    ...mapActions('Accounts', ['toggleAccountStatus', 'deleteAccounts']),
    ...mapMutations('Accounts', ['RESET_ACCOUNTS_SELECTION'])
  }
}
</script>
