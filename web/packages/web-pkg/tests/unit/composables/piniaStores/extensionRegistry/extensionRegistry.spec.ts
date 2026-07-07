import {
  ActionExtension,
  CustomComponentExtension,
  Extension,
  ExtensionPoint,
  SidebarNavExtension,
  SidebarPanelExtension,
  useExtensionRegistry
} from '../../../../../src'
import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { createPinia, setActivePinia } from 'pinia'
import { computed, ref, unref } from 'vue'
import { mock } from 'vitest-mock-extended'

describe('useExtensionRegistry', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('register and request extensions', () => {
    describe('querying extensions throws an error', () => {
      it('if the extensionPoint has no id', () => {
        getWrapper({
          setup: (instance) => {
            const extensionPointMock: ExtensionPoint<ActionExtension> = mock<
              ExtensionPoint<ActionExtension>
            >({
              id: '',
              extensionType: 'action'
            })
            expect(() => instance.requestExtensions(extensionPointMock)).toThrowError()
          }
        })
      })
      it('if the extensionPoint has no extensionType', () => {
        getWrapper({
          setup: (instance) => {
            const extensionPointMock: ExtensionPoint<any> = mock<ExtensionPoint<any>>({
              id: 'some.unique.id',
              extensionType: ''
            })
            expect(() => instance.requestExtensions(extensionPointMock)).toThrowError()
          }
        })
      })
    })
    describe('querying extensions has an empty result', () => {
      it('if no extensions are registered', () => {
        getWrapper({
          setup: (instance) => {
            const extensionPointMock: ExtensionPoint<CustomComponentExtension> = mock<
              ExtensionPoint<CustomComponentExtension>
            >({
              id: 'some.unique.id',
              extensionType: 'customComponent'
            })
            const result = instance.requestExtensions(extensionPointMock)
            expect(result.length).toBe(0)
          }
        })
      })
      it('if no matching extensions are found', () => {
        const extensionPoint: ExtensionPoint<CustomComponentExtension> = mock<
          ExtensionPoint<CustomComponentExtension>
        >({
          id: 'extension-point-id',
          extensionType: 'customComponent'
        })
        const extensions = computed<Extension[]>(() =>
          ['foo-1', 'foo-2'].map((id) =>
            mock<Extension>({
              id,
              type: 'customComponent',
              extensionPointIds: [extensionPoint.id]
            })
          )
        )

        getWrapper({
          setup: (instance) => {
            instance.registerExtensions(extensions)

            const extensionPointWrongType: ExtensionPoint<SidebarPanelExtension<any, any, any>> =
              mock<ExtensionPoint<SidebarPanelExtension<any, any, any>>>({
                id: 'extension-point-id',
                extensionType: 'sidebarPanel'
              })
            const result1 = instance.requestExtensions(extensionPointWrongType)
            expect(result1.length).toBe(0)

            const extensionPointWrongId: ExtensionPoint<CustomComponentExtension> = mock<
              ExtensionPoint<CustomComponentExtension>
            >({
              id: 'some-other-extension-point-id',
              extensionType: 'customComponent'
            })
            const result2 = instance.requestExtensions(extensionPointWrongId)
            expect(result2.length).toBe(0)
          }
        })
      })
    })

    it('can query extensions by extension point', () => {
      const extensionPoint: ExtensionPoint<CustomComponentExtension> = mock<
        ExtensionPoint<CustomComponentExtension>
      >({
        id: 'extension-point-id',
        extensionType: 'customComponent'
      })
      const extensionIds = ['foo-1', 'foo-2', 'foo-3']
      const extensions = computed(() => [
        ...extensionIds.map((id) =>
          mock<Extension>({
            id,
            type: 'customComponent',
            extensionPointIds: [extensionPoint.id]
          })
        ),
        mock<Extension>({
          id: 'foo-4',
          type: 'customComponent',
          extensionPointIds: ['some-other-extension-point-id']
        })
      ])

      getWrapper({
        setup: (instance) => {
          instance.registerExtensions(extensions)

          const result = instance.requestExtensions(extensionPoint)
          expect(result.map((e) => e.id)).toEqual(extensionIds)
        }
      })
    })

    it('unregisters extensions', () => {
      const extensionPoint: ExtensionPoint<CustomComponentExtension> = mock<
        ExtensionPoint<CustomComponentExtension>
      >({
        id: 'extension-point-id',
        extensionType: 'customComponent'
      })

      const extension = mock<Extension>({
        id: 'foo-1',
        type: 'customComponent',
        extensionPointIds: [extensionPoint.id]
      })
      const extensions = computed(() => [extension])

      getWrapper({
        setup: (instance) => {
          instance.registerExtensions(extensions)

          const result1 = instance.requestExtensions(extensionPoint)
          expect(result1.length).toBe(1)

          instance.unregisterExtensions([extension.id])

          const result2 = instance.requestExtensions(extensionPoint)
          expect(result2.length).toBe(0)
        }
      })
    })
  })

  describe('register and get extensionPoints', () => {
    describe('querying extension points has an empty result', () => {
      it('if no extension points are registered', () => {
        getWrapper({
          setup: (instance) => {
            const result = instance.getExtensionPoints()
            expect(result.length).toBe(0)
          }
        })
      })

      it('if no matching extension points are found', () => {
        const extensionPoint = mock<ExtensionPoint<Extension>>({
          id: 'foo-1',
          extensionType: 'customComponent'
        })
        const extensionPoints = computed<ExtensionPoint<Extension>[]>(() => [extensionPoint])

        getWrapper({
          setup: (instance) => {
            instance.registerExtensionPoints(extensionPoints)

            const result = instance.getExtensionPoints({ extensionType: 'customComponent' })
            expect(result.length).toBe(1)
          }
        })
      })
    })

    it('can query extension points by type', () => {
      const extensionPoints = computed<ExtensionPoint<Extension>[]>(() => {
        return [
          mock<ExtensionPoint<CustomComponentExtension>>({
            id: 'foo-1',
            extensionType: 'customComponent'
          }),
          mock<ExtensionPoint<SidebarPanelExtension<any, any, any>>>({
            id: 'foo-2',
            extensionType: 'sidebarPanel'
          }),
          mock<ExtensionPoint<CustomComponentExtension>>({
            id: ' foo-3',
            extensionType: 'customComponent'
          })
        ]
      })

      getWrapper({
        setup: (instance) => {
          instance.registerExtensionPoints(extensionPoints)

          const result1 = instance.getExtensionPoints({ extensionType: 'customComponent' })
          expect(result1.map((ep) => ep.id)).toEqual([
            unref(extensionPoints)[0].id,
            unref(extensionPoints)[2].id
          ])

          const result2 = instance.getExtensionPoints({ extensionType: 'sidebarPanel' })
          expect(result2.map((ep) => ep.id)).toEqual([unref(extensionPoints)[1].id])
        }
      })
    })

    it('unregisters extension points', () => {
      const extensionPoint = mock<ExtensionPoint<CustomComponentExtension>>({
        id: 'foo-1',
        extensionType: 'customComponent'
      })

      const extensionPoints = computed<ExtensionPoint<Extension>[]>(() => [extensionPoint])

      getWrapper({
        setup: (instance) => {
          instance.registerExtensionPoints(extensionPoints)

          const result1 = instance.getExtensionPoints({ extensionType: 'customComponent' })
          expect(result1.length).toBe(1)

          instance.unregisterExtensionPoints([extensionPoint.id])

          const result2 = instance.getExtensionPoints({ extensionType: 'customComponent' })
          expect(result2.length).toBe(0)
        }
      })
    })
  })

  describe('rebuild', () => {
    it('returns a non-vault route when scope is not vault', () => {
      const sidebarExtension = mock<SidebarNavExtension>({
        id: 'sidebar-non-vault',
        type: 'sidebarNav',
        navItem: {
          name: 'Files',
          route: 'files'
        }
      })

      getWrapper({
        setup: (instance) => {
          instance.registerExtensions(ref<Extension[]>([sidebarExtension]))

          instance.rebuild({ route: ref({ params: { scope: 'projects' } }) })

          const rebuiltExtension = unref(unref(instance.extensions)[0])[0] as SidebarNavExtension
          expect(rebuiltExtension.navItem.route).toBe('/files')
        }
      })
    })

    it('returns a vault-prefixed route when scope is vault', () => {
      const sidebarExtension = mock<SidebarNavExtension>({
        id: 'sidebar-vault',
        type: 'sidebarNav',
        navItem: {
          name: 'Files',
          route: 'files'
        }
      })

      getWrapper({
        setup: (instance) => {
          instance.registerExtensions(ref<Extension[]>([sidebarExtension]))

          instance.rebuild({ route: ref({ params: { scope: 'vault' } }) })

          const rebuiltExtension = unref(unref(instance.extensions)[0])[0] as SidebarNavExtension
          expect(rebuiltExtension.navItem.route).toBe('/vault/files')
        }
      })
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (instance: ReturnType<typeof useExtensionRegistry>) => void
}) {
  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useExtensionRegistry()
        setup(instance)
      },
      { pluginOptions: { pinia: false } }
    )
  }
}
