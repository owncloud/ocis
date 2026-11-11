import Drop from './OcDrop.vue'
import { getSizeClass } from '../../helpers'
import { shallowMount, mount } from '../../testing'

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

  describe('focus on close via item click', () => {
    const mountCloseOnClick = () => {
      document.body.innerHTML = ''
      const wrapper = mount(
        {
          template: `<div>
            <button id="trigger">trigger</button>
            <input id="elsewhere" />
            <oc-drop toggle="#trigger" close-on-click>
              <button id="item">item</button>
              <button id="moves-focus" @click="moveFocus">moves focus</button>
            </oc-drop>
          </div>`,
          components: { 'oc-drop': Drop },
          methods: {
            moveFocus() {
              document.getElementById('elsewhere').focus()
            }
          }
        },
        { attachTo: document.body }
      )
      const open = async () => {
        await wrapper.find('#trigger').trigger('click')
        await new Promise((resolve) => setTimeout(resolve, 100))
      }
      return { wrapper, open }
    }

    it('returns focus to the trigger so it does not fall back to the page', async () => {
      const { open } = mountCloseOnClick()
      await open()
      const item = document.getElementById('item')
      item.focus()
      item.click()

      expect(document.activeElement).toBe(document.getElementById('trigger'))
    })

    it('keeps focus where the item handler moved it', async () => {
      const { open } = mountCloseOnClick()
      await open()
      const item = document.getElementById('moves-focus')
      item.focus()
      item.click()

      expect(document.activeElement).toBe(document.getElementById('elsewhere'))
    })

    it("hides the trigger's own tooltip instead of leaving it open when refocused", async () => {
      const { open } = mountCloseOnClick()
      await open()
      const trigger = document.getElementById('trigger') as HTMLElement & {
        tooltip?: { hide: () => void }
      }
      const hide = vi.fn()
      trigger.tooltip = { hide }

      const item = document.getElementById('item')
      item.focus()
      item.click()

      expect(hide).toHaveBeenCalled()
    })
  })

  describe('focus leaving the drop', () => {
    const mountWithFocusableItem = () => {
      document.body.innerHTML = ''
      const wrapper = mount(
        {
          template: `<div>
            <button id="trigger">trigger</button>
            <input id="elsewhere" />
            <oc-drop toggle="#trigger">
              <button id="item">item</button>
            </oc-drop>
          </div>`,
          components: { 'oc-drop': Drop }
        },
        { attachTo: document.body }
      )
      const tippy = wrapper.findComponent({ name: 'oc-drop' }).vm.tippyInstance
      const wait = () => new Promise((resolve) => setTimeout(resolve, 100))
      const open = async () => {
        await wrapper.find('#trigger').trigger('click')
        await wait()
      }
      return { wrapper, tippy, open, wait }
    }

    it('hides when focus moves outside of the drop', async () => {
      const { tippy, open } = mountWithFocusableItem()
      await open()
      document.getElementById('item').focus()
      document.getElementById('elsewhere').focus()

      expect(tippy.state.isVisible).toBe(false)
    })

    it('closes instead of reopening when the trigger is clicked while focus is inside', async () => {
      const { wrapper, tippy, open, wait } = mountWithFocusableItem()
      await open()
      document.getElementById('item').focus()

      // clicking the trigger focuses it before the click toggles the drop
      document.getElementById('trigger').focus()
      await wrapper.find('#trigger').trigger('click')
      await wait()

      expect(tippy.state.isVisible).toBe(false)
    })
  })
})
