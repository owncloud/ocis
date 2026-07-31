<template>
  <InnerSideBar
    v-if="isOpen"
    ref="sidebar"
    class="files-side-bar"
    :is-open="isOpen"
    :active-panel="activePanel"
    :available-panels="availablePanels"
    :panel-context="panelContext"
    :loading="isLoading"
    v-bind="$attrs"
    data-custom-key-bindings-disabled="true"
    @select-panel="setActiveSideBarPanel"
    @close="closeSideBar"
  >
    <template #rootHeader>
      <file-info
        v-if="isFileHeaderVisible"
        class="sidebar-panel__file_info"
        :is-sub-panel-active="false"
      />
      <space-info v-else-if="isSpaceHeaderVisible" class="sidebar-panel__space_info" />
    </template>
    <template #subHeader>
      <file-info
        v-if="isFileHeaderVisible"
        class="sidebar-panel__file_info"
        :is-sub-panel-active="true"
      />
      <space-info v-else-if="isSpaceHeaderVisible" class="sidebar-panel__space_info" />
    </template>
  </InnerSideBar>
</template>

<script lang="ts" setup>
import { computed, provide, readonly, ref, unref, watch } from 'vue'
import PQueue from 'p-queue'
import { SideBarPanelContext } from '../SideBar/types'
import InnerSideBar from '../SideBar/SideBar.vue'
import SpaceInfo from './Spaces/SpaceInfo.vue'
import FileInfo from './Files/FileInfo.vue'
import {
  isLocationCommonActive,
  isLocationSharesActive,
  isLocationSpacesActive
} from '../../router'
import {
  SidebarPanelExtension,
  SideBarEventTopics,
  useClientService,
  useEventBus,
  useRouter,
  useActiveLocation,
  useExtensionRegistry,
  useSelectedResources,
  useSpacesStore,
  useSharesStore,
  useResourcesStore,
  useConfigStore,
  useAppsStore,
  useCanListShares,
  useCanListVersions
} from '../../composables'
import {
  isProjectSpaceResource,
  SpaceResource,
  Resource,
  ShareRole,
  call,
  isCollaboratorShare,
  isLinkShare,
  isShareSpaceResource,
  isIncomingShareResource,
  isPersonalSpaceResource
} from '@ownclouders/web-client'
import { storeToRefs } from 'pinia'
import { useTask } from 'vue-concurrency'
import { ListPermissionsSpaceRootSelectEnum } from '@ownclouders/web-client/graph/generated'

interface Props {
  isOpen: boolean
  activePanel?: string
  space?: SpaceResource
}
const { isOpen, activePanel = null, space = null } = defineProps<Props>()
const router = useRouter()
const clientService = useClientService()
const extensionRegistry = useExtensionRegistry()
const eventBus = useEventBus()
const spacesStore = useSpacesStore()
const sharesStore = useSharesStore()
const configStore = useConfigStore()
const appsStore = useAppsStore()
const { canListShares } = useCanListShares()
const { canListVersions } = useCanListVersions()

const resourcesStore = useResourcesStore()
const { currentFolder } = storeToRefs(resourcesStore)

const loadedResource = ref<Resource>()
const versions = ref<Resource[]>([])

const availableInternalShareRoles = ref<ShareRole[]>([])
const availableExternalShareRoles = ref<ShareRole[]>([])

const { selectedResources } = useSelectedResources()

const isMetaDataLoading = ref(false)

const isLoading = computed(() => {
  return unref(isMetaDataLoading) || loadVersionsTask.isRunning
})

const panelContext = computed<SideBarPanelContext<SpaceResource, Resource, Resource>>(() => {
  if (unref(selectedResources).length === 0) {
    return {
      root: space,
      parent: null,
      items: unref(currentFolder)?.id ? [unref(currentFolder)] : []
    }
  }
  return {
    root: space,
    parent: unref(currentFolder),
    items: unref(selectedResources)
  }
})

const isSharedWithMeLocation = useActiveLocation(isLocationSharesActive, 'files-shares-with-me')
const isSharedWithOthersLocation = useActiveLocation(
  isLocationSharesActive,
  'files-shares-with-others'
)
const isSharedViaLinkLocation = useActiveLocation(isLocationSharesActive, 'files-shares-via-link')
const isProjectsLocation = isLocationSpacesActive(router, 'files-spaces-projects')
const isFavoritesLocation = useActiveLocation(isLocationCommonActive, 'files-common-favorites')
const isSearchLocation = useActiveLocation(isLocationCommonActive, 'files-common-search')

