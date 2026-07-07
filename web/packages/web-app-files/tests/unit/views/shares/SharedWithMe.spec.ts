import SharedWithMe from '../../../../src/views/shares/SharedWithMe.vue'
import { useResourcesViewDefaults } from '../../../../src/composables'
import {
  queryItemAsString,
  InlineFilterOption,
  useSort,
  useOpenWithDefaultApp,
  ItemFilter,
  AppBar
} from '@ownclouders/web-pkg'
import { useResourcesViewDefaultsMock } from '../../../../tests/mocks/useResourcesViewDefaultsMock'
import { ref } from 'vue'
import { defaultStubs, RouteLocation } from '@ownclouders/web-test-helpers'
import { useSortMock } from '../../../../tests/mocks/useSortMock'
import { mock } from 'vitest-mock-extended'
import { defaultPlugins, mount, defaultComponentMocks } from '@ownclouders/web-test-helpers'
import { ShareTypes, IncomingShareResource, ShareType } from '@ownclouders/web-client'
import SharedWithMeSection from '../../../../src/components/Shares/SharedWithMeSection.vue'

vi.mock('../../../../src/composables/resourcesViewDefaults')
vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useSort: vi.fn().mockImplementation(() => useSortMock()),
  queryItemAsString: vi.fn(),
  useRouteQuery: vi.fn(),
  useOpenWithDefaultApp: vi.fn()
}))

describe('SharedWithMe view', () => {
  it('appBar always present', () => {
    const { wrapper } = getMountedWrapper()
    expect(wrapper.find('app-bar-stub').exists()).toBeTruthy()
  })
  it('sideBar always present', () => {
    const { wrapper } = getMountedWrapper()
    expect(wrapper.find('file-side-bar-stub').exists()).toBeTruthy()
  })
  describe('different files view states', () => {
    it('shows the loading spinner during loading', () => {
      const { wrapper } = getMountedWrapper({ loading: true })
      expect(wrapper.find('oc-spinner-stub').exists()).toBeTruthy()
    })
    it('does not show the loading spinner after loading finished', () => {
      const { wrapper } = getMountedWrapper()
      expect(wrapper.find('oc-spinner-stub').exists()).toBeFalsy()
    })
  })
  describe('open with default app', () => {
    it('gets called if given via route query param', async () => {
      const { wrapper, mocks } = getMountedWrapper({ openWithDefaultAppQuery: 'true' })
      await (wrapper.vm as any).loadResourcesTask.last
      expect(mocks.openWithDefaultApp).toHaveBeenCalled()
    })
    it('gets not called if not given via route query param', async () => {
      const { wrapper, mocks } = getMountedWrapper()
      await (wrapper.vm as any).loadResourcesTask.last
      expect(mocks.openWithDefaultApp).not.toHaveBeenCalled()
    })
  })
  describe('filter', () => {
    describe('share visibility', () => {
      it('shows filter', () => {
        const { wrapper } = getMountedWrapper()
        expect(wrapper.find('.share-visibility-filter').exists()).toBeTruthy()
        expect(wrapper.find('item-filter-inline-stub').exists()).toBeTruthy()
      })
      it('shows all visible shares', () => {
        const { wrapper } = getMountedWrapper()
        expect(wrapper.findAll('shared-with-me-section-stub').length).toBe(1)
        expect(
          wrapper
            .findComponent<typeof SharedWithMeSection>('shared-with-me-section-stub')
            .props('title')
        ).toEqual('Shares')
      })
      it('shows all hidden shares', async () => {
        const { wrapper } = getMountedWrapper()
        ;(wrapper.vm as any).setAreHiddenFilesShown(mock<InlineFilterOption>({ name: 'hidden' }))
        await wrapper.vm.$nextTick()
        expect(wrapper.findAll('shared-with-me-section-stub').length).toBe(1)
        expect(
          wrapper
            .findComponent<typeof SharedWithMeSection>('shared-with-me-section-stub')
            .props('title')
        ).toEqual('Hidden Shares')
      })
    })
    describe('share type', () => {
      it('shows all available share types as filter option', () => {
        const shareType1 = ShareTypes.user
        const shareType2 = ShareTypes.group
        const { wrapper } = getMountedWrapper({
          files: [
            mock<IncomingShareResource>({ shareTypes: [shareType1.value] }),
            mock<IncomingShareResource>({ shareTypes: [shareType2.value] })
          ]
        })
        const filterItems = wrapper
          .findComponent<typeof ItemFilter>('.share-type-filter')
          .props('items') as ShareType[]

        expect(wrapper.find('.share-type-filter').exists()).toBeTruthy()
        expect(filterItems[0].value).toEqual(shareType1.value)
        expect(filterItems[1].value).toEqual(shareType2.value)
      })
    })
    describe('shared by', () => {
      it('shows all available collaborators as filter option', () => {
        const collaborator1 = { id: 'user1', displayName: 'user1' }
        const collaborator2 = { id: 'user2', displayName: 'user2' }
        const { wrapper } = getMountedWrapper({
          files: [
            mock<IncomingShareResource>({
              sharedBy: [collaborator1],
              shareTypes: [ShareTypes.user.value]
            }),
            mock<IncomingShareResource>({
              sharedBy: [collaborator2],
              shareTypes: [ShareTypes.user.value]
            })
          ]
        })
        const filterItems = wrapper
          .findComponent<typeof ItemFilter>('.shared-by-filter')
          .props('items')
        expect(wrapper.find('.shared-by-filter').exists()).toBeTruthy()
        expect(filterItems).toEqual([collaborator1, collaborator2])
      })
    })
    describe('search', () => {
      it('shows filter', () => {
        const { wrapper } = getMountedWrapper()
        expect(wrapper.find('.search-filter').exists()).toBeTruthy()
      })
      it('filters shares accordingly by name', async () => {
        const { wrapper } = getMountedWrapper({
          files: [
            mock<IncomingShareResource>({
              name: 'share1',
              hidden: false,
              shareTypes: [ShareTypes.user.value]
            }),
            mock<IncomingShareResource>({
              name: 'share2',
              hidden: false,
              shareTypes: [ShareTypes.user.value]
            })
          ]
        })

        await wrapper.vm.$nextTick()
        ;(wrapper.vm as any).filterTerm = 'share1'
        expect((wrapper.vm as any).items.find(({ name }) => name === 'share1')).toBeDefined()
        expect((wrapper.vm as any).items.find(({ name }) => name === 'share2')).toBeUndefined()
      })
    })
  })
})

