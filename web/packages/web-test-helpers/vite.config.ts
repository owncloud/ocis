import { resolve } from 'path'
import { defineConfig } from 'vite'
import dts from 'vite-plugin-dts'
import pkg from './package.json' assert { type: 'json' }
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  build: {
    lib: {
      entry: resolve(__dirname, 'src/index.ts'),
      name: 'web-test-helpers',
      // `.cjs` for the umd build: the package is ESM (`"type": "module"`), so a `.js` file
      // can't be `require`d - which is what publishConfig.exports points the `require`
      // condition at.
      fileName: (format) => `web-test-helpers.${format}.${format === 'umd' ? 'cjs' : 'js'}`
    },
    rollupOptions: {
      // `id.startsWith` (not exact match) so subpath imports like '@ownclouders/web-pkg/src/testing'
      // are externalized too, not inlined into this package's own bundle.
      external: (id) =>
        Object.keys(pkg.peerDependencies).some((dep) => id === dep || id.startsWith(`${dep}/`))
    }
  },
  plugins: [vue(), dts({ include: ['src'], outDir: 'dist/types', insertTypesEntry: true })]
})
