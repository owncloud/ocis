import { getAbilities } from '../../../../src/services/auth/abilities'

describe('getAbilities', () => {
  it('gets no abilities if empty permissions given', () => {
    const abilities = getAbilities([])
    expect(abilities.length).toBe(0)
  })
  it('gets correct abilities for subject "Account"', function () {
    const abilities = getAbilities(['Accounts.ReadWrite.all'])
    const expectedActions = ['create-all', 'delete-all', 'read-all', 'update-all']
    expect(abilities).toEqual(expectedActions.map((action) => ({ action, subject: 'Account' })))
  })
  it('gets correct abilities for subject "Group"', function () {
    const abilities = getAbilities(['Groups.ReadWrite.all'])
    const expectedActions = ['create-all', 'delete-all', 'read-all', 'update-all']
    expect(abilities).toEqual(expectedActions.map((action) => ({ action, subject: 'Group' })))
  })
  it.each([
    { permissions: [''], expectedActions: [] },
    { permissions: ['Favorites.List.own'], expectedActions: ['read'] },
    { permissions: ['Favorites.Write.own'], expectedActions: ['create', 'update'] }
  ])('gets correct abilities for subject "Favorites"', function (data) {
    const abilities = getAbilities(data.permissions)
    const expectedResult = data.expectedActions.map((action) => ({ action, subject: 'Favorite' }))
    expect(abilities).toEqual(expectedResult)
  })
  it('gets correct abilities for subject "Language"', function () {
    const abilities = getAbilities(['Language.ReadWrite.all'])
    const expectedActions = ['read-all', 'update-all']
    expect(abilities).toEqual(expectedActions.map((action) => ({ action, subject: 'Language' })))
  })
  it('gets correct abilities for subject "Logo"', function () {
    const abilities = getAbilities(['Logo.Write.all'])
    const expectedActions = ['update-all']
    expect(abilities).toEqual(expectedActions.map((action) => ({ action, subject: 'Logo' })))
  })
  it('gets correct abilities for subject "Role"', function () {
    const abilities = getAbilities(['Roles.ReadWrite.all'])
    const expectedActions = ['create-all', 'delete-all', 'read-all', 'update-all']
    expect(abilities).toEqual(expectedActions.map((action) => ({ action, subject: 'Role' })))
  })
  it('gets correct abilities for subject "PublicLink"', function () {
    const abilities = getAbilities(['PublicLink.Write.all'])
    const expectedActions = ['create-all']
    expect(abilities).toEqual(expectedActions.map((action) => ({ action, subject: 'PublicLink' })))
  })
  it('gets correct abilities for subject "Share"', function () {
    const abilities = getAbilities(['Shares.Write.all'])
    const expectedActions = ['create-all', 'update-all']
    expect(abilities).toEqual(expectedActions.map((action) => ({ action, subject: 'Share' })))
  })
  it('gets correct abilities for subject "Setting"', function () {
    const abilities = getAbilities(['Settings.ReadWrite.all'])
    const expectedActions = ['read-all', 'update-all']
    expect(abilities).toEqual(expectedActions.map((action) => ({ action, subject: 'Setting' })))
  })
  it.each([
    { permissions: [''], expectedActions: [] },
    { permissions: ['Drives.Create.all'], expectedActions: ['create-all'] },
    { permissions: ['Drives.List.all'], expectedActions: ['read-all'] },
    { permissions: ['Drives.ReadWriteEnabled.all'], expectedActions: ['delete-all', 'update-all'] },
    { permissions: ['Drives.ReadWriteProjectQuota.all'], expectedActions: ['set-quota-all'] }
  ])('gets correct abilities for subject "Drive"', function (data) {
    const abilities = getAbilities(data.permissions)
    const expectedResult = data.expectedActions.map((action) => ({ action, subject: 'Drive' }))
    expect(abilities).toEqual(expectedResult)
  })
})
