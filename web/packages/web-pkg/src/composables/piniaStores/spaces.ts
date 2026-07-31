import { defineStore } from 'pinia'
import { computed, ref, unref } from 'vue'
import {
  buildShareSpaceResource,
  isMountPointSpaceResource,
  SpaceDeletedState,
  SpaceResource
} from '@ownclouders/web-client'
import { Graph } from '@ownclouders/web-client/graph'
import {
  buildSpace,
  extractStorageId,
  isPersonalSpaceResource,
  isProjectSpaceResource
} from '@ownclouders/web-client'
import type { CollaboratorShare, MountPointSpaceResource, ShareRole } from '@ownclouders/web-client'
import { useUserStore } from './user'
import { ConfigStore, useConfigStore } from './config'
import { useSharesStore } from './shares'

// sort space members with higher permissions (managers) at the top
export const sortSpaceMembers = (shares: CollaboratorShare[]) => {
  return shares.sort((a, b) => b.permissions.length - a.permissions.length)
}

export const getSpacesByType = async ({
  graphClient,
  driveType,
  configStore,
  graphRoles,
  signal
}: {
  graphClient: Graph
  driveType: string
  configStore: ConfigStore
  graphRoles: Record<string, ShareRole>
  signal?: AbortSignal
}) => {
  const mountpoints = await graphClient.drives.listMyDrives(
    graphRoles,
    {
      orderBy: 'name asc',
      filter: `driveType eq ${driveType}`
    },
    { signal }
  )
  if (!mountpoints.length) {
    return []
  }

  const enabledMountpoints = mountpoints.filter(
    (space) =>
      !isPersonalSpaceResource(space) || space.root.deleted?.state !== SpaceDeletedState.Trashed
  )

  if (driveType !== 'mountpoint' || !configStore.options.routing?.fullShareOwnerPaths) {
    return enabledMountpoints
  }

  const rootSpaceDriveAliasMapping: Record<string, string> = {}
  enabledMountpoints.forEach((space) => {
    const { rootId, driveAlias } = space.root.remoteItem
    rootSpaceDriveAliasMapping[rootId] = driveAlias
  })

  const rootSpaces = Object.entries(rootSpaceDriveAliasMapping).map(([id, driveAlias]) =>
    // FIXME: create proper buildRootSpace (or whatever function)
    buildSpace(
      {
        id: extractStorageId(id),
        name: driveAlias, // FIXME: set a proper name
        driveType: driveAlias.split('/')[0], // FIXME: can we retrieve this from api?
        driveAlias,
        path: '/',
        serverUrl: configStore.serverUrl
      },
      graphRoles
    )
  )

  return [...enabledMountpoints, ...rootSpaces]
}

