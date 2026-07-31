import GroupsList from '../../../../src/components/Groups/GroupsList.vue'
import {
  defaultComponentMocks,
  defaultPlugins,
  mount,
  shallowMount
} from '@ownclouders/web-test-helpers'
import { displayPositionedDropdown, eventBus, queryItemAsString } from '@ownclouders/web-pkg'
import { SideBarEventTopics } from '@ownclouders/web-pkg'
import { useGroupSettingsStore } from '../../../../src/composables'
import { Group } from '@ownclouders/web-client/graph/generated'

const getGroupMocks = () =>
  [
    { id: '1', members: [] },
    { id: '2', members: [] }
  ] as Group[]

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  queryItemAsString: vi.fn(),
  displayPositionedDropdown: vi.fn()
}))

vi.mock('vue-router', () => ({
  useRoute: vi.fn().mockReturnValue({ query: {} }),
  useRouter: vi.fn()
}))

describe('GroupsList', () => {
  describe('method "orderBy"', () => {
    it('should return an ascending ordered list while desc is set to false', () => {
      const { wrapper } = getWrapper()

      expect(
        wrapper.vm.orderBy(
          [{ displayName: 'users' }, { displayName: 'admins' }],
          'displayName',
          false
        )
      ).toEqual([{ displayName: 'admins' }, { displayName: 'users' }])
    })
    it('should return an descending ordered list while desc is set to true', () => {
      const { wrapper } = getWrapper()

      expect(
        wrapper.vm.orderBy(
          [{ displayName: 'admins' }, { displayName: 'users' }],
          'displayName',
          true
        )
      ).toEqual([{ displayName: 'users' }, { displayName: 'admins' }])
    })
  })

  describe('method "filter"', () => {
    it('should return a list containing record admins if search term is "ad"', () => {
      const { wrapper } = getWrapper()

      expect(
        wrapper.vm.filter([{ displayName: 'users' }, { displayName: 'admins' }], 'ad')
      ).toEqual([{ displayName: 'admins' }])
    })
    it('should return an an empty list if search term does not match any entry', () => {
      const { wrapper } = getWrapper()

      expect(
        wrapper.vm.filter([{ displayName: 'admins' }, { displayName: 'users' }], 'ownClouders')
      ).toEqual([])
    })
  })
  it('should show the context menu on right click', async () => {
    const groups = getGroupMocks()
    const spyDisplayPositionedDropdown = vi.mocked(displayPositionedDropdown)
    const { wrapper } = getWrapper({ mountType: mount, groups })
    await wrapper.find(`[data-item-id="${groups[0].id}"]`).trigger('contextmenu')
    expect(spyDisplayPositionedDropdown).toHaveBeenCalledTimes(1)
  })
  it('should show the context menu on context menu button click', async () => {
    const groups = getGroupMocks()
    const spyDisplayPositionedDropdown = vi.mocked(displayPositionedDropdown)
    const { wrapper } = getWrapper({ mountType: mount, groups })
    await wrapper.find('.groups-table-btn-action-dropdown').trigger('click')
    expect(spyDisplayPositionedDropdown).toHaveBeenCalledTimes(1)
  })
  it('should show the group details on details button click', async () => {
    const groups = getGroupMocks()
    const eventBusSpy = vi.spyOn(eventBus, 'publish')
    const { wrapper } = getWrapper({ mountType: mount, groups })
    await wrapper.find('.groups-table-btn-details').trigger('click')
    expect(eventBusSpy).toHaveBeenCalledWith(SideBarEventTopics.open)
  })
  describe('toggle selection', () => {
    describe('selectGroups method', () => {
      it('selects all groups', () => {
        const groups = getGroupMocks()
        const { wrapper } = getWrapper({ mountType: shallowMount, groups })
        wrapper.vm.selectGroups(groups)
        const { setSelectedGroups } = useGroupSettingsStore()
        expect(setSelectedGroups).toHaveBeenCalledWith(groups)
      })
    })
    describe('selectGroup method', () => {
      it('selects a group', () => {
        const groups = getGroupMocks()
        const { wrapper } = getWrapper({ mountType: shallowMount, groups: [groups[0]] })
        wrapper.vm.selectGroup(groups[0])
        const { addSelectedGroup } = useGroupSettingsStore()
        expect(addSelectedGroup).toHaveBeenCalledWith(groups[0])
      })
      it('de-selects a selected group', () => {
        const groups = getGroupMocks()
        const { wrapper } = getWrapper({
          mountType: shallowMount,
          groups: [groups[0]],
          selectedGroups: [groups[0]]
        })
        wrapper.vm.selectGroup(groups[0])
        const { setSelectedGroups } = useGroupSettingsStore()
        expect(setSelectedGroups).toHaveBeenCalledWith([])
      })
    })
    describe('unselectAllGroups method', () => {
      it('de-selects all selected groups', () => {
        const groups = getGroupMocks()
        const { wrapper } = getWrapper({
          mountType: shallowMount,
          groups: [groups[0]],
          selectedGroups: [groups[0]]
        })
        wrapper.vm.unselectAllGroups()
        const { setSelectedGroups } = useGroupSettingsStore()
        expect(setSelectedGroups).toHaveBeenCalledWith([])
      })
    })
  })
})

function getWrapper({
  mountType = shallowMount,
  groups = [],
  selectedGroups = []
}: { mountType?: typeof mount; groups?: Group[]; selectedGroups?: Group[] } = {}) {
  vi.mocked(queryItemAsString).mockImplementationOnce(() => '1')
  vi.mocked(queryItemAsString).mockImplementationOnce(() => '100')
  const mocks = defaultComponentMocks()

  return {
    wrapper: mountType(GroupsList, {
      props: {
        headerPosition: 0
      },
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
          OcCheckbox: true
        }
      }
    })
  }
}
