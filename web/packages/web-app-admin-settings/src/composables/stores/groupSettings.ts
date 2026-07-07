import { defineStore } from 'pinia'
import { ref, unref } from 'vue'
import { Group } from '@ownclouders/web-client/graph/generated'

export const useGroupSettingsStore = defineStore('groupSettings', () => {
  const groups = ref<Group[]>([])
  const selectedGroups = ref<Group[]>([])

  const setGroups = (data: Group[]) => {
    groups.value = data
  }

  const upsertGroup = (group: Group) => {
    const existing = unref(groups).find(({ id }) => id === group.id)
    if (existing) {
      Object.assign(existing, group)
      return
    }
    unref(groups).push({ ...group, members: [] })
  }

  const removeGroups = (values: Group[]) => {
    groups.value = unref(groups).filter((group) => !values.find(({ id }) => id === group.id))
  }

  const setSelectedGroups = (data: Group[]) => {
    selectedGroups.value = data
  }

  const addSelectedGroup = (data: Group) => {
    unref(selectedGroups).push(data)
  }

  const reset = () => {
    groups.value = []
    selectedGroups.value = []
  }

  return {
    groups,
    upsertGroup,
    setGroups,
    removeGroups,
    reset,
    selectedGroups,
    addSelectedGroup,
    setSelectedGroups
  }
})

export type GroupSettingsStore = ReturnType<typeof useGroupSettingsStore>
