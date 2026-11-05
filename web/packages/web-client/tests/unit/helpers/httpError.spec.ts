/**
 * @vitest-environment node
 */
import { createHttpError } from '../../../src/helpers/httpError'

describe('createHttpError', () => {
  it('keeps the status code and the request id of the failed request', () => {
    const error = createHttpError({ message: 'Bad Gateway', statusCode: 502, xReqId: 'req-1' })

    expect(error.statusCode).toBe(502)
    expect(error.response.headers.get('x-request-id')).toBe('req-1')
  })

  it('sets no request id when the request got no response', () => {
    const error = createHttpError({
      message: 'Failed to fetch',
      statusCode: undefined,
      xReqId: undefined
    })

    expect(error.response.headers.get('x-request-id')).toBeNull()
  })
})
