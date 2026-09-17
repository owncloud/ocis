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
  /** extra bookkeeping on success, such as recording the last successful request time */
  onSuccess?: () => void
}

/**
 * Builds the `onResponse` handler that keeps maintenance mode in sync with what the server
 * answers. Only a successful response clears maintenance mode; an unrelated error (a 404, a
 * transient 500, ...) isn't evidence the server is healthy again, so it leaves the flag as is.
 */
export function maintenanceResponseHandler(
  onSetMaintenance: (value: boolean) => void,
  { onSuccess }: MaintenanceHandlerOptions = {}
) {
  return ({ response, status, requestUrl }: OnResponseArgs): void => {
    if (response?.ok) {
      onSetMaintenance(false)
      onSuccess?.()
      return
    }

    if (shouldResponseTriggerMaintenance(status, requestUrl)) {
      onSetMaintenance(true)
    }
  }
}
