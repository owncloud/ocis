import { ListSettingsBundles, ListSettingsValues } from '../client/settings'

const state = {
  config: null,
  initialized: false,
  settingsBundles: []
}

const getters = {
  config: state => state.config,
  initialized: state => state.initialized,
  settingsBundles: state => state.settingsBundles,
  extensions: state => {
    return [...new Set(Array.from(state.settingsBundles).map(bundle => bundle.identifier.extension))].sort()
  },
  getSettingsBundlesByExtension: state => extension => {
    return state.settingsBundles.filter(bundle => bundle.identifier.extension === extension).sort((b1,b2) => {
      return b1.identifier.bundleKey.localeCompare(b2.identifier.bundleKey)
    })
  }
}

const mutations = {
  SET_INITIALIZED (state, value) {
    state.initialized = value
  },
  SET_SETTINGS_BUNDLES (state, payload) {
    state.settingsBundles = payload
  },
  LOAD_CONFIG (state, config) {
    state.config = config
  }
}

const actions = {
  loadConfig ({ commit }, config) {
    commit('LOAD_CONFIG', config)
  },

  async initialize({ commit, dispatch }) {
    await Promise.all([
      dispatch('fetchSettingsBundles'),
      // dispatch('fetchSettingsValues')
    ])
    commit('SET_INITIALIZED', true)
  },

  async fetchSettingsBundles ({ commit, dispatch, getters }) {
    const response = await ListSettingsBundles({
      $domain: getters.config.url,
      identifierExtension: "ocis-accounts"
    })
    console.log(response)
    if (response.status === 200) {
      // the settings markup has implicit typing. inject an explicit type variable here
      const settingsBundles = response.data.settingsBundles
      if (settingsBundles) {
        settingsBundles.forEach(bundle => {
          bundle.settings.forEach(setting => {
            if (setting['intValue']) {
              setting.type = 'number'
            } else if (setting['stringValue']) {
              setting.type = 'string'
            } else if (setting['boolValue']) {
              setting.type = 'boolean'
            } else if (setting['singleChoiceValue']) {
              setting.type = 'singleChoice'
            } else if (setting['multiChoiceValue']) {
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

  async fetchSettingsValues ({ commit, dispatch, getters }) {
    const response = await ListSettingsValues({
      $domain: getters.config.url,
      identifierAccountUuid: "5681371F-4A6E-43BC-8BB5-9C9237FA9C58",
      identifierExtension: "ocis-accounts",
      identifierBundleKey: "notifications"
    })
    console.log(response)
    if (response.status === 200) {
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
