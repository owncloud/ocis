import {
  defaultComponentMocks,
  defaultPlugins,
  mockHttpResponse,
  shallowMount
} from '@ownclouders/web-test-helpers'
import { FederatedConnection } from '@ownclouders/web-pkg'
import ConnectionsPanel from '../../../src/views/ConnectionsPanel.vue'

describe('ConnectionsPanel', () => {
  describe('deleteConnection', () => {
    it('sends the user id and idp as the request body', async () => {
      const { wrapper, mocks } = getWrapper()
      const connection = { id: 'c1', user_id: 'einstein', idp: 'https://idp.example.org' }

      await (wrapper.vm as any).deleteConnection(connection)

      expect(mocks.$clientService.httpAuthenticated.delete).toHaveBeenCalledWith(
        '/sciencemesh/delete-accepted-user',
        { user_id: 'einstein', idp: 'https://idp.example.org' }
      )
    })
  })
})

function getWrapper() {
  const mocks = defaultComponentMocks()
  mocks.$clientService.httpAuthenticated.delete.mockResolvedValue(mockHttpResponse({}))

  return {
    mocks,
    wrapper: shallowMount(ConnectionsPanel, {
      props: { connections: [] as FederatedConnection[], loading: false },
      global: {
        plugins: [...defaultPlugins()],
        mocks,
        provide: mocks
      }
    })
  }
}