function getMountedWrapper({
  mocks = {},
  loading = false,
  files = [],
  openWithDefaultAppQuery = ''
}: {
  mocks?: Record<string, unknown>
  files?: IncomingShareResource[]
  loading?: boolean
  openWithDefaultAppQuery?: string
} = {}) {
  vi.mocked(useResourcesViewDefaults).mockImplementation(() =>
    useResourcesViewDefaultsMock({
      paginatedResources: ref(files),
      areResourcesLoading: ref(loading)
    })
  )
  vi.mocked(useSort).mockImplementation((options) => useSortMock({ items: ref(options.items) }))
  // selected share types
  vi.mocked(queryItemAsString).mockImplementationOnce(() => undefined)
  // selected shared by
  vi.mocked(queryItemAsString).mockImplementationOnce(() => undefined)
  // openWithDefaultAppQuery
  vi.mocked(queryItemAsString).mockImplementationOnce(() => openWithDefaultAppQuery)

  const openWithDefaultApp = vi.fn()
  vi.mocked(useOpenWithDefaultApp).mockReturnValue({ openWithDefaultApp })

  const defaultMocks = {
    ...defaultComponentMocks({
      currentRoute: mock<RouteLocation>({ name: 'files-shares-with-me' })
    }),
    ...(mocks && mocks),
    openWithDefaultApp
  }

  return {
    mocks: defaultMocks,
    wrapper: mount(SharedWithMe, {
      global: {
        components: {
          AppBar
        },
        plugins: [...defaultPlugins()],
        mocks: defaultMocks,
        stubs: { ...defaultStubs, itemFilterInline: true, ItemFilter: true }
      }
    })
  }
}
