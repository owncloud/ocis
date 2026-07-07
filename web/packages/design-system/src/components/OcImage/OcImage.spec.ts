import OcImage from './OcImage.vue'
import { mount } from '@ownclouders/web-test-helpers'

// @vitest-environment jsdom
describe('OcImage', () => {
  function getWrapper(props = {}) {
    return mount(OcImage, {
      props: {
        src: 'http://someimage.jpg',
        ...props
      }
    })
  }
  it('should set the provided src', () => {
    const wrapper = getWrapper()
    expect(wrapper.attributes('src')).toBe('http://someimage.jpg')
  })
  it('should set the provided title for image', () => {
    const wrapper = getWrapper({ title: 'test title' })
    expect(wrapper.attributes('title')).toBe('test title')
  })
  it.each(['eager', 'lazy'])('should set the provided loading type for image', (loadingType) => {
    const wrapper = getWrapper({ loadingType: loadingType })
    expect(wrapper.attributes('loading')).toBe(loadingType)
  })
  describe('when alt is set', () => {
    const wrapper = getWrapper({ alt: 'test alt text' })
    it('should set the provided alt for image', () => {
      expect(wrapper.attributes('alt')).toBe('test alt text')
    })
    it('should set aria hidden property to "false"', () => {
      expect(wrapper.attributes('aria-hidden')).toBe('false')
    })
  })
  describe('when alt is not set', () => {
    it('should set aria hidden property to "true"', () => {
      const wrapper = getWrapper()
      expect(wrapper.attributes('aria-hidden')).toBe('true')
    })
  })
})
