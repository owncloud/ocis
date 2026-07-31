/**
 * This composable is used for various vault-related functionality.
 */
interface VaultComposable {
  /**
   * Checks whether the user is currently in the vault.
   * This is not reactive value because a full page reload is required to switch between regular and vault mode.
   */
  isInVault: boolean
}

/**
 * Checks whether the given URL is a vault URL.
 * @param url - The URL to check.
 * @returns Whether the URL is a vault URL.
 */
function getIsVaultUrl(url: string): boolean {
  if (url.startsWith('/vault')) {
    return true
  }

  const hashIndex = url.indexOf('#')
  if (hashIndex !== -1) {
    const hash = url.slice(hashIndex)
    return /^#\/vault(?:\/|$)/.test(hash)
  }

  return false
}

export function useVault(): VaultComposable {
  const isInVault = (() => {
    const { pathname, hash } = window.location

    if (getIsVaultUrl(pathname + hash)) {
      return true
    }

    // During OIDC callback the current URL is /web-oidc-callback, not the vault path.
    // The intended post-login destination is stored in sessionStorage/localStorage under
    // this key by the auth service (web-runtime/src/services/auth/userManager.ts).
    if (pathname === '/web-oidc-callback') {
      const postLoginRedirectUrlKey = 'oc.postLoginRedirectUrl'
      const postLoginRedirectUrl =
        sessionStorage.getItem(postLoginRedirectUrlKey) ||
        localStorage.getItem(postLoginRedirectUrlKey)

      if (postLoginRedirectUrl) {
        return getIsVaultUrl(postLoginRedirectUrl)
      }
    }

    return false
  })()

  return {
    isInVault
  }
}
