import 'regenerator-runtime/runtime'
import SettingsApp from './components/SettingsApp.vue'
import store from './store'
import translations from './../l10n/translations.json'

// just a dummy function to trick gettext tools
function $gettext (msg) {
  return msg
}

const appInfo = {
  name: $gettext('Settings'),
  id: 'settings',
  icon: 'application',
  isFileEditor: false,
  extensions: []
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
    },
    menu: 'user'
  }
]

export default {
  appInfo,
  store,
  routes,
  navItems,
  translations
}
