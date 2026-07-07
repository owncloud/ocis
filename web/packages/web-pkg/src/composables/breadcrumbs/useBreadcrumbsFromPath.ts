import { eventBus } from '../../services/eventBus'
import { RouteLocation } from 'vue-router'
import omit from 'lodash-es/omit'
import { BreadcrumbItem } from '@ownclouders/design-system/helpers'
import { v4 as uuidV4 } from 'uuid'
import { SpaceResource } from '@ownclouders/web-client'
import { urlJoin } from '@ownclouders/web-client'
import { useGetMatchingSpace } from '../spaces'
import { Ref, ref, unref } from 'vue'
import { AncestorMetaData, AncestorMetaDataValue } from '../../types'

export const useBreadcrumbsFromPath = () => {
  const { isResourceAccessible } = useGetMatchingSpace()

  const breadcrumbsFromPath = ({
    route,
    space,
    resourcePath,
    ancestorMetaData = ref({})
  }: {
    route: RouteLocation
    space: Ref<SpaceResource>
    resourcePath: string
    ancestorMetaData?: Ref<AncestorMetaData>
  }): BreadcrumbItem[] => {
    const pathSplit = (p = '') => p.split('/').filter(Boolean)
    const current = pathSplit(route.path)
    const resource = pathSplit(resourcePath)

    return resource.map((text, i) => {
      const relativePath = urlJoin(...resource.slice(0, i + 1), { leadingSlash: true })
      const isAccessible = isResourceAccessible({ space: unref(space), path: relativePath })

      let ancestor: AncestorMetaDataValue
      if (isAccessible) {
        // use ancestor to retrieve fileId
        ancestor = unref(ancestorMetaData)[relativePath]
      }

      return {
        id: uuidV4(),
        allowContextActions: true,
        text,
        ...(isAccessible && {
          to: {
            path: '/' + [...current].splice(0, current.length - resource.length + i + 1).join('/'),
            query: {
              ...omit(route.query, 'page', 'fileId'),
              ...(ancestor && { fileId: ancestor.id })
            }
          }
        }),
        isStaticNav: false
      } as BreadcrumbItem
    })
  }

  const concatBreadcrumbs = (...items: BreadcrumbItem[]): BreadcrumbItem[] => {
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

  return { breadcrumbsFromPath, concatBreadcrumbs }
}
