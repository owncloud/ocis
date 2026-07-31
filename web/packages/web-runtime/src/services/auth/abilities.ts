import { SubjectRawRule } from '@casl/ability'
import { AbilityActions, AbilitySubjects } from '@ownclouders/web-client'

export const getAbilities = (
  permissions: string[]
): SubjectRawRule<AbilityActions, AbilitySubjects, any>[] => {
  const abilities: Record<string, SubjectRawRule<AbilityActions, AbilitySubjects, any>[]> = {
    'Accounts.ReadWrite.all': [
      { action: 'create-all', subject: 'Account' },
      { action: 'delete-all', subject: 'Account' },
      { action: 'read-all', subject: 'Account' },
      { action: 'update-all', subject: 'Account' }
    ],
    'Favorites.List.own': [{ action: 'read', subject: 'Favorite' }],
    'Favorites.Write.own': [
      { action: 'create', subject: 'Favorite' },
      { action: 'update', subject: 'Favorite' }
    ],
    'Groups.ReadWrite.all': [
      { action: 'create-all', subject: 'Group' },
      { action: 'delete-all', subject: 'Group' },
      { action: 'read-all', subject: 'Group' },
      { action: 'update-all', subject: 'Group' }
    ],
    'Language.ReadWrite.all': [
      { action: 'read-all', subject: 'Language' },
      { action: 'update-all', subject: 'Language' }
    ],
    'Logo.Write.all': [{ action: 'update-all', subject: 'Logo' }],
    'PublicLink.Write.all': [{ action: 'create-all', subject: 'PublicLink' }],
    'ReadOnlyPublicLinkPassword.Delete.all': [
      { action: 'delete-all', subject: 'ReadOnlyPublicLinkPassword' }
    ],
    'Roles.ReadWrite.all': [
      { action: 'create-all', subject: 'Role' },
      { action: 'delete-all', subject: 'Role' },
      { action: 'read-all', subject: 'Role' },
      { action: 'update-all', subject: 'Role' }
    ],
    'Shares.Write.all': [
      { action: 'create-all', subject: 'Share' },
      { action: 'update-all', subject: 'Share' }
    ],
    'Settings.ReadWrite.all': [
      { action: 'read-all', subject: 'Setting' },
      { action: 'update-all', subject: 'Setting' }
    ],
    'Drives.Create.all': [{ action: 'create-all', subject: 'Drive' }],
    'Drives.ReadWriteEnabled.all': [
      { action: 'delete-all', subject: 'Drive' },
      { action: 'update-all', subject: 'Drive' }
    ],
    'Drives.List.all': [{ action: 'read-all', subject: 'Drive' }],
    'Drives.ReadWriteProjectQuota.all': [{ action: 'set-quota-all', subject: 'Drive' }],
    'VaultMode.ReadWriteEnabled.own': [{ action: 'read-all', subject: 'Vault' }]
  }

  return Object.keys(abilities).reduce((acc, permission) => {
    if (permissions.includes(permission)) {
      acc.push(...abilities[permission])
    }
    return acc
  }, [])
}
