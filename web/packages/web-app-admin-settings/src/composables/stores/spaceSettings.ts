import { defineStore } from 'pinia'
import { ref, unref } from 'vue'
import { SpaceResource } from '@ownclouders/web-client'

export const useSpaceSettingsStore = defineStore('spaceSettings', () => {
  const spaces = ref<SpaceResource[]>([])
  const selectedSpaces = ref<SpaceResource[]>([])

  const setSpaces = (data: SpaceResource[]) => {
    spaces.value = data
  }

  const upsertSpace = (space: SpaceResource) => {
    const existing = unref(spaces).find(({ id }) => id === space.id)
    if (existing) {
      Object.assign(existing, space)
      return
    }
    unref(spaces).push(space)
  }

  const removeSpaces = (values: SpaceResource[]) => {
    spaces.value = unref(spaces).filter((space) => !values.find(({ id }) => id === space.id))
  }

  const setSelectedSpaces = (data: SpaceResource[]) => {
    selectedSpaces.value = data
  }

  const addSelectedSpace = (data: SpaceResource) => {
    unref(selectedSpaces).push(data)
  }

  const reset = () => {
    spaces.value = []
    selectedSpaces.value = []
  }

  return {
    spaces,
    setSpaces,
    upsertSpace,
    removeSpaces,
    reset,
    selectedSpaces,
    addSelectedSpace,
    setSelectedSpaces
  }
})

export type SpaceSettingsStore = ReturnType<typeof useSpaceSettingsStore>
