import { computed } from 'vue'
import { useSpacesStore } from '../piniaStores'

export const useSpacesLoading = () => {
  const spacesStore = useSpacesStore()
  const areSpacesLoading = computed(
    () => !spacesStore.spacesInitialized || spacesStore.spacesLoading
  )
  return {
    areSpacesLoading
  }
}
