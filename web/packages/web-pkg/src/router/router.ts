import { defineComponent } from 'vue'

/**
 * we need to inject the vue files into the route builders,
 * this is because we also import the provided helpers from other js|ts files
 * like mixins, rollup seems to have a problem to import files which contain vue file imports
 * into js files which then again get imported by other vue files...
 */

type Component = ReturnType<typeof defineComponent>

export interface RouteComponents {
  App: Component
  Favorites: Component
  FilesDrop: Component
  SearchResults: Component
  Shares: {
    SharedWithMe: Component
    SharedWithOthers: Component
    SharedViaLink: Component
  }
  Spaces: {
    DriveResolver: Component
    Projects: Component
  }
  Trash: {
    Overview: Component
  }
}
