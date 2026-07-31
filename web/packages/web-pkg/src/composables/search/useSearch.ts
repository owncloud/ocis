import { computed, unref } from 'vue'
import { SearchResult } from '../../components'
import { DavProperties } from '@ownclouders/web-client/webdav'
import { call, urlJoin } from '@ownclouders/web-client'
import { useClientService } from '../clientService'
import { isProjectSpaceResource } from '@ownclouders/web-client'
import {
  useCapabilityStore,
  useConfigStore,
  useResourcesStore,
  useSpacesStore
} from '../piniaStores'
import { SearchResource } from '@ownclouders/web-client'
import { useTask } from 'vue-concurrency'

export const useSearch = () => {
  const configStore = useConfigStore()
  const clientService = useClientService()
  const spacesStore = useSpacesStore()
  const resourcesStore = useResourcesStore()
  const capabilityStore = useCapabilityStore()

  const fullTextSearchEnabled = computed(() => capabilityStore.searchContent?.enabled)
  const areHiddenFilesShown = computed(() => resourcesStore.areHiddenFilesShown)
  const projectSpaces = computed(() => spacesStore.spaces.filter(isProjectSpaceResource))
  const getProjectSpace = (id: string) => {
    return unref(projectSpaces).find((s) => s.id === id)
  }

  const searchTask = useTask(function* (signal, term: string, searchLimit: number = null) {
    if (configStore.options.routing.fullShareOwnerPaths) {
      yield spacesStore.loadMountPoints({ graphClient: clientService.graphAuthenticated, signal })
    }

    if (!term) {
      return {
        totalResults: null,
        values: []
      }
    }

    const { resources, totalResults } = yield* call(
      clientService.webdav.search(term, {
        searchLimit,
        davProperties: DavProperties.Default,
        signal
      })
    )

    return {
      totalResults,
      values: resources
        .map((resource) => {
          const matchingSpace = getProjectSpace(resource.parentFolderId)
          const data = (matchingSpace ? matchingSpace : resource) as SearchResource

          if (configStore.options.routing.fullShareOwnerPaths && data.remoteItemPath) {
            data.path = urlJoin(data.remoteItemPath, data.path)
          }

          return { id: data.id, data }
        })
        .filter(({ data }) => {
          // filter results if hidden files shouldn't be shown due to settings
          return !data.name.startsWith('.') || unref(areHiddenFilesShown)
        })
    }
  }).restartable()

  const search = async (term: string, searchLimit: number = null): Promise<SearchResult> => {
    return await searchTask.perform(term, searchLimit)
  }

  const buildSearchTerm = ({
    term,
    isTitleOnlySearch,
    tags,
    lastModified,
    mediaType,
    scope,
    useScope,
    isVault
  }: {
    term: string
    isTitleOnlySearch?: boolean
    tags?: string
    lastModified?: string
    mediaType?: string
    scope?: string
    useScope?: boolean
    isVault?: boolean
  }) => {
    const query: string[] = []

    const humanSearchTerm = term
    const useFullTextSearch = unref(fullTextSearchEnabled) && !isTitleOnlySearch

    if (!!humanSearchTerm) {
      let nameQuery = `name:"*${humanSearchTerm}*"`

      if (useFullTextSearch) {
        nameQuery = `(name:"*${humanSearchTerm}*" OR content:"${humanSearchTerm}")`
      }

      query.push(nameQuery)
    }

    if (useScope && scope) {
      query.push(`scope:${scope}`)
    }

    if (tags) {
      const tagArr = tags.split('+').map((t) => `"${t}"`)
      query.push(`tag:(${tagArr.join(' OR ')})`)
    }

    if (lastModified) {
      query.push(`mtime:${lastModified}`)
    }

    if (mediaType) {
      const mediatypes = mediaType.split('+').map((t) => `"${t}"`)
      query.push(`mediatype:(${mediatypes.join(' OR ')})`)
    }

    if (isVault) {
      query.push('vault:true')
    }

    return query
      .sort((a, b) => Number(a.startsWith('scope:')) - Number(b.startsWith('scope:')))
      .join(' AND ')
  }

  return {
    search,
    buildSearchTerm
  }
}

export type SearchFunction = ReturnType<typeof useSearch>['search']
