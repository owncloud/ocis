import translations from '../l10n/translations.json'
import General from './views/General.vue'
import Users from './views/Users.vue'
import Groups from './views/Groups.vue'
import Spaces from './views/Spaces.vue'
import { Ability, urlJoin } from '@ownclouders/web-client'
import {
  ApplicationInformation,
  AppMenuItemExtension,
  AppNavigationItem,
  defineWebApplication,
  useAbility,
  useAuthService,
  useCapabilityStore,
  useUserStore
} from '@ownclouders/web-pkg'
import { RouteRecordRaw } from 'vue-router'
import { computed } from 'vue'

// just a dummy function to trick gettext tools
function $gettext(msg: string) {
  return msg
}

const appId = 'admin-settings'

function getNextAvailableRoute(ability: Ability) {
  if (ability.can('read-all', 'Setting')) {
    return { name: 'admin-settings-general' }
  }

  if (ability.can('read-all', 'Account')) {
    return { name: 'admin-settings-users' }
  }

  if (ability.can('read-all', 'Group')) {
    return { name: 'admin-settings-groups' }
  }

  if (ability.can('read-all', 'Drive')) {
    return { name: 'admin-settings-spaces' }
  }

  throw Error('Insufficient permissions')
}

async function requireAcr(redirectUrl: string) {
  const capabilityStore = useCapabilityStore()
  if (!capabilityStore.authMfaEnabled) {
    return
  }

  const authService = useAuthService()
  await authService.requireAcr(capabilityStore.authMfaRequiredLevelname, redirectUrl)
}

export const routes = ({ $ability }: { $ability: Ability }): RouteRecordRaw[] => [
  {
    path: '/',
    redirect: () => {
      return { name: 'admin-settings-general' }
    }
  },
  {
    path: '/general',
    name: 'admin-settings-general',
    component: General,
    beforeEnter: async (to, from, next) => {
      if (!$ability.can('read-all', 'Setting')) {
        return next(getNextAvailableRoute($ability))
      }

      await requireAcr(to.fullPath)

      next()
    },
    meta: {
      authContext: 'user',
      title: $gettext('General')
    }
  },
  {
    path: '/users',
    name: 'admin-settings-users',
    component: Users,
    beforeEnter: async (to, from, next) => {
      if (!$ability.can('read-all', 'Account')) {
        return next(getNextAvailableRoute($ability))
      }

      await requireAcr(to.fullPath)

      next()
    },
    meta: {
      authContext: 'user',
      title: $gettext('Users')
    }
  },
  {
    path: '/groups',
    name: 'admin-settings-groups',
    component: Groups,
    beforeEnter: async (to, from, next) => {
      if (!$ability.can('read-all', 'Group')) {
        return next(getNextAvailableRoute($ability))
      }

      await requireAcr(to.fullPath)

      next()
    },
    meta: {
      authContext: 'user',
      title: $gettext('Groups')
    }
  },
  {
    path: '/spaces',
    name: 'admin-settings-spaces',
    component: Spaces,
    beforeEnter: async (to, from, next) => {
      if (!$ability.can('read-all', 'Drive')) {
        return next(getNextAvailableRoute($ability))
      }

      await requireAcr(to.fullPath)

      next()
    },
    meta: {
      authContext: 'user',
      title: $gettext('Spaces')
    }
  }
]

export const navItems = ({ $ability }: { $ability: Ability }): AppNavigationItem[] => [
  {
    name: $gettext('General'),
    icon: 'settings-4',
    route: {
      path: `/${appId}/general?`
    },
    isVisible: () => {
      return $ability.can('read-all', 'Setting')
    },
    priority: 10
  },
  {
    name: $gettext('Users'),
    icon: 'user',
    route: {
      path: `/${appId}/users?`
    },
    isVisible: () => {
      return $ability.can('read-all', 'Account')
    },
    priority: 20
  },
  {
    name: $gettext('Groups'),
    icon: 'group-2',
    route: {
      path: `/${appId}/groups?`
    },
    isVisible: () => {
      return $ability.can('read-all', 'Group')
    },
    priority: 30
  },
  {
    name: $gettext('Spaces'),
    icon: 'layout-grid',
    route: {
      path: `/${appId}/spaces?`
    },
    isVisible: () => {
      return $ability.can('read-all', 'Drive')
    },
    priority: 40
  }
]

export default defineWebApplication({
  setup() {
    const { can } = useAbility()
    const userStore = useUserStore()

    const appInfo: ApplicationInformation = {
      name: $gettext('Admin Settings'),
      id: appId,
      icon: 'settings-4',
      color: '#2b2b2b'
    }

    const menuItems = computed<AppMenuItemExtension[]>(() => {
      const items: AppMenuItemExtension[] = []

      const menuItemAvailable =
        userStore.user &&
        (can('read-all', 'Setting') ||
          can('read-all', 'Account') ||
          can('read-all', 'Group') ||
          can('read-all', 'Drive'))

      if (menuItemAvailable) {
        items.push({
          id: `app.${appInfo.id}.menuItem`,
          type: 'appMenuItem',
          label: () => appInfo.name,
          color: appInfo.color,
          icon: appInfo.icon,
          priority: 40,
          path: urlJoin(appInfo.id)
        })
      }

      return items
    })

    return {
      appInfo,
      routes,
      navItems,
      translations,
      extensions: menuItems
    }
  }
})
