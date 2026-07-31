import { useService } from '../service'
import { NavigationFailure } from 'vue-router'

export interface AuthServiceInterface {
  handleAuthError(route: any, options?: { forceLogout?: boolean }): any
  signinSilent(): Promise<unknown>
  logoutUser(): Promise<void | NavigationFailure>
  getRefreshToken(): Promise<string>
  requireAcr(acrValue: string, redirectUrl: string): Promise<void>
}

export const useAuthService = (): AuthServiceInterface => {
  return useService('$authService')
}
