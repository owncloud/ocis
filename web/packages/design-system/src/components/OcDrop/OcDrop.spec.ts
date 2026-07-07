import { shallowMount, mount } from '@ownclouders/web-test-helpers'
import Drop from './OcDrop.vue'
import { getSizeClass } from '../../helpers'

const dom = ({ position = 'auto', mode = 'click', paddingSize = 'medium' } = {}) => {
  document.body.innerHTML = ''
  const wrapper = mount(
    {
      template:
        '<div><p id="trigger">trigger</p><oc-drop :position="position" :mode="mode" :padding-size="paddingSize" toggle="#trigger">show</oc-drop></div>',
      components: { 'oc-drop': Drop }
    },
    {
      attachTo: document.body,
      data: () => ({ position, mode, paddingSize })
    }
  )
  const drop = wrapper.findComponent({ name: 'oc-drop' })
  const tippy = drop.vm.tippyInstance

  return { wrapper, drop, tippy }
}

describe('OcDrop', () => {
  it('handles dropId prop', () => {
    for (let i = 0; i < 5; i++) {
      const wrapper = shallowMount(Drop)
      expect(wrapper.attributes().id).toBe(`oc-drop-${i + 1}`)
    }

    for (let i = 0; i < 5; i++) {
      const id = `custom-drop-id-${i}`
      const wrapper = shallowMount(Drop, {
        props: {
          dropId: id
        }
      })
      expect(wrapper.attributes().id).toBe(id)
    }
  })

  it.each(['xsmall', 'small', 'medium', 'large', 'xlarge', 'xxlarge', 'remove'])(
    'handles padding size prop for value %s',
    (size) => {
      const { drop } = dom({ paddingSize: size })
      expect(drop.html().includes(`oc-p-${getSizeClass(size)}`)).toBeTruthy()
    }
  )

  describe('tippy', () => {
    it('inits tippy', () => {
      const { wrapper, drop, tippy } = dom()

      expect(tippy).toBeTruthy()
      expect(tippy.reference).toBe(wrapper.find('#trigger').element)
      expect(tippy.props.content).toBe(drop.vm.$refs.drop)
    })

    it('updates tippy', async () => {
      const { wrapper, tippy } = dom()

      await wrapper.setData({
        position: 'left',
        mode: 'hover'
      })

      expect(tippy.props.placement).toBe('left')
      expect(tippy.props.trigger).toBe('mouseenter focus')
    })

    it('renders tippy', async () => {
      const { wrapper } = dom()
      const trigger = wrapper.find('#trigger')
      const wait = async () => {
        await wrapper.vm.$nextTick()
        return new Promise((resolve) => setTimeout(resolve, 100))
      }

      await trigger.trigger('click') // show
      await wait()
      expect(wrapper.findComponent(Drop).exists()).toBeTruthy()
      expect(trigger.attributes()['aria-expanded']).toBe('true')
      expect(wrapper.element).toMatchSnapshot()

      await trigger.trigger('click') // hide
      await wait()
      expect(trigger.attributes()['aria-expanded']).toBe('false')
      expect(wrapper.element).toMatchSnapshot()

      await wrapper.setData({
        mode: 'hover'
      })

      await trigger.trigger('mouseenter') // show
      await wait()
      expect(trigger.attributes()['aria-expanded']).toBe('true')
      expect(wrapper.element).toMatchSnapshot()
    })
  })
})
