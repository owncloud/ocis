import { ref, nextTick } from 'vue'
import { useIsVisible } from './index'
import { mount } from '@ownclouders/web-test-helpers'

const mockIntersectionObserver = () => {
  const enable = () => {
    const mock = {
      observe: vi.fn(),
      disconnect: vi.fn(),
      unobserve: vi.fn()
    }

    window.IntersectionObserver = vi.fn().mockImplementation(function () {
      return mock
    })

    return {
      mock,
      callback: (args: unknown[], fastForward = 0) => {
        ;(window.IntersectionObserver as any).mock.calls[0][0](args)
        vi.advanceTimersByTime(fastForward)
      }
    }
  }

  const disable = () => {
    delete window.IntersectionObserver
  }

  return { enable, disable }
}

const createWrapper = (options = {}) =>
  mount({
    template: `
      <div>
      <div ref="target">{{ isVisible }}</div>
      </div>`,
    setup: () => {
      const target = ref<HTMLElement>()
      const { isVisible } = useIsVisible({ ...options, target })

      return {
        isVisible,
        target
      }
    }
  })

describe('useIsVisible', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  const { enable: enableIntersectionObserver, disable: disableIntersectionObserver } =
    mockIntersectionObserver()

  it('is visible by default if browser does not support IntersectionObserver', () => {
    disableIntersectionObserver()
    const wrapper = createWrapper()
    expect((wrapper.vm.$refs.target as any).innerHTML).toBe('true')
  })

  it('observes the target', async () => {
    const { mock: observerMock } = enableIntersectionObserver()
    createWrapper()
    await nextTick()

    expect(observerMock.observe).toHaveBeenCalledTimes(1)
  })

  it('only shows once and then gets unobserved if the the composable is in the default show mode', async () => {
    const { mock: observerMock, callback: observerCallback } = enableIntersectionObserver()
    const wrapper = createWrapper()

    await nextTick()
    expect((wrapper.vm.$refs.target as any).innerHTML).toBe('false')

    observerCallback([{ isIntersecting: true }])
    await nextTick()
    expect((wrapper.vm.$refs.target as any).innerHTML).toBe('true')
    expect(observerMock.unobserve).toHaveBeenCalledTimes(1)
  })

  it('shows and hides multiple times if the the composable is in showHide mode', async () => {
    const { mock: observerMock, callback: observerCallback } = enableIntersectionObserver()
    const wrapper = createWrapper({ mode: 'showHide' })

    await nextTick()
    expect((wrapper.vm.$refs.target as any).innerHTML).toBe('false')

    observerCallback([{ isIntersecting: true }])
    await nextTick()
    expect((wrapper.vm.$refs.target as any).innerHTML).toBe('true')
    expect(observerMock.unobserve).toHaveBeenCalledTimes(0)
  })

  it('disconnects the observer before component gets unmounted', () => {
    const { mock: observerMock } = enableIntersectionObserver()
    const wrapper = createWrapper()

    expect(observerMock.disconnect).toHaveBeenCalledTimes(0)
    wrapper.unmount()
    expect(observerMock.disconnect).toHaveBeenCalledTimes(1)
  })
})
