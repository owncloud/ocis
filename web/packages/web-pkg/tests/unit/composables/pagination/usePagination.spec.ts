import { ref, unref } from 'vue'
import { usePagination } from '../../../../src/composables'
import { getComposableWrapper } from '@ownclouders/web-test-helpers'

describe('usePagination', () => {
  describe('computed items', () => {
    const items = [1, 2, 3, 4, 5, 6]

    it.each([
      { currentPage: 1, itemsPerPage: 100, expected: [1, 2, 3, 4, 5, 6] },
      { currentPage: 1, itemsPerPage: 2, expected: [1, 2] },
      { currentPage: 2, itemsPerPage: 2, expected: [3, 4] }
    ])('returns proper paginated items', ({ currentPage, itemsPerPage, expected }) => {
      getWrapper({
        setup: ({ items }) => {
          expect(unref(items)).toEqual(expected)
        },
        items,
        currentPage,
        itemsPerPage
      })
    })
  })
  describe('computed total', () => {
    it.each([
      { itemCount: 1, itemsPerPage: 100, expected: 1 },
      { itemCount: 101, itemsPerPage: 100, expected: 2 },
      { itemCount: 201, itemsPerPage: 100, expected: 3 }
    ])('returns proper total pages', ({ itemCount, itemsPerPage, expected }) => {
      const items = Array(itemCount).fill(1)
      getWrapper({
        setup: ({ total }) => {
          expect(unref(total)).toEqual(expected)
        },
        items,
        currentPage: 1,
        itemsPerPage
      })
    })
  })
})

function getWrapper({
  setup,
  items,
  currentPage,
  itemsPerPage
}: {
  setup: (instance: ReturnType<typeof usePagination>) => void
  items: number[]
  currentPage: number
  itemsPerPage: number
}) {
  return {
    wrapper: getComposableWrapper(() => {
      const instance = usePagination({
        items: ref(items),
        page: currentPage,
        perPage: itemsPerPage,
        perPageStoragePrefix: 'unit-tests'
      })
      setup(instance)
    })
  }
}
