import { AppNavigationItem, ExtensionRegistry, SidebarNavExtension } from '@ownclouders/web-pkg'

export interface NavItem extends Omit<AppNavigationItem, 'name'> {
  name: string
  active: boolean
}

export const getExtensionNavItems = ({
  extensionRegistry,
  appId
}: {
  extensionRegistry: ExtensionRegistry
  appId: string
}) =>
  extensionRegistry
    .requestExtensions<SidebarNavExtension>({
      id: `app.${appId}.navItems`,
      extensionType: 'sidebarNav'
    })
    .map(({ navItem }) => navItem)
    .filter((n) => !Object.hasOwn(n, 'isVisible') || n.isVisible())
