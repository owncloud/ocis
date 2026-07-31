import { mock } from 'vitest-mock-extended'
import { Resource, SpaceResource } from '@ownclouders/web-client'
import { getIndicators } from '../../../src/helpers/statusIndicators'
import { User } from '@ownclouders/web-client/graph/generated'
import { AncestorMetaDataValue } from '../../../src/types'
import { ResourceIndicator } from '../../../src/helpers'
import { createTestingPinia } from '@ownclouders/web-test-helpers'
import {
  ResourceIndicatorExtension,
  useExtensionRegistry
} from '../../../src/composables/piniaStores/extensionRegistry'

describe('status indicators', () => {
  const user = mock<User>()

  createTestingPinia()

  describe('indicator extensions', () => {
    it('should be requested from the extension registry', () => {
      const space = mock<SpaceResource>({ driveType: 'project' })
      const resource = mock<Resource>({ id: 'resource' })

      const { requestExtensions } = useExtensionRegistry()
      vi.mocked(requestExtensions<ResourceIndicatorExtension>).mockReturnValue([
        {
          id: 'test.files.resource-indicator.stub',
          type: 'resourceIndicator',
          extensionPointIds: ['global.files.resource-indicator'],
          getResourceIndicators: (resource: Resource): ResourceIndicator[] => {
            return [
              {
                id: 'some-id',
                accessibleDescription: 'some accessible description',
                label: 'some label',
                icon: 'check_box_outline_blank',
                fillType: 'line',
                type: 'some-type',
                category: 'system'
              }
            ]
          }
        } satisfies ResourceIndicatorExtension
      ])

      const indicators = getIndicators({ space, resource, ancestorMetaData: {}, user })

      expect(requestExtensions).toHaveBeenCalled()
      expect(indicators.some(({ id }) => id === 'some-id')).toBeTruthy()
    })
  })

  describe('locked indicator', () => {
    it.each([true, false])('should only be present if the file is locked', (locked) => {
      const space = mock<SpaceResource>({ id: 'space' })
      const resource = mock<Resource>({ id: 'resource', locked })
      const indicators = getIndicators({ space, resource, ancestorMetaData: {}, user })

      expect(indicators.some(({ type }) => type === 'resource-locked')).toBe(locked)
    })
  })

  describe('processing indicator', () => {
    it.each([true, false])('should only be present if the file is processing', (processing) => {
      const space = mock<SpaceResource>({ id: 'space' })
      const resource = mock<Resource>({ id: 'resource', processing })
      const indicators = getIndicators({ space, resource, ancestorMetaData: {}, user })

      expect(indicators.some(({ type }) => type === 'resource-processing')).toBe(processing)
    })
  })

  describe('sharing indicators', () => {
    it('should not be present if the user is not a member of the project space', () => {
      const space = mock<SpaceResource>({ driveType: 'project', isMember: () => false })
      const resource = mock<Resource>({ id: 'resource', shareTypes: [0, 3] })
      const indicators = getIndicators({ space, resource, ancestorMetaData: {}, user })

      expect(indicators.some(({ category }) => category === 'sharing')).toBeFalsy()
    })
    it("should not be present in another user's personal space", () => {
      const space = mock<SpaceResource>({ driveType: 'personal', isOwner: () => false })
      const resource = mock<Resource>({ id: 'resource', shareTypes: [0, 3] })
      const indicators = getIndicators({ space, resource, ancestorMetaData: {}, user })

      expect(indicators.some(({ category }) => category === 'sharing')).toBeFalsy()
    })
    it('should not be present in a share space', () => {
      const space = mock<SpaceResource>({ driveType: 'share' })
      const resource = mock<Resource>({ id: 'resource', shareTypes: [0, 3] })
      const indicators = getIndicators({ space, resource, ancestorMetaData: {}, user })

      expect(indicators.some(({ category }) => category === 'sharing')).toBeFalsy()
    })
    it('should not be present in a public space', () => {
      const space = mock<SpaceResource>({ driveType: 'public' })
      const resource = mock<Resource>({ id: 'resource', shareTypes: [0, 3] })
      const indicators = getIndicators({ space, resource, ancestorMetaData: {}, user })

      expect(indicators.some(({ category }) => category === 'sharing')).toBeFalsy()
    })
    it('should be present for direct collaborator and link shares', () => {
      const space = mock<SpaceResource>({ driveType: 'project', isMember: () => true })
      const resource = mock<Resource>({ id: 'resource', shareTypes: [0, 3] })
      const indicators = getIndicators({ space, resource, ancestorMetaData: {}, user })

      expect(
        indicators.some(({ type, category }) => category === 'sharing' && 'link-direct')
      ).toBeTruthy()
      expect(
        indicators.some(({ type, category }) => category === 'sharing' && 'user-direct')
      ).toBeTruthy()
    })
    it('should be present for indirect collaborator and link shares', () => {
      const ancestorMetaData = {
        '/': mock<AncestorMetaDataValue>({ shareTypes: [0, 3] })
      }

      const space = mock<SpaceResource>({ driveType: 'project', isMember: () => true })
      const resource = mock<Resource>({ id: 'resource', shareTypes: [] })
      const indicators = getIndicators({ space, resource, ancestorMetaData, user })

      expect(
        indicators.some(({ type, category }) => category === 'sharing' && type === 'link-indirect')
      ).toBeTruthy()
      expect(
        indicators.some(({ type, category }) => category === 'sharing' && type === 'user-indirect')
      ).toBeTruthy()
    })
  })
})
