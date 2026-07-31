import { setUser as sentrySetUser } from '@sentry/vue'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { User } from '@ownclouders/web-client/graph/generated'

export const useUserStore = defineStore('user', () => {
  const user = ref<User>()

  const setUser = (data: User) => {
    user.value = data
    sentrySetUser({ username: data.onPremisesSamAccountName })
  }

  const reset = () => {
    user.value = null
  }

  return {
    user,
    setUser,
    reset
  }
})

export type UserStore = ReturnType<typeof useUserStore>
