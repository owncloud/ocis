import {
  maintenanceResponseHandler,
  shouldResponseTriggerMaintenance
} from '../../../src/helpers/maintenance'
import type { OnResponseArgs } from '../../../src/http'

describe('shouldResponseTriggerMaintenance', () => {
  it('is true for a 503 on a regular endpoint', () => {
    expect(shouldResponseTriggerMaintenance(503, '/remote.php/webdav/foo')).toBe(true)
  })

  it('is false for a 503 on an excluded endpoint', () => {
    expect(
      shouldResponseTriggerMaintenance(
        503,
        'ocs/v2.php/apps/notifications/api/v1/notifications/sse'
      )
    ).toBe(false)
  })

  it('is false for any other status', () => {
    expect(shouldResponseTriggerMaintenance(404, '/remote.php/webdav/foo')).toBe(false)
    expect(shouldResponseTriggerMaintenance(500, '/remote.php/webdav/foo')).toBe(false)
    expect(shouldResponseTriggerMaintenance(200, '/remote.php/webdav/foo')).toBe(false)
  })
})

describe('maintenanceResponseHandler', () => {
  const args = (overrides: Partial<OnResponseArgs>): OnResponseArgs => ({
    response: null,
    status: 500,
    requestUrl: '/remote.php/webdav/foo',
    ...overrides
  })

  it('clears maintenance mode on a successful response', () => {
    const onSetMaintenance = vi.fn()
    const handler = maintenanceResponseHandler(onSetMaintenance)

    handler(args({ response: { ok: true } as Response, status: 200 }))

    expect(onSetMaintenance).toHaveBeenCalledWith(false)
  })

  it('runs onSuccess bookkeeping only on a successful response', () => {
    const onSuccess = vi.fn()
    const handler = maintenanceResponseHandler(vi.fn(), { onSuccess })

    handler(args({ response: { ok: true } as Response, status: 200 }))
    expect(onSuccess).toHaveBeenCalledTimes(1)

    handler(args({ response: { ok: false } as Response, status: 404 }))
    expect(onSuccess).toHaveBeenCalledTimes(1)
  })

  it('sets maintenance mode on a 503', () => {
    const onSetMaintenance = vi.fn()
    const handler = maintenanceResponseHandler(onSetMaintenance)

    handler(args({ response: { ok: false } as Response, status: 503 }))

    expect(onSetMaintenance).toHaveBeenCalledWith(true)
  })

  it('clears maintenance mode on an unrelated error, since the server did answer', () => {
    const onSetMaintenance = vi.fn()
    const handler = maintenanceResponseHandler(onSetMaintenance)

    handler(args({ response: { ok: false } as Response, status: 404 }))
    expect(onSetMaintenance).toHaveBeenCalledWith(false)

    onSetMaintenance.mockClear()
    handler(args({ response: { ok: false } as Response, status: 500 }))
    expect(onSetMaintenance).toHaveBeenCalledWith(false)
  })

  it('clears maintenance mode for a 503 on an excluded endpoint', () => {
    const onSetMaintenance = vi.fn()
    const handler = maintenanceResponseHandler(onSetMaintenance)

    handler(
      args({
        response: { ok: false } as Response,
        status: 503,
        requestUrl: 'ocs/v2.php/apps/notifications/api/v1/notifications/sse'
      })
    )

    expect(onSetMaintenance).toHaveBeenCalledWith(false)
  })

  it('leaves maintenance mode untouched on a transport failure', () => {
    const onSetMaintenance = vi.fn()
    const handler = maintenanceResponseHandler(onSetMaintenance)

    handler(args({ response: null, status: 500 }))

    expect(onSetMaintenance).not.toHaveBeenCalled()
  })
})
