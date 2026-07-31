import {
  CollaboratorShare,
  LinkShare,
  Share,
  ShareRole,
  ShareTypes,
  isProjectSpaceResource
} from '@ownclouders/web-client'
import { defineStore } from 'pinia'
import { Ref, ref, unref } from 'vue'
import {
  AddLinkOptions,
  AddShareOptions,
  DeleteLinkOptions,
  DeleteShareOptions,
  UpdateLinkOptions,
  UpdateShareOptions
} from './types'
import { useResourcesStore } from '../resources'
import { useThemeStore } from '../theme'
import { Permission, UnifiedRoleDefinition } from '@ownclouders/web-client/graph/generated'

export const useSharesStore = defineStore('shares', () => {
  const resourcesStore = useResourcesStore()
  const { getRoleIcon: getThemeRoleIcon } = useThemeStore()
  const loading = ref(false)
  const collaboratorShares = ref<CollaboratorShare[]>([]) as Ref<CollaboratorShare[]>
  const linkShares = ref<LinkShare[]>([]) as Ref<LinkShare[]>
  const graphRoles = ref<Record<string, ShareRole>>({}) as Ref<Record<string, ShareRole>>
  const hasLoadingFailed = ref(false)

  const setGraphRoles = (values: UnifiedRoleDefinition[]) => {
    graphRoles.value = values.reduce<Record<string, ShareRole>>((acc, role) => {
      acc[role.id] = {
        ...role,
        icon: getThemeRoleIcon(role)
      }
      return acc
    }, {})
  }

  const upsertCollaboratorShare = (share: CollaboratorShare) => {
    const existingShare = unref(collaboratorShares).find(({ id }) => id === share.id)

    if (existingShare) {
      Object.assign(existingShare, share)
      return
    }

    unref(collaboratorShares).push(share)
  }

  const setCollaboratorShares = (values: CollaboratorShare[]) => {
    collaboratorShares.value = values
  }

  const addCollaboratorShares = (values: CollaboratorShare[]) => {
    unref(collaboratorShares).push(...values)
  }

  const removeCollaboratorShare = (share: CollaboratorShare) => {
    collaboratorShares.value = unref(collaboratorShares).filter(({ id }) => id !== share.id)
  }

  const pruneShares = () => {
    collaboratorShares.value = []
    linkShares.value = []
    loading.value = undefined
  }

  // remove loaded shares that are not within the current path
  const removeOrphanedShares = () => {
    const ancestorIds = Object.values(resourcesStore.ancestorMetaData).map(({ id }) => id)

    if (!ancestorIds.length) {
      collaboratorShares.value = []
      linkShares.value = []
      return
    }

    unref(collaboratorShares).forEach((share) => {
      if (!ancestorIds.includes(share.resourceId)) {
        removeCollaboratorShare(share)
      }
    })

    unref(linkShares).forEach((share) => {
      if (!ancestorIds.includes(share.resourceId)) {
        removeLinkShare(share)
      }
    })
  }

  const setLinkShares = (values: LinkShare[]) => {
    linkShares.value = values
  }

  const upsertLinkShare = (share: LinkShare) => {
    const existingShare = unref(linkShares).find(({ id }) => id === share.id)

    if (existingShare) {
      Object.assign(existingShare, share)
      return
    }

    unref(linkShares).push(share)
  }

  const removeLinkShare = (share: LinkShare) => {
    linkShares.value = unref(linkShares).filter(({ id }) => id !== share.id)
  }

  const setLoading = (value: boolean) => {
    loading.value = value
  }

  const updateFileShareTypes = (id: string) => {
    const computeShareTypes = (shares: Share[]) => {
      const shareTypes = new Set<number>()
      shares.forEach((share) => {
        shareTypes.add(share.shareType)
      })
      return Array.from(shareTypes)
    }

    const file = [...resourcesStore.resources, resourcesStore.currentFolder].find(
      (f) => f?.id === id
    )
    if (!file || isProjectSpaceResource(file)) {
      return
    }

    const allShares = [...unref(collaboratorShares), ...unref(linkShares)]
    resourcesStore.updateResourceField({
      id: file.id,
      field: 'shareTypes',
      value: computeShareTypes(allShares.filter((s) => !s.indirect))
    })

    const ancestorEntry = resourcesStore.getAncestorById(id)
    if (ancestorEntry) {
      resourcesStore.updateAncestorField({
        path: ancestorEntry.path,
        field: 'shareTypes',
        value: computeShareTypes(allShares.filter((s) => !s.indirect))
      })
    }
  }

  const addShare = async ({ clientService, space, resource, options }: AddShareOptions) => {
    const client = clientService.graphAuthenticated.permissions
    const share = await client.createInvite(space.id, resource.id, options, unref(graphRoles))

    addCollaboratorShares([share])
    updateFileShareTypes(resource.id)
    return share
  }

  const updateShare = async ({
    clientService,
    space,
    resource,
    collaboratorShare,
    options
  }: UpdateShareOptions) => {
    const client = clientService.graphAuthenticated.permissions

    const payload = {
      roles: options.roles,
      expirationDateTime: options.expirationDateTime
    } satisfies Permission

    const share = await client.updatePermission<CollaboratorShare>(
      space.id,
      resource.id,
      collaboratorShare.id,
      payload,
      unref(graphRoles)
    )

    upsertCollaboratorShare(share)
    return share
  }

  const deleteShare = async ({
    clientService,
    space,
    resource,
    collaboratorShare
  }: DeleteShareOptions) => {
    const client = clientService.graphAuthenticated.permissions

    await client.deletePermission(space.id, resource.id, collaboratorShare.id)

    removeCollaboratorShare(collaboratorShare)
    updateFileShareTypes(resource.id)
  }

  const addLink = async ({ clientService, space, resource, options }: AddLinkOptions) => {
    const client = clientService.graphAuthenticated.permissions
    const link = await client.createLink(space.id, resource.id, options)

    const selectedFiles = resourcesStore.selectedResources
    const fileIsSelected =
      selectedFiles.some(({ fileId }) => fileId === resource.fileId) ||
      (selectedFiles.length === 0 && resourcesStore.currentFolder.fileId === resource.fileId)

    upsertLinkShare(link)
    updateFileShareTypes(resource.id)

    if (!fileIsSelected) {
      // we might need to update the share types for the ancestor resource as well
      const ancestor = resourcesStore.ancestorMetaData[resource.path] ?? null
      if (ancestor) {
        const { shareTypes } = ancestor
        if (!shareTypes.includes(ShareTypes.link.value)) {
          resourcesStore.updateAncestorField({
            path: ancestor.path,
            field: 'shareTypes',
            value: [...shareTypes, ShareTypes.link.value]
          })
        }
      }
    }

    return link
  }

  const updateLink = async ({
    clientService,
    space,
    resource,
    linkShare,
    options
  }: UpdateLinkOptions) => {
    const client = clientService.graphAuthenticated.permissions
    let link: LinkShare

    if (Object.hasOwn(options, 'password')) {
      link = await client.setPermissionPassword(space.id, resource.id, linkShare.id, {
        password: options.password
      })

      linkShare.hasPassword = !!options.password
    } else {
      const payload = {
        link: {
          ...(options.type && { type: options.type }),
          ...(options.displayName && {
            '@libre.graph.displayName': options.displayName
          })
        },
        ...(Object.hasOwn(options, 'expirationDateTime') && {
          expirationDateTime: options.expirationDateTime
        })
      } satisfies Permission

      link = await client.updatePermission<LinkShare>(space.id, resource.id, linkShare.id, payload)
    }

    upsertLinkShare(link)
    return link
  }

  const deleteLink = async ({ clientService, space, resource, linkShare }: DeleteLinkOptions) => {
    const client = clientService.graphAuthenticated.permissions
    await client.deletePermission(space.id, resource.id, linkShare.id)

    removeLinkShare(linkShare)
    updateFileShareTypes(resource.id)
  }

  const setHasLoadingFailed = (value: boolean) => {
    hasLoadingFailed.value = value
  }

  async function fetchShareRolesDefinitions({ clientService }): Promise<UnifiedRoleDefinition[]> {
    const result = await clientService.graphAuthenticated.permissions.listRoleDefinitions()

    return result
  }

  return {
    loading,
    collaboratorShares,
    linkShares,
    graphRoles,

    setGraphRoles,
    setLoading,
    setCollaboratorShares,
    setLinkShares,
    removeOrphanedShares,

    pruneShares,
    addShare,
    updateShare,
    deleteShare,

    addLink,
    updateLink,
    deleteLink,

    hasLoadingFailed,
    setHasLoadingFailed,
    fetchShareRolesDefinitions
  }
})

export type SharesStore = ReturnType<typeof useSharesStore>