const closeSideBar = () => {
  eventBus.publish(SideBarEventTopics.close)
}
const setActiveSideBarPanel = (panelName: string) => {
  eventBus.publish(SideBarEventTopics.setActivePanel, panelName)
}
const isFileHeaderVisible = computed(() => {
  return (
    unref(panelContext).items?.length === 1 && !isProjectSpaceResource(unref(panelContext).items[0])
  )
})
const isSpaceHeaderVisible = computed(() => {
  return (
    unref(panelContext).items?.length === 1 && isProjectSpaceResource(unref(panelContext).items[0])
  )
})

const isShareLocation = computed(() => {
  return (
    unref(isSharedWithMeLocation) ||
    unref(isSharedWithOthersLocation) ||
    unref(isSharedViaLinkLocation)
  )
})
const isFlatFileList = computed(() => {
  return unref(isShareLocation) || unref(isSearchLocation) || unref(isFavoritesLocation)
})

const availablePanels = computed(() =>
  extensionRegistry
    .requestExtensions<SidebarPanelExtension<SpaceResource, Resource, Resource>>({
      id: 'global.files.sidebar',
      extensionType: 'sidebarPanel'
    })
    .map((e) => e.panel)
)

const loadVersionsTask = useTask(function* (signal, resource: Resource) {
  versions.value = yield clientService.webdav.listFileVersions(resource.id, { signal })
})

const loadSharesTask = useTask(function* (signal, resource: Resource) {
  try {
    sharesStore.setLoading(true)
    sharesStore.removeOrphanedShares()
    sharesStore.setHasLoadingFailed(false)

    const { collaboratorShares: collaboratorCache, linkShares: linkCache } = sharesStore
    const client = clientService.graphAuthenticated.permissions

    let driveId = space?.id
    if (isShareSpaceResource(space)) {
      const matchingMountPoint = yield spacesStore.getMountPointForSpace({
        graphClient: clientService.graphAuthenticated,
        space,
        signal
      })
      if (matchingMountPoint) {
        driveId = matchingMountPoint.root.remoteItem.rootId
      }
    }

    // load direct shares
    const { shares, allowedRoles } = yield* call(
      client.listPermissions(driveId, resource.fileId, sharesStore.graphRoles, {}, { signal })
    )

    const loadedCollaboratorShares = shares.filter(isCollaboratorShare)
    const loadedLinkShares = shares.filter(isLinkShare)

    const rolesArray = Object.values(sharesStore.graphRoles)
    availableInternalShareRoles.value =
      allowedRoles?.map((r) => {
        return {
          ...r,
          icon: rolesArray.find((role) => role.id === r.id)?.icon
        }
      }) || []

    // load external share roles
    if (appsStore.isAppEnabled('open-cloud-mesh')) {
      const { allowedRoles } = yield* call(
        client.listPermissions(
          driveId,
          resource.fileId,
          sharesStore.graphRoles,
          {
            filter: `@libre.graph.permissions.roles.allowedValues/rolePermissions/any(p:contains(p/condition, '@Subject.UserType=="Federated"'))`,
            select: [ListPermissionsSpaceRootSelectEnum.LibreGraphPermissionsRolesAllowedValues]
          },
          { signal }
        )
      )

      availableExternalShareRoles.value =
        allowedRoles?.map((r) => {
          return {
            ...r,
            icon: rolesArray.find((role) => role.id === r.id)?.icon
          }
        }) || []
    }

    // use cache for indirect shares
    const useCache = !unref(isFlatFileList) && !unref(isProjectsLocation)
    if (useCache) {
      collaboratorCache.forEach((share) => {
        if (loadedCollaboratorShares.some((s) => s.id === share.id)) {
          return
        }

        loadedCollaboratorShares.push({ ...share, indirect: true })
      })

      linkCache.forEach((share) => {
        if (loadedLinkShares.some((s) => s.id === share.id)) {
          return
        }

        loadedLinkShares.push({ ...share, indirect: true })
      })
    }

    if (isLocationCommonActive(router, 'files-common-search')) {
      yield resourcesStore.loadAncestorMetaData({
        folder: unref(resource),
        space: unref(space),
        client: clientService.webdav,
        signal
      })
    }

    // gather all ancestors we need to load shares for (indirect shares, space members)
    const cachedIds = [...collaboratorCache, ...linkCache].map(({ resourceId }) => resourceId)
    const ancestorIds = Object.values(resourcesStore.ancestorMetaData)
      .filter(({ id, path }) => {
        if (id === resource.id || cachedIds.includes(id)) {
          // share already cached
          return false
        }
        if (isIncomingShareResource(resource)) {
          // incoming shares don't have ancestors because they are root elements themselves
          return false
        }
        if (isPersonalSpaceResource(space)) {
          // filter out personal space roots since they don't have shares
          return path !== '/'
        }
        return true
      })
      .map(({ id }) => id)

    if (
      unref(isFlatFileList) &&
      isProjectSpaceResource(space) &&
      !isProjectSpaceResource(resource)
    ) {
      // add project space to ancestors in flat file list where we don't have ancestors
      // to display space members in the sidebar
      ancestorIds.push(space.id)
    }

    const queue = new PQueue({
      concurrency: configStore.options.concurrentRequests.shares.list
    })

    const promises = [...new Set(ancestorIds)].map((id) => {
      return queue.add(() =>
        clientService.graphAuthenticated.permissions
          .listPermissions(driveId, id, sharesStore.graphRoles, {}, { signal })
          .then((result) => {
            const indirectShares = result.shares.map((s) => ({ ...s, indirect: true }))
            loadedCollaboratorShares.push(...indirectShares.filter(isCollaboratorShare))
            loadedLinkShares.push(...indirectShares.filter(isLinkShare))
          })
      )
    })

    yield Promise.allSettled(promises)
    sharesStore.setCollaboratorShares(loadedCollaboratorShares)
    sharesStore.setLinkShares(loadedLinkShares)
  } catch (error) {
    console.error(error)
    sharesStore.setHasLoadingFailed(true)
  } finally {
    sharesStore.setLoading(false)
  }
}).restartable()

