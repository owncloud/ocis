/**
 * generic error implementation which captures the stackTrace and has custom naming
 */
export class RuntimeError extends Error {
  name = 'RuntimeError'
  constructor(message?: unknown, ...additional: unknown[]) {
    super([message, ...additional].filter(Boolean).join(', '))

    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, this.constructor)
    }
  }
}

/**
 * error dedicated to api type errors
 */
export class ApiError extends RuntimeError {
  name = 'ApiError'
}
