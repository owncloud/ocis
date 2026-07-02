import { useClientService } from '@ownclouders/web-pkg'

export function useMaintenanceMode() {
  let timeoutId = null

  const clientService = useClientService()

  /**
   * Starts a timer that checks for maintenance mode every minute.
   * Since the maintenance mode is asserted by a request that returns a 503 status code, we can just call any endpoint.
   * Response is parsed in the axios response interceptor.
   */
  const startCheckingMaintenanceMode = async () => {
    try {
      await clientService.ocs.getCapabilities()
    } catch (error) {
      console.error(error)
    }

    timeoutId = setTimeout(startCheckingMaintenanceMode, 60_000)
  }

  const stopCheckingMaintenanceMode = () => {
    if (!timeoutId) {
      return
    }

    clearTimeout(timeoutId)
    timeoutId = null
  }

  return {
    startCheckingMaintenanceMode,
    stopCheckingMaintenanceMode
  }
}
