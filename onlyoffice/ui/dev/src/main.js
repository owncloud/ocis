import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import { sync } from 'vuex-router-sync'

// Import the Design System
import ODS from 'owncloud-design-system'
import 'owncloud-design-system/dist/system/system.css'

Vue.config.productionTip = false
Vue.use(ODS)

const registerStoreModule = app => {
  if (app.store.default) {
    return store.registerModule(app.appInfo.name, app.store.default)
  }

  return store.registerModule(app.appInfo.name, app.store)
}

const mount = () => {
  new Vue({
    router,
    store,
    render: h => h(App)
  }).$mount('#app')
}

const loadExtension = extension => {
  // Redirect to default path
  const routes = [
    extension.navItems && {
      path: '/',
      redirect: () => extension.navItems[0].route
    }
  ]

  if (!extension.appInfo) {
    console.error('Tried to load an extension with missing appInfo')
  }

  if (extension.routes) {
    // rewrite relative app routes by adding their corresponding appId as prefix
    extension.routes.forEach(
      r => (r.path = `/${encodeURI(extension.appInfo.id)}${r.path}`)
    )

    // adjust routes in nav items
    if (extension.navItems) {
      extension.navItems.forEach(nav => {
        const r = extension.routes.find(function (element) {
          return element.name === nav.route.name
        })

        if (r) {
          r.meta = r.meta || {}
          r.meta.pageTitle = nav.name
          nav.route.path = nav.route.path || r.path
        } else {
          console.error(`Unknown route name ${nav.route.name}`)
        }
      })
    }

    routes.push(extension.routes)
  }

  if (extension.store) {
    registerStoreModule(extension)
  }

  router.addRoutes(routes.flat())
  sync(store, router)
  mount()
};

(() => {
  // eslint-disable-next-line no-undef
  requirejs([window.location.origin + '/ocis-onlyoffice.js'], loadExtension)
})()
