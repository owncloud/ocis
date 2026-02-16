import { defineConfig, splitVendorChunkPlugin } from "vite";
import react from "@vitejs/plugin-react";
import checker from "vite-plugin-checker";
import legacy from "@vitejs/plugin-legacy";


const addScriptCSPNoncePlaceholderPlugin = () => {
  return {
    name: "add-script-nonce-placeholderP-plugin",
    apply: "build",
    transformIndexHtml: {
      order: "post",
      handler(htmlData) {

        return htmlData.replaceAll(
          /<script nomodule>/gi,
          `<script nomodule nonce="__CSP_NONCE__">`
        ).replaceAll(
          /<script type="module">/gi,
          `<script type="module" nonce="__CSP_NONCE__">`
        ).replaceAll(
          /<script nomodule crossorigin id="vite-legacy-entry"/gi,
          `<script nomodule crossorigin id="vite-legacy-entry" nonce="__CSP_NONCE__"`
        );
      },
    },
  };
};

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
      legacy({
        targets: ['edge 18'],
      }),
      env.mode !== "test" &&
        checker({
          typescript: true,
          eslint: {
            lintCommand: 'eslint --max-warnings=0 src',
          },
        }),
      splitVendorChunkPlugin(),
      addScriptCSPNoncePlaceholderPlugin(),
    ],
    test: {
      globals: true,
      environment: 'jsdom',
      setupFiles: './tests/setup.js',
    },
  };
});
