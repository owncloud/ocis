import { Router, RouteRecordRaw } from 'vue-router'
import { App, Component } from 'vue'
import { AppNavigationItem } from '@ownclouders/web-pkg'

/** RuntimeApi defines the publicly available runtime api */
export interface RuntimeApi {
  announceRoutes: (routes: RouteRecordRaw[]) => void
  announceNavigationItems: (navigationItems: AppNavigationItem[]) => void
  requestRouter: () => Router
  openPortal: (
    instance: App,
    toApp: string,
    toPortal: string,
    order: number,
    components: Component[]
  ) => void
}
