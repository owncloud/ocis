import DesignSystem from '@ownclouders/design-system'
import { createGettext } from 'vue3-gettext'
import { App, Plugin, h } from 'vue'
import { abilitiesPlugin } from '@casl/vue'
import { createMongoAbility } from '@casl/ability'
import { AbilityRule } from '@ownclouders/web-client'
import { PiniaMockOptions, createMockStore } from './mocks'

export interface DefaultPluginsOptions {
  abilities?: AbilityRule[]
  designSystem?: boolean
  gettext?: boolean
  pinia?: boolean
  piniaOptions?: PiniaMockOptions
  getTextDefaultLanguage?: string
}

export const defaultPlugins = ({
  abilities = [],
  designSystem = true,
  gettext = true,
  pinia = true,
  piniaOptions = {},
  getTextDefaultLanguage = 'en'
}: DefaultPluginsOptions = {}): Plugin[] => {
  const plugins = []

  plugins.push({
    install(app: App) {
      app.use(abilitiesPlugin, createMongoAbility(abilities))
    }
  })

  if (designSystem) {
    plugins.push(DesignSystem as unknown as Plugin)
  }

  if (gettext) {
    plugins.push(
      createGettext({ translations: {}, silent: true, defaultLanguage: getTextDefaultLanguage })
    )
  } else {
    plugins.push({
      install(app: App) {
        // mock `v-translate` directive
        app.directive('translate', {
          mounted: () => undefined
        })
      }
    })
  }

  if (pinia) {
    plugins.push(createMockStore(piniaOptions))
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
