import { useScrollTo } from '@ownclouders/web-pkg'

export const useScrollToMock = (
  options: Partial<ReturnType<typeof useScrollTo>> = {}
): ReturnType<typeof useScrollTo> => {
  return {
    scrollToResource: vi.fn(),
    scrollToResourceFromRoute: vi.fn(),
    ...options
  }
}
