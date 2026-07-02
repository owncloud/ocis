import { mock, mockDeep } from 'vitest-mock-extended'
import { Language } from 'vue3-gettext'
import { UploadResourceConflict } from '../../../../../src/helpers/resource'
import { ResolveStrategy, OcUppyFile, useResourcesStore } from '@ownclouders/web-pkg'
import { Resource } from '@ownclouders/web-client'
import { createTestingPinia } from '@ownclouders/web-test-helpers'

const getResourceConflictInstance = ({
  currentFiles = [mockDeep<Resource>()]
}: {
  currentFiles?: Resource[]
} = {}) => {
  createTestingPinia({ initialState: { resources: { resources: currentFiles } } })
  const resourcesStore = useResourcesStore()
  return new UploadResourceConflict(resourcesStore, mock<Language>())
}

describe('upload helper', () => {
  describe('method "getConflicts"', () => {
    it('should return file and folder conflicts', () => {
      const fileName = 'someFile.txt'
      const folderName = 'someFolder'
      const currentFiles = [
        mockDeep<Resource>({ name: fileName }),
        mockDeep<Resource>({ name: folderName })
      ]
      const filesToUpload = [
        mockDeep<OcUppyFile>({ name: fileName, meta: { relativePath: '', relativeFolder: '' } }),
        mockDeep<OcUppyFile>({
          name: 'anotherFile',
          meta: { relativePath: `/${folderName}/anotherFile` }
        })
      ]
      const instance = getResourceConflictInstance({ currentFiles })
      const conflicts = instance.getConflicts(filesToUpload)

      expect(conflicts.length).toBe(2)
      expect(conflicts).toContainEqual({ name: fileName, type: 'file' })
      expect(conflicts).toContainEqual({ name: folderName, type: 'folder' })
    })
  })
  describe('method "displayOverwriteDialog"', () => {
    it.each([ResolveStrategy.REPLACE, ResolveStrategy.KEEP_BOTH])(
      'should return all files if user chooses replace or keep both for all',
      async (strategy) => {
        const OcUppyFile = mockDeep<OcUppyFile>({
          name: 'test',
          meta: {
            relativeFolder: ''
          },
          xhrUpload: { endpoint: 'https://example.com/upload' }
        })
        const conflict = {
          name: OcUppyFile.name,
          type: 'file'
        }

        const instance = getResourceConflictInstance()
        const resolveFileConflictMethod = vi.fn(() =>
          Promise.resolve({ strategy, doForAllConflicts: true })
        )
        instance.resolveFileExists = resolveFileConflictMethod

        const result = await instance.displayOverwriteDialog([OcUppyFile], [conflict])
        expect(result.length).toBe(1)
        expect(result).toEqual([OcUppyFile])
      }
    )
    it('should return no files if user chooses skip for all', async () => {
      const OcUppyFile = mockDeep<OcUppyFile>({ name: 'test' })
      const conflict = { name: OcUppyFile.name, type: 'file' }

      const instance = getResourceConflictInstance()

      const resolveFileConflictMethod = vi.fn(() =>
        Promise.resolve({ strategy: ResolveStrategy.SKIP, doForAllConflicts: true })
      )
      instance.resolveFileExists = resolveFileConflictMethod

      const result = await instance.displayOverwriteDialog([OcUppyFile], [conflict])
      expect(result.length).toBe(0)
    })
    it('should show dialog once if do for all conflicts is ticked', async () => {
      const OcUppyFileOne = mockDeep<OcUppyFile>({ name: 'test' })
      const OcUppyFileTwo = mockDeep<OcUppyFile>({ name: 'test2' })
      const conflictOne = { name: OcUppyFileOne.name, type: 'file' }
      const conflictTwo = { name: OcUppyFileTwo.name, type: 'file' }

      const instance = getResourceConflictInstance()
      const resolveFileConflictMethod = vi.fn(() =>
        Promise.resolve({ strategy: ResolveStrategy.REPLACE, doForAllConflicts: true })
      )
      instance.resolveFileExists = resolveFileConflictMethod

      await instance.displayOverwriteDialog(
        [OcUppyFileOne, OcUppyFileTwo],
        [conflictOne, conflictTwo]
      )

      expect(resolveFileConflictMethod).toHaveBeenCalledTimes(1)
    })
    it('should show dialog twice if do for all conflicts is ticked and folders and files are uploaded', async () => {
      const OcUppyFileOne = mockDeep<OcUppyFile>({ name: 'test' })
      const OcUppyFileTwo = mockDeep<OcUppyFile>({ name: 'folder' })
      const conflictOne = {
        name: OcUppyFileOne.name,
        type: 'file',
        meta: { relativeFolder: '/' }
      }
      const conflictTwo = { name: OcUppyFileTwo.name, type: 'folder' }

      const instance = getResourceConflictInstance()
      instance.resolveFileExists = vi.fn(() =>
        Promise.resolve({ strategy: ResolveStrategy.REPLACE, doForAllConflicts: true })
      )

      await instance.displayOverwriteDialog(
        [OcUppyFileOne, OcUppyFileTwo],
        [conflictOne, conflictTwo]
      )

      expect(instance.resolveFileExists).toHaveBeenCalledTimes(2)
    })
  })
})
