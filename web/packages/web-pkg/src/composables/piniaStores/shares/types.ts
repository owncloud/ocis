import { Resource } from '@ownclouders/web-client'
import { ClientService } from '../../../services'
import { CollaboratorShare, LinkShare, SpaceResource } from '@ownclouders/web-client'
import { DriveItemCreateLink, DriveItemInvite } from '@ownclouders/web-client/graph/generated'

export interface AddShareOptions {
  clientService: ClientService
  space: SpaceResource
  resource: Resource
  options: DriveItemInvite
}

export interface UpdateShareOptions {
  clientService: ClientService
  space: SpaceResource
  resource: Resource
  collaboratorShare: CollaboratorShare
  options: DriveItemInvite
}

export interface DeleteShareOptions {
  clientService: ClientService
  space: SpaceResource
  resource: Resource
  collaboratorShare: CollaboratorShare
}

export interface AddLinkOptions {
  clientService: ClientService
  space: SpaceResource
  resource: Resource
  options: DriveItemCreateLink
}

export interface UpdateLinkOptions {
  clientService: ClientService
  space: SpaceResource
  resource: Resource
  linkShare: LinkShare
  options: Omit<DriveItemCreateLink, '@libre.graph.quickLink'>
}

export interface DeleteLinkOptions {
  clientService: ClientService
  space: SpaceResource
  resource: Resource
  linkShare: LinkShare
}
