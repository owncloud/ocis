import { checkResponseStatus, request } from '../http'
import { User } from '../../types'
import join from 'join-path'
import { getSpaceIdBySpaceName } from '../graph'
import convert from 'xml-js'
import _ from 'lodash-es/object.js'
import { createTagsForResource } from '../graph/utils'

export const folderExists = async ({
  user,
  path
}: {
  user: User
  path: string
}): Promise<boolean> => {
  const getResponse = await request({
    method: 'GET',
    path,
    user: user
  })

  return getResponse.status === 200
}

const createFolder = async ({
  user,
  folder,
  webDavEndPathToRoot // the root of the WebDAV path. This is `spaces/<space-id>` for ocis or `files/<user>` for oC10
}: {
  user: User
  folder: string
  webDavEndPathToRoot: string
}): Promise<void> => {
  const paths = folder.split('/')

  let parentFolder = ''
  for (const resource of paths) {
    const path = join('dav', webDavEndPathToRoot, parentFolder, resource)
    // check if the folder exists already or not
    const folderExist = await folderExists({ user, path })
    if (folderExist === false) {
      const response = await request({
        method: 'MKCOL',
        path,
        user: user
      })
      checkResponseStatus(response, 'Failed while creating folder')
    }
    parentFolder = join(parentFolder, resource)
  }
}

const createFile = async ({
  user,
  pathToFile,
  content,
  webDavEndPathToRoot, // the root of the WebDAV path. This is `spaces/<space-id>` for ocis or `files/<user>` for oC10
  mtimeDeltaDays
}: {
  user: User
  pathToFile: string
  content: string | Buffer
  webDavEndPathToRoot: string
  mtimeDeltaDays?: string
}): Promise<void> => {
  const today = new Date()
  const response = await request({
    method: 'PUT',
    path: join('dav', webDavEndPathToRoot, pathToFile),
    body: content,
    user: user,
    header: mtimeDeltaDays
      ? { 'X-OC-Mtime': today.getTime() / 1000 + parseInt(mtimeDeltaDays) * 86400 }
      : {}
  })

  checkResponseStatus(response, `Failed while uploading file '${pathToFile}' in personal space`)
}

export const uploadFileInPersonalSpace = async ({
  user,
  pathToFile,
  content,
  mtimeDeltaDays
}: {
  user: User
  pathToFile: string
  content: string | Buffer
  mtimeDeltaDays?: string
}): Promise<void> => {
  const webDavEndPathToRoot =
    'spaces/' + (await getSpaceIdBySpaceName({ user, spaceType: 'personal' }))
  await createFile({ user, pathToFile, content, webDavEndPathToRoot, mtimeDeltaDays })
}

export const createFolderInsideSpaceBySpaceName = async ({
  user,
  folder,
  spaceName
}: {
  user: User
  folder: string
  spaceName: string
}): Promise<void> => {
  const webDavEndPathToRoot =
    'spaces/' + (await getSpaceIdBySpaceName({ user, spaceType: 'project', spaceName }))
  await createFolder({ user, folder, webDavEndPathToRoot })
}

export const createFolderInsidePersonalSpace = async ({
  user,
  folder
}: {
  user: User
  folder: string
}): Promise<void> => {
  const webDavEndPathToRoot =
    'spaces/' + (await getSpaceIdBySpaceName({ user, spaceType: 'personal' }))
  await createFolder({ user, folder, webDavEndPathToRoot })
}

export const uploadFileInsideSpaceBySpaceName = async ({
  user,
  pathToFile,
  spaceName,
  content = ''
}: {
  user: User
  pathToFile: string
  spaceName: string
  content?: string | Buffer
}): Promise<void> => {
  const webDavEndPathToRoot =
    'spaces/' + (await getSpaceIdBySpaceName({ user, spaceType: 'project', spaceName }))
  await createFile({ user, pathToFile, content, webDavEndPathToRoot })
}

export const getDataOfFileInsideSpace = async ({
  user,
  pathToFileName,
  spaceType,
  spaceName
}: {
  user: User
  pathToFileName: string
  spaceType: string
  spaceName: string
}) => {
  const body =
    '<?xml version="1.0"?>\n' +
    '<d:propfind  xmlns:d="DAV:" xmlns:oc="http://owncloud.org/ns">\n' +
    '  <d:prop>\n' +
    '    <oc:permissions />\n' +
    '    <oc:favorite />\n' +
    '    <oc:fileid />\n' +
    '    <oc:file-parent />\n' +
    '    <oc:name />\n' +
    '    <oc:owner-id />\n' +
    '    <oc:owner-display-name />\n' +
    '    <oc:shareid />\n' +
    '    <oc:shareroot />\n' +
    '    <oc:share-types />\n' +
    '    <oc:privatelink />\n' +
    '    <d:getcontentlength />\n' +
    '    <oc:size />\n' +
    '    <d:getlastmodified />\n' +
    '    <d:getetag />\n' +
    '    <d:getcontenttype />\n' +
    '    <d:resourcetype />\n' +
    '    <oc:downloadURL />\n' +
    '  </d:prop>\n' +
    '</d:propfind>'
  const response = await request({
    method: 'PROPFIND',
    path: join(
      'dav',
      'spaces',
      await getSpaceIdBySpaceName({ user, spaceType, spaceName }),
      pathToFileName
    ),
    body: body,
    user: user
  })
  checkResponseStatus(response, `Failed while getting information of file ${pathToFileName}`)
  const fileData = JSON.parse(convert.xml2json(await response.text(), { compact: true }))
  return _.get(fileData, '[d:multistatus][d:response]')
}

