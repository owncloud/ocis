// Marks a pending vault MFA step-up so the redirect back from the IdP can be recognized.
// sessionStorage: survives the IdP round trip, scoped to the tab.
export const vaultStepUpAttemptKey = 'oc_vaultMfaStepUpAttempt'
// Set when a step-up failed and the page is being reloaded outside the vault,
// so the message can be shown after the reload.
export const vaultStepUpFailedKey = 'oc_vaultMfaStepUpFailed'

/**
 * Whether the page is currently being reloaded outside the vault after a failed MFA step-up.
 * Requests with the non-MFA token (e.g. loading vault spaces) are expected to fail meanwhile.
 */
export const isLeavingVaultAfterFailedStepUp = (): boolean => {
  try {
    return sessionStorage.getItem(vaultStepUpFailedKey) !== null
  } catch {
    return false
  }
}
