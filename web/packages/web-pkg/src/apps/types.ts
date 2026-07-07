import { App, ComponentCustomProperties, Ref } from 'vue'
import { RouteLocationRaw, Router, RouteRecordRaw } from 'vue-router'
import { Extension, ExtensionPoint } from '../composables/piniaStores'
import { IconFillType } from '../helpers'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { Translations } from 'vue3-gettext'
import { FileActionOptions } from '../composables/actions/types'

export interface AppReadyHookArgs {
  globalProperties: ComponentCustomProperties & Record<string, any>
  router: Router
  instance?: App
  portal?: any
}

export interface AppNavigationItem {
  isActive?: () => boolean
  activeFor?: { name?: string; path?: string }[]
  isVisible?: () => boolean
  fillType?: IconFillType
  icon?: string
  name: string | (() => string)
  route?: RouteLocationRaw
  handler?: () => void
  priority?: number
}

export type AppConfigObject = Record<string, any>

export interface ApplicationFileExtension {
  app?: string
  extension?: string
  createFileHandler?: (arg: {
    fileName: string
    space: SpaceResource
    currentFolder: Resource
  }) => Promise<Resource>
  hasPriority?: boolean
  label?: string | (() => string)
  name?: string
  icon?: string
  mimeType?: string
  newFileMenu?: {
    menuTitle: () => string
    isVisible?: ({ currentFolder }: { currentFolder: Resource }) => boolean
  }
  routeName?: string
  secureView?: boolean
  customHandler?: (
    fileActionOptions: FileActionOptions,
    extension: string,
    appFileExtension: ApplicationFileExtension
  ) => Promise<void> | void
}

/** ApplicationInformation describes required information of an application */
export interface ApplicationInformation {
  color?: string
  id?: string
  name?: string
  icon?: string
  iconFillType?: IconFillType
  iconColor?: string
  img?: string
  meta?: {
    fileSizeLimit?: number
  }
  extensions?: ApplicationFileExtension[]
  defaultExtension?: string
  translations?: Translations
  /** Asserts whether the app has any route which works as an editor */
  hasEditor?: boolean
}

/**
 * ApplicationTranslations is a map of language keys to translations
 */
export interface ApplicationTranslations {
  [lang: string]: {
    [key: string]: string | string[]
  }
}

/** ClassicApplicationScript reflects classic application script structure */
export interface ClassicApplicationScript {
  appInfo?: Omit<ApplicationInformation, 'hasEditor'>
  routes?: ((args: ComponentCustomProperties) => RouteRecordRaw[]) | RouteRecordRaw[]
  navItems?: ((args: ComponentCustomProperties) => AppNavigationItem[]) | AppNavigationItem[]
  translations?: Translations
  extensions?: Ref<Extension[]>
  extensionPoints?: Ref<ExtensionPoint<any>[]>
  initialize?: () => void
  ready?: (args: AppReadyHookArgs) => Promise<void> | void
  mounted?: (args: AppReadyHookArgs) => void
  // TODO: move this to its own type
  setup?: (args: { applicationConfig: AppConfigObject }) => ClassicApplicationScript
}

export type ApplicationSetupOptions = { applicationConfig: AppConfigObject }

export const defineWebApplication = (args: {
  setup: (options: ApplicationSetupOptions) => ClassicApplicationScript
}) => {
  return args
}
