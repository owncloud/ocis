import { CapabilityStore } from '@ownclouders/web-pkg'
import { mock } from 'vitest-mock-extended'
import { AuthService } from '../../../../src/services/auth/authService'
import { createTestingPinia } from '@ownclouders/web-test-helpers'
import { useModals } from '@ownclouders/web-pkg'

// The MFA session must terminate once it actually expires, no matter how the user interacts
// with the expiry warning — dismissing it, or just leaving it open and ignoring it.
describe('AuthService MFA expiry', () => {
  afterEach(() => {
    vi.useRealTimers()
  })

  const setUp = (sessionDuration = 3600) => {
    vi.useFakeTimers()
    createTestingPinia({ stubActions: false })

    const authService = new AuthService()
    const capabilityStore = mock<CapabilityStore>({
      vaultEnabled: true,
      authMfaSessionDuration: sessionDuration
    })
    const language = { $gettext: (s: string) => s } as any

    authService.initialize(null, null, null, null, language, null, null, capabilityStore, null)

    const logoutSpy = vi.spyOn(authService, 'logoutUser')

    // Mirrors what updateMfaExpiryTimer() does: store the deadline and arm the hard logout.
    ;(authService as any).armMfaTimer(Math.floor(Date.now() / 1000) + sessionDuration)

    // Simulates the mfaExpiryWorker's onExpiring callback firing.
    ;(authService as any).showMfaExpiryWarning()

    return { authService, logoutSpy, sessionDuration }
  }

  it('terminates the session once it expires, even if the expiry warning was dismissed', () => {
    const { logoutSpy, sessionDuration } = setUp()

    const modalStore = useModals()
    const modal = modalStore.activeModal
    expect(modal.title).toBe('Session expiring')

    modal.onCancel() // user clicks "Dismiss"
    modalStore.removeModal(modal.id)

    expect(logoutSpy).not.toHaveBeenCalled()

    vi.advanceTimersByTime(sessionDuration * 1000)

    expect(logoutSpy).toHaveBeenCalled()
  })

  it('terminates the session once it expires, even if the expiry warning is left untouched', () => {
    const { logoutSpy, sessionDuration } = setUp()

    vi.advanceTimersByTime(sessionDuration * 1000)

    expect(logoutSpy).toHaveBeenCalled()
  })
})
