import { computed } from 'vue'
import { CustomComponentExtension, Extension } from '@ownclouders/web-pkg'
import SearchBar from './portals/SearchBar.vue'

const searchBarExtension: CustomComponentExtension = {
  id: 'com.github.owncloud.web.search.search-bar',
  type: 'customComponent',
  extensionPointIds: ['app.runtime.header.center'],
  content: SearchBar
}

export const extensions = () => {
  return computed<Extension[]>(() => {
    return [searchBarExtension]
  })
}
