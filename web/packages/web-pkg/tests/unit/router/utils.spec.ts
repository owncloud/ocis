import {
  isLocationActive,
  isLocationActiveDirector,
  createLocation
} from '../../../src/router/utils'
import { RouteLocation, RouteLocationNamedRaw, Router } from 'vue-router'
import { mock } from 'vitest-mock-extended'
import { ref } from 'vue'

describe('utils', () => {
  describe('isLocationActive', () => {
    it('returns true if one location is active', () => {
      const fakeRouter = mock<Router>({
        currentRoute: ref({ name: 'foo' }),
        resolve: (r: RouteLocationNamedRaw) =>
          mock<RouteLocation & { href: string }>({ href: r.name.toString() })
      })

      expect(isLocationActive(fakeRouter, mock<RouteLocationNamedRaw>({ name: 'foo' }))).toBe(true)
      expect(
        isLocationActive(
          fakeRouter,
          mock<RouteLocationNamedRaw>({ name: 'foo' }),
          mock<RouteLocationNamedRaw>({ name: 'bar' })
        )
      ).toBe(true)
    })

    it('returns false if all locations inactive', () => {
      const fakeRouter = mock<Router>({
        currentRoute: ref({ name: 'foo' }),
        resolve: (r: RouteLocationNamedRaw) =>
          mock<RouteLocation & { href: string }>({ href: r.name.toString() })
      })

      expect(isLocationActive(fakeRouter, mock<RouteLocationNamedRaw>({ name: 'bar' }))).toBe(false)
      expect(
        isLocationActive(
          fakeRouter,
          mock<RouteLocationNamedRaw>({ name: 'bar' }),
          mock<RouteLocationNamedRaw>({ name: 'baz' })
        )
      ).toBe(false)
    })
  })

  describe('isLocationActiveDirector', () => {
    test('director can be created and be used to check active locations', () => {
      const fakeRouter = mock<Router>({
        currentRoute: ref({ name: 'unknown' }),
        resolve: (r: RouteLocationNamedRaw) =>
          mock<RouteLocation & { href: string }>({ href: r.name.toString() })
      })

      const isFilesLocationActive = isLocationActiveDirector(
        mock<RouteLocationNamedRaw>({ name: 'foo' }),
        mock<RouteLocationNamedRaw>({ name: 'bar' }),
        mock<RouteLocationNamedRaw>({ name: 'baz' })
      )
      expect(isFilesLocationActive(fakeRouter)).toBe(false)

      fakeRouter.currentRoute.value.name = 'bar'

      expect(isFilesLocationActive(fakeRouter)).toBe(true)
      expect(isFilesLocationActive(fakeRouter, 'foo', 'bar')).toBe(true)
    })

    test('director closure only allows to check known locations and throws if unknown', () => {
      const fakeRouter = mock<Router>({
        currentRoute: ref({ name: 'baz' }),
        resolve: (r: RouteLocationNamedRaw) =>
          mock<RouteLocation & { href: string }>({ href: r.name.toString() })
      })

      const isFilesLocationActive = isLocationActiveDirector(
        mock<RouteLocationNamedRaw>({ name: 'foo' }),
        mock<RouteLocationNamedRaw>({ name: 'bar' })
      )
      expect(() => isFilesLocationActive(fakeRouter, 'unknown')).toThrow()
    })
  })

  describe('createLocationDirector', () => {
    test('creates a location and handle arguments', () => {
      const testLocation = createLocation(
        'foo',
        mock<RouteLocationNamedRaw>({
          params: { foo: 'foo-param-value' },
          query: { bar: 'bar-query-value' }
        })
      )
      expect(testLocation.name).toBe('foo')
      expect(testLocation.params.foo).toBe('foo-param-value')
      expect(testLocation.query.bar).toBe('bar-query-value')
    })
  })
})
