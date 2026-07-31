import OcDropzone from './OcDropzone.vue'
import { mount } from '@ownclouders/web-test-helpers'

describe('OcDropzone', () => {
  const selectors = {
    dropzone: '.oc-dropzone'
  }
  describe('when slot is provided', () => {
    it('should render slot html', () => {
      const wrapper = mount(OcDropzone, {
        slots: {
          default: "<p class='test-class'>Drag and drop to upload content into current folder</p>"
        }
      })
      const slotElement = wrapper.find(`${selectors.dropzone} p`)
      expect(slotElement.exists()).toBeTruthy()
      expect(slotElement.attributes('class')).toBe('test-class')
      expect(slotElement.text()).toBe('Drag and drop to upload content into current folder')
    })
  })
})
