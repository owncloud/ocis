import { App, Plugin } from 'vue'
import { abilitiesPlugin } from '@casl/vue'
import { createMongoAbility } from '@casl/ability'
import { AbilityRule } from '@ownclouders/web-client'
import {
  defaultPlugins as designSystemPlugins,
  DesignSystemPluginsOptions
} from '@ownclouders/design-system/testing'
import { createMockStore, PiniaMockOptions } from './pinia'

// Extends the design system's mount plugins (design system itself + gettext + a `<router-link>`
// stub) with what web-pkg's own components need on top: casl abilities and mocked pinia stores.
// The design system half is defined once, in `@ownclouders/design-system/testing`, because that
// package needs it for its own specs and cannot depend on this one.
export interface DefaultPluginsOptions extends DesignSystemPluginsOptions {
  abilities?: AbilityRule[]
  pinia?: boolean
  piniaOptions?: PiniaMockOptions
}

export const defaultPlugins = ({
  abilities = [],
  pinia = true,
  piniaOptions = {},
  ...designSystemOptions
}: DefaultPluginsOptions = {}): Plugin[] => [
  {
    install(app: App) {
      app.use(abilitiesPlugin, createMongoAbility(abilities))
    }
  },
  ...designSystemPlugins(designSystemOptions),
  ...(pinia ? [createMockStore(piniaOptions)] : [])
]
