import { mock } from 'vitest-mock-extended'
import { ref } from 'vue'
import { FileContext, useAppDefaults, AppConfigObject } from '@ownclouders/web-pkg'
import { FileResource, Resource } from '@ownclouders/web-client'
import { GetFileContentsResponse } from '@ownclouders/web-client/webdav'

export const useAppDefaultsMock = (
  options: Partial<ReturnType<typeof useAppDefaults>> = {}
): ReturnType<typeof useAppDefaults> => {
  return {
    isPublicLinkContext: ref(false),
    currentFileContext: ref(mock<FileContext>()),
    applicationConfig: ref(mock<AppConfigObject>()),
    closeApp: vi.fn(),
    replaceInvalidFileRoute: vi.fn(),
    getUrlForResource: vi.fn(),
    revokeUrl: vi.fn(),
    getFileInfo: vi.fn().mockImplementation(() => mock<Resource>()),
    getFileContents: vi.fn().mockImplementation(() => mock<GetFileContentsResponse>({ body: '' })),
    putFileContents: vi.fn().mockImplementation(() => mock<FileResource>()),
    isFolderLoading: ref(false),
    activeFiles: ref([]),
    loadFolderForFileContext: vi.fn(),
    makeRequest: vi.fn().mockResolvedValue({ status: 200 }),
    closed: ref(false),
    ...options
  }
}
