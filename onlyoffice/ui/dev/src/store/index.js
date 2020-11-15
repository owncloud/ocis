import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    configuration: {
      server: 'http://localhost:9105'
    },
    user: {
      token: 'test'
    }
  },
  getters: {
    configuration: state => state.configuration,
    user: state => state.user
  },
  mutations: {},
  actions: {},
  modules: {}
})
