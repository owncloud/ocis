import { ActionExtension, useEmbedMode } from '@ownclouders/web-pkg'
import QuickActions from '../../../../src/components/FilesList/QuickActions.vue'
import { defaultComponentMocks, defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { useExtensionRegistry } from '@ownclouders/web-pkg'
import { mock } from 'vitest-mock-extended'
import { ref } from 'vue'
import { Resource } from '@ownclouders/web-client'
import { quickActionsExtensionPoint } from '../../../../src/extensionPoints'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useEmbedMode: vi.fn()
}))

const collaboratorAction = {
  isVisible: vi.fn(() => true),
  handler: vi.fn(),
  icon: 'group-add',
  id: 'collaborators',
  name: 'show-shares',
  label: () => 'Add people'
}

const permanentLinkAction = {
  isVisible: vi.fn(() => false),
  handler: vi.fn(),
  icon: 'link-add',
  id: 'permanent-link',
  name: 'copy-permanent-link',
  label: () => 'Copy permanent link'
}

const testItem = {
  id: '1',
  icon: 'file',
  name: 'lorem.txt',
  path: '/lorem.txt',
  size: '12220',
  spaceId: '1'
} as Resource

describe('QuickActions', () => {
  describe('when multiple actions are provided', () => {
    const { wrapper } = getWrapper()

    it('should display all action buttons where "displayed" is set to true', () => {
      const actionButtons = wrapper.findAll('.oc-button')
      // there are two items provided as actions, with only one item set to display
      expect(actionButtons.length).toBe(1)

      const actionButton = actionButtons.at(0)
      const iconEl = actionButton.find('oc-icon-stub')

      expect(actionButton.exists()).toBeTruthy()
      expect(actionButton.attributes().class).toContain('files-quick-action-show-shares')
      expect(iconEl.exists()).toBeTruthy()
      expect(iconEl.attributes().name).toBe('group-add')
      expect(actionButton.attributes('aria-label')).toBe('Add people')
    })

    it('should not display action buttons where "displayed" is set to false', () => {
      const linkActionButton = wrapper.find('.files-quick-action-copy-permanent-link')

      expect(linkActionButton.exists()).toBeFalsy()
    })
  })

  describe('action handler', () => {
    it('should call action handler on click', async () => {
      const { wrapper } = getWrapper()
      const handlerAction = collaboratorAction.handler.mockImplementation(() => undefined)

      const actionButton = wrapper.find('.oc-button')
      await actionButton.trigger('click')
      expect(handlerAction).toHaveBeenCalledTimes(1)
    })
  })

  it('does not show actions in embed mode', () => {
    const { wrapper } = getWrapper({ embedModeEnabled: true })
    expect(wrapper.findAll('.oc-button').length).toBe(0)
  })
})

function getWrapper({ embedModeEnabled = false } = {}) {
  const plugins = defaultPlugins()

  vi.mocked(useEmbedMode).mockReturnValue(
    mock<ReturnType<typeof useEmbedMode>>({ isEnabled: ref(embedModeEnabled) })
  )

  const { requestExtensions } = useExtensionRegistry()
  vi.mocked(requestExtensions).mockReturnValue([
    mock<ActionExtension>({
      extensionPointIds: [quickActionsExtensionPoint.id],
      action: collaboratorAction
    }),
    mock<ActionExtension>({
      extensionPointIds: [quickActionsExtensionPoint.id],
      action: permanentLinkAction
    })
  ])

  return {
    wrapper: shallowMount(QuickActions, {
      props: {
        item: testItem
      },
      global: {
        stubs: { OcButton: false },
        mocks: { ...defaultComponentMocks() },
        plugins
      }
    })
  }
}
