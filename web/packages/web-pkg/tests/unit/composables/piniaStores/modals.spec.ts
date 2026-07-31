import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { useModals } from '../../../../src/composables/piniaStores'
import { createPinia, setActivePinia } from 'pinia'

describe('useModals', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  describe('method "dispatchModal"', () => {
    it('adds a modal to the stack of modals and sets it active', () => {
      getWrapper({
        setup: (instance) => {
          const data = { title: 'test' }
          const modal = instance.dispatchModal(data)

          expect(modal.id).toBeDefined()
          expect(modal.title).toEqual(data.title)
          expect(instance.activeModal).toEqual(modal)

          const modal2 = instance.dispatchModal(data)
          expect(instance.activeModal).toEqual(modal2)
        }
      })
    })
  })
  describe('method "updateModal"', () => {
    it('updates a modal with new data', () => {
      getWrapper({
        setup: (instance) => {
          const modal = instance.dispatchModal({ title: 'test' })
          const newTitle = 'new title'
          instance.updateModal(modal.id, 'title', newTitle)
          expect(instance.activeModal.title).toEqual(newTitle)
        }
      })
    })
  })
  describe('method "removeModal"', () => {
    it('removes an existing modal and sets another existing modal active', () => {
      getWrapper({
        setup: (instance) => {
          const modal = instance.dispatchModal({ title: 'test' })
          const modal2 = instance.dispatchModal({ title: 'test2' })

          expect(instance.modals.length).toBe(2)
          expect(instance.activeModal).toEqual(modal2)

          instance.removeModal(modal2.id)
          expect(instance.modals.length).toBe(1)
          expect(instance.activeModal).toEqual(modal)
        }
      })
    })
  })
  describe('method "setModalActive"', () => {
    it('moves a modal to the first position of the modal stack, making it active', () => {
      getWrapper({
        setup: (instance) => {
          const modal = instance.dispatchModal({ title: 'test' })
          const modal2 = instance.dispatchModal({ title: 'test2' })

          expect(instance.activeModal.id).toEqual(modal2.id)
          instance.setModalActive(modal.id)
          expect(instance.activeModal.id).toEqual(modal.id)
        }
      })
    })
  })
})

function getWrapper({ setup }: { setup: (instance: ReturnType<typeof useModals>) => void }) {
  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useModals()
        setup(instance)
      },
      { pluginOptions: { pinia: false } }
    )
  }
}
