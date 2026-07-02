import CompareSaveDialog from '../../../../src/components/SideBar/CompareSaveDialog.vue'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'

describe('CompareSaveDialog', () => {
  describe('computed method "unsavedChanges"', () => {
    it('should be false if objects are equal', () => {
      const { wrapper } = getWrapper({
        propsData: {
          originalObject: { id: '1', displayName: 'jan' },
          compareObject: { id: '1', displayName: 'jan' }
        }
      })
      expect((wrapper.vm as any).unsavedChanges).toBeFalsy()
    })

    it('should be true if objects are not equal', () => {
      const { wrapper } = getWrapper({
        propsData: {
          originalObject: { id: '1', displayName: 'jan' },
          compareObject: { id: '1', displayName: 'janina' }
        }
      })
      expect((wrapper.vm as any).unsavedChanges).toBeTruthy()
    })
  })
})

function getWrapper({ propsData = {} } = {}) {
  return {
    wrapper: shallowMount(CompareSaveDialog, {
      props: {
        originalObject: {},
        compareObject: {},
        ...propsData
      },
      global: {
        plugins: [...defaultPlugins()]
      }
    })
  }
}
