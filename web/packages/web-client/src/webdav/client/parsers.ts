import { parseXML, prepareFileFromProps } from 'webdav'
import { XMLParser } from 'fast-xml-parser'
import { WebDavResponseResource, WebDavResponseTusSupport } from '../../helpers'
import { urlJoin } from '../../utils'
import { DavErrorCode } from '../constants'
import { join, normalize } from 'path'

export const parseTusHeaders = (headers: Headers) => {
  const result: WebDavResponseTusSupport = {}

  const version = headers.get('tus-version')
  if (!version) {
    return null
  }

  result.version = version.split(',')
  if (headers.get('tus-extension')) {
    result.extension = headers.get('tus-extension').split(',')
  }
  if (headers.get('tus-resumable')) {
    result.resumable = headers.get('tus-resumable')
  }
  if (headers.get('tus-max-size')) {
    result.maxSize = parseInt(headers.get('tus-max-size'), 10)
  }
  return result
}

export const parseMultiStatus = async (xmlBody: string, remoteBasePath: string) => {
  const parseFileName = (name: string) => {
    const decoded = decodeURIComponent(name)
    const prefix = normalize(join(remoteBasePath, 'dav'))
    if (name?.startsWith(prefix)) {
      // strip out '/dav/' from the beginning
      return urlJoin(decoded.replace(prefix, ''), {
        leadingSlash: true,
        trailingSlash: false
      })
    }

    return decoded
  }

  const parsedXML = await parseXML(xmlBody)

  return parsedXML.multistatus.response.map(({ href, propstat }) => {
    const data = {
      ...prepareFileFromProps(propstat.prop, parseFileName(href), true),
      processing: propstat.status === 'HTTP/1.1 425 TOO EARLY'
    }

    if (data.props.name) {
      data.props.name = data.props.name.toString()
    }

    return data
  }) as unknown as WebDavResponseResource[]
}

export const parseError = (xmlBody: string): { message: string; errorCode: DavErrorCode } => {
  const parser = new XMLParser()
  const errorObj = { message: 'Unknown error', errorCode: undefined }

  try {
    const parsed = parser.parse(xmlBody)
    if (!parsed['d:error']) {
      return errorObj
    }
    if (parsed['d:error']['s:message']) {
      const message = parsed['d:error']['s:message']
      if (typeof message === 'string') {
        errorObj.message = message
      }
    }
    if (parsed['d:error']['s:errorcode']) {
      const errorCode = parsed['d:error']['s:errorcode']
      if (typeof errorCode === 'string') {
        errorObj.errorCode = errorCode
      }
    }
  } catch {
    return errorObj
  }

  return errorObj
}
