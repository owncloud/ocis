import Groups from '../../../src/views/Groups.vue'
import { mock, mockDeep } from 'vitest-mock-extended'
import { ClientService } from '@ownclouders/web-pkg'
import { defaultComponentMocks, defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import { Group } from '@ownclouders/web-client/graph/generated'

const selectors = { batchActionsStub: 'batch-actions-stub', searchInput: '#groups-filter' }
const getClientServiceMock = () => {
  const clientService = mockDeep<ClientService>()
  clientService.graphAuthenticated.groups.listGroups.mockResolvedValue([
    mock<Group>({ id: '1', displayName: 'users', groupTypes: [] })
  ])
  return clientService
}
vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useAppDefaults: vi.fn()
}))

const { mockRouterPush, mockRoute } = vi.hoisted(() => {
  const route = { query: {} }
  const mockRouterPush = vi.fn((newRoute) => {
    if (newRoute.query) {
      route.query = { ...route.query, ...newRoute.query }
    }
    return Promise.resolve()
  })
  return { mockRouterPush, mockRoute: route }
})

vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => mockRoute),
  useRouter: vi.fn().mockReturnValue({ push: mockRouterPush })
}))

describe('Groups view', () => {
  describe('computed method "sideBarAvailablePanels"', () => {
    describe('EditPanel', () => {
      it('should be available when one group is selected', () => {
        const { wrapper } = getWrapper()
        expect(
          (wrapper.vm as any).sideBarAvailablePanels
            .find(({ name }) => name === 'EditPanel')
            .isVisible({ items: [{ id: '1' }] })
        ).toBeTruthy()
      })
      it('should not be available when multiple groups are selected', () => {
        const { wrapper } = getWrapper()
        expect(
          (wrapper.vm as any).sideBarAvailablePanels
            .find(({ name }) => name === 'EditPanel')
            .isVisible({ items: [{ id: '1' }, { id: '2' }] })
        ).toBeFalsy()
      })
      it('should not be available when one read-only group is selected', () => {
        const { wrapper } = getWrapper()
        expect(
          (wrapper.vm as any).sideBarAvailablePanels
            .find(({ name }) => name === 'EditPanel')
            .isVisible({ items: [{ id: '1', groupTypes: ['ReadOnly'] }] })
        ).toBeFalsy()
      })
    })
    describe('DetailsPanel', () => {
      it('should contain DetailsPanel when no group is selected', () => {
        const { wrapper } = getWrapper()
        expect(
          (wrapper.vm as any).sideBarAvailablePanels
            .find(({ name }) => name === 'DetailsPanel')
            .isVisible({ items: [] })
        ).toBeTruthy()
      })
    })
  })

  describe('batch actions', () => {
    it('do not display when no group selected', async () => {
      const { wrapper } = getWrapper()
      await (wrapper.vm as any).loadResourcesTask.last
      expect(wrapper.find(selectors.batchActionsStub).exists()).toBeFalsy()
    })
    it('display when one group selected', async () => {
      const { wrapper } = getWrapper({ selectedGroups: [{ id: '1' }] })
      await (wrapper.vm as any).loadResourcesTask.last
      await wrapper.vm.$nextTick()
      expect(wrapper.find(selectors.batchActionsStub).exists()).toBeTruthy()
    })
    it('display when more than one groups selected', async () => {
      const { wrapper } = getWrapper({ selectedGroups: [{ id: '1' }, { id: '2' }] })
      await (wrapper.vm as any).loadResourcesTask.last
      await wrapper.vm.$nextTick()
      expect(wrapper.find(selectors.batchActionsStub).exists()).toBeTruthy()
    })
  })

  describe('search', () => {
    it('should search for groups when the search term changes', async () => {
      const { wrapper, mocks } = getWrapper()
      await (wrapper.vm as any).loadResourcesTask.last
      const searchInput = wrapper.find(selectors.searchInput)
      await searchInput.setValue('test')
      await searchInput.trigger('keydown.enter')
      expect(mockRouterPush).toHaveBeenCalledWith({
        query: {
          q_displayName: 'test',
          page: '1'
        }
      })
      expect(mocks.$clientService.graphAuthenticated.groups.listGroups).toHaveBeenCalledWith(
        {
          orderBy: ['displayName'],
          expand: ['members'],
          search: 'test'
        },
        { signal: expect.any(AbortSignal) }
      )
    })
  })

  describe('provide/inject', () => {
    it('provides the first selected group', async () => {
      const selectedGroup = mock<Group>({ id: '1', displayName: 'test-group', groupTypes: [] })
      const { wrapper } = getWrapper({ selectedGroups: [selectedGroup] })
      await (wrapper.vm as any).loadResourcesTask.last

      const providedGroup = (wrapper.vm as any).$.provides.group
      expect(providedGroup.value).toEqual(selectedGroup)
    })

    it('provides undefined when no group is selected', async () => {
      const { wrapper } = getWrapper({ selectedGroups: [] })
      await (wrapper.vm as any).loadResourcesTask.last

      const providedGroup = (wrapper.vm as any).$.provides.group
      expect(providedGroup.value).toBeUndefined()
    })

    it('provides the first group when multiple groups are selected', async () => {
      const firstGroup = mock<Group>({ id: '1', displayName: 'first-group', groupTypes: [] })
      const secondGroup = mock<Group>({ id: '2', displayName: 'second-group', groupTypes: [] })
      const { wrapper } = getWrapper({ selectedGroups: [firstGroup, secondGroup] })
      await (wrapper.vm as any).loadResourcesTask.last

      const providedGroup = (wrapper.vm as any).$.provides.group
      expect(providedGroup.value).toEqual(firstGroup)
    })
  })
})

function getWrapper({
  clientService = getClientServiceMock(),
  groups = [],
  selectedGroups = []
}: { clientService?: ClientService; groups?: Group[]; selectedGroups?: Group[] } = {}) {
  const mocks = { ...defaultComponentMocks(), $clientService: clientService }

  return {
    wrapper: mount(Groups, {
      global: {
        plugins: [
          ...defaultPlugins({
            piniaOptions: {
              groupSettingsStore: { groups, selectedGroups }
            }
          })
        ],
        mocks,
        provide: mocks,
        stubs: {
          AppLoadingSpinner: true,
          NoContentMessage: true,
          GroupsList: {
            template: '<div><slot name="filter"></slot></div>'
          },
          OcBreadcrumb: true,
          BatchActions: true
        }
      }
    }),
    mocks
  }
}
