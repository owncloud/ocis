<template>
  <oc-grid key="selected-accounts-info" gutter="small" class="uk-flex-middle">
    <span v-text="selectionInfoText" />
    <span>|</span>
    <div>
      <oc-button v-text="$gettext('Clear selection')" variation="raw" @click="RESET_ACCOUNTS_SELECTION" />
    </div>
    <oc-grid gutter="small" id="accounts-batch-actions">
      <div v-for="action in actions" :key="action.label">
        <oc-alert v-if="isConfirmationInProgress[action.id]" :variation="action.confirmation.variation || 'default'" noClose class="tmp-alert-fixes">
          <span>{{ action.confirmation.message }}</span>
          <oc-button size="small" :id="action.confirmation.cancel.id" @click="action.confirmation.cancel.handler" :variation="action.confirmation.cancel.variation || 'default'">
            {{ action.confirmation.cancel.label }}
          </oc-button>
          <oc-button size="small" :id="action.confirmation.confirm.id" @click="action.confirmation.confirm.handler" :variation="action.confirmation.confirm.variation || 'default'">
            {{ action.confirmation.confirm.label }}
          </oc-button>
        </oc-alert>
        <oc-button v-else :id="action.id" @click="action.handler" :variation="action.variation || 'default'" :icon="action.icon">
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
            variation: 'secondary',
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

<style scoped>
.tmp-alert-fixes {
  padding: 5px 10px 4px !important;
  border-radius: 3px !important;
  background-color: #fff !important;
  border: 1px solid rgb(224, 0, 0) !important;
  color: rgb(224, 0, 0) !important;

  font-size: 1.125rem !important;
  font-weight: 600 !important;
  line-height: 1.4 !important;
}
</style>
