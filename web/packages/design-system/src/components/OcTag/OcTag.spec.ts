import { shallowMount } from '@ownclouders/web-test-helpers'
import Tag from './OcTag.vue'

describe('OcTag', () => {
  it('uses correct component when type is specified', () => {
    const wrapper = shallowMount(Tag, {
      props: {
        type: 'button'
      }
    })

    expect(wrapper.element.tagName.toLowerCase()).toMatch('button')
    expect(wrapper.html()).toMatchSnapshot()
  })

  it('emits a click event', async () => {
    const wrapper = shallowMount(Tag, {
      props: {
        type: 'a'
      }
    })

    await wrapper.trigger('click')
    expect(wrapper.emitted().click).toBeTruthy()
  })

  describe('Tag chip fill color', () => {
    it('emits no inline style at all when fillColor is not set', () => {
      // OcTag is a generic badge — it also renders the app-store "most recent" pill, the
      // primary/active instance badges and the tag overflow counter. Fill must stay opt-in so
      // none of those pick up a background.
      const wrapper = shallowMount(Tag, {
        props: {
          type: 'button'
        }
      })

      expect(wrapper.attributes('style')).toBeUndefined()
    })

    it('applies fillColor to rendered output when fillColor is set', () => {
      const wrapper = shallowMount(Tag, {
        props: {
          type: 'button',
          fillColor: 'var(--oc-color-tag-3)'
        }
      })

      const element = wrapper.element as HTMLElement
      const style = element.getAttribute('style')
      expect(style).toBeTruthy()
      expect(style).toContain('var(--oc-color-tag-3)')
    })

    it('applies fillColor independently of other props', () => {
      const wrapper = shallowMount(Tag, {
        props: {
          type: 'span',
          size: 'large',
          rounded: true,
          fillColor: 'var(--oc-color-tag-5)'
        }
      })

      const classes = wrapper.classes()
      expect(classes).toContain('oc-tag-l')
      expect(classes).toContain('oc-tag-rounded')
      expect(wrapper.html()).toContain('var(--oc-color-tag-5)')
    })

    it('applies labelColor as the text colour alongside the fill', () => {
      const wrapper = shallowMount(Tag, {
        props: {
          fillColor: 'var(--oc-color-tag-3)',
          labelColor: '#000000'
        }
      })

      const style = wrapper.attributes('style')
      expect(style).toContain('background-color: var(--oc-color-tag-3)')
      expect(style).toContain('color: #000000')
    })

    it('leaves the text colour inherited when labelColor is not given', () => {
      // A fill without a label colour must not guess. `--oc-color-text-inverse` is white in the
      // light theme and black in the dark one, which is backwards for pale light-theme fills.
      const wrapper = shallowMount(Tag, {
        props: {
          fillColor: 'var(--oc-color-tag-3)'
        }
      })

      expect(wrapper.attributes('style')).not.toContain('color: var(--oc-color-text-inverse)')
    })

    it('marks a filled tag, so nested links and icons take the chip colour', () => {
      // Setting `color` on the chip is not enough on its own: a nested `<a>` carries the global
      // anchor colour (styles/theme/oc-text.scss:21), `.oc-button-raw` paints its own text and
      // icon, and `.oc-tag .oc-icon > svg` pins a muted fill. The class is what the stylesheet
      // hangs the `inherit`/`currentColor` overrides on.
      const wrapper = shallowMount(Tag, {
        props: { fillColor: 'var(--oc-color-tag-3)', labelColor: '#ffffff' }
      })

      expect(wrapper.classes()).toContain('oc-tag-filled')
    })

    it('does not mark an unfilled tag, leaving the generic badge untouched', () => {
      expect(shallowMount(Tag).classes()).not.toContain('oc-tag-filled')
    })
  })
})
