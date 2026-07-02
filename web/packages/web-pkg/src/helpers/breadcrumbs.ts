import { eventBus } from '../services'
import { RouteLocation } from 'vue-router'
import omit from 'lodash-es/omit'
import { BreadcrumbItem } from '@ownclouders/design-system/helpers'
import { v4 as uuidV4 } from 'uuid'

export const breadcrumbsFromPath = (
  currentRoute: RouteLocation,
  resourcePath: string
): BreadcrumbItem[] => {
  const pathSplit = (p = '') => p.split('/').filter(Boolean)
  const current = pathSplit(currentRoute.path)
  const resource = pathSplit(resourcePath)

  return resource.map(
    (text, i) =>
      ({
        id: uuidV4(),
        allowContextActions: true,
        text,
        to: {
          path: '/' + [...current].splice(0, current.length - resource.length + i + 1).join('/'),
          query: omit(currentRoute.query, 'fileId', 'page') // TODO: we need the correct fileId in the query. until we have that we must omit it because otherwise we would correct the path to the one of the (wrong) fileId.
        },
        isStaticNav: false
      }) as BreadcrumbItem
  )
}

export const concatBreadcrumbs = (...items: BreadcrumbItem[]): BreadcrumbItem[] => {
  const last = items.pop()

  return [
    ...items,
    {
      id: uuidV4(),
      allowContextActions: last.allowContextActions,
      text: last.text,
      onClick: () => eventBus.publish('app.files.list.load'),
      isTruncationPlaceholder: last.isTruncationPlaceholder,
      isStaticNav: last.isStaticNav
    }
  ]
}
