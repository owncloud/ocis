import { nextTick, computed, unref, Ref } from 'vue'
import { fileList } from '../../helpers/ui'
import {
  usePagination,
  useSort,
  SortDir,
  SortField,
  useRouteName,
  useResourcesStore,
  folderService
} from '@ownclouders/web-pkg'
import { useSideBar } from '@ownclouders/web-pkg'
import { queryItemAsString, useRouteQuery } from '@ownclouders/web-pkg'
import {
  determineResourceTableSortFields,
  determineResourceTilesSortFields,
  translateSortFields
} from '@ownclouders/web-pkg'
import { Task } from 'vue-concurrency'
import { Resource } from '@ownclouders/web-client'
import { useSelectedResources, SelectedResourcesResult } from '@ownclouders/web-pkg'
import { ReadOnlyRef } from '@ownclouders/web-pkg'
import {
  useFileListHeaderPosition,
  useViewMode,
  useViewSize,
  FolderViewModeConstants
} from '@ownclouders/web-pkg'

import { ScrollToResult, useScrollTo } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'

interface ResourcesViewDefaultsOptions<T, U extends any[]> {
  loadResourcesTask?: Task<T, U>
}

type ResourcesViewDefaultsResult<T extends Resource, TT, TU extends any[]> = {
  fileListHeaderY: Ref<any>
  refreshFileListHeaderPosition(): void
  loadResourcesTask: Task<TT, TU>
  areResourcesLoading: ReadOnlyRef<boolean>
  storeItems: ReadOnlyRef<T[]>
  sortFields: ReadOnlyRef<SortField[]>
  paginatedResources: Ref<T[]>
  paginationPages: ReadOnlyRef<number>
  paginationPage: ReadOnlyRef<number>
  handleSort({ sortBy, sortDir }: { sortBy: string; sortDir: SortDir }): void
  sortBy: ReadOnlyRef<string>
  sortDir: ReadOnlyRef<SortDir>
  viewMode: ReadOnlyRef<string>
  viewSize: ReadOnlyRef<number>
  selectedResources: Ref<Resource[]>
  selectedResourcesIds: Ref<string[]>
  isResourceInSelection(resource: Resource): boolean

  isSideBarOpen: Ref<boolean>
  sideBarActivePanel: Ref<string>
} & SelectedResourcesResult &
  ScrollToResult

export const useResourcesViewDefaults = <T extends Resource, TT, TU extends any[]>(
  options: ResourcesViewDefaultsOptions<TT, TU> = {}
): ResourcesViewDefaultsResult<T, TT, TU> => {
  const loadResourcesTask = options.loadResourcesTask || folderService.getTask()
  const areResourcesLoading = computed(() => {
    return loadResourcesTask.isRunning || !loadResourcesTask.last
  })

  const language = useGettext()
  const resourcesStore = useResourcesStore()
  const storeItems = computed(() => resourcesStore.activeResources) as unknown as Ref<T[]>

  const { refresh: refreshFileListHeaderPosition, y: fileListHeaderY } = useFileListHeaderPosition()

  const currentRoute = useRouteName()
  const currentViewModeQuery = useRouteQuery(
    `${unref(currentRoute)}-${FolderViewModeConstants.queryName}`,
    FolderViewModeConstants.defaultModeName
  )
  const currentViewMode = computed((): string => queryItemAsString(currentViewModeQuery.value))
  const viewMode = useViewMode(currentViewMode)

  const currentTilesSizeQuery = useRouteQuery('tiles-size', '1')
  const currentTilesSize = computed((): string => String(currentTilesSizeQuery.value))
  const viewSize = useViewSize(currentTilesSize)

  const sortFields = computed((): SortField[] => {
    if (unref(viewMode) === FolderViewModeConstants.name.tiles) {
      return translateSortFields(determineResourceTilesSortFields(unref(storeItems)[0]), language)
    }
    return determineResourceTableSortFields(unref(storeItems)[0])
  })

  const { sortBy, sortDir, items, handleSort } = useSort<T>({
    items: storeItems,
    fields: sortFields
  })
  const {
    items: paginatedResources,
    total: paginationPages,
    page: paginationPage
  } = usePagination<T>({ items, perPageStoragePrefix: 'files' })

  const accentuateItem = async (id: string) => {
    await nextTick()
    fileList.accentuateItem(id)
  }
  resourcesStore.$onAction((action) => {
    if (action.name === 'upsertResource') {
      accentuateItem(action.args[0].id)
    }
  })

  return {
    fileListHeaderY,
    refreshFileListHeaderPosition,
    loadResourcesTask,
    areResourcesLoading,
    storeItems,
    sortFields,
    viewMode,
    viewSize,
    paginatedResources,
    paginationPages,
    paginationPage,
    handleSort,
    sortBy,
    sortDir,
    ...useSelectedResources(),
    ...useSideBar(),
    ...useScrollTo()
  }
}
