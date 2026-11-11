import { graph } from '../../../src/graph'
import { FetchClient } from '../../../src/http'

const respondWith = (payload: unknown) => {
  const fetchMock = vi.fn().mockResolvedValue(
    new Response(JSON.stringify(payload), {
      status: 200,
      headers: { 'Content-Type': 'application/json' }
    })
  )
  vi.stubGlobal('fetch', fetchMock)
  return fetchMock
}

const client = () => graph('https://host', new FetchClient())

describe('graph users', () => {
  /**
   * oCIS returns `attributes` on users, filled from OCIS_USER_SEARCH_DISPLAYED_ATTRIBUTES, and
   * the sharing autocomplete displays them. The upstream spec does not declare the field, and
   * the generator rebuilds every response from the declared fields only, so it has to be
   * declared during generation or it gets dropped on the way out of the client.
   */
  it('keeps the oCIS-only attributes when listing users', async () => {
    respondWith({
      value: [
        {
          id: 'alice',
          displayName: 'Alice',
          onPremisesSamAccountName: 'alice',
          attributes: ['Engineering', 'Vienna']
        }
      ]
    })

    const [user] = await client().users.listUsers({})

    expect(user.attributes).toEqual(['Engineering', 'Vienna'])
  })

  it('keeps the oCIS-only attributes when getting a single user', async () => {
    respondWith({
      id: 'alice',
      displayName: 'Alice',
      onPremisesSamAccountName: 'alice',
      attributes: ['Engineering']
    })

    const user = await client().users.getUser('alice')

    expect(user.attributes).toEqual(['Engineering'])
  })

  it('omits the attributes when the server does not send any', async () => {
    respondWith({ id: 'alice', displayName: 'Alice', onPremisesSamAccountName: 'alice' })

    const user = await client().users.getUser('alice')

    expect(user.attributes).toBeUndefined()
  })

  /**
   * The generator materializes every declared field, so a response that omits `accountEnabled`
   * still yields an own property holding `undefined`. Consumers therefore cannot tell an unset
   * field apart with `'accountEnabled' in user` and have to compare against `undefined`.
   */
  it('reports an unset accountEnabled as undefined', async () => {
    respondWith({ id: 'alice', displayName: 'Alice', onPremisesSamAccountName: 'alice' })

    const user = await client().users.getUser('alice')

    expect(user.accountEnabled).toBeUndefined()
  })

  it.each([true, false])('passes an accountEnabled of %s through', async (accountEnabled) => {
    respondWith({
      id: 'alice',
      displayName: 'Alice',
      onPremisesSamAccountName: 'alice',
      accountEnabled
    })

    const user = await client().users.getUser('alice')

    expect(user.accountEnabled).toBe(accountEnabled)
  })
})
