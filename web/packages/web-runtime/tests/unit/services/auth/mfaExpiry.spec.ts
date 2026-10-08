import { CapabilityStore } from '@ownclouders/web-pkg'
import { mock } from 'vitest-mock-extended'
import { AuthService } from '../../../../src/services/auth/authService'
import { createTestingPinia } from '@ownclouders/web-test-helpers'
import { useModals } from '@ownclouders/web-pkg'

// Dismissing the MFA expiry warning must not keep the session alive forever.
describe('AuthService MFA expiry', () => {
  afterEach(() => {
    vi.useRealTimers()
  })

  it('terminates the session once the MFA session actually expires, even if the expiry warning was dismissed', () => {
    vi.useFakeTimers()
    createTestingPinia({ stubActions: false })

    const authService = new AuthService()
    const sessionDuration = 3600
    const capabilityStore = mock<CapabilityStore>({
      vaultEnabled: true,
      authMfaSessionDuration: sessionDuration
    })
    const language = { $gettext: (s: string) => s } as any

    authService.initialize(null, null, null, null, language, null, null, capabilityStore, null)

    const logoutSpy = vi.spyOn(authService, 'logoutUser')

    // Mirrors what updateMfaExpiryTimer() would have stored before arming the worker.
    ;(authService as any).mfaExpiresAt = Math.floor(Date.now() / 1000) + sessionDuration

    // Simulates the mfaExpiryWorker's onExpiring callback firing.
    ;(authService as any).showMfaExpiryWarning()

    const modalStore = useModals()
    const modal = modalStore.activeModal
    expect(modal.title).toBe('Session expiring')

    modal.onCancel() // user clicks "Dismiss"
    modalStore.removeModal(modal.id)

    expect(logoutSpy).not.toHaveBeenCalled()

    vi.advanceTimersByTime(sessionDuration * 1000)

    expect(logoutSpy).toHaveBeenCalled()
  })
})
