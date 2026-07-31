import { mock, mockDeep } from 'vitest-mock-extended'
import {
  IncomingShareResource,
  OutgoingShareResource,
  Resource,
  ShareResource,
  ShareRole,
  ShareTypes
} from '../../../../src/helpers'
import {
  buildCollaboratorShare,
  buildIncomingShareResource,
  buildLinkShare,
  buildOutgoingShareResource,
  getShareResourcePermissions,
  getShareResourceRoles,
  isIncomingShareResource,
  isOutgoingShareResource,
  isShareResource
} from '../../../../src/helpers/share/functions'
import {
  DriveItem,
  Identity,
  Permission,
  UnifiedRoleDefinition,
  User
} from '../../../../src/graph/generated'
import { urlJoin } from '../../../../src'

describe('share helper functions', () => {
  describe('isShareResource', () => {
    it('returns true for shares based on "sharedWith" property', () => {
      const resource = mock<ShareResource>({ sharedWith: mock<ShareResource['sharedWith']>() })
      expect(isShareResource(resource)).toBeTruthy()
    })
    it('returns false for regular resources based on "sharedWith" property', () => {
      const resource = mock<Resource>()
      expect(isShareResource(resource)).toBeFalsy()
    })
  })

  describe('isOutgoingShareResource', () => {
    it('returns true for outgoing shares', () => {
      const resource = mock<OutgoingShareResource>({
        outgoing: true,
        sharedWith: mock<ShareResource['sharedWith']>()
      })
      expect(isOutgoingShareResource(resource)).toBeTruthy()
    })
    it('returns false for incoming shares', () => {
      const resource = mock<IncomingShareResource>({
        outgoing: false,
        sharedWith: mock<ShareResource['sharedWith']>()
      })
      expect(isOutgoingShareResource(resource)).toBeFalsy()
    })
  })

  describe('isIncomingShareResource', () => {
    it('returns true for incoming shares', () => {
      const resource = mock<IncomingShareResource>({
        outgoing: false,
        sharedWith: mock<ShareResource['sharedWith']>()
      })
      expect(isIncomingShareResource(resource)).toBeTruthy()
    })
    it('returns false for outgoing shares', () => {
      const resource = mock<OutgoingShareResource>({
        outgoing: true,
        sharedWith: mock<ShareResource['sharedWith']>()
      })
      expect(isIncomingShareResource(resource)).toBeFalsy()
    })
  })

  describe('getShareResourceRoles', () => {
    it("returns all roles from a drive item's permissions that are also included in the graphRoles", () => {
      const driveItem = mockDeep<DriveItem>()
      driveItem.remoteItem.permissions = [{ roles: ['1', '2'] }, { roles: ['1', '3'] }]
      const graphRoles = { '1': mock<ShareRole>({ id: '1' }), '4': mock<ShareRole>({ id: '4' }) }

      const result = getShareResourceRoles({ driveItem, graphRoles })

      expect(result.length).toBe(1)
      expect(result[0].id).toEqual('1')
    })
  })

  describe('getShareResourcePermissions', () => {
    it('returns permissions based on the given graph share roles', () => {
      const permissions = ['view', 'edit']
      const shareRoles = [
        { rolePermissions: [{ allowedResourceActions: [permissions[0]] }] },
        { rolePermissions: [{ allowedResourceActions: [permissions[1]] }] }
      ] as UnifiedRoleDefinition[]

      const result = getShareResourcePermissions({ driveItem: undefined, shareRoles })

      expect(result).toEqual(permissions)
    })
    it('returns permissions based on a drive item if no graph share roles given', () => {
      const permissions = ['view', 'edit']
      const driveItem = mockDeep<DriveItem>()
      driveItem.remoteItem.permissions = [
        { '@libre.graph.permissions.actions': [permissions[0]] },
        { '@libre.graph.permissions.actions': [permissions[1]] }
      ]

      const result = getShareResourcePermissions({ driveItem, shareRoles: [] })

      expect(result).toEqual(permissions)
    })
  })

  describe('buildIncomingShareResource', () => {
    const driveItem = mockDeep<DriveItem>({ id: 'driveItemId', name: 'driveItemName' })
    const sharedBy = { id: '1', displayName: 'user1' } as Identity
    const sharedWith = { id: '2', displayName: 'user2' } as Identity
    driveItem.remoteItem.permissions = [
      {
        roles: ['1', '2'],
        invitation: { invitedBy: { user: sharedBy } },
        grantedToV2: { user: sharedWith }
      }
    ]

    const graphRoles = {
      '1': mock<ShareRole>({ id: '1', rolePermissions: [{ allowedResourceActions: ['view'] }] }),
      '2': mock<ShareRole>({ id: '1', rolePermissions: [{ allowedResourceActions: ['view'] }] })
    }

    it('sets ids based on the drive item, its first permission and parent reference', () => {
      const result = buildIncomingShareResource({ driveItem, graphRoles, serverUrl: '' })

      expect(result.id).toEqual(driveItem.id)
      expect(result.fileId).toEqual(driveItem.remoteItem.id)
      expect(result.remoteItemId).toEqual(driveItem.remoteItem.id)
      expect(result.driveId).toEqual(driveItem.parentReference.driveId)
      expect(result.parentFolderId).toEqual(driveItem.parentReference.id)
    })
    it.each([true, false])('correctly detects if the resource is a folder', (isFolder) => {
      const item = { ...driveItem }
      item.folder = isFolder ? mock<DriveItem['folder']>() : undefined
      const result = buildIncomingShareResource({ driveItem: item, graphRoles, serverUrl: '' })

      expect(result.isFolder).toEqual(isFolder)
      expect(result.type).toEqual(isFolder ? 'folder' : 'file')
    })
    it('sets outgoing to false', () => {
      const result = buildIncomingShareResource({ driveItem, graphRoles, serverUrl: '' })
      expect(result.outgoing).toBeFalsy()
    })
    it('sets sharedBy based on the permission invitation', () => {
      const result = buildIncomingShareResource({ driveItem, graphRoles, serverUrl: '' })
      expect(result.sharedBy).toEqual([sharedBy])
    })
    it('sets sharedBy based on the permission invitation', () => {
      const result = buildIncomingShareResource({ driveItem, graphRoles, serverUrl: '' })
      expect(result.sharedWith).toEqual([{ ...sharedWith, shareType: ShareTypes.user.value }])
    })
    it('constructs a private link', () => {
      const serverUrl = 'https://example.com'
      const result = buildIncomingShareResource({ driveItem, graphRoles, serverUrl })
      expect(result.privateLink).toEqual(urlJoin(serverUrl, 'f', driveItem.remoteItem.id))
    })
  })

  describe('buildOutgoingShareResource', () => {
    const driveItem = mockDeep<DriveItem>({ id: 'driveItemId', name: 'driveItemName' })
    driveItem.parentReference.path = ''
    const sharedBy = { id: '1', displayName: 'user1' } as Identity
    const sharedWith = { id: '2', displayName: 'user2' } as Identity
    driveItem.permissions = [
      {
        roles: ['1', '2'],
        invitation: { invitedBy: { user: sharedBy } },
        grantedToV2: { user: sharedWith }
      }
    ]
    const user = { id: '1', displayName: 'user1' } as User

    it('sets ids based on the drive item, its first permission and parent reference', () => {
      const result = buildOutgoingShareResource({ driveItem, user, serverUrl: '' })

      expect(result.id).toEqual(driveItem.id)
      expect(result.fileId).toEqual(driveItem.id)
      expect(result.driveId).toEqual(driveItem.parentReference.driveId)
      expect(result.parentFolderId).toEqual(driveItem.parentReference.id)
    })
    it('sets outgoing to true', () => {
      const result = buildOutgoingShareResource({ driveItem, user, serverUrl: '' })
      expect(result.outgoing).toBeTruthy()
    })
    it('sets the path based on the parent reference path and the drive item name', () => {
      const result = buildOutgoingShareResource({ driveItem, user, serverUrl: '' })
      expect(result.path).toEqual(`${driveItem.parentReference.path}/${driveItem.name}`)
    })
    it.each([true, false])('correctly detects if the resource is a folder', (isFolder) => {
      const item = { ...driveItem }
      item.folder = isFolder ? mock<DriveItem['folder']>() : undefined
      const result = buildOutgoingShareResource({ driveItem: item, user, serverUrl: '' })

      expect(result.isFolder).toEqual(isFolder)
      expect(result.type).toEqual(isFolder ? 'folder' : 'file')
    })
    it('constructs a private link', () => {
      const serverUrl = 'https://example.com'
      const result = buildOutgoingShareResource({ driveItem, user, serverUrl })
      expect(result.privateLink).toEqual(urlJoin(serverUrl, 'f', driveItem.id))
    })
  })

  describe('buildCollaboratorShare', () => {
    const graphRoles = {
      '1': mock<ShareRole>({ id: '1', rolePermissions: [{ allowedResourceActions: ['view'] }] }),
      '2': mock<ShareRole>({ id: '1', rolePermissions: [{ allowedResourceActions: ['view'] }] })
    }

    const resourceId = '1'

    it('sets ids based on the permission and the given resource id', () => {
      const graphPermission = mock<Permission>({ '@libre.graph.permissions.actions': [] })

      const result = buildCollaboratorShare({
        graphPermission,
        graphRoles,
        resourceId
      })

      expect(result.id).toEqual(graphPermission.id)
      expect(result.resourceId).toEqual(resourceId)
    })
    describe('share type', () => {
      it('is user type if grantedToV2 includes a user', () => {
        const graphPermission = mock<Permission>({
          '@libre.graph.permissions.actions': [],
          grantedToV2: { user: {}, group: undefined },
          link: undefined
        })

        const result = buildCollaboratorShare({
          graphPermission,
          graphRoles,
          resourceId
        })

        expect(result.shareType).toEqual(ShareTypes.user.value)
      })
      it('is group type if grantedToV2 includes a group', () => {
        const graphPermission = mock<Permission>({
          '@libre.graph.permissions.actions': [],
          grantedToV2: { user: undefined, group: {} },
          link: undefined
        })

        const result = buildCollaboratorShare({
          graphPermission,
          graphRoles,
          resourceId
        })

        expect(result.shareType).toEqual(ShareTypes.group.value)
      })
      it('is external type if grantedToV2 includes a user that is external', () => {
        const graphPermission = mock<Permission>({
          '@libre.graph.permissions.actions': [],
          grantedToV2: { user: { '@libre.graph.userType': 'Federated' }, group: undefined },
          link: undefined
        })

        const result = buildCollaboratorShare({
          graphPermission,
          graphRoles,
          resourceId
        })

        expect(result.shareType).toEqual(ShareTypes.remote.value)
      })
    })
    describe('permissions', () => {
      it('sets permissions if given directly via property', () => {
        const permissions = ['view', 'edit']
        const graphPermission = mock<Permission>({
          '@libre.graph.permissions.actions': permissions
        })

        const result = buildCollaboratorShare({
          graphPermission,
          graphRoles,
          resourceId
        })

        expect(result.permissions).toEqual(permissions)
      })
      it('sets permissions from the graph roles as fallback', () => {
        const graphPermission = mock<Permission>({
          '@libre.graph.permissions.actions': undefined,
          roles: [graphRoles['1'].id]
        })

        const result = buildCollaboratorShare({
          graphPermission,
          graphRoles,
          resourceId
        })

        expect(result.permissions).toEqual(
          graphRoles['1'].rolePermissions.flatMap(
            ({ allowedResourceActions }) => allowedResourceActions
          )
        )
      })
    })
  })

  describe('buildLinkShare', () => {
    const resourceId = '1'

    it('sets ids based on the permission and the given resource id', () => {
      const graphPermission = mock<Permission>({ '@libre.graph.permissions.actions': [] })
      const result = buildLinkShare({ graphPermission, resourceId })

      expect(result.id).toEqual(graphPermission.id)
      expect(result.resourceId).toEqual(resourceId)
    })
    it('sets the sharing link type', () => {
      const graphPermission = mock<Permission>({ '@libre.graph.permissions.actions': [] })
      const result = buildLinkShare({ graphPermission, resourceId })

      expect(result.shareType).toEqual(ShareTypes.link.value)
    })
  })
})
