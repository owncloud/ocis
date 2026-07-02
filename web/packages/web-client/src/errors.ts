import { DavErrorCode } from './webdav'

export class HttpError extends Error {
  public response: Response
  public statusCode: number

  constructor(message: string, response: Response, statusCode: number = null) {
    super(message)
    this.response = response
    this.statusCode = statusCode
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
