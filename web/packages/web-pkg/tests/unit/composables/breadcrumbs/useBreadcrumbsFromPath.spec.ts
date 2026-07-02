import { mock } from 'vitest-mock-extended'
import { SpaceResource } from '@ownclouders/web-client'
import { RouteLocation, getComposableWrapper } from '@ownclouders/web-test-helpers'
import { useBreadcrumbsFromPath } from '../../../../src/composables/breadcrumbs'
import { ref } from 'vue'
import { AncestorMetaData } from '../../../../src/types'

describe('useBreadcrumbsFromPath', () => {
  describe('builds an array of breadcrumbitems', () => {
    it('from a path', () => {
      const wrapper = getWrapper()
      const { breadcrumbsFromPath } = wrapper.vm as ReturnType<typeof useBreadcrumbsFromPath>
      const breadCrumbs = breadcrumbsFromPath({
        route: { path: '/files/spaces/personal/home/test' } as RouteLocation,
        space: ref(mock<SpaceResource>()),
        resourcePath: '/test',
        ancestorMetaData: ref(mock<AncestorMetaData>())
      })
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
      const wrapper = getWrapper()
      const { breadcrumbsFromPath, concatBreadcrumbs } = wrapper.vm as ReturnType<
        typeof useBreadcrumbsFromPath
      >
      const initialBreadCrumbs = [{ text: 'Foo' }, { text: 'Bar' }]
      const breadCrumbsFromPath = breadcrumbsFromPath({
        route: { path: '/app/foo/bar?all=500' } as RouteLocation,
        space: ref(mock<SpaceResource>()),
        resourcePath: '/bar',
        ancestorMetaData: ref(mock<AncestorMetaData>())
      })
      const result = concatBreadcrumbs(...initialBreadCrumbs, ...breadCrumbsFromPath)
      expect(result[0]).toMatchObject({ text: 'Foo' })
      expect(result[1]).toMatchObject({ text: 'Bar' })
    })
  })
})

function getWrapper() {
  return getComposableWrapper(() => {
    return useBreadcrumbsFromPath()
  })
}
