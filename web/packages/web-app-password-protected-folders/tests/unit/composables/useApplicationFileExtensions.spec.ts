import { defaultComponentMocks, getComposableWrapper } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import { Resource } from '@ownclouders/web-client'
import { useApplicationFileExtensions } from '../../../src/composables/useApplicationFileExtensions'
import { shareType } from '../../../../design-system/src/utils/shareType'

describe('application file extensions', () => {
  it('should hide the new file menu when the current folder is a link share', () => {
    const currentFolder = mock<Resource>({ path: '/current/folder', shareTypes: [shareType.link] })
    getWrapper({
      setup(instance) {
        expect(instance.at(0).newFileMenu.isVisible({ currentFolder })).toBe(false)
      }
    })
  })

  it('should not hide the new file menu when the current folder is not a link share', () => {
    const currentFolder = mock<Resource>({ path: '/current/folder', shareTypes: [] })
    getWrapper({
      setup(instance) {
        expect(instance.at(0).newFileMenu.isVisible({ currentFolder })).toBe(true)
      }
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (
    instance: ReturnType<typeof useApplicationFileExtensions>,
    mocks: ReturnType<typeof defaultComponentMocks>
  ) => void
}) {
  const mocks = defaultComponentMocks()

  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useApplicationFileExtensions()
        setup(instance, mocks)
      },
      {
        mocks,
        provide: mocks
      }
    )
  }
}
