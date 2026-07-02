import { defaultPlugins, shallowMount, RouteLocation } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { RouteLocationPathRaw, RouterLink } from 'vue-router'
import Pagination from './OcPagination.vue'

const defaultProps = {
  pages: 5,
  currentPage: 3,
  currentRoute: mock<RouteLocation>({ name: 'files' })
}

const selectors = {
  listItemPage: '.oc-pagination-list-item-page',
  listItemCurrent: '.oc-pagination-list-item-current',
  listItemPrevious: '.oc-pagination-list-item-prev',
  listItemNext: '.oc-pagination-list-item-next',
  listItemEllipsis: '.oc-pagination-list-item-ellipsis',
  listItemLink: '.oc-pagination-list-item-link'
}

describe('OcPagination', () => {
  it('displays all pages', () => {
    const wrapper = getWrapper()

    expect(wrapper.findAll(selectors.listItemPage).length).toBe(5)
    expect(wrapper.findAll(selectors.listItemCurrent).length).toBe(1)
  })

  it('displays prev and next links', () => {
    const wrapper = getWrapper()

    expect(wrapper.find(selectors.listItemPrevious).exists()).toBeTruthy()
    expect(wrapper.find(selectors.listItemNext).exists()).toBeTruthy()
  })

  it('hides prev link if the current page is the first page', () => {
    const wrapper = getWrapper({ currentPage: 1 })

    expect(wrapper.find(selectors.listItemPrevious).exists()).toBeFalsy()
    expect(wrapper.find(selectors.listItemNext).exists()).toBeTruthy()
  })

  it('hides next link if the current page is the last page', () => {
    const wrapper = getWrapper({ currentPage: 5 })

    expect(wrapper.find(selectors.listItemPrevious).exists()).toBeTruthy()
    expect(wrapper.find(selectors.listItemNext).exists()).toBeFalsy()
  })

  it('truncates pages if maxDisplayed prop is set', () => {
    const wrapper = getWrapper({
      pages: 10,
      currentPage: 5,
      maxDisplayed: 3
    })

    expect(wrapper.findAll(selectors.listItemEllipsis).length).toBe(2)
    expect(wrapper.findAll(selectors.listItemLink).length).toBe(4)
    expect(wrapper.findAll(selectors.listItemCurrent).length).toBe(1)
  })

  it('does not truncates pages if length of pages is the same as with ...', async () => {
    const wrapper = getWrapper({
      pages: 4,
      currentPage: 1,
      maxDisplayed: 3
    })

    expect(wrapper.find(selectors.listItemEllipsis).exists()).toBeFalsy()
    expect(wrapper.find(selectors.listItemPrevious).exists()).toBeFalsy()
    expect(wrapper.find(selectors.listItemNext).exists()).toBeTruthy()

    await wrapper.setProps({ currentPage: 2 })

    expect(wrapper.find(selectors.listItemEllipsis).exists()).toBeFalsy()
    expect(wrapper.find(selectors.listItemPrevious).exists()).toBeTruthy()
    expect(wrapper.find(selectors.listItemNext).exists()).toBeTruthy()

    await wrapper.setProps({ currentPage: 4 })

    expect(wrapper.find(selectors.listItemEllipsis).exists()).toBeFalsy()
    expect(wrapper.find(selectors.listItemPrevious).exists()).toBeTruthy()
    expect(wrapper.find(selectors.listItemNext).exists()).toBeFalsy()

    await wrapper.setProps({ pages: 10 })

    expect(wrapper.find(selectors.listItemEllipsis).exists()).toBeTruthy()
  })

  it("doesn't show ellipsis if maxDisplayed prop is set but no pages are removed", () => {
    const wrapper = getWrapper({ maxDisplayed: 3 })

    expect(wrapper.findAll(selectors.listItemEllipsis).length).toBe(0)
    expect(wrapper.findAll(selectors.listItemLink).length).toBe(4)
    expect(wrapper.findAll(selectors.listItemCurrent).length).toBe(1)
  })

  it('builds correct prev and next links', () => {
    const wrapper = getWrapper({ pages: 10, currentPage: 6 })

    const prevPage = wrapper
      .findComponent<typeof RouterLink>(selectors.listItemPrevious)
      .props('to')
    const nextPage = wrapper.findComponent<typeof RouterLink>(selectors.listItemNext).props('to')

    expect((prevPage as RouteLocationPathRaw).query?.page).toBe(5)
    expect((nextPage as RouteLocationPathRaw).query?.page).toBe(7)
  })
})

function getWrapper(props = {}) {
  return shallowMount(Pagination, {
    props: { ...defaultProps, ...props },
    global: { plugins: [...defaultPlugins()] }
  })
}
