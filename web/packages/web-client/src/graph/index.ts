import { AxiosInstance } from 'axios'
import { Configuration } from './generated'
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

export const graph = (baseURI: string, axiosClient: AxiosInstance): Graph => {
  const url = new URL(baseURI)
  url.pathname = [...url.pathname.split('/'), 'graph'].filter(Boolean).join('/')
  const config = new Configuration({
    basePath: url.href
  })

  return <Graph>{
    activities: ActivitiesFactory({ axiosClient, config }),
    applications: ApplicationsFactory({ axiosClient, config }),
    tags: TagsFactory({ axiosClient, config }),
    drives: DrivesFactory({ axiosClient, config }),
    driveItems: DriveItemsFactory({ axiosClient, config }),
    users: UsersFactory({ axiosClient, config }),
    groups: GroupsFactory({ axiosClient, config }),
    permissions: PermissionsFactory({ axiosClient, config })
  }
}
