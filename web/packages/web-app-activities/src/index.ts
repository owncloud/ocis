import translations from '../l10n/translations.json'
import { useGettext } from 'vue3-gettext'
import { computed } from 'vue'
import { AppMenuItemExtension, defineWebApplication, Extension } from '@ownclouders/web-pkg'
import { urlJoin } from '@ownclouders/web-client'
import { RouteRecordRaw } from 'vue-router'
import { APPID } from './appid'

export default defineWebApplication({
  setup() {
    const { $gettext } = useGettext()

    const appInfo = {
      name: $gettext('Activities'),
      id: APPID,
      icon: 'pulse',
      color: '#887ef1'
    }

    const routes: RouteRecordRaw[] = [
      {
        path: '/',
        name: 'root',
        component: () => import('./LayoutContainer.vue'),
        redirect: urlJoin(appInfo.id, 'list'),
        meta: {
          authContext: 'user'
        },
        children: [
          {
            path: 'list',
            name: 'list',
            component: () => import('./views/App.vue'),
            meta: {
              authContext: 'user',
              title: $gettext('Activities')
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

      result.push(menuItemExtension)

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
