import { navItems, routes } from '../../src/index'
import { Ability } from '@ownclouders/web-client'
import { createTestingPinia } from '@ownclouders/web-test-helpers'
import { mock } from 'vitest-mock-extended'
import * as pkg from '@ownclouders/web-pkg'
import { AuthService } from '../../../web-runtime/src/services/auth/authService'

const mockRequireAcr = vi.fn()
vi.spyOn(pkg, 'useAuthService').mockReturnValue(mock<AuthService>({ requireAcr: mockRequireAcr }))

const getAbilityMock = (hasPermission: boolean) => mock<Ability>({ can: () => hasPermission })

describe('admin settings index', () => {
  beforeEach(() => {
    createTestingPinia({
      initialState: {
        capabilities: {
          capabilities: { auth: { mfa: { enabled: true, levelnames: ['advanced'] } } }
        }
      }
    })
  })

  describe('navItems', () => {
    describe('general', () => {
      it.each([true, false])('should be enabled according to the permissions', (enabled) => {
        expect(
          navItems({ $ability: getAbilityMock(enabled) })
            .find((n) => n.name === 'General')
            .isVisible()
        ).toBe(enabled)
      })
    })
    describe('user management', () => {
      it.each([true, false])('should be enabled according to the permissions', (enabled) => {
        expect(
          navItems({ $ability: getAbilityMock(enabled) })
            .find((n) => n.name === 'Users')
            .isVisible()
        ).toBe(enabled)
      })
    })
    describe('group management', () => {
      it.each([true, false])('should be enabled according to the permissions', (enabled) => {
        expect(
          navItems({ $ability: getAbilityMock(enabled) })
            .find((n) => n.name === 'Groups')
            .isVisible()
        ).toBe(enabled)
      })
    })
    describe('space management', () => {
      it.each([true, false])('should be enabled according to the permissions', (enabled) => {
        expect(
          navItems({ $ability: getAbilityMock(enabled) })
            .find((n) => n.name === 'Spaces')
            .isVisible()
        ).toBe(enabled)
      })
    })
  })
  describe('routes', () => {
    describe('default-route "/"', () => {
      it('should redirect to general', () => {
        const ability = mock<Ability>()
        ability.can.mockReturnValueOnce(true)
        const route = routes({ $ability: ability }).find((n) => n.path === '/')
        expect((route.redirect as any)().name).toEqual('admin-settings-general')
      })
    })
    it.each([
      { can: vi.fn(() => true), redirect: null },
      {
        can: vi.fn((_, subject) => {
          if (subject === 'Group') {
            return true
          }

          return false
        }),
        redirect: { name: 'admin-settings-groups' }
      }
    ])('redirects "/general" with sufficient permissions', async ({ can, redirect }) => {
      const ability = mock<Ability>({ can })
      const route = routes({ $ability: ability }).find((n) => n.path === '/general')
      const nextMock = vi.fn()
      await (route.beforeEnter as any)({}, {}, nextMock)
      const args = [...(redirect ? [redirect] : [])]
      expect(nextMock).toHaveBeenCalledWith(...args)
    })
    it.each([
      { can: vi.fn(() => true), redirect: null },
      {
        can: vi.fn((_, subject) => {
          if (subject === 'Drive') {
            return true
          }

          return false
        }),
        redirect: { name: 'admin-settings-spaces' }
      }
    ])('redirects "/users" with sufficient permissions', async ({ can, redirect }) => {
      const ability = mock<Ability>({ can })
      const route = routes({ $ability: ability }).find((n) => n.path === '/users')
      const nextMock = vi.fn()
      await (route.beforeEnter as any)({}, {}, nextMock)
      const args = [...(redirect ? [redirect] : [])]
      expect(nextMock).toHaveBeenCalledWith(...args)
    })
    it.each([
      { can: vi.fn(() => true), redirect: null },
      {
        can: vi.fn((_, subject) => {
          if (subject === 'Setting') {
            return true
          }

          return false
        }),
        redirect: { name: 'admin-settings-general' }
      }
    ])('redirects "/groups" with sufficient permissions', async ({ can, redirect }) => {
      const ability = mock<Ability>({ can })
      const route = routes({ $ability: ability }).find((n) => n.path === '/groups')
      const nextMock = vi.fn()
      await (route.beforeEnter as any)({}, {}, nextMock)
      const args = [...(redirect ? [redirect] : [])]
      expect(nextMock).toHaveBeenCalledWith(...args)
    })
    it.each([
      { can: vi.fn(() => true), redirect: null },
      {
        can: vi.fn((_, subject) => {
          if (subject === 'Account') {
            return true
          }

          return false
        }),
        redirect: { name: 'admin-settings-users' }
      }
    ])('redirects "/spaces" with sufficient permissions', async ({ can, redirect }) => {
      const ability = mock<Ability>({ can })
      const route = routes({ $ability: ability }).find((n) => n.path === '/spaces')
      const nextMock = vi.fn()
      await (route.beforeEnter as any)({}, {}, nextMock)
      const args = [...(redirect ? [redirect] : [])]
      expect(nextMock).toHaveBeenCalledWith(...args)
    })
    it.each(['/general', '/users', '/groups', '/spaces'])(
      'should throw an error if permissions are insufficient',
      async (path) => {
        const ability = mock<Ability>({ can: vi.fn(() => false) })
        const route = routes({ $ability: ability }).find((n) => n.path === path)
        const nextMock = vi.fn()
        await expect(() => (route.beforeEnter as any)({}, {}, nextMock)).rejects.toThrowError(
          'Insufficient permissions'
        )
      }
    )

    describe('requireAcr', () => {
      it.each(['/general', '/users', '/groups', '/spaces'])(
        'should call requireAcr if MFA is enabled when path is %s',
        async (path) => {
          const ability = mock<Ability>({ can: vi.fn(() => true) })
          const route = routes({ $ability: ability }).find((n) => n.path === path)
          await (route.beforeEnter as any)({ fullPath: path }, {}, vi.fn())
          expect(mockRequireAcr).toHaveBeenCalledWith('advanced', path)
        }
      )
    })

    describe('requireAcr', () => {
      it.each(['/general', '/users', '/groups', '/spaces'])(
        'should not call requireAcr if MFA is disabled when path is %s',
        async (path) => {
          createTestingPinia({
            initialState: {
              capabilities: {
                capabilities: { auth: { mfa: { enabled: false, levelnames: ['advanced'] } } }
              }
            }
          })

          const ability = mock<Ability>({ can: vi.fn(() => true) })
          const route = routes({ $ability: ability }).find((n) => n.path === path)
          await (route.beforeEnter as any)({}, {}, vi.fn())
          expect(mockRequireAcr).not.toHaveBeenCalled()
        }
      )
    })
  })
})
