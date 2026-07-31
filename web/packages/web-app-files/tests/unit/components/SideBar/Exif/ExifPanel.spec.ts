import ExifPanel from '../../../../../src/components/SideBar/Exif/ExifPanel.vue'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Photo, Image, GeoCoordinates } from '@ownclouders/web-client/graph/generated'
import { Resource } from '@ownclouders/web-client'
import { formatDateFromISO } from '@ownclouders/web-pkg'

describe('Exif SideBar Panel', () => {
  describe('photo metadata', () => {
    const keys = [
      'cameraMake',
      'cameraModel',
      'focalLength',
      'fNumber',
      'exposureTime',
      'iso',
      'orientation',
      'takenDateTime'
    ]
    it.each(keys)('shows value in panel for key "%s"', (key) => {
      const resource = mock<Resource>({
        photo: mock<Photo>({
          cameraMake: 'Canon',
          cameraModel: 'Canon EOS 1300D',
          focalLength: 222,
          fNumber: 5.6,
          exposureNumerator: 1,
          exposureDenominator: 20,
          iso: 320,
          orientation: 1,
          takenDateTime: '2017-02-11T14:54:50Z'
        })
      })
      const expectedValues = {
        cameraMake: resource.photo.cameraMake,
        cameraModel: resource.photo.cameraModel,
        focalLength: `${resource.photo.focalLength} mm`,
        fNumber: `f/${resource.photo.fNumber}`,
        exposureTime: `${resource.photo.exposureNumerator}/${resource.photo.exposureDenominator}`,
        iso: resource.photo.iso.toString(),
        orientation: resource.photo.orientation.toString(),
        takenDateTime: formatDateFromISO(resource.photo.takenDateTime, 'en')
      }
      const { wrapper } = createWrapper({ resource })
      expect(wrapper.find(`[data-testid="exif-panel-${key}"]`).text()).toBe(expectedValues[key])
    })
    it.each(keys)('shows "-" in panel if key "%s" has no value in provided data', (key) => {
      const emptyPhotoResourceMock = mock<Resource>({})
      const { wrapper } = createWrapper({ resource: emptyPhotoResourceMock })
      expect(wrapper.find(`[data-testid="exif-panel-${key}"]`).text()).toEqual('-')
    })
  })
  describe('image metadata', () => {
    const dimensionsSelector = '[data-testid="exif-panel-dimensions"]'
    it('shows "width x height" if both props provided in data', () => {
      const resource = mock<Resource>({
        image: mock<Image>({
          width: 5000,
          height: 3000
        })
      })
      const { wrapper } = createWrapper({ resource })
      expect(wrapper.find(dimensionsSelector).text()).toEqual('5000x3000')
    })
    it.each([
      { width: 5000, height: undefined },
      { width: undefined, height: 3000 },
      { width: undefined, height: undefined }
    ])('shows "-" if width or height is missing in data', (options) => {
      const resource = mock<Resource>({
        image: mock<Image>(options)
      })
      const { wrapper } = createWrapper({ resource })
      expect(wrapper.find(dimensionsSelector).text()).toEqual('-')
    })
    it.each([5, 6, 7, 8])(
      'flips width and height because photo orientation "%s" indicates portrait mode',
      (orientation) => {
        const resource = mock<Resource>({
          image: mock<Image>({ width: 5000, height: 3000 }),
          photo: mock<Photo>({ orientation })
        })
        const { wrapper } = createWrapper({ resource })
        expect(wrapper.find(dimensionsSelector).text()).toBe('3000x5000')
      }
    )
  })
  describe('location metadata', () => {
    const locationSelector = '[data-testid="exif-panel-location"]'
    it('shows "latitude, longitude" if both props provided in data', () => {
      const resource = mock<Resource>({
        location: mock<GeoCoordinates>({
          latitude: 51.30044714422953,
          longitude: 7.373170282627126
        })
      })
      const { wrapper } = createWrapper({ resource })
      expect(wrapper.find(locationSelector).text()).toBe(
        `${resource.location.latitude}, ${resource.location.longitude}`
      )
    })
    it.each([
      { latitude: 51.30044714422953, longitude: undefined },
      { latitude: undefined, longitude: 7.373170282627126 },
      { latitude: undefined, longitude: undefined }
    ])('shows "-" if latitude or longitude is missing in data', ({ latitude, longitude }) => {
      const resource = mock<Resource>({
        location: mock<GeoCoordinates>({ latitude, longitude })
      })
      const { wrapper } = createWrapper({ resource })
      expect(wrapper.find(locationSelector).text()).toEqual('-')
    })
  })
})

function createWrapper({ resource }: { resource: Resource }) {
  return {
    wrapper: shallowMount(ExifPanel, {
      global: {
        plugins: [...defaultPlugins({})],
        provide: {
          resource
        }
      }
    })
  }
}
