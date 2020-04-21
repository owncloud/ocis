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
  let path = '/api/v0/bundles'
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
  return '/api/v0/bundles'
}
export const ListSettingsBundles_TYPE = function() {
  return 'get'
}
export const ListSettingsBundlesURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/api/v0/bundles'
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
  let path = '/api/v0/bundles'
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
  return '/api/v0/bundles'
}
export const CreateSettingsBundle_TYPE = function() {
  return 'post'
}
export const CreateSettingsBundleURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/api/v0/bundles'
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
 * @param key - 
 */
export const GetSettingsBundle = function(parameters = {}) {
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  const config = parameters.$config
  let path = '/api/v0/bundles/{extension}/{key}'
  let body
  let queryParameters = {}
  let form = {}
  path = path.replace('{extension}', `${parameters['extension']}`)
  if (parameters['extension'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: extension'))
  }
  path = path.replace('{key}', `${parameters['key']}`)
  if (parameters['key'] === undefined) {
    return Promise.reject(new Error('Missing required  parameter: key'))
  }
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    });
  }
  return request('get', domain + path, body, queryParameters, form, config)
}
export const GetSettingsBundle_RAW_URL = function() {
  return '/api/v0/bundles/{extension}/{key}'
}
export const GetSettingsBundle_TYPE = function() {
  return 'get'
}
export const GetSettingsBundleURL = function(parameters = {}) {
  let queryParameters = {}
  const domain = parameters.$domain ? parameters.$domain : getDomain()
  let path = '/api/v0/bundles/{extension}/{key}'
  path = path.replace('{extension}', `${parameters['extension']}`)
  path = path.replace('{key}', `${parameters['key']}`)
  if (parameters.$queryParameters) {
    Object.keys(parameters.$queryParameters).forEach(function(parameterName) {
      queryParameters[parameterName] = parameters.$queryParameters[parameterName]
    })
  }
  let keys = Object.keys(queryParameters)
  return domain + path + (keys.length > 0 ? '?' + (keys.map(key => key + '=' + encodeURIComponent(queryParameters[key])).join('&')) : '')
}