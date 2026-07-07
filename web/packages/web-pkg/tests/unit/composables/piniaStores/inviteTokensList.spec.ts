import { createPinia, setActivePinia } from 'pinia'

import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import { useInviteTokensListStore } from '../../../../src/composables/piniaStores'

describe('useInviteTokensList', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  test('to add new token', () => {
    const token = {
      id: 'token1',
      token: 'token1',
      expiration: new Date('2021-05-03T00:00:00.000Z'),
      expirationSeconds: 1620000000,
      description: 'desc1'
    }
    getWrapper({
      setup: (instance) => {
        instance.addToken(token)
        expect(instance.getTokensList()).toEqual([token])
      }
    })
  })

  test('to add new token to list of existing tokens', () => {
    const tokensList = [
      {
        id: 'token1',
        token: 'token1',
        expiration: new Date('2021-05-03T00:00:00.000Z'),
        expirationSeconds: 1620000000,
        description: 'desc1'
      },
      {
        id: 'token2',
        token: 'token2',
        expiration: new Date('2021-05-03T00:00:00.000Z'),
        expirationSeconds: 1620000000,
        description: 'desc2'
      }
    ]
    const tokenToAdd = {
      id: 'token3',
      token: 'token3',
      expiration: new Date('2021-05-03T00:00:00.000Z'),
      expirationSeconds: 1620000000,
      description: 'desc3'
    }

    getWrapper({
      setup: (instance) => {
        instance.setTokensList(tokensList)
        instance.addToken(tokenToAdd)
        expect(instance.getTokensList()).toHaveLength(3)
        expect(instance.getTokensList()).toContainEqual(tokenToAdd)
      }
    })
  })

  test('to set tokens list', () => {
    const tokensList = [
      {
        id: 'token1',
        token: 'token1',
        expiration: new Date('2021-05-03T00:00:00.000Z'),
        expirationSeconds: 1620000000,
        description: 'desc1'
      },
      {
        id: 'token2',
        token: 'token2',
        expiration: new Date('2021-05-03T00:00:00.000Z'),
        expirationSeconds: 1620000000,
        description: 'desc2'
      }
    ]

    getWrapper({
      setup: (instance) => {
        instance.setTokensList(tokensList)
        expect(instance.getTokensList()).toEqual(tokensList)
      }
    })
  })

  test('to set last and get created token', () => {
    getWrapper({
      setup: (instance) => {
        instance.setLastCreatedToken('token@example.com')
        expect(instance.getLastCreatedToken()).toEqual('token@example.com')
      }
    })
  })
})

function getWrapper({
  setup
}: {
  setup: (instance: ReturnType<typeof useInviteTokensListStore>) => void
}) {
  return {
    wrapper: getComposableWrapper(
      () => {
        const instance = useInviteTokensListStore()
        setup(instance)
      },
      { pluginOptions: { pinia: false } }
    )
  }
}
