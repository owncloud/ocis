import {
  defineConfig,
  mergeConfig,
  Plugin,
  searchForWorkspaceRoot,
  UserConfig,
  ViteDevServer
} from 'vite'
import vue from '@vitejs/plugin-vue'
import EnvironmentPlugin from 'vite-plugin-environment'
import { viteStaticCopy } from 'vite-plugin-static-copy'
import { treatAsCommonjs } from 'vite-plugin-treat-umd-as-commonjs'
import visualizer from 'rollup-plugin-visualizer'
import compression from 'rollup-plugin-gzip'
import { nodePolyfills } from 'vite-plugin-node-polyfills'

import { basename, join } from 'path'
import { existsSync, readdirSync, readFileSync } from 'fs'

// build config
import packageJson from './package.json'
import { compilerOptions } from './vite.config.common'
import { getUserAgentRegex } from 'browserslist-useragent-regexp'
import browserslistToEsbuild from 'browserslist-to-esbuild'
import fetch from 'node-fetch'
import { Agent } from 'https'

// @ts-ignore
import ejs from 'ejs'

const dist = process.env.DIST_DIR || 'dist'

const buildConfig = {
  requirejs: {},
  cdn: process.env.CDN === 'true',
  documentation_url: process.env.DOCUMENTATION_URL,
  ...(process.env.REQUIRE_TIMEOUT && {
    requirejs: { waitSeconds: parseInt(process.env.REQUIRE_TIMEOUT) }
  })
}

const projectRootDir = searchForWorkspaceRoot(process.cwd())
const { version } = packageJson
const supportedBrowsersRegex = getUserAgentRegex({ allowHigherVersions: true })

const stripScssMarker = '/* STYLES STRIP IMPORTS MARKER */'

// determine inputs
const input = readdirSync('packages').reduce(
  (acc, i) => {
    if (!i.startsWith('web-app')) {
      return acc
    }
    for (const extension of ['js', 'ts']) {
      const root = join('packages', i, 'src', `index.${extension}`)
      if (existsSync(root)) {
        acc[i as keyof typeof acc] = root
        break
      }
    }
    return acc
  },
  {
    'index.html': 'index.html',
    'oidc-silent-redirect.html': 'oidc-silent-redirect.html',
    'oidc-callback.html': 'oidc-callback.html'
  }
)

const getJson = async (url: string) => {
  return (
    await fetch(url, {
      ...(url.startsWith('https:') && {
        agent: new Agent({ rejectUnauthorized: false })
      })
    })
  ).json()
}

type ConfigJsonResponseBody = {
  options: Record<string, any>
}

const getConfigJson = async (url: string) => {
  return (await getJson(url)) as ConfigJsonResponseBody
}

export const historyModePlugins = () =>
  [
    {
      name: 'base-href',
      transformIndexHtml: {
        handler() {
          return [
            {
              injectTo: 'head-prepend',
              tag: 'base',
              attrs: {
                href: '/'
              }
            }
          ]
        }
      }
    }
  ] as const

