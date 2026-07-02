import { RouteComponents } from './router'
import { RouteLocationNamedRaw, RouteRecordRaw } from 'vue-router'
import { createLocation, $gettext, isLocationActiveDirector } from './utils'

type commonTypes = 'files-common-favorites' | 'files-common-search'

export const createLocationCommon = (name: commonTypes, location = {}): RouteLocationNamedRaw =>
  createLocation(name, location)

export const locationFavorites = createLocationCommon('files-common-favorites')
export const locationSearch = createLocationCommon('files-common-search')

export const isLocationCommonActive = isLocationActiveDirector<commonTypes>(
  locationFavorites,
  locationSearch
)

export const buildRoutes = (components: RouteComponents): RouteRecordRaw[] => [
  {
    path: '/search',
    component: components.App,
    children: [
      {
        name: locationSearch.name,
        path: 'list/:page?',
        component: components.SearchResults,
        meta: {
          authContext: 'user',
          title: $gettext('Search results'),
          contextQueryItems: [
            'term',
            'provider',
            'q_tags',
            'q_lastModified',
            'q_titleOnly',
            'q_mediaType',
            'scope',
            'useScope'
          ]
        }
      }
    ]
  },
  {
    path: '/favorites',
    component: components.App,
    children: [
      {
        name: locationFavorites.name,
        path: '',
        component: components.Favorites,
        meta: {
          authContext: 'user',
          title: $gettext('Favorite files')
        }
      }
    ]
  }
]
