import { RouteLocation } from '@ownclouders/web-test-helpers'
import { breadcrumbsFromPath, concatBreadcrumbs } from '../../../src/helpers/breadcrumbs'

describe('builds an array of breadcrumbitems', () => {
  it('from a path', () => {
    const breadCrumbs = breadcrumbsFromPath(
      { path: '/files/spaces/personal/home/test' } as RouteLocation,
      '/test'
    )
    expect(breadCrumbs).toEqual([
      {
        id: expect.anything(),
        isStaticNav: false,
        allowContextActions: true,
        text: 'test',
        to: { path: '/files/spaces/personal/home/test', query: {} }
      }
    ])
  })

  it('from an array of breadcrumbitems', () => {
    const initialBreadCrumbs = [{ text: 'Foo' }, { text: 'Bar' }]
    const breadCrumbsFromPath = breadcrumbsFromPath(
      { path: '/app/foo/bar?all=500' } as RouteLocation,
      '/bar'
    )
    const result = concatBreadcrumbs(...initialBreadCrumbs, ...breadCrumbsFromPath)
    expect(result[0]).toMatchObject({ text: 'Foo' })
    expect(result[1]).toMatchObject({ text: 'Bar' })
  })
})
