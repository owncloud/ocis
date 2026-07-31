import {
  defaultComponentMocks,
  defaultPlugins,
  defaultStubs,
  mount
} from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import ContextActions from '../../../../src/components/FilesList/ContextActions.vue'

import {
  useFileActionsEnableSync,
  useFileActionsCopyPermanentLink,
  useFileActionsRename,
  useFileActionsCopy
} from '../../../../src/composables'
import { computed } from 'vue'
import { Action } from '../../../../src/composables/actions'

// vi.mock('../../../../src/composables/actions/files', async (importOriginal) => {
//   const original = await importOriginal()
//   return createMockActionComposables(importOriginal())
// })

describe.skip('ContextActions', () => {
  describe('menu sections', () => {
    it('do not render when no action enabled', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.findAll('action-menu-item-stub').length).toBe(0)
    })

    it('render enabled actions', () => {
      const enabledComposables = [
        useFileActionsEnableSync,
        useFileActionsCopyPermanentLink,
        useFileActionsRename,
        useFileActionsCopy
      ]
      for (const composable of enabledComposables) {
        vi.mocked(composable).mockImplementation(() => ({
          actions: computed(() => [mock<Action>({ isVisible: () => true })])
        }))
      }

      const { wrapper } = getWrapper()
      expect(wrapper.findAll('action-menu-item-stub').length).toBe(enabledComposables.length)
    })
  })
})

function getWrapper() {
  const mocks = {
    ...defaultComponentMocks()
  }
  return {
    mocks,
    wrapper: mount(ContextActions, {
      props: {
        actionOptions: {
          space: mock<SpaceResource>(),
          resources: [mock<Resource>()]
        }
      },
      global: {
        mocks,
        provide: { ...mocks, currentSpace: mock<SpaceResource>() },
        stubs: { ...defaultStubs, 'action-menu-item': true },
        plugins: [...defaultPlugins()]
      }
    })
  }
}
