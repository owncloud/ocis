import { defineStore } from 'pinia'
import { AppStoreRepository } from '../types'
import { ref } from 'vue'
import { APPID } from '../appid'

export const useRepositoriesStore = defineStore(`${APPID}-repositories`, () => {
  const repositories = ref<AppStoreRepository[]>([])

  const setRepositories = (repos: AppStoreRepository[]) => {
    repositories.value = repos
  }

  return {
    repositories,
    setRepositories
  }
})
