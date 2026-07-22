import { DOMWrapper } from '@vue/test-utils'
import SkipTo from '../../../src/components/SkipTo.vue'
import { shallowMount } from '@ownclouders/web-test-helpers'

const selectors = {
  skipButton: '.skip-button'
}

describe('SkipTo component', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('should render provided text in the slot', () => {
    const targetElement = {
      setAttribute: vi.fn(),
      focus: vi.fn(),
      scrollIntoView: vi.fn()
    }
    vi.spyOn(document, 'getElementById').mockReturnValue(targetElement as unknown as HTMLElement)

    const { wrapper } = getShallowWrapper()
    const skipButton: DOMWrapper<Element> = wrapper.find(selectors.skipButton)

    expect(skipButton.text()).toEqual('Skip to main')
  })

  it('should look up the target element at click time and focus it', async () => {
    const targetElement = {
      setAttribute: vi.fn(),
      focus: vi.fn(),
      scrollIntoView: vi.fn()
    }
    const getElementByIdSpy = vi
      .spyOn(document, 'getElementById')
      .mockReturnValue(targetElement as unknown as HTMLElement)

    const { wrapper } = getShallowWrapper()

    // the target element is not looked up before the button is clicked
    expect(getElementByIdSpy).not.toHaveBeenCalled()

    await wrapper.find(selectors.skipButton).trigger('click')

    expect(getElementByIdSpy).toHaveBeenCalledWith('main-content')
    expect(targetElement.setAttribute).toHaveBeenCalledWith('tabindex', '-1')
    expect(targetElement.focus).toHaveBeenCalledTimes(1)
    expect(targetElement.scrollIntoView).toHaveBeenCalledTimes(1)
  })

  it('re-resolves the target on every click, picking up a target added after mount', async () => {
    const getElementByIdSpy = vi.spyOn(document, 'getElementById').mockReturnValue(null)

    const { wrapper } = getShallowWrapper()

    await wrapper.find(selectors.skipButton).trigger('click')
    expect(getElementByIdSpy).toHaveBeenCalledTimes(1)

    const targetElement = {
      setAttribute: vi.fn(),
      focus: vi.fn(),
      scrollIntoView: vi.fn()
    }
    getElementByIdSpy.mockReturnValue(targetElement as unknown as HTMLElement)

    await wrapper.find(selectors.skipButton).trigger('click')

    expect(getElementByIdSpy).toHaveBeenCalledTimes(2)
    expect(targetElement.focus).toHaveBeenCalledTimes(1)
  })

  it('should do nothing when the target element does not exist', async () => {
    vi.spyOn(document, 'getElementById').mockReturnValue(null)

    const { wrapper } = getShallowWrapper()

    await expect(wrapper.find(selectors.skipButton).trigger('click')).resolves.not.toThrow()
  })
})

function getShallowWrapper() {
  return {
    wrapper: shallowMount(SkipTo, {
      props: {
        target: 'main-content'
      },
      slots: {
        default: 'Skip to main'
      }
    })
  }
}
