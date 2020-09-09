<template>
  <oc-table-row>
    <oc-table-cell colspan="9">
      <oc-grid gutter="small">
        <label>
          <oc-text-input
            type="text"
            v-model="username"
            :error-message="usernameError"
            :placeholder="$gettext('Username')"
            @input="checkUsername"
          />
        </label>
        <label>
          <oc-text-input
            type="email"
            v-model="email"
            :error-message="emailError"
            :placeholder="$gettext('Email')"
            @input="checkEmail"
          />
        </label>
        <label class="uk-flex uk-flex-middle">
          <oc-text-input
            :type="passwordInputType"
            v-model="password"
            :error-message="passwordError"
            :placeholder="$gettext('Password')"
            class="uk-margin-xsmall-right"
            @input="checkPassword"
          />
          <oc-button variation="raw" :aria-label="$gettext('Display password')" @click="togglePasswordVisibility">
            <oc-icon name="remove_red_eye" aria-hidden="true" size="small" />
          </oc-button>
        </label>
        <div>
          <oc-button v-text="$gettext('Cancel')" @click="emitCancel" class="uk-margin-xsmall-right" />
          <oc-button v-text="$gettext('Create')" variation="primary" @click="createAccount" />
        </div>
      </oc-grid>
    </oc-table-cell>
  </oc-table-row>
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
    passwordInputType: 'password'
  }),

  methods: {
    ...mapActions('Accounts', ['createNewAccount']),

    emitCancel () {
      this.$emit('cancel')
    },
    createAccount () {
      this.checkUsername()
      this.checkEmail()
      this.checkPassword()

      if (this.usernameError !== '' || this.emailError !== '' || this.passwordError !== '') {
        return
      }

      this.createNewAccount({ username: this.username, email: this.email, password: this.password })
    },

    checkUsername () {
      if (isEmpty(this.username)) {
        debounce(this.usernameError = this.$gettext('Username cannot be empty'), 500)

        return
      }

      this.usernameError = ''
    },

    checkEmail () {
      if (isEmpty(this.email)) {
        debounce(this.emailError = this.$gettext('Email cannot be empty'), 500)

        return
      }

      if (!isEmail(this.email)) {
        debounce(this.emailError = this.$gettext('Invalid email address'), 500)

        return
      }

      this.emailError = ''
    },

    checkPassword () {
      // Later on some restrictions might be applied here
      if (isEmpty(this.password)) {
        debounce(this.passwordError = this.$gettext('Password cannot be empty'), 500)

        return
      }

      this.passwordError = ''
    },

    togglePasswordVisibility () {
      this.passwordInputType === 'password'
        ? this.passwordInputType = 'text'
        : this.passwordInputType = 'password'
    }
  }
}
</script>
