import { createTestingPinia, getComposableWrapper } from '@ownclouders/web-test-helpers'
import { RawConfig, useAppsStore, useConfigStore } from '../../../../src/composables/piniaStores'

describe('useConfigStore', () => {
  beforeEach(() => {
    createTestingPinia({ stubActions: false })
  })

  it('has initial options as defaults', () => {
    getWrapper({
      setup: (instance) => {
        expect(Object.keys(instance.options).length).toBeGreaterThan(0)
      }
    })
  })

  describe('method "loadConfig"', () => {
    it('sets given options and overwrites defaults', () => {
      getWrapper({
        setup: (instance) => {
          expect(instance.options.contextHelpersReadMore).toBeTruthy()

          const data = { options: { contextHelpersReadMore: false } } as RawConfig
          instance.loadConfig(data)

          expect(instance.options.contextHelpersReadMore).toBeFalsy()
        }
      })
    })
    it('loads config for external apps', () => {
      getWrapper({
        setup: (instance) => {
          const externalApp = { id: '1', path: '/foo', config: { foo: 'bar' } }
          const data = {
            server: 'https://foo.bar',
            theme: undefined,
            options: { contextHelpersReadMore: false },
            external_apps: [externalApp]
          } as RawConfig

          instance.loadConfig(data)
          const appsStore = useAppsStore()

          expect(appsStore.loadExternalAppConfig).toHaveBeenCalledWith({
            appId: externalApp.id,
            config: externalApp.config
          })
        }
      })
    })
  })

  describe('serverUrl', () => {
    it('defaults to "window.location.origin" if no server url given', () => {
      getWrapper({
        setup: (instance) => {
          const data = { theme: '' } as RawConfig
          instance.loadConfig(data)

          expect(instance.serverUrl).toEqual(`${window.location.origin}/`)
        }
      })
    })
  })
})

function getWrapper({ setup }: { setup: (instance: ReturnType<typeof useConfigStore>) => void }) {
  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useConfigStore()
        setup(instance)
      },
      { pluginOptions: { pinia: false } }
    )
  }
}
