import { DavErrorCode } from './webdav'

export class HttpError extends Error {
  public response: Response
  public statusCode: number
  /** parsed response body, read once before the error is thrown */
  public data?: unknown

  constructor(
    message: string,
    response: Response,
    statusCode: number = null,
    data?: unknown
  ) {
    super(message)
    this.response = response
    this.statusCode = statusCode
    this.data = data
  }
}

export class DavHttpError extends HttpError {
  public errorCode: DavErrorCode

  constructor(
    message: string,
    errorCode: DavErrorCode,
    response: Response,
    statusCode: number = null
  ) {
    super(message, response, statusCode)
    this.errorCode = errorCode
  }
}
