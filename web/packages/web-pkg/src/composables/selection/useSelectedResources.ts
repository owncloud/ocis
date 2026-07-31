import { computed, unref, WritableComputedRef, Ref } from 'vue'
import { Resource } from '@ownclouders/web-client'
import { useGetMatchingSpace } from '../spaces'
import { SpaceResource } from '@ownclouders/web-client'
import { useResourcesStore } from '../piniaStores'

export interface SelectedResourcesResult {
  selectedResources: Ref<Resource[]>
  selectedResourcesIds: Ref<string[]>
  isResourceInSelection(resource: Resource): boolean
  selectedResourceSpace?: Ref<SpaceResource>
}

export const useSelectedResources = (): SelectedResourcesResult => {
  const { getMatchingSpace } = useGetMatchingSpace()
  const resourcesStore = useResourcesStore()

  const selectedResources: WritableComputedRef<Resource[]> = computed({
    get(): Resource[] {
      return resourcesStore.selectedResources
    },
    set(resources) {
      resourcesStore.setSelection(resources.map(({ id }) => id))
    }
  })
  const selectedResourcesIds: WritableComputedRef<string[]> = computed({
    get(): string[] {
      return resourcesStore.selectedIds
    },
    set(selectedIds) {
      resourcesStore.setSelection(selectedIds)
    }
  })

  const isResourceInSelection = (resource: Resource): boolean => {
    return unref(selectedResourcesIds).includes(resource.id)
  }

  const selectedResourceSpace = computed(() => {
    if (unref(selectedResources).length !== 1) {
      return null
    }
    const resource = unref(selectedResources)[0]
    return getMatchingSpace(resource)
  })

  return {
    selectedResources,
    selectedResourcesIds,
    isResourceInSelection,
    selectedResourceSpace
  }
}
