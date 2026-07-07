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
