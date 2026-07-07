import { unref } from 'vue'
import { useLinkTypes } from '../../../../src/composables/links/useLinkTypes'
import { defaultComponentMocks, getComposableWrapper } from '@ownclouders/web-test-helpers'
import { SharingLinkType } from '@ownclouders/web-client/graph/generated'
import { AbilityRule } from '@ownclouders/web-client'
import { Capabilities } from '@ownclouders/web-client/ocs'
import { PasswordEnforcedForCapability } from '@ownclouders/web-client/ocs'

describe('useLinkTypes', () => {
  it('should be valid', () => {
    expect(useLinkTypes).toBeDefined()
  })
  describe('computed "defaultLinkType"', () => {
    it('is viewer', () => {
      getWrapper({
        setup: ({ defaultLinkType }) => {
          expect(unref(defaultLinkType)).toBe(SharingLinkType.View)
        }
      })
    })
  })
  describe('method "getAvailableLinkTypes"', () => {
    it('is empty if the user cannot create public links', () => {
      getWrapper({
        abilities: [],
        setup: ({ getAvailableLinkTypes }) => {
          expect(getAvailableLinkTypes({ isFolder: true })).toEqual([])
        }
      })
    })
    it('returns all available types for folders accordingly', () => {
      getWrapper({
        abilities: [{ action: 'create-all', subject: 'PublicLink' }],
        setup: ({ getAvailableLinkTypes }) => {
          expect(getAvailableLinkTypes({ isFolder: true })).toEqual([
            SharingLinkType.View,
            SharingLinkType.Edit,
            SharingLinkType.CreateOnly
          ])
        }
      })
    })
    it('returns all available types for files accordingly', () => {
      getWrapper({
        abilities: [{ action: 'create-all', subject: 'PublicLink' }],
        setup: ({ getAvailableLinkTypes }) => {
          expect(getAvailableLinkTypes({ isFolder: false })).toEqual([
            SharingLinkType.View,
            SharingLinkType.Edit
          ])
        }
      })
    })
  })
  describe('method "getLinkRoleByType"', () => {
    it.each([
      SharingLinkType.Internal,
      SharingLinkType.View,
      SharingLinkType.Upload,
      SharingLinkType.Edit,
      SharingLinkType.CreateOnly
    ])('returns the link role by id', (type) => {
      getWrapper({
        setup: ({ getLinkRoleByType, linkShareRoles }) => {
          expect(getLinkRoleByType(type)).toEqual(linkShareRoles.find(({ id }) => id === type))
        }
      })
    })
  })
  describe('method "isPasswordEnforcedForLinkType"', () => {
    it('returns true for view type if set via capabilities', () => {
      getWrapper({
        passwordEnforcedFor: { read_only: true },
        setup: ({ isPasswordEnforcedForLinkType }) => {
          expect(isPasswordEnforcedForLinkType(SharingLinkType.View)).toBeTruthy()
        }
      })
    })
    it('returns true for upload type if set via capabilities', () => {
      getWrapper({
        passwordEnforcedFor: { upload_only: true },
        setup: ({ isPasswordEnforcedForLinkType }) => {
          expect(isPasswordEnforcedForLinkType(SharingLinkType.Upload)).toBeTruthy()
        }
      })
    })
    it('returns true for create only type if set via capabilities', () => {
      getWrapper({
        passwordEnforcedFor: { read_write: true },
        setup: ({ isPasswordEnforcedForLinkType }) => {
          expect(isPasswordEnforcedForLinkType(SharingLinkType.CreateOnly)).toBeTruthy()
        }
      })
    })
    it('returns true for edit type if set via capabilities', () => {
      getWrapper({
        passwordEnforcedFor: { read_write_delete: true },
        setup: ({ isPasswordEnforcedForLinkType }) => {
          expect(isPasswordEnforcedForLinkType(SharingLinkType.Edit)).toBeTruthy()
        }
      })
    })
  })
})

function getWrapper({
  setup,
  abilities = [],
  defaultPermissions = undefined,
  passwordEnforcedFor = undefined
}: {
  setup: (instance: ReturnType<typeof useLinkTypes>) => void
  abilities?: AbilityRule[]
  defaultPermissions?: number
  passwordEnforcedFor?: PasswordEnforcedForCapability
}) {
  const mocks = defaultComponentMocks()

  const capabilities = {
    files_sharing: {
      public: {
        default_permissions: defaultPermissions,
        password: { enforced_for: passwordEnforcedFor }
      }
    }
  } satisfies Partial<Capabilities['capabilities']>

  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useLinkTypes()
        setup(instance)
      },
      {
        mocks,
        pluginOptions: { abilities, piniaOptions: { capabilityState: { capabilities } } }
      }
    )
  }
}