export const getIdOfFileInsideSpace = async ({
  user,
  pathToFileName,
  spaceType,
  spaceName
}: {
  user: User
  pathToFileName: string
  spaceType: string
  spaceName: string
}): Promise<string> => {
  const fileDataResponse = await getDataOfFileInsideSpace({
    user,
    pathToFileName,
    spaceType,
    spaceName
  })
  // when there is a file inside the folder response comes as
  // [ [Object], [Object] ], so handel this case
  if (fileDataResponse.constructor.name === 'Array') {
    for (const key in fileDataResponse) {
      if (fileDataResponse[key]['d:propstat'][0]['d:prop']['oc:name']._text === pathToFileName) {
        return _.get(fileDataResponse[key], '[d:propstat][0][d:prop][oc:fileid]')._text
      }
    }
  } else {
    const propstat = _.get(fileDataResponse, '[d:propstat]')
    if (Array.isArray(propstat)) {
      return _.get(propstat, '[0][d:prop][oc:fileid]')._text
    }
    return _.get(propstat, '[d:prop][oc:fileid]')._text
  }
}

export const addTagToResource = async ({
  user,
  resource,
  tags
}: {
  user: User
  resource: string
  tags: string
}): Promise<void> => {
  const resourceId = await getIdOfFileInsideSpace({
    user,
    pathToFileName: resource,
    spaceType: 'personal',
    spaceName: user.displayName
  })
  const tagNames = tags.split(',').map((tag) => tag.trim())
  await createTagsForResource({ user, resourceId, tags: tagNames })
}

export const listSpaceResources = async ({
  user,
  spaceType,
  spaceName = '',
  folder = ''
}: {
  user: User
  spaceType: string
  spaceName?: string
  folder?: string
}): Promise<object> => {
  const propBody =
    '<?xml version="1.0"?>\n' +
    '<d:propfind  xmlns:d="DAV:" xmlns:oc="http://owncloud.org/ns">\n' +
    '  <d:prop>\n' +
    '    <oc:id />\n' +
    '    <oc:name />\n' +
    '  </d:prop>\n' +
    '</d:propfind>'
  const spaceId = await getSpaceIdBySpaceName({ user, spaceType, spaceName })

  const response = await request({
    method: 'PROPFIND',
    path: join('dav', 'spaces', spaceId, folder),
    body: propBody,
    user: user
  })

  checkResponseStatus(response, `Failed to list resources of space '${spaceName}'`)
  const xmlData = await response.text()

  const data = JSON.parse(convert.xml2json(xmlData, { compact: true }))
  let entries = _.get(data, '[d:multistatus][d:response]')
  if (!Array.isArray(entries)) {
    entries = [entries]
  }

  const resources = {}
  for (const entry of entries) {
    const resourceName = _.get(entry, '[d:propstat][d:prop][oc:name]')._text
    const resourceId = _.get(entry, '[d:propstat][d:prop][oc:id]')._text

    // do not include root folder
    if (!resourceId.endsWith(spaceId.split('$')[1])) {
      resources[resourceId] = resourceName
    }
  }
  return resources
}

export const deleteSpaceResource = async ({
  user,
  spaceType = '',
  spaceName = '',
  folder = '',
  fileId = ''
}: {
  user: User
  spaceType?: string
  spaceName?: string
  folder?: string
  fileId?: string
}): Promise<void> => {
  const path = ['dav', 'spaces']
  let errMessage = `Failed to delete resource '${folder}'`
  if (!fileId) {
    const spaceId = await getSpaceIdBySpaceName({ user, spaceType, spaceName })
    path.push(spaceId)
    path.push(folder)
  } else {
    path.push(fileId)
    errMessage = `Failed to delete resource with id '${fileId}'`
  }

  const response = await request({
    method: 'DELETE',
    path: join(...path),
    user: user
  })

  checkResponseStatus(response, errMessage)
}

export const emptyTrashbin = async ({
  user,
  spaceType,
  spaceName = ''
}: {
  user: User
  spaceType: string
  spaceName?: string
}): Promise<void> => {
  const spaceId = await getSpaceIdBySpaceName({ user, spaceType, spaceName })

  const response = await request({
    method: 'DELETE',
    path: join('dav', 'spaces', 'trash-bin', spaceId),
    user: user
  })

  checkResponseStatus(response, `Failed to empty trashbin of space '${spaceName}'`)
}
