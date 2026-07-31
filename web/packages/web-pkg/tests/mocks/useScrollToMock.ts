import { useScrollTo } from '../../src/composables/scrollTo'

export const useScrollToMock = (
  options: Partial<ReturnType<typeof useScrollTo>> = {}
): ReturnType<typeof useScrollTo> => {
  return {
    scrollToResource: vi.fn(),
    scrollToResourceFromRoute: vi.fn(),
    ...options
  }
}
