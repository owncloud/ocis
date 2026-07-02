import { useRouteQuery } from '../../../../src/composables'
import { Ref, nextTick, computed, ComputedRef, unref } from 'vue'
import { getComposableWrapper, createRouter } from '@ownclouders/web-test-helpers'

describe('useRouteQuery', () => {
  it('is reactive', async () => {
    const router = createRouter({
      routes: [{ path: '/', redirect: null }]
    })
    router.push('/')
    await router.isReady()

    let fooQuery: Ref
    let fooValue: ComputedRef

    const mocks = { $router: router }
    const wrapper = getComposableWrapper(
      () => {
        fooQuery = useRouteQuery('foo', '1')
        fooValue = computed(() => parseInt(fooQuery.value as string))

        return {
          fooQuery,
          fooValue
        }
      },
      {
        mocks,
        provide: mocks,
        template: `<div><div id="fooQuery">{{ fooQuery }}</div><div id="fooValue">{{ fooValue }}</div></div>`
      }
    )

    expect(wrapper.find('#fooQuery').element.innerHTML).toBe('1')
    expect(wrapper.find('#fooValue').element.innerHTML).toBe('1')
    expect(typeof fooQuery.value).toBe('string')
    expect(fooQuery.value).toBe('1')
    expect(typeof fooValue.value).toBe('number')
    expect(fooValue.value).toBe(1)

    fooQuery.value = '2'

    // FIXME: why do we have to wait for so many ticks?
    // Why don't we have to do that for any other expectation in the whole file?!
    for (let i = 0; i < 33; i++) {
      await nextTick()
    }

    expect(wrapper.find('#fooQuery').element.innerHTML).toBe('2')
    expect(wrapper.find('#fooValue').element.innerHTML).toBe('2')

    await router.push({ path: '/', query: { foo: '3' } })
    await nextTick()
    expect(wrapper.find('#fooQuery').element.innerHTML).toBe('3')
    expect(wrapper.find('#fooValue').element.innerHTML).toBe('3')

    await router.push({ path: '/', query: {} })
    await nextTick()
    expect(wrapper.find('#fooQuery').element.innerHTML).toBe('1')
    expect(wrapper.find('#fooValue').element.innerHTML).toBe('1')
  })

  it('has a default value if route query is not set', () => {
    const router = createRouter()

    const mocks = { $router: router }
    getComposableWrapper(
      async () => {
        const fooQuery = useRouteQuery('foo', 'defaultValue')

        expect(fooQuery.value).toBe('defaultValue')

        router.push({ path: '/', query: { foo: 'foo-1' } })
        await nextTick()
        expect(fooQuery.value).toBe('foo-1')

        router.push({})
        expect(fooQuery.value).toBe('defaultValue')
      },
      { mocks, provide: mocks }
    )
  })

  it('should update on route query change', () => {
    const router = createRouter()

    const mocks = { $router: router }
    getComposableWrapper(
      async () => {
        const fooQuery = useRouteQuery('foo')
        const barQuery = useRouteQuery('bar')

        router.push({ path: '/', query: { foo: 'foo-1' } })
        await nextTick()
        expect(fooQuery.value).toBe('foo-1')
        expect(barQuery.value).toBeFalsy()

        router.push({ query: { foo: 'foo-2' } })
        expect(fooQuery.value).toBe('foo-2')
        expect(barQuery.value).toBeFalsy()

        router.push({ query: { foo: 'foo-3', bar: 'bar-1' } })
        expect(fooQuery.value).toBe('foo-3')
        expect(barQuery.value).toBe('bar-1')
      },
      { mocks, provide: mocks }
    )
  })

  it('should be undefined if route changes and query is not present', () => {
    const router = createRouter()

    const mocks = { $router: router }
    getComposableWrapper(
      async () => {
        const fooQuery = useRouteQuery('foo')

        router.push({ path: '/home', query: { foo: 'bar' } })
        await nextTick()
        expect(fooQuery.value).toBe('bar')

        router.push({ path: '/sub' })
        expect(fooQuery.value).toBeUndefined()
      },
      { mocks, provide: mocks }
    )
  })

  it('should update route query', () => {
    const router = createRouter()

    const mocks = { $router: router }
    getComposableWrapper(
      async () => {
        const fooQuery = useRouteQuery('foo')

        router.push({ path: '/home' })
        expect(fooQuery.value).toBeUndefined()

        fooQuery.value = 'changedThroughRef'
        await nextTick()
        expect(unref(router.currentRoute).query.foo).toBe('changedThroughRef')
      },
      { mocks, provide: mocks }
    )
  })
})
