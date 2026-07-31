import { computed, unref } from 'vue'
import { Resource } from '@ownclouders/web-client'
import { queryItemAsString } from '../appDefaults/useAppNavigation'
import { eventBus } from '../../services'
import { useRouteQuery } from '../router'
import { SideBarEventTopics } from '../sideBar'
import { useResourcesStore } from '../piniaStores'
import { isIncomingShareResource } from '@ownclouders/web-client'

export interface ScrollToResult {
  scrollToResource(
    resourceId: Resource['id'],
    options?: { forceScroll?: boolean; topbarElement: string }
  ): void
  scrollToResourceFromRoute(resources: Resource[], topbarElement: string): void
}

export const useScrollTo = (): ScrollToResult => {
  const scrollToQuery = useRouteQuery('scrollTo')
  const detailsQuery = useRouteQuery('details')
  const resourcesStore = useResourcesStore()

  const scrollTo = computed(() => {
    return queryItemAsString(unref(scrollToQuery))
  })
  const details = computed(() => {
    return queryItemAsString(unref(detailsQuery))
  })

  const scrollToResource = (
    resourceId: Resource['id'],
    options: { forceScroll?: boolean; topbarElement: string } = {
      forceScroll: false,
      topbarElement: null
    }
  ) => {
    const resourceElement = document.querySelectorAll(
      `[data-item-id='${resourceId}']`
    )[0] as HTMLElement

    if (!resourceElement) {
      eventBus.publish('app.files.navigate.page', {
        resourceId,
        forceScroll: options.forceScroll,
        topbarElement: options.topbarElement
      })
      return
    }

    // files-app-bar, admin-settings-app-bar
    const topbarElement = document.getElementById(options.topbarElement)
    // topbar height + th height + height of one row = offset needed when scrolling top
    const topOffset = topbarElement?.offsetHeight || 0 + resourceElement.offsetHeight * 2

    if (
      resourceElement.getBoundingClientRect().bottom > window.innerHeight ||
      resourceElement.getBoundingClientRect().top < topOffset ||
      options.forceScroll
    ) {
      resourceElement.scrollIntoView({ behavior: 'smooth', block: 'center' })
    }
  }

  eventBus.subscribe(
    'app.files.navigate.scrollTo',
    (data: { resourceId: string; forceScroll: boolean; topbarElement: string }) =>
      scrollToResource(data.resourceId, {
        forceScroll: data.forceScroll,
        topbarElement: data.topbarElement
      })
  )

  const scrollToResourceFromRoute = (resources: Resource[], topbarElement: null) => {
    if (!unref(scrollTo) || !resources.length) {
      return
    }

    const resource = unref(resources).find((r) => {
      if (isIncomingShareResource(r)) {
        return r.remoteItemId === unref(scrollTo)
      }
      return r.id === unref(scrollTo)
    })

    if (resource && resource.processing !== true) {
      resourcesStore.setSelection([resource.id])
      scrollToResource(resource.id, { forceScroll: true, topbarElement })

      if (unref(details)) {
        eventBus.publish(SideBarEventTopics.openWithPanel, unref(details))
      }
    }
  }

  return {
    scrollToResource,
    scrollToResourceFromRoute
  }
}
