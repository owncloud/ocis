import { join, resolve } from 'path'
import { defineConfig, searchForWorkspaceRoot } from 'vite'
import dts from 'vite-plugin-dts'
import { nodePolyfills } from 'vite-plugin-node-polyfills'
import vue from '@vitejs/plugin-vue'
import pkg from './package.json' assert { type: 'json' }

const projectRootDir = searchForWorkspaceRoot(process.cwd())
const external = [
  ...Object.keys(pkg.dependencies),
  // Subpaths aren't covered by the exact-match entries above. The testing entry has to keep this
  // one external: inlining it would ship a second copy of the design system and of
  // @vue/test-utils, so consumers would get components and wrappers that aren't identity-equal to
  // the ones from their own `@ownclouders/design-system` / `@vue/test-utils`.
  '@ownclouders/design-system/testing'
]

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
      entry: {
        index: resolve(__dirname, 'src/index.ts'),
        testing: resolve(__dirname, 'src/testing/index.ts')
      },
      name: 'web-pkg',
      fileName: (format, entryName) => {
        const base = entryName === 'index' ? 'web-pkg' : `web-pkg-${entryName}`
        return format === 'es' ? `${base}.js` : `${base}.umd.cjs`
      }
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
