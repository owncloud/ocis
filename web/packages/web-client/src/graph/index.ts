import { Configuration } from './generated'
import { FetchClient } from '../http'
import { undeclaredParams } from './types'
import { type GraphUsers, UsersFactory } from './users'
import { type GraphGroups, GroupsFactory } from './groups'
import { ApplicationsFactory, GraphApplications } from './applications'
import { DrivesFactory, GraphDrives } from './drives'
import { DriveItemsFactory, GraphDriveItems } from './driveItems'
import { TagsFactory, GraphTags } from './tags'
import { ActivitiesFactory, GraphActivities } from './activities'
import { PermissionsFactory, GraphPermissions } from './permissions'

export interface Graph {
  activities: GraphActivities
  applications: GraphApplications
  tags: GraphTags
  drives: GraphDrives
  driveItems: GraphDriveItems
  users: GraphUsers
  groups: GraphGroups
  permissions: GraphPermissions
}

export const graph = (baseURI: string, httpClient: FetchClient): Graph => {
  const url = new URL(baseURI)
  url.pathname = [...url.pathname.split('/'), 'graph'].filter(Boolean).join('/')
  const config = new Configuration({
    basePath: url.href,
    // Route every generated request through the core so header injection, maintenance
    // detection and HttpError-on-non-2xx all apply. Because the core throws on non-2xx,
    // callers keep seeing HttpError (with `statusCode` and `data`) rather than the
    // generated ResponseError.
    fetchApi: (input: RequestInfo | URL, init?: RequestInit) => {
      const params = (init as Record<symbol, Record<string, string>>)?.[undeclaredParams]
      return httpClient.fetch(String(input), {
        method: init?.method,
        headers: Object.fromEntries(new Headers(init?.headers).entries()),
        body: init?.body,
        signal: init?.signal ?? undefined,
        ...(params && { params })
      })
    }
  })

  return <Graph>{
    activities: ActivitiesFactory({ httpClient, config }),
    applications: ApplicationsFactory({ httpClient, config }),
    tags: TagsFactory({ httpClient, config }),
    drives: DrivesFactory({ httpClient, config }),
    driveItems: DriveItemsFactory({ httpClient, config }),
    users: UsersFactory({ httpClient, config }),
    groups: GroupsFactory({ httpClient, config }),
    permissions: PermissionsFactory({ httpClient, config })
  }
}
