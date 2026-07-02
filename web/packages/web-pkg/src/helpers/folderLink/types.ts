import { Resource } from '@ownclouders/web-client'

export interface CreateTargetRouteOptions {
  path: string
  fileId?: string | number
  resource: Resource
}
