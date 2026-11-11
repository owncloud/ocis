import { computed, markRaw } from 'vue'
import { CustomComponentExtension, Extension } from '@ownclouders/web-pkg'
import SearchBar from './portals/SearchBar.vue'

const searchBarExtension: CustomComponentExtension = {
  id: 'com.github.owncloud.web.search.search-bar',
  type: 'customComponent',
  extensionPointIds: ['app.runtime.header.center'],
  content: markRaw(SearchBar)
}

export const extensions = () => {
  return computed<Extension[]>(() => {
    return [searchBarExtension]
  })
}
