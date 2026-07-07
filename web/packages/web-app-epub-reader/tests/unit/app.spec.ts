import {
  PartialComponentProps,
  defaultPlugins,
  getOcSelectOptions,
  mount,
  nextTicks
} from '@ownclouders/web-test-helpers'
import App from '../../src/App.vue'
import { useLocalStorage } from '@ownclouders/web-pkg'
import { Resource } from '@ownclouders/web-client'
import { mock } from 'vitest-mock-extended'
import { ref } from 'vue'

vi.mock('@ownclouders/web-pkg', async (importOriginal) => ({
  ...(await importOriginal<any>()),
  useLocalStorage: vi.fn()
}))

vi.mock('epubjs', () => ({
  __esModule: true,
  default: vi.fn(() => {
    return {
      loaded: {
        navigation: Promise.resolve({
          toc: [
            { id: '1', label: 'Chapter 1', href: 'c1' },
            { id: '2', label: 'Chapter 2', href: 'c2' }
          ]
        })
      },
      renderTo: vi.fn(() => ({
        on: vi.fn(),
        themes: {
          register: vi.fn(),
          select: vi.fn(),
          fontSize: vi.fn()
        },
        display: vi.fn(),
        prev: vi.fn(),
        next: vi.fn()
      }))
    }
  })
}))

