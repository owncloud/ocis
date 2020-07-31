import {
  // eslint-disable-next-line camelcase
  AccountsService_ListAccounts
} from '../client/accounts'
import axios from 'axios'

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
    injectAuthToken(rootGetters)
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

  async fetchRoles ({ commit, rootGetters }) {
    const headers = new Headers()

    headers.append('Authorization', 'Bearer ' + rootGetters.getToken)

    let roles = await fetch(`${rootGetters.configuration.server}/api/v0/settings/roles-list`, {
      method: 'POST',
      mode: 'cors',
      headers: this.headers,
      body: JSON.stringify({
        account_uuid: rootGetters.user.id
      })
    })

    roles = await roles.json()
    commit('SET_ROLES', roles.bundles)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}

function injectAuthToken (rootGetters) {
  axios.interceptors.request.use(config => {
    if (typeof config.headers.Authorization === 'undefined') {
      const token = rootGetters.user.token
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
    }
    return config
  })
}
