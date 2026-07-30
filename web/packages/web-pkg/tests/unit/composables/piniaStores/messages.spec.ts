import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { useMessages } from '../../../../src/composables/piniaStores'
import { createPinia, setActivePinia } from 'pinia'

describe('useMessages', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('method "triggerLatestAction"', () => {
    it('returns false when there are no messages', () => {
      getWrapper({
        setup: (instance) => {
          expect(instance.triggerLatestAction()).toBe(false)
        }
      })
    })

    it('returns false when no message has an action', () => {
      getWrapper({
        setup: (instance) => {
          instance.showMessage({ title: 'plain message' })
          expect(instance.triggerLatestAction()).toBe(false)
        }
      })
    })

    it('runs the action of the most recent message that has one and dismisses it', () => {
      getWrapper({
        setup: (instance) => {
          const onClick = vi.fn()
          instance.showMessage({ title: 'plain message' })
          const withAction = instance.showMessage({
            title: 'undoable message',
            actions: [{ label: 'Undo', onClick }]
          })

          const triggered = instance.triggerLatestAction()

          expect(triggered).toBe(true)
          expect(onClick).toHaveBeenCalledTimes(1)
          expect(instance.messages.find((m) => m.id === withAction.id)).toBeUndefined()
        }
      })
    })

    it('skips messages without an action in favor of an earlier one that has one', () => {
      getWrapper({
        setup: (instance) => {
          const onClick = vi.fn()
          const withAction = instance.showMessage({
            title: 'undoable message',
            actions: [{ label: 'Undo', onClick }]
          })
          instance.showMessage({ title: 'later plain message' })

          const triggered = instance.triggerLatestAction()

          expect(triggered).toBe(true)
          expect(onClick).toHaveBeenCalledTimes(1)
          expect(instance.messages.find((m) => m.id === withAction.id)).toBeUndefined()
        }
      })
    })
  })
})

function getWrapper({ setup }: { setup: (instance: ReturnType<typeof useMessages>) => void }) {
  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useMessages()
        setup(instance)
      },
      { pluginOptions: { pinia: false } }
    )
  }
}
