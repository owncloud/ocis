import { HttpError, httpHeaders, type HttpResponse } from '@ownclouders/web-client'

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
  headers: httpHeaders(new Headers(headers))
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

/**
 * @deprecated use {@link mockHttpResponse}. Kept so that suites written against the axios
 * era keep compiling; the envelope it returns is the same one.
 */
export const mockAxiosResolve = <T>(data: T = {} as T): HttpResponse<T> => mockHttpResponse(data)

/**
 * @deprecated use {@link mockHttpError}, which carries a status and a body. This rejects
 * with a bare `Error`, exactly as it did before.
 */
export const mockAxiosReject = <T = never>(message = ''): Promise<T> =>
  Promise.reject(new Error(message))
