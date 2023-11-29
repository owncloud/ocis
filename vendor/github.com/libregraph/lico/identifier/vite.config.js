import { defineConfig, splitVendorChunkPlugin } from 'vite';
import react from '@vitejs/plugin-react';
import checker from 'vite-plugin-checker';

export default defineConfig((env) => {
  return {
    build: {
      outDir: 'build',
      assetsDir: 'static/assets',
      manifest: 'asset-manifest.json',
      sourcemap: true,
    },
    base: './',
    server: {
      port: 3001,
      strictPort: true,
      host: '127.0.0.1',
      hmr: {
        protocol: 'ws',
        host: '127.0.0.1',
        clientPort: 3001,
      },
    },
    plugins: [
      react(),
      env.mode !== 'test' && checker({
        typescript: true,
        eslint: {
          lintCommand: 'eslint --max-warnings=0 src',
        },
      }),
      splitVendorChunkPlugin(),
    ],
    test: {
      globals: true,
      environment: 'jsdom',
      setupFiles: './tests/setup.js',
    },
  };
});
