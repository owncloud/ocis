import { defaultComponentMocks, mount, defaultPlugins } from '@ownclouders/web-test-helpers'
import { mock, mockDeep } from 'vitest-mock-extended'
import { Resource } from '@ownclouders/web-client'
import { ClientService, eventBus, useCapabilityStore, useMessages } from '@ownclouders/web-pkg'
import { OcSelect } from '@ownclouders/design-system/components'
import { storeToRefs } from 'pinia'
import { unref } from 'vue'

vi.hoisted(() => {
  vi.doMock('@ownclouders/web-pkg', async (importOriginal) => ({
    ...(await importOriginal<any>()),
    useTagColor: vi.fn(() => ({
      tagColorIndex: vi.fn((name: string) => {
        const tagColorMap: Record<string, number> = { invoice: 7, project: 3 }
        return tagColorMap[name] ?? -1
      })
    }))
  }))
})

import TagsSelect from '../../../../src/components/SideBar/Details/TagsSelect.vue'

describe('Tag Select', () => {
  it('show tags input form if loaded successfully', () => {
    const resource = mock<Resource>({ tags: [] })
    const { wrapper } = createWrapper(resource)
    expect(wrapper.find('.tags-select').exists()).toBeTruthy()
  })

  it('all available tags are selectable', async () => {
    const tags = 'a,b,c'
    const resource = mock<Resource>({ tags: [] })
    const clientService = mockDeep<ClientService>()
    clientService.graphAuthenticated.tags.listTags.mockResolvedValueOnce(tags.split(','))

    const { wrapper } = createWrapper(resource, clientService)
    await (wrapper.vm as any).loadAvailableTagsTask.last
    expect(
      (wrapper.findComponent<typeof OcSelect>('vue-select-stub').props() as any).options
    ).toEqual([{ label: 'a' }, { label: 'b' }, { label: 'c' }])
  })

  describe('save method', () => {
    it('publishes the "save"-event', async () => {
      const eventStub = vi.spyOn(eventBus, 'publish')
      const tags = ['a', 'b']
      const resource = mock<Resource>({ tags: tags })
      const { wrapper } = createWrapper(resource, mockDeep<ClientService>(), false)
      await (wrapper.vm as any).save(tags)
      expect(eventStub).toHaveBeenCalled()
    })
  })

  test.each<[string[], { label: string }[], string[]]>([
    [['a', 'b'], [{ label: 'c' }], ['c']],
    [['a', 'b'], [{ label: 'a' }, { label: 'b' }, { label: 'c' }], ['c']],
    [
      ['a', 'b'],
      [{ label: 'a' }, { label: 'b' }, { label: 'c' }, { label: 'd' }],
      ['c', 'd']
    ]
  ])(
    'resource with the initial tags %s and selected tags %s adds %s',
    async (resourceTags, selectedTags, expected) => {
      const resource = mock<Resource>({ tags: resourceTags })
      const clientService = mockDeep<ClientService>()
      const stub = clientService.graphAuthenticated.tags.assignTags.mockResolvedValue(undefined)
      const { wrapper } = createWrapper(resource, clientService, false)

      ;(wrapper.vm as any).selectedTags = selectedTags

      await (wrapper.vm as any).save(selectedTags)

      if (expected.length) {
        expect(stub).toHaveBeenCalledWith(
          expect.objectContaining({
            tags: expected
          })
        )
      } else {
        expect(stub).not.toHaveBeenCalled()
      }
    }
  )

  test.each<[string[], { label: string }[], string[]]>([
    [['a', 'b'], [{ label: 'a' }], ['b']],
    [['a', 'b'], [{ label: 'a' }, { label: 'b' }, { label: 'c' }], []],
    [['a', 'b'], [], ['a', 'b']]
  ])(
    'resource with the initial tags %s and selected tags %s removes %s',
    async (resourceTags, selectedTags, expected) => {
      const resource = mock<Resource>({ tags: resourceTags })
      const clientService = mockDeep<ClientService>()
      const stub = clientService.graphAuthenticated.tags.unassignTags.mockResolvedValue(undefined)
      const { wrapper } = createWrapper(resource, clientService, false)

      ;(wrapper.vm as any).selectedTags = selectedTags

      await (wrapper.vm as any).save(selectedTags)

      if (expected.length) {
        expect(stub).toHaveBeenCalledWith(
          expect.objectContaining({
            tags: expected
          })
        )
      } else {
        expect(stub).not.toHaveBeenCalled()
      }
    }
  )

  it('shows message on failure', async () => {
    vi.spyOn(console, 'error').mockImplementation(() => undefined)
    const clientService = mockDeep<ClientService>()
    const assignTagsStub = clientService.graphAuthenticated.tags.assignTags.mockRejectedValue(
      new Error()
    )
    const resource = mock<Resource>({ tags: ['a'] })
    const eventStub = vi.spyOn(eventBus, 'publish')
    const { wrapper } = createWrapper(resource, clientService)
    ;(wrapper.vm as any).selectedTags.push({ label: 'b' })
    await (wrapper.vm as any).save((wrapper.vm as any).selectedTags)
    expect(assignTagsStub).toHaveBeenCalled()
    expect(eventStub).not.toHaveBeenCalled()
    const { showErrorMessage } = useMessages()
    expect(showErrorMessage).toHaveBeenCalledTimes(1)
  })

  it('does not accept tags consisting of blanks only', () => {
    const { wrapper } = createWrapper(mock<Resource>({ tags: [] }))
    const option = (wrapper.vm as any).createOption(' ')
    expect(option.error).toBeDefined()
    expect(option.selectable).toBeFalsy()
  })

  it('should not accept tags longer than max tag length', () => {
    const { wrapper } = createWrapper(mock<Resource>({ tags: [] }))

    const capabilitiesStore = useCapabilityStore()
    const { graphTagsMaxTagLength } = storeToRefs(capabilitiesStore)

    const option = (wrapper.vm as any).createOption('a'.repeat(unref(graphTagsMaxTagLength) + 1))

    expect(option.error).toBeDefined()
    expect(option.selectable).toBeFalsy()
  })

  describe('Tag chip colours', () => {
    it('a selected chip receives the fill colour for its tag name', async () => {
      const tagName = 'invoice'
      const expectedColour = 'var(--oc-color-tag-7)'
      const resource = mock<Resource>({ tags: [tagName] })
      const clientService = mockDeep<ClientService>()
      clientService.graphAuthenticated.tags.listTags.mockResolvedValue([])

      const { wrapper } = createWrapper(resource, clientService, false) // Don't stub VueSelect

      // Wait for the component to render with the selected tag
      await wrapper.vm.$nextTick()

      // OcTag turns the colour index into inline styles: background-color plus the label
      // colour that goes with it
      const ocTags = wrapper.findAll('.tags-select-tag')
      expect(ocTags.length).toBeGreaterThan(0)

      // The fill colour is applied to the tag's background-color style
      const tagElement = ocTags[0]
      const styleAttribute = tagElement.attributes('style') || ''
      expect(styleAttribute).toContain(`background-color: ${expectedColour}`)
    })

    it('a dropdown option chip receives the fill colour for its label', async () => {
      const tagName = 'project'
      const expectedColour = 'var(--oc-color-tag-3)'
      const resource = mock<Resource>({ tags: [] })
      const clientService = mockDeep<ClientService>()
      clientService.graphAuthenticated.tags.listTags.mockResolvedValueOnce([tagName])

      const mocks = { ...defaultComponentMocks(), $clientService: clientService }
      mocks.$clientService.graphAuthenticated.tags.listTags.mockResolvedValue([])
      const wrapper = mount(TagsSelect, {
        global: {
          plugins: [...defaultPlugins()],
          mocks,
          provide: { ...mocks },
          stubs: { CompareSaveDialog: true } // Don't stub VueSelect so it renders the option template
        },
        props: {
          resource
        }
      })

      // Wait for available tags to load
      await (wrapper.vm as any).loadAvailableTagsTask.last
      await wrapper.vm.$nextTick()

      // vue-select only renders the #option slot while its dropdown is open, and OcSelect
      // additionally gates that on a click of the select itself (see `dropdownEnabled`).
      await wrapper.find('.oc-select').trigger('click')
      await wrapper.find('input.vs__search').trigger('focus')
      await wrapper.vm.$nextTick()

      const optionTags = wrapper
        .findAll('.tags-select-tag')
        .filter((tag) => tag.element.closest('.vs__dropdown-menu') !== null)
      expect(optionTags.length).toBe(1)
      expect(optionTags[0].attributes('style')).toContain(`background-color: ${expectedColour}`)
    })
  })

  describe('Tag chip icon removal', () => {
    it('the price-tag-3 icon is no longer rendered in the selected chip template', async () => {
      const tagName = 'invoice'
      const resource = mock<Resource>({ tags: [tagName] })
      const clientService = mockDeep<ClientService>()
      clientService.graphAuthenticated.tags.listTags.mockResolvedValue([])

      const mocks = { ...defaultComponentMocks(), $clientService: clientService }
      mocks.$clientService.graphAuthenticated.tags.listTags.mockResolvedValue([])
      const wrapper = mount(TagsSelect, {
        global: {
          plugins: [...defaultPlugins()],
          mocks,
          provide: { ...mocks },
          stubs: { CompareSaveDialog: true }
        },
        props: {
          resource
        }
      })

      // Wait for the component to render
      await wrapper.vm.$nextTick()

      // The selected chip should not contain the price-tag-3 icon
      const selectedTags = wrapper.findAll('.tags-select-tag')
      expect(selectedTags.length).toBeGreaterThan(0)

      // Get the HTML of the tag to check for price-tag-3 reference
      const tagHTML = selectedTags[0].html()

      // Verify no price-tag-3 icon is rendered
      expect(tagHTML).not.toContain('price-tag-3')
    })
  })
})

function createWrapper(
  resource: Resource,
  clientService = mockDeep<ClientService>(),
  stubVueSelect = true
) {
  const mocks = { ...defaultComponentMocks(), $clientService: clientService }
  mocks.$clientService.graphAuthenticated.tags.listTags.mockResolvedValue([])
  return {
    wrapper: mount(TagsSelect, {
      global: {
        plugins: [...defaultPlugins()],
        mocks,
        provide: { ...mocks },
        stubs: { VueSelect: stubVueSelect, CompareSaveDialog: true }
      },
      props: {
        resource
      }
    })
  }
}
