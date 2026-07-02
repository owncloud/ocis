import translations from '../l10n/translations.json'
import Resolve from './views/Resolve.vue'

// just a dummy function to trick gettext tools
function $gettext(msg: string) {
  return msg
}

const appInfo = {
  name: $gettext('Webfinger'),
  id: 'webfinger',
  icon: 'fingerprint'
}

const routes = () => [
  {
    name: 'webfinger-root',
    path: '/',
    redirect: () => {
      return { name: 'webfinger-resolve' }
    },
    meta: {
      authContext: 'anonymous'
    }
  },
  {
    path: '/resolve',
    name: 'webfinger-resolve',
    component: Resolve,
    meta: {
      authContext: 'idp',
      title: $gettext('Resolve destination'),
      entryPoint: true
    }
  }
]

export default {
  appInfo,
  routes,
  translations
}
