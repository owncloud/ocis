import TrashOverview from '../../../../src/views/trash/Overview.vue'
import {
  defaultComponentMocks,
  defaultPlugins,
  defaultStubs,
  mount,
  RouteLocation
} from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { nextTick } from 'vue'
import { SpaceResource } from '@ownclouders/web-client'
import { SortDir, AppBar } from '@ownclouders/web-pkg'
import { OcTable } from '@ownclouders/design-system/components'

const spaceMocks = [
  {
    id: '1',
    name: 'admin',
    disabled: false,
    driveType: 'personal',
    getDriveAliasAndItem: () => '1',
    isOwner: () => true
  },
  {
    id: '2',
    name: 'Project space 1',
    disabled: false,
    driveType: 'project',
    getDriveAliasAndItem: () => '2',
    isOwner: () => false
  },
  {
    id: '3',
    name: 'Project space 2',
    disabled: false,
    driveType: 'project',
    getDriveAliasAndItem: () => '3',
    isOwner: () => false
  }
] as unknown as SpaceResource[]

describe('TrashOverview', () => {
  it('should render no content message if no spaces exist', async () => {
    const { wrapper } = getWrapper({ spaces: [] })
    await (wrapper.vm as any).loadResourcesTask.last
    expect(wrapper.find('no-content-message-stub').exists()).toBeTruthy()
  })
  it('should navigate to single space trash if only one space exists', () => {
    const { mocks } = getWrapper({ spaces: [spaceMocks[0]] })
    expect(mocks.$router.push).toHaveBeenCalledWith({
      name: 'files-trash-generic',
      params: { driveAliasAndItem: spaceMocks[0].getDriveAliasAndItem(undefined) },
      query: {}
    })
  })
  describe('view states', () => {
    it('shows the loading spinner during loading', () => {
      const { wrapper } = getWrapper()
      expect(wrapper.find('oc-spinner-stub').exists()).toBeTruthy()
    })
    it('should render spaces list', async () => {
      const { wrapper } = getWrapper()
      await (wrapper.vm as any).loadResourcesTask.last
      expect(wrapper.html()).toMatchSnapshot()
    })
  })
  describe('sorting', () => {
    it('sorts by property name', async () => {
      const { wrapper } = getWrapper()
      await (wrapper.vm as any).loadResourcesTask.last
      let sortedSpaces = []

      ;(wrapper.vm as any).sortBy = 'name'
      await nextTick()
      sortedSpaces = wrapper.findComponent<typeof OcTable>({ name: 'oc-table' }).props()
        .data as SpaceResource[]
      expect(sortedSpaces.map((s) => s.id)).toEqual([
        spaceMocks[0].id,
        spaceMocks[1].id,
        spaceMocks[2].id
      ])
      ;(wrapper.vm as any).sortDir = SortDir.Desc
      await nextTick()
      sortedSpaces = wrapper.findComponent<typeof OcTable>({ name: 'oc-table' }).props()
        .data as SpaceResource[]
      expect(sortedSpaces.map((s) => s.id)).toEqual([
        spaceMocks[0].id,
        spaceMocks[2].id,
        spaceMocks[1].id
      ])
    })
    it('should set the sort parameters accordingly when calling "handleSort"', () => {
      const { wrapper } = getWrapper({ spaces: [spaceMocks[0]] })
      const sortBy = 'members'
      const sortDir = SortDir.Desc
      ;(wrapper.vm as any).handleSort({ sortBy, sortDir })
      expect((wrapper.vm as any).sortBy).toEqual(sortBy)
      expect((wrapper.vm as any).sortDir).toEqual(sortDir)
    })
  })
  describe('filtering', () => {
    it('shows only filtered spaces if filter applied', async () => {
      const { wrapper } = getWrapper()
      ;(wrapper.vm as any).filterTerm = 'admin'
      await nextTick()
      expect((wrapper.vm as any).displaySpaces.length).toEqual(1)
      expect((wrapper.vm as any).displaySpaces[0].id).toEqual(spaceMocks[0].id)
    })
  })
})

function getWrapper({ spaces = spaceMocks }: { spaces?: SpaceResource[] } = {}) {
  const mocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({ name: 'trash-overview' })
    })
  }

  return {
    mocks,
    wrapper: mount(TrashOverview, {
      global: {
        components: {
          AppBar
        },
        stubs: { ...defaultStubs, NoContentMessage: true },
        mocks,
        provide: mocks,
        plugins: [...defaultPlugins({ piniaOptions: { spacesState: { spaces } } })]
      }
    })
  }
}
