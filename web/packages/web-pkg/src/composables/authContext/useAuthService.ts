import { useService } from '../service'
import { NavigationFailure } from 'vue-router'

export interface AuthServiceInterface {
  handleAuthError(route: any, options?: { forceLogout?: boolean }): any
  signinSilent(): Promise<unknown>
  logoutUser(): Promise<void | NavigationFailure>
  getRefreshToken(): Promise<string>
  /**
   * @returns false if the required `acr` could not be reached, access must not be granted then
   */
  requireAcr(acrValue: string, redirectUrl: string): Promise<boolean>
}

export const useAuthService = (): AuthServiceInterface => {
  return useService('$authService')
}
