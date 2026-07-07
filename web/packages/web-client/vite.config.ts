import { join, resolve } from 'path'
import { defineConfig, searchForWorkspaceRoot } from 'vite'
import dts from 'vite-plugin-dts'
import { nodePolyfills } from 'vite-plugin-node-polyfills'

const projectRootDir = searchForWorkspaceRoot(process.cwd())

export default defineConfig({
  resolve: {
    alias: {
      crypto: join(projectRootDir, 'polyfills/crypto.js')
    }
  },
  build: {
    lib: {
      entry: {
        'web-client': resolve(__dirname, 'src/index.ts'),
        'web-client/graph': resolve(__dirname, 'src/graph/index.ts'),
        'web-client/graph/generated': resolve(__dirname, 'src/graph/generated/index.ts'),
        'web-client/ocs': resolve(__dirname, 'src/ocs/index.ts'),
        'web-client/sse': resolve(__dirname, 'src/sse/index.ts'),
        'web-client/webdav': resolve(__dirname, 'src/webdav/index.ts')
      }
    }
  },
  plugins: [
    nodePolyfills({
      exclude: ['crypto']
    }),
    dts({
      include: ['src'],
      outDir: 'dist/types',
      insertTypesEntry: true
    })
  ]
})
