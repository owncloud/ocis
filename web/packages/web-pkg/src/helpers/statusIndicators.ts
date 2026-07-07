import { ShareTypes } from '@ownclouders/web-client'
import { eventBus } from '../services'
import { SideBarEventTopics } from '../composables/sideBar'
import { useExtensionRegistry } from '../composables/piniaStores/extensionRegistry'
import { Resource } from '@ownclouders/web-client'
import { AncestorMetaData } from '../types'
import {
  SpaceResource,
  isPersonalSpaceResource,
  isProjectSpaceResource
} from '@ownclouders/web-client'
import { User } from '@ownclouders/web-client/graph/generated'
import { IconFillType } from './resource'
import { resourceIndicatorExtensionPoint } from '../extensionPoints'

// dummy to trick gettext string extraction into recognizing strings
const $gettext = (str: string): string => {
  return str
}

export type ResourceIndicatorCategory = 'system' | 'sharing'

export interface ResourceIndicator {
  id: string
  accessibleDescription: string
  label: string
  icon: string
  fillType: IconFillType
  type: string
  category: ResourceIndicatorCategory
  handler?: (resource: Resource) => void
}

const isUserShare = (shareTypes: number[]) => {
  return ShareTypes.containsAnyValue(ShareTypes.authenticated, shareTypes ?? [])
}

const isLinkShare = (shareTypes: number[]) => {
  return ShareTypes.containsAnyValue(ShareTypes.unauthenticated, shareTypes ?? [])
}

const shareUserIconDescribedBy = ({ isDirect }: { isDirect: boolean }) => {
  return isDirect
    ? $gettext('This item is directly shared with others.')
    : $gettext('This item is shared with others through one of the parent folders.')
}

const shareLinkDescribedBy = ({ isDirect }: { isDirect: boolean }) => {
  return isDirect
    ? $gettext('This item is directly shared via links.')
    : $gettext('This item is shared via links through one of the parent folders.')
}

const getUserIndicator = ({
  resource,
  isDirect
}: {
  resource: Resource
  isDirect: boolean
}): ResourceIndicator => {
  return {
    id: `files-sharing-${resource.getDomSelector()}`,
    accessibleDescription: shareUserIconDescribedBy({ isDirect }),
    label: $gettext('Show invited people'),
    icon: 'group',
    category: 'sharing',
    type: isDirect ? 'user-direct' : 'user-indirect',
    fillType: 'line',
    handler: () => {
      eventBus.publish(SideBarEventTopics.openWithPanel, 'sharing#peopleShares')
    }
  }
}

const getLinkIndicator = ({
  resource,
  isDirect
}: {
  resource: Resource
  isDirect: boolean
}): ResourceIndicator => {
  return {
    id: `file-link-${resource.getDomSelector()}`,
    accessibleDescription: shareLinkDescribedBy({ isDirect }),
    label: $gettext('Show links'),
    icon: 'link',
    category: 'sharing',
    type: isDirect ? 'link-direct' : 'link-indirect',
    fillType: 'line',
    handler: () => {
      eventBus.publish(SideBarEventTopics.openWithPanel, 'sharing#linkShares')
    }
  }
}

const getLockedIndicator = ({ resource }: { resource: Resource }): ResourceIndicator => {
  return {
    id: `resource-locked-${resource.getDomSelector()}`,
    accessibleDescription: $gettext('Item locked'),
    label: $gettext('This item is locked'),
    icon: 'lock',
    category: 'system',
    type: 'resource-locked',
    fillType: 'line'
  }
}

const getProcessingIndicator = ({ resource }: { resource: Resource }): ResourceIndicator => {
  return {
    id: `resource-processing-${resource.getDomSelector()}`,
    accessibleDescription: $gettext('Item in processing'),
    label: $gettext('This item is in processing'),
    icon: 'loop-right',
    category: 'system',
    type: 'resource-processing',
    fillType: 'line'
  }
}

export const getIndicators = ({
  space,
  resource,
  ancestorMetaData,
  user
}: {
  space: SpaceResource
  resource: Resource
  ancestorMetaData: AncestorMetaData
  user: User
}): ResourceIndicator[] => {
  const indicators: ResourceIndicator[] = []

  if (resource.locked) {
    indicators.push(getLockedIndicator({ resource }))
  }

  if (resource.processing) {
    indicators.push(getProcessingIndicator({ resource }))
  }

  const shareIndicatorsAllowed =
    (isProjectSpaceResource(space) && space.isMember(user)) ||
    (isPersonalSpaceResource(space) && space.isOwner(user))

  if (shareIndicatorsAllowed) {
    const parentShareTypes = Object.values(ancestorMetaData).reduce<number[]>((acc, data) => {
      acc.push(...(data.shareTypes || []))
      return acc
    }, [])

    const isDirectUserShare = isUserShare(resource.shareTypes)
    if (isDirectUserShare || isUserShare(parentShareTypes)) {
      indicators.push(getUserIndicator({ resource, isDirect: isDirectUserShare }))
    }

    const isDirectLinkShare = isLinkShare(resource.shareTypes)
    if (isDirectLinkShare || isLinkShare(parentShareTypes)) {
      indicators.push(getLinkIndicator({ resource, isDirect: isDirectLinkShare }))
    }
  }

  ;(useExtensionRegistry().requestExtensions(resourceIndicatorExtensionPoint) || []).forEach(
    (extension) => {
      const extensionIndicators = extension.getResourceIndicators(resource)
      if (extensionIndicators) {
        indicators.push(...extensionIndicators)
      }
    }
  )

  return indicators
}