export const useSpacesStore = defineStore('spaces', () => {
  const userStore = useUserStore()
  const configStore = useConfigStore()
  const sharesStore = useSharesStore()

  const spaces = ref<SpaceResource[]>([])
  const currentSpace = ref<SpaceResource>()
  const spacesInitialized = ref(false)
  const mountPointsInitialized = ref(false)
  const spacesLoading = ref(false)

  const personalSpace = computed(() => {
    return unref(spaces).find((s) => isPersonalSpaceResource(s) && s.isOwner(userStore.user))
  })

  const setSpacesInitialized = (value: boolean) => {
    spacesInitialized.value = value
  }

  const setMountPointsInitialized = (value: boolean) => {
    mountPointsInitialized.value = value
  }

  const setSpacesLoading = (value: boolean) => {
    spacesLoading.value = value
  }

  const setCurrentSpace = (space: SpaceResource) => {
    currentSpace.value = space
  }

  const getSpaceMembers = (space: SpaceResource) => {
    // only project spaces have members
    if (!isProjectSpaceResource(space)) {
      return []
    }
    const members = sharesStore.collaboratorShares.filter((c) => c.resourceId === space.id)
    return sortSpaceMembers(members)
  }

  const addSpaces = (s: SpaceResource[]) => {
    unref(spaces).push(...s)
  }

  const removeSpace = (space: SpaceResource) => {
    spaces.value = unref(spaces).filter(({ id }) => id !== space.id)
  }

  const getSpace = (id: string) => {
    return unref(spaces).find((s) => id == s.id)
  }

  const getMountPointForSpace = async ({
    graphClient,
    space,
    signal
  }: {
    graphClient: Graph
    space: SpaceResource
    signal?: AbortSignal
  }): Promise<MountPointSpaceResource> => {
    await loadMountPoints({ graphClient, signal })

    // even if the resource has been shared via multiple permissions (e.g. directly via user and a group)
    // we only care about one matching mount point since the remote item contains all permissions
    return unref(spaces).find(
      (s) => isMountPointSpaceResource(s) && s.root?.remoteItem?.id === space.id
    )
  }

  const createShareSpace = ({
    driveAliasPrefix,
    id,
    shareName
  }: {
    driveAliasPrefix: 'share' | 'ocm-share'
    id: string
    shareName: string
  }) => {
    const space = buildShareSpaceResource({
      driveAliasPrefix,
      id,
      shareName,
      serverUrl: configStore.serverUrl
    })
    addSpaces([space])
    return space
  }

  const upsertSpace = (space: SpaceResource) => {
    const existingSpace = unref(spaces).find(({ id }) => id === space.id)
    if (existingSpace) {
      Object.assign(existingSpace, space)
      return
    }
    addSpaces([space])
  }

  const updateSpaceField = <T extends SpaceResource>({
    id,
    field,
    value
  }: {
    id: T['id']
    field: keyof T
    value: T[keyof T]
  }) => {
    const space = unref(spaces).find((space) => id === space.id) as T
    if (space) {
      space[field] = value
    }
  }

  const loadSpaces = async ({
    graphClient,
    isInVault
  }: {
    graphClient: Graph
    isInVault: boolean
  }) => {
    spacesLoading.value = true
    try {
      /**
       * FIXME: this is bad for two reasons:
       * 1. fetching by specific drive type is bad because if more drive types are being added it needs additional code.
       *    as soon as the backend allows to filter by `driveType neq virtual` we want to use that here.
       * 2. fetching the mountpoint drives only on first access is kind of error prone, because mount points are
       *    trying to be accessed in multiple code locations. all of them need to check now if mountpoints need to be
       *    fetched first. but at the moment fetching mountpoints is kind of expensive, so we need to accept that for now.
       */
      const [personalSpaces, projectSpaces] = await Promise.all([
        getSpacesByType({
          graphClient,
          driveType: 'personal',
          configStore,
          graphRoles: sharesStore.graphRoles
        }),
        getSpacesByType({
          graphClient,
          driveType: 'project',
          configStore,
          graphRoles: sharesStore.graphRoles
        })
      ])

      addSpaces([...personalSpaces, ...projectSpaces])
      spacesInitialized.value = true
    } finally {
      spacesLoading.value = false
    }
  }

  const loadMountPoints = async ({
    graphClient,
    signal
  }: {
    graphClient: Graph
    signal?: AbortSignal
  }) => {
    // fetching mount points is particularly expensive, so we do that only on first access.
    if (unref(mountPointsInitialized)) {
      return
    }
    try {
      const mountPointSpaces = await getSpacesByType({
        graphClient,
        driveType: 'mountpoint',
        configStore,
        graphRoles: sharesStore.graphRoles,
        signal
      })
      addSpaces(mountPointSpaces)
    } finally {
      mountPointsInitialized.value = true
    }
  }

  const reloadProjectSpaces = async ({
    graphClient,
    isInVault,
    signal
  }: {
    graphClient: Graph
    isInVault: boolean
    signal?: AbortSignal
  }) => {
    const projectSpaces = await getSpacesByType({
      graphClient,
      driveType: 'project',
      configStore,
      graphRoles: sharesStore.graphRoles,
      signal
    })
    spaces.value = unref(spaces).filter((s) => !isProjectSpaceResource(s))
    addSpaces(projectSpaces)
  }

  const getSpacesByName = (name: string): SpaceResource[] => {
    const matchingSpaces = unref(spaces).filter((s) => s.name === name)
    return matchingSpaces
  }

  return {
    spaces,
    spacesInitialized,
    mountPointsInitialized,
    spacesLoading,
    currentSpace,
    personalSpace,

    getSpace,
    createShareSpace,
    setSpacesInitialized,
    setMountPointsInitialized,
    setSpacesLoading,
    setCurrentSpace,
    getSpaceMembers,
    getMountPointForSpace,

    addSpaces,
    removeSpace,
    upsertSpace,
    updateSpaceField,
    loadSpaces,
    loadMountPoints,
    reloadProjectSpaces,
    getSpacesByName
  }
})

export type SpacesStore = ReturnType<typeof useSpacesStore>
