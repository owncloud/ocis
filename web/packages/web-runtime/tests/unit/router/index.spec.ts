import { createRouter, createWebHashHistory, createWebHistory } from 'vue-router'

// @vitest-environment jsdom
describe('buildUrl', () => {
  it.each`
    location                                                                          | base                         | path            | expected
    ${'https://localhost:8080/index.php/apps/web/index.html#/files/list/all'}         | ${''}                        | ${'/login'}     | ${'https://localhost:8080/index.php/apps/web#/login'}
    ${'https://localhost:8080/index.php/apps/web/index.html#/files/list/all'}         | ${''}                        | ${'/login/foo'} | ${'https://localhost:8080/index.php/apps/web#/login/foo'}
    ${'https://localhost:8080/////index.php/apps/web/////index.html#/files/list/all'} | ${''}                        | ${'/login/foo'} | ${'https://localhost:8080/index.php/apps/web#/login/foo'}
    ${'https://localhost:8080/index.php/apps/web/#/login'}                            | ${''}                        | ${'/bar.html'}  | ${'https://localhost:8080/index.php/apps/web/bar.html'}
    ${'https://localhost:8080/index.php/apps/web/#/login'}                            | ${'/index.php/apps/web/foo'} | ${'/bar'}       | ${'https://localhost:8080/index.php/apps/web/foo/bar'}
    ${'https://localhost:8080/index.php/apps/web/#/login'}                            | ${'/index.php/apps/web/foo'} | ${'/bar.html'}  | ${'https://localhost:8080/index.php/apps/web/foo/bar.html'}
    ${'https://localhost:9200/#/files/list/all'}                                      | ${''}                        | ${'/login/foo'} | ${'https://localhost:9200/#/login/foo'}
    ${'https://localhost:9200/#/files/list/all'}                                      | ${''}                        | ${'/bar.html'}  | ${'https://localhost:9200/bar.html'}
    ${'https://localhost:9200/files/list/all'}                                        | ${'/'}                       | ${'/login/foo'} | ${'https://localhost:9200/login/foo'}
    ${'https://localhost:9200/files/list/all'}                                        | ${'/foo'}                    | ${'/login/foo'} | ${'https://localhost:9200/foo/login/foo'}
    ${'https://localhost:9200/files/list/all'}                                        | ${'/'}                       | ${'/bar.html'}  | ${'https://localhost:9200/bar.html'}
    ${'https://localhost:9200/files/list/all'}                                        | ${'/foo'}                    | ${'/bar.html'}  | ${'https://localhost:9200/foo/bar.html'}
  `('$path -> $expected', async ({ location, base, path, expected }) => {
    delete window.location
    window.location = new URL(location) as any

    document.querySelectorAll('base').forEach((e) => e.remove())

    if (base) {
      const baseElement = document.createElement('base')
      baseElement.href = base
      document.getElementsByTagName('head')[0].appendChild(baseElement)
    }

    const { buildUrl } = await import('@ownclouders/web-pkg/src/helpers/router/buildUrl')
    vi.resetModules()

    // hide warnings for non-existent routes
    vi.spyOn(console, 'warn').mockImplementation(() => undefined)

    const router = createRouter({
      routes: [
        {
          path: '/login',
          component: {}
        }
      ],
      history: (base && createWebHistory(base)) || createWebHashHistory()
    })

    expect(buildUrl(router, path)).toBe(expected)
  })
})
