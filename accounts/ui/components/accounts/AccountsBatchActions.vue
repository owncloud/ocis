<template>
  <oc-grid key="selected-accounts-info" gutter="small" class="uk-flex-middle">
    <span v-text="selectionInfoText" />
    <span>|</span>
    <div>
      <oc-button v-text="$gettext('Clear selection')" appearance="raw" @click="RESET_ACCOUNTS_SELECTION" />
    </div>
    <oc-grid gutter="small" id="accounts-batch-actions">
      <div v-for="action in actions" :key="action.label">
        <div v-if="isConfirmationInProgress[action.id]" :variation="action.confirmation.variation || 'primary'" noClose class="uk-flex uk-flex-middle tmp-alert-fixes">
          <span>{{ action.confirmation.message }}</span>
          <oc-button :id="action.confirmation.cancel.id" @click="action.confirmation.cancel.handler" :variation="action.confirmation.cancel.variation || 'passive'">
            {{ action.confirmation.cancel.label }}
          </oc-button>
          <oc-button :id="action.confirmation.confirm.id" @click="action.confirmation.confirm.handler" :variation="action.confirmation.confirm.variation || 'primary'">
            {{ action.confirmation.confirm.label }}
          </oc-button>
        </div>
        <oc-button v-else :id="action.id" @click="action.handler" :variation="action.variation || 'primary'" :icon="action.icon">
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
  data: () => {
    return {
      isConfirmationInProgress: {}
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
          id: 'accounts-batch-action-enable',
          label: this.$gettext('Activate'),
          icon: 'ready',
          handler: () => this.setAccountActivated(true)
        })
      }

      if (isAnyAccountEnabled) {
        actions.push({
          id: 'accounts-batch-action-disable',
          label: this.$gettext('Block'),
          icon: 'deprecated',
          handler: () => this.setAccountActivated(false)
        })
      }

      const idDeleteAction = 'accounts-batch-action-delete'
      actions.push({
        id: idDeleteAction,
        label: this.$gettext('Delete'),
        icon: 'delete',
        variation: 'danger',
        handler: () => this.showConfirmationRequest(idDeleteAction),
        confirmation: {
          variation: 'danger',
          message: this.$ngettext(
            'Delete the selected account?',
            'Delete the selected accounts?',
            this.numberOfSelectedAccounts
          ),
          cancel: {
            id: 'accounts-batch-action-delete-cancel',
            label: this.$gettext('Cancel'),
            handler: () => this.hideConfirmationRequest(idDeleteAction)
          },
          confirm: {
            id: 'accounts-batch-action-delete-confirm',
            label: this.$gettext('Confirm'),
            variation: 'danger',
            handler: this.deleteAccounts
          }
        }
      })

      return actions
    }
  },
  methods: {
    ...mapActions('Accounts', ['setAccountActivated', 'deleteAccounts']),
    ...mapMutations('Accounts', ['RESET_ACCOUNTS_SELECTION']),
    showConfirmationRequest (actionId) {
      this.isConfirmationInProgress = { ...this.isConfirmationInProgress, [actionId]: true }
    },
    hideConfirmationRequest (actionId) {
      this.isConfirmationInProgress = { ...this.isConfirmationInProgress, [actionId]: false }
    }
  }
}
</script>

<style lang="scss" scoped>
.tmp-alert-fixes {
  color: rgb(224, 0, 0) !important;

  font-size: 1.125rem !important;
  font-weight: 600 !important;
  line-height: 1.4 !important;
}
.tmp-alert-fixes > *:not(:last-child) {
  margin-right: 8px;
}
.tmp-alert-fixes > button {
  padding: 0.2rem 0.5rem;
}
</style>
