import { HttpError, type HttpResponse } from '@ownclouders/web-client'

/**
 * Builds the envelope HttpClient resolves with.
 */
export const mockHttpResponse = <T>(
  data: T = {} as T,
  {
    status = 200,
    statusText = 'OK',
    headers = {}
  }: {
    status?: number
    statusText?: string
    headers?: Record<string, string>
  } = {}
): HttpResponse<T> => ({
  data,
  status,
  statusText,
  headers: new Headers(headers)
})

/**
 * Builds a rejected promise carrying the HttpError the fetch core throws. Note that
 * callers branch on `statusCode`, never `status`.
 */
export const mockHttpError = (
  status = 500,
  data: unknown = undefined,
  message = ''
): Promise<never> =>
  Promise.reject(new HttpError(message, new Response(null, { status }), status, data))
