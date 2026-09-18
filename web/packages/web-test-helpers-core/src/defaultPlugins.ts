import { createGettext } from 'vue3-gettext'
import { App, Plugin, h } from 'vue'
import { abilitiesPlugin } from '@casl/vue'
import { createMongoAbility } from '@casl/ability'
import { AbilityRule } from '@ownclouders/web-client'
import { getDesignSystemPlugin, getMockStoreFactory } from './registry'

// Generic over the pinia mock-store options type so this package never has to name
// `PiniaMockOptions` (defined in @ownclouders/web-pkg) - see registry.ts for why.
export interface DefaultPluginsOptions<TPiniaOptions = Record<string, unknown>> {
  abilities?: AbilityRule[]
  designSystem?: boolean
  gettext?: boolean
  pinia?: boolean
  piniaOptions?: TPiniaOptions
  getTextDefaultLanguage?: string
}

export const defaultPlugins = <TPiniaOptions = Record<string, unknown>>({
  abilities = [],
  designSystem = true,
  gettext = true,
  pinia = true,
  piniaOptions = {} as TPiniaOptions,
  getTextDefaultLanguage = 'en'
}: DefaultPluginsOptions<TPiniaOptions> = {}): Plugin[] => {
  const plugins = []

  plugins.push({
    install(app: App) {
      app.use(abilitiesPlugin, createMongoAbility(abilities))
    }
  })

  if (designSystem) {
    plugins.push(getDesignSystemPlugin())
  }

  if (gettext) {
    plugins.push(
      createGettext({ translations: {}, silent: true, defaultLanguage: getTextDefaultLanguage })
    )
  }

  if (pinia) {
    plugins.push(getMockStoreFactory()(piniaOptions))
  }

  plugins.push({
    install(app: App) {
      app.component('RouterLink', {
        name: 'RouterLink',
        props: {
          tag: { type: String, default: 'a' },
          to: { type: [String, Object], default: '' }
        },
        setup(props) {
          let path = props.to

          if (!!path && typeof path !== 'string') {
            path = props.to.path || props.to.name

            if (props.to.params) {
              path += '/' + Object.values(props.to.params).join('/')
            }

            if (props.to.query) {
              path += '?' + Object.values(props.to.query).join('&')
            }
          }

          return () => h(props.tag, { attrs: { href: path } })
        }
      })
    }
  })

  return plugins
}
