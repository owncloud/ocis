import {
  // eslint-disable-next-line camelcase
  BundleService_ListBundles,
  // eslint-disable-next-line camelcase
  ValueService_SaveValue
} from '../client/settings'
import axios from 'axios'
import keyBy from 'lodash/keyBy'

const state = {
  config: null,
  initialized: false,
  bundles: {}
}

const getters = {
  config: state => state.config,
  initialized: state => state.initialized,
  extensions: state => {
    return [...new Set(Object.values(state.bundles).map(bundle => bundle.extension))].sort()
  },
  getBundlesByExtension: state => extension => {
    return Object.values(state.bundles)
      .filter(bundle => bundle.extension === extension)
      .sort((b1, b2) => {
        return b1.name.localeCompare(b2.name)
      })
  },
  getServerForJsClient: (state, getters, rootState, rootGetters) => rootGetters.configuration.server.replace(/\/$/, '')
}

const mutations = {
  SET_INITIALIZED (state, value) {
    state.initialized = value
  },
  SET_BUNDLES (state, bundles) {
    state.bundles = keyBy(bundles, 'id')
  },
  LOAD_CONFIG (state, config) {
    state.config = config
  }
}

const actions = {
  // Used by ocis-web.
  loadConfig ({ commit }, config) {
    commit('LOAD_CONFIG', config)
  },

  async initialize ({ commit, dispatch }) {
    await dispatch('fetchBundles')
    commit('SET_INITIALIZED', true)
  },

  async fetchBundles ({ commit, dispatch, getters, rootGetters }) {
    injectAuthToken(rootGetters)
    try {
      const response = await BundleService_ListBundles({
        $domain: getters.getServerForJsClient,
        body: {}
      })
      if (response.status === 201) {
        // the settings markup has implicit typing. inject an explicit type variable here
        const bundles = response.data.bundles
        if (bundles) {
          bundles.forEach(bundle => {
            bundle.settings.forEach(setting => {
              if (setting.intValue) {
                setting.type = 'number'
              } else if (setting.stringValue) {
                setting.type = 'string'
              } else if (setting.boolValue) {
                setting.type = 'boolean'
              } else if (setting.singleChoiceValue) {
                setting.type = 'singleChoice'
              } else if (setting.multiChoiceValue) {
                setting.type = 'multiChoice'
              } else {
                setting.type = 'unknown'
              }
            })
          })
          commit('SET_BUNDLES', bundles)
        } else {
          commit('SET_BUNDLES', [])
        }
      }
    } catch (err) {
      dispatch('showMessage', {
        title: 'Failed to fetch bundles.',
        status: 'danger'
      }, { root: true })
    }
  },

  async saveValue ({ commit, dispatch, getters, rootGetters }, { setting, payload }) {
    injectAuthToken(rootGetters)
    try {
      const response = await ValueService_SaveValue({
        $domain: getters.getServerForJsClient,
        body: {
          value: payload
        }
      })
      if (response.status === 201 && response.data.value) {
        commit('SET_SETTINGS_VALUE', response.data.value, { root: true })
      }
    } catch (e) {
      dispatch('showMessage', {
        title: `Failed to save »${setting.displayName}«.`,
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
