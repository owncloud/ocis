import type { OnResponseArgs } from '../http'

/**
 * List of all API endpoints that should not trigger a maintenance mode warning even when they return 503 status code.
 */
const MAINTENANCE_EXCLUDED_ENDPOINTS = ['ocs/v2.php/apps/notifications/api/v1/notifications/sse']

export function shouldResponseTriggerMaintenance(responseStatus: number, requestUrl: string) {
  if (responseStatus === 503 && !MAINTENANCE_EXCLUDED_ENDPOINTS.includes(requestUrl)) {
    return true
  }

  return false
}

export interface MaintenanceHandlerOptions {
  /**
   * Whether a non-2xx response that is *not* a maintenance signal clears maintenance mode.
   *
   * The two clients have always disagreed here, and both behaviours are kept as they were.
   * The webdav client clears (its axios error interceptor called
   * `onSetMaintenance(shouldResponseTriggerMaintenance(...))` unconditionally, so any non-503
   * error reset the flag); `ClientService` does not (its `#handleAxiosError` only ever set the
   * flag to `true`). The inconsistency predates the fetch migration — unifying it would change
   * when the maintenance banner disappears, which is a product decision, not a refactor.
   */
  clearOnUnrelatedError?: boolean
  /** extra bookkeeping on success, such as recording the last successful request time */
  onSuccess?: () => void
}

/**
 * Builds the `onResponse` handler that keeps maintenance mode in sync with what the server
 * answers. Only a successful response clears maintenance mode unconditionally.
 */
export function maintenanceResponseHandler(
  onSetMaintenance: (value: boolean) => void,
  { clearOnUnrelatedError = false, onSuccess }: MaintenanceHandlerOptions = {}
) {
  return ({ response, status, requestUrl }: OnResponseArgs): void => {
    if (response?.ok) {
      onSetMaintenance(false)
      onSuccess?.()
      return
    }

    const isMaintenance = shouldResponseTriggerMaintenance(status, requestUrl)
    if (isMaintenance || clearOnUnrelatedError) {
      onSetMaintenance(isMaintenance)
    }
  }
}
