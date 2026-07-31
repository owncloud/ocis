import focusMixin from '../../../src/mixins/focusMixin'
import { defineComponent } from 'vue'
import { mount } from '@ownclouders/web-test-helpers'

const Component = defineComponent({
  name: 'DummyComponent',
  mixins: [focusMixin],
  template: `
          <ul>
            <li v-for="index in 10" :key="index">
              <a v-bind:id="'item-' + index" tabindex="0">{{index}}</a>
            </li>
          </ul>
        `
})

function getWrapper() {
  return {
    wrapper: mount(Component, { attachTo: document.body })
  }
}

const { wrapper } = getWrapper()
const wrapperComponent = wrapper.findComponent({ name: 'DummyComponent' })
const item1 = wrapper.get('#item-1')
const item2 = wrapper.get('#item-2')
const item3 = wrapper.get('#item-3')
const item4 = wrapper.get('#item-4')
const item5 = wrapper.get('#item-5')
const item6 = wrapper.get('#item-6')
const item7 = wrapper.get('#item-7')
const item8 = wrapper.get('#item-8')
const item9 = wrapper.get('#item-9')
const item10 = wrapper.get('#item-10')
const focus = wrapperComponent.vm?.focus

describe('focusMixin', () => {
  // trap -----------------------
  // #item-1  || ---  x ||  3 <--
  // #item-2  || -->  4 ||  x ---
  // #item-3  || -->  1 ||  5 <--
  // #item-4  || -->  6 ||  2 <--
  // #item-5  || -->  3 ||  7 <--
  // #item-6  || -->  8 ||  4 <--
  // #item-7  || -->  5 ||  9 <--
  // #item-8  || --> 10 ||  6 <--
  // #item-9  || -->  7 || 10 <--
  // #item-10 || -->  9 ||  8 <--
  it('records and replays focus events', () => {
    focus({ from: item2.element, to: item4.element })
    expect(document.activeElement.id).toBe(item4.element.id)

    focus({ to: item6.element })
    expect(document.activeElement.id).toBe(item6.element.id)

    focus({ to: item8.element })
    expect(document.activeElement.id).toBe(item8.element.id)

    focus({ to: item10.element })
    expect(document.activeElement.id).toBe(item10.element.id)

    focus({ to: item9.element })
    expect(document.activeElement.id).toBe(item9.element.id)

    focus({ to: item7.element })
    expect(document.activeElement.id).toBe(item7.element.id)

    focus({ to: item5.element })
    expect(document.activeElement.id).toBe(item5.element.id)

    focus({ to: item3.element })
    expect(document.activeElement.id).toBe(item3.element.id)

    focus({ to: item1.element })
    expect(document.activeElement.id).toBe(item1.element.id)

    focus({ revert: true })
    expect(document.activeElement.id).toBe(item3.element.id)

    focus({ revert: true })
    expect(document.activeElement.id).toBe(item5.element.id)

    focus({ revert: true })
    expect(document.activeElement.id).toBe(item7.element.id)

    focus({ revert: true })
    expect(document.activeElement.id).toBe(item9.element.id)

    focus({ revert: true })
    expect(document.activeElement.id).toBe(item10.element.id)

    focus({ revert: true })
    expect(document.activeElement.id).toBe(item8.element.id)

    focus({ revert: true })
    expect(document.activeElement.id).toBe(item6.element.id)

    focus({ revert: true })
    expect(document.activeElement.id).toBe(item4.element.id)

    focus({ revert: true })
    expect(document.activeElement.id).toBe(item2.element.id)
  })

  // trap ----------------------
  // #item-2  || -->  4 || x ---
  // #item-4  || -->  6 || x ---
  // #item-6  || ---  x || 4 <--
  // restart trap --------------
  // #item-1  || -->  8 || x ---
  // #item-8  || --> 10 || 1 <--
  // #item-10 || ---  x || 8 <--
  it('can be restarted', () => {
    focus({ from: item2.element, to: item4.element })
    expect(document.activeElement.id).toBe(item4.element.id)

    focus({ to: item6.element })
    expect(document.activeElement.id).toBe(item6.element.id)

    focus({ revert: true })
    expect(document.activeElement.id).toBe(item4.element.id)

    focus({ from: item1.element, to: item8.element })
    expect(document.activeElement.id).toBe(item8.element.id)

    focus({ to: item10.element })
    expect(document.activeElement.id).toBe(item10.element.id)

    focus({ revert: true })
    expect(document.activeElement.id).toBe(item8.element.id)

    focus({ revert: true })
    expect(document.activeElement.id).toBe(item1.element.id)
  })
})
