import {
  ListSettingsBundles,
  SaveSettingsValue
} from '../client/settings'
import axios from 'axios'

const state = {
  config: null,
  initialized: false,
  settingsBundles: {}
}

const getters = {
  config: state => state.config,
  initialized: state => state.initialized,
  extensions: state => {
    return Array.from(state.settingsBundles.keys()).sort()
  },
  getSettingsBundlesByExtension: state => extension => {
    if (state.settingsBundles.has(extension)) {
      return Array.from(state.settingsBundles.get(extension).values()).sort((b1, b2) => {
        return b1.identifier.bundleKey.localeCompare(b2.identifier.bundleKey)
      })
    }
    return []
  }
}

const mutations = {
  SET_INITIALIZED (state, value) {
    state.initialized = value
  },
  SET_SETTINGS_BUNDLES (state, settingsBundles) {
    const map = new Map()
    Array.from(settingsBundles).forEach(bundle => {
      if (!map.has(bundle.identifier.extension)) {
        map.set(bundle.identifier.extension, new Map())
      }
      map.get(bundle.identifier.extension).set(bundle.identifier.bundleKey, bundle)
    })
    state.settingsBundles = map
  },
  LOAD_CONFIG (state, config) {
    state.config = config
  }
}

const actions = {
  loadConfig ({ commit }, config) {
    commit('LOAD_CONFIG', config)
  },

  async initialize ({ commit, dispatch }) {
    await dispatch('fetchSettingsBundles')
    commit('SET_INITIALIZED', true)
  },

  async fetchSettingsBundles ({ commit, dispatch, getters, rootGetters }) {
    injectAuthToken(rootGetters)
    const response = await ListSettingsBundles({
      $domain: getters.config.url,
      body: {}
    })
    if (response.status === 201) {
      // the settings markup has implicit typing. inject an explicit type variable here
      const settingsBundles = response.data.settingsBundles
      if (settingsBundles) {
        settingsBundles.forEach(bundle => {
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
        commit('SET_SETTINGS_BUNDLES', settingsBundles)
      } else {
        commit('SET_SETTINGS_BUNDLES', [])
      }
    } else {
      dispatch('showMessage', {
        title: 'Failed to fetch settings bundles.',
        desc: response.statusText,
        status: 'danger'
      }, { root: true })
    }
  },

  async saveSettingsValue ({ commit, dispatch, getters, rootGetters }, payload) {
    injectAuthToken(rootGetters)
    const response = await SaveSettingsValue({
      $domain: getters.config.url,
      body: {
        settingsValue: payload
      }
    })
    if (response.status === 201) {
      if (response.data.settingsValue) {
        commit('SET_SETTINGS_VALUE', response.data.settingsValue, { root: true })
      }
    } else {
      dispatch('showMessage', {
        title: 'Failed to save settings value.',
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
