import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import {
  IncomingShareResource,
  Resource,
  SpaceResource,
  TrashResource
} from '@ownclouders/web-client'
import { useCanListVersions } from '../../../../src/composables/resources'

describe('useCanListVersions', () => {
  describe('canListVersions', () => {
    it('returns true for files when user has sufficient permissions in space', () => {
      getWrapper({
        setup: ({ canListVersions }) => {
          const space = mock<SpaceResource>({ canListVersions: () => true })
          const resource = mock<Resource>({ type: 'file' })
          const canList = canListVersions({ space, resource })
          expect(canList).toBeTruthy()
        }
      })
    })
    it('returns false for folders', () => {
      getWrapper({
        setup: ({ canListVersions }) => {
          const space = mock<SpaceResource>({ canListVersions: () => true })
          const resource = mock<Resource>({ type: 'folder' })
          const canList = canListVersions({ space, resource })
          expect(canList).toBeFalsy()
        }
      })
    })
    it('returns false for space resources', () => {
      getWrapper({
        setup: ({ canListVersions }) => {
          const space = mock<SpaceResource>({ canListVersions: () => true })
          const resource = mock<SpaceResource>({ type: 'space' })
          const canList = canListVersions({ space, resource })
          expect(canList).toBeFalsy()
        }
      })
    })
    it('returns false for trash resources', () => {
      getWrapper({
        setup: ({ canListVersions }) => {
          const space = mock<SpaceResource>({ canListVersions: () => true })
          const resource = mock<TrashResource>({ type: 'file', ddate: '' })
          const canList = canListVersions({ space, resource })
          expect(canList).toBeFalsy()
        }
      })
    })
    it('returns false when user does not have sufficient permissions in space', () => {
      getWrapper({
        setup: ({ canListVersions }) => {
          const space = mock<SpaceResource>({ canListVersions: () => false })
          const resource = mock<Resource>({ type: 'file' })
          const canList = canListVersions({ space, resource })
          expect(canList).toBeFalsy()
        }
      })
    })

    it('should use resource permissions instead of space permissions for received shares', () => {
      getWrapper({
        setup: ({ canListVersions }) => {
          const space = mock<SpaceResource>({ canListVersions: () => false })
          const resource = mock<IncomingShareResource>({
            type: 'file',
            isReceivedShare: () => true,
            canListVersions: () => true
          })
          const canList = canListVersions({ space, resource })
          expect(canList).toBeTruthy()
        }
      })
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (instance: ReturnType<typeof useCanListVersions>) => void
}) {
  return {
    wrapper: getComposableWrapper(() => {
      const instance = useCanListVersions()
      setup(instance)
    })
  }
}
