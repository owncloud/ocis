import { buildSpace } from '../../../../src/helpers/space'
import { mock } from 'vitest-mock-extended'
import { Ability, GraphSharePermission, ShareRole } from '@ownclouders/web-client'
import { Drive, User } from '@ownclouders/web-client/graph/generated'

const noPermissionsRole = mock<ShareRole>({ id: '1', rolePermissions: [] })
const canUploadRole = mock<ShareRole>({
  id: '2',
  rolePermissions: [
    {
      condition: 'exists @Resource.Root',
      allowedResourceActions: [GraphSharePermission.createUpload]
    }
  ]
})
const canDeletePermissionsRole = mock<ShareRole>({
  id: '3',
  rolePermissions: [
    {
      condition: 'exists @Resource.Root',
      allowedResourceActions: [GraphSharePermission.deletePermissions]
    }
  ]
})
const canCreatePermissionsRole = mock<ShareRole>({
  id: '4',
  rolePermissions: [
    {
      condition: 'exists @Resource.Root',
      allowedResourceActions: [GraphSharePermission.createPermissions]
    }
  ]
})

const graphRoles = {
  [noPermissionsRole.id]: noPermissionsRole,
  [canUploadRole.id]: canUploadRole,
  [canDeletePermissionsRole.id]: canDeletePermissionsRole,
  [canCreatePermissionsRole.id]: canCreatePermissionsRole
}

