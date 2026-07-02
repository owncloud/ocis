import { PasswordPolicyService } from '../../../src/services'
import { createTestingPinia } from '@ownclouders/web-test-helpers'
import { Language } from 'vue3-gettext'
import { PasswordPolicyCapability } from '@ownclouders/web-client/ocs'
import { useCapabilityStore } from '../../../src/composables/piniaStores'
import { describe } from 'vitest'

describe('PasswordPolicyService', () => {
  describe('policy', () => {
    describe('mustNotBeEmpty rule', () => {
      it('is present when "enforcePassword" is set', () => {
        {
          const { passwordPolicyService, store } = getWrapper({} as PasswordPolicyCapability)
          passwordPolicyService.initialize(store)
          expect(
            Object.keys(passwordPolicyService.getPolicy({ enforcePassword: true }).rules)
          ).toEqual(['mustNotBeEmpty'])
        }
      })
      it('is not present when "enforcePassword" is set but other rules apply', () => {
        {
          const { passwordPolicyService, store } = getWrapper({
            min_characters: 2
          } as PasswordPolicyCapability)
          passwordPolicyService.initialize(store)
          expect(
            Object.keys(passwordPolicyService.getPolicy({ enforcePassword: true }).rules)
          ).toEqual(['atLeastCharacters'])
        }
      })
      it('is not present when "enforcePassword" is not set', () => {
        {
          const { passwordPolicyService, store } = getWrapper({} as PasswordPolicyCapability)
          passwordPolicyService.initialize(store)
          expect(Object.keys(passwordPolicyService.getPolicy({}).rules)).toEqual([])
        }
      })
    })
    describe('contains the rules according to the capability', () => {
      it.each([
        [{ min_characters: 2 } as PasswordPolicyCapability, ['atLeastCharacters']],
        [
          { min_lowercase_characters: 2 } as PasswordPolicyCapability,
          ['atLeastLowercaseCharacters']
        ],
        [
          { min_uppercase_characters: 2 } as PasswordPolicyCapability,
          ['atLeastUppercaseCharacters']
        ],
        [{ min_digits: 2 } as PasswordPolicyCapability, ['atLeastDigits']],
        [{ min_digits: 2 } as PasswordPolicyCapability, ['atLeastDigits']],
        [{ min_special_characters: 2 } as PasswordPolicyCapability, ['mustContain']],
        [
          {
            min_characters: 2,
            min_lowercase_characters: 2,
            min_uppercase_characters: 2,
            min_digits: 2,
            min_special_characters: 2
          } as PasswordPolicyCapability,
          [
            'atLeastCharacters',
            'atLeastUppercaseCharacters',
            'atLeastLowercaseCharacters',
            'atLeastDigits',
            'mustContain'
          ]
        ]
      ])('capability "%s"', (capability: PasswordPolicyCapability, expected: Array<string>) => {
        const { passwordPolicyService, store } = getWrapper(capability)
        passwordPolicyService.initialize(store)
        expect(Object.keys(passwordPolicyService.getPolicy().rules)).toEqual(expected)
      })
    })
    describe('method "check"', () => {
      describe('test the password correctly against te defined rules', () => {
        it.each([
          [{} as PasswordPolicyCapability, ['', 'o'], true, [false, true]],
          [{} as PasswordPolicyCapability, ['', 'o'], false, [true, true]],
          [
            { min_characters: 2 } as PasswordPolicyCapability,
            ['', 'o', 'ow', 'ownCloud'],
            false,
            [false, false, true, true]
          ],
          [
            { min_lowercase_characters: 2 } as PasswordPolicyCapability,
            ['', 'o', 'oWNCLOUD', 'ownCloud'],
            false,
            [false, false, false, true]
          ],
          [
            { min_uppercase_characters: 2 } as PasswordPolicyCapability,
            ['', 'o', 'ownCloud', 'ownCLoud'],
            false,
            [false, false, false, true]
          ],
          [
            { min_digits: 2 } as PasswordPolicyCapability,
            ['', '1', 'ownCloud1', 'ownCloud12'],
            false,
            [false, false, false, true]
          ],
          [
            { min_special_characters: 2 } as PasswordPolicyCapability,
            ['', '!', 'ownCloud!', 'ownCloud!#'],
            false,
            [false, false, false, true]
          ],
          [
            {
              min_characters: 8,
              min_lowercase_characters: 2,
              min_uppercase_characters: 2,
              min_digits: 2,
              min_special_characters: 2,
              max_characters: 72
            } as PasswordPolicyCapability,
            ['öwnCloud', 'öwnCloudää', 'öwnCloudää12', 'öwnCloudäÄ12#!'],
            false,
            [false, false, false, true]
          ]
        ])(
          'capability "%s, passwords "%s"',
          (
            capability: PasswordPolicyCapability,
            passwords: Array<string>,
            enforcePassword: boolean,
            expected: Array<boolean>
          ) => {
            const { passwordPolicyService, store } = getWrapper(capability)
            passwordPolicyService.initialize(store)
            const policy = passwordPolicyService.getPolicy({ enforcePassword })
            for (let i = 0; i < passwords.length; i++) {
              expect(policy.check(passwords[i])).toEqual(expected[i])
            }
          }
        )
      })
    })
  })

  describe('password generator', () => {
    it('should generate a password with default rules when no capability is set', () => {
      const { passwordPolicyService, store } = getWrapper({})
      passwordPolicyService.initialize(store)

      const password = passwordPolicyService.generatePassword()

      expect(password.length).toBe(12)
      expect((password.match(/[a-z]/g) || []).length).toBeGreaterThanOrEqual(2)
      expect((password.match(/[A-Z]/g) || []).length).toBeGreaterThanOrEqual(2)
      expect((password.match(/[0-9]/g) || []).length).toBeGreaterThanOrEqual(2)

      const specialCharsRegex = /[!#$%&'()*+,\-./:;<=>?@[\\\]^_`{|}~]/g
      expect((password.match(specialCharsRegex) || []).length).toBeGreaterThanOrEqual(2)
    })

    it('should generate a password with the specified minimum character counts', () => {
      const capability = {
        min_characters: 16,
        min_lowercase_characters: 3,
        min_uppercase_characters: 4,
        min_digits: 5,
        min_special_characters: 2
      }

      const { passwordPolicyService, store } = getWrapper(capability)
      passwordPolicyService.initialize(store)

      const password = passwordPolicyService.generatePassword()

      expect(password.length).toBe(16)
      expect((password.match(/[a-z]/g) || []).length).toBeGreaterThanOrEqual(3)
      expect((password.match(/[A-Z]/g) || []).length).toBeGreaterThanOrEqual(4)
      expect((password.match(/[0-9]/g) || []).length).toBeGreaterThanOrEqual(5)

      const specialCharsRegex = /[!#$%&'()*+,\-./:;<=>?@[\\\]^_`{|}~]/g
      expect((password.match(specialCharsRegex) || []).length).toBeGreaterThanOrEqual(2)
    })

    it('should generate a password with length based on sum of minimum requirements if greater than min_characters', () => {
      const capability = {
        min_characters: 10,
        min_lowercase_characters: 3,
        min_uppercase_characters: 4,
        min_digits: 5,
        min_special_characters: 2
      }

      const { passwordPolicyService, store } = getWrapper(capability)
      passwordPolicyService.initialize(store)

      const password = passwordPolicyService.generatePassword()

      expect(password.length).toBe(14)
    })

    it('should generate different passwords on subsequent calls', () => {
      const { passwordPolicyService, store } = getWrapper({} as PasswordPolicyCapability)
      passwordPolicyService.initialize(store)

      const password1 = passwordPolicyService.generatePassword()
      const password2 = passwordPolicyService.generatePassword()

      expect(password1).not.toBe(password2)
    })
  })
})

const getWrapper = (capability: PasswordPolicyCapability) => {
  createTestingPinia()
  const store = useCapabilityStore()
  store.capabilities.password_policy = capability

  return {
    store,
    passwordPolicyService: new PasswordPolicyService({
      language: { current: 'en' } as Language
    })
  }
}
