import { DOMWrapper } from '@vue/test-utils'
import SkipTo from '../../../src/components/SkipTo.vue'
import { shallowMount } from '@ownclouders/web-test-helpers'
;(document as any).getElementById = vi.fn(() => ({
  setAttribute: vi.fn(),
  focus: vi.fn(),
  scrollIntoView: vi.fn()
}))

const selectors = {
  skipButton: '.skip-button'
}

describe('SkipTo component', () => {
  const spySkipToTarget = vi.spyOn(SkipTo.methods, 'skipToTarget')

  let wrapper: ReturnType<typeof getShallowWrapper>['wrapper']
  let skipButton: DOMWrapper<Element>
  beforeEach(() => {
    wrapper = getShallowWrapper().wrapper
    skipButton = wrapper.find(selectors.skipButton)
  })

  it('should render provided text in the slot', () => {
    expect(skipButton.text()).toEqual('Skip to main')
  })
  it('should call "skipToTarget" method on click', async () => {
    await skipButton.trigger('click')

    expect(spySkipToTarget).toHaveBeenCalledTimes(1)
  })
})

function getShallowWrapper() {
  return {
    wrapper: shallowMount(SkipTo, {
      props: {
        target: ''
      },
      slots: {
        default: 'Skip to main'
      }
    })
  }
}
