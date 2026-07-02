import { captureException } from '@sentry/vue'

export interface SettingsValue {
  identifier: {
    bundle: string
    extension: string
    setting: string
  }
  value: {
    accountUuid: string
    bundleId: string
    id: string
    resource: {
      type: string
    }
    settingId: string
    boolValue?: boolean
    listValue?: {
      values: {
        stringValue: string
      }[]
    }
    collectionValue?: {
      values: {
        key: string
        boolValue: boolean
      }[]
    }
    stringValue?: string
  }
}

interface SettingsBundleSetting {
  description: string
  displayName: string
  id: string
  name: string
  resource: {
    type: string
  }
  singleChoiceValue?: {
    options: Record<string, any>[]
  }
  multiChoiceCollectionValue?: {
    options: {
      value: {
        boolValue: {
          default?: boolean
        }
      }
      key: string
      displayValue: string
      attribute?: 'disabled'
    }[]
  }
  boolValue?: Record<string, any>
}

export interface SettingsBundle {
  displayName: string
  extension: string
  id: string
  name: string
  resource: {
    type: string
  }
  settings: SettingsBundleSetting[]
  type: string
  roleId?: string
}

export interface LanguageOption {
  label: string
  value: string
}

/** IDs of notifications setting bundles */
export enum SettingsNotificationBundle {
  ShareCreated = '872d8ef6-6f2a-42ab-af7d-f53cc81d7046',
  ShareRemoved = 'd7484394-8321-4c84-9677-741ba71e1f80',
  ShareExpired = 'e1aa0b7c-1b0f-4072-9325-c643c89fee4e',
  SpaceShared = '694d5ee1-a41c-448c-8d14-396b95d2a918',
  SpaceUnshared = '26c20e0e-98df-4483-8a77-759b3a766af0',
  SpaceMembershipExpired = '7275921e-b737-4074-ba91-3c2983be3edd',
  SpaceDisabled = 'eb5c716e-03be-42c6-9ed1-1105d24e109f',
  SpaceDeleted = '094ceca9-5a00-40ba-bb1a-bbc7bccd39ee',
  PostprocessingStepFinished = 'fe0a3011-d886-49c8-b797-33d02fa426ef',
  ScienceMeshInviteTokenGenerated = 'b441ffb1-f5ee-4733-a08f-48d03f6e7f22'
}

/** IDs of email notifications setting bundles */
export enum SettingsEmailNotificationBundle {
  EmailSendingInterval = '08dec2fe-3f97-42a9-9d1b-500855e92f25'
}

// We need the type specified here because e.g. includes method would otherwise complain about it
export const SETTINGS_NOTIFICATION_BUNDLE_IDS: string[] = [
  SettingsNotificationBundle.ShareCreated,
  SettingsNotificationBundle.ShareRemoved,
  SettingsNotificationBundle.ShareExpired,
  SettingsNotificationBundle.SpaceShared,
  SettingsNotificationBundle.SpaceUnshared,
  SettingsNotificationBundle.SpaceMembershipExpired,
  SettingsNotificationBundle.SpaceDisabled,
  SettingsNotificationBundle.SpaceDeleted,
  SettingsNotificationBundle.PostprocessingStepFinished,
  SettingsNotificationBundle.ScienceMeshInviteTokenGenerated
]

export const SETTINGS_EMAIL_NOTIFICATION_BUNDLE_IDS: string[] = [
  SettingsEmailNotificationBundle.EmailSendingInterval
]

function getSettingsDefaultValue(setting: SettingsBundleSetting) {
  if (setting.singleChoiceValue) {
    const [option] = setting.singleChoiceValue.options

    return {
      value: option.value.stringValue,
      displayValue: option.displayValue
    }
  }

  if (setting.multiChoiceCollectionValue) {
    return setting.multiChoiceCollectionValue.options.reduce((acc, curr) => {
      acc[curr.key] = curr.value.boolValue.default

      return acc
    }, {})
  }

  const error = new Error('Unsupported setting value')

  console.error(error)
  captureException(error)

  return null
}

export function getSettingsValue(
  setting: SettingsBundleSetting,
  valueList: SettingsValue[]
): boolean | string | null | { [key: string]: boolean } | { value: string; displayValue: string } {
  const { value } = valueList.find((v) => v.identifier.setting === setting.name) || {}

  if (!value) {
    return getSettingsDefaultValue(setting)
  }

  if (value.collectionValue) {
    return setting.multiChoiceCollectionValue.options.reduce((acc, curr) => {
      const val = value.collectionValue.values.find((v) => v.key === curr.key)

      if (val) {
        acc[curr.key] = val.boolValue
        return acc
      }

      acc[curr.key] = curr.value.boolValue.default
      return acc
    }, {})
  }

  if (value.stringValue) {
    const option = setting.singleChoiceValue.options.find(
      (o) => o.value.stringValue === value.stringValue
    )

    return { value: value.stringValue, displayValue: option?.displayValue || value.stringValue }
  }
}
