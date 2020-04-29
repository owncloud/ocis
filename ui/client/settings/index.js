/* eslint-disable */
import axios from 'axios'
import qs from 'qs'
let domain = ''
export const getDomain = () => {
  return domain
}
export const setDomain = ($domain) => {
  domain = $domain
}
export const request = (method, url, body, queryParameters, form, config) => {
  method = method.toLowerCase()
  let keys = Object.keys(queryParameters)
  let queryUrl = url
  if (keys.length > 0) {
    queryUrl = url + '?' + qs.stringify(queryParameters)
  }
  // let queryUrl = url+(keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
  if (body) {
    return axios[method](queryUrl, body, config)
  } else if (method === 'get') {
    return axios[method](queryUrl, config)
  } else {
    return axios[method](queryUrl, qs.stringify(form), config)
  }
}
/*==========================================================
 *                    
 ==========================================================*/
/**
 * 
 * request: SaveSettingsBundle
 * url: SaveSettingsBundleURL
 * method: SaveSettingsBundle_TYPE
 * raw_url: SaveSettingsBundle_RAW_URL
 * @param body - 
 */
export const SaveSettingsBundle = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config
  let path = '/api/v0/settings/bundles'
  let body
  let queryParameters = {}
  let form = {}
  if (parameters['body'] !== undefined) {
    body = parameters['body']
  }
  if (parameters['body'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: body'))
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('post', domain + path, body, queryParameters, form, config)
}
export const SaveSettingsBundle_RAW_URL = function() {
  return '/api/v0/settings/bundles'
}
export const SaveSettingsBundle_TYPE = function() {
  return 'post'
}
export const SaveSettingsBundleURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/api/v0/settings/bundles'
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}
/**
 * 
 * request: ListSettingsBundles
 * url: ListSettingsBundlesURL
 * method: ListSettingsBundles_TYPE
 * raw_url: ListSettingsBundles_RAW_URL
 * @param identifierExtension - 
 * @param identifierBundleKey - 
 * @param identifierSettingKey - 
 * @param identifierAccountUuid - 
 */
