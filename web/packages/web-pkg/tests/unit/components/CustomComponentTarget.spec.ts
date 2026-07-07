import CustomComponentTarget from '../../../src/components/CustomComponentTarget.vue'
import { defaultComponentMocks, defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import {
  CustomComponentExtension,
  Extension,
  ExtensionPoint,
  ExtensionPreferenceItem,
  useExtensionPreferencesStore,
  useExtensionRegistry
} from '../../../src'
import { mock } from 'vitest-mock-extended'
import { h } from 'vue'

const selectors = {
  target: '[data-testid="custom-component-target"]'
}

describe('CustomComponentTarget', () => {
  const mockExtensionPointSingle = mock<ExtensionPoint<Extension>>({
    id: 'dummy-extension-point-single',
    extensionType: 'customComponent',
    multiple: false
  })
  const mockExtensionPointMulti = mock<ExtensionPoint<Extension>>({
    id: 'dummy-extension-point-multi',
    extensionType: 'customComponent',
    multiple: true
  })

  describe('no extensions match the extension point', () => {
    it.each([mockExtensionPointSingle, mockExtensionPointMulti])(
      'renders nothing',
      (extensionPoint) => {
        const { wrapper } = getWrapper({
          extensionPoint,
          extensions: []
        })
        expect(wrapper.find(selectors.target).exists()).toBeFalsy()
      }
    )
  })

  describe('exactly 1 extension matches the extension point', () => {
    it.each([mockExtensionPointSingle, mockExtensionPointMulti])(
      'renders 1 component',
      (extensionPoint) => {
        const extensionId = 'custom-1'
        const { wrapper } = getWrapper({
          extensionPoint,
          extensions: [createExtensionMock(extensionId, extensionPoint.id)]
        })
        expect(wrapper.find(selectors.target).exists()).toBeTruthy()
      }
    )
  })

  describe('multiple extensions match the extension point', () => {
    describe('extension point allows only 1 extension', () => {
      it('renders 1 component, respecting the user preference', () => {
        const extensionMocks = ['custom-1', 'custom-2'].map((id) =>
          createExtensionMock(id, mockExtensionPointSingle.id)
        )
        const { wrapper } = getWrapper({
          extensionPoint: mockExtensionPointSingle,
          extensions: extensionMocks,
          preference: mock<ExtensionPreferenceItem>({
            extensionPointId: mockExtensionPointSingle.id,
            selectedExtensionIds: [extensionMocks[1].id]
          })
        })
        expect(wrapper.findAll(selectors.target).length).toBe(1)
      })
    })

    describe('multiple extensions match the extension point', () => {
      it('renders n components', () => {
        const extensionMocks = ['custom-1', 'custom-2'].map((id) =>
          createExtensionMock(id, mockExtensionPointMulti.id)
        )
        const { wrapper } = getWrapper({
          extensionPoint: mockExtensionPointMulti,
          extensions: extensionMocks
        })
        expect(wrapper.findAll(selectors.target).length).toBe(extensionMocks.length)
      })
    })
  })
})

function getWrapper({
  extensionPoint,
  extensions,
  preference
}: {
  extensionPoint: ExtensionPoint<Extension>
  extensions: Extension[]
  preference?: ExtensionPreferenceItem
}) {
  const plugins = defaultPlugins()

  const { getExtensionPreference } = useExtensionPreferencesStore()
  vi.mocked(getExtensionPreference).mockReturnValue(preference)

  const { requestExtensions } = useExtensionRegistry()
  vi.mocked(requestExtensions).mockReturnValue(extensions)

  const mocks = defaultComponentMocks()

  return {
    mocks,
    wrapper: mount(CustomComponentTarget, {
      props: {
        extensionPoint
      },
      global: {
        plugins,
        mocks,
        provide: mocks,
        stubs: { OcCheckbox: true }
      }
    })
  }
}

function createExtensionMock(id: string, extensionPointId: string) {
  return mock<CustomComponentExtension>({
    id,
    type: 'customComponent',
    extensionPointIds: [extensionPointId],
    content: () => [
      h('p', {
        innerHTML: `hello from ${id}`,
        'data-testid': 'custom-component-target'
      })
    ]
  })
}
