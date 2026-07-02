import { mount } from '@ownclouders/web-test-helpers'
import OcIcon from './OcIcon.vue'
import { glob } from 'glob'
import path from 'path'

describe('OcIcon', () => {
  test('render and update', async () => {
    const icons = glob.sync(path.resolve(__dirname, '..', 'assets', 'icons', '*.svg')).map((p) => ({
      name: path.basename(p, '.svg'),
      path: path.relative(__dirname, p)
    }))

    const wrapper = mount(OcIcon)
    await wrapper.vm.$nextTick()
    const inlineSvg = wrapper.findComponent({ name: 'inline-svg' })

    for (const icon of icons) {
      await wrapper.setProps({ name: icon.name })
      await inlineSvg.emitted('unloaded')
      await inlineSvg.emitted('loaded')
      expect(wrapper.find('text').text()).toBe(icon.name)
    }
  })
})
