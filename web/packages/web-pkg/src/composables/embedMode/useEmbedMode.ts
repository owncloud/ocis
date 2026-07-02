import { computed, unref } from 'vue'
import { useConfigStore } from '../piniaStores'
import { Resource } from '@ownclouders/web-client'
import { LocationQuery } from '../router'

export interface embedModeFilePickMessageData {
  resource: Resource
  locationQuery: LocationQuery
}

export interface embedModeLocationPickMessageData {
  resources: Resource[]
  fileName?: string
  locationQuery?: LocationQuery
}

export const useEmbedMode = () => {
  const configStore = useConfigStore()

  const isEnabled = computed(() => configStore.options.embed?.enabled)

  const isLocationPicker = computed(() => {
    return configStore.options.embed?.target === 'location'
  })

  const chooseFileName = computed(() => {
    return configStore.options.embed?.chooseFileName
  })

  const chooseFileNameSuggestion = computed(() => {
    return configStore.options.embed?.chooseFileNameSuggestion
  })

  const isFilePicker = computed(() => {
    return configStore.options.embed?.target === 'file'
  })

  const fileTypes = computed(() => {
    return configStore.options.embed?.fileTypes
  })

  const messagesTargetOrigin = computed(() => configStore.options.embed?.messagesOrigin)

  const isDelegatingAuthentication = computed(
    () => unref(isEnabled) && configStore.options.embed?.delegateAuthentication
  )

  const delegateAuthenticationOrigin = computed(
    () => configStore.options.embed?.delegateAuthenticationOrigin
  )

  const postMessage = <Payload>(name: string, data?: Payload): void => {
    const options: WindowPostMessageOptions = {}

    if (unref(messagesTargetOrigin)) {
      options.targetOrigin = unref(messagesTargetOrigin)
    }

    window.parent.postMessage({ name, data }, options)
  }

  const verifyDelegatedAuthenticationOrigin = (eventOrigin: string): boolean => {
    if (!unref(delegateAuthenticationOrigin)) {
      return true
    }

    return unref(delegateAuthenticationOrigin) === eventOrigin
  }

  const verifyMessageOrigin = (eventOrigin: string): boolean => {
    const allowedOrigins = [window.location.origin]

    if (unref(messagesTargetOrigin)) {
      allowedOrigins.push(unref(messagesTargetOrigin))
    }

    return allowedOrigins.includes(eventOrigin)
  }

  return {
    isEnabled,
    isLocationPicker,
    chooseFileName,
    chooseFileNameSuggestion,
    isFilePicker,
    messagesTargetOrigin,
    isDelegatingAuthentication,
    fileTypes,
    postMessage,
    verifyDelegatedAuthenticationOrigin,
    verifyMessageOrigin
  }
}
