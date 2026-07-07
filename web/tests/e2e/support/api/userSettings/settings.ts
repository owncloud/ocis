import join from 'join-path'

import { checkResponseStatus, request } from '../http'
import { User } from '../../types'

interface BundleSetting {
  id: string
  name: string
  description: string
  boolValue?: {
    default?: boolean
    label?: string
  }
}

interface Setting {
  value: {
    identifier: {
      extension: string
      bundle: string
      setting: string
    }
    value: SettingValue
  }
}

interface SettingValue {
  id: string
  bundleId: string
  settingId: string
  accountUuid: string
  resource: {
    type: string
  }
  boolValue?: boolean
  stringValue?: string
  listValue?: {
    values: { stringValue: string }[]
  }
  collectionValue?: {
    values: { key: string; boolValue: boolean }[]
  }
}

const settings = {
  'auto-accept-shares': 'ec3ed4a3-3946-4efc-8f9f-76d38b12d3a9',
  language: 'aa8cfbe5-95d4-4f7e-a032-c3c01f5f062f'
}

const getSettingId = (setting: string): string => {
  return settings[setting as keyof typeof settings]
}

const getProfileBundleSettings = async (user: User) => {
  const response = await request({
    method: 'POST',
    path: join('api', 'v0', 'settings', 'bundle-get'),
    body: JSON.stringify({ bundleId: '2a506de7-99bd-4f0d-994e-c38e72c28fd9' }),
    user
  })
  checkResponseStatus(response, 'Failed get profile bundle')
  const { bundle } = (await response.json()) as { bundle: { settings: BundleSetting[] } }
  return bundle.settings
}

const getSingleProfileSetting = async (user: User, setting: string): Promise<BundleSetting> => {
  const settings = await getProfileBundleSettings(user)
  const settingObj = settings.find((s) => s.name === setting)
  if (!settingObj) {
    throw new Error(`Setting '${setting}' not found`)
  }
  return settingObj
}

export const configureAutoAcceptShare = async ({
  user,
  state
}: {
  user: User
  state: boolean
}): Promise<void> => {
  const body = JSON.stringify({
    value: {
      accountUuid: 'me',
      bundleId: '2a506de7-99bd-4f0d-994e-c38e72c28fd9',
      settingId: getSettingId('auto-accept-shares'),
      resource: {
        type: 'TYPE_USER'
      },
      boolValue: state
    }
  })
  const response = await request({
    method: 'POST',
    path: join('api', 'v0', 'settings', 'values-save'),
    body,
    user
  })
  checkResponseStatus(response, 'Failed while disabling auto-accept share')
}

export const changeLanguage = async ({
  user,
  language
}: {
  user: User
  language: string
}): Promise<void> => {
  const response = await request({
    method: 'PATCH',
    path: join('graph', 'v1.0', 'me'),
    body: JSON.stringify({ preferredLanguage: language }),
    user
  })
  checkResponseStatus(response, 'Failed change language: ' + language)
}

export const getSettingValue = async ({
  user,
  setting
}: {
  user: User
  setting: string
}): Promise<SettingValue | null> => {
  const body = JSON.stringify({
    accountUuid: 'me',
    settingId: getSettingId(setting)
  })

  const response = await request({
    method: 'POST',
    path: join('api', 'v0', 'settings', 'values-get-by-unique-identifiers'),
    body,
    user
  })

  if (response.status === 404) {
    return null
  }

  checkResponseStatus(response, 'Failed get setting: ' + setting)
  const settingValue = (await response.json()) as Setting
  return settingValue.value.value
}

export const getAutoAcceptSharesValue = async (user: User): Promise<boolean> => {
  const settingValue = await getSettingValue({ user, setting: 'auto-accept-shares' })
  if (settingValue === null) {
    // return default value
    const defaultSetting = await getSingleProfileSetting(user, 'auto-accept-shares')
    return defaultSetting.boolValue.default
  }
  return settingValue.boolValue
}