export default defineConfig(({ mode, command }) => {
  const production = mode === 'production'

  /**
     When setting `OWNCLOUD_WEB_CONFIG_URL` make sure to configure the oauth/oidc client


     # oCIS
     For oCIS instances you can use `./dev/docker/ocis.idp.config.yaml`.
     In docker setups you need to mount it to `/etc/ocis/idp.yaml`.
     E.g. with docker-compose you could add a volume to the ocis container like this:
     - /home/youruser/projects/oc-web/dev/docker/ocis.idp.config.yaml:/etc/ocis/idp.yaml

     To use the oCIS deployment examples start vite like this:
     OWNCLOUD_WEB_CONFIG_URL="https://ocis.owncloud.test/config.json" pnpm vite

     */
  const configUrl =
    process.env.OWNCLOUD_WEB_CONFIG_URL || 'https://host.docker.internal:9200/config.json'

  const config: UserConfig = {
    ...(!production && {
      server: {
        port: 9201,
        ...(process.env.VITEST !== 'true' && {
          https: {
            key: readFileSync('./dev/docker/traefik/certificates/server.key'),
            cert: readFileSync('./dev/docker/traefik/certificates/server.crt')
          },
          proxy: {
            '/themes': {
              target: 'https://host.docker.internal:9200',
              changeOrigin: true,
              secure: false // allow self-signed certs
            }
          }
        })
      }
    })
  }

  return mergeConfig(
    {
      base: '',
      publicDir: 'packages/web-container',
      build: {
        // TODO: Vue3: We currently cannot inline styles of components because @vite/plugin-vue2 does not support it
        // c.f. https://github.com/vitejs/vite-plugin-vue2/issues/18
        // That's why we need to put all styles of our monorepo apps into a monolithic css file for now
        // Once the above issue is resolved or we switch to @vitejs/plugin-vue, we can remove the `cssCodeSplit` setting here
        cssCodeSplit: false,
        rollupOptions: {
          preserveEntrySignatures: 'strict',
          input,
          output: {
            dir: dist,
            chunkFileNames: join('js', 'chunks', `[name]-[hash].mjs`),
            entryFileNames: join('js', '[name]-[hash].mjs')
          }
        },
        target: browserslistToEsbuild()
      },
      server: {
        host: 'host.docker.internal',
        strictPort: true
      },
      css: {
        preprocessorOptions: {
          scss: {
            additionalData: `
              @use "sass:math";
              @import "${projectRootDir}/packages/design-system/src/styles/styles";${stripScssMarker}
            `,
            silenceDeprecations: ['legacy-js-api', 'import']
          }
        }
      },
      resolve: {
        dedupe: ['vue3-gettext'],
        alias: {
          crypto: join(projectRootDir, 'polyfills/crypto.js')
        }
      },
      plugins: [
        nodePolyfills({
          exclude: ['crypto']
        }),

        // We need to "undefine" `define` which is set by requirejs loaded in index.html
        treatAsCommonjs(),

        // In order to avoid multiple definitions of the global styles we import via additionalData into every component
        // we also insert a marker, so we can remove the global definitions after processing.
        // The downside of this approach is that @extend does not work because it modifies the global styles, thus we emit
        // a warning if `@extend` is used in the code base.
        {
          name: '@ownclouders/vite-plugin-strip-css',
          transform(src: string, id: string) {
            if (id.endsWith('.vue') && !id.includes('node_modules') && src.includes('@extend')) {
              console.warn(
                'You are using @extend in your component. This is likely not working in your styles. Please use mixins instead.',
                id.replace(`${projectRootDir}/`, '')
              )
            }
            if (id.includes('lang.scss')) {
              const split = src.split(stripScssMarker)
              const newSrc = split[split.length - 1]

              return {
                code: newSrc,
                map: null
              }
            }
          }
        },
        EnvironmentPlugin({
          PACKAGE_VERSION: version
        }),
        vue({
          template: {
            compilerOptions
          }
        }),
        viteStaticCopy({
          targets: (() => {
            const targets = [
              ...['fonts', 'icons'].map((name) => ({
                src: `packages/design-system/src/assets/${name}/*`,
                dest: `${name}`
              })),
              {
                src: 'node_modules/requirejs/require.js',
                dest: 'js'
              }
            ]

            // in development this is handled by the proxy
            if (production) {
              targets.push({
                src: `./packages/web-runtime/themes/*`,
                dest: `themes`
              })
            }

            return targets
          })()
        }),
        {
          name: '@ownclouders/vite-plugin-runtime-config',
          configureServer(server: ViteDevServer) {
            server.middlewares.use(async (request, response, next) => {
              if (request.url === '/config.json') {
                try {
                  const configJson = await getConfigJson(configUrl)
                  response.statusCode = 200
                  response.setHeader('Content-Type', 'application/json')
                  response.end(JSON.stringify(configJson))
                } catch (e) {
                  response.statusCode = 502
                  response.setHeader('Content-Type', 'application/json')
                  response.end(JSON.stringify(e))
                }
                return
              }
              next()
            })
          }
        },
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
        },
        {
          name: 'ejs',
          transformIndexHtml: {
            order: 'pre',
            handler(html, { filename }) {
              if (basename(filename) !== 'index.html') {
                return
              }
              return ejs.render(html, {
                data: {
                  buildConfig,

                  title: process.env.TITLE || 'ownCloud',
                  compilationTimestamp: new Date().getTime(),
                  supportedBrowsersRegex: supportedBrowsersRegex
                }
              })
            }
          }
        },
        {
          name: 'import-map',
          transformIndexHtml: {
            handler(html, { bundle, filename }) {
              if (basename(filename) !== 'index.html') {
                return
              }

              // Build an import map for loading internal (as in: shipped and built within this mono repo) apps
              let moduleNames: string[]
              let buildModulePath: any
              if (bundle) {
                moduleNames = Object.keys(bundle)
                // We are in production mode here and need to provide paths relative to the module that contains the import, i.e. web-runtime-*.mjs
                // so it works when oC Web is hosted in a sub folder, e.g. when using the oC 10 integration app
                buildModulePath = (moduleName: string) => moduleName.replace('js/', './')
              } else {
                // We are in development mode here, so we can just use absolute module paths
                moduleNames = Object.keys(input)
                buildModulePath = (moduleName: string) => `/packages/${moduleName}/src/index`
              }

              const re = new RegExp(/(web-app-.*)/)
              const map = Object.fromEntries(
                moduleNames
                  .map((m) => {
                    const appName = re.exec(bundle?.[m]?.name || m)?.[1]
                    if (appName) {
                      return [appName, buildModulePath(m)]
                    }
                  })
                  .filter(Boolean)
              )
              return [
                {
                  tag: 'script',
                  children: `window.WEB_APPS_MAP = ${JSON.stringify(map)}`
                }
              ]
            }
          }
        },
        ...(command === 'serve' ? historyModePlugins() : []),
        compression(),
        process.env.REPORT !== 'true'
          ? null
          : visualizer({
              filename: join('dist', 'report.html')
            })
      ] as Plugin[]
    },
    config
  )
})
