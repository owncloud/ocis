import { Key, Modifier, useKeyboardActions } from '../../../../src/composables/keyboardActions'
import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { ref } from 'vue'

describe('useKeyboardActions', () => {
  it('should be valid', () => {
    expect(useKeyboardActions).toBeDefined()
  })

  it('should bind keys', () => {
    const wrapper = getWrapper()
    const { keyboardActions } = wrapper.vm

    keyboardActions.bindKeyAction({ primary: Key.A }, () => undefined)
    expect(keyboardActions.actions.value.length).toBe(1)

    wrapper.unmount()
  })

  it('should be possible remove keys', () => {
    const wrapper = getWrapper()
    const { keyboardActions } = wrapper.vm

    const keyActionIndex = keyboardActions.bindKeyAction({ primary: Key.A }, () => undefined)

    expect(keyboardActions.actions.value.length).toBe(1)

    keyboardActions.removeKeyAction(keyActionIndex)
    expect(keyboardActions.actions.value.length).toBe(0)

    wrapper.unmount()
  })

  it('should be possible execute callback on key event', () => {
    const wrapper = getWrapper()
    const { keyboardActions } = wrapper.vm
    const counter = ref(0)

    const increment = () => {
      counter.value += 1
    }

    // primary key
    keyboardActions.bindKeyAction({ primary: Key.A }, increment)

    const event = new KeyboardEvent('keydown', { key: 'a' })
    document.dispatchEvent(event)

    expect(counter.value).toBe(1)

    // primary key + modifier
    keyboardActions.bindKeyAction({ modifier: Modifier.Ctrl, primary: Key.A }, increment)

    const eventWithModifier = new KeyboardEvent('keydown', { key: 'a', ctrlKey: true })
    document.dispatchEvent(eventWithModifier)

    expect(counter.value).toBe(2)

    wrapper.unmount()
  })

  it('should not execute callback on key event if disallowed modifier is present', () => {
    const wrapper = getWrapper()
    const { keyboardActions } = wrapper.vm
    const counter = ref(0)

    const increment = () => {
      counter.value += 1
    }

    keyboardActions.bindKeyAction({ modifier: Modifier.Ctrl, primary: Key.A }, increment)

    const eventWithModifier = new KeyboardEvent('keydown', {
      key: 'a',
      ctrlKey: true,
      shiftKey: true
    })
    document.dispatchEvent(eventWithModifier)

    expect(counter.value).toBe(0)

    wrapper.unmount()
  })

  it('should not execute callback when checkbox input is focused', () => {
    const wrapper = getWrapper()
    const { keyboardActions } = wrapper.vm
    const counter = ref(0)

    keyboardActions.bindKeyAction({ primary: Key.Space }, () => {
      counter.value += 1
    })

    const checkbox = document.createElement('input')
    checkbox.type = 'checkbox'
    document.body.appendChild(checkbox)
    checkbox.focus()

    const event = new KeyboardEvent('keydown', { key: ' ' })
    document.dispatchEvent(event)

    expect(counter.value).toBe(0)

    checkbox.remove()
    wrapper.unmount()
  })

  it('should execute modifier callback when checkbox input is focused', () => {
    const wrapper = getWrapper()
    const { keyboardActions } = wrapper.vm
    const counter = ref(0)

    keyboardActions.bindKeyAction({ modifier: Modifier.Ctrl, primary: Key.C }, () => {
      counter.value += 1
    })

    const checkbox = document.createElement('input')
    checkbox.type = 'checkbox'
    document.body.appendChild(checkbox)
    checkbox.focus()

    const event = new KeyboardEvent('keydown', { key: 'c', ctrlKey: true })
    document.dispatchEvent(event)

    expect(counter.value).toBe(1)

    checkbox.remove()
    wrapper.unmount()
  })
})

function getWrapper() {
  return getComposableWrapper(() => {
    const keyboardActions = useKeyboardActions()
    return { keyboardActions }
  })
}
