import { mount, VueWrapper } from '@vue/test-utils'
import { defineComponent, nextTick } from 'vue'
import { createRouter as _createRouter, createMemoryHistory, RouterOptions } from 'vue-router'
import { defaultPlugins, DefaultPluginsOptions } from './defaultPlugins'
export { mount, shallowMount, flushPromises, VueWrapper } from '@vue/test-utils'

vi.spyOn(console, 'warn').mockImplementation(() => undefined)

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

export type { RouteLocation } from 'vue-router'
export { RouterLinkStub } from '@vue/test-utils'
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

type DefinedComponent = new (...args: any[]) => any
export type ComponentProps<T extends DefinedComponent> = InstanceType<T>['$props']
export type PartialComponentProps<T extends DefinedComponent> = Partial<ComponentProps<T>>
