import { mock } from 'vitest-mock-extended'
import { unref } from 'vue'
import { useFileActionsEnableSync } from '../../../../../src/composables/actions/files/useFileActionsEnableSync'
import { IncomingShareResource } from '@ownclouders/web-client'
import {
  defaultComponentMocks,
  getComposableWrapper,
  RouteLocation
} from '@ownclouders/web-test-helpers'

const sharesWithMeLocation = 'files-shares-with-me'
const sharesWithOthersLocation = 'files-shares-with-others'

describe('acceptShare', () => {
  describe('computed property "actions"', () => {
    describe('isVisible property of returned element', () => {
      it.each([
        { resources: [{ syncEnabled: false }] as IncomingShareResource[], expectedStatus: true },
        { resources: [{ syncEnabled: true }] as IncomingShareResource[], expectedStatus: false }
      ])(
        `should be set according to the resource syncEnabled state if the route name is "${sharesWithMeLocation}"`,
        (inputData) => {
          getWrapper({
            setup: () => {
              const { actions } = useFileActionsEnableSync()

              const resources = inputData.resources
              expect(unref(actions)[0].isVisible({ space: null, resources })).toBe(
                inputData.expectedStatus
              )
            }
          })
        }
      )
      it.each([
        { syncEnabled: false } as IncomingShareResource,
        { syncEnabled: true } as IncomingShareResource
      ])(
        `should be set as false if the route name is other than "${sharesWithMeLocation}"`,
        (resource) => {
          getWrapper({
            routeName: sharesWithOthersLocation,
            setup: () => {
              const { actions } = useFileActionsEnableSync()

              expect(
                unref(actions)[0].isVisible({ space: null, resources: [resource] })
              ).toBeFalsy()
            }
          })
        }
      )
    })
  })
})

function getWrapper({
  setup,
  routeName = sharesWithMeLocation
}: {
  setup: (instance: ReturnType<typeof useFileActionsEnableSync>) => void
  routeName?: string
}) {
  const mocks = defaultComponentMocks({ currentRoute: mock<RouteLocation>({ name: routeName }) })
  return {
    wrapper: getComposableWrapper(setup, {
      mocks,
      provide: mocks
    })
  }
}
