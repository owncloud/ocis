import ActionMenuItem from '../../../../src/components/ContextActions/ActionMenuItem.vue'
import { FileAction } from '../../../../src/composables/actions'
import { defaultPlugins, mount, shallowMount } from '@ownclouders/web-test-helpers'

const selectors = {
  handler: '[data-testid="action-handler"]',
  icon: '[data-testid="action-icon"]',
  img: '[data-testid="action-img"]',
  label: '[data-testid="action-label"]',
  srHint: '[data-testid="action-sr-hint"]',
  ariaLabel: '[aria-label="foo"]'
}

const fileActions = {
  download: {
    name: 'download-file',
    icon: 'file-download',
    handler: vi.fn(),
    label: () => 'Download',
    class: 'oc-files-actions-download-file-trigger'
  } as unknown as FileAction,
  navigate: {
    name: 'navigate',
    icon: 'folder-open',
    route: () => ({ name: 'files-personal' }),
    label: () => 'Open Folder',
    class: 'oc-files-actions-navigate-trigger'
  } as unknown as FileAction,
  redirect: {
    name: 'redirect',
    icon: 'external',
    href: () => 'https://owncloud.com',
    label: () => 'Visit ownCloud',
    class: 'oc-files-actions-redirect-trigger'
  } as unknown as FileAction
}

describe('ActionMenuItem component', () => {
  it('renders an icon if there is one defined in the action', () => {
    const action = fileActions.download
    const { wrapper } = getWrapper(action)
    expect(wrapper.find(selectors.icon).exists()).toBeTruthy()
    expect(wrapper.find(selectors.icon).attributes().name).toBe(action.icon)
  })
  it('renders an image if there is one defined in the action', () => {
    const action = { ...fileActions.download, img: 'https://owncloud.tld/img.png' }
    const { wrapper } = getWrapper(action)
    expect(wrapper.find(selectors.icon).exists()).toBeFalsy()
    expect(wrapper.find(selectors.img).exists()).toBeTruthy()
    expect(wrapper.find(selectors.img).attributes().src).toBe(action.img)
  })
  it('renders the action label', () => {
    const action = fileActions.download
    const { wrapper } = getWrapper(action)
    expect(wrapper.find(selectors.label).exists()).toBeTruthy()
    expect(wrapper.find(selectors.label).text()).toBe(action.label())
  })
  it('renders a tooltip for a disabled action', () => {
    const action = { ...fileActions.download, disabledTooltip: () => 'Foo', isDisabled: () => true }
    const { wrapper } = getWrapper(action)

    expect(wrapper.find(`.${fileActions.download.class}`).attributes().arialabel).toBe(
      action.disabledTooltip()
    )
    expect(wrapper.find(`.${fileActions.download.class}`).attributes().disabled).toBeTruthy()
  })
  describe('component is of type oc-button', () => {
    it('calls the action handler on button click', async () => {
      const action = fileActions.download
      const spyHandler = action.handler
      const { wrapper } = getWrapper(action, mount)
      const button = wrapper.find(selectors.handler)
      expect(button.exists()).toBeTruthy()
      expect(button.element.tagName).toBe('BUTTON')
      await button.trigger('click')
      expect(spyHandler).toHaveBeenCalled()
    })
  })
  describe('component is of type router-link', () => {
    it('has a link', () => {
      const action = fileActions.navigate
      const { wrapper } = getWrapper(action, mount)
      const link = wrapper.find(selectors.handler)
      expect(link.exists()).toBeTruthy()
      expect(link.element.tagName).toBe('ROUTER-LINK-STUB')
      // FIXME: to.name cannot be accessed, because the attributes().to holds a string containing `[object Object]`.
      // That doesn't allow checking the name of the router-link.
      // expect(link.attributes().href).toBe(action.route)
    })
  })
  describe('component is of type a', () => {
    it('has a href', () => {
      const action = fileActions.redirect
      const { wrapper } = getWrapper(action, mount)
      const link = wrapper.find(selectors.handler)
      expect(link.exists()).toBeTruthy()
      expect(link.element.tagName).toBe('A')
      expect(link.attributes().href).toBe(action.href())
    })
  })
})

function getWrapper(action: FileAction, mountType = shallowMount) {
  return {
    wrapper: mountType(ActionMenuItem, {
      props: {
        action,
        actionOptions: {},
        items: []
      },
      global: {
        renderStubDefaultSlot: true,
        stubs: {
          'router-link': true
        },
        plugins: [...defaultPlugins()]
      }
    })
  }
}
