import { PartialComponentProps, defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Resource } from '@ownclouders/web-client'
import BatchActions from '../../../src/components/BatchActions.vue'
import { Action, ActionMenuItem } from '../../../src'

const selectors = {
  actionMenuItemStub: 'action-menu-item-stub',
  batchActionsSquashed: '.oc-appbar-batch-actions-squashed'
}

describe('BatchActions', () => {
  describe('menu sections', () => {
    it('do not render when no action enabled', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.findAll(selectors.actionMenuItemStub).length).toBe(0)
    })

    it('render enabled actions', () => {
      const actions = [{} as Action]
      const { wrapper } = getWrapper({ props: { actions } })
      expect(wrapper.findAll(selectors.actionMenuItemStub).length).toBe(actions.length)
    })
  })
  describe('limited screen space', () => {
    it('adds the squashed-class when limited screen space is available', () => {
      const { wrapper } = getWrapper({ props: { limitedScreenSpace: true } })
      expect(wrapper.find(selectors.batchActionsSquashed).exists()).toBeTruthy()
    })
    it('correctly tells the action item component to show tooltips when limited screen space is available', () => {
      const { wrapper } = getWrapper({
        props: { actions: [{} as Action], limitedScreenSpace: true }
      })
      expect(
        wrapper.findComponent<typeof ActionMenuItem>(selectors.actionMenuItemStub).props()
          .showTooltip
      ).toBeTruthy()
    })
  })
})

function getWrapper(
  { props }: { props?: PartialComponentProps<typeof BatchActions> } = { props: {} }
) {
  return {
    wrapper: mount(BatchActions, {
      props: {
        items: [mock<Resource>()],
        actions: [],
        actionOptions: {},
        ...props
      },
      global: {
        stubs: { 'action-menu-item': true },
        plugins: [...defaultPlugins()]
      }
    })
  }
}
