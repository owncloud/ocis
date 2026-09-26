import { defineComponent, nextTick } from 'vue'
import { createRouter as _createRouter, createMemoryHistory, RouterOptions } from 'vue-router'
import {
  mount,
  shallowMount,
  flushPromises,
  VueWrapper,
  RouterLinkStub,
  routerLinkStubPlugin
} from '@ownclouders/design-system/testing'
import type {
  RouteLocation,
  ComponentProps,
  PartialComponentProps
} from '@ownclouders/design-system/testing'
import { defaultPlugins, DefaultPluginsOptions } from './defaultPlugins'

// Generic mount helpers and prop types live in `@ownclouders/design-system/testing`, since that
// package sits below this one and cannot import from it. Re-exported here so consumers get the
// full set from a single entrypoint.
export { mount, shallowMount, flushPromises, VueWrapper, RouterLinkStub, routerLinkStubPlugin }
export type { RouteLocation, ComponentProps, PartialComponentProps }

export const getComposableWrapper = <T>(
  setup: (...args: any[]) => T,
  {
    mocks = undefined,
    provide = undefined,
    template = undefined,
    pluginOptions = undefined
  }: {
    mocks?: Record<string, unknown>
    provide?: Record<string, unknown>
    template?: string
    pluginOptions?: DefaultPluginsOptions
  } = {}
) => {
  return mount(
    defineComponent({
      setup,
      template: template ? template : '<div></div>'
    }),
    {
      global: {
        plugins: [...defaultPlugins(pluginOptions)],
        ...(mocks && { mocks }),
        ...(provide && { provide })
      }
    }
  )
}

export const getOcSelectOptions = async (
  wrapper: VueWrapper<unknown>,
  selector: string,
  options = { close: true }
) => {
  const selectElement = wrapper.find(selector)
  await selectElement.find('input').trigger('click')
  await selectElement.find('.vs__dropdown-toggle').trigger('mousedown')

  const optionElements = selectElement.findAll<HTMLOptionElement>('.vs__dropdown-option')

  if (options.close) {
    await selectElement.find('.vs__search').trigger('blur')
  }

  return optionElements
}

export const createRouter = (options?: Partial<RouterOptions>) =>
  _createRouter({
    history: createMemoryHistory(),
    routes: [],
    strict: false,
    ...options
  })

export const writable = <T>(value: Readonly<T>): T => {
  return value as T
}

export const sleep = (ms: number) => {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

export const nextTicks = async (amount: number) => {
  for (let i = 0; i < amount - 1; i++) {
    await nextTick()
  }
}
