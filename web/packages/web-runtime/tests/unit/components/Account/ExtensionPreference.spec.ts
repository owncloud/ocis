import ExtensionPreference from '../../../../src/components/Account/ExtensionPreference.vue'
import {
  defaultComponentMocks,
  defaultPlugins,
  getOcSelectOptions,
  mount
} from '@ownclouders/web-test-helpers'
import {
  Extension,
  ExtensionPoint,
  useExtensionPreferencesStore,
  useExtensionRegistry
} from '@ownclouders/web-pkg'
import { mock } from 'vitest-mock-extended'

const selectors = {
  dropdown: '.extension-preference'
}

describe('ExtensionPreference component', () => {
  afterEach(() => {
    localStorage.clear()
  })

  it('renders a dropdown for an extension point', () => {
    const extensionPoint = mock<ExtensionPoint<Extension>>({
      id: 'test-extension-point',
      multiple: false
    })
    const { wrapper } = getWrapper({ extensionPoint })
    expect(wrapper.find('.v-select').exists()).toBeTruthy()
  })

  describe('extensionPoint with multiple=false', () => {
    const extensionPoint = mock<ExtensionPoint<Extension>>({
      id: 'test-extension-point',
      multiple: false,
      defaultExtensionId: 'foo-2'
    })
    const extensions = [
      mock<Extension>({
        id: 'foo-1',
        userPreference: {
          optionLabel: 'Foo 1'
        }
      }),
      mock<Extension>({
        id: 'foo-2',
        userPreference: {
          optionLabel: 'Foo 2'
        }
      })
    ]
    it('renders extensions as dropdown options', async () => {
      const { wrapper } = getWrapper({ extensionPoint, extensions })
      const options = await getOcSelectOptions(wrapper, selectors.dropdown)
      expect(options.length).toBe(2)
      extensions.forEach((extension) => {
        expect(options.some((option) => option.text() === extension.userPreference.optionLabel))
          .toBeTruthy
      })
    })
    it('renders the default extension first in the options list', async () => {
      const { wrapper } = getWrapper({ extensionPoint, extensions })
      const options = await getOcSelectOptions(wrapper, selectors.dropdown)
      const defaultExtension = extensions.find((e) => e.id === extensionPoint.defaultExtensionId)
      expect(options[0].text()).toBe(defaultExtension.userPreference.optionLabel)
    })
    it('selecting an extension updates the extension preference store', async () => {
      const { wrapper } = getWrapper({ extensionPoint, extensions })
      const options = await getOcSelectOptions(wrapper, selectors.dropdown)
      await options[1].trigger('click')
      const preferences = useExtensionPreferencesStore()
      const preference = preferences.getExtensionPreference(extensionPoint.id, [
        extensionPoint.defaultExtensionId
      ])
      const clickedExtension = extensions.find(
        (extension) => extension.userPreference.optionLabel === options[1].text()
      )
      expect(preference.selectedExtensionIds).toEqual([clickedExtension.id])
    })
  })

  // FIXME: Add tests for extensionPoint with multiple=true as soon as that's fully supported.
  describe('extensionPoint with multiple=true', () => {
    it.todo('renders extensions as dropdown options')
    it.todo('selecting extensions updates the extension preference store')
  })
})

function getWrapper({
  extensionPoint,
  extensions = []
}: {
  extensionPoint: ExtensionPoint<Extension>
  extensions?: Extension[]
}) {
  const plugins = defaultPlugins({
    piniaOptions: {
      stubActions: false
    }
  })

  const { requestExtensions } = useExtensionRegistry()
  vi.mocked(requestExtensions).mockReturnValue(extensions)

  const mocks = {
    ...defaultComponentMocks()
  }

  return {
    mocks,
    wrapper: mount(ExtensionPreference, {
      props: {
        extensionPoint
      },
      global: {
        plugins,
        mocks,
        provide: mocks,
        stubs: { VueSelect: false }
      }
    })
  }
}
