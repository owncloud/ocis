import { Page } from '@playwright/test'
import { dav, graph, external } from '../../api'
import { config } from '../../../config'
import { User } from '../../types'
import { objects } from '../../index'

export interface navigateToDetailsPanelOfResourceArgs {
  page: Page
  resource: string
  detailsPanel: string
  user: User
  space: string
}

export interface openResourceViaUrlArgs {
  page: Page
  resource?: string
  user: User
  space?: string
  editorName?: string
  client?: string
}

export const navigateToDetailsPanelOfResource = async (
  args: navigateToDetailsPanelOfResourceArgs
): Promise<void> => {
  const { page, resource, detailsPanel, user, space } = args
  const fileId = await getTheFileIdOfSpaceFile(user, space, resource)
  const fullUrl = `${config.baseUrl}/f/${fileId}?details=${detailsPanel}`
  await page.goto(fullUrl)
}

export const openResourceViaUrl = async (args: openResourceViaUrlArgs) => {
  // accessibility check is not needed here because it opens in external app
  const { page, resource, user, space, editorName, client = '' } = args
  const fileId = await getTheFileIdOfSpaceFile(user, space, resource)
  let fullUrl

  switch (client) {
    case 'desktop':
      fullUrl = `${config.baseUrl}/external/open-with-web/?appName=${editorName}&fileId=${fileId}`
      break
    case 'mobile':
      fullUrl = await external.getOpenWithWebUrl({ user, fileId, editorName })
      break
    default:
      fullUrl = `${config.baseUrl}/f/${fileId}`
  }
  await page.goto(fullUrl)
}

export const openSpaceViaUrl = async (args: openResourceViaUrlArgs) => {
  const { page, user, space } = args
  let spaceName = null
  let spaceType = null
  if (space.toLowerCase() === 'personal') {
    spaceName = user.displayName
    spaceType = space.toLowerCase()
  } else {
    spaceName = space
    spaceType = 'project'
  }
  const fileId = await graph.getSpaceIdBySpaceName({ user, spaceType, spaceName })
  const fullUrl = `${config.baseUrl}/f/${fileId}`
  await page.goto(fullUrl)
}

const getTheFileIdOfSpaceFile = async (
  user: User,
  space: string,
  pathToFileName: string
): Promise<string> => {
  let spaceName = null
  let spaceType = null
  if (space.toLowerCase() === 'personal') {
    spaceName = user.displayName
    spaceType = space.toLowerCase()
  } else {
    spaceName = space
    spaceType = 'project'
  }
  return await dav.getIdOfFileInsideSpace({
    user,
    pathToFileName,
    spaceType,
    spaceName
  })
}

export const navigateToNonExistingPage = async ({ page }: { page: Page }) => {
  await page.goto(`${config.baseUrl}/'a-non-existing-page'`)
}
export const waitForNotFoundPageToBeVisible = async ({ page }: { page: Page }) => {
  await page.locator('.page-not-found').waitFor()
  await objects.a11y.Accessibility.assertNoSevereA11yViolations(
    page,
    ['pageNotFound'],
    'page-not-found page'
  )
}
