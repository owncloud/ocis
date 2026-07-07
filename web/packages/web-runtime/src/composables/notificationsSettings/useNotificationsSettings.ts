import { computed, Ref, unref } from 'vue'
import {
  getSettingsValue,
  SETTINGS_EMAIL_NOTIFICATION_BUNDLE_IDS,
  SETTINGS_NOTIFICATION_BUNDLE_IDS,
  SettingsBundle,
  SettingsValue
} from '../../helpers/settings'

export const useNotificationsSettings = (
  valueList: Ref<SettingsValue[]>,
  bundle: Ref<SettingsBundle>
) => {
  const values = computed(() => {
    if (!unref(bundle)) {
      return {}
    }

    return unref(bundle).settings.reduce((acc, curr) => {
      if (!SETTINGS_NOTIFICATION_BUNDLE_IDS.includes(curr.id)) {
        return acc
      }

      acc[curr.id] = getSettingsValue(curr, unref(valueList))

      return acc
    }, {})
  })

  const options = computed<SettingsBundle['settings']>(() => {
    if (!unref(bundle)) {
      return []
    }

    return unref(bundle).settings.filter(({ id }) => SETTINGS_NOTIFICATION_BUNDLE_IDS.includes(id))
  })

  const emailOptions = computed<SettingsBundle['settings']>(() => {
    if (!unref(bundle)) {
      return []
    }

    return unref(bundle).settings.filter(({ id }) =>
      SETTINGS_EMAIL_NOTIFICATION_BUNDLE_IDS.includes(id)
    )
  })

  const emailValues = computed(() => {
    if (!unref(bundle)) {
      return {}
    }

    return unref(bundle).settings.reduce((acc, curr) => {
      if (!SETTINGS_EMAIL_NOTIFICATION_BUNDLE_IDS.includes(curr.id)) {
        return acc
      }

      acc[curr.id] = getSettingsValue(curr, unref(valueList))

      return acc
    }, {})
  })

  return { values, options, emailOptions, emailValues }
}
