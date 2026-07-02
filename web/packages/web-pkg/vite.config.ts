import { join, resolve } from 'path'
import { defineConfig, searchForWorkspaceRoot } from 'vite'
import dts from 'vite-plugin-dts'
import { nodePolyfills } from 'vite-plugin-node-polyfills'
import vue from '@vitejs/plugin-vue'
import pkg from './package.json' assert { type: 'json' }

const projectRootDir = searchForWorkspaceRoot(process.cwd())
const external = [...Object.keys(pkg.dependencies)]

export default defineConfig({
  resolve: {
    alias: {
      crypto: join(projectRootDir, 'polyfills/crypto.js')
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `
          @use "sass:math";
          @import "${projectRootDir}/packages/design-system/src/styles/styles";
        `,
        silenceDeprecations: ['legacy-js-api', 'import']
      }
    }
  },
  build: {
    lib: {
      entry: resolve(__dirname, 'src/index.ts'),
      name: 'web-pkg',
      fileName: 'web-pkg'
    },
    rollupOptions: {
      external
    }
  },
  plugins: [
    vue(),
    nodePolyfills({
      exclude: ['crypto']
    }),
    dts({ exclude: ['**/tests'], include: ['src'], outDir: 'dist/types', insertTypesEntry: true }),
    {
      name: '@ownclouders/vite-plugin-docs',
      transform(src, id) {
        if (id.includes('type=docs')) {
          return {
            code: 'export default {}',
            map: null
          }
        }
      }
    }
  ]
})
