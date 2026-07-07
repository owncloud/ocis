import { cacheService } from '../../../src/services'
import { Cache } from '../../../src/helpers/cache'

describe('cache', () => {
  describe('cacheService', () => {
    test('filePreview', () => {
      const filePreviewCache = cacheService.filePreview
      expect(filePreviewCache).toBeInstanceOf(Cache)
    })

    test('avatarUrl', () => {
      const avatarUrlCache = cacheService.avatarUrl
      expect(avatarUrlCache).toBeInstanceOf(Cache)
    })
  })
})
