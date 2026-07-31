import { Resource } from '@ownclouders/web-client'
import { canBeMoved } from '../../../src/helpers/permissions'

describe('permissions helper', () => {
  describe('canBeMoved function', () => {
    it.each([
      {
        name: 'forest.jpg',
        isReceivedShare: false,
        isMounted: false,
        canBeDeleted: true,
        parentPath: ''
      },
      {
        name: 'forest.jpg',
        isReceivedShare: false,
        isMounted: false,
        canBeDeleted: true,
        parentPath: 'folder'
      },
      {
        name: 'forest.jpg',
        isReceivedShare: true,
        isMounted: false,
        canBeDeleted: true,
        parentPath: 'folder'
      },
      {
        name: 'forest.jpg',
        isReceivedShare: false,
        isMounted: true,
        canBeDeleted: true,
        parentPath: 'folder'
      }
    ])(
      'should return true if the given resource can be deleted and if it is not mounted in root',
      (input) => {
        // resources are supposed to be external if it is a received share or is mounted
        // resources are supposed to be mountedInRoot if its parentPath is an empty string and resource is external
        expect(
          canBeMoved(
            {
              name: input.name,
              isReceivedShare: () => input.isReceivedShare,
              isMounted: () => input.isMounted,
              canBeDeleted: () => input.canBeDeleted
            } as unknown as Resource,
            input.parentPath
          )
        ).toBeTruthy()
      }
    )
    it.each([
      {
        name: 'forest.jpg',
        isReceivedShare: false,
        isMounted: false,
        canBeDeleted: false,
        parentPath: ''
      },
      {
        name: 'forest.jpg',
        isReceivedShare: false,
        isMounted: false,
        canBeDeleted: false,
        parentPath: 'folder'
      },
      {
        name: 'forest.jpg',
        isReceivedShare: true,
        isMounted: false,
        canBeDeleted: false,
        parentPath: 'folder'
      },
      {
        name: 'forest.jpg',
        isReceivedShare: false,
        isMounted: true,
        canBeDeleted: false,
        parentPath: 'folder'
      },
      {
        name: 'forest.jpg',
        isReceivedShare: false,
        isMounted: true,
        canBeDeleted: true,
        parentPath: ''
      },
      {
        name: 'forest.jpg',
        isReceivedShare: true,
        isMounted: false,
        canBeDeleted: true,
        parentPath: ''
      },
      {
        name: 'forest.jpg',
        isReceivedShare: true,
        isMounted: true,
        canBeDeleted: true,
        parentPath: ''
      },
      {
        name: 'forest.psec',
        isReceivedShare: false,
        isMounted: false,
        canBeDeleted: true,
        parentPath: ''
      }
    ])(
      'should return false if the given resource cannot be deleted or if it is mounted in root',
      (input) => {
        expect(
          canBeMoved(
            {
              name: input.name,
              isReceivedShare: () => input.isReceivedShare,
              isMounted: () => input.isMounted,
              canBeDeleted: () => input.canBeDeleted
            } as unknown as Resource,
            input.parentPath
          )
        ).toBeFalsy()
      }
    )
  })
})
