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
  failed: false,
  accounts: {},
  roles: null,
  selectedAccounts: []
}

const getters = {
  config: state => state.config,
  isInitialized: state => state.initialized,
  hasFailed: state => state.failed,
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
  SET_FAILED (state, value) {
    state.failed = value
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

  async initialize ({ commit, dispatch, getters }) {
    await Promise.all([
      dispatch('fetchAccounts'),
      dispatch('fetchRoles')
    ])
    if (!getters.hasFailed) {
      commit('SET_INITIALIZED', true)
    }
  },

  async fetchAccounts ({ commit, rootGetters }) {
    injectAuthToken(rootGetters.user.token)
    try {
      const response = await AccountsService_ListAccounts({
        $domain: rootGetters.configuration.server.replace(/\/$/, ''),
        body: {}
      })
      if (response.status === 201) {
        const accounts = response.data.accounts
        commit('SET_ACCOUNTS', accounts || [])
        return
      }
    } catch (e) {
    }
    commit('SET_FAILED', true)
  },

  async fetchRoles ({ commit, rootGetters }) {
    injectAuthToken(rootGetters.user.token)
    try {
      const response = await RoleService_ListRoles({
        $domain: rootGetters.configuration.server.replace(/\/$/, ''),
        body: {}
      })
      if (response.status === 201) {
        const roles = response.data.bundles
        commit('SET_ROLES', roles || [])
        return
      }
    } catch (e) {
    }
    commit('SET_FAILED', true)
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

      try {
        const response = await AccountsService_UpdateAccount({
          $domain: rootGetters.configuration.server.replace(/\/$/, ''),
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
          failedAccounts.push({ account: account.username })
        }
      } catch (error) {
        failedAccounts.push({ account: account.username })
      }
    }

    if (failedAccounts.length > 0) {
      let errorTitle = ''
      if (failedAccounts.length === 1) {
        errorTitle = activated ? 'Failed to activate account.' : 'Failed to block account.'
      } else {
        errorTitle = activated ? 'Failed to activate accounts.' : 'Failed to block accounts.'
      }
      dispatch('showMessage', {
        title: errorTitle,
        status: 'danger'
      }, { root: true })
      return Promise.resolve(false)
    }

    commit('RESET_ACCOUNTS_SELECTION')
    return Promise.resolve(true)
  },
  async createNewAccount ({ rootGetters, commit, dispatch }, account) {
    injectAuthToken(rootGetters.user.token)

    try {
      const response = await AccountsService_CreateAccount({
        $domain: rootGetters.configuration.server.replace(/\/$/, ''),
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
        return Promise.resolve(true)
      }
    } catch (error) {
      dispatch('showMessage', {
        title: 'Failed to create account.',
        status: 'danger'
      }, { root: true })
      return Promise.reject(error)
    }
    return Promise.resolve(false)
  },

  async deleteAccounts ({ rootGetters, state, commit, dispatch }) {
    const failedAccounts = []

    injectAuthToken(rootGetters.user.token)

    for (const account of state.selectedAccounts) {
      try {
        const response = await AccountsService_DeleteAccount({
          $domain: rootGetters.configuration.server.replace(/\/$/, ''),
          body: {
            id: account.id
          }
        })
        if (response.status === 201 || response.status === 204) {
          commit('DELETE_ACCOUNT', account.id)
        } else {
          failedAccounts.push({ account: account.username })
        }
      } catch (error) {
        failedAccounts.push({ account: account.username })
      }
    }

    if (failedAccounts.length > 0) {
      const errorTitle = failedAccounts.length === 1 ? 'Failed to delete account.' : 'Failed to delete accounts.'
      dispatch('showMessage', {
        title: errorTitle,
        status: 'danger'
      }, { root: true })
      return Promise.resolve(false)
    }

    commit('RESET_ACCOUNTS_SELECTION')
    return Promise.resolve(true)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
