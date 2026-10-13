import { useDisableRootOverscroll } from '../../../../src/composables/rootOverscroll'
import { getComposableWrapper } from '../../../../src/testing'

describe('useDisableRootOverscroll', () => {
  const getRootOverscroll = () =>
    document.documentElement.style.getPropertyValue('overscroll-behavior-y')

  afterEach(() => {
    document.documentElement.style.removeProperty('overscroll-behavior-y')
  })

  it('disables the vertical overscroll of the page root while mounted', () => {
    const wrapper = getComposableWrapper(() => useDisableRootOverscroll())
    expect(getRootOverscroll()).toBe('none')
    wrapper.unmount()
  })

  it('removes the setting again when unmounted', () => {
    const wrapper = getComposableWrapper(() => useDisableRootOverscroll())
    wrapper.unmount()
    expect(getRootOverscroll()).toBe('')
  })

  it('restores a previously set value when unmounted', () => {
    document.documentElement.style.setProperty('overscroll-behavior-y', 'contain')
    const wrapper = getComposableWrapper(() => useDisableRootOverscroll())
    expect(getRootOverscroll()).toBe('none')
    wrapper.unmount()
    expect(getRootOverscroll()).toBe('contain')
  })
})
