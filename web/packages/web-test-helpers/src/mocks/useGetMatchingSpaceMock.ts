import { mock } from 'vitest-mock-extended'
import { useGetMatchingSpace } from '@ownclouders/web-pkg'
import { SpaceResource } from '@ownclouders/web-client'

export const useGetMatchingSpaceMock = (
  options: Partial<ReturnType<typeof useGetMatchingSpace>> = {}
): ReturnType<typeof useGetMatchingSpace> => {
  return {
    getInternalSpace() {
      return mock<SpaceResource>()
    },
    getMatchingSpace() {
      return mock<SpaceResource>()
    },
    isResourceAccessible() {
      return false
    },
    isPersonalSpaceRoot() {
      return false
    },
    ...options
  }
}
