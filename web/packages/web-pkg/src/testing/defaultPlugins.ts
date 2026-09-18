import { Plugin } from 'vue'
import {
  defaultPlugins as baseDefaultPlugins,
  DefaultPluginsOptions as BaseDefaultPluginsOptions
} from '@ownclouders/web-test-helpers-core'
import { PiniaMockOptions } from './pinia'

// Retypes the generic `defaultPlugins` from `@ownclouders/web-test-helpers-core` so that
// `piniaOptions` here is `PiniaMockOptions` (defined in this package) rather than
// `Record<string, unknown>`. Core can't name `PiniaMockOptions` itself - see its
// `defaultPlugins.ts` for why.
export type DefaultPluginsOptions = BaseDefaultPluginsOptions<PiniaMockOptions>

export const defaultPlugins = (options: DefaultPluginsOptions = {}): Plugin[] =>
  baseDefaultPlugins<PiniaMockOptions>(options)
