import { resolve } from 'path'
import { defineConfig, searchForWorkspaceRoot } from 'vite'
import dts from 'vite-plugin-dts'
import vue from '@vitejs/plugin-vue'
import pkg from './package.json' assert { type: 'json' }

const projectRootDir = searchForWorkspaceRoot(process.cwd())

export default defineConfig({
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
  resolve: {
    alias: {
      // this is a hack to mock it for the tests... won't be a problem if you don't plan to use the ODS standalone
      webfontloader: resolve(__dirname, './../../tests/unit/stubs/webfontloader.ts')
    }
  },
  build: {
    lib: {
      entry: {
        'design-system': resolve(__dirname, 'src/index.ts'),
        'design-system/components': resolve(__dirname, 'src/components/index.ts'),
        'design-system/composables': resolve(__dirname, 'src/composables/index.ts'),
        'design-system/helpers': resolve(__dirname, 'src/helpers/index.ts')
      }
    },
    rollupOptions: {
      external: [
        ...Object.keys(pkg.dependencies).filter(
          (dep) =>
            // include vue-select because there is something off with its module type
            dep !== 'vue-select' &&
            // include webfontloader to mock it for the tests... won't be a problem if you don't plan to use the ODS standalone
            dep !== 'webfontloader'
        ),
        ...Object.keys(pkg.peerDependencies),
        '**/tests',
        '**/*.spec.ts'
      ]
    }
  },
  plugins: [
    vue(),
    dts({
      exclude: ['**/tests', '**/*.spec.ts'],
      include: ['src'],
      outDir: 'dist/types',
      insertTypesEntry: true
    }),
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
