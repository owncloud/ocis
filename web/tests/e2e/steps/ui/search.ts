import { objects } from '../../support'
import { substitute } from '../../support/utils'
import { getWorld } from '../../environment/world'
import { searchFilter } from '../../environment/constants'

export async function userShouldSeeMessageOnSearchResult({
  stepUser,
  message
}: {
  stepUser: string
  message: string
}): Promise<boolean> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const searchObject = new objects.applicationFiles.Search({ page })
  const actualMessage = await searchObject.getSearchResultMessage()
  return actualMessage === substitute(message)
}

export async function userFiltersSearchResultWithTag({
  stepUser,
  tag
}: {
  stepUser: string
  tag: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const searchObject = new objects.applicationFiles.Search({ page })
  await searchObject.selectTagFilter({ tag })
}

export async function userClearsFilter({
  stepUser,
  filter
}: {
  stepUser: string
  filter:
    | typeof searchFilter.mediaType
    | typeof searchFilter.tags
    | typeof searchFilter.lastModified
    | typeof searchFilter.fullText
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const searchObject = new objects.applicationFiles.Search({ page })
  await searchObject.clearFilter({ filter })
}

export async function userEnablesTitleOnlySearch({
  stepUser
}: {
  stepUser: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const searchObject = new objects.applicationFiles.Search({ page })
  await searchObject.toggleSearchTitleOnly({ enableOrDisable: 'enable' })
}

export async function userFiltersSearchByMediaType({
  stepUser,
  mediaType
}: {
  stepUser: string
  mediaType: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const searchObject = new objects.applicationFiles.Search({ page })
  await searchObject.selectMediaTypeFilter({ mediaType })
}

export async function userFiltersSearchByLastModifiedDate({
  stepUser,
  lastModified
}: {
  stepUser: string
  lastModified: string
}): Promise<void> {
  const world = getWorld()
  const { page } = world.actorsEnvironment.getActor({ key: stepUser })
  const searchObject = new objects.applicationFiles.Search({ page })
  await searchObject.selectlastModifiedFilter({ lastModified })
}