const selectors = {
  increaseFontSize: '.epub-reader-controls-font-size-increase',
  decreaseFontSize: '.epub-reader-controls-font-size-decrease',
  resetFontSize: '.epub-reader-controls-font-size-reset',
  chaptersListItem: '.epub-reader-chapters-list-item',
  chaptersSelect: '.epub-reader-controls-chapters-select',
  navigateLeft: '.epub-reader-navigate-left',
  navigateRight: '.epub-reader-navigate-right'
}
describe('Epub reader app', () => {
  it('renders correctly', async () => {
    const { wrapper } = getWrapper()
    await nextTicks(2)
    expect(wrapper.html()).toMatchSnapshot()
  })
  describe('theme', () => {
    it('sets the theme based on current theme setting', async () => {
      const { wrapper } = getWrapper({ localStorageGeneral: { fontSizePercentage: 50 } })
      await nextTicks(2)
      expect((wrapper.vm as any).rendition.themes.select).toHaveBeenCalledWith('light')
    })
  })
  describe('font size', () => {
    it('initializes with default font size percentage', async () => {
      const { wrapper } = getWrapper()
      await nextTicks(2)
      expect((wrapper.vm as any).rendition.themes.fontSize).toHaveBeenCalledWith('100%')
    })
    it('initializes with local storage font size when set', async () => {
      const { wrapper } = getWrapper({ localStorageGeneral: { fontSizePercentage: 50 } })
      await nextTicks(2)
      expect((wrapper.vm as any).rendition.themes.fontSize).toHaveBeenCalledWith('50%')
    })
    describe('increase font size button', () => {
      it('increases font size when clicked', async () => {
        const { wrapper } = getWrapper()
        await nextTicks(2)
        await wrapper.find(selectors.increaseFontSize).trigger('click')
        expect((wrapper.vm as any).rendition.themes.fontSize).toHaveBeenCalledWith('110%')
      })
      it('is disabled when "MAX_FONT_SIZE_PERCENTAGE" is reached', () => {
        const { wrapper } = getWrapper({ localStorageGeneral: { fontSizePercentage: 150 } })
        expect(
          wrapper.find<HTMLButtonElement>(selectors.increaseFontSize).element.disabled
        ).toBeTruthy()
      })
    })
    describe('decrease font size button', () => {
      it('decreases font size when clicked', async () => {
        const { wrapper } = getWrapper()
        await nextTicks(2)
        await wrapper.find(selectors.decreaseFontSize).trigger('click')
        expect((wrapper.vm as any).rendition.themes.fontSize).toHaveBeenCalledWith('90%')
      })
      it('is disabled when "MIN_FONT_SIZE_PERCENTAGE" is reached', () => {
        const { wrapper } = getWrapper({ localStorageGeneral: { fontSizePercentage: 50 } })
        expect(
          wrapper.find<HTMLButtonElement>(selectors.decreaseFontSize).element.disabled
        ).toBeTruthy()
      })
    })
    describe('reset font size button', () => {
      it('resets font size when clicked', async () => {
        const { wrapper } = getWrapper({ localStorageGeneral: { fontSizePercentage: 50 } })
        await nextTicks(2)
        await wrapper.find(selectors.resetFontSize).trigger('click')
        expect((wrapper.vm as any).rendition.themes.fontSize).toHaveBeenCalledWith('100%')
      })
      it('shows the current font size', async () => {
        const { wrapper } = getWrapper()
        await nextTicks(2)
        await wrapper.find(selectors.decreaseFontSize).trigger('click')
        expect(wrapper.find(selectors.resetFontSize).text()).toBe('90%')
      })
    })
  })
  describe('location', () => {
    it('initializes with local storage location when set', async () => {
      const { wrapper } = getWrapper({
        localStorageResource: {
          currentLocation: { start: { cfi: 'epubcfi(/6/4!/4/4/14/2/150/2/1:23)' } }
        }
      })
      await nextTicks(2)
      expect((wrapper.vm as any).rendition.display).toHaveBeenCalledWith(
        'epubcfi(/6/4!/4/4/14/2/150/2/1:23)'
      )
    })
  })
  describe('chapters', () => {
    describe('chapters list', () => {
      it('renders correctly', async () => {
        const { wrapper } = getWrapper()
        await nextTicks(2)
        const chapterElements = wrapper.findAll(selectors.chaptersListItem)
        expect(chapterElements.length).toEqual(2)
        expect(chapterElements[0].text()).toEqual('Chapter 1')
        expect(chapterElements[1].text()).toEqual('Chapter 2')
      })
      it('calls method "display" when item is clicked', async () => {
        const { wrapper } = getWrapper()
        await nextTicks(2)
        const chapterElements = wrapper.findAll(selectors.chaptersListItem)
        await chapterElements[1].find('.oc-button').trigger('click')
        expect((wrapper.vm as any).rendition.display).toHaveBeenCalledWith('c2')
      })
    })
    describe('chapters select', () => {
      it('renders correctly', async () => {
        const { wrapper } = getWrapper()
        await nextTicks(2)
        const chapterElements = await getOcSelectOptions(wrapper, selectors.chaptersSelect)
        expect(chapterElements.length).toEqual(2)
        expect(chapterElements[0].text()).toEqual('Chapter 1')
        expect(chapterElements[1].text()).toEqual('Chapter 2')
      })
      it('calls method "display" when item is clicked', async () => {
        const { wrapper } = getWrapper()
        await nextTicks(2)
        const chapterElements = await getOcSelectOptions(wrapper, selectors.chaptersSelect, {
          close: false
        })
        await chapterElements[1].trigger('click')
        expect((wrapper.vm as any).rendition.display).toHaveBeenCalledWith('c2')
      })
    })
  })
  describe('navigate', () => {
    describe('keyboard navigation', () => {
      it('calls method "prev" when left arrow key is pressed', async () => {
        const { wrapper } = getWrapper()
        await nextTicks(2)
        const keyboardEvent = new KeyboardEvent('keydown', { key: 'ArrowLeft' })
        document.dispatchEvent(keyboardEvent)
        expect((wrapper.vm as any).rendition.prev).toHaveBeenCalled()
      })
      it('calls method "next" when right arrow key is pressed', async () => {
        const { wrapper } = getWrapper()
        await nextTicks(2)
        const keyboardEvent = new KeyboardEvent('keydown', { key: 'ArrowRight' })
        document.dispatchEvent(keyboardEvent)
        expect((wrapper.vm as any).rendition.next).toHaveBeenCalled()
      })
    })
    describe('navigate left button', () => {
      it('calls method "prev" when clicked', async () => {
        const { wrapper } = getWrapper()
        await nextTicks(2)
        await wrapper.find(selectors.navigateLeft).trigger('click')
        expect((wrapper.vm as any).rendition.prev).toHaveBeenCalled()
      })
    })
    describe('navigate right button', () => {
      it('calls method "next" when clicked', async () => {
        const { wrapper } = getWrapper()
        await nextTicks(2)
        await wrapper.find(selectors.navigateRight).trigger('click')
        expect((wrapper.vm as any).rendition.next).toHaveBeenCalled()
      })
    })
  })
})

function getWrapper({
  propsData = {},
  localStorageGeneral = {},
  localStorageResource = {}
}: {
  propsData?: PartialComponentProps<typeof App>
  localStorageGeneral?: Record<string, unknown>
  localStorageResource?: Record<string, unknown>
} = {}) {
  vi.mocked(useLocalStorage<unknown>).mockImplementationOnce(() => ref(localStorageGeneral))
  vi.mocked(useLocalStorage<unknown>).mockImplementationOnce(() => ref(localStorageResource))

  return {
    wrapper: mount(App, {
      props: {
        currentContent: '',
        resource: mock<Resource>({
          id: '1'
        }),
        ...propsData
      },
      global: {
        plugins: [...defaultPlugins()]
      }
    })
  }
}
