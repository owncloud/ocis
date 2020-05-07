import {BundleService_ListSettingsBundles, ValueService_ListSettingsValues} from '../client/settings'

const state = {
  config: null,
  initialized: false,
  settingsBundles: {},
  settingsValues: {}
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
  },
  getSettingsValueByIdentifier: state => ({extension, bundleKey, settingKey}) => {
    if (state.settingsValues.has(extension)
      && state.settingsValues.get(extension).has(bundleKey)
      && state.settingsValues.get(extension).get(bundleKey).has(settingKey)) {
      return state.settingsValues.get(extension).get(bundleKey).get(settingKey)
    }
    return null
  }
}

const mutations = {
  SET_INITIALIZED(state, value) {
    state.initialized = value
  },
  SET_SETTINGS_BUNDLES(state, payload) {
    const map = new Map()
    Array.from(payload).forEach(bundle => {
      if (!map.has(bundle.identifier.extension)) {
        map.set(bundle.identifier.extension, new Map())
      }
      map.get(bundle.identifier.extension).set(bundle.identifier.bundleKey, bundle)
    })
    state.settingsBundles = map
  },
  SET_SETTINGS_VALUES(state, payload) {
    const map = new Map()
    Array.from(payload).forEach(value => {
      if (!map.has(value.identifier.extension)) {
        map.set(value.identifier.extension, new Map())
      }
      if (!map.get(value.identifier.extension).has(value.identifier.bundleKey)) {
        map.get(value.identifier.extension).set(value.identifier.bundleKey, new Map())
      }
      map.get(value.identifier.extension).get(value.identifier.bundleKey).set(value.identifier.settingKey, value)
    })
    state.settingsValues = map
  },
  LOAD_CONFIG(state, config) {
    state.config = config
  }
}

const actions = {
  loadConfig({commit}, config) {
    commit('LOAD_CONFIG', config)
  },

  async initialize({commit, dispatch}) {
    await Promise.all([
      dispatch('fetchSettingsBundles'),
      dispatch('fetchSettingsValues')
    ])
    commit('SET_INITIALIZED', true)
  },

  async fetchSettingsBundles({commit, dispatch, getters}) {
    const response = await BundleService_ListSettingsBundles({
      $domain: getters.config.url,
      body: {}
    })
    if (response.status === 201) {
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
      }, {root: true})
    }
  },

  async fetchSettingsValues({commit, dispatch, getters}) {
    const response = await ValueService_ListSettingsValues({
      $domain: getters.config.url,
      body: {
        identifier: {
          account_uuid: "me"
        }
      }
    })
    if (response.status === 201) {
      const settingsValues = response.data.settingsValues
      if (settingsValues) {
        commit('SET_SETTINGS_VALUES', settingsValues)
      } else {
        commit('SET_SETTINGS_VALUES', [])
      }
    } else {
      dispatch('showMessage', {
        title: 'Failed to fetch settings values.',
        desc: response.statusText,
        status: 'danger'
      }, {root: true})
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
