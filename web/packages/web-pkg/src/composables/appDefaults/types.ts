import { MaybeRef } from '../../utils'
import { LocationQuery } from '../router'
import { SpaceResource } from '@ownclouders/web-client'
import { RouteParams } from 'vue-router'

export interface FileContext {
  path: MaybeRef<string>
  driveAliasAndItem: MaybeRef<string>
  space: MaybeRef<SpaceResource>
  item: MaybeRef<string>
  itemId: MaybeRef<string>
  fileName: MaybeRef<string>
  routeName: MaybeRef<string>
  routeParams: MaybeRef<RouteParams>
  routeQuery: MaybeRef<LocationQuery>
}
