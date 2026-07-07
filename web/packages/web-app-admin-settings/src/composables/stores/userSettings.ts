import { defineStore } from 'pinia'
import { ref, unref } from 'vue'
import { User } from '@ownclouders/web-client/graph/generated'

export const useUserSettingsStore = defineStore('userSettings', () => {
  const users = ref<User[]>([])
  const selectedUsers = ref<User[]>([])

  const setUsers = (data: User[]) => {
    users.value = data
  }

  const upsertUser = (user: User) => {
    const existing = unref(users).find(({ id }) => id === user.id)
    if (existing) {
      Object.assign(existing, user)
      return
    }
    unref(users).push(user)
  }

  const removeUsers = (values: User[]) => {
    users.value = unref(users).filter((user) => !values.find(({ id }) => id === user.id))
  }

  const setSelectedUsers = (data: User[]) => {
    selectedUsers.value = data
  }

  const addSelectedUser = (data: User) => {
    unref(selectedUsers).push(data)
  }

  const reset = () => {
    users.value = []
    selectedUsers.value = []
  }

  return {
    users,
    setUsers,
    upsertUser,
    removeUsers,
    reset,
    selectedUsers,
    addSelectedUser,
    setSelectedUsers
  }
})

export type UserSettingsStore = ReturnType<typeof useUserSettingsStore>
