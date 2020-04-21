import { ListSettingsBundles } from '../client/settings'

const state = {
  config: null,
  settingsBundles: []
}

const getters = {
  config: state => state.config,
  settingsBundles: state => state.settingsBundles
}

const mutations = {
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

  async fetchSettingsBundles ({ commit, dispatch, getters }) {
    const response = await ListSettingsBundles({
      $domain: getters.config.url
    })
    if (response.status === 200) {
      console.log(response.data)
      commit('SET_SETTINGS_BUNDLES', response.data)
    } else {
      dispatch('showMessage', {
        title: 'Failed to fetch settings bundles.',
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
