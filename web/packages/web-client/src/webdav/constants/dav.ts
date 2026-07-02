import { Audio, GeoCoordinates, Image, Photo } from '../../graph/generated'

export abstract class DavPermission {
  static readonly Shared: string = 'S'
  static readonly Shareable: string = 'R'
  static readonly Mounted: string = 'M'
  static readonly Deletable: string = 'D'
  static readonly Renameable: string = 'N'
  static readonly Moveable: string = 'V'
  static readonly Updateable: string = 'NV'
  static readonly FileUpdateable: string = 'W'
  static readonly FolderCreateable: string = 'CK'
  static readonly Deny: string = 'Z'
  static readonly SecureView: string = 'X'
}

export type DavErrorCode =
  | 'ERR_LISTING_MEMBERS_NOT_ALLOWED'
  | 'ERR_INVALID_CREDENTIALS'
  | 'ERR_MISSING_BASIC_AUTH'
  | 'ERR_MISSING_BEARER_AUTH'
  | 'ERR_FILE_NOT_FOUND_IN_ROOT'

export enum DavMethod {
  copy = 'COPY',
  delete = 'DELETE',
  lock = 'LOCK',
  mkcol = 'MKCOL',
  move = 'MOVE',
  propfind = 'PROPFIND',
  proppatch = 'PROPPATCH',
  put = 'PUT',
  report = 'REPORT',
  unlock = 'UNLOCK'
}

type M<V, T> = {
  value: V
  type: T
}

const def = <V, T>(v: V): M<V, T> => ({
  value: v,
  type: null
})
const defStringOrNumber = <V>(v: V) => def<V, string | number>(v)
const defString = <V>(v: V) => def<V, string>(v)
const defNumber = <V>(v: V) => def<V, number>(v)
const defStringArray = <V>(v: V) => def<V, string[]>(v)

const DavPropertyMapping = {
  Permissions: defString('permissions' as const),
  IsFavorite: defNumber('favorite' as const),
  FileId: defString('fileid' as const),
  FileParent: defString('file-parent' as const),
  Name: defString('name' as const),
  OwnerId: defString('owner-id' as const),
  OwnerDisplayName: defString('owner-display-name' as const),
  PrivateLink: defString('privatelink' as const),
  ContentLength: defNumber('getcontentlength' as const),
  ContentSize: defNumber('size' as const),
  LastModifiedDate: defString('getlastmodified' as const),
  Tags: defStringOrNumber('tags' as const),
  Audio: {
    value: 'audio',
    type: null as Audio
  },
  Location: {
    value: 'location',
    type: null as GeoCoordinates
  },
  Image: {
    value: 'image',
    type: null as Image
  },
  Photo: {
    value: 'photo',
    type: null as Photo
  },
  ETag: defString('getetag' as const),
  MimeType: defString('getcontenttype' as const),
  ResourceType: defStringArray('resourcetype' as const),
  LockDiscovery: { value: 'lockdiscovery', type: null as Record<string, unknown> },
  LockOwner: defString('owner' as const),
  LockTime: defString('locktime' as const),
  ActiveLock: {
    value: 'activelock',
    type: null as Record<string, unknown>
  },
  DownloadURL: defString('downloadURL' as const),
  Highlights: defString('highlights' as const),
  MetaPathForUser: defString('meta-path-for-user' as const),
  RemoteItemId: defString('remote-item-id' as const),

  ShareId: defString('shareid' as const),
  ShareRoot: defString('shareroot' as const),
  ShareTypes: { value: 'share-types', type: null as Record<string, number[]> },
  SharePermissions: defString('share-permissions' as const),

  TrashbinOriginalFilename: defString('trashbin-original-filename' as const),
  TrashbinOriginalLocation: defString('trashbin-original-location' as const),
  TrashbinDeletedDate: defString('trashbin-delete-datetime' as const),

  PublicLinkItemType: defString('public-link-item-type' as const),
  PublicLinkPermission: defString('public-link-permission' as const),
  PublicLinkExpiration: defString('public-link-expiration' as const),
  PublicLinkShareDate: defString('public-link-share-datetime' as const),
  PublicLinkShareOwner: defString('public-link-share-owner' as const),
  SignatureAuth: {
    value: 'signature-auth',
    type: null as Record<'signature' | 'expiration', string>
  },
  SpaceId: defString('spaceid' as const)
} as const satisfies Record<string, M<unknown, unknown>>

type DavPropertyMappingType = typeof DavPropertyMapping

export const DavProperty = Object.fromEntries(
  Object.entries(DavPropertyMapping).map(([key, value]) => [key, value.value])
) as {
  [K in keyof DavPropertyMappingType as K]: DavPropertyMappingType[K]['value']
}

export type DavFileInfoResponse = {
  [K in keyof DavPropertyMappingType as DavPropertyMappingType[K]['value']]: DavPropertyMappingType[K]['type']
}

export type DavPropertyValue = (typeof DavProperty)[keyof typeof DavProperty]

export abstract class DavProperties {
  static readonly Default: DavPropertyValue[] = [
    DavProperty.Permissions,
    DavProperty.IsFavorite,
    DavProperty.FileId,
    DavProperty.FileParent,
    DavProperty.Name,
    DavProperty.LockDiscovery,
    DavProperty.ActiveLock,
    DavProperty.OwnerId,
    DavProperty.OwnerDisplayName,
    DavProperty.RemoteItemId,
    DavProperty.ShareRoot,
    DavProperty.ShareTypes,
    DavProperty.PrivateLink,
    DavProperty.ContentLength,
    DavProperty.ContentSize,
    DavProperty.LastModifiedDate,
    DavProperty.ETag,
    DavProperty.MimeType,
    DavProperty.ResourceType,
    DavProperty.DownloadURL,
    DavProperty.Tags,
    DavProperty.Audio,
    DavProperty.Location,
    DavProperty.Image,
    DavProperty.Photo,
    DavProperty.SpaceId
  ]

  static readonly PublicLink: DavPropertyValue[] = DavProperties.Default.concat([
    DavProperty.PublicLinkItemType,
    DavProperty.PublicLinkPermission,
    DavProperty.PublicLinkExpiration,
    DavProperty.PublicLinkShareDate,
    DavProperty.PublicLinkShareOwner,
    DavProperty.SignatureAuth
  ])

  static readonly Trashbin: DavPropertyValue[] = [
    DavProperty.ContentLength,
    DavProperty.ResourceType,
    DavProperty.TrashbinOriginalLocation,
    DavProperty.TrashbinOriginalFilename,
    DavProperty.TrashbinDeletedDate,
    DavProperty.Permissions,
    DavProperty.FileParent
  ]

  // these dav properties are dav standard and don't live in the oc namespace
  static readonly DavNamespace: DavPropertyValue[] = [
    DavProperty.ContentLength,
    DavProperty.LastModifiedDate,
    DavProperty.ETag,
    DavProperty.MimeType,
    DavProperty.ResourceType,
    DavProperty.LockDiscovery,
    DavProperty.ActiveLock
  ]
}
