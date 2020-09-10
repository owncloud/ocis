/* eslint-disable camelcase */
import {
  AccountsService_ListAccounts,
  AccountsService_UpdateAccount,
  AccountsService_CreateAccount,
  AccountsService_DeleteAccount
} from '../client/accounts'
import { RoleService_ListRoles } from '../client/settings'
/* eslint-enable camelcase */
import { injectAuthToken } from '../helpers/auth'

const state = {
  config: null,
  initialized: false,
  accounts: {},
  roles: null,
  selectedAccounts: []
}

const getters = {
  config: state => state.config,
  isInitialized: state => state.initialized,
  getAccountsSorted: state => {
    return Object.values(state.accounts).sort((a1, a2) => {
      if (a1.onPremisesSamAccountName === a2.onPremisesSamAccountName) {
        return a1.id.localeCompare(a2.id)
      }
      return a1.onPremisesSamAccountName.localeCompare(a2.onPremisesSamAccountName)
    })
  },
  areAllAccountsSelected: state => state.accounts.length === state.selectedAccounts.length,
  isAnyAccountSelected: state => state.selectedAccounts.length > 0
}

const mutations = {
  LOAD_CONFIG (state, config) {
    state.config = config
  },
  SET_INITIALIZED (state, value) {
    state.initialized = value
  },
  SET_ACCOUNTS (state, accounts) {
    state.accounts = accounts
  },
  SET_ROLES (state, roles) {
    state.roles = roles
  },
  TOGGLE_SELECTION_ACCOUNT (state, account) {
    const accountIndex = state.selectedAccounts.indexOf(account)

    accountIndex > -1 ? state.selectedAccounts.splice(accountIndex, 1) : state.selectedAccounts.push(account)
  },
  SET_SELECTED_ACCOUNTS (state, accounts) {
    state.selectedAccounts = accounts
  },

  UPDATE_ACCOUNT (state, updatedAccount) {
    const accountIndex = state.accounts.findIndex(account => account.id === updatedAccount.id)

    state.accounts.splice(accountIndex, 1, updatedAccount)
  },

  RESET_ACCOUNTS_SELECTION (state) {
    state.selectedAccounts = []
  },

  PUSH_NEW_ACCOUNT (state, account) {
    state.accounts.push(account)
  },

  DELETE_ACCOUNT (state, accountId) {
    const accountIndex = state.accounts.findIndex(account => account.id === accountId)

    state.accounts.splice(accountIndex, 1)
  }
}

const actions = {
  loadConfig ({ commit }, config) {
    commit('LOAD_CONFIG', config)
  },

  async initialize ({ commit, dispatch }) {
    await dispatch('fetchAccounts')
    await dispatch('fetchRoles')
    commit('SET_INITIALIZED', true)
  },

  async fetchAccounts ({ commit, dispatch, rootGetters }) {
    injectAuthToken(rootGetters.user.token)
    const response = await AccountsService_ListAccounts({
      $domain: rootGetters.configuration.server,
      body: {}
    })
    if (response.status === 201) {
      const accounts = response.data.accounts
      commit('SET_ACCOUNTS', accounts || [])
    } else {
      dispatch('showMessage', {
        title: 'Failed to fetch accounts.',
        desc: response.statusText,
        status: 'danger'
      }, { root: true })
    }
  },

  async fetchRoles ({ commit, dispatch, rootGetters }) {
    injectAuthToken(rootGetters.user.token)

    const response = await RoleService_ListRoles({
      $domain: rootGetters.configuration.server,
      body: {}
    })

    if (response.status === 201) {
      const roles = response.data.bundles

      commit('SET_ROLES', roles || [])
    } else {
      dispatch('showMessage', {
        title: 'Failed to fetch roles.',
        desc: response.statusText,
        status: 'danger'
      }, { root: true })
    }
  },

  toggleSelectionAll ({ commit, getters, state }) {
    getters.areAllAccountsSelected ? commit('RESET_ACCOUNTS_SELECTION') : commit('SET_SELECTED_ACCOUNTS', [...state.accounts])
  },

  async setAccountActivated ({ commit, dispatch, state, rootGetters }, activated) {
    const failedAccounts = []
    injectAuthToken(rootGetters.user.token)

    for (const account of state.selectedAccounts) {
      if (account.accountEnabled === activated) {
        continue
      }

      const response = await AccountsService_UpdateAccount({
        $domain: rootGetters.configuration.server,
        body: {
          account: {
            id: account.id,
            accountEnabled: activated
          },
          update_mask: {
            paths: ['AccountEnabled']
          }
        }
      })

      if (response.status === 201) {
        commit('UPDATE_ACCOUNT', { ...account, accountEnabled: activated })
      } else {
        failedAccounts.push({ account: account.displayName, statusText: response.statusText })
      }
    }

    if (failedAccounts.length === 1) {
      const failedMessageTitle = activated ? 'Failed to activate account.' : 'Failed to block account.'

      dispatch('showMessage', {
        title: failedMessageTitle,
        desc: failedAccounts[0].statusText,
        status: 'danger'
      }, { root: true })
    }

    if (failedAccounts.length > 1) {
      const failedMessageTitle = activated ? 'Failed to activate accounts.' : 'Failed to block accounts.'
      const failedMessageDesc = activated ? 'Could not activate multiple accounts.' : 'Could not block multiple accounts.'

      dispatch('showMessage', {
        title: failedMessageTitle,
        desc: failedMessageDesc,
        status: 'danger'
      }, { root: true })
    }

    commit('RESET_ACCOUNTS_SELECTION')
  },
  async createNewAccount ({ rootGetters, commit, dispatch }, account) {
    injectAuthToken(rootGetters.user.token)

    const response = await AccountsService_CreateAccount({
      $domain: rootGetters.configuration.server,
      body: {
        account: {
          on_premises_sam_account_name: account.username,
          preferred_name: account.username,
          mail: account.email,
          password_profile: {
            password: account.password
          },
          account_enabled: true,
          display_name: account.username
        }
      }
    })

    if (response.status === 201) {
      commit('PUSH_NEW_ACCOUNT', response.data)
    } else {
      dispatch('showMessage', {
        title: 'Failed to create account',
        desc: response.statusText,
        status: 'danger'
      }, { root: true })
    }
  },

  async deleteAccounts ({ rootGetters, state, commit, dispatch }) {
    const failedAccounts = []

    injectAuthToken(rootGetters.user.token)

    for (const account of state.selectedAccounts) {
      const response = await AccountsService_DeleteAccount({
        $domain: rootGetters.configuration.server,
        body: {
          id: account.id
        }
      })

      if (response.status === 201 || response.status === 204) {
        commit('DELETE_ACCOUNT', account.id)
      } else {
        failedAccounts.push({ account: account.diisplayName, statusText: response.statusText })
      }
    }

    if (failedAccounts.length === 1) {
      dispatch('showMessage', {
        title: 'Failed to delete account',
        desc: failedAccounts[0].statusText,
        status: 'danger'
      }, { root: true })
    }

    if (failedAccounts.length > 1) {
      dispatch('showMessage', {
        title: 'Failed to delete accounts',
        desc: 'Could not delete multiple accounts',
        status: 'danger'
      }, { root: true })
    }

    commit('RESET_ACCOUNTS_SELECTION')
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