describe('buildSpace', () => {
  const id = '1'

  const getSpace = ({ role, permissions }: { role: ShareRole; permissions: string[] }) => {
    return buildSpace(
      {
        special: [],
        root: {
          permissions: [
            {
              roles: role ? [role.id] : [],
              grantedToV2: { user: { id } },
              ...(permissions.length && { '@libre.graph.permissions.actions': permissions })
            }
          ]
        }
      } as Drive,
      graphRoles
    )
  }

  describe('canUpload', () => {
    it.each([
      { permissions: [GraphSharePermission.createUpload], role: undefined, expectedResult: true },
      { permissions: [], role: canUploadRole, expectedResult: true },
      { permissions: [], role: noPermissionsRole, expectedResult: false }
    ])(
      'behaves accordingly to the given role and permissions',
      ({ permissions, role, expectedResult }) => {
        const space = getSpace({ role, permissions })
        expect(space.canUpload({ user: mock<User>({ id, memberOf: [] }) })).toBe(expectedResult)
      }
    )
  })

  describe('canDownload', () => {
    it('is always true', () => {
      const space = getSpace({ role: noPermissionsRole, permissions: [] })
      expect(space.canDownload()).toBeTruthy()
    })
  })

  describe('canBeDeleted', () => {
    it.each([
      {
        userCan: false,
        permissions: [GraphSharePermission.deletePermissions],
        role: undefined,
        disabled: true,
        expectedResult: true
      },
      {
        userCan: false,
        permissions: [],
        role: canDeletePermissionsRole,
        disabled: true,
        expectedResult: true
      },
      {
        userCan: false,
        permissions: [],
        role: noPermissionsRole,
        disabled: true,
        expectedResult: false
      },
      {
        userCan: true,
        permissions: [],
        role: noPermissionsRole,
        disabled: true,
        expectedResult: true
      },
      {
        userCan: true,
        permissions: [],
        role: noPermissionsRole,
        disabled: false,
        expectedResult: false
      }
    ])(
      'behaves accordingly to the given role, permissions, abilities and disabled state',
      ({ permissions, role, expectedResult, userCan, disabled }) => {
        const ability = mock<Ability>({ can: () => userCan })
        const space = getSpace({ role, permissions })
        space.disabled = disabled
        expect(space.canBeDeleted({ user: mock<User>({ id, memberOf: [] }), ability })).toBe(
          expectedResult
        )
      }
    )
  })

  describe('canRename', () => {
    it.each([
      {
        userCan: false,
        permissions: [GraphSharePermission.deletePermissions],
        role: undefined,
        expectedResult: true
      },
      {
        userCan: false,
        permissions: [],
        role: canDeletePermissionsRole,
        expectedResult: true
      },
      {
        userCan: false,
        permissions: [],
        role: noPermissionsRole,
        expectedResult: false
      },
      {
        userCan: true,
        permissions: [],
        role: noPermissionsRole,
        expectedResult: true
      }
    ])(
      'behaves accordingly to the given role, permissions and abilities',
      ({ permissions, role, expectedResult, userCan }) => {
        const ability = mock<Ability>({ can: () => userCan })
        const space = getSpace({ role, permissions })
        expect(space.canRename({ user: mock<User>({ id, memberOf: [] }), ability })).toBe(
          expectedResult
        )
      }
    )
  })

  describe('canEditDescription', () => {
    it.each([
      {
        userCan: false,
        permissions: [GraphSharePermission.deletePermissions],
        role: undefined,
        expectedResult: true
      },
      {
        userCan: false,
        permissions: [],
        role: canDeletePermissionsRole,
        expectedResult: true
      },
      {
        userCan: false,
        permissions: [],
        role: noPermissionsRole,
        expectedResult: false
      },
      {
        userCan: true,
        permissions: [],
        role: noPermissionsRole,
        expectedResult: true
      }
    ])(
      'behaves accordingly to the given role, permissions and abilities',
      ({ permissions, role, expectedResult, userCan }) => {
        const ability = mock<Ability>({ can: () => userCan })
        const space = getSpace({ role, permissions })
        expect(space.canEditDescription({ user: mock<User>({ id, memberOf: [] }), ability })).toBe(
          expectedResult
        )
      }
    )
  })

  describe('canShare', () => {
    it.each([
      {
        permissions: [GraphSharePermission.createPermissions],
        role: undefined,
        expectedResult: true
      },
      { permissions: [], role: canCreatePermissionsRole, expectedResult: true },
      { permissions: [], role: noPermissionsRole, expectedResult: false }
    ])(
      'behaves accordingly to the given role and permissions',
      ({ permissions, role, expectedResult }) => {
        const space = getSpace({ role, permissions })
        expect(space.canShare({ user: mock<User>({ id, memberOf: [] }) })).toBe(expectedResult)
      }
    )
  })

  describe('canRestore', () => {
    it.each([
      {
        userCan: false,
        permissions: [GraphSharePermission.deletePermissions],
        role: undefined,
        disabled: true,
        expectedResult: true
      },
      {
        userCan: false,
        permissions: [],
        role: canDeletePermissionsRole,
        disabled: true,
        expectedResult: true
      },
      {
        userCan: false,
        permissions: [],
        role: noPermissionsRole,
        disabled: true,
        expectedResult: false
      },
      {
        userCan: true,
        permissions: [],
        role: noPermissionsRole,
        disabled: true,
        expectedResult: true
      },
      {
        userCan: true,
        permissions: [],
        role: noPermissionsRole,
        disabled: false,
        expectedResult: false
      }
    ])(
      'behaves accordingly to the given role, permissions, abilities and disabled state',
      ({ permissions, role, expectedResult, userCan, disabled }) => {
        const ability = mock<Ability>({ can: () => userCan })
        const space = getSpace({ role, permissions })
        space.disabled = disabled
        expect(space.canRestore({ user: mock<User>({ id, memberOf: [] }), ability })).toBe(
          expectedResult
        )
      }
    )
  })

  describe('canDisable', () => {
    it.each([
      {
        userCan: false,
        permissions: [GraphSharePermission.deletePermissions],
        role: undefined,
        disabled: false,
        expectedResult: true
      },
      {
        userCan: false,
        permissions: [],
        role: canDeletePermissionsRole,
        disabled: false,
        expectedResult: true
      },
      {
        userCan: false,
        permissions: [],
        role: noPermissionsRole,
        disabled: false,
        expectedResult: false
      },
      {
        userCan: true,
        permissions: [],
        role: noPermissionsRole,
        disabled: false,
        expectedResult: true
      },
      {
        userCan: true,
        permissions: [],
        role: noPermissionsRole,
        disabled: true,
        expectedResult: false
      }
    ])(
      'behaves accordingly to the given role, permissions, abilities and disabled state',
      ({ permissions, role, expectedResult, userCan, disabled }) => {
        const ability = mock<Ability>({ can: () => userCan })
        const space = getSpace({ role, permissions })
        space.disabled = disabled
        expect(space.canDisable({ user: mock<User>({ id, memberOf: [] }), ability })).toBe(
          expectedResult
        )
      }
    )
  })

  describe('canEditImage', () => {
    it.each([
      {
        permissions: [GraphSharePermission.deletePermissions],
        role: undefined,
        disabled: false,
        expectedResult: true
      },
      {
        permissions: [],
        role: canDeletePermissionsRole,
        disabled: false,
        expectedResult: true
      },
      {
        permissions: [],
        role: noPermissionsRole,
        disabled: false,
        expectedResult: false
      },
      {
        permissions: [],
        role: noPermissionsRole,
        disabled: true,
        expectedResult: false
      }
    ])(
      'behaves accordingly to the given role, permissions and disabled state',
      ({ permissions, role, expectedResult, disabled }) => {
        const space = getSpace({ role, permissions })
        space.disabled = disabled
        expect(space.canEditImage({ user: mock<User>({ id, memberOf: [] }) })).toBe(expectedResult)
      }
    )
  })
  describe('canEditReadme', () => {
    it.each([
      {
        permissions: [GraphSharePermission.deletePermissions],
        role: undefined,
        disabled: false,
        expectedResult: true
      },
      {
        permissions: [],
        role: canDeletePermissionsRole,
        disabled: false,
        expectedResult: true
      },
      {
        permissions: [],
        role: noPermissionsRole,
        disabled: false,
        expectedResult: false
      },
      {
        permissions: [],
        role: noPermissionsRole,
        disabled: true,
        expectedResult: false
      }
    ])(
      'behaves accordingly to the given role, permissions and disabled state',
      ({ permissions, role, expectedResult, disabled }) => {
        const space = getSpace({ role, permissions })
        space.disabled = disabled
        expect(space.canEditReadme({ user: mock<User>({ id, memberOf: [] }) })).toBe(expectedResult)
      }
    )
  })
})
