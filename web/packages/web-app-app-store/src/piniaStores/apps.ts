import { defineStore } from 'pinia'
import { App, AppStoreRepository, RawAppListSchema } from '../types'
import { ref, unref } from 'vue'
import { APPID } from '../appid'
import { useRepositoriesStore } from './repositories'

export const useAppsStore = () => {
  const repositoriesStore = useRepositoriesStore()

  return defineStore(`${APPID}-apps`, () => {
    const apps = ref<App[]>([])

    const getById = (id: string) => {
      return unref(apps).find((app) => app.id === id)
    }

    const loadApps = async () => {
      const loadAppsByRepo = async (repo: AppStoreRepository): Promise<App[]> => {
        try {
          const data = await fetch(repo.url)
          const appsListData = await data.json()
          const appsList = RawAppListSchema.parse(appsListData)
          return appsList.apps
            .map((app) => {
              return {
                ...app,
                repository: repo,
                mostRecentVersion: app.versions[0]
              }
            })
            .sort((a, b) => {
              return a.name.toLowerCase().localeCompare(b.name.toLowerCase())
            })
        } catch (e) {
          console.error(e)
          return []
        }
      }

      const loadAppsPromises: Promise<App[]>[] = []
      for (const repo of repositoriesStore.repositories) {
        loadAppsPromises.push(loadAppsByRepo(repo))
      }
      apps.value = (await Promise.all(loadAppsPromises)).flat()
    }

    return {
      apps,
      getById,
      loadApps
    }
  })()
}
