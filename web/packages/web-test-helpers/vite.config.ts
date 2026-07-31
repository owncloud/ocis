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
      fileName: (format) => `web-test-helpers.${format}.js`
    },
    rollupOptions: {
      external: [...Object.keys(pkg.dependencies), ...Object.keys(pkg.peerDependencies)]
    }
  },
  plugins: [vue(), dts({ include: ['src'], outDir: 'dist/types', insertTypesEntry: true })]
})
