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
     * The spec marks every `Quota` field read-only, which the generator honours by omitting
     * them from its serializers. Without stripping the annotation before generating, the
     * quota would silently leave as `{}` and the drive would keep its old limit.
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
