import { isNaN, isNumber, isObject, isString, isBoolean } from 'lodash-es'
import { Language } from 'vue3-gettext'

export interface PasswordPolicyRuleOptions {
  minLength?: number
  maxLength?: number
  characters?: string
}

export interface PasswordPolicyRuleExplained {
  code: string
  message: string
  helperMessage?: string
  format: (number | string)[]
  verified?: boolean
}

export interface PasswordPolicyRule {
  assert(options: PasswordPolicyRuleOptions, password: string): boolean

  explain(options: PasswordPolicyRuleOptions, verified?: boolean): PasswordPolicyRuleExplained

  missing(options: PasswordPolicyRuleOptions, password: string): PasswordPolicyRuleExplained

  validate(options?: PasswordPolicyRuleOptions): boolean
}

export class MustNotBeEmptyRule implements PasswordPolicyRule {
  protected $gettext

  constructor({ $gettext }: Language) {
    this.$gettext = $gettext
  }

  explain(options: PasswordPolicyRuleOptions, verified: boolean): PasswordPolicyRuleExplained {
    return {
      code: 'mustNotBeEmpty',
      message: this.$gettext('Must not be empty'),
      format: [],
      ...(isBoolean(verified) && { verified })
    }
  }

  assert(options: PasswordPolicyRuleOptions, password: string): boolean {
    return password.length > 0
  }

  validate(): boolean {
    return true
  }

  missing(options: PasswordPolicyRuleOptions, password: string): PasswordPolicyRuleExplained {
    return this.explain(options, this.assert(options, password))
  }
}

export class MustContainRule implements PasswordPolicyRule {
  protected $gettext

  constructor({ $gettext }: Language) {
    this.$gettext = $gettext
  }

  explain(options: PasswordPolicyRuleOptions, verified: boolean): PasswordPolicyRuleExplained {
    return {
      code: 'mustContain',
      helperMessage: this.$gettext(
        'Valid special characters: %{characters}',
        {
          characters: options.characters
        },
        true
      ),
      message: this.$gettext('%{param1}+ special characters'),
      format: [options.minLength],
      ...(isBoolean(verified) && { verified })
    }
  }

  assert(options: PasswordPolicyRuleOptions, password: string) {
    const charsCount = Array.from(password).filter((char) =>
      options.characters.includes(char)
    ).length

    return charsCount >= options.minLength
  }

  validate(options: PasswordPolicyRuleOptions): boolean {
    if (!isObject(options)) {
      throw new Error('options should be an object')
    }

    if (!isNumber(options.minLength) || isNaN(options.minLength)) {
      throw new Error('minLength should be a non-zero number')
    }

    if (!isString(options.characters)) {
      throw new Error('characters should be a character sequence')
    }

    return true
  }

  missing(options: PasswordPolicyRuleOptions, password: string) {
    return this.explain(options, this.assert(options, password))
  }
}

export class AtLeastBaseRule implements PasswordPolicyRule {
  protected $gettext

  constructor({ $gettext }: Language) {
    this.$gettext = $gettext
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  assert(options: PasswordPolicyRuleOptions, password: string): boolean {
    throw new Error('Method not implemented.')
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  explain(options: PasswordPolicyRuleOptions, verified?: boolean): PasswordPolicyRuleExplained {
    throw new Error('Method not implemented.')
  }

  validate(options: PasswordPolicyRuleOptions): boolean {
    if (!isObject(options)) {
      throw new Error('options should be an object')
    }

    if (!isNumber(options.minLength) || isNaN(options.minLength)) {
      throw new Error('minLength should be a non-zero number')
    }

    return true
  }

  missing(options: PasswordPolicyRuleOptions, password: string): PasswordPolicyRuleExplained {
    return this.explain(options, this.assert(options, password))
  }
}

export class AtLeastCharactersRule extends AtLeastBaseRule implements PasswordPolicyRule {
  constructor(args: Language) {
    super(args)
  }

  explain(options: PasswordPolicyRuleOptions, verified: boolean): PasswordPolicyRuleExplained {
    return {
      code: 'atLeastCharacters',
      message: this.$gettext('%{param1}+ letters'),
      format: [options.minLength],
      ...(isBoolean(verified) && { verified })
    }
  }

  assert(options: PasswordPolicyRuleOptions, password: string): boolean {
    return password.length >= options.minLength
  }
}

export class AtLeastUppercaseCharactersRule extends AtLeastBaseRule {
  constructor(args: Language) {
    super(args)
  }

  explain(options: PasswordPolicyRuleOptions, verified: boolean): PasswordPolicyRuleExplained {
    return {
      code: 'atLeastUppercaseCharacters',
      message: this.$gettext('%{param1}+ uppercase letters'),
      format: [options.minLength],
      ...(isBoolean(verified) && { verified })
    }
  }

  assert(options: PasswordPolicyRuleOptions, password: string): boolean {
    const uppercaseCount = (password || '').match(/[A-Z\xC0-\xD6\xD8-\xDE]/g)?.length
    return uppercaseCount >= options.minLength
  }
}

export class AtLeastLowercaseCharactersRule extends AtLeastBaseRule {
  constructor(args: Language) {
    super(args)
  }

  explain(options: PasswordPolicyRuleOptions, verified: boolean): PasswordPolicyRuleExplained {
    return {
      code: 'atLeastLowercaseCharacters',
      message: this.$gettext('%{param1}+ lowercase letters'),
      format: [options.minLength],
      ...(isBoolean(verified) && { verified })
    }
  }

  assert(options: PasswordPolicyRuleOptions, password: string): boolean {
    const lowercaseCount = (password || '').match(/[a-z\xDF-\xF6\xF8-\xFF]/g)?.length
    return lowercaseCount >= options.minLength
  }
}

export class AtLeastDigitsRule extends AtLeastBaseRule {
  constructor(args: Language) {
    super(args)
  }

  explain(options: PasswordPolicyRuleOptions, verified: boolean): PasswordPolicyRuleExplained {
    return {
      code: 'atLeastDigits',
      message: this.$gettext('%{param1}+ numbers'),
      format: [options.minLength],
      ...(isBoolean(verified) && { verified })
    }
  }

  assert(options: PasswordPolicyRuleOptions, password: string): boolean {
    const digitCount = (password || '').match(/\d/g)?.length
    return digitCount >= options.minLength
  }
}
