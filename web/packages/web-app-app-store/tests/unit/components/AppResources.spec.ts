import { mount } from '@ownclouders/web-test-helpers'
import AppResources from '../../../src/components/AppResources.vue'
import { App, AppResource } from '../../../src/types'
import { mock } from 'vitest-mock-extended'

const resource1: AppResource = {
  url: 'https://trololo.whatever',
  icon: 'github',
  label: 'GitHub'
}
const resource2: AppResource = {
  url: 'https://wololo',
  label: 'Wololo'
}
const resource3: AppResource = {
  url: 'https://some.url',
  icon: 'file',
  label: ''
}
const resource4: AppResource = {
  label: 'Wololo',
  url: ''
}
const resources = [resource1, resource2, resource3, resource4]

const selectors = {
  link: '[data-testid="resource-link"]',
  icon: '[data-testid="resource-icon"]',
  label: '[data-testid="resource-label"]'
}

describe('AppResources.vue', () => {
  it('renders only resources with a valid URL and a label', () => {
    const { wrapper } = getWrapper()
    expect(wrapper.findAll(selectors.link).length).toBe(2)
  })
  it('renders a link per resource including an icon if present', () => {
    const { wrapper } = getWrapper()

    const link1 = wrapper.findAll(selectors.link).at(0)
    expect(link1.exists()).toBeTruthy()
    expect(link1.attributes().href).toBe(resource1.url)
    expect(link1.find(selectors.label).text()).toBe(resource1.label)
    expect(link1.find(selectors.icon).exists()).toBeTruthy()
    if (link1.find(selectors.icon).exists()) {
      expect(link1.find(selectors.icon).attributes().name).toBe(resource1.icon)
    }

    const link2 = wrapper.findAll(selectors.link).at(1)
    expect(link2.exists()).toBeTruthy()
    expect(link2.attributes().href).toBe(resource2.url)
    expect(link2.find(selectors.label).text()).toBe(resource2.label)
    expect(link2.find(selectors.icon).exists()).toBeFalsy()
  })
})

const getWrapper = () => {
  const app = { ...mock<App>({}), resources }

  return {
    wrapper: mount(AppResources, {
      props: { app }
    })
  }
}
