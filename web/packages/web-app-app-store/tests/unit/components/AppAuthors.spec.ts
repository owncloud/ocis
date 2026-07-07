import AppAuthors from '../../../src/components/AppAuthors.vue'
import { mock } from 'vitest-mock-extended'
import { App, AppAuthor } from '../../../src/types'
import { mount } from '@ownclouders/web-test-helpers'

const author1: AppAuthor = {
  name: 'John Doe',
  url: 'https://johndoe.com'
}
const author2: AppAuthor = {
  name: 'Jane Doe'
}
const author3: AppAuthor = {
  name: 'Wololo Priest',
  url: 'wololo'
}
const author4: AppAuthor = {
  name: 'Trololo',
  url: 'trololo'
}
const authors = [author1, author2, author3, author4]

const selectors = {
  item: '.app-author-item',
  link: '[data-testid="author-link"]',
  label: '[data-testid="author-label"]'
}

describe('AppAuthors.vue', () => {
  it('renders only authors with name and valid or empty url', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.findAll(selectors.item).length).toBe(2)
  })
  it('renders authors as link when they have a url and name', () => {
    const { wrapper } = getWrapper()
    const author = wrapper.findAll(selectors.item).at(0)
    expect(author.exists()).toBeTruthy()
    const link = author.find(selectors.link)
    expect(link.exists()).toBeTruthy()
    expect(link.attributes().href).toBe(author1.url)
    expect(link.text()).toBe(author1.name)
  })
  it('renders authors as span when they only have a name', () => {
    const { wrapper } = getWrapper()
    const author = wrapper.findAll(selectors.item).at(1)
    expect(author.find(selectors.link).exists()).toBeFalsy()
    expect(author.find(selectors.label).text()).toBe(author2.name)
  })
})

const getWrapper = () => {
  const app = { ...mock<App>({}), authors }

  return {
    wrapper: mount(AppAuthors, {
      props: { app }
    })
  }
}
