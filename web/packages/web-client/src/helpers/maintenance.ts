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
 * answers. Maintenance is the explicit 503 signal, not general unhealthiness, so any other real
 * response (2xx, 4xx, or a 5xx that isn't the signal) clears it. Only a transport-level failure
 * (no response at all) is left untouched, since that's not evidence either way.
 */
export function maintenanceResponseHandler(
  onSetMaintenance: (value: boolean) => void,
  { onSuccess }: MaintenanceHandlerOptions = {}
) {
  return ({ response, status, requestUrl }: OnResponseArgs): void => {
    if (!response) {
      return
    }

    if (shouldResponseTriggerMaintenance(status, requestUrl)) {
      onSetMaintenance(true)
      return
    }

    onSetMaintenance(false)
    if (response.ok) {
      onSuccess?.()
    }
  }
}
