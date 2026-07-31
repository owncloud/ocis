import { defineStore } from 'pinia'
import { ref } from 'vue'

type Token = {
  id: string
  token: string
  link?: string
  expiration?: Date
  expirationSeconds?: number
  description?: string
}

export const useInviteTokensListStore = defineStore('inviteTokensList', () => {
  const tokensList = ref<Token[]>([])
  const lastCreatedToken = ref('')

  const setTokensList = (tokens: Token[]) => {
    tokensList.value = tokens
  }
  const addToken = (token: Token) => {
    tokensList.value.push(token)
  }
  const getTokensList = () => {
    return tokensList.value
  }

  const getLastCreatedToken = () => {
    return lastCreatedToken.value
  }
  const setLastCreatedToken = (token: string) => {
    lastCreatedToken.value = token
  }

  return {
    addToken,
    setTokensList,
    getTokensList,
    setLastCreatedToken,
    getLastCreatedToken
  }
})

export type inviteTokensListStore = ReturnType<typeof useInviteTokensListStore>
