import { defineProject, mergeConfig } from 'vitest/config'
import configShared from '../../tests/unit/config/vitest.config'

export default mergeConfig(
  configShared,
  defineProject({
    test: {
      name: 'web-runtime'
    }
  })
)
