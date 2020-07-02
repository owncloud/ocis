<template>
  <div>
    <div class="uk-container uk-padding">
      <h1>
        Accounts
      </h1>
      <template v-if="isInitialized">
        <accounts-list :accounts="accounts" />
      </template>
      <oc-loader v-else />
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import AccountsList from './accounts/AccountsList.vue'
export default {
  name: 'App',
  components: { AccountsList },
  computed: {
    ...mapGetters('Accounts', ['isInitialized', 'getAccountsSorted']),
    accounts () {
      return this.getAccountsSorted
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
