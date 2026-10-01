import { App, Plugin, h } from 'vue'
import { createGettext } from 'vue3-gettext'
import DesignSystem from '../index'

/**
 * Test helpers for mounting components that render design system components (`<oc-button>`,
 * `<oc-icon>`, ...). Those are registered globally by this package's Vue plugin, so anything
 * that mounts them has to install it first.
 *
 * They live here rather than in a dedicated test-helper package because `design-system` sits at
 * the bottom of the workspace package graph: it cannot depend on a package that installs it
 * without forming a cycle - which is exactly what {design-system, web-pkg, web-test-helpers}
 * used to be.
 *
 * `@ownclouders/web-pkg/src/testing` builds on this and adds what its own components need
 * (pinia stores, casl abilities); `@ownclouders/web-test-helpers` re-exports that in turn.
 * Only `design-system`'s own specs should import from here - everything else uses
 * `@ownclouders/web-test-helpers`.
 */

export { mount, shallowMount, flushPromises, VueWrapper, RouterLinkStub } from '@vue/test-utils'
export type { RouteLocation } from 'vue-router'

type DefinedComponent = new (...args: any[]) => any
export type ComponentProps<T extends DefinedComponent> = InstanceType<T>['$props']
export type PartialComponentProps<T extends DefinedComponent> = Partial<ComponentProps<T>>

export interface DesignSystemPluginsOptions {
  designSystem?: boolean
  gettext?: boolean
  getTextDefaultLanguage?: string
}

/**
 * Stubs `<router-link>` so components can render links without a real router instance.
 */
export const routerLinkStubPlugin: Plugin = {
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
}

export const defaultPlugins = ({
  designSystem = true,
  gettext = true,
  getTextDefaultLanguage = 'en'
}: DesignSystemPluginsOptions = {}): Plugin[] => {
  const plugins: Plugin[] = []

  if (designSystem) {
    plugins.push(DesignSystem as unknown as Plugin)
  }

  if (gettext) {
    plugins.push(
      createGettext({ translations: {}, silent: true, defaultLanguage: getTextDefaultLanguage })
    )
  }

  plugins.push(routerLinkStubPlugin)

  return plugins
}
