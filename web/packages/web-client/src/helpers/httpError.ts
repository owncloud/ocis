import { HttpError } from '../errors'

/**
 * Create a HttpError based on a given message, status code and x-request-id.
 */
export const createHttpError = ({
  message,
  statusCode,
  xReqId
}: {
  message: string
  statusCode: number
  xReqId: string
}) => {
  const response = new Response(undefined, {
    headers: { 'x-request-id': xReqId },
    status: statusCode
  })
  return new HttpError(message, response, statusCode)
}
