<template>
  <div>
    <oc-grid v-if="isFormInProgress" gutter="small">
      <oc-text-input
        id="accounts-new-account-input-username"
        type="text"
        v-model="formData.username"
        :error-message="formValidation.usernameError"
        :label="$gettext('Username')"
        :disabled="isRequestInProgress"
        @keydown.enter="createAccount"
      />
      <oc-text-input
        id="accounts-new-account-input-email"
        type="email"
        v-model="formData.email"
        :error-message="formValidation.emailError"
        :label="$gettext('Email')"
        :disabled="isRequestInProgress"
        @keydown.enter="createAccount"
      />
      <oc-text-input
        id="accounts-new-account-input-password"
        type="password"
        v-model="formData.password"
        :error-message="formValidation.passwordError"
        :label="$gettext('Password')"
        :disabled="isRequestInProgress"
        @keydown.enter="createAccount"
      />
      <div class="uk-flex">
        <oc-button
          class="oc-mr-s oc-mb-s"
          v-text="$gettext('Cancel')"
          @click="cancelForm"
          :disabled="isRequestInProgress"
        />
        <oc-button
          id="accounts-new-account-button-confirm"
          class="oc-mr-s oc-mb-s"
          variation="primary"
          appearance="filled"
          :disabled="isRequestInProgress"
          @click="createAccount"
          gap-size="small"
          :class="{ 'border-ods-tmp-fix': !isRequestInProgress }"
        >
          <oc-spinner
            v-if="isRequestInProgress"
            key="account-creation-in-progress"
            size="small"
            aria-hidden="true"
          />
          <span v-text="isRequestInProgress ? $gettext('Creating') : $gettext('Create')" />
        </oc-button>
      </div>
    </oc-grid>
    <oc-grid v-else gutter="small">
      <div>
        <oc-button
          id="accounts-new-account-trigger"
          key="create-accounts-button"
          variation="primary"
          appearance="filled"
          gap-size="small"
          @click="setFormInProgress(true)"
        >
          <oc-icon name="add" />
          <translate>Create new account</translate>
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
  name: 'AccountsCreate',

  data: () => ({
    isFormInProgress: false,
    isRequestInProgress: false,
    formData: {
      username: '',
      email: '',
      password: ''
    },
    formValidation: {
      usernameError: '',
      emailError: '',
      passwordError: ''
    }
  }),

  methods: {
    ...mapActions('Accounts', ['createNewAccount']),

    setFormInProgress (inProgress) {
      this.isFormInProgress = inProgress
    },

    cancelForm () {
      this.isRequestInProgress = false
      this.setFormInProgress(false)
      this.formData = {
        username: '',
        email: '',
        password: ''
      }
      this.formValidation = {
        usernameError: '',
        emailError: '',
        passwordError: ''
      }
    },

    createAccount () {
      // note: use bitwise AND because we want all checks to be performed
      if (!(this.checkUsername() & this.checkEmail() & this.checkPassword())) {
        return
      }

      this.isRequestInProgress = true
      this.createNewAccount(this.formData)
        .then((success) => {
          if (success) {
            this.cancelForm()
          }
        })
        .finally(() => {
          this.isRequestInProgress = false
        })
    },

    checkUsername () {
      if (isEmpty(this.formData.username)) {
        debounce(this.formValidation.usernameError = this.$gettext('Username cannot be empty'), 500)
        return false
      }
      // hacky check: we want to allow emails and the username part of emails as username
      if (!isEmail(this.formData.username) && !isEmail(this.formData.username + '@validate.it')) {
        debounce(this.formValidation.usernameError = this.$gettext('Invalid username'), 500)
        return false
      }

      this.formValidation.usernameError = ''
      return true
    },

    checkEmail () {
      if (isEmpty(this.formData.email)) {
        debounce(this.formValidation.emailError = this.$gettext('Email cannot be empty'), 500)
        return false
      }

      if (!isEmail(this.formData.email)) {
        debounce(this.formValidation.emailError = this.$gettext('Invalid email address'), 500)
        return false
      }

      this.formValidation.emailError = ''
      return true
    },

    checkPassword () {
      // Later on some restrictions might be applied here
      if (isEmpty(this.formData.password)) {
        debounce(this.formValidation.passwordError = this.$gettext('Password cannot be empty'), 500)
        return false
      }

      this.formValidation.passwordError = ''
      return true
    }
  },
  onDestroy () {
    this.cancelForm()
  }
}
</script>

<style>
#accounts-new-account-button-confirm > span {
  display: flex;
  align-items: center;
}
</style>
