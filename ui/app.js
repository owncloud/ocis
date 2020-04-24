import 'regenerator-runtime/runtime'
import SettingsApp from './components/SettingsApp.vue'
import store from './store'

// just a dummy function to trick gettext tools
function $gettext(msg) {
  return msg
}

const appInfo = {
  name: $gettext('Settings'),
  id: 'settings',
  icon: 'application',
  isFileEditor: false,
  extensions: [],
  config: {
    url: 'http://localhost:9190'
  }
}

const routes = [
  {
    name: 'settings',
    path: '/:extension?',
    components: {
      app: SettingsApp
    }
  }
]

const navItems = [
  {
    name: $gettext('Settings'),
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
