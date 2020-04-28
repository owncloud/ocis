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
 * request: ListSettingsBundles
 * url: ListSettingsBundlesURL
 * method: ListSettingsBundles_TYPE
 * raw_url: ListSettingsBundles_RAW_URL
 * @param extension - 
 */
export const ListSettingsBundles = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config
  let path = '/api/v0/settings/bundles'
  let body
  let queryParameters = {}
  let form = {}
  if (parameters['extension'] !== undefined) {
    queryParameters['extension'] = parameters['extension']
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const ListSettingsBundles_RAW_URL = function() {
  return '/api/v0/settings/bundles'
}
export const ListSettingsBundles_TYPE = function() {
  return 'get'
}
export const ListSettingsBundlesURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/api/v0/settings/bundles'
  if (parameters['extension'] !== undefined) {
    queryParameters['extension'] = parameters['extension']
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
 * request: CreateSettingsBundle
 * url: CreateSettingsBundleURL
 * method: CreateSettingsBundle_TYPE
 * raw_url: CreateSettingsBundle_RAW_URL
 * @param body - 
 */
export const CreateSettingsBundle = function(parameters = {}) {
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
export const CreateSettingsBundle_RAW_URL = function() {
  return '/api/v0/settings/bundles'
}
export const CreateSettingsBundle_TYPE = function() {
  return 'post'
}
export const CreateSettingsBundleURL = function(parameters = {}) {
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
 * request: GetSettingsBundle
 * url: GetSettingsBundleURL
 * method: GetSettingsBundle_TYPE
 * raw_url: GetSettingsBundle_RAW_URL
 * @param extension - 
 * @param bundleKey - 
 */
export const GetSettingsBundle = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config
  let path = '/api/v0/settings/bundles/{extension}/{bundle_key}'
  let body
  let queryParameters = {}
  let form = {}
  path = path.replace('{extension}', `${parameters['extension']}`)
  if (parameters['extension'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: extension'))
  }
  path = path.replace('{bundle_key}', `${parameters['bundleKey']}`)
  if (parameters['bundleKey'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: bundleKey'))
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const GetSettingsBundle_RAW_URL = function() {
  return '/api/v0/settings/bundles/{extension}/{bundle_key}'
}
export const GetSettingsBundle_TYPE = function() {
  return 'get'
}
export const GetSettingsBundleURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/api/v0/settings/bundles/{extension}/{bundle_key}'
  path = path.replace('{extension}', `${parameters['extension']}`)
  path = path.replace('{bundle_key}', `${parameters['bundleKey']}`)
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
 * request: GetSettingsValue
 * url: GetSettingsValueURL
 * method: GetSettingsValue_TYPE
 * raw_url: GetSettingsValue_RAW_URL
 * @param accountUuid - 
 * @param extension - 
 * @param bundleKey - 
 * @param settingKey - 
 */
export const GetSettingsValue = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config
  let path = '/api/v0/settings/values/{account_uuid}/{extension}/{bundle_key}/{setting_key}'
  let body
  let queryParameters = {}
  let form = {}
  path = path.replace('{account_uuid}', `${parameters['accountUuid']}`)
  if (parameters['accountUuid'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: accountUuid'))
  }
  path = path.replace('{extension}', `${parameters['extension']}`)
  if (parameters['extension'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: extension'))
  }
  path = path.replace('{bundle_key}', `${parameters['bundleKey']}`)
  if (parameters['bundleKey'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: bundleKey'))
  }
  path = path.replace('{setting_key}', `${parameters['settingKey']}`)
  if (parameters['settingKey'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: settingKey'))
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const GetSettingsValue_RAW_URL = function() {
  return '/api/v0/settings/values/{account_uuid}/{extension}/{bundle_key}/{setting_key}'
}
export const GetSettingsValue_TYPE = function() {
  return 'get'
}
export const GetSettingsValueURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/api/v0/settings/values/{account_uuid}/{extension}/{bundle_key}/{setting_key}'
  path = path.replace('{account_uuid}', `${parameters['accountUuid']}`)
  path = path.replace('{extension}', `${parameters['extension']}`)
  path = path.replace('{bundle_key}', `${parameters['bundleKey']}`)
  path = path.replace('{setting_key}', `${parameters['settingKey']}`)
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}