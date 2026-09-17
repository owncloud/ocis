import { graph } from '../../../src/graph'
import { FetchClient } from '../../../src/http'

const drive = {
  id: 'storage-id',
  name: 'Alice',
  webUrl: 'https://host/f/storage-id',
  quota: { total: 500, used: 0 }
}

describe('graph drives', () => {
  let fetchMock: ReturnType<typeof vi.fn>

  beforeEach(() => {
    fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify(drive), {
        status: 200,
        headers: { 'Content-Type': 'application/json' }
      })
    )
    vi.stubGlobal('fetch', fetchMock)
  })

  describe('updateDrive', () => {
    /**
     * A `readOnly: true` in the spec makes the generator omit the field from its serializers,
     * so a field oCIS does read on write but that the spec marks read-only leaves as `{}` or
     * disappears entirely — the request still succeeds, it just does nothing. Type checking
     * cannot catch that, hence these tests assert on the serialized body.
     */
    it('sends the quota in the request body', async () => {
      await graph('https://host', new FetchClient()).drives.updateDrive(
        drive.id,
        { name: drive.name, quota: { total: 500 } },
        {}
      )

      const body = JSON.parse(fetchMock.mock.calls[0][1].body)
      expect(body.quota).toEqual({ total: 500 })
    })

    it.each(['image', 'readme'])('sends the id of the %s special folder', async (name) => {
      await graph('https://host', new FetchClient()).drives.updateDrive(
        drive.id,
        { name: drive.name, special: [{ specialFolder: { name }, id: 'file-id' }] },
        {}
      )

      const body = JSON.parse(fetchMock.mock.calls[0][1].body)
      expect(body.special).toEqual([{ specialFolder: { name }, id: 'file-id' }])
    })

    it('returns the updated quota on the space', async () => {
      const space = await graph('https://host', new FetchClient()).drives.updateDrive(
        drive.id,
        { name: drive.name, quota: { total: 500 } },
        {}
      )

      expect(space.spaceQuota).toEqual({ total: 500, used: 0 })
    })
  })
})
