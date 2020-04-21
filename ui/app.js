import SettingsApp from './components/SettingsApp.vue'
import store from './store'

const appInfo = {
  name: 'Settings',
  id: 'settings',
  icon: 'gear',
  isFileEditor: false,
  extensions: [],
  config: {
    url: 'http://localhost:9190'
  }
}

const routes = [
  {
    name: 'settings',
    path: '/',
    components: {
      app: SettingsApp
    }
  }
]

const navItems = [
  {
    name: 'Settings',
    iconMaterial: appInfo.icon,
    route: {
      name: 'settings',
      path: `/${appInfo.id}/`
    }
  }
]

export default {
  appInfo,
  store,
  routes,
  navItems
}
