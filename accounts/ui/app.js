import 'regenerator-runtime/runtime'
import App from './components/App.vue'
import store from './store'

const appInfo = {
  name: 'Accounts',
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
    name: 'Accounts',
    iconMaterial: appInfo.icon,
    route: {
      name: 'accounts',
      path: `/${appInfo.id}/`
    },
    menu: 'user'
  }
]

export default {
  appInfo,
  routes,
  navItems,
  store
}