export const ListSettingsBundles = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config
  let path = '/api/v0/settings/bundles/{identifier.extension}'
  let body
  let queryParameters = {}
  let form = {}
  path = path.replace('{identifier.extension}', `${parameters['identifierExtension']}`)
  if (parameters['identifierExtension'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: identifierExtension'))
  }
  if (parameters['identifierBundleKey'] !== undefined) {
    queryParameters['identifier.bundle_key'] = parameters['identifierBundleKey']
  }
  if (parameters['identifierSettingKey'] !== undefined) {
    queryParameters['identifier.setting_key'] = parameters['identifierSettingKey']
  }
  if (parameters['identifierAccountUuid'] !== undefined) {
    queryParameters['identifier.account_uuid'] = parameters['identifierAccountUuid']
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const ListSettingsBundles_RAW_URL = function() {
  return '/api/v0/settings/bundles/{identifier.extension}'
}
export const ListSettingsBundles_TYPE = function() {
  return 'get'
}
export const ListSettingsBundlesURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/api/v0/settings/bundles/{identifier.extension}'
  path = path.replace('{identifier.extension}', `${parameters['identifierExtension']}`)
  if (parameters['identifierBundleKey'] !== undefined) {
    queryParameters['identifier.bundle_key'] = parameters['identifierBundleKey']
  }
  if (parameters['identifierSettingKey'] !== undefined) {
    queryParameters['identifier.setting_key'] = parameters['identifierSettingKey']
  }
  if (parameters['identifierAccountUuid'] !== undefined) {
    queryParameters['identifier.account_uuid'] = parameters['identifierAccountUuid']
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}
/**
 * 
 * request: GetSettingsBundle
 * url: GetSettingsBundleURL
 * method: GetSettingsBundle_TYPE
 * raw_url: GetSettingsBundle_RAW_URL
 * @param identifierExtension - 
 * @param identifierBundleKey - 
 * @param identifierSettingKey - 
 * @param identifierAccountUuid - 
 */
export const GetSettingsBundle = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config
  let path = '/api/v0/settings/bundles/{identifier.extension}/{identifier.bundle_key}'
  let body
  let queryParameters = {}
  let form = {}
  path = path.replace('{identifier.extension}', `${parameters['identifierExtension']}`)
  if (parameters['identifierExtension'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: identifierExtension'))
  }
  path = path.replace('{identifier.bundle_key}', `${parameters['identifierBundleKey']}`)
  if (parameters['identifierBundleKey'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: identifierBundleKey'))
  }
  if (parameters['identifierSettingKey'] !== undefined) {
    queryParameters['identifier.setting_key'] = parameters['identifierSettingKey']
  }
  if (parameters['identifierAccountUuid'] !== undefined) {
    queryParameters['identifier.account_uuid'] = parameters['identifierAccountUuid']
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const GetSettingsBundle_RAW_URL = function() {
  return '/api/v0/settings/bundles/{identifier.extension}/{identifier.bundle_key}'
}
export const GetSettingsBundle_TYPE = function() {
  return 'get'
}
export const GetSettingsBundleURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/api/v0/settings/bundles/{identifier.extension}/{identifier.bundle_key}'
  path = path.replace('{identifier.extension}', `${parameters['identifierExtension']}`)
  path = path.replace('{identifier.bundle_key}', `${parameters['identifierBundleKey']}`)
  if (parameters['identifierSettingKey'] !== undefined) {
    queryParameters['identifier.setting_key'] = parameters['identifierSettingKey']
  }
  if (parameters['identifierAccountUuid'] !== undefined) {
    queryParameters['identifier.account_uuid'] = parameters['identifierAccountUuid']
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}
/**
 * 
 * request: SaveSettingsValue
 * url: SaveSettingsValueURL
 * method: SaveSettingsValue_TYPE
 * raw_url: SaveSettingsValue_RAW_URL
 * @param body - 
 */
export const SaveSettingsValue = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config
  let path = '/api/v0/settings/values'
  let body
  let queryParameters = {}
  let form = {}
  if (parameters['body'] !== undefined) {
    body = parameters['body']
  }
  if (parameters['body'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: body'))
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('post', domain + path, body, queryParameters, form, config)
}
export const SaveSettingsValue_RAW_URL = function() {
  return '/api/v0/settings/values'
}
export const SaveSettingsValue_TYPE = function() {
  return 'post'
}
export const SaveSettingsValueURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/api/v0/settings/values'
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}
/**
 * 
 * request: ListSettingsValues
 * url: ListSettingsValuesURL
 * method: ListSettingsValues_TYPE
 * raw_url: ListSettingsValues_RAW_URL
 * @param identifierAccountUuid - 
 * @param identifierExtension - 
 * @param identifierBundleKey - 
 * @param identifierSettingKey - 
 */
export const ListSettingsValues = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config
  let path = '/api/v0/settings/values/{identifier.account_uuid}/{identifier.extension}/{identifier.bundle_key}'
  let body
  let queryParameters = {}
  let form = {}
  path = path.replace('{identifier.account_uuid}', `${parameters['identifierAccountUuid']}`)
  if (parameters['identifierAccountUuid'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: identifierAccountUuid'))
  }
  path = path.replace('{identifier.extension}', `${parameters['identifierExtension']}`)
  if (parameters['identifierExtension'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: identifierExtension'))
  }
  path = path.replace('{identifier.bundle_key}', `${parameters['identifierBundleKey']}`)
  if (parameters['identifierBundleKey'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: identifierBundleKey'))
  }
  if (parameters['identifierSettingKey'] !== undefined) {
    queryParameters['identifier.setting_key'] = parameters['identifierSettingKey']
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const ListSettingsValues_RAW_URL = function() {
  return '/api/v0/settings/values/{identifier.account_uuid}/{identifier.extension}/{identifier.bundle_key}'
}
export const ListSettingsValues_TYPE = function() {
  return 'get'
}
export const ListSettingsValuesURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/api/v0/settings/values/{identifier.account_uuid}/{identifier.extension}/{identifier.bundle_key}'
  path = path.replace('{identifier.account_uuid}', `${parameters['identifierAccountUuid']}`)
  path = path.replace('{identifier.extension}', `${parameters['identifierExtension']}`)
  path = path.replace('{identifier.bundle_key}', `${parameters['identifierBundleKey']}`)
  if (parameters['identifierSettingKey'] !== undefined) {
    queryParameters['identifier.setting_key'] = parameters['identifierSettingKey']
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}
/**
 * 
 * request: GetSettingsValue
 * url: GetSettingsValueURL
 * method: GetSettingsValue_TYPE
 * raw_url: GetSettingsValue_RAW_URL
 * @param identifierAccountUuid - 
 * @param identifierExtension - 
 * @param identifierBundleKey - 
 * @param identifierSettingKey - 
 */
export const GetSettingsValue = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config
  let path = '/api/v0/settings/values/{identifier.account_uuid}/{identifier.extension}/{identifier.bundle_key}/{identifier.setting_key}'
  let body
  let queryParameters = {}
  let form = {}
  path = path.replace('{identifier.account_uuid}', `${parameters['identifierAccountUuid']}`)
  if (parameters['identifierAccountUuid'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: identifierAccountUuid'))
  }
  path = path.replace('{identifier.extension}', `${parameters['identifierExtension']}`)
  if (parameters['identifierExtension'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: identifierExtension'))
  }
  path = path.replace('{identifier.bundle_key}', `${parameters['identifierBundleKey']}`)
  if (parameters['identifierBundleKey'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: identifierBundleKey'))
  }
  path = path.replace('{identifier.setting_key}', `${parameters['identifierSettingKey']}`)
  if (parameters['identifierSettingKey'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: identifierSettingKey'))
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const GetSettingsValue_RAW_URL = function() {
  return '/api/v0/settings/values/{identifier.account_uuid}/{identifier.extension}/{identifier.bundle_key}/{identifier.setting_key}'
}
export const GetSettingsValue_TYPE = function() {
  return 'get'
}
export const GetSettingsValueURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/api/v0/settings/values/{identifier.account_uuid}/{identifier.extension}/{identifier.bundle_key}/{identifier.setting_key}'
  path = path.replace('{identifier.account_uuid}', `${parameters['identifierAccountUuid']}`)
  path = path.replace('{identifier.extension}', `${parameters['identifierExtension']}`)
  path = path.replace('{identifier.bundle_key}', `${parameters['identifierBundleKey']}`)
  path = path.replace('{identifier.setting_key}', `${parameters['identifierSettingKey']}`)
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}