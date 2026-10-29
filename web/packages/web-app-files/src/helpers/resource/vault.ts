import { Resource, VAULT_STORAGE_PROVIDER_ID } from '@ownclouders/web-client'

// A resource lives in the vault if its storage provider is the vault one. Space resources carry
// the full drive id in storageId, plain resources only the provider id, hence the split.
export const isVaultResource = (resource: Resource, vaultEnabled: boolean): boolean => {
  const storageId = resource?.storageId
  if (!vaultEnabled || typeof storageId !== 'string') {
    return false
  }
  return storageId.split('$')[0] === VAULT_STORAGE_PROVIDER_ID
}
