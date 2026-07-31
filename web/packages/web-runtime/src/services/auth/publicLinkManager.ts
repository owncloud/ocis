import { AuthStore, CapabilityStore, ClientService } from '@ownclouders/web-pkg'
import { PublicLinkType } from '@ownclouders/web-client'

export interface PublicLinkManagerOptions {
  clientService: ClientService
  authStore: AuthStore
  capabilityStore: CapabilityStore
}

export class PublicLinkManager {
  private clientService: ClientService
  private authStore: AuthStore
  private capabilityStore: CapabilityStore

  constructor(options: PublicLinkManagerOptions) {
    this.clientService = options.clientService
    this.authStore = options.authStore
    this.capabilityStore = options.capabilityStore
  }

  private static buildStorageKey(token: string, suffix: string): string {
    return `oc.publicLink.${token}.${suffix}`
  }

  clear(token: string) {
    ;['resolved', 'passwordRequired', 'password'].forEach((key) => {
      sessionStorage.removeItem(PublicLinkManager.buildStorageKey(token, key))
    })
    this.authStore.clearPublicLinkContext()
  }

  isResolved(token: string): boolean {
    const resolved = sessionStorage.getItem(PublicLinkManager.buildStorageKey(token, 'resolved'))
    return resolved === 'true'
  }

  setResolved(token: string, resolved: boolean): void {
    sessionStorage.setItem(PublicLinkManager.buildStorageKey(token, 'resolved'), resolved + '')
  }

  setType(token: string, type: PublicLinkType): void {
    sessionStorage.setItem(PublicLinkManager.buildStorageKey(token, 'type'), type)
  }

  getType(token: string): PublicLinkType {
    return sessionStorage.getItem(
      PublicLinkManager.buildStorageKey(token, 'type')
    ) as PublicLinkType
  }

  isPasswordRequired(token: string): boolean {
    const passwordRequired = sessionStorage.getItem(
      PublicLinkManager.buildStorageKey(token, 'passwordRequired')
    )
    return passwordRequired === 'true'
  }

  setPasswordRequired(token: string, required: boolean): void {
    sessionStorage.setItem(
      PublicLinkManager.buildStorageKey(token, 'passwordRequired'),
      required + ''
    )
  }

  getPassword(token: string): string {
    const password = sessionStorage.getItem(PublicLinkManager.buildStorageKey(token, 'password'))
    if (password) {
      try {
        return Buffer.from(password, 'base64').toString()
      } catch {
        this.clear(token)
      }
    }
    return ''
  }

  setPassword(token: string, password: string): void {
    if (password.length) {
      const encodedPassword = Buffer.from(password).toString('base64')
      sessionStorage.setItem(PublicLinkManager.buildStorageKey(token, 'password'), encodedPassword)
    } else {
      sessionStorage.removeItem(PublicLinkManager.buildStorageKey(token, 'password'))
    }
  }

  async updateContext(token: string) {
    if (!this.isResolved(token)) {
      return
    }
    if (this.authStore.publicLinkContextReady && this.authStore.publicLinkToken === token) {
      return
    }

    let password
    if (this.isPasswordRequired(token)) {
      password = this.getPassword(token)
    }

    try {
      await this.fetchCapabilities()
    } catch (e) {
      console.error(e)
    }

    this.authStore.setPublicLinkContext({
      publicLinkToken: token,
      publicLinkPassword: password,
      publicLinkContextReady: true,
      publicLinkType: this.getType(token),
      publicLinkPasswordRequired: this.isPasswordRequired(token)
    })
  }

  clearContext() {
    this.authStore.clearPublicLinkContext()
  }

  private async fetchCapabilities(): Promise<void> {
    if (this.capabilityStore.isInitialized) {
      return
    }
    const client = this.clientService.ocs
    const response = await client.getCapabilities()
    this.capabilityStore.setCapabilities(response)
  }
}
