import OcLogo from './OcLogo.vue'
import { shallowMount, mount } from '@ownclouders/web-test-helpers'
import OcImage from '../OcImage/OcImage.vue'

describe('OcLogo', () => {
  const requiredProps = {
    src: 'http://some-image-link.jpg',
    alt: 'test alt'
  }
  function getWrapper(props = {}) {
    return shallowMount(OcLogo, {
      props: {
        ...requiredProps,
        ...props
      },
      global: {
        stubs: { 'oc-img': true }
      }
    })
  }
  const wrapper = getWrapper()
  it('should set the provided src to image element', () => {
    const imageElement = wrapper.find('oc-img-stub')
    expect(imageElement.attributes('src')).toBe('http://some-image-link.jpg')
  })
  it('should set the provided alt to image element', () => {
    const imageElement = wrapper.find('oc-img-stub')
    expect(imageElement.attributes('alt')).toBe('test alt')
  })
  it('should add provided class to the image element', () => {
    const component = {
      template: "<oc-logo src='image-link' alt='imageText' class='test-class'></oc-logo>",
      components: { OcLogo },
      name: 'TestOcLogo'
    }
    const wrapper = mount(component, { global: { stubs: { 'oc-img': OcImage } } })
    expect(wrapper.attributes('class')).toContain('test-class')
  })
})
