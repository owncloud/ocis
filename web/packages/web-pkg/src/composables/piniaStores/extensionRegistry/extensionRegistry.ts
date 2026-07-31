import { defineStore } from 'pinia'
import { ref, Ref, unref } from 'vue'
import { useConfigStore } from '../config'
import { Extension, ExtensionPoint, ExtensionType, SidebarNavExtension } from './types'

const withScopePrefix = (routePath: string, isVaultScope: boolean) => {
  const normalizedPath = routePath.replace(/^\/vault(?=\/|$)/, '')
  const basePath = normalizedPath.startsWith('/') ? normalizedPath : `/${normalizedPath}`

  if (!isVaultScope) {
    return basePath
  }

  return basePath === '/' ? '/vault' : `/vault${basePath}`
}

const mapNavRoute = (navRoute: SidebarNavExtension['navItem']['route'], isVaultScope: boolean) => {
  if (typeof navRoute === 'string') {
    return withScopePrefix(navRoute, isVaultScope)
  }

  if (!navRoute || typeof navRoute !== 'object' || !('path' in navRoute)) {
    return navRoute
  }

  if (typeof navRoute.path !== 'string') {
    return navRoute
  }

  return {
    ...navRoute,
    path: withScopePrefix(navRoute.path, isVaultScope)
  }
}

export const useExtensionRegistry = defineStore('extensionRegistry', () => {
  const configStore = useConfigStore()

  const extensions = ref<Ref<Extension[]>[]>([])
  const rebuild = ({ route }) => {
    const isVaultScope = unref(route).params?.scope === 'vault'

    extensions.value = unref(extensions).map((extension) =>
      ref(
        unref(extension).map((ext) => {
          if (ext.type !== 'sidebarNav') {
            return ext
          }

          const sidebarExtension = ext as SidebarNavExtension
          return {
            ...sidebarExtension,
            navItem: {
              ...sidebarExtension.navItem,
              route: mapNavRoute(sidebarExtension.navItem.route, isVaultScope)
            }
          }
        })
      )
    )
  }

  const registerExtensions = (e: Ref<Extension[]>) => {
    extensions.value.push(e)
  }
  const unregisterExtensions = (ids: string[]) => {
    extensions.value = unref(extensions)
      .map((e) => ref(unref(e).filter(({ id }) => !ids.includes(id))))
      .filter((e) => unref(e).length)
  }
  const requestExtensions = <T extends Extension>(extensionPoint: ExtensionPoint<T>) => {
    if (!extensionPoint.id || !extensionPoint.extensionType) {
      throw new Error('ExtensionPoint must have an id and an extensionType')
    }

    return unref(extensions).flatMap((e) =>
      unref(e).filter(
        (e) =>
          e.type === extensionPoint.extensionType &&
          !configStore.options.disabledExtensions.includes(e.id) &&
          (!e.extensionPointIds || e.extensionPointIds?.includes(extensionPoint.id))
      )
    ) as T[]
  }

  const extensionPoints = ref<Ref<ExtensionPoint<Extension>[]>[]>([])
  const registerExtensionPoints = <T extends Extension>(e: Ref<ExtensionPoint<T>[]>) => {
    extensionPoints.value.push(e)
  }
  const unregisterExtensionPoints = (ids: string[]) => {
    extensionPoints.value = unref(extensionPoints)
      .map((e) => ref(unref(e).filter(({ id }) => !ids.includes(id))))
      .filter((e) => unref(e).length)
  }
  const getExtensionPoints = <T extends ExtensionPoint<Extension>>(
    options: {
      extensionType?: ExtensionType
    } = {}
  ) => {
    return unref(extensionPoints).flatMap(
      (e) =>
        unref(e).filter((e) => {
          if (
            Object.hasOwn(options, 'extensionType') &&
            e.extensionType !== options.extensionType
          ) {
            return false
          }
          return true
        }) as T[]
    )
  }

  return {
    extensions,
    registerExtensions,
    unregisterExtensions,
    requestExtensions,
    extensionPoints,
    registerExtensionPoints,
    unregisterExtensionPoints,
    getExtensionPoints,
    rebuild
  }
})

export type ExtensionRegistry = ReturnType<typeof useExtensionRegistry>
