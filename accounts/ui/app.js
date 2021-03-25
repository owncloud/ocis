import 'regenerator-runtime/runtime'
import App from './components/App.vue'
import store from './store'
import translations from './../l10n/translations.json'

// just a dummy function to trick gettext tools
function $gettext (msg) {
  return msg
}

const appInfo = {
  name: $gettext('Accounts'),
  id: 'accounts',
  icon: 'text-vcard',
  isFileEditor: false,
  extensions: []
}

const routes = [
  {
    name: 'accounts',
    path: '/',
    components: {
      app: App
    }
  }
]

const navItems = [
  {
    name: $gettext('Accounts'),
    iconMaterial: appInfo.icon,
    route: {
      name: 'accounts',
      path: `/${appInfo.id}/`
    },
    menu: 'apps'
  }
]

export default {
  appInfo,
  routes,
  navItems,
  store,
  translations
}
