import { Extension, ExtensionPoint, useExtensionPreferencesStore } from '../../../../../src'
import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { createPinia, setActivePinia } from 'pinia'
import { mock } from 'vitest-mock-extended'

describe('useExtensionPreferencesStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })
  afterEach(() => {
    localStorage.clear()
  })

  it('has no preferences initially', () => {
    getWrapper({
      setup: (instance) => {
        expect(Object.keys(instance.extensionPreferences).length).toBe(0)
      }
    })
  })

  describe('set and get selected extensions', () => {
    it('can set and get selected extensions', () => {
      getWrapper({
        setup: (instance) => {
          const extensionPointId = 'extension-point-id'
          const extensionId = 'extension-id'
          instance.setSelectedExtensionIds(extensionPointId, [extensionId])
          const extensionPreference = instance.getExtensionPreference(extensionPointId, [
            'fallback'
          ])
          expect(extensionPreference.selectedExtensionIds).toEqual([extensionId])
        }
      })
    })
    it('keeps preferences unique by extension point id', () => {
      getWrapper({
        setup: (instance) => {
          const extensionPointId = 'extension-point-id'
          instance.setSelectedExtensionIds(extensionPointId, ['foo-1'])
          instance.setSelectedExtensionIds(extensionPointId, ['foo-2'])
          instance.setSelectedExtensionIds(extensionPointId, ['foo-3'])
          const extensionPreference = instance.getExtensionPreference(extensionPointId, [
            'fallback'
          ])
          expect(extensionPreference.selectedExtensionIds).toEqual(['foo-3'])
        }
      })
    })
    it('uses the provided fallback value for unknown extension points', () => {
      getWrapper({
        setup: (instance) => {
          const extensionPointId = 'extension-point-id'
          const fallbackId = 'fallback'
          const extensionPreference = instance.getExtensionPreference(extensionPointId, [
            fallbackId
          ])
          expect(extensionPreference.selectedExtensionIds).toEqual([fallbackId])
        }
      })
    })
  })

  describe('extractDefaultExtensionIds', () => {
    describe('extension points with multiple=true', () => {
      it('returns all provided extension ids', () => {
        getWrapper({
          setup: (instance) => {
            const extensionPoint = mock<ExtensionPoint<Extension>>({
              id: 'extension-point-id',
              multiple: true
            })
            const extensionIds = ['foo-1', 'foo-2']
            const extensionMocks = extensionIds.map((id) => mock<Extension>({ id }))
            expect(instance.extractDefaultExtensionIds(extensionPoint, extensionMocks)).toEqual(
              extensionIds
            )
          }
        })
      })
    })
    describe('extension points with multiple=false', () => {
      it('returns the default extension id from the extension point', () => {
        getWrapper({
          setup: (instance) => {
            const extensionPoint = mock<ExtensionPoint<Extension>>({
              id: 'extension-point-id',
              multiple: false,
              defaultExtensionId: 'foo-1'
            })
            expect(instance.extractDefaultExtensionIds(extensionPoint, [])).toEqual(['foo-1'])
          }
        })
      })
      it('returns an empty array if the extension point has no default extension id', () => {
        getWrapper({
          setup: (instance) => {
            const extensionPoint = mock<ExtensionPoint<Extension>>({
              id: 'extension-point-id',
              multiple: false,
              defaultExtensionId: undefined
            })
            expect(instance.extractDefaultExtensionIds(extensionPoint, [])).toEqual([])
          }
        })
      })
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (instance: ReturnType<typeof useExtensionPreferencesStore>) => void
}) {
  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useExtensionPreferencesStore()
        setup(instance)
      },
      { pluginOptions: { pinia: false } }
    )
  }
}
