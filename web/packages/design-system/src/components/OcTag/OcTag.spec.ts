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
    it('emits no inline style at all when no colour index is given', () => {
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

    it('fills itself from the theme colour the index points at', () => {
      const wrapper = shallowMount(Tag, {
        props: {
          type: 'button',
          colorIndex: 3
        }
      })

      const style = wrapper.attributes('style')
      expect(style).toContain('background-color: var(--oc-color-tag-3)')
      expect(style).toContain('color: var(--oc-color-tag-3-text,currentColor)')
    })

    it('takes the label colour with the fill, so no caller works out contrast', () => {
      // The `-text` half of the pair is derived once per theme by the theme store, from the
      // theme's own text colours. A chip only ever names a slot.
      const wrapper = shallowMount(Tag, { props: { colorIndex: 7 } })

      expect(wrapper.attributes('style')).toContain(
        'color: var(--oc-color-tag-7-text,currentColor)'
      )
    })

    it('colours itself independently of other props', () => {
      const wrapper = shallowMount(Tag, {
        props: {
          type: 'span',
          size: 'large',
          rounded: true,
          colorIndex: 5
        }
      })

      const classes = wrapper.classes()
      expect(classes).toContain('oc-tag-l')
      expect(classes).toContain('oc-tag-rounded')
      expect(wrapper.attributes('style')).toContain('var(--oc-color-tag-5)')
    })

    it('stays unfilled for the -1 a theme without tag colours yields', () => {
      // `useTagColor` returns -1 when the theme ships no `tagColorsList`, and a chip must then
      // look exactly like the plain badge rather than reference a var that resolves to nothing.
      const wrapper = shallowMount(Tag, { props: { colorIndex: -1 } })

      expect(wrapper.attributes('style')).toBeUndefined()
      expect(wrapper.classes()).not.toContain('oc-tag-filled')
    })

    it('marks a filled tag, so nested links and icons take the chip colour', () => {
      // Setting `color` on the chip is not enough on its own: a nested `<a>` carries the global
      // anchor colour (styles/theme/oc-text.scss:21), `.oc-button-raw` paints its own text and
      // icon, and `.oc-tag .oc-icon > svg` pins a muted fill. The class is what the stylesheet
      // hangs the `inherit`/`currentColor` overrides on.
      const wrapper = shallowMount(Tag, { props: { colorIndex: 3 } })

      expect(wrapper.classes()).toContain('oc-tag-filled')
    })

    it('does not mark an unfilled tag, leaving the generic badge untouched', () => {
      expect(shallowMount(Tag).classes()).not.toContain('oc-tag-filled')
    })
  })
})
