import { Language } from 'vue3-gettext'
import {
  AtLeastCharactersRule,
  AtLeastDigitsRule,
  AtLeastLowercaseCharactersRule,
  AtLeastUppercaseCharactersRule,
  MustContainRule,
  MustNotBeEmptyRule
} from './rules'
import { PasswordPolicyCapability } from '@ownclouders/web-client/ocs'
import { CapabilityStore } from '../../composables'

// @ts-ignore
import { PasswordPolicy } from 'password-sheriff'

interface GeneratePasswordRules {
  length: number
  minLowercaseCharacters: number
  minUppercaseCharacters: number
  minSpecialCharacters: number
  minDigits: number
}

export class PasswordPolicyService {
  private readonly language: Language
  private capability: PasswordPolicyCapability
  private policy: PasswordPolicy
  private generatePasswordRules: GeneratePasswordRules

  constructor({ language }: { language: Language }) {
    this.language = language
  }

  public initialize(capabilityStore: CapabilityStore) {
    this.capability = capabilityStore.passwordPolicy
    this.buildGeneratePasswordRules()
  }

  private hasRules(): boolean {
    return (
      !!this.capability.min_characters ||
      !!this.capability.min_lowercase_characters ||
      !!this.capability.min_uppercase_characters ||
      !!this.capability.min_digits ||
      !!this.capability.min_special_characters
    )
  }

  private buildGeneratePasswordRules(): void {
    const DEFAULT_LENGTH = 12
    const DEFAULT_MIN_LOWERCASE_CHARACTERS = 2
    const DEFAULT_MIN_UPPERCASE_CHARACTERS = 2
    const DEFAULT_MIN_SPECIAL_CHARACTERS = 2
    const DEFAULT_MIN_DIGITS = 2

    this.generatePasswordRules = {
      length: Math.max(
        this.capability.min_characters || 0,
        (this.capability.min_lowercase_characters || 0) +
          (this.capability.min_uppercase_characters || 0) +
          (this.capability.min_digits || 0) +
          (this.capability.min_special_characters || 0),
        DEFAULT_LENGTH
      ),
      minLowercaseCharacters: Math.max(
        this.capability.min_lowercase_characters || 0,
        DEFAULT_MIN_LOWERCASE_CHARACTERS
      ),
      minUppercaseCharacters: Math.max(
        this.capability.min_uppercase_characters || 0,
        DEFAULT_MIN_UPPERCASE_CHARACTERS
      ),
      minSpecialCharacters: Math.max(
        this.capability.min_special_characters || 0,
        DEFAULT_MIN_SPECIAL_CHARACTERS
      ),
      minDigits: Math.max(this.capability.min_digits || 0, DEFAULT_MIN_DIGITS)
    }
  }

  private buildPolicy({ enforcePassword = false } = {}): void {
    const ruleset = {
      atLeastCharacters: new AtLeastCharactersRule({ ...this.language }),
      mustNotBeEmpty: new MustNotBeEmptyRule({ ...this.language }),
      atLeastUppercaseCharacters: new AtLeastUppercaseCharactersRule({ ...this.language }),
      atLeastLowercaseCharacters: new AtLeastLowercaseCharactersRule({ ...this.language }),
      atLeastDigits: new AtLeastDigitsRule({ ...this.language }),
      mustContain: new MustContainRule({ ...this.language })
    }
    const rules = {} as Record<string, unknown>

    if (enforcePassword && !this.hasRules()) {
      rules.mustNotBeEmpty = {}
    }

    if (this.capability.min_characters) {
      rules.atLeastCharacters = { minLength: this.capability.min_characters }
    }

    if (this.capability.min_uppercase_characters) {
      rules.atLeastUppercaseCharacters = {
        minLength: this.capability.min_uppercase_characters
      }
    }

    if (this.capability.min_lowercase_characters) {
      rules.atLeastLowercaseCharacters = {
        minLength: this.capability.min_lowercase_characters
      }
    }

    if (this.capability.min_digits) {
      rules.atLeastDigits = { minLength: this.capability.min_digits }
    }

    if (this.capability.min_special_characters) {
      rules.mustContain = {
        minLength: this.capability.min_special_characters,
        characters: ' "!#\\$%&\'()*+,-./:;<=>?@[\\]^_`{|}~"'
      }
    }

    this.policy = new PasswordPolicy(rules, ruleset)
  }

  public getPolicy({ enforcePassword = false } = {}): PasswordPolicy {
    this.buildPolicy({ enforcePassword })
    return this.policy
  }

  public generatePassword(): string {
    const lowerChars = 'abcdefghijklmnopqrstuvwxyz'
    const upperChars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
    const numberChars = '0123456789'
    const symbolChars = "!#\\$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

    const totalMinChars =
      this.generatePasswordRules.minUppercaseCharacters +
      this.generatePasswordRules.minLowercaseCharacters +
      this.generatePasswordRules.minDigits +
      this.generatePasswordRules.minSpecialCharacters

    if (totalMinChars > this.generatePasswordRules.length) {
      throw new Error('Sum of minimum character requirements exceeds password length')
    }

    let passwdArray: string[] = []

    const getRandomCharsFromSet = (charSet: string, count: number): string[] => {
      const setLimit = 256 - (256 % charSet.length)
      const result: string[] = []

      for (let i = 0; i < count; i++) {
        let randval: number

        do {
          randval = window.crypto.getRandomValues(new Uint8Array(1))[0]
        } while (randval >= setLimit)

        result.push(charSet[randval % charSet.length])
      }

      return result
    }

    passwdArray = passwdArray.concat(
      getRandomCharsFromSet(lowerChars, this.generatePasswordRules.minLowercaseCharacters)
    )
    passwdArray = passwdArray.concat(
      getRandomCharsFromSet(upperChars, this.generatePasswordRules.minUppercaseCharacters)
    )
    passwdArray = passwdArray.concat(
      getRandomCharsFromSet(numberChars, this.generatePasswordRules.minDigits)
    )
    passwdArray = passwdArray.concat(
      getRandomCharsFromSet(symbolChars, this.generatePasswordRules.minSpecialCharacters)
    )

    const allChars = lowerChars + upperChars + numberChars + symbolChars
    const remaining = this.generatePasswordRules.length - passwdArray.length

    passwdArray = passwdArray.concat(getRandomCharsFromSet(allChars, remaining))

    for (let i = passwdArray.length - 1; i > 0; i--) {
      const setLimit = 256 - (256 % (i + 1))
      let randval: number
      do {
        randval = window.crypto.getRandomValues(new Uint8Array(1))[0]
      } while (randval >= setLimit)
      const j = randval % (i + 1)
      ;[passwdArray[i], passwdArray[j]] = [passwdArray[j], passwdArray[i]]
    }

    return passwdArray.join('')
  }
}
