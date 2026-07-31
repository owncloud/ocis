import { XMLBuilder } from 'fast-xml-parser'
import { DavProperties, DavPropertyValue } from '../constants'

const getNamespacedDavProps = (
  obj: Partial<Record<DavPropertyValue, unknown>>,
  extraProps: string[]
) => {
  return Object.fromEntries(
    Object.entries(obj).map(([name, value]) => {
      if (extraProps.includes(name)) {
        return [name, value || '']
      }

      const davNamespace = DavProperties.DavNamespace.includes(name as unknown as DavPropertyValue)
      const propName = davNamespace ? `d:${name}` : `oc:${name}`
      return [propName, value || '']
    })
  )
}

export const buildPropFindBody = (
  properties: DavPropertyValue[] = [],
  {
    pattern,
    filterRules,
    limit = 0,
    extraProps = []
  }: {
    pattern?: string
    filterRules?: Partial<Record<DavPropertyValue, unknown>>
    limit?: number
    extraProps: string[]
  }
): string => {
  let bodyType = 'd:propfind'
  if (pattern) {
    bodyType = 'oc:search-files'
  }

  if (filterRules) {
    bodyType = 'oc:filter-files'
  }

  const object = properties.reduce((obj, item) => Object.assign(obj, { [item]: null }), {})
  const props = getNamespacedDavProps(object, extraProps)

  const xmlObj = {
    [bodyType]: {
      'd:prop': props,
      '@@xmlns:d': 'DAV:',
      '@@xmlns:oc': 'http://owncloud.org/ns',
      ...(pattern && {
        'oc:search': { 'oc:pattern': pattern, 'oc:limit': limit }
      }),
      ...(filterRules && {
        'oc:filter-rules': getNamespacedDavProps(filterRules, [])
      })
    }
  }

  const builder = new XMLBuilder({
    format: true,
    ignoreAttributes: false,
    attributeNamePrefix: '@@',
    suppressEmptyNode: true
  })

  return builder.build(xmlObj)
}

export const buildPropPatchBody = (
  properties: Partial<Record<DavPropertyValue, unknown>>
): string => {
  const xmlObj = {
    'd:propertyupdate': {
      'd:set': { 'd:prop': getNamespacedDavProps(properties, []) },
      '@@xmlns:d': 'DAV:',
      '@@xmlns:oc': 'http://owncloud.org/ns'
    }
  }

  const builder = new XMLBuilder({
    format: true,
    ignoreAttributes: false,
    attributeNamePrefix: '@@',
    suppressEmptyNode: true
  })

  return builder.build(xmlObj)
}
