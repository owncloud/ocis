/* eslint-disable camelcase */
import { AccountsService_ListAccounts } from '../client/accounts'
import { RoleService_ListRoles } from '../client/settings'
/* eslint-enable camelcase */
import { injectAuthToken } from '../helpers/auth'

const state = {
  config: null,
  initialized: false,
  accounts: {},
  roles: null
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
  }
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
      body: {
        accountUuid: 'me'
      }
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
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
