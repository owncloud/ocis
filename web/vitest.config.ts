import { defineConfig } from 'vitest/config'
import path from 'path'

const root = path.resolve(__dirname, './')

export default defineConfig({
  test: {
    projects: ['packages/*/vitest.config.ts'],
    coverage: {
      provider: 'v8',
      reportsDirectory: `${root}/coverage`,
      reporter: 'lcov',
      include: ['packages/*/src/**/*.ts', 'packages/*/src/**/*.vue']
    }
  }
})