watch(
  () => [...unref(panelContext).items, isOpen],
  async () => {
    if (unref(panelContext).items?.length !== 1) {
      return
    }

    if (!isOpen) {
      versions.value = []
      return
    }

    const resource = unref(panelContext).items[0]

    if (loadVersionsTask.isRunning) {
      loadVersionsTask.cancelAll()
    }

    if (!canListVersions({ space, resource })) {
      return
    }

    try {
      await loadVersionsTask.perform(resource)
    } catch (e) {
      console.error(e)
    }
  },
  { immediate: true, deep: true }
)

watch(
  () => [...unref(panelContext).items, isOpen],
  async () => {
    if (!isOpen) {
      sharesStore.pruneShares()
      loadedResource.value = null
      return
    }
    if (unref(panelContext).items?.length !== 1) {
      // don't load additional metadata for empty or multi-select contexts
      return
    }
    const resource = unref(panelContext).items[0]
    isMetaDataLoading.value = true
    if (canListShares({ space, resource })) {
      try {
        if (loadSharesTask.isRunning) {
          loadSharesTask.cancelAll()
        }

        loadSharesTask.perform(resource)
      } catch (e) {
        console.error(e)
      }
    }

    if (!unref(isShareLocation)) {
      loadedResource.value = resource
      isMetaDataLoading.value = false
      return
    }

    // shared resources look different, hence we need to fetch the actual resource here
    try {
      const webDavResource = await clientService.webdav.getFileInfo(space, {
        path: resource.path
      })

      // make sure props from the share (=resource) are available on the merged resource
      const mergedResource = {
        ...webDavResource,
        ...resource,
        tags: webDavResource.tags // tags are always [] in Graph API, hence take them from webdav
      }
      loadedResource.value = mergedResource
    } catch (error) {
      loadedResource.value = resource
      console.error(error)
    }
    isMetaDataLoading.value = false
  },
  {
    deep: true,
    immediate: true
  }
)

provide('resource', readonly(loadedResource))
provide('versions', readonly(versions))
provide(
  'space',
  computed(() => space)
)
provide(
  'activePanel',
  computed(() => activePanel)
)
provide('availableInternalShareRoles', readonly(availableInternalShareRoles))
provide('availableExternalShareRoles', readonly(availableExternalShareRoles))
</script>

<style lang="scss">
.files-side-bar {
  z-index: 3;

  .sidebar-panel {
    &__file_info {
      padding: var(--oc-space-small) var(--oc-space-small) 0 var(--oc-space-small);
    }
  }

  ._clipboard-success-animation {
    animation-name: _clipboard-success-animation;
    animation-duration: 0.8s;
    animation-timing-function: ease-out;
    animation-fill-mode: both;
  }
}

@keyframes _clipboard-success-animation {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.9;
  }
  100% {
    opacity: 0;
  }
}
</style>
