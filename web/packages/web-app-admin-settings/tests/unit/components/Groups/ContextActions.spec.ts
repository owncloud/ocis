import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Resource } from '@ownclouders/web-client'
import ContextActions from '../../../../src/components/Groups/ContextActions.vue'
import { useGroupActionsDelete, useGroupActionsEdit } from '../../../../src/composables/actions'
import { computed, ref } from 'vue'
import { Action } from '@ownclouders/web-pkg'

function createMockActionComposables(module: Record<string, unknown>) {
  const mockModule: Record<string, any> = {}
  for (const m of Object.keys(module)) {
    mockModule[m] = vi.fn(() => ({ actions: ref([]) }))
  }
  return mockModule
}

vi.mock('@ownclouders/web-pkg', async (importOriginal) => {
  const original = await importOriginal<any>()
  return createMockActionComposables(original)
})

vi.mock(
  'web-app-admin-settings/src/composables/actions/groups/useGroupActionsDelete',
  async (importOriginal) => {
    const original = await importOriginal<any>()
    return createMockActionComposables(original)
  }
)
vi.mock(
  'web-app-admin-settings/src/composables/actions/groups/useGroupActionsEdit',
  async (importOriginal) => {
    const original = await importOriginal<any>()
    return createMockActionComposables(original)
  }
)

const selectors = {
  actionMenuItemStub: 'action-menu-item-stub'
}

describe.skip('ContextActions', () => {
  describe('menu sections', () => {
    it('do not render when no action enabled', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.findAll(selectors.actionMenuItemStub).length).toBe(0)
    })

    it('render enabled actions', () => {
      const enabledComposables = [useGroupActionsDelete, useGroupActionsEdit]
      vi.mocked(useGroupActionsDelete).mockImplementation(() => ({
        actions: computed(() => [mock<Action>({ isVisible: () => true })]),
        deleteGroups: null
      }))
      vi.mocked(useGroupActionsEdit).mockImplementation(() => ({
        actions: computed(() => [mock<Action>({ isVisible: () => true })])
      }))
      const { wrapper } = getWrapper()
      expect(wrapper.findAll(selectors.actionMenuItemStub).length).toBe(enabledComposables.length)
    })
  })
})

function getWrapper() {
  return {
    wrapper: mount(ContextActions, {
      props: {
        actionOptions: {
          resources: [mock<Resource>()]
        }
      },
      global: {
        stubs: { 'action-menu-item': true },
        plugins: [...defaultPlugins()]
      }
    })
  }
}
