import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string>()
  const idpContextReady = ref(false)
  const userContextReady = ref(false)
  const publicLinkToken = ref<string>()
  const publicLinkPassword = ref<string>()
  const publicLinkPasswordRequired = ref<boolean>()
  const publicLinkType = ref<string>()
  const publicLinkContextReady = ref(false)

  const setAccessToken = (value: string) => {
    accessToken.value = value
  }
  const setIdpContextReady = (value: boolean) => {
    idpContextReady.value = value
  }
  const setUserContextReady = (value: boolean) => {
    userContextReady.value = value
  }
  const setPublicLinkContext = (context: {
    publicLinkToken: string
    publicLinkPassword: string
    publicLinkType: string
    publicLinkContextReady: boolean
    publicLinkPasswordRequired: boolean
  }) => {
    publicLinkToken.value = context.publicLinkToken
    publicLinkPassword.value = context.publicLinkPassword
    publicLinkType.value = context.publicLinkType
    publicLinkContextReady.value = context.publicLinkContextReady
    publicLinkPasswordRequired.value = context.publicLinkPasswordRequired
  }

  const clearUserContext = () => {
    setAccessToken(null)
    setIdpContextReady(null)
    setUserContextReady(null)
  }

  const clearPublicLinkContext = () => {
    setPublicLinkContext({
      publicLinkToken: null,
      publicLinkPassword: null,
      publicLinkType: null,
      publicLinkContextReady: false,
      publicLinkPasswordRequired: null
    })
  }

  return {
    accessToken,
    idpContextReady,
    userContextReady,
    publicLinkToken,
    publicLinkPassword,
    publicLinkType,
    publicLinkContextReady,
    publicLinkPasswordRequired,

    setAccessToken,
    setIdpContextReady,
    setUserContextReady,
    setPublicLinkContext,
    clearUserContext,
    clearPublicLinkContext
  }
})

export type AuthStore = ReturnType<typeof useAuthStore>
