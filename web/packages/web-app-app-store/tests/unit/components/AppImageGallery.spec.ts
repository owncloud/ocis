import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import AppImageGallery from '../../../src/components/AppImageGallery.vue'
import { App, AppBadge, AppImage, BADGE_COLORS } from '../../../src/types'

const coverImageWithUrl: AppImage = {
  url: 'https://example.com/cover.jpg',
  caption: 'Cover image'
}
const coverImageWithoutUrl: AppImage = {
  caption: 'Trololo',
  url: ''
}
const screenshot1: AppImage = {
  url: 'https://example.com/screenshot1.jpg',
  caption: 'Screenshot 1'
}
const screenshot2: AppImage = {
  url: 'https://example.com/screenshot2.jpg',
  caption: 'Screenshot 2'
}
const screenshot3: AppImage = {
  url: 'https://example.com/screenshot3.jpg',
  caption: 'Screenshot 3'
}
const screenshots = [screenshot1, screenshot2, screenshot3]

const selectors = {
  badge: '.app-image-ribbon',
  image: '.app-image img',
  imageFallback: '.app-image .fallback-icon',
  pagination: '.app-image-navigation',
  paginationPrev: '[data-testid="prev-image"]',
  paginationNext: '[data-testid="next-image"]',
  paginationSet: '[data-testid="set-image"]'
}

describe('AppImageGallery.vue', () => {
  describe('badges', () => {
    it('renders no ribbon if the app has no badge', () => {
      const { wrapper } = getWrapper({})
      expect(wrapper.find(selectors.badge).exists()).toBeFalsy()
    })
    it('renders a ribbon if the app has a badge', () => {
      const badge = { label: 'New', color: BADGE_COLORS[1] }
      const { wrapper } = getWrapper({ badge })
      expect(wrapper.find(selectors.badge).exists()).toBeTruthy()
      expect(wrapper.find(selectors.badge).text()).toBe(badge.label)
      expect(
        wrapper.find(selectors.badge).element.className.includes(`app-image-ribbon-${badge.color}`)
      ).toBeTruthy()
    })
  })
  describe('current image', () => {
    it('renders oc-img if the image has a url', () => {
      const { wrapper } = getWrapper({ coverImage: coverImageWithUrl })
      expect(wrapper.find(selectors.image).exists()).toBeTruthy()
      expect(wrapper.find(selectors.image).attributes().src).toBe(coverImageWithUrl.url)
      expect(wrapper.find(selectors.imageFallback).exists()).toBeFalsy()
    })
    it('renders oc-icon if the image has no url', () => {
      const { wrapper } = getWrapper({ coverImage: coverImageWithoutUrl })
      expect(wrapper.find(selectors.image).exists()).toBeFalsy()
      expect(wrapper.find(selectors.imageFallback).exists()).toBeTruthy()
    })
    it('renders oc-icon if the app has no coverImage field', () => {
      const { wrapper } = getWrapper({ coverImage: null })
      expect(wrapper.find(selectors.image).exists()).toBeFalsy()
      expect(wrapper.find(selectors.imageFallback).exists()).toBeTruthy()
    })
  })
  describe('navigation', () => {
    it('has no navigation if there is only a single image', () => {
      const { wrapper } = getWrapper({ showPagination: true, coverImage: coverImageWithUrl })
      expect(wrapper.find(selectors.pagination).exists()).toBeFalsy()
    })
    it('has no navigation if it is disabled via prop', () => {
      const { wrapper } = getWrapper({
        showPagination: false,
        coverImage: coverImageWithUrl,
        screenshots
      })
      expect(wrapper.find(selectors.pagination).exists()).toBeFalsy()
    })
    describe('is visible', () => {
      it('has a pagination container', () => {
        const { wrapper } = getWrapper({
          showPagination: true,
          coverImage: coverImageWithUrl,
          screenshots
        })
        expect(wrapper.find(selectors.pagination).exists()).toBeTruthy()
      })
      it('has a prev button which cycles through images backwards', async () => {
        const { wrapper } = getWrapper({
          showPagination: true,
          coverImage: coverImageWithUrl,
          screenshots
        })
        const button = wrapper.find(selectors.paginationPrev)
        expect(button.exists()).toBeTruthy()
        const images = [coverImageWithUrl, ...screenshots]
        for (let i = 1; i <= images.length; i++) {
          await button.trigger('click')
          expect(wrapper.find(selectors.image).attributes().src).toBe(images[images.length - i].url)
        }
      })
      it('has a next button which cycles through images forward', async () => {
        const { wrapper } = getWrapper({
          showPagination: true,
          coverImage: coverImageWithUrl,
          screenshots
        })
        const button = wrapper.find(selectors.paginationNext)
        const images = [coverImageWithUrl, ...screenshots]
        expect(button.exists()).toBeTruthy()
        for (let i = 1; i <= images.length; i++) {
          await button.trigger('click')
          expect(wrapper.find(selectors.image).attributes().src).toBe(images[i % images.length].url)
        }
      })
      it('has a set-button per image which changes the current image', async () => {
        const { wrapper } = getWrapper({
          showPagination: true,
          coverImage: coverImageWithUrl,
          screenshots
        })
        const buttons = wrapper.findAll(selectors.paginationSet)
        const images = [coverImageWithUrl, ...screenshots]
        expect(buttons.length).toBe(images.length)
        const indices = [2, 0, 1, 3, 3, 3, 0]
        for (let i = 0; i < indices.length; i++) {
          await buttons[indices[i]].trigger('click')
          expect(wrapper.find(selectors.image).attributes().src).toBe(images[indices[i]].url)
        }
      })
    })
  })
})

const getWrapper = ({
  showPagination,
  badge,
  coverImage,
  screenshots = []
}: {
  showPagination?: boolean
  badge?: AppBadge
  coverImage?: AppImage
  screenshots?: AppImage[]
}) => {
  const app = { ...mock<App>({}), badge, coverImage, screenshots }

  return {
    wrapper: mount(AppImageGallery, {
      props: {
        app,
        showPagination
      },
      global: {
        plugins: defaultPlugins()
      }
    })
  }
}
