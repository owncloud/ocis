import translations from '../l10n/translations.json'
import { useGettext } from 'vue3-gettext'
import { computed, unref } from 'vue'
import {
  AppMenuItemExtension,
  defineWebApplication,
  Extension,
  useAbility,
  useUserStore
} from '@ownclouders/web-pkg'
import { urlJoin } from '@ownclouders/web-client'
import { RouteRecordRaw } from 'vue-router'
import { useRepositoriesStore } from './piniaStores'
import { AppStoreConfigSchema, AppStoreRepository } from './types'
import { APPID } from './appid'

export default defineWebApplication({
  setup({ applicationConfig }) {
    const { $gettext } = useGettext()
    const { can } = useAbility()
    const userStore = useUserStore()
    const repositoryStore = useRepositoriesStore()

    const defaultRepositories: AppStoreRepository[] = [
      {
        name: 'awesome-ocis',
        url: 'https://raw.githubusercontent.com/owncloud/awesome-ocis/main/webApps/apps.json'
      }
    ]
    if (applicationConfig?.repositories) {
      const { repositories } = AppStoreConfigSchema.parse(applicationConfig)
      repositoryStore.setRepositories(repositories || defaultRepositories)
    } else {
      repositoryStore.setRepositories(defaultRepositories)
    }

    const appInfo = {
      name: $gettext('App Store'),
      id: APPID,
      icon: 'store',
      color: '#ff6961'
    }

    const hasPermission = computed(() => {
      // TODO: which permission(s) do we need to check here?
      return userStore.user && can('read-all', 'Setting')
    })

    const routes: RouteRecordRaw[] = [
      {
        path: '/',
        name: 'root',
        component: () => import('./LayoutContainer.vue'),
        redirect: urlJoin(appInfo.id, 'list'),
        beforeEnter: (to, from, next) => {
          if (!unref(hasPermission)) {
            return next({ path: '/' })
          }
          next()
        },
        meta: {
          authContext: 'user'
        },
        children: [
          {
            path: 'list',
            name: 'list',
            component: () => import('./views/AppList.vue'),
            meta: {
              authContext: 'user',
              title: $gettext('App Store')
            }
          },
          {
            path: 'app/:appId',
            name: 'details',
            component: () => import('./views/AppDetails.vue'),
            meta: {
              authContext: 'user',
              title: $gettext('App Details')
            }
          }
        ]
      }
    ]

    const menuItemExtension: AppMenuItemExtension = {
      id: `app.${appInfo.id}.menuItem`,
      type: 'appMenuItem',
      label: () => appInfo.name,
      color: appInfo.color,
      icon: appInfo.icon,
      priority: 30,
      path: urlJoin(appInfo.id)
    }
    const extensions = computed(() => {
      const result: Extension[] = []

      if (unref(hasPermission)) {
        result.push(menuItemExtension)
      }

      return result
    })

    return {
      appInfo,
      routes,
      translations,
      extensions
    }
  }
})
