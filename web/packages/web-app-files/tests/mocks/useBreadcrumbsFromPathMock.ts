import { useBreadcrumbsFromPath } from '@ownclouders/web-pkg'

export const useBreadcrumbsFromPathMock = (
  options: Partial<ReturnType<typeof useBreadcrumbsFromPath>> = {}
): ReturnType<typeof useBreadcrumbsFromPath> => {
  return {
    breadcrumbsFromPath: vi.fn(() => []),
    concatBreadcrumbs: vi.fn((...args) => args),
    ...options
  }
}
