/// <reference types="vite/client" />

import { Route, Router } from 'vue-router'

// This file must have at least one export or import on top-level
export {}

declare module 'vue' {
  interface ComponentCustomProperties {
    $router: Router
    $route: Route
  }
}
