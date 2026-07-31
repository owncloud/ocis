import { defineComponent } from 'vue'
import { IconFillType } from '../../helpers/resource'
import { Item } from '@ownclouders/web-client'

export interface SideBarPanelContext<R extends Item, P extends Item, T extends Item> {
  root?: R
  parent?: P
  items?: T[]
}

export interface SideBarPanel<R extends Item, P extends Item, T extends Item> {
  name: string
  icon: string
  iconFillType?: IconFillType
  title(context: SideBarPanelContext<R, P, T>): string
  isVisible(context: SideBarPanelContext<R, P, T>): boolean
  component: ReturnType<typeof defineComponent>
  componentAttrs?(context: SideBarPanelContext<R, P, T>): any
  /**
   * defines if the panel is a `root` level panel in the right sidebar.
   * In the long run this should be configured by admins or even users, as it's ideally not to be decided by extension developer.
   */
  isRoot?(context: SideBarPanelContext<R, P, T>): boolean
}
