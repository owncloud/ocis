<template>
  <div class="uk-flex uk-flex-top">
    <oc-grid gutter="small">
      <label>
        <oc-text-input
          id="accounts-new-account-input-username"
          type="text"
          v-model="username"
          :error-message="usernameError"
          :placeholder="$gettext('Username')"
          :disabled="isInProgress"
          @keydown.enter="createAccount"
        />
      </label>
      <label>
        <oc-text-input
          id="accounts-new-account-input-email"
          type="email"
          v-model="email"
          :error-message="emailError"
          :placeholder="$gettext('Email')"
          :disabled="isInProgress"
          @keydown.enter="createAccount"
        />
      </label>
      <label class="uk-margin-xsmall-right">
        <oc-text-input
          id="accounts-new-account-input-password"
          type="password"
          v-model="password"
          :error-message="passwordError"
          :placeholder="$gettext('Password')"
          :disabled="isInProgress"
          @keydown.enter="createAccount"
        />
      </label>
      <div>
        <oc-button
          v-text="$gettext('Cancel')"
          @click="emitClose"
          class="uk-margin-xsmall-right"
          :disabled="isInProgress"
        />
        <oc-button
          id="accounts-new-account-button-confirm"
          variation="primary"
          :disabled="isInProgress"
          @click="createAccount"
        >
          <oc-spinner
            v-if="isInProgress"
            key="account-creation-in-progress"
            size="xsmall"
            class="uk-margin-xsmall-right"
            aria-hidden="true"
          />
          <span v-text="isInProgress ? $gettext('Creating') : $gettext('Create')" />
        </oc-button>
      </div>
    </oc-grid>
  </div>
</template>

<script>
import isEmail from 'validator/es/lib/isEmail'
import isEmpty from 'validator/es/lib/isEmpty'
import debounce from 'debounce'
import { mapActions } from 'vuex'

export default {
  name: 'AccountsListNewAccountRow',

  data: () => ({
    username: '',
    usernameError: '',
    email: '',
    emailError: '',
    password: '',
    passwordError: '',
    isInProgress: false
  }),

  methods: {
    ...mapActions('Accounts', ['createNewAccount']),

    emitClose () {
      this.$emit('close')
    },
    createAccount () {
      if (!(this.checkUsername() & this.checkEmail() & this.checkPassword())) {
        return
      }

      this.isInProgress = true
      this.createNewAccount({ username: this.username, email: this.email, password: this.password }).finally(() => {
        this.isInProgress = false
      })
      this.emitClose()
    },

    checkUsername () {
      if (isEmpty(this.username)) {
        debounce(this.usernameError = this.$gettext('Username cannot be empty'), 500)
        return false
      }

      this.usernameError = ''
      return true
    },

    checkEmail () {
      if (isEmpty(this.email)) {
        debounce(this.emailError = this.$gettext('Email cannot be empty'), 500)
        return false
      }

      if (!isEmail(this.email)) {
        debounce(this.emailError = this.$gettext('Invalid email address'), 500)
        return false
      }

      this.emailError = ''
      return true
    },

    checkPassword () {
      // Later on some restrictions might be applied here
      if (isEmpty(this.password)) {
        debounce(this.passwordError = this.$gettext('Password cannot be empty'), 500)
        return false
      }

      this.passwordError = ''
      return true
    }
  }
}
</script>

<style>
#accounts-new-account-button-confirm > span {
  display: flex;
  align-items: center;
}
</style>
