import { computed } from 'vue'

export const useIsSearchActive = () =>
  computed(() => !!document.getElementById('files-global-search-options'))
