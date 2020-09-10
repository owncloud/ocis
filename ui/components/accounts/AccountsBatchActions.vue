<template>
  <oc-grid key="selected-accounts-info" gutter="small" class="uk-flex-middle">
    <span v-text="selectionInfoText" />
    <span>|</span>
    <div>
      <oc-button v-text="$gettext('Clear selection')" variation="raw" @click="RESET_ACCOUNTS_SELECTION" />
    </div>
    <oc-grid gutter="small" id="accounts-batch-actions">
      <div v-for="action in actions" :key="action.label">
        <oc-button :id="action.id" @click="action.handler" :variation="action.variation || 'default'" :icon="action.icon">
          {{ action.label }}
        </oc-button>
      </div>
    </oc-grid>
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
          label: this.$gettext('Activate'),
          icon: 'ready',
          handler: () => this.setAccountActivated(true)
        })
      }

      if (isAnyAccountEnabled) {
        actions.push({
          id: 'accounts-actions-dropdown-action-disable',
          label: this.$gettext('Block'),
          icon: 'deprecated',
          handler: () => this.setAccountActivated(false)
        })
      }

      actions.push({
        id: 'accounts-actions-dropdown-action-delete',
        label: this.$gettext('Delete'),
        icon: 'delete',
        handler: this.deleteAccounts
      })

      return actions
    }
  },
  methods: {
    ...mapActions('Accounts', ['setAccountActivated', 'deleteAccounts']),
    ...mapMutations('Accounts', ['RESET_ACCOUNTS_SELECTION'])
  }
}
</script>
