import 'regenerator-runtime/runtime'
import App from './components/App.vue'

const appInfo = {
  name: 'Accounts',
  id: 'accounts',
  icon: 'text-vcard',
  isFileEditor: false,
  extensions: [],
  config: {
    url: 'https://localhost:9200'
  }
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
    }
  }
]

export default {
  appInfo,
  routes,
  navItems
}
